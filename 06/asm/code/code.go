package code

import (
	"log"
	"strconv"
)

func C(dest string, comp string, jump string) uint16 {
	d := Dest(dest)
	c := Comp(comp)
	j := Jump(jump)
	return 0b111<<13 | c<<6 | d<<3 | j
}

func A(symbol string) uint16 {
	// string -> uint16
	ui64, err := strconv.ParseUint(symbol, 10, 16)
	s := uint16(ui64)
	if err != nil {
		log.Fatalf("Can't parse A symbol %v : %v", symbol, err)
	}
	return 0b0<<15 | s
}

func Dest(dest string) uint16 {
	switch dest {
	case "null":
		return 0b000
	case "M":
		return 0b001
	case "D":
		return 0b010
	case "MD":
		return 0b011
	case "A":
		return 0b100
	case "AM":
		return 0b101
	case "AD":
		return 0b110
	case "AMD":
		return 0b111
	}
	log.Fatalf("Illegal dest : %v", dest)
	return 0
}

func Comp(comp string) uint16 {
	switch comp {
	// a=0
	case "0":
		return 0b0101010
	case "1":
		return 0b0111111
	case "-1":
		return 0b0111010
	case "D":
		return 0b0001100
	case "A":
		return 0b0110000
	case "!D":
		return 0b0001100
	case "!A":
		return 0b0110001
	case "-D":
		return 0b0001111
	case "-A":
		return 0b0110011
	case "D+1":
		return 0b0011111
	case "A+1":
		return 0b0110111
	case "D-1":
		return 0b0001110
	case "A-1":
		return 0b0110010
	case "D+A":
		return 0b0000010
	case "D-A":
		return 0b0010011
	case "A-D":
		return 0b0000111
	case "D&A":
		return 0b0000000
	case "D|A":
		return 0b0010101
	// a=1
	case "M":
		return 0b1110000
	case "!M":
		return 0b1110001
	case "-M":
		return 0b1110011
	case "M+1":
		return 0b1110111
	case "M-1":
		return 0b1110010
	case "D+M":
		return 0b1000010
	case "D-M":
		return 0b1010011
	case "M-D":
		return 0b1000111
	case "D&M":
		return 0b1000000
	case "D|M":
		return 0b1010101
	}
	log.Fatalf("Illegal comp : %v", comp)
	return 0
}

func Jump(jump string) uint16 {
	switch jump {
	case "null":
		return 0b000
	case "JGT":
		return 0b001
	case "JEQ":
		return 0b010
	case "JGE":
		return 0b011
	case "JLT":
		return 0b100
	case "JNE":
		return 0b101
	case "JLE":
		return 0b110
	case "JMP":
		return 0b111
	}
	log.Fatalf("Illegal jump : %v", jump)
	return 0
}
