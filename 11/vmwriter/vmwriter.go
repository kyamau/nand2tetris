package vmwriter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Stack struct {
	stack []string
	sp    int
}

func NewStack() *Stack {
	initialSize := 20
	return &Stack{make([]string, initialSize), 0}
}

func (s *Stack) Push(e string) {
	s.stack[s.sp] = e
	s.sp++
}

// Pop. Panic if the stack is empty.
func (s *Stack) Pop() string {
	poped := s.stack[s.sp-1]
	s.sp--
	return poped
}

// See the top of the stack. Panic if the stack is empty.
func (s *Stack) Top() string {
	return s.stack[s.sp-1]
}

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
	vmWriter := VMWriter{[]string{}}
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

func IfGotoCode(label string) string {
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

type LabelManager struct {
	counter    map[string]int
	ifStack    Stack
	whileStack Stack
}

func NewLabelManager() LabelManager {
	return LabelManager{map[string]int{"while": -1, "if": -1}, *NewStack(), *NewStack()}
}

func (l *LabelManager) StartWhile() {
	l.counter["while"]++
	l.whileStack.Push(strconv.Itoa(l.counter["while"]))
}

func (l *LabelManager) EndWhile() {
	l.whileStack.Pop()
}

func (l *LabelManager) WhileExpLabel() string {
	return fmt.Sprintf("WHILE_EXP%s", l.whileStack.Top())
}

func (l *LabelManager) WhileEndLabel() string {
	return fmt.Sprintf("WHILE_END%s", l.whileStack.Top())
}

func (l *LabelManager) StartIf() {
	l.counter["if"]++
	l.ifStack.Push(strconv.Itoa(l.counter["if"]))
}

func (l *LabelManager) EndIf() {
	l.ifStack.Pop()
}

func (l *LabelManager) IfTrueLabel() string {
	return fmt.Sprintf("IF_TRUE%s", l.ifStack.Top())
}

func (l *LabelManager) IfFalseLabel() string {
	return fmt.Sprintf("IF_FALSE%s", l.ifStack.Top())
}

func (l *LabelManager) IfEndLabel() string {
	return fmt.Sprintf("IF_END%s", l.ifStack.Top())
}
