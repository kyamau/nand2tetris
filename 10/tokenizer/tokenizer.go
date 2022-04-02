package tokenizer

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Tokenizer struct {
	src     string
	tokens  []Token
	current int
}

type Token interface {
	String() string
	Type() string
	Pos() []int
}

type GenericToken struct {
	token     string
	tokenType string
	pos       []int
}

func (t *GenericToken) String() string {
	return t.token
}

func (t *GenericToken) Type() string {
	return t.tokenType
}

func (t *GenericToken) Pos() []int {
	return t.pos
}

const (
	STR_CONST  = "stringConstant"
	SYMBOL     = "symbol"
	INT_CONST  = "integerConstant"
	IDENTIFIER = "identifier"
	KEYWORD    = "keyword"
)

type StrConst struct {
	*GenericToken
}

func NewStrConst(s string, pos []int) *StrConst {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "\n", "")
	return &StrConst{&GenericToken{token: s, tokenType: "stringConstant", pos: pos}}
}

type Symbol struct {
	*GenericToken
}

func NewSymbol(s string, pos []int) *Symbol {
	return &Symbol{&GenericToken{token: s, tokenType: SYMBOL, pos: pos}}
}

const (
	CBRACKET_L    = "{"
	CBRACKET_R    = "}"
	PARENTHESES_L = "("
	PARENTHESES_R = ")"
	SBRACKET_L    = "["
	SBRACKET_R    = "]"
	PERIOD        = "."
	COMMA         = ","
	SEMICOLON     = ";"
	PLUS          = "+"
	MINUS         = "-"
	ASTERISK      = "*"
	SLASH         = "/"
	AND           = "&"
	OR            = "|"
	LESS          = "<"
	GREATER       = ">"
	EQUAL         = "="
	TILDA         = "~"
)

type IntConst struct {
	*GenericToken
	value int
}

func (t *IntConst) Int() int {
	return t.value
}

func NewIntConst(s string, pos []int) (*IntConst, error) {
	value, err := strconv.Atoi(s)
	if value < 0 || value > 32768 {
		return nil, errors.New(fmt.Sprintf("Integer constant must be in [0, 32768]: %v", value))
	}
	return &IntConst{GenericToken: &GenericToken{token: s, tokenType: INT_CONST, pos: pos}, value: value}, err
}

type Identifier struct {
	*GenericToken
}

func NewIdentifier(s string, pos []int) *Identifier {
	return &Identifier{&GenericToken{token: s, tokenType: "identifier", pos: pos}}
}

type Keyword struct {
	*GenericToken
}

func NewKeyword(s string, pos []int) *Keyword {
	return &Keyword{&GenericToken{token: s, tokenType: KEYWORD, pos: pos}}
}

const (
	CLASS       = "class"
	CONSTRUCTOR = "constructor"
	FUNCTION    = "function"
	METHOD      = "method"
	FIELD       = "field"
	STATIC      = "static"
	VAR         = "var"
	INT         = "int"
	CHAR        = "char"
	BOOLEAN     = "boolean"
	VOID        = "void"
	TRUE        = "true"
	FALSE       = "false"
	NULL        = "null"
	THIS        = "this"
	LET         = "let"
	DO          = "do"
	IF          = "if"
	ELSE        = "else"
	WHILE       = "while"
	RETURN      = "return"
)

func (t *Tokenizer) HasMoreTokens() bool {
	return len(t.tokens)-1 > t.current
}

func (t *Tokenizer) Current() Token {
	return t.tokens[t.current]
}

func (t *Tokenizer) Advance() error {
	if !t.HasMoreTokens() {
		return errors.New("Couldn't advance. No more tokens.")
	}
	t.current++
	return nil
}

func (t *Tokenizer) LookAhead(offset int) (Token, error) {
	if t.current+offset > len(t.tokens)-1 {
		return nil, errors.New("Look ahead faild. Index out of range")
	}
	return t.tokens[t.current+offset], nil
}

var singleLineComment *regexp.Regexp = regexp.MustCompile(`(//).*`)
var multiLineComment *regexp.Regexp = regexp.MustCompile(`(?s)(/\*\*).*?(\*/)`)
var emptyLine *regexp.Regexp = regexp.MustCompile(`(?m)^\n`)

func preprocess(src string) string {
	// CRLF -> LF
	src = strings.ReplaceAll(src, "\r\n", "\n")
	src = singleLineComment.ReplaceAllString(src, "")
	src = multiLineComment.ReplaceAllString(src, "")
	src = emptyLine.ReplaceAllString(src, "")
	return src
}

var tokenRegexp *regexp.Regexp = regexp.MustCompile(`(?P<strConst>"[^"]+")|(?P<idOrKeyword>[a-zA-Z_][a-zA-Z0-9_]*)|(?P<symbol>[{}\(\)\[\]\.\,;\+\-\*\/&\|<>=~])|(?P<intConst>[0-9]+)`)

var lf *regexp.Regexp = regexp.MustCompile(`\n`)

func lineCountCummulativeMap(src string) map[int]int {
	matches := lf.FindAllStringIndex(src, -1)
	m := make(map[int]int, 0)
	for _, j := range matches {
		m[j[0]] += 1
	}
	return m
}

// Return array of [row,column] at each character index
// Row and column start with 1.
func charIndex2LineColumnArray(src string) [][]int {
	matchedLFs := lf.FindAllStringIndex(src, -1)
	lineColumn := make([][]int, len(src))
	nLF := 0
	posLF := matchedLFs[nLF][0]
	prevPosLF := -1
	for i := 0; i < len(lineColumn); i++ {
		if posLF < i {
			nLF++
			prevPosLF = posLF
			if nLF < len(matchedLFs) {
				posLF = matchedLFs[nLF][0]
			} else {
				// No more LFs.
				posLF = len(src)
			}
		}
		currentLine := nLF + 1
		currentColumn := i - prevPosLF // How many chars after previous LF
		lineColumn[i] = []int{currentLine, currentColumn}
	}
	return lineColumn
}

func tokenize(src string) ([]Token, error) {
	charIndex2LineColumn := charIndex2LineColumnArray(src)
	tokens := make([]Token, 0)
	matchStrings := tokenRegexp.FindAllStringSubmatch(src, -1) // Use capturing groups.
	matchIndices := tokenRegexp.FindAllStringIndex(src, -1)
	groupNames := tokenRegexp.SubexpNames()
	for i, matchString := range matchStrings {
		lineColumn := charIndex2LineColumn[matchIndices[i][0]]
		for j, name := range groupNames[1:] { // SubexpNames()[0] is always empty.
			m := matchString[j+1]
			if m != "" {
				var t Token
				switch name {
				case "strConst":
					t = NewStrConst(m, lineColumn)
				case "idOrKeyword":
					switch m {
					case CLASS, CONSTRUCTOR, FUNCTION, METHOD, FIELD, STATIC, VAR, INT, CHAR, BOOLEAN, VOID, TRUE, FALSE, NULL, THIS, LET, DO, IF, ELSE, WHILE, RETURN:
						t = NewKeyword(m, lineColumn)
					default:
						t = NewIdentifier(m, lineColumn)
					}
				case "symbol":
					t = NewSymbol(m, lineColumn)
				case "intConst":
					var err error
					t, err = NewIntConst(m, lineColumn)
					if err != nil {
						return nil, fmt.Errorf("Invalid integer constant: %w", err)
					}
				default:
					return nil, errors.New(fmt.Sprintf("Unknown token: %v", m))
				}
				tokens = append(tokens, t)
			}
		}
	}
	return tokens, nil
}

type TokensXml struct {
	XMLName xml.Name `xml:"tokens"`
	Tokens  []Token
}

func (tokensXml TokensXml) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "tokens"
	e.EncodeToken(start)
	for _, t := range tokensXml.Tokens {
		e.EncodeElement(fmt.Sprintf(" %v ", t.String()), xml.StartElement{Name: xml.Name{Local: t.Type()}})
	}
	e.EncodeToken(start.End())
	return nil
}

func (t *Tokenizer) XML() string {
	tokensXml := TokensXml{Tokens: t.tokens}
	buf, _ := xml.MarshalIndent(tokensXml, "", "  ")
	return string(buf)
}

func (t *Tokenizer) Tokenize() error {
	t.src = preprocess(t.src)
	tokens, err := tokenize(t.src)
	if err != nil {
		return fmt.Errorf("Failed to tokenize: %v", err)
	}
	t.tokens = tokens
	return nil
}

func NewTokenizer(r io.Reader) (*Tokenizer, error) {
	b, err := ioutil.ReadAll(r)
	src := string(b)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .jack: %v", err)
	}
	p := &Tokenizer{src, nil, 0}
	fmt.Println("")
	return p, nil
}
