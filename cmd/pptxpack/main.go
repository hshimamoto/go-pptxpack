// MIT License Copyright (C) 2022 Hiroshi Shimamoto
package main

import (
	"fmt"
	"os"

	"github.com/hshimamoto/go-pptxpack"
)

func usage() {
	fmt.Printf("pptxpack <dir> <pptx file>\n")
	fmt.Printf("pptxpack -d <dir> <pptx file>\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) <= 2 {
		usage()
		return
	}
	if os.Args[1] == "-d" {
		if len(os.Args) != 4 {
			usage()
			return
		}
		p, err := pptxpack.New(os.Args[2])
		if err != nil {
			fmt.Printf("New: %v\n", err)
			return
		}
		err = p.Unpack(os.Args[3])
		if err != nil {
			fmt.Printf("Unpack: %v\n", err)
			return
		}
		return
	}
	p, err := pptxpack.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Open: %v\n", err)
		return
	}
	err = p.Pack(os.Args[2])
	if err != nil {
		fmt.Printf("Pack: %v\n", err)
		return
	}
}
