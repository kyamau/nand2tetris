package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"vm/codewriter"
	"vm/parser"
)

func Compile(r io.Reader, vmName string, bootstrap bool) string {
	p, err := parser.NewParser(r)
	if err != nil {
		log.Fatalf("Couldn't initialize parser : %v", err)
	}

	if !p.HasMoreCommands() {
		return ""
	}
	p.Advance()

	var b bytes.Buffer
	if bootstrap {
		b.WriteString(codewriter.Bootstrap())
	}
	for ; ; p.Advance() {
		cmdType := p.CommandType()
		switch cmdType {
		case parser.C_PUSH, parser.C_POP:
			segment := p.Arg1()
			arg2 := p.Arg2()
			index, err := strconv.Atoi(arg2)
			if err != nil {
				log.Fatalf("Argument of push must be integer : %v", index)
			}
			if segment != "static" {
				c := codewriter.WritePushPop(cmdType, segment, index)
				b.WriteString(c)
			} else {
				c := codewriter.WritePushPopStatic(cmdType, segment, index, vmName)
				b.WriteString(c)
			}
		case parser.C_LABEL:
			label := p.Arg1()
			c := codewriter.WriteLabel(label)
			b.WriteString(c)
		case parser.C_GOTO:
			label := p.Arg1()
			c := codewriter.WriteGoto(label)
			b.WriteString(c)
		case parser.C_IF:
			label := p.Arg1()
			c := codewriter.WriteIf(label)
			b.WriteString(c)
		case parser.C_FUNCTION:
			name := p.Arg1()
			arg2 := p.Arg2()
			nLocals, err := strconv.Atoi(arg2)
			if err != nil {
				log.Fatalf("2nd argument of function must be integer : %v", arg2)
			}
			c := codewriter.WriteFunction(name, nLocals)
			b.WriteString(c)
		case parser.C_RETURN:
			c := codewriter.WriteReturn()
			b.WriteString(c)
		case parser.C_CALL:
			name := p.Arg1()
			arg2 := p.Arg2()
			nArgs, err := strconv.Atoi(arg2)
			if err != nil {
				log.Fatalf("2nd argument of call must be integer : %v", arg2)
			}
			c := codewriter.WriteCall(name, nArgs)
			b.WriteString(c)
		case parser.C_ARITHMETIC:
			op, err := parser.ALOperatorFromString(p.Current())
			if err != nil {
				log.Fatalf("Invalid operator : %v", op)
			}
			c := codewriter.WriteArithmetic(op)
			b.WriteString(c)
		}
		if !p.HasMoreCommands() {
			return b.String()
		}
	}
}

func main() {
	var (
		bootstrap = flag.Bool("bootstrap", true, "Write bootstrap code or not")
	)
	flag.Parse()
	args := flag.Args()
	if flag.NArg() < 1 {
		exe, _ := os.Executable()
		fmt.Fprintf(os.Stderr, "Usage: %v <.vm dir> -bootstrap=<true/false>]\n", filepath.Base(exe))
		os.Exit(1)
	}

	vmDirPath, _ := filepath.Abs(args[0])
	var asm bytes.Buffer

	err := filepath.Walk(vmDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".vm" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err)
				return err
			}
			vmName := info.Name()[0 : len(info.Name())-3]
			asm.WriteString(Compile(bytes.NewReader(b), vmName, *bootstrap))
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Couldn't read .vm in the directory : %v", err)
	}

	asmPath := filepath.Base(vmDirPath) + ".asm"
	err = ioutil.WriteFile(asmPath, []byte(asm.String()), 644)
	if err != nil {
		log.Fatalf("Couldn't write .asm : %v, %v", asmPath, err)
	}
}
