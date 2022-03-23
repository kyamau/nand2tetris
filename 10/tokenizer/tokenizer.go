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
}

const (
	STR_CONST  = "stringConstant"
	SYMBOL     = "symbol"
	INT_CONST  = "integerConstant"
	IDENTIFIER = "identifier"
	KEYWORD    = "keyword"
)

type StrConst struct {
	token     string
	tokenType string
}

func (t *StrConst) String() string {
	return t.token
}

func (t *StrConst) Type() string {
	return t.tokenType
}

func NewStrConst(s string) *StrConst {
	s = strings.ReplaceAll(s, "\"", "")
	s = strings.ReplaceAll(s, "\n", "")
	return &StrConst{token: s, tokenType: "stringConstant"}
}

type Symbol struct {
	token     string
	tokenType string
}

func (t *Symbol) String() string {
	return t.token
}

func (t *Symbol) Type() string {
	return t.tokenType
}

func NewSymbol(s string) *Symbol {
	return &Symbol{token: s, tokenType: SYMBOL}
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
	token     string
	tokenType string
	value     int
}

func (t *IntConst) String() string {
	return t.token
}

func (t *IntConst) Int() int {
	return t.value
}

func (t *IntConst) Type() string {
	return t.tokenType
}

func NewIntConst(s string) (*IntConst, error) {
	value, err := strconv.Atoi(s)
	if value < 0 || value > 32768 {
		return nil, errors.New(fmt.Sprintf("Integer constant must be in [0, 32768]: %v", value))
	}
	return &IntConst{token: s, tokenType: INT_CONST, value: value}, err
}

type Identifier struct {
	token     string
	tokenType string
}

func (t *Identifier) String() string {
	return t.token
}

func (t *Identifier) Type() string {
	return t.tokenType
}

func NewIdentifier(s string) *Identifier {
	return &Identifier{token: s, tokenType: "identifier"}
}

type Keyword struct {
	token     string
	tokenType string
}

func (t *Keyword) String() string {
	return t.token
}

func (t *Keyword) Type() string {
	return t.tokenType
}

func NewKeyword(s string) *Keyword {
	return &Keyword{token: s, tokenType: KEYWORD}
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

func (t *Tokenizer) TokenType() string {
	return t.Current().Type()
}

func (t *Tokenizer) Advance() error {
	if !t.HasMoreTokens() {
		return errors.New("Couldn't advance. No more tokens.")
	}
	t.current++
	return nil
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

func tokenize(src string) ([]Token, error) {
	tokens := make([]Token, 0)
	matchs := tokenRegexp.FindAllStringSubmatch(src, -1)
	groupNames := tokenRegexp.SubexpNames()
	for _, match := range matchs {
		for i, name := range groupNames[1:] {
			m := match[i+1]
			if m != "" {
				var t Token
				switch name {
				case "strConst":
					t = NewStrConst(m)
				case "idOrKeyword":
					switch m {
					case CLASS, CONSTRUCTOR, FUNCTION, METHOD, FIELD, STATIC, VAR, INT, CHAR, BOOLEAN, VOID, TRUE, FALSE, NULL, THIS, LET, DO, IF, ELSE, WHILE, RETURN:
						t = NewKeyword(m)
					default:
						t = NewIdentifier(m)
					}
				case "symbol":
					t = NewSymbol(m)
				case "intConst":
					var err error
					t, err = NewIntConst(m)
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
