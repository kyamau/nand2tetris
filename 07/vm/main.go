package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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
		fmt.Errorf("Usage: %v <.vm>", exe)
	}

	vmPath := os.Args[1]
	f, err := os.Open(vmPath)
	if err != nil {
		log.Fatalf("Couldn't open file : %v", vmPath)
	}
	c := Compile(f)
	fmt.Println(c)
}
