package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"compiler/parser"
	"compiler/tokenizer"
)

var (
	tokenizeOnly = flag.Bool("tokenize", false, "Tokenization only mode")
	parse        = flag.Bool("parse", false, "Tokenization + Parsing mode")
)

func compile(srcPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("Failed to open .jack: %v", err)
	}

	// Tokenize
	tokenizer, err := tokenizer.NewTokenizer(f)
	if err != nil {
		log.Fatalf("Failed to initialize tokenizer: %v", err)
	}
	err = tokenizer.Tokenize()
	if err != nil {
		return fmt.Errorf("Failed to tokenize: src=%v: %v", srcPath, err)
	}

	tokenXML := tokenizer.XML()
	base := filepath.Base(srcPath)
	tokenFilename := base[:strings.LastIndex(base, ".")] + "T.xml.out"
	tokenDstPath := filepath.Join(filepath.Dir(srcPath), tokenFilename)
	if os.Getenv("LOGLEVEL") == "debug" {
		log.Printf("Tokenize output path=%v\n", tokenDstPath)
	}
	err = ioutil.WriteFile(tokenDstPath, []byte(tokenXML), 0666)
	if err != nil {
		return err
	} else if *tokenizeOnly {
		return nil
	}

	// Parse
	parser := parser.NewParser(*tokenizer)
	err = parser.Parse()
	if err != nil {
		return fmt.Errorf("Failed to parse: src=%v: %v", srcPath, err)
	}

	treeXML := parser.XML()
	treeFileName := base[:strings.LastIndex(base, ".")] + ".xml.out"
	treeDstPath := filepath.Join(filepath.Dir(srcPath), treeFileName)
	if os.Getenv("LOGLEVEL") == "debug" {
		log.Printf("Parse output path=%v\n", treeDstPath)
	}

	err = ioutil.WriteFile(treeDstPath, []byte(treeXML), 0666)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	args := flag.Args()
	if flag.NArg() < 1 {
		exe, _ := os.Executable()
		fmt.Fprintf(os.Stderr, "Usage: %v <.jack/.jack dir> [-tokenize | -parse]\n", filepath.Base(exe))
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
			log.Fatalf("Failed to compile %v: %v", srcPath, err)
		}
	}
}
