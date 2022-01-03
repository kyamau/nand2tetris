package main

import (
	"asm/code"
	"asm/parser"
	"asm/symbol_table"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var symbolTable *symbol_table.SymbolTable = symbol_table.NewSymbolTable()

func Compile(r io.Reader) []uint16 {
	p, err := parser.NewParser(r)
	if err != nil {
		log.Fatalf("Couldn't initialize parser : %v", err)
	}

	if !p.HasMoreCommands() {
		// No commands
		return make([]uint16, 0)
	}
	p.Advance()

	// First path
	//log.Println("First path")
	romAddress := uint16(0)
	for ; ; p.Advance() {
		cmdType := p.CommandType()
		//log.Printf("Command=%v", p.Current())
		switch cmdType {
		case parser.L_COMMAND:
			label := p.Label()
			symbolTable.AddLable(label, romAddress)
		}
		if !p.HasMoreCommands() {
			break
		}
		romAddress++
	}

	p.ResetCurrent()

	// Second path
	//log.Println("Second path")
	obj := make([]uint16, 0)
	for ; ; p.Advance() {
		cmdType := p.CommandType()
		//log.Printf("Command=%v", p.Current())

		switch cmdType {
		case parser.A_COMMAND:
			symbol := p.Symbol()
			// Variable
			if _, err := strconv.Atoi(symbol); err != nil {
				if !symbolTable.ExistVariable(symbol) {
					symbolTable.AddVariable(symbol)
				}
				address := symbolTable.GetAddress(symbol)
				p.RewriteSymbolToAddress(address)
			}
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
