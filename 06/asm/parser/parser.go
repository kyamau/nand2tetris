package parser

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type (
	CommandType int
)

const (
	A_COMMAND CommandType = iota
	C_COMMAND CommandType = iota
	L_COMMAND CommandType = iota
)

var comment *regexp.Regexp = regexp.MustCompile(`(//).*`)
var spaceTab *regexp.Regexp = regexp.MustCompile(`[\t ]`)
var emptyLine *regexp.Regexp = regexp.MustCompile(`(?m)^\n`)

//var a_command *regexp.Regexp = regexp.MustCompile(`^@[a-zA-Z0-9_\.\$:]+$`)
//var label *regexp.Regexp = regexp.MustCompile(`^\([a-zA-Z0-9_\.\$:]+\)$`)

type Command struct {
	command     string
	commandType CommandType
}

type Parser struct {
	commands        []string
	current         int
	hasMoreCommands bool
}

func (p *Parser) HasMoreCommands() bool {
	return len(p.commands)-1 > p.current
}

func (p *Parser) Advance() {
	if !p.HasMoreCommands() {
		log.Fatal("No more commands")
	}
	p.current++
}

func (p *Parser) CommandType() CommandType {
	cmd := p.Current()
	if cmd[0] == '@' {
		return A_COMMAND
	} else if cmd[0] == '(' {
		return L_COMMAND
	} else {
		return C_COMMAND
	}
}

func (p *Parser) Symbol() string {
	cmd := p.Current()
	if p.CommandType() != A_COMMAND {
		log.Fatalf("Can't get symbol from command other than A : %v", cmd)
	}
	return cmd[1:]
}

func (p *Parser) Dest() string {
	cmd := p.Current()
	if p.CommandType() != C_COMMAND {
		log.Fatalf("Can't get dest from command other than C : %v", cmd)
	}
	spl := strings.Split(cmd, "=")
	if len(spl) == 1 {
		return "null"
	} else {
		return spl[0]
	}
}

func (p *Parser) Comp() string {
	cmd := p.Current()
	if p.CommandType() != C_COMMAND {
		log.Fatalf("Can't get comp from command other than C : %v", cmd)
	}
	spl := strings.Split(cmd, "=")
	if len(spl) == 1 {
		return strings.Split(spl[0], ";")[0]
	} else {
		return strings.Split(spl[1], ";")[0]
	}
}

func (p *Parser) Jump() string {
	cmd := p.Current()
	if p.CommandType() != C_COMMAND {
		log.Fatalf("Can't get jump from command other than C : %v", cmd)
	}
	spl := strings.Split(cmd, ";")
	if len(spl) == 1 {
		return "null"
	} else {
		return spl[1]
	}
}

func (p *Parser) Current() string {
	return p.commands[p.current]
}

func removeIrrelevants(lines []string) []string {
	ret := make([]string, 0)
	for _, l := range lines {
		l = comment.ReplaceAllString(l, "")
		l = spaceTab.ReplaceAllString(l, "")
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
		return nil, fmt.Errorf("reading asm code failed : %v", err)
	}
	lines := strings.Split(s, "\r\n")
	lines = removeIrrelevants(lines)
	p := &Parser{commands: lines, current: -1}
	return p, nil
}
