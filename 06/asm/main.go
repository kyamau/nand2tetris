package main

import (
	"asm/code"
	"asm/parser"
	"asm/symbol_table"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
		exe, _ := os.Executable()
		fmt.Fprintf(os.Stderr, "Usage: %v <.asm>\n", filepath.Base(exe))
		os.Exit(1)
	}
	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Couldn't open .asm : %v", err)
	}

	obj := Compile(f)
	hackPath := filepath.Base(path[:len(path)-len(filepath.Ext(path))]) + ".hack"

	hackb := ""
	for i := 0; i < len(obj); i++ {
		hackb += fmt.Sprintf("%016b\r\n", obj[i])
	}
	err = ioutil.WriteFile(hackPath, []byte(hackb), 644)
	if err != nil {
		log.Fatalf("Couldn't write .hack : %v, %v", hackPath, err)
	}
}
