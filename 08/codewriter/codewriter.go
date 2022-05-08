package codewriter

import (
	"bytes"
	"fmt"
	"log"

	"vm/parser"
)

var labelIndex int

// Bootstrap code
// SP=256
// call Sys.init
func Bootstrap() []string {
	var code []string
	code = append(code, "@256")
	code = append(code, "D=A")
	code = append(code, "@SP")
	code = append(code, "M=D")
	code = append(code, WriteCall("Sys.init", 0)...)
	return code
}

func pushD() []string {
	var code []string
	code = append(code, "@SP // Push the value at the address in D")
	code = append(code, "A=M")
	code = append(code, "M=D")
	code = append(code, "@SP")
	code = append(code, "M=M+1")
	return code
}

func popToD() []string {
	var code []string
	code = append(code, "@SP // Pop to the address in D")
	code = append(code, "M=M-1")
	code = append(code, "A=M")
	code = append(code, "D=M")
	return code
}

func setTrueOrFalseToD(comp string, jump string) []string {
	var code []string
	code = append(code, fmt.Sprintf("@TRUE%v // Set true or false to D", labelIndex))
	code = append(code, fmt.Sprintf("%v;%v", comp, jump))

	// False: set 0 to D
	code = append(code, "@0 // False: set 0 to D")
	code = append(code, "D=A")
	code = append(code, fmt.Sprintf("@TFEND%v", labelIndex))
	code = append(code, "0;JMP")

	// True: set -1 to D
	code = append(code, fmt.Sprintf("(TRUE%v)", labelIndex))
	code = append(code, "@1 // True: set -1 to D")
	code = append(code, "D=-A")

	code = append(code, fmt.Sprintf("(TFEND%v)", labelIndex))
	labelIndex++
	return code
}

func WriteArithmetic(op parser.ALOperator) []string {
	var code []string

	// Pop operand y from the stack to R13
	c := popToD()
	c[0] += fmt.Sprintf("// [Start:WriteArithmetic(%v)]", op)
	code = append(code, c...)
	code = append(code, "@13 // Pop y to R13")
	code = append(code, "M=D")

	// If op is a binary operator, Pop operand x from the stack to R14
	switch op {
	case parser.ADD, parser.SUB, parser.EQ, parser.GT, parser.LT, parser.AND, parser.OR:
		code = append(code, popToD()...)
		code = append(code, "@14 // Pop x to R14")
		code = append(code, "M=D")
	}

	// Calculate and load the result to D
	switch op {
	case parser.ADD:
		code = append(code, "@14 // add")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D+M")
	case parser.SUB:
		code = append(code, "@14 // sub")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D-M")
	case parser.NEG:
		code = append(code, "@13 // neg")
		code = append(code, "D=-M")
	case parser.EQ:
		code = append(code, "@14 // eq")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D-M")
		code = append(code, setTrueOrFalseToD("D", "JEQ")...) // x-y==0
	case parser.GT:
		code = append(code, "@14 // gt")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D-M")
		code = append(code, setTrueOrFalseToD("D", "JGT")...) // x-y>0
	case parser.LT:
		code = append(code, "@14 // lt")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D-M")
		code = append(code, setTrueOrFalseToD("D", "JLT")...) // x-y<0
	case parser.AND:
		code = append(code, "@14 // and")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D&M")
	case parser.OR:
		code = append(code, "@14 // or")
		code = append(code, "D=M")
		code = append(code, "@13")
		code = append(code, "D=D|M")
	case parser.NOT:
		code = append(code, "@13 // not")
		code = append(code, "D=!M")
	}
	// Push D to the stack
	code = append(code, pushD()...)

	return code
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

func setAddressToD(segment string, index int) []string {
	var code []string
	segsym := segment2Symbol(segment)
	code = append(code, fmt.Sprintf("@%v // Set segment + index address to D", index))
	code = append(code, "D=A")
	code = append(code, fmt.Sprintf("@%v", segsym))
	switch segment {
	case "local", "argument", "this", "that":
		code = append(code, "D=D+M")
	case "pointer", "temp":
		code = append(code, "D=D+A")
	}
	return code
}

func WritePushPop(cmdType parser.CommandType, segment string, index int) []string {
	var code []string
	switch cmdType {
	case parser.C_POP:

		// Pop to R13
		c := popToD()
		c[0] += fmt.Sprintf("[Start:WritePushPop - pop(%v, %v, %v)] ", cmdType, segment, index)
		code = append(code, c...)
		code = append(code, "@13 // Load poped value to R13")
		code = append(code, "M=D")

		// Calculate segment + index and set the address to R14
		code = append(code, setAddressToD(segment, index)...)
		code = append(code, "@14 // Load segment + index address to R14")
		code = append(code, "M=D")

		// Write the value in R13 to the address in R14
		code = append(code, "@13 // Write the value in R13 to the address in R14")
		code = append(code, "D=M")
		code = append(code, "@14")
		code = append(code, "A=M")
		code = append(code, "M=D")

	case parser.C_PUSH:
		// Load to D
		switch segment {
		case "constant":
			code = append(code, fmt.Sprintf("@%v // [Start:WritePushPop - push(%v, %v, %v)]", index, cmdType, segment, index))
			code = append(code, "D=A")
		default:
			c := setAddressToD(segment, index)
			c[0] += fmt.Sprintf("[Start:WritePushPop - push(%v, %v, %v)]", cmdType, segment, index)
			code = append(code, c...)
			code = append(code, "A=D")
			code = append(code, "D=M")
		}

		// Push
		code = append(code, pushD()...)
	}
	return code

}

func WritePushPopStatic(cmdType parser.CommandType, segment string, index int, vmName string) []string {
	var code []string
	switch cmdType {
	case parser.C_POP:

		// Pop to R13
		c := popToD()
		c[0] += fmt.Sprintf("// [Start:WritePushPopStatic - pop(%v, %v, %v)] ", cmdType, segment, index)
		code = append(code, c...)
		code = append(code, "@13 // Load poped value to R13")
		code = append(code, "M=D")

		// Set static variable's address to D
		code = append(code, fmt.Sprintf("@%v.%v", vmName, index))
		code = append(code, "D=A")
		code = append(code, "@14")
		code = append(code, "M=D")

		// Write the value in R13 to the address in R14
		code = append(code, "@13 // Write the value in R13 to the address in R14")
		code = append(code, "D=M")
		code = append(code, "@14")
		code = append(code, "A=M")
		code = append(code, "M=D")

	case parser.C_PUSH:
		// Set static variable's address to D
		code = append(code, fmt.Sprintf("@%v.%v", vmName, index))
		code = append(code, "D=M")

		// Push
		code = append(code, pushD()...)

	}
	return code
}

func WriteLabel(label string) []string {
	var code []string
	code = append(code, fmt.Sprintf("(%v)", label))
	return code
}

var labelCnt int = 0

func generateUniqueLabel(label string) string {
	ret := fmt.Sprintf("%v%v", label, labelCnt)
	labelCnt++
	return ret
}

func WriteGoto(label string) []string {
	var code []string
	code = append(code, fmt.Sprintf("@%v // [Start:WriteGoto(%v)]", label, label))
	code = append(code, "0;JMP")
	return code
}

func WriteGotoA() string {
	var code bytes.Buffer
	code.WriteString("0;JMP // [Start:WriteGotoA()")
	return code.String()
}

func WriteIf(label string) []string {
	var code []string
	code = append(code, popToD()...)
	code = append(code, fmt.Sprintf("@%v // [Start:WriteIf(%v)]", label, label))
	code = append(code, "D;JNE")
	return code
}

func WriteIfJLE(label string) []string {
	var code []string
	code = append(code, popToD()...)
	code = append(code, fmt.Sprintf("@%v // [Start:WriteIfJEQ(%v)]", label, label))
	code = append(code, "D;JLE")
	return code
}

func WriteFunction(name string, nLocals int) []string {
	var code []string
	code = append(code, WriteLabel(name)...)
	// Initialize local variables
	for i := 0; i < nLocals; i++ {
		code = append(code, WritePushPop(parser.C_PUSH, "constant", 0)...)
	}
	return code
}

// Write code for return.
// See the chapter 8 slide p43- https://drive.google.com/file/d/1lBsaO5XKLkUgrGY6g6vLMsiZo6rWxlYJ/view
// NOTE: Contract between caller and callee
// - A return value(must exist) had to be pushed by callee on the top of the stack. See WriteCall().
// - A return address had to be pushed by caller on LCL-5.
func WriteReturn() []string {
	var code []string

	// R15 = *(LCL-5)
	// Before copying the return value on *ARG(the top of the calle's frame), we have to memorize the return address first.
	// Because if the function don't have any argurements, the top of the frame is the return address and will be overwritten by the return value.
	code = append(code, "@LCL")
	code = append(code, "D=M-1")
	code = append(code, "D=D-1")
	code = append(code, "D=D-1")
	code = append(code, "D=D-1")
	code = append(code, "D=D-1")
	code = append(code, "A=D")
	code = append(code, "D=M")
	code = append(code, "@R15") // return address
	code = append(code, "M=D")

	// Pop the return value to *ARG(the top of the caller's frame)
	code = append(code, WritePushPop(parser.C_POP, "argument", 0)...)

	// SP = ARG+1
	code = append(code, "@ARG")
	code = append(code, "D=M")
	code = append(code, "@SP")
	code = append(code, "M=D+1")

	// R13 is just a counter
	code = append(code, "@LCL // [Start:WriteReturn] R13 = LCL")
	code = append(code, "D=M")
	code = append(code, "@R13")
	code = append(code, "M=D")

	// Restore caller registers
	// THAT = *(LCL-1)
	code = append(code, "@R13")
	code = append(code, "M=M-1")
	code = append(code, "A=M")
	code = append(code, "D=M")
	code = append(code, "@THAT")
	code = append(code, "M=D")

	// THIS = *(LCL-2)
	code = append(code, "@R13")
	code = append(code, "M=M-1")
	code = append(code, "A=M")
	code = append(code, "D=M")
	code = append(code, "@THIS")
	code = append(code, "M=D")

	// ARG = *(LCL-3)
	code = append(code, "@R13")
	code = append(code, "M=M-1")
	code = append(code, "A=M")
	code = append(code, "D=M")
	code = append(code, "@ARG")
	code = append(code, "M=D")

	// LCL = *(LCL-4)
	code = append(code, "@R13")
	code = append(code, "M=M-1")
	code = append(code, "A=M")
	code = append(code, "D=M")
	code = append(code, "@LCL")
	code = append(code, "M=D")

	// Jump to the return address in R15
	code = append(code, "@R15")
	code = append(code, "A=M")
	code = append(code, WriteGotoA())
	return code
}

// Write code for call.
// See chapter 8 slide p32-.
// https://drive.google.com/file/d/1lBsaO5XKLkUgrGY6g6vLMsiZo6rWxlYJ/view
func WriteCall(name string, nArgs int) []string {
	var code []string
	// Push return address
	retLabel := generateUniqueLabel("RET")
	code = append(code, fmt.Sprintf("@%v // [Start:WriteCall(%v,%v)] Push return address", retLabel, name, nArgs))
	code = append(code, "D=A")
	code = append(code, pushD()...)

	// Save LCL
	code = append(code, "@LCL")
	code = append(code, "D=M")
	code = append(code, pushD()...)

	// Save ARG
	code = append(code, "@ARG")
	code = append(code, "D=M")
	code = append(code, pushD()...)

	// Save THIS
	code = append(code, "@THIS")
	code = append(code, "D=M")
	code = append(code, pushD()...)

	// Save THAT
	code = append(code, "@THAT")
	code = append(code, "D=M")
	code = append(code, pushD()...)

	// ARG = SP-n-5
	code = append(code, "@SP")
	code = append(code, "D=M")
	for i := 0; i < 5+nArgs; i++ {
		code = append(code, "D=D-1")
	}
	code = append(code, "@ARG")
	code = append(code, "M=D")

	// LCL = SP
	code = append(code, "@SP")
	code = append(code, "D=M")
	code = append(code, "@LCL")
	code = append(code, "M=D")

	// goto f
	code = append(code, WriteGoto(name)...)

	// label for return
	code = append(code, WriteLabel(retLabel)...)
	return code
}
