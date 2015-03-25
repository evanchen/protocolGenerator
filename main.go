package main

import (
	"flag"
	_ "fmt"
	"proto"
)

func main() {
	var srcPath, tarPath string
	flag.StringVar(&srcPath, "src", "./", "protocol source files path")
	flag.StringVar(&tarPath, "tar", "./", "protocol target files path")
	flag.Parse()

	//if !filepath.IsAbs(srcPath) {
	//	log.Fatalf("src: '%s' must be absolute path.\n", srcPath)
	//}

	//if !filepath.IsAbs(tarPath) {
	//	log.Fatalf("tar: '%s' must be absolute path.\n", tarPath)
	//}
	proto.Start(srcPath, tarPath)
}
