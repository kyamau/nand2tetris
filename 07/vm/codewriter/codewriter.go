package codewriter

import (
	"bytes"
	"fmt"
	"log"

	"vm/parser"
)

func WriteArithmetic(op parser.ALOperator) string {
	var code bytes.Buffer
	// POP from the stack to R13
	code.WriteString("@SP\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@13\r\n") // Load poped item to R13
	code.WriteString("M=D\r\n")
	// POP from the stack to R14
	code.WriteString("@SP\r\n")
	code.WriteString("M=M-1\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("D=M\r\n")
	code.WriteString("@14\r\n") // Load poped item to R14
	code.WriteString("M=D\r\n")
	switch op {

	// Caluculate and load the result to D
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
	code.WriteString("@SP\r\n")
	code.WriteString("A=M\r\n")
	code.WriteString("M=D\r\n")
	code.WriteString("@SP\r\n")
	code.WriteString("M=M+1\r\n")

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

func WritePushPop(cmdType parser.CommandType, segment string, index int) string {
	var code bytes.Buffer
	switch cmdType {
	case parser.C_POP:

		// Load from the stack and SP--
		code.WriteString("@SP\r\n")
		code.WriteString("M=M-1\r\n")
		code.WriteString("A=M\r\n")
		code.WriteString("D=M\r\n")
		code.WriteString("@13\r\n") // Load poped item to R13
		code.WriteString("M=D\r\n")

		// Calculate segment + i
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
		code.WriteString("@14\r\n") // Load destination address to R14
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
			code.WriteString("A=D\r\n")
			code.WriteString("D=M\r\n")
		}

		// Push
		code.WriteString("@SP\r\n")
		code.WriteString("A=M\r\n")
		code.WriteString("M=D\r\n")
		code.WriteString("@SP\r\n")
		code.WriteString("M=M+1\r\n")
	}
	return code.String()

}
