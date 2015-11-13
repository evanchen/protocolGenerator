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

type member struct {
	mname string
	mtype string
}

type proto struct {
	name    string
	members []*member
}

var protoFile_ch = make(chan string)
var finish_ch = make(chan bool)
var protoNameMap = make(map[string]byte)
var logger = log.New(os.Stderr, "", log.Lshortfile)

func Generate(srcPath, tarPath string) {
	go parseProtoFile(tarPath)
	walkThrough(srcPath)

	<-finish_ch
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
				finish_ch <- true
				break
			}
			protoFile = filepath.ToSlash(protoFile)
			fmt.Printf("parsing file: %s\n\n", protoFile)

			fh, err := os.Open(protoFile)
			if err != nil {
				logger.Fatalf("parseProtoFile: %s \nerror: %s\n", protoFile, err.Error())
			}
			protoArray := readFile(fh, protoFile)
			fh.Close()

			//start writing protocol file
			//get new file name
			tarFileName := filepath.Base(protoFile)
			tarFileName = strings.Split(tarFileName, ".")[0]
			tarFileName = fmt.Sprintf("%sProto.go", tarFileName)
			tarFileName = filepath.Join(tarPath, tarFileName)
			wh, err2 := os.Create(tarFileName) //rw truncate
			if err2 != nil {
				logger.Fatalf("failed to create protocol file: %s\n", err2.Error())
			}
			whwriter := bufio.NewWriter(wh)

			//header
			printHeader(whwriter)
			whwriter.Flush()
			for _, v := range protoArray {
				fmt.Fprintf(whwriter, "\n//===========================protocol %s===========================\n", v.name)
				v.printBody(whwriter)
				whwriter.Flush()
			}
			wh.Close()
		}
	}
}

func readFile(fh *os.File, fileName string) []*proto {
	fhreader := bufio.NewReader(fh)
	lineNum := 0
	var protoBlock *proto
	var memberConflict map[string]byte
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
				members: make([]*member, 0),
			}
			memberConflict = make(map[string]byte)
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
			memberConflict = nil
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

			line = strings.Title(line)
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

			//does member name exist
			memberName := strings.TrimSpace(valPair[0])
			memberName = strings.Title(memberName)
			_, ok := memberConflict[memberName]
			if ok {
				FatalErr(fileName, line, lineNum, "proto member name redefined!")
			}

			//now the member type
			memberType := strings.TrimSpace(valPair[1])
			if strings.HasPrefix(memberType, "[") { //is it an array
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
				_, isTypeProto := protoNameMap[strings.Title(memberType)]
				if isTypeProto {
					memberType = strings.Title(memberType)
				}
				memberType = fmt.Sprintf("[]%s", strings.TrimSpace(memberType))
			} else {
				_, isTypeProto := protoNameMap[strings.Title(memberType)]
				if isTypeProto {
					memberType = strings.Title(memberType)
				}
			}

			m := &member{
				mname: memberName,
				mtype: memberType,
			}
			protoBlock.members = append(protoBlock.members, m)
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

func printHeader(w io.Writer) {
	fmt.Fprintln(w, "//This file is automatically created by protocolGenerator.")
	fmt.Fprintln(w, "//https://github.com/evanchen/protocolGenerator")
	fmt.Fprintln(w, "//Any manual changes are not suggested.\n")

	fmt.Fprintln(w, "package protocol\n")
}

func (this *proto) printBody(w io.Writer) {
	//structure
	fmt.Fprintf(w, "type %s struct {\n", this.name)
	for _, v := range this.members {
		fmt.Fprintf(w, "	%s %s\n", v.mname, v.mtype)
	}
	fmt.Fprintln(w, "}\n")
	//func CreateXXX()
	fmt.Fprintf(w, "func Create%s() *%s {\n	obj := &%s{}\n	return obj\n}\n\n", this.name, this.name, this.name)
	//func Marshal()
	fmt.Fprintf(w, "func (this *%s) Marshal() ([]byte) {\n", this.name)
	fmt.Fprintln(w, "	buf := make([]byte,0,16)")
	for _, v := range this.members {
		_, ok := protoNameMap[v.mtype]
		if ok {
			fmt.Fprintf(w, "	buf = append(buf,this.%s.Marshal()...)\n", v.mname)
		} else {
			arrPos := strings.Index(v.mtype, "]")
			if arrPos > 0 { //is it an array
				arrType := v.mtype[arrPos+1:]
				fmt.Fprintf(w, "	buf = append(buf,Encode_array_%s(this.%s)...)\n", arrType, v.mname)
			} else {
				fmt.Fprintf(w, "	buf = append(buf,Encode_%s(this.%s)...)\n", v.mtype, v.mname)
			}
		}
	}
	fmt.Fprintln(w, "	return buf")
	fmt.Fprintln(w, "}\n")

	//func Unmarshal()
	fmt.Fprintf(w, "func (this *%s) Unmarshal(Data []byte) ([]byte) {\n", this.name)
	for _, v := range this.members {
		_, ok := protoNameMap[v.mtype]
		if ok {
			fmt.Fprintf(w, "	Data = this.%s.Unmarshal(Data)\n", v.mname)
		} else {
			arrPos := strings.Index(v.mtype, "]")
			if arrPos > 0 { //is it an array
				arrType := v.mtype[arrPos+1:]
				fmt.Fprintf(w, "	this.%s,Data = Decode_array_%s(Data)\n", v.mname, arrType)
			} else {
				fmt.Fprintf(w, "	this.%s,Data = Decode_%s(Data)\n", v.mname, v.mtype)
			}
		}
	}
	fmt.Fprintln(w, "	return Data")
	fmt.Fprintln(w, "}\n")

	//protocol array
	for _, v := range this.members {
		arrPos := strings.Index(v.mtype, "]")
		if arrPos > 0 { //is it an array
			arrType := v.mtype[arrPos+1:]
			_, ok := protoNameMap[arrType]
			if ok {
				//encode array func
				fmt.Fprintf(w, "func Encode_array_%s(%s []%s) ([]byte) {\n", arrType, v.mname, arrType)
				fmt.Fprintln(w, "	buf := make([]byte,0,16)")
				fmt.Fprintf(w, "	size := uint16(len(%s))\n", v.mname)
				fmt.Fprintln(w, "	buf = append(buf,Encode_uint16(size)...)")
				fmt.Fprintf(w, "	for _,obj := range %s {\n", v.mname)
				fmt.Fprintln(w, "		buf = append(buf,obj.Marshal()...)")
				fmt.Fprintln(w, "	}")
				fmt.Fprintln(w, "	return buf")
				fmt.Fprintln(w, "}\n")

				//decode array func
				fmt.Fprintf(w, "func Decode_array_%s(Data []byte) ([]%s,[]byte) {\n", arrType, arrType)
				fmt.Fprintln(w, "	var size uint16")
				fmt.Fprintln(w, "	size,Data = Decode_uint16(Data)")
				fmt.Fprintf(w, "	%s := make([]%s,0,size)\n", v.mname, arrType)
				fmt.Fprintf(w, "	var obj *%s\n", arrType)
				fmt.Fprintln(w, "	for i := uint16(0); i < size; i++ {")
				fmt.Fprintf(w, "		obj = &%s{}\n", arrType)
				fmt.Fprintln(w, "		Data = obj.Unmarshal(Data)")
				fmt.Fprintf(w, "		%s = append(%s,*obj)\n", v.mname, v.mname)
				fmt.Fprintln(w, "	}")
				fmt.Fprintf(w, "	return %s,Data\n", v.mname)
				fmt.Fprintln(w, "}\n")
			}
		}
	}

}
