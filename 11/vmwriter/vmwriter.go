package vmwriter

import (
	"fmt"
	"os"
	"strings"
)

type VMWriter struct {
	lines []string
}

func (w *VMWriter) Add(code string) {
	w.lines = append(w.lines, code)
}

func (w *VMWriter) Code() []string {
	return w.lines
}

func NewVMWriter() (*VMWriter, error) {
	vmWriter := VMWriter{make([]string, 0)}
	return &vmWriter, nil
}

func WriteCode(vmCode []string, filepath string) error {
	err := os.WriteFile(filepath, []byte(strings.Join(vmCode, "\n")), 0666)
	if err != nil {
		return err
	}
	return nil
}

func Push(segment string, index int) string {
	return fmt.Sprintf("push %s %d", segment, index)
}

func Pop(segment string, index int) string {
	return fmt.Sprintf("pop %s %d", segment, index)
}

func Arithmetic(command string) string {
	return command
}

func Label(label string) string {
	return fmt.Sprintf("label %s", label)
}

func Goto(label string) string {
	return fmt.Sprintf("goto %s", label)
}

func If(label string) string {
	return fmt.Sprintf("if-goto %s", label)
}

func Function(name string, nLocals int) string {
	return fmt.Sprintf("function %s %d", name, nLocals)
}

func Call(name string, nArgs int) string {
	return fmt.Sprintf("call %s %d", name, nArgs)
}

func Return(name string, nArgs int) string {
	return fmt.Sprintf("return")
}
