package protocol

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var protoFile_ch = make(chan string)
var finishch = make(chan bool)

type proto struct {
	name    string
	members map[string]string
}

var protoNameMap = make(map[string]byte)
var logger = log.New(os.Stderr, "", log.Lshortfile)

func Generate(srcPath, tarPath string) {
	go parseProtoFile(tarPath)
	walkThrough(srcPath)

	<-finishch
}

func walkThrough(srcPath string) {
	defer close(protoFile_ch)

	if err := filepath.Walk(srcPath, walkFunc); err != nil {
		logger.Fatalf("walkThrough() error: %s\n", err.Error())
	}
}

func walkFunc(path string, info os.FileInfo, err error) error {
	if !info.IsDir() {
		if filepath.Ext(info.Name()) == ".proto" {
			protoFile_ch <- path
		}
	}
	return err
}

func parseProtoFile(tarPath string) {

	for {
		select {
		case protoFile, ok := <-protoFile_ch:
			if !ok { //finish reading
				finishch <- true
				break
			}
			fmt.Printf("parsing file: %s....\n\n", protoFile)

			fh, err := os.Open(protoFile)
			defer fh.Close()

			if err != nil {
				logger.Fatalf("parseProtoFile: %s \nerror: %s\n", err.Error())
			}

			protoArray := readFile(fh, protoFile)

			//test print
			for _, v := range protoArray {
				v.print()
				fmt.Println("==========================")
			}
		}
	}
}

func readFile(fh *os.File, fileName string) []*proto {
	fhreader := bufio.NewReader(fh)
	lineNum := 0
	var protoBlock *proto
	isCloseBlock := true
	protoArray := make([]*proto, 0)

	for {
		rline, _, err := fhreader.ReadLine() //"/r/n" or "/n" is removed
		line := string(rline)

		if err != nil {
			if err == io.EOF {
				break
			}

			logger.Fatalf("readFile error: %s\n", err.Error())
		}

		lineNum++

		//remove comments
		commentPos := strings.Index(line, "//")
		if commentPos != -1 {
			line = line[:commentPos]
		}
		line = strings.TrimSpace(line)

		if len(line) <= 0 {
			continue
		}

		//start parsing a new proto block
		blockStartPos := strings.Index(line, "{")
		if blockStartPos != -1 {
			if !isCloseBlock || protoBlock != nil {
				FatalErr(fileName, line, lineNum, "proto block starts twice!")
			}
			protoBlock = &proto{
				name:    "",
				members: make(map[string]string),
			}
			isCloseBlock = false
		}

		//end parsing a proto block
		blockEndPos := strings.Index(line, "}")
		if blockEndPos != -1 {
			if isCloseBlock || protoBlock == nil {
				FatalErr(fileName, line, lineNum, "proto block closes twice!")
			}

			//proto ending
			if blockEndPos != 0 {
				FatalErr(fileName, line, lineNum, "proto ending should be just \"}\"!")
			}

			protoArray = append(protoArray, protoBlock)
			protoBlock = nil
			isCloseBlock = true

			continue
		}

		//now focus on the proto block parsing
		if protoBlock == nil {
			FatalErr(fileName, line, lineNum, "proto block does not start yet!")
		}

		//proto name
		if blockStartPos != -1 {
			if blockStartPos <= 0 {
				FatalErr(fileName, line, lineNum, "no proto name!")
			}
			line = line[:blockStartPos]
			line = strings.TrimSpace(line)
			if strings.ContainsAny(line, " ") {
				FatalErr(fileName, line, lineNum, "proto name should not contain any spaces!")
			}

			_, ok := protoNameMap[line]
			if ok {
				FatalErr(fileName, line, lineNum, "proto name redefined!")
			}

			protoBlock.name = line
			protoNameMap[line] = 1
			continue
		}

		//proto members
		memberList := strings.Split(line, ",")
		if len(memberList) <= 0 {
			FatalErr(fileName, line, lineNum, "no proto member!")
		}

		for _, val := range memberList {
			val = strings.TrimSpace(val)
			if len(val) <= 0 {
				continue
			}
			valPair := strings.Split(val, ":")
			if len(valPair) != 2 {
				FatalErr(fileName, line, lineNum, "proto member error!")
			}

			memberName := strings.TrimSpace(valPair[0])
			memberType := strings.TrimSpace(valPair[1])
			_, ok := protoBlock.members[memberName]
			if ok {
				FatalErr(fileName, line, lineNum, "proto member name redefined!")
			}

			//now the member type
			if strings.HasPrefix(memberType, "[") {
				if !strings.HasSuffix(memberType, "]") {
					FatalErr(fileName, line, lineNum, "proto member arry type error!")
				} else if strings.ContainsAny(memberType[:len(memberType)-1], "]") { //one more "]" ?
					FatalErr(fileName, line, lineNum, "proto member arry type error!")
				}

				if strings.ContainsAny(memberType[1:], "[") { //one more "[" ?
					FatalErr(fileName, line, lineNum, "proto member arry type error!")
				}
				memberType = strings.TrimLeft(memberType, "[")
				memberType = strings.TrimRight(memberType, "]")
				memberType = fmt.Sprintf("[]%s", strings.TrimSpace(memberType))
			}

			protoBlock.members[memberName] = memberType
		}
	}

	return protoArray
}

func FatalErr(fileName, line string, lineNum int, reason string) {
	l := fmt.Sprintf("parsing proto file: %s:%d\nerror str: %s\nreason: %s\n",
		fileName, lineNum, line, reason)
	logger.Output(2, l)
	os.Exit(1)
}

func (this *proto) print() {
	fmt.Printf("type %s struct {\n", this.name)
	for k, v := range this.members {
		fmt.Printf("	%s %s\n", k, v)
	}
	fmt.Println("}")
}
