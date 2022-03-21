package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"compiler/tokenizer"
)

const (
	SEP = "\r\n"
)

func compile(srcPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("Failed to open .jack: %v", err)
	}
	tokenizer, err := tokenizer.NewTokenizer(f)
	if err != nil {
		log.Fatalf("Failed to initialize tokenizer: %v", err)
	}
	err = tokenizer.Tokenize()
	if err != nil {
		return err
	}
	xml := tokenizer.XML()
	fmt.Print(xml)
	return err
}

func main() {
	flag.Parse()

	args := flag.Args()
	if flag.NArg() < 1 {
		exe, _ := os.Executable()
		fmt.Fprintf(os.Stderr, "Usage: %v <.jack/.jack dir>]\n", filepath.Base(exe))
		os.Exit(1)
	}

	srcPath, _ := filepath.Abs(args[0])
	finfo, err := os.Stat(srcPath)
	if err != nil {
		log.Fatalf("Couldn't read %v", srcPath)
	}
	if finfo.IsDir() {
		err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(info.Name()) == ".jack" {
				err = compile(path)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to compile: %v", err)
		}
	} else if filepath.Ext(srcPath) == ".jack" {
		err = compile(srcPath)
		if err != nil {
			log.Fatalf("Failed to compile: %v", err)
		}
	}
}
