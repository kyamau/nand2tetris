package codewriter

import (
	"bytes"
	"fmt"
	"log"

	"vm/parser"
)

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

func WriteArithmetic(op parser.ALOperator) string {
	var code bytes.Buffer
	// POP from the stack to R13
	popToD(&code)
	code.WriteString("@13\r\n")
	code.WriteString("M=D\r\n")
	// POP from the stack to R14
	popToD(&code)
	code.WriteString("@14\r\n")
	code.WriteString("M=D\r\n")

	// Caluculate and load the result to D
	switch op {
	case parser.ADD:
		code.WriteString("@14\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D+M\r\n")
	case parser.SUB:
		code.WriteString("@14\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n")
		code.WriteString("D=D-M\r\n")
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
