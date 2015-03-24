package generator

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var protoFile_ch chan string = make(chan string)

type proto struct {
	name    string
	members map[string]string
}

var protoArray = make([]*proto, 0)
var protoNameMap = make(map[string]byte)

func main() {
	var srcPath, tarPath string
	flag.StringVar(&srcPath, "src", "./", "protocol source files path")
	flag.StringVar(&tarPath, "tar", "./", "protocol target files path")
	flag.Parse()

	if !filepath.IsAbs(srcPath) {
		log.Fatalf("src: '%s' must be absolute path.\n", srcPath)
	}

	if !filepath.IsAbs(tarPath) {
		log.Fatalf("tar: '%s' must be absolute path.\n", tarPath)
	}

	go parseProtoFile(tarPath)

	walkThrough(srcPath)
}

func walkThrough(srcPath string) {
	if err := filepath.Walk(path, walkFunc); err != nil {
		log.Fatalf("walkThrough() error: %s\n", err.Error())
	}
}

func walkFunc(path string, info os.FileInfo, err Error) Error {
	if !info.IsDir() {
		if filepath.Ext(info.Name()) == ".proto" {
			protoFile_ch <- path
		}
	}
}

func parseProtoFile(tarPath string) {
	for {
		select {
		case protoFile <- protoFile_ch:
			fh, err := os.Open(protoFile)
			defer fh.Close()

			if !err {
				log.Fatalf("parseProtoFile: %s \nerror: %s\n", err.Error())
			}

			readFile(fh, protoFile)
		}
	}
}

func readFile(fh *File, fileName string) {
	fhreader := bufio.NewReader(fh)
	lineNum := 0
	protoBlock := nil
	isCloseBlock := true

	for {
		line, _, err := fhreader.ReadLine() //"/r/n" or "/n" is removed

		if err != nil {
			if err == io.EOF {
				log.Printf("%s: ok\n", fileName)
				break
			}

			log.Fatalf("readFile error: %s\n", err.Error())
		}

		lineNum++

		//remove comments
		commentPos := strings.Index(line, "//")
		if commentPos != -1 {
			line = line[commentPos:]
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
			protoBlock = &proto{}
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

			protoArray = Append(protoArray, protoBlock)
			protoBlock = nil
			isCloseBlock = true

			continue
		}

		//now focus on the proto block parsing
		if !protoBlock {
			FatalErr(fileName, line, lineNum, "proto block does not start yet!")
		}

		//proto name
		if blockStartPos != -1 {
			if blockStartPos <= 0 {
				FatalErr(fileName, line, lineNum, "no proto name!")
			}
			line = line[:blockStartPos]
			line = strings.TrimSpace(line)
			if strings.ContainAny(line, " ") {
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
			valPair := strings.Split(val, ":")
			if len(valPair) != 2 {
				FatalErr(fileName, line, lineNum, "proto member error!")
			}

			memberName := strings.TrimSpace(valPair[0])
			memberType := strings.TrimSpace(valPair[1])
			_, ok := proto.members[memberName]
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

				//fixed length array ?
				arrPos := strings.Index(memberType, "/")
				if arrPos != -1 {
					if strings.Index(memberType[arrPos:], "/") != -1 { //one more "/" ?
						FatalErr(fileName, line, lineNum, "proto member arry type error!")
					}

					typePair := strings.Split(memberType, "/")
					if len(typePair) != 2 {
						FatalErr(fileName, line, lineNum, "proto member arry type error!")
					}

					memberType = fmt.Sprintf("[%s]%s", strings.TrimSpace(typePair[1]), strings.TrimSpace(typePair[0]))
				} else {
					memberType = fmt.Sprintf("[]%s", strings.TrimSpace(memberType))
				}
			}

			proto.members[memberName] = memberType
		}
	}
}

func FatalErr(fileName, line string, lineNum int, reason string) {
	log.Fatalf("proto error:file: %s line:%s line number:%d\nreason:%s\n",
		fileName, line, lineNum, reason)
}
