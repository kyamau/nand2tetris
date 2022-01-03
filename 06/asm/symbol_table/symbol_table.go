package symbol_table

import (
	"fmt"
	"log"
)

const baseAddress uint16 = 1024

type SymbolTable struct {
	table  map[string]uint16
	offset uint16
}

func (t *SymbolTable) addSystemSymbol(newSymbol string, address uint16) {
	t.table[newSymbol] = address
}

func (t *SymbolTable) AddVariable(newSymbol string) uint16 {
	address := baseAddress + t.offset
	t.table[newSymbol] = address
	t.offset++
	return address
}

func (t *SymbolTable) AddLable(newSymbol string, romAddress uint16) {
	t.table[newSymbol] = romAddress
}

func (t *SymbolTable) ExistVariable(symbol string) bool {
	_, ok := t.table[symbol]
	return ok
}

func (t *SymbolTable) GetAddress(symbol string) uint16 {
	ret, ok := t.table[symbol]
	if !ok {
		log.Fatalf("No such variable in symbol table, %v", symbol)
	}
	return ret
}

func NewSymbolTable() *SymbolTable {
	t := SymbolTable{table: map[string]uint16{}}
	t.addSystemSymbol("SP", 0)
	t.addSystemSymbol("LCL", 1)
	t.addSystemSymbol("ARG", 2)
	t.addSystemSymbol("THIS", 3)
	t.addSystemSymbol("THAT", 4)
	for i := uint16(0); i < 15; i++ {
		t.addSystemSymbol(fmt.Sprintf("R%d", i), i)
	}
	t.addSystemSymbol("SCREEN", 16384)
	t.addSystemSymbol("KBD", 24576)
	return &t
}
