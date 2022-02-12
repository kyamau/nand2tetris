package codewriter

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"vm/parser"
)

var labelIndex int

func pushD(code *bytes.Buffer, comment ...string) {
	str := strings.Join(comment, ",")
	code.WriteString(fmt.Sprintf("@SP // %vPush the value at the address in D\r\n", str))
	code.WriteString("A=M\r\n")
	code.WriteString("M=D\r\n")
	code.WriteString("@SP\r\n")
	code.WriteString("M=M+1\r\n")
}

func popToD(code *bytes.Buffer, comment ...string) {
	str := strings.Join(comment, ",")
	code.WriteString(fmt.Sprintf("@SP // %vPop to the address in D\r\n", str))
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
}

func setTrueOrFalseToD(code *bytes.Buffer, comp string, jump string) {
	code.WriteString(fmt.Sprintf("@TRUE%v // Set true or false to D\r\n", labelIndex))
	code.WriteString(fmt.Sprintf("%v;%v\r\n", comp, jump))

	// False: set 0 to D
	code.WriteString("@0 // False: set 0 to D\r\n")
	code.WriteString("D=A\r\n")
	code.WriteString(fmt.Sprintf("@TFEND%v\r\n", labelIndex))
	code.WriteString("0;JMP\r\n")

	// True: set -1 to D
	code.WriteString(fmt.Sprintf("(TRUE%v)\r\n", labelIndex))
	code.WriteString("@1 // True: set -1 to D\r\n")
	code.WriteString("D=-A\r\n")

	code.WriteString(fmt.Sprintf("(TFEND%v)\r\n", labelIndex))
	labelIndex++
}

func WriteArithmetic(op parser.ALOperator) string {
	var code bytes.Buffer

	// Pop operand y from the stack to R13
	popToD(&code, fmt.Sprintf("[Start:WriteArithmetic(%v)]", op))
	code.WriteString("@13 // Pop y to R13\r\n")
	code.WriteString("M=D\r\n")

	// If op is a binary operator, Pop operand x from the stack to R14
	switch op {
	case parser.ADD, parser.SUB, parser.EQ, parser.GT, parser.LT, parser.AND, parser.OR:
		popToD(&code)
		code.WriteString("@14 // Pop x to R14\r\n")
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

func setAddressToD(code *bytes.Buffer, segment string, index int, comment ...string) {
	str := strings.Join(comment, ",")
	segsym := segment2Symbol(segment)
	code.WriteString(fmt.Sprintf("@%v // %v Set segment + index address to D\r\n", index, str))
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
		popToD(&code, fmt.Sprintf("[Start:WritePushPop - pop(%v, %v, %v)] ", cmdType, segment, index))
		code.WriteString("@13 // Load poped value to R13\r\n")
		code.WriteString("M=D\r\n")

		// Calculate segment + index and set the address to R14
		setAddressToD(&code, segment, index)
		code.WriteString("@14 // Load segment + index address to R14\r\n")
		code.WriteString("M=D\r\n")

		// Write the value in R13 to the address in R14
		code.WriteString("@13 // Write the value in R13 to the address in R14\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@14\r\n")
		code.WriteString("A=M\r\n")
		code.WriteString("M=D\r\n")

	case parser.C_PUSH:
		// Load to D
		switch segment {
		case "constant":
			code.WriteString(fmt.Sprintf("@%v // [Start:WritePushPop - push(%v, %v, %v)]\r\n", index, cmdType, segment, index))
			code.WriteString("D=A\r\n")
		default:
			setAddressToD(&code, segment, index, fmt.Sprintf("[Start:WritePushPop - push(%v, %v, %v)]", cmdType, segment, index))
			code.WriteString("A=D\r\n")
			code.WriteString("D=M\r\n")
		}

		// Push
		pushD(&code)
	}
	return code.String()

}

func WriteLabel(label string) string {
	var code bytes.Buffer
	code.WriteString(fmt.Sprintf("(%v)\r\n", label))
	return code.String()
}

func WriteGoto(label string) string {
	var code bytes.Buffer
	code.WriteString(fmt.Sprintf("@%v // [Start:WriteGoto(%v)]\r\n", label, label))
	code.WriteString("0;JMP\r\n")
	return code.String()
}

func WriteGotoA() string {
	var code bytes.Buffer
	code.WriteString("0;JMP // [Start:WriteGotoA()\r\n")
	return code.String()
}

func WriteIf(label string) string {
	var code bytes.Buffer
	popToD(&code)
	code.WriteString(fmt.Sprintf("@%v // [Start:WriteIf(%v)]\r\n", label, label))
	code.WriteString("D;JNE\r\n")
	return code.String()
}

func WriteFunction(name string, nLocals int) string {
	var code bytes.Buffer
	code.WriteString(WriteLabel(name))
	// Initialize local variables
	for i := 0; i < nLocals; i++ {
		code.WriteString(WritePushPop(parser.C_PUSH, "constant", 0))
	}

	return code.String()
}

func WriteReturn() string {
	var code bytes.Buffer
	// CALLEE_FRAME = LCL
	code.WriteString("@LCL // [Start:WriteReturn] R13 = LCL\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@R15\r\n")
	code.WriteString("M=D\r\n")

	// Pop the return value to ARG and set SP to ARG+1
	code.WriteString(WritePushPop(parser.C_POP, "argument", 0))
	code.WriteString("@ARG\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@SP\r\n")
	code.WriteString("M=D+1\r\n")

	// Restore caller registers
	// THAT = *(CALLEE_FRAME-1)
	code.WriteString("@R15\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@THAT\r\n")
	code.WriteString("M=D\r\n")

	// THIS = *(CALLEE_FRAME-2)
	code.WriteString("@R15\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@THIS\r\n")
	code.WriteString("M=D\r\n")

	// ARG = *(CALLEE_FRAME-3)
	code.WriteString("@R15\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@ARG\r\n")
	code.WriteString("M=D\r\n")

	// LCL = *(CALLEE_FRAME-4)
	code.WriteString("@R15\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@LCL\r\n")
	code.WriteString("M=D\r\n")

	// RET = *(CALLEE_FRAME-5)
	code.WriteString("@R15\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	// Jump to RET
	code.WriteString(WriteGotoA())
	return code.String()
}
