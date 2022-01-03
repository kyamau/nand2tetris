package main

import (
	"asm/code"
	"asm/parser"
	"fmt"
	"io"
	"log"
	"os"
)

func Compile(r io.Reader) []uint16 {
	p, err := parser.NewParser(r)
	if err != nil {
		log.Fatalf("Couldn't initialize parser : %v", err)
	}

	obj := make([]uint16, 0)
	if !p.HasMoreCommands() {
		return obj
	}
	for p.Advance(); ; p.Advance() {
		cmdType := p.CommandType()
		cmd := p.Current()
		log.Printf("current command=%v", cmd)

		switch cmdType {
		case parser.A_COMMAND:
			obj = append(obj, code.A(p.Symbol()))
		case parser.C_COMMAND:
			obj = append(obj, code.C(p.Dest(), p.Comp(), p.Jump()))
		case parser.L_COMMAND:
		}
		if !p.HasMoreCommands() {
			break
		}
	}
	return obj
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: Assembler <.asm>")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Couldn't open .asm : %v", err)
	}
	Compile(f)
}
