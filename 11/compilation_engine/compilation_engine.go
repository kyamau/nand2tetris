package compilation_engine

import (
	"bytes"
	. "compiler/symbol_table"
	. "compiler/tokenizer"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
)

type CompilationEngine struct {
	t      Tokenizer
	root   Elem
	tables []*SymbolTable // tables[0] for class, tables[1] for method
}

func NewCompilationEngine(t Tokenizer) *CompilationEngine {
	return &CompilationEngine{t, nil, make([]*SymbolTable, 2)}
}

func (ce *CompilationEngine) Compile() error {
	var err error
	ce.root, err = ce.compileClass()
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) ClassTable() *SymbolTable {
	return ce.tables[0]
}

func (ce *CompilationEngine) MethodTable() *SymbolTable {
	return ce.tables[1]
}

var emptyXmlElem *regexp.Regexp = regexp.MustCompile(`( +)(<[a-zA-Z]+>)(</[a-zA-Z]+>)`)

// Change format of empty element to Nand2Tetris's one
// before: <expressionList></expressionList>
// after : <expressionList>
//         </expressionList>
func format(xmlStr string) string {
	return emptyXmlElem.ReplaceAllString(xmlStr, "$1$2\n$1$3")
}

func (ce *CompilationEngine) XML() string {
	buf, _ := xml.MarshalIndent(ce.root, "", "  ")
	xmlStr := format(string(buf)) + "\n"
	return xmlStr
}

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

func (ce *CompilationEngine) NewTokenElemCurrent() Elem {
	e := NewTokenElem(ce.t.Current())
	return e
}

func compileError(err error, token Token) error {
	return fmt.Errorf("line=%v, column=%v: %v", token.Pos()[0], token.Pos()[1], err)
}

func (ce *CompilationEngine) validateCurrent(tokenType string, tokenString string) error {
	if ce.t.Current().Type() != tokenType || ce.t.Current().String() != tokenString {
		return compileError(fmt.Errorf("want: type=%v, string=%v, got: type=%v, string=%v", tokenType, tokenString, ce.t.Current().Type(), ce.t.Current().String()), ce.t.Current())
	}
	return nil
}

func (ce *CompilationEngine) validateCurrentWithList(tokenType string, tokenStrings []string) error {
	for _, token := range tokenStrings {
		if ce.t.Current().Type() == tokenType && ce.t.Current().String() == token {
			return nil
		}
	}
	return compileError(fmt.Errorf("want: type=%v, string=%v, got: type=%v, string=%v", tokenType, tokenStrings, ce.t.Current().Type(), ce.t.Current().String()), ce.t.Current())
}

func (ce *CompilationEngine) validateCurrentType(tokenType string) error {
	if ce.t.Current().Type() != tokenType {
		return compileError(fmt.Errorf("want: type=%v, got: type=%v, string=%v", tokenType, ce.t.Current().Type(), ce.t.Current().String()), ce.t.Current())
	}
	return nil
}

func (ce *CompilationEngine) validateCurrentIsTypeToken() error {
	if ce.isCurrentTypeToken() {
		return nil
	} else {
		return compileError(fmt.Errorf("want: type token, got: type=%v string=%v", ce.t.Current().Type(), ce.t.Current().String()), ce.t.Current())
	}
}

func (ce *CompilationEngine) isCurrentEqualTo(tokenType string, tokenString string) bool {
	return ce.t.Current().Type() == tokenType && ce.t.Current().String() == tokenString
}

func (ce *CompilationEngine) isCurrentStringEqualTo(tokenString string) bool {
	return ce.t.Current().String() == tokenString
}

func (ce *CompilationEngine) isCurrentTypeEqualTo(tokenType string) bool {
	return ce.t.Current().String() == tokenType
}

func (ce *CompilationEngine) isCurrentTypeToken() bool {
	curType := ce.t.Current().Type()
	curStr := ce.t.Current().String()
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

func (ce *CompilationEngine) compileClass() (Elem, error) {
	class := NewSyntaxElem("class")
	// class
	err := ce.validateCurrent(KEYWORD, CLASS)
	if err != nil {
		return nil, fmt.Errorf("Invalid class declaration: %v", err)
	}
	class.AddChild(ce.NewTokenElemCurrent())

	// className
	ce.t.Advance()
	err = ce.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, fmt.Errorf("Class name wasn't an identifier: %v", err)
	}
	// Add symbol table for class
	ce.tables[0] = NewEmptyClassSymbolTable(ce.t.Current().String())
	class.AddChild(ce.NewTokenElemCurrent())

	// {
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, fmt.Errorf("Class declaration didn't start with {: %v", err)
	}
	class.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	for !ce.isCurrentEqualTo(SYMBOL, "}") {

		curStr := ce.t.Current().String()
		curType := ce.t.Current().Type()

		// classVarDec
		if curType == KEYWORD && (curStr == "static" || curStr == "field") {
			varDec, err := ce.compileClassVarDec()
			if err != nil {
				return nil, err
			}
			class.AddChild(varDec)

			// subroutine
		} else if curType == KEYWORD && (curStr == "constructor" || curStr == "function" || curStr == "method") {
			subroutine, err := ce.compileSubroutine()
			if err != nil {
				return nil, err
			}
			class.AddChild(subroutine)
		} else {
			return nil, compileError(fmt.Errorf("Reached end of code"), ce.t.Current())
		}
		ce.t.Advance()
	}

	// }
	class.AddChild(ce.NewTokenElemCurrent())
	return class, nil
}

func (ce *CompilationEngine) compileClassVarDec() (Elem, error) {
	classVarDec := NewSyntaxElem("classVarDec")

	var varType string
	var varKind string

	// static or field
	isStatic := ce.isCurrentEqualTo(KEYWORD, "static")
	isField := ce.isCurrentEqualTo(KEYWORD, "field")
	if !isStatic && !isField {
		return nil, compileError(errors.New("Invalid class var declaration."), ce.t.Current())
	} else if isStatic {
		varKind = STATIC
	} else {
		varKind = FIELD
	}
	classVarDec.AddChild(ce.NewTokenElemCurrent())

	// type
	ce.t.Advance()
	err := ce.validateCurrentIsTypeToken()
	if err != nil {
		return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, ce.t.Current()))
	}
	varType = ce.t.Current().String()
	classVarDec.AddChild(ce.NewTokenElemCurrent())

	for {
		// varName
		ce.t.Advance()
		err = ce.validateCurrentType(IDENTIFIER)
		if err != nil {
			return nil, fmt.Errorf("Invalid var name: %v ", err)
		}
		// Add entry for the last symbol table
		ce.ClassTable().Define(ce.t.Current().String(), varType, varKind)
		classVarDec.AddChild(ce.NewTokenElemCurrent())

		next, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, compileError(err, ce.t.Current())
		}
		if !(next.Type() == SYMBOL && next.String() == ",") {
			break
		}

		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, ",")
		if err != nil {
			return nil, err
		}
		classVarDec.AddChild(ce.NewTokenElemCurrent())
	}
	// ;
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, fmt.Errorf("Class var dec must end with ;: %v", compileError(err, ce.t.Current()))
	}
	classVarDec.AddChild(ce.NewTokenElemCurrent())

	return classVarDec, nil
}

func (ce *CompilationEngine) compileSubroutine() (Elem, error) {
	subroutineDec := NewSyntaxElem("subroutineDec")

	// constructor, function, or method
	err := ce.validateCurrentWithList(KEYWORD, []string{"constructor", "function", "method"})
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(ce.NewTokenElemCurrent())

	// void or type name
	ce.t.Advance()
	err1, err2 := ce.validateCurrentType(KEYWORD), ce.validateCurrentType(IDENTIFIER)
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", ce.t.Current().Type())
	}
	subroutineDec.AddChild(ce.NewTokenElemCurrent())

	// subroutineName
	ce.t.Advance()
	err = ce.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	symbolTableName := ce.t.Current().String()
	subroutineDec.AddChild(ce.NewTokenElemCurrent())

	// (
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(ce.NewTokenElemCurrent())

	// Add symbol table for this method
	ce.tables[1] = NewEmptyMethodSymbolTable(symbolTableName)

	ce.t.Advance()
	parameterList, err := ce.compileParameterList()
	if err != nil {
		return nil, err
	}
	subroutineDec.AddChild(parameterList)

	// )
	err = ce.validateCurrent(SYMBOL, ")")
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(ce.NewTokenElemCurrent())

	// subroutineBody
	ce.t.Advance()
	subroutineBody, err := ce.compileSubroutineBody()
	if err != nil {
		return nil, err
	}
	subroutineDec.AddChild(subroutineBody)
	return subroutineDec, nil
}

func (ce *CompilationEngine) compileSubroutineBody() (Elem, error) {
	subroutineBody := NewSyntaxElem("subroutineBody")
	err := ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, compileError(err, ce.t.Current())
	}
	subroutineBody.AddChild(ce.NewTokenElemCurrent())

	for {
		a, err := ce.t.LookAhead(1)
		if a.Type() != KEYWORD || a.String() != "var" {
			break
		}

		ce.t.Advance()
		varDec, err := ce.compileVarDec()
		if err != nil {
			return nil, err
		}
		subroutineBody.AddChild(varDec)
	}

	ce.t.Advance()
	statements, err := ce.compileStatements()
	if err != nil {
		return nil, fmt.Errorf("Failed to compile statements: %v", err)
	}
	subroutineBody.AddChild(statements)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "}")
	if err != nil {
		return nil, compileError(err, ce.t.Current())
	}
	subroutineBody.AddChild(ce.NewTokenElemCurrent())

	return subroutineBody, nil
}

func (ce *CompilationEngine) compileVarDec() (Elem, error) {
	varDec := NewSyntaxElem("varDec")

	// var
	err := ce.validateCurrent(KEYWORD, "var")
	if err != nil {
		return nil, compileError(err, ce.t.Current())
	}
	varDec.AddChild(ce.NewTokenElemCurrent())

	// type
	ce.t.Advance()
	err = ce.validateCurrentIsTypeToken()
	if err != nil {
		return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, ce.t.Current()))
	}
	varType := ce.t.Current().String()
	varDec.AddChild(ce.NewTokenElemCurrent())

	// varName
	ce.t.Advance()
	varName := ce.t.Current().String()
	err = ce.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, ce.t.Current()))
	}

	// Add var to symbol table
	ce.MethodTable().Define(varName, varType, "var")
	varDec.AddChild(ce.NewTokenElemCurrent())

	for {
		next, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, compileError(err, ce.t.Current())
		}
		if !(next.String() == ",") {
			break
		}
		// ,
		ce.t.Advance()
		varDec.AddChild(ce.NewTokenElemCurrent())

		// varName
		ce.t.Advance()
		err = ce.validateCurrentType(IDENTIFIER)
		if err != nil {
			return nil, fmt.Errorf("Invalid type declaration: %v", compileError(err, ce.t.Current()))
		}
		varDec.AddChild(ce.NewTokenElemCurrent())
	}
	// ;
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, fmt.Errorf("Var dec must end with ;: %v", compileError(err, ce.t.Current()))
	}
	varDec.AddChild(ce.NewTokenElemCurrent())

	return varDec, nil
}

func (ce *CompilationEngine) compileParameterList() (Elem, error) {
	parameterList := NewSyntaxElem("parameterList")

	if !ce.isCurrentTypeToken() {
		return parameterList, nil
	}

	for {
		err := ce.validateCurrentIsTypeToken()
		if err != nil {
			return nil, err
		}
		varType := ce.t.Current().String()
		varKind := "argument"
		parameterList.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
		err = ce.validateCurrentType(IDENTIFIER)
		if err != nil {
			return nil, err
		}
		varName := ce.t.Current().String()
		// Add arguments to the symbol table
		ce.MethodTable().Define(varName, varType, varKind)
		parameterList.AddChild(ce.NewTokenElemCurrent())

		aheadToken, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		if aheadToken.String() != "," {
			break
		}
		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, ",")
		if err != nil {
			return nil, err
		}
		parameterList.AddChild(ce.NewTokenElemCurrent())
		ce.t.Advance()
	}
	ce.t.Advance()
	return parameterList, nil
}

func (ce *CompilationEngine) compileStatements() (Elem, error) {
	statements := NewSyntaxElem("statements")
	contd := true

	if ce.t.Current().Type() == SYMBOL && ce.t.Current().String() == "}" {
		ce.t.Backward()
		return statements, nil
	}

	switch ce.t.Current().String() {
	case "let", "if", "while", "do", "return":
	default:
		// Empty statement
		return statements, nil
	}
	for contd {
		switch ce.t.Current().String() {
		case "let":
			statement, err := ce.compileLet()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile let statement: %v", err)
			}
			statements.AddChild(statement)

		case "if":
			statement, err := ce.compileIf()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile if statement: %v", err)
			}
			statements.AddChild(statement)
		case "while":
			statement, err := ce.compileWhile()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile while statement: %v", err)
			}
			statements.AddChild(statement)

		case "do":
			statement, err := ce.compileDo()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile do statement: %v", err)
			}
			statements.AddChild(statement)
		case "return":
			statement, err := ce.compileReturn()
			if err != nil {
				return nil, fmt.Errorf("Faile to compile return statement: %v", err)
			}
			statements.AddChild(statement)
		}
		a, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		switch a.String() {
		case "let", "if", "while", "do", "return":
		default:
			contd = false
		}
		if contd == true {
			ce.t.Advance()
		}
	}
	return statements, nil

}

func (ce *CompilationEngine) compileLet() (Elem, error) {
	let := NewSyntaxElem("letStatement")
	err := ce.validateCurrent(KEYWORD, "let")
	if err != nil {
		return nil, err
	}
	let.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	err = ce.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, err
	}
	let.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	if ce.isCurrentEqualTo(SYMBOL, "[") {
		let.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
		expression, err := ce.compileExpression()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile expression in right side: %v", err)
		}
		let.AddChild(expression)

		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, "]")
		if err != nil {
			return nil, err
		}
		let.AddChild(ce.NewTokenElemCurrent())
		ce.t.Advance()
	}

	err = ce.validateCurrent(SYMBOL, "=")
	if err != nil {
		return nil, err
	}
	let.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	expression, err := ce.compileExpression()
	if err != nil {
		return nil, fmt.Errorf("Failed to compile expression in left side: %v", err)
	}
	let.AddChild(expression)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, err
	}
	let.AddChild(ce.NewTokenElemCurrent())

	return let, nil
}

// Start: do
// End:   ;
func (ce *CompilationEngine) compileDo() (Elem, error) {
	dost := NewSyntaxElem("doStatement")
	err := ce.validateCurrent(KEYWORD, "do")
	if err != nil {
		return nil, err
	}
	dost.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	err = ce.compileSubroutineCall(dost)
	if err != nil {
		return nil, fmt.Errorf("Failed to compile subroutine call: %v", err)
	}

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, fmt.Errorf("Invalid ; %v", err)
	}
	dost.AddChild(ce.NewTokenElemCurrent())
	return dost, nil
}

// Start: while
// End:   }
func (ce *CompilationEngine) compileWhile() (Elem, error) {
	whilest := NewSyntaxElem("whileStatement")
	err := ce.validateCurrent(KEYWORD, "while")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	expression, err := ce.compileExpression()
	if err != nil {
		return nil, err
	}
	whilest.AddChild(expression)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ")")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	statement, err := ce.compileStatements()
	if err != nil {
		return nil, err
	}
	whilest.AddChild(statement)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "}")
	if err != nil {
		return nil, err
	}
	whilest.AddChild(ce.NewTokenElemCurrent())

	return whilest, nil

}

func (ce *CompilationEngine) compileReturn() (Elem, error) {
	returnst := NewSyntaxElem("returnStatement")

	err := ce.validateCurrent(KEYWORD, "return")
	if err != nil {
		return nil, err
	}
	returnst.AddChild(ce.NewTokenElemCurrent())

	a, err := ce.t.LookAhead(1)
	if err != nil {
		return nil, err
	}
	if a.Type() == SYMBOL && a.String() == ";" {
		ce.t.Advance()
		returnst.AddChild(ce.NewTokenElemCurrent())
		return returnst, nil
	}
	ce.t.Advance()
	expression, err := ce.compileExpression()
	if err != nil {
		return nil, err
	}
	returnst.AddChild(expression)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, err
	}
	returnst.AddChild(ce.NewTokenElemCurrent())
	return returnst, nil
}

// Start: if
// End    }
func (ce *CompilationEngine) compileIf() (Elem, error) {
	ifst := NewSyntaxElem("ifStatement")
	err := ce.validateCurrent(KEYWORD, "if")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	expression, err := ce.compileExpression()
	if err != nil {
		return nil, err
	}
	ifst.AddChild(expression)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ")")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	statements, err := ce.compileStatements()
	if err != nil {
		return nil, err
	}
	ifst.AddChild(statements)

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "}")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(ce.NewTokenElemCurrent())

	a, err := ce.t.LookAhead(1)
	if err != nil {
		return nil, err
	}

	if a.Type() == KEYWORD && a.String() == "else" {
		ce.t.Advance()
		ifst.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, "{")
		if err != nil {
			return nil, err
		}
		ifst.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
		statements, err := ce.compileStatements()
		if err != nil {
			return nil, err
		}
		ifst.AddChild(statements)

		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, "}")
		if err != nil {
			return nil, fmt.Errorf("Failed to close } in compileIf %v: ", err)
		}
		ifst.AddChild(ce.NewTokenElemCurrent())
	}
	return ifst, nil
}

func (ce *CompilationEngine) compileExpression() (Elem, error) {
	expression := NewSyntaxElem("expression")
	for {
		term, err := ce.compileTerm()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile term %v: ", err)
		}
		expression.AddChild(term)

		a, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		// op
		if !isOp(a) {
			break
		}
		ce.t.Advance()
		expression.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
	}
	return expression, nil
}

// Start: (
// End:   )
func (ce *CompilationEngine) compileExpressionList() (Elem, error) {
	expressionList := NewSyntaxElem("expressionList")
	if ce.t.Current().Type() == SYMBOL && ce.t.Current().String() == ")" {
		ce.t.Backward()
		return expressionList, nil
	}
	for {
		expression, err := ce.compileExpression()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile expression: %v", err)
		}
		expressionList.AddChild(expression)
		a, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		if !(a.Type() == SYMBOL && a.String() == ",") {
			break
		}
		// ,
		ce.t.Advance()
		expressionList.AddChild(ce.NewTokenElemCurrent())
		ce.t.Advance()
	}
	return expressionList, nil
}

func (ce *CompilationEngine) compileTerm() (Elem, error) {
	term := NewSyntaxElem("term")

	cur := ce.t.Current()
	if cur.Type() == INT_CONST || cur.Type() == STR_CONST {
		// integerConstant or stringConstant
		term.AddChild(ce.NewTokenElemCurrent())
	} else if isKeywordConstant(ce.t.Current()) {
		// keywordConstant
		term.AddChild(ce.NewTokenElemCurrent())
	} else if isUnaryOp(ce.t.Current()) {
		// UnaryOp term
		term.AddChild(ce.NewTokenElemCurrent())
		ce.t.Advance()
		term2, err := ce.compileTerm()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile UnaryOp term: %v", err)
		}
		term.AddChild(term2)
	} else if cur.Type() == SYMBOL && cur.String() == "(" {
		// ( expression )
		term.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
		expression, err := ce.compileExpression()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile expression in '( expression )': %v", err)
		}
		term.AddChild(expression)

		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, ")")
		if err != nil {
			return nil, err
		}
		term.AddChild(ce.NewTokenElemCurrent())

	} else if cur.Type() == IDENTIFIER {
		// subroutine call or array or var
		a, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, err
		}
		if a.Type() == SYMBOL && a.String() == "(" {
			// subroutineName ( expressionList )
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, "(")
			if err != nil {
				return nil, err
			}
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			expressionList, err := ce.compileExpressionList()
			if err != nil {
				return nil, err
			}
			term.AddChild(expressionList)

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, ")")
			if err != nil {
				return nil, err
			}
		} else if a.Type() == SYMBOL && a.String() == "[" {
			// varName [ expression ]
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, "[")
			if err != nil {
				return nil, err
			}
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			expression, err := ce.compileExpression()
			if err != nil {
				return nil, err
			}
			term.AddChild(expression)

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, "]")
			if err != nil {
				return nil, err
			}
			term.AddChild(ce.NewTokenElemCurrent())
		} else if a.Type() == SYMBOL && a.String() == "." {
			// (className | varName).subroutineName(expressionList)
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, ".")
			if err != nil {
				return nil, err
			}
			term.AddChild(ce.NewTokenElemCurrent())

			// subroutineName
			ce.t.Advance()
			err = ce.validateCurrentType(IDENTIFIER)
			if err != nil {
				return nil, err
			}
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, "(")
			if err != nil {
				return nil, err
			}
			term.AddChild(ce.NewTokenElemCurrent())

			ce.t.Advance()
			expressionList, err := ce.compileExpressionList()
			term.AddChild(expressionList)
			if err != nil {
				return nil, fmt.Errorf("Failed to compile expression List in subroutine call: %v", err)
			}
			// }

			ce.t.Advance()
			err = ce.validateCurrent(SYMBOL, ")")
			if err != nil {
				return nil, fmt.Errorf("Failed to compile ) in subroutine call: %v", err)
			}
			term.AddChild(ce.NewTokenElemCurrent())

		} else {
			// varName
			term.AddChild(ce.NewTokenElemCurrent())
		}
	}
	return term, nil
}

// Start: subroutineName
// End:   )
func (ce *CompilationEngine) compileSubroutineCall(e Elem) error {
	err := ce.validateCurrentType(IDENTIFIER)
	if err != nil {
		return err
	}
	e.AddChild(ce.NewTokenElemCurrent())

	a1, err := ce.t.LookAhead(1)
	if err != nil {
		return err
	}

	if a1.String() == "." {
		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, ".")
		if err != nil {
			return err
		}
		e.AddChild(ce.NewTokenElemCurrent())
		ce.t.Advance()
		err = ce.validateCurrentType(IDENTIFIER)
		if err != nil {
			return err
		}
		e.AddChild(ce.NewTokenElemCurrent())
	}
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "(")
	if err != nil {
		return err
	}
	e.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	expressionList, err := ce.compileExpressionList()
	if err != nil {
		return err
	}
	e.AddChild(expressionList)
	ce.t.Advance()
	// }
	err = ce.validateCurrent(SYMBOL, ")")
	if err != nil {
		return err
	}
	e.AddChild(ce.NewTokenElemCurrent())
	return nil
}
