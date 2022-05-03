package vmwriter

import (
	"fmt"
	"os"
	"strings"
)

type Stack struct {
	stack []string
	sp    int
}

func NewStack() *Stack {
	return &Stack{[]string{}, -1}
}

func (s *Stack) Push(e string) {
	s.stack = append(s.stack, e)
	s.sp++
}

func (s *Stack) Pop() (string, bool) {
	if s.sp-1 < 0 {
		return "", false
	}
	poped := s.stack[s.sp-1]
	s.stack = s.stack[:s.sp-1]
	s.sp--
	return poped, true
}

type VMWriter struct {
	lines         []string
	OperatorStack Stack
}

func (w *VMWriter) Add(code string) {
	w.lines = append(w.lines, code)
}

func (w *VMWriter) Code() []string {
	return w.lines
}

func NewVMWriter() (*VMWriter, error) {
	vmWriter := VMWriter{[]string{}, Stack{}}
	return &vmWriter, nil
}

func WriteCode(vmCode []string, filepath string) error {
	err := os.WriteFile(filepath, []byte(strings.Join(vmCode, "\n")), 0666)
	if err != nil {
		return err
	}
	return nil
}

func PushCode(segment string, index int) string {
	return fmt.Sprintf("push %s %d", segment, index)
}

func PopCode(segment string, index int) string {
	return fmt.Sprintf("pop %s %d", segment, index)
}

func ArithmeticCode(command string) string {
	return command
}

func LabelCode(label string) string {
	return fmt.Sprintf("label %s", label)
}

func GotoCode(label string) string {
	return fmt.Sprintf("goto %s", label)
}

func IfCode(label string) string {
	return fmt.Sprintf("if-goto %s", label)
}

func FunctionCode(name string, nLocals int) string {
	return fmt.Sprintf("function %s %d", name, nLocals)
}

func CallCode(name string, nArgs int) string {
	return fmt.Sprintf("call %s %d", name, nArgs)
}

func ReturnCode() string {
	return fmt.Sprintf("return")
}
