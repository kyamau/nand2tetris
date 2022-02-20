package main

import (
	"bytes"
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

func Compile(r io.Reader) string {
	p, err := parser.NewParser(r)
	if err != nil {
		log.Fatalf("Couldn't initialize parser : %v", err)
	}

	if !p.HasMoreCommands() {
		return ""
	}
	p.Advance()

	var b bytes.Buffer
	for ; ; p.Advance() {
		cmdType := p.CommandType()
		switch cmdType {
		case parser.C_PUSH, parser.C_POP:
			arg1 := p.Arg1()
			arg2s := p.Arg2()
			arg2, err := strconv.Atoi(arg2s)
			if err != nil {
				log.Fatalf("Argument of push must be integer : %v", arg2)
			}
			c := codewriter.WritePushPop(cmdType, arg1, arg2)
			b.WriteString(c)
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
	if len(os.Args) != 2 {
		exe, _ := os.Executable()
		fmt.Fprintf(os.Stderr, "Usage: %v <.vm dir>\n", filepath.Base(exe))
		os.Exit(1)
	}

	vmDirPath, _ := filepath.Abs(os.Args[1])
	var commands bytes.Buffer

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
			commands.Write(b)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Couldn't read .vm in the directory : %v", err)
	}

	c := Compile(bytes.NewReader(commands.Bytes()))
	asmPath := filepath.Base(vmDirPath) + ".asm"
	err = ioutil.WriteFile(asmPath, []byte(c), 644)
	if err != nil {
		log.Fatalf("Couldn't write .asm : %v, %v", asmPath, err)
	}
}
