package codewriter

import (
	"bytes"
	"fmt"
	"log"

	"vm/parser"
)

var labelIndex int

func pushD(code *bytes.Buffer) {
	code.WriteString("@SP\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("M=D\r\n")
	code.WriteString("@SP\r\n")
	code.WriteString("M=M+1\r\n")
}

func popToD(code *bytes.Buffer) {
	code.WriteString("@SP\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
}

func setTrueOrFalseToD(code *bytes.Buffer, comp string, jump string) {
	code.WriteString(fmt.Sprintf("@TRUE%v // Set true or false to D\r\n", labelIndex))
	code.WriteString(fmt.Sprintf("%v;%v\r\n", comp, jump))

	// False: set 0 to D
	code.WriteString("@0\r\n")
	code.WriteString("D=A\r\n")
	code.WriteString(fmt.Sprintf("@TFEND%v\r\n", labelIndex))
	code.WriteString("0;JMP\r\n")

	// True: set -1 to D
	code.WriteString(fmt.Sprintf("(TRUE%v)\r\n", labelIndex))
	code.WriteString("@1\r\n")
	code.WriteString("D=-A\r\n")

	code.WriteString(fmt.Sprintf("(TFEND%v)\r\n", labelIndex))
	labelIndex++
}

func WriteArithmetic(op parser.ALOperator) string {
	var code bytes.Buffer

	// Pop operand y from the stack to R13
	popToD(&code)
	code.WriteString("@13\r\n")
	code.WriteString("M=D\r\n")

	// If op is a binary operator, Pop operand x from the stack to R14
	switch op {
	case parser.ADD, parser.SUB, parser.EQ, parser.GT, parser.LT, parser.AND, parser.OR:
		popToD(&code)
		code.WriteString("@14\r\n")
		code.WriteString("M=D\r\n")
	}

	// Caluculate and load the result to D
	switch op {
	case parser.ADD:
		code.WriteString("@14 // add\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D+M\r\n")
	case parser.SUB:
		code.WriteString("@14 // sub\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D-M\r\n")
	case parser.NEG:
		code.WriteString("@13 // neg\r\n")
		code.WriteString("D=-M\r\n")
	case parser.EQ:
		code.WriteString("@14 // eq\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D-M\r\n")
		setTrueOrFalseToD(&code, "D", "JEQ") // x-y==0
	case parser.GT:
		code.WriteString("@14 // gt\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D-M\r\n")
		setTrueOrFalseToD(&code, "D", "JGT") // x-y>0
	case parser.LT:
		code.WriteString("@14 // lt\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D-M\r\n")
		setTrueOrFalseToD(&code, "D", "JLT") // x-y<0
	case parser.AND:
		code.WriteString("@14 // and\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D&M\r\n")
	case parser.OR:
		code.WriteString("@14 // or\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D|M\r\n")
	case parser.NOT:
		code.WriteString("@13 // not\r\n")
		code.WriteString("D=!M\r\n")
	}
	// Push D to the stack
	pushD(&code)

	return code.String()
}

func segment2Symbol(segment string) string {
	switch segment {
	case "local":
		return "LCL"
	case "argument":
		return "ARG"
	case "this":
		return "THIS"
	case "that":
		return "THAT"
	case "pointer":
		return "3"
	case "temp":
		return "5"
	case "static":
		return "16"
	}
	log.Fatalf("Invalid segment : %v", segment)
	return ""
}

func setAddressToD(code *bytes.Buffer, segment string, index int) {
	segsym := segment2Symbol(segment)
	code.WriteString(fmt.Sprintf("@%v\r\n", index))
	code.WriteString("D=A\r\n")
	code.WriteString(fmt.Sprintf("@%v\r\n", segsym))
	switch segment {
	case "local", "argument", "this", "that":
		code.WriteString("D=D+M\r\n")
	case "pointer", "temp":
		code.WriteString("D=D+A\r\n")
	}
}

func WritePushPop(cmdType parser.CommandType, segment string, index int) string {
	var code bytes.Buffer
	switch cmdType {
	case parser.C_POP:

		// Pop to R13
		popToD(&code)
		code.WriteString("@13\r\n") // Load poped item to R13
		code.WriteString("M=D\r\n")

		// Calculate segment + i and set it to R14
		setAddressToD(&code, segment, index)
		code.WriteString("@14\r\n") // Load poped item to R13
		code.WriteString("M=D\r\n")

		// Write the poped value to the destination
		code.WriteString("@13\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@14\r\n")
		code.WriteString("A=M\r\n")
		code.WriteString("M=D\r\n")

	case parser.C_PUSH:
		// Load to D
		switch segment {
		case "constant":
			code.WriteString(fmt.Sprintf("@%v\r\n", index))
			code.WriteString("D=A\r\n")
		default:
			setAddressToD(&code, segment, index)
			code.WriteString("A=D\r\n")
			code.WriteString("D=M\r\n")
		}

		// Push
		pushD(&code)
	}
	return code.String()

}
