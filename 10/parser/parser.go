package parser

import (
	"bytes"
	. "compiler/tokenizer"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
)

type Parser struct {
	t    Tokenizer
	root Elem
}

func NewParser(t Tokenizer) *Parser {
	return &Parser{t, nil}
}

func (p *Parser) Parse() error {
	var err error
	p.root, err = p.compileClass()
	if err != nil {
		return err
	}
	return nil
}

var emptyXmlElem *regexp.Regexp = regexp.MustCompile(`( +)(<[a-zA-Z]+>)(</[a-zA-Z]+>)`)

// Change format of empty element to Nand2Tetris's one
// before: <expressionList></expressionList>
// after : <expressionList>
//         </expressionList>
func format(xmlStr string) string {
	return emptyXmlElem.ReplaceAllString(xmlStr, "$1$2\n$1$3")
}

func (p *Parser) XML() string {
	buf, _ := xml.MarshalIndent(p.root, "", "  ")
	xmlStr := format(string(buf))
	return xmlStr
}

const (
	TOKEN_ELEM  = "TOKEN"
	SYNTAX_ELEM = "SYNTAX"
)

type Elem interface {
	AddChild(c Elem)
	MarshalXML(enc *xml.Encoder, start xml.StartElement) error
	String() string
}

type BaseElem struct {
	elemName string
	children []Elem
}

func (e *BaseElem) AddChild(c Elem) {
	e.children = append(e.children, c)
}

func (e *BaseElem) String() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("elemName=%v\n", e.elemName))
	for _, child := range e.children {
		b.WriteString(child.String())
	}
	return b.String()
}

type TokenElem struct {
	*BaseElem
	token Token
}

func (e *TokenElem) String() string {
	return fmt.Sprintf("elemName=%v, tokenString=%v\n", e.elemName, e.token.String())
}

func (e *TokenElem) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	enc.EncodeElement(fmt.Sprintf(" %v ", e.token.String()), xml.StartElement{Name: xml.Name{Local: e.elemName}})
	return nil
}

type SyntaxElem struct {
	*BaseElem
}

func (e *SyntaxElem) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = e.elemName
	enc.EncodeToken(start)
	for _, child := range e.children {
		child.MarshalXML(enc, start)
	}
	enc.EncodeToken(start.End())
	return nil
}

func NewTokenElem(token Token) Elem {
	e := TokenElem{BaseElem: &BaseElem{elemName: token.Type(), children: make([]Elem, 0)}, token: token}
	return &e
}

func NewSyntaxElem(name string) Elem {
	e := SyntaxElem{BaseElem: &BaseElem{elemName: name, children: make([]Elem, 0)}}
	return &e
}

func (p *Parser) NewTokenElemCurrent() Elem {
	e := NewTokenElem(p.t.Current())
	return e
}

func compileError(err error, token Token) error {
	return fmt.Errorf("line=%v, column=%v: %v", token.Pos()[0], token.Pos()[1], err)
}

func (p *Parser) validateCurrent(tokenType string, tokenString string) error {
	if p.t.Current().Type() != tokenType || p.t.Current().String() != tokenString {
		return compileError(fmt.Errorf("want: type=%v, string=%v, got: type=%v, string=%v", tokenType, tokenString, p.t.Current().Type(), p.t.Current().String()), p.t.Current())
	}
	return nil
}

func (p *Parser) validateCurrentWithList(tokenType string, tokenStrings []string) error {
	for _, token := range tokenStrings {
		if p.t.Current().Type() == tokenType && p.t.Current().String() == token {
			return nil
		}
	}
	return compileError(fmt.Errorf("want: type=%v, string=%v, got: type=%v, string=%v", tokenType, tokenStrings, p.t.Current().Type(), p.t.Current().String()), p.t.Current())
}

func (p *Parser) validateCurrentType(tokenType string) error {
	if p.t.Current().Type() != tokenType {
		return compileError(fmt.Errorf("want: type=%v, got: type=%v, string=%v", tokenType, p.t.Current().Type(), p.t.Current().String()), p.t.Current())
	}
	return nil
}

func (p *Parser) validateCurrentIsTypeToken() error {
	if p.isCurrentTypeToken() {
		return nil
	} else {
		return compileError(fmt.Errorf("want: type token, got: type=%v string=%v", p.t.Current().Type(), p.t.Current().String()), p.t.Current())
	}
}

func (p *Parser) isCurrentEqualTo(tokenType string, tokenString string) bool {
	return p.t.Current().Type() == tokenType && p.t.Current().String() == tokenString
}

func (p *Parser) isCurrentStringEqualTo(tokenString string) bool {
	return p.t.Current().String() == tokenString
}

func (p *Parser) isCurrentTypeEqualTo(tokenType string) bool {
	return p.t.Current().String() == tokenType
}

func (p *Parser) isCurrentTypeToken() bool {
	curType := p.t.Current().Type()
	curStr := p.t.Current().String()
	if (curType == KEYWORD && (curStr == "int" || curStr == "char" || curStr == "boolean")) || curType == IDENTIFIER {
		return true
	} else {
		return false
	}
}

func isKeywordConstant(token Token) bool {
	if token.Type() != KEYWORD {
		return false
	}
	switch token.String() {
	case "true", "false", "null", "this":
		return true
	default:
		return false
	}
}

func isOp(token Token) bool {
	if token.Type() != SYMBOL {
		return false
	}
	switch token.String() {
	case "+", "-", "*", "/", "&", "|", "<", ">", "=":
		return true
	default:
		return false
	}
}

func isUnaryOp(token Token) bool {
	if token.Type() != SYMBOL {
		return false
	}
	switch token.String() {
	case "-", "~":
		return true
	default:
		return false
	}
}

func (p *Parser) compileClass() (Elem, error) {
	class := NewSyntaxElem("class")
	// class
	err := p.validateCurrent(KEYWORD, CLASS)
	if err != nil {
		return nil, fmt.Errorf("Invalid class declaration: %v", err)
	}
	class.AddChild(p.NewTokenElemCurrent())

	// className
	p.t.Advance()
	err = p.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, fmt.Errorf("Class name wasn't an identifier: %v", err)
	}
	class.AddChild(p.NewTokenElemCurrent())

	// {
	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, fmt.Errorf("Class declaration didn't start with {: %v", err)
	}
	class.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	for !p.isCurrentEqualTo(SYMBOL, "}") {

		curStr := p.t.Current().String()
		curType := p.t.Current().Type()

		// classVarDec
		if curType == KEYWORD && (curStr == "static" || curStr == "field") {
			varDec, err := p.compileClassVarDec()
			if err != nil {
				return nil, err
			}
			class.AddChild(varDec)

			// subroutine
		} else if curType == KEYWORD && (curStr == "constructor" || curStr == "function" || curStr == "method") {
			subroutine, err := p.compileSubroutine()
			if err != nil {
				return nil, err
			}
			class.AddChild(subroutine)
		} else {
			return nil, compileError(fmt.Errorf("Reached end of code"), p.t.Current())
		}
		p.t.Advance()
	}

	// }
	class.AddChild(p.NewTokenElemCurrent())
	return class, nil
}

func (p *Parser) compileClassVarDec() (Elem, error) {
	classVarDec := NewSyntaxElem("classVarDec")

	// static or field
	if !p.isCurrentEqualTo(KEYWORD, "static") && !p.isCurrentEqualTo(KEYWORD, "field") {
		return nil, compileError(errors.New("Invalid class var declaration."), p.t.Current())
	}
	classVarDec.AddChild(p.NewTokenElemCurrent())

	// type
	p.t.Advance()
	err := p.validateCurrentIsTypeToken()
	if err != nil {
		return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, p.t.Current()))
	}
	classVarDec.AddChild(p.NewTokenElemCurrent())
	for {
		// varName
		p.t.Advance()
		err = p.validateCurrentType(IDENTIFIER)
		if err != nil {
			return nil, fmt.Errorf("Invalid var name: %v ", err)
		}
		classVarDec.AddChild(p.NewTokenElemCurrent())

		next, err := p.t.LookAhead(1)
		if err != nil {
			return nil, compileError(err, p.t.Current())
		}
		if !(next.Type() == SYMBOL && next.String() == ",") {
			break
		}

		p.t.Advance()
		err = p.validateCurrent(SYMBOL, ",")
		if err != nil {
			return nil, err
		}
		classVarDec.AddChild(p.NewTokenElemCurrent())
	}
	// ;
	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, fmt.Errorf("Class var dec must end with ;: %v", compileError(err, p.t.Current()))
	}
	classVarDec.AddChild(p.NewTokenElemCurrent())

	return classVarDec, nil
}

func (p *Parser) compileSubroutine() (Elem, error) {
	subroutineDec := NewSyntaxElem("subroutineDec")

	// constructor, function, or method
	err := p.validateCurrentWithList(KEYWORD, []string{"constructor", "function", "method"})
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(p.NewTokenElemCurrent())

	// void or type name
	p.t.Advance()
	err1, err2 := p.validateCurrentType(KEYWORD), p.validateCurrentType(IDENTIFIER)
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", p.t.Current().Type())
	}
	subroutineDec.AddChild(p.NewTokenElemCurrent())

	// subroutineName
	p.t.Advance()
	err = p.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(p.NewTokenElemCurrent())

	// (
	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	parameterList, err := p.compileParameterList()
	if err != nil {
		return nil, err
	}
	subroutineDec.AddChild(parameterList)

	// )
	err = p.validateCurrent(SYMBOL, ")")
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(p.NewTokenElemCurrent())

	// subroutineBody
	p.t.Advance()
	subroutineBody, err := p.compileSubroutineBody()
	if err != nil {
		return nil, err
	}
	subroutineDec.AddChild(subroutineBody)
	return subroutineDec, nil
}

func (p *Parser) compileSubroutineBody() (Elem, error) {
	subroutineBody := NewSyntaxElem("subroutineBody")
	err := p.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, compileError(err, p.t.Current())
	}
	subroutineBody.AddChild(p.NewTokenElemCurrent())

	for {
		a, err := p.t.LookAhead(1)
		if a.Type() != KEYWORD || a.String() != "var" {
			break
		}

		p.t.Advance()
		varDec, err := p.compileVarDec()
		if err != nil {
			return nil, err
		}
		subroutineBody.AddChild(varDec)
	}

	p.t.Advance()
	statements, err := p.compileStatements()
	if err != nil {
		return nil, fmt.Errorf("Failed to compile statements: %v", err)
	}
	subroutineBody.AddChild(statements)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "}")
	if err != nil {
		return nil, compileError(err, p.t.Current())
	}
	subroutineBody.AddChild(p.NewTokenElemCurrent())

	return subroutineBody, nil
}

func (p *Parser) compileVarDec() (Elem, error) {
	varDec := NewSyntaxElem("varDec")

	// var
	err := p.validateCurrent(KEYWORD, "var")
	if err != nil {
		return nil, compileError(err, p.t.Current())
	}
	varDec.AddChild(p.NewTokenElemCurrent())

	// type
	p.t.Advance()
	err = p.validateCurrentIsTypeToken()
	if err != nil {
		return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, p.t.Current()))
	}
	varDec.AddChild(p.NewTokenElemCurrent())

	// varName
	p.t.Advance()
	err = p.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, p.t.Current()))
	}
	varDec.AddChild(p.NewTokenElemCurrent())

	for {
		next, err := p.t.LookAhead(1)
		if err != nil {
			return nil, compileError(err, p.t.Current())
		}
		if !(next.String() == ",") {
			break
		}
		// ,
		p.t.Advance()
		varDec.AddChild(p.NewTokenElemCurrent())

		// varName
		p.t.Advance()
		err = p.validateCurrentType(IDENTIFIER)
		if err != nil {
			return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, p.t.Current()))
		}
		varDec.AddChild(p.NewTokenElemCurrent())
	}
	// ;
	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, fmt.Errorf("Var dec must end with ;: %v", compileError(err, p.t.Current()))
	}
	varDec.AddChild(p.NewTokenElemCurrent())

	return varDec, nil
}

func (p *Parser) compileParameterList() (Elem, error) {
	parameterList := NewSyntaxElem("parameterList")

	if !p.isCurrentTypeToken() {
		return parameterList, nil
	}

	for {
		err := p.validateCurrentIsTypeToken()
		if err != nil {
			return nil, err
		}
		parameterList.AddChild(p.NewTokenElemCurrent())

		p.t.Advance()
		err = p.validateCurrentType(IDENTIFIER)
		if err != nil {
			return nil, err
		}
		parameterList.AddChild(p.NewTokenElemCurrent())

		aheadToken, err := p.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		if aheadToken.String() != "," {
			break
		}
		p.t.Advance()
		err = p.validateCurrent(SYMBOL, ",")
		if err != nil {
			return nil, err
		}
		parameterList.AddChild(p.NewTokenElemCurrent())
		p.t.Advance()
	}
	p.t.Advance()
	return parameterList, nil
}

func (p *Parser) compileStatements() (Elem, error) {
	statements := NewSyntaxElem("statements")
	contd := true

	if p.t.Current().Type() == SYMBOL && p.t.Current().String() == "}" {
		p.t.Backward()
		return statements, nil
	}

	switch p.t.Current().String() {
	case "let", "if", "while", "do", "return":
	default:
		// Empty statement
		return statements, nil
	}
	for contd {
		switch p.t.Current().String() {
		case "let":
			statement, err := p.compileLet()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile let statement: %v", err)
			}
			statements.AddChild(statement)

		case "if":
			statement, err := p.compileIf()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile if statement: %v", err)
			}
			statements.AddChild(statement)
		case "while":
			statement, err := p.compileWhile()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile while statement: %v", err)
			}
			statements.AddChild(statement)

		case "do":
			statement, err := p.compileDo()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile do statement: %v", err)
			}
			statements.AddChild(statement)
		case "return":
			statement, err := p.compileReturn()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile return statement: %v", err)
			}
			statements.AddChild(statement)
		}
		a, err := p.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		switch a.String() {
		case "let", "if", "while", "do", "return":
		default:
			contd = false
		}
		if contd == true {
			p.t.Advance()
		}
	}
	return statements, nil

}

func (p *Parser) compileLet() (Elem, error) {
	let := NewSyntaxElem("letStatement")
	err := p.validateCurrent(KEYWORD, "let")
	if err != nil {
		return nil, err
	}
	let.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	err = p.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, err
	}
	let.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	if p.isCurrentEqualTo(SYMBOL, "[") {
		let.AddChild(p.NewTokenElemCurrent())

		p.t.Advance()
		expression, err := p.compileExpression()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile expression in right side: %v", err)
		}
		let.AddChild(expression)

		p.t.Advance()
		err = p.validateCurrent(SYMBOL, "]")
		if err != nil {
			return nil, err
		}
		let.AddChild(p.NewTokenElemCurrent())
		p.t.Advance()
	}

	err = p.validateCurrent(SYMBOL, "=")
	if err != nil {
		return nil, err
	}
	let.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	expression, err := p.compileExpression()
	if err != nil {
		return nil, fmt.Errorf("Failed to compile expression in left side: %v", err)
	}
	let.AddChild(expression)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, err
	}
	let.AddChild(p.NewTokenElemCurrent())

	return let, nil
}

// Start: do
// End:   ;
func (p *Parser) compileDo() (Elem, error) {
	dost := NewSyntaxElem("doStatement")
	err := p.validateCurrent(KEYWORD, "do")
	if err != nil {
		return nil, err
	}
	dost.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	err = p.compileSubroutineCall(dost)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile subroutine call: %v", err)
	}

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, fmt.Errorf("Invalid ; %v", err)
	}
	dost.AddChild(p.NewTokenElemCurrent())
	return dost, nil
}

// Start: while
// End:   }
func (p *Parser) compileWhile() (Elem, error) {
	whilest := NewSyntaxElem("whileStatement")
	err := p.validateCurrent(KEYWORD, "while")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	expression, err := p.compileExpression()
	if err != nil {
		return nil, err
	}
	whilest.AddChild(expression)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ")")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	statement, err := p.compileStatements()
	if err != nil {
		return nil, err
	}
	whilest.AddChild(statement)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "}")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(p.NewTokenElemCurrent())

	return whilest, nil

}

func (p *Parser) compileReturn() (Elem, error) {
	returnst := NewSyntaxElem("returnStatement")

	err := p.validateCurrent(KEYWORD, "return")
	if err != nil {
		return nil, err
	}
	returnst.AddChild(p.NewTokenElemCurrent())

	a, err := p.t.LookAhead(1)
	if err != nil {
		return nil, err
	}
	if a.Type() == SYMBOL && a.String() == ";" {
		p.t.Advance()
		returnst.AddChild(p.NewTokenElemCurrent())
		return returnst, nil
	}
	p.t.Advance()
	expression, err := p.compileExpression()
	if err != nil {
		return nil, err
	}
	returnst.AddChild(expression)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, err
	}
	returnst.AddChild(p.NewTokenElemCurrent())
	return returnst, nil
}

// Start: if
// End    }
func (p *Parser) compileIf() (Elem, error) {
	ifst := NewSyntaxElem("ifStatement")
	err := p.validateCurrent(KEYWORD, "if")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	expression, err := p.compileExpression()
	if err != nil {
		return nil, err
	}
	ifst.AddChild(expression)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, ")")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	statements, err := p.compileStatements()
	if err != nil {
		return nil, err
	}
	ifst.AddChild(statements)

	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "}")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(p.NewTokenElemCurrent())

	a, err := p.t.LookAhead(1)
	if err != nil {
		return nil, err
	}

	if a.Type() == KEYWORD && a.String() == "else" {
		p.t.Advance()
		ifst.AddChild(p.NewTokenElemCurrent())

		p.t.Advance()
		err = p.validateCurrent(SYMBOL, "{")
		if err != nil {
			return nil, err
		}
		ifst.AddChild(p.NewTokenElemCurrent())

		p.t.Advance()
		statements, err := p.compileStatements()
		if err != nil {
			return nil, err
		}
		ifst.AddChild(statements)

		p.t.Advance()
		err = p.validateCurrent(SYMBOL, "}")
		if err != nil {
			return nil, fmt.Errorf("Failed to close } in compileIf %v: ", err)
		}
		ifst.AddChild(p.NewTokenElemCurrent())
	}
	return ifst, nil
}

func (p *Parser) compileExpression() (Elem, error) {
	expression := NewSyntaxElem("expression")
	for {
		term, err := p.compileTerm()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile term %v: ", err)
		}
		expression.AddChild(term)

		a, err := p.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		// op
		if !isOp(a) {
			break
		}
		p.t.Advance()
		expression.AddChild(p.NewTokenElemCurrent())

		p.t.Advance()
	}
	return expression, nil
}

// Start: (
// End:   )
func (p *Parser) compileExpressionList() (Elem, error) {
	expressionList := NewSyntaxElem("expressionList")
	if p.t.Current().Type() == SYMBOL && p.t.Current().String() == ")" {
		p.t.Backward()
		return expressionList, nil
	}
	for {
		expression, err := p.compileExpression()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile expression: %v", err)
		}
		expressionList.AddChild(expression)
		a, err := p.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		if !(a.Type() == SYMBOL && a.String() == ",") {
			break
		}
		// ,
		p.t.Advance()
		expressionList.AddChild(p.NewTokenElemCurrent())
		p.t.Advance()
	}
	return expressionList, nil
}

func (p *Parser) compileTerm() (Elem, error) {
	term := NewSyntaxElem("term")

	cur := p.t.Current()
	if cur.Type() == INT_CONST || cur.Type() == STR_CONST {
		// integerConstant or stringConstant
		term.AddChild(p.NewTokenElemCurrent())
	} else if isKeywordConstant(p.t.Current()) {
		// keywordConstant
		term.AddChild(p.NewTokenElemCurrent())
	} else if isUnaryOp(p.t.Current()) {
		// UnaryOp term
		term.AddChild(p.NewTokenElemCurrent())
		p.t.Advance()
		term2, err := p.compileTerm()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile UnaryOp term: %v", err)
		}
		term.AddChild(term2)
	} else if cur.Type() == SYMBOL && cur.String() == "(" {
		// ( expression )
		term.AddChild(p.NewTokenElemCurrent())

		p.t.Advance()
		expression, err := p.compileExpression()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile expression in '( expression )': %v", err)
		}
		term.AddChild(expression)

		p.t.Advance()
		err = p.validateCurrent(SYMBOL, ")")
		if err != nil {
			return nil, err
		}
		term.AddChild(p.NewTokenElemCurrent())

	} else if cur.Type() == IDENTIFIER {
		// subroutine call or array or var
		a, err := p.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		if a.Type() == SYMBOL && a.String() == "(" {
			// subroutineName ( expressionList )
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, "(")
			if err != nil {
				return nil, err
			}
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			expressionList, err := p.compileExpressionList()
			if err != nil {
				return nil, err
			}
			term.AddChild(expressionList)

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, ")")
			if err != nil {
				return nil, err
			}
		} else if a.Type() == SYMBOL && a.String() == "[" {
			// varName [ expression ]
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, "[")
			if err != nil {
				return nil, err
			}
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			expression, err := p.compileExpression()
			if err != nil {
				return nil, err
			}
			term.AddChild(expression)

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, "]")
			if err != nil {
				return nil, err
			}
			term.AddChild(p.NewTokenElemCurrent())
		} else if a.Type() == SYMBOL && a.String() == "." {
			// (className | varName).subroutineName(expressionList)
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, ".")
			if err != nil {
				return nil, err
			}
			term.AddChild(p.NewTokenElemCurrent())

			// subroutineName
			p.t.Advance()
			err = p.validateCurrentType(IDENTIFIER)
			if err != nil {
				return nil, err
			}
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, "(")
			if err != nil {
				return nil, err
			}
			term.AddChild(p.NewTokenElemCurrent())

			p.t.Advance()
			expressionList, err := p.compileExpressionList()
			term.AddChild(expressionList)
			if err != nil {
				return nil, fmt.Errorf("Failed to compile expression List in subroutine call: %v", err)
			}
			// }

			p.t.Advance()
			err = p.validateCurrent(SYMBOL, ")")
			if err != nil {
				return nil, fmt.Errorf("Failed to compile ) in subroutine call: %v", err)
			}
			term.AddChild(p.NewTokenElemCurrent())

		} else {
			// varName
			term.AddChild(p.NewTokenElemCurrent())
		}
	}
	return term, nil
}

// Start: subroutineName
// End:   )
func (p *Parser) compileSubroutineCall(e Elem) error {
	err := p.validateCurrentType(IDENTIFIER)
	if err != nil {
		return err
	}
	e.AddChild(p.NewTokenElemCurrent())

	a1, err := p.t.LookAhead(1)
	if err != nil {
		return err
	}

	if a1.String() == "." {
		p.t.Advance()
		err = p.validateCurrent(SYMBOL, ".")
		if err != nil {
			return err
		}
		e.AddChild(p.NewTokenElemCurrent())
		p.t.Advance()
		err = p.validateCurrentType(IDENTIFIER)
		if err != nil {
			return err
		}
		e.AddChild(p.NewTokenElemCurrent())
	}
	p.t.Advance()
	err = p.validateCurrent(SYMBOL, "(")
	if err != nil {
		return err
	}
	e.AddChild(p.NewTokenElemCurrent())

	p.t.Advance()
	expressionList, err := p.compileExpressionList()
	if err != nil {
		return err
	}
	e.AddChild(expressionList)
	p.t.Advance()
	// }
	err = p.validateCurrent(SYMBOL, ")")
	if err != nil {
		return err
	}
	e.AddChild(p.NewTokenElemCurrent())
	return nil
}
