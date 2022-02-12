package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Parser struct {
	commands        []string
	current         int
	hasMoreCommands bool
}

type (
	CommandType int
)

const (
	C_ARITHMETIC CommandType = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
	UNKNOWN_COMMAND
)

type (
	ALOperator int
)

const (
	ADD ALOperator = iota
	SUB
	NEG
	EQ
	GT
	LT
	AND
	OR
	NOT
	UNKNOWN_ALOPERATOR
)

func ALOperatorFromString(s string) (ALOperator, error) {
	switch s {
	case "add":
		return ADD, nil
	case "sub":
		return SUB, nil
	case "neg":
		return NEG, nil
	case "eq":
		return EQ, nil
	case "gt":
		return GT, nil
	case "lt":
		return LT, nil
	case "and":
		return AND, nil
	case "or":
		return OR, nil
	case "not":
		return NOT, nil
	}
	return UNKNOWN_ALOPERATOR, fmt.Errorf("Not an arithmetic operator : %v", s)
}

var comment *regexp.Regexp = regexp.MustCompile(`(//).*`)
var spaceTabTrim *regexp.Regexp = regexp.MustCompile(`^[\t ]+|[\t ]+$`)
var emptyLine *regexp.Regexp = regexp.MustCompile(`(?m)^\n`)

func (p *Parser) HasMoreCommands() bool {
	return len(p.commands)-1 > p.current
}

func (p *Parser) Advance() {
	if !p.HasMoreCommands() {
		log.Fatal("No more commands")
	}
	p.current++
}

func (p *Parser) Current() string {
	return p.commands[p.current]
}

// Implement only C_ARITHMETIC, C_PUSH, C_POP for the project 07
func (p *Parser) CommandType() CommandType {
	cmdLine := p.Current()
	cmd := strings.Split(cmdLine, " ")[0]

	switch cmd {
	case "pop":
		return C_POP
	case "push":
		return C_PUSH
	case "label":
		return C_LABEL
	case "goto":
		return C_GOTO
	case "if-goto":
		return C_IF
	case "function":
		return C_FUNCTION
	case "return":
		return C_RETURN
	}
	// Is arithmetic operator?
	_, err := ALOperatorFromString(cmd)
	if err == nil {
		return C_ARITHMETIC
	}

	log.Fatalf("Can't interpret operator : %v", err)
	return UNKNOWN_COMMAND
}

func (p *Parser) Arg1() string {
	if p.CommandType() == C_RETURN {
		log.Fatalf("Don't call Arg1() for C_RETURN : %v", p.Current())
	}
	cmdLine := p.Current()
	tokens := strings.Split(cmdLine, " ")
	if p.CommandType() == C_ARITHMETIC {
		return tokens[0]
	}
	return tokens[1]
}

func (p *Parser) Arg2() string {
	if p.CommandType() == C_RETURN {
		log.Fatalf("Don't call Arg2() other than C_PUSH, C_POP, C_FUNCTION, C_CALL : %v", p.Current())
	}
	cmdLine := p.Current()
	tokens := strings.Split(cmdLine, " ")
	return tokens[2]
}

func removeIrrelvants(lines []string) []string {
	ret := make([]string, 0)
	for _, l := range lines {
		l = comment.ReplaceAllString(l, "")
		l = spaceTabTrim.ReplaceAllString(l, "")
		if len(l) > 0 {
			ret = append(ret, l)
		}
	}
	return ret
}

func NewParser(r io.Reader) (*Parser, error) {
	b, err := ioutil.ReadAll(r)
	s := string(b)
	if err != nil {
		return nil, fmt.Errorf("Reading vm code failed : %v", err)
	}
	lines := strings.Split(s, "\r\n")
	lines = removeIrrelvants(lines)
	p := &Parser{commands: lines, current: -1}
	return p, nil
}
