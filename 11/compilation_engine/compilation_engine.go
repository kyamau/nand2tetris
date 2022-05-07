package compilation_engine

import (
	. "compiler/symbol_table"
	. "compiler/tokenizer"
	. "compiler/vmwriter"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type CompilationEngine struct {
	t             *Tokenizer
	root          Elem
	tables        []*SymbolTable // tables[0] for class, tables[1] for subroutine. tables[1] will be cleared at every subroutine declaration.
	labelManager  *LabelManager  // lable manager. It will be cleared at every subroutine declaration
	operatorStack *Stack
	vmwriter      *VMWriter
}

func NewCompilationEngine(t *Tokenizer, vmWriter *VMWriter) *CompilationEngine {
	return &CompilationEngine{t, nil, make([]*SymbolTable, 2), NewLabelManager(), NewStack(), vmWriter}
}

func (ce *CompilationEngine) Compile() error {
	var err error
	ce.root, err = ce.compileClass()
	if err != nil {
		return err
	}
	return nil
}

func (ce *CompilationEngine) classTable() *SymbolTable {
	return ce.tables[0]
}

func (ce *CompilationEngine) subroutineTable() *SymbolTable {
	return ce.tables[1]
}

func (ce *CompilationEngine) initializeSubroutineTable(subroutineName string) {
	ce.tables[1] = NewSymbolTable(subroutineName)
}

func (ce *CompilationEngine) initializeLabelManager() {
	ce.labelManager = NewLabelManager()
}

func (ce *CompilationEngine) resolveVariableInSubroutine(varName string) (*SymbolTable, bool) {
	if _, ok := ce.subroutineTable().KindOf(varName); ok {
		return ce.subroutineTable(), ok
	}
	_, ok := ce.classTable().KindOf(varName)
	return ce.classTable(), ok
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

func (ce *CompilationEngine) WriteCode(filepath string) error {
	err := WriteCode(ce.vmwriter.Code(), filepath)
	if err != nil {
		return fmt.Errorf("Failed to write VM code: %v", err)
	}
	return nil
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

func (ce *CompilationEngine) isCurrent(tokenType string, tokenString string) bool {
	return ce.t.Current().Type() == tokenType && ce.t.Current().String() == tokenString
}

func (ce *CompilationEngine) isCurrentString(tokenString string) bool {
	return ce.t.Current().String() == tokenString
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
	class.AddChild(ce.NewTokenElemCurrent())
	// Add symbol table for class
	ce.tables[0] = NewSymbolTable(ce.t.Current().String())

	// {
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, fmt.Errorf("Class declaration didn't start with {: %v", err)
	}
	class.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	for !ce.isCurrent(SYMBOL, "}") {

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
	isStatic := ce.isCurrent(KEYWORD, "static")
	isField := ce.isCurrent(KEYWORD, "field")
	if isStatic {
		varKind = "static"
	} else if isField {
		varKind = "field"
	} else {
		return nil, compileError(errors.New("Invalid class var declaration."), ce.t.Current())
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
		classVarDec.AddChild(ce.NewTokenElemCurrent())
		// Add entry for the last symbol table
		ce.classTable().Define(ce.t.Current().String(), varType, varKind)

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
	subroutineType := ce.t.Current().String()

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
	subroutineDec.AddChild(ce.NewTokenElemCurrent())
	subroutineName := ce.classTable().Name() + "." + ce.t.Current().String()
	symbolTableName := subroutineName

	// Initialize symbol table and label manager for this subroutine
	ce.initializeSubroutineTable(symbolTableName)
	ce.initializeLabelManager()

	// (
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "(")
	if err != nil {
		return nil, fmt.Errorf("Invalid subroutine declaration: %v", err)
	}
	subroutineDec.AddChild(ce.NewTokenElemCurrent())

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
	subroutineBody, err := ce.compileSubroutineBody(subroutineType, subroutineName)
	if err != nil {
		return nil, err
	}
	subroutineDec.AddChild(subroutineBody)

	return subroutineDec, nil
}

func (ce *CompilationEngine) compileSubroutineBody(subroutineType string, subroutineName string) (Elem, error) {
	subroutineBody := NewSyntaxElem("subroutineBody")
	err := ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, compileError(err, ce.t.Current())
	}
	subroutineBody.AddChild(ce.NewTokenElemCurrent())

	nLocals := 0
	for {
		a, err := ce.t.LookAhead(1)
		if a.Type() != KEYWORD || a.String() != "var" {
			break
		}

		ce.t.Advance()
		varDec, nLocalsInVarDec, err := ce.compileVarDec()
		if err != nil {
			return nil, err
		}
		subroutineBody.AddChild(varDec)
		nLocals += nLocalsInVarDec
	}

	ce.vmwriter.Add(FunctionCode(subroutineName, nLocals))

	// Allocate a new object
	if subroutineType == "constructor" {
		nFields := ce.classTable().VarCount("field")
		ce.vmwriter.Add(PushCode("constant", nFields))
		ce.vmwriter.Add(CallCode("Memory.alloc", 1))
		// Set the pointer of the instance to the pointer segment.
		// This will be passed to the left side variable in let statement.
		ce.vmwriter.Add(PopCode("pointer", 0))
	}

	// Set the instance(passed as 1st parameter implicitly) to the pointer
	if subroutineType == "method" {
		ce.vmwriter.Add(PushCode("argument", 0))
		ce.vmwriter.Add(PopCode("pointer", 0))
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

// It also returns the number of local variables in the declaration
func (ce *CompilationEngine) compileVarDec() (Elem, int, error) {
	varDec := NewSyntaxElem("varDec")

	// var
	err := ce.validateCurrent(KEYWORD, "var")
	if err != nil {
		return nil, -1, compileError(err, ce.t.Current())
	}
	varDec.AddChild(ce.NewTokenElemCurrent())

	// type
	ce.t.Advance()
	err = ce.validateCurrentIsTypeToken()
	if err != nil {
		return nil, -1, fmt.Errorf("Invalid type declaration: %v", compileError(err, ce.t.Current()))
	}
	varDec.AddChild(ce.NewTokenElemCurrent())
	varType := ce.t.Current().String()

	// varName
	ce.t.Advance()
	err = ce.validateCurrentType(IDENTIFIER)
	if err != nil {
		return nil, -1, fmt.Errorf("Invalid var declaration: %v", compileError(err, ce.t.Current()))
	}
	varDec.AddChild(ce.NewTokenElemCurrent())

	// Add var to symbol table
	varName := ce.t.Current().String()
	ce.subroutineTable().Define(varName, varType, "var")

	nLocals := 1

	// If there are more variables, continue reading
	for {
		next, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, -1, compileError(err, ce.t.Current())
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
			return nil, -1, fmt.Errorf("Invalid type declaration: %v", compileError(err, ce.t.Current()))
		}
		varDec.AddChild(ce.NewTokenElemCurrent())

		// Add var to symbol table
		varName = ce.t.Current().String()
		ce.subroutineTable().Define(varName, varType, "var")

		nLocals++
	}
	// ;
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, ";")
	if err != nil {
		return nil, -1, fmt.Errorf("Var dec must end with ;: %v", compileError(err, ce.t.Current()))
	}
	varDec.AddChild(ce.NewTokenElemCurrent())

	return varDec, nLocals, nil
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
		ce.subroutineTable().Define(varName, varType, varKind)
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
	varName := ce.t.Current().String()

	var table *SymbolTable
	var ok bool
	if table, ok = ce.resolveVariableInSubroutine(varName); ok {
	} else {
		return nil, fmt.Errorf("Variable %s is not defined.", varName)
	}
	varIndex, _ := table.IndexOf(varName)
	varKind, _ := table.KindOf(varName)
	isArray := false

	ce.t.Advance()
	// Array
	if ce.isCurrent(SYMBOL, "[") {
		isArray = true
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
		// array head + index
		ce.vmwriter.Add(PushCode(varKindToSegment(varKind), varIndex))
		ce.vmwriter.Add("add")
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

	if isArray {
		ce.vmwriter.Add(PopCode("temp", 0))    // Store the result of the left expression
		ce.vmwriter.Add(PopCode("pointer", 1)) // Store the pointer to the array head + index
		ce.vmwriter.Add(PushCode("temp", 0))
		ce.vmwriter.Add(PopCode("that", 0))
	} else { //{ varType == "int" || varType == "boolean" || varType == "char" {
		ce.vmwriter.Add(PopCode(varKindToSegment(varKind), varIndex))
		//TODO add other kind
	}
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

	// Pop the return value and discard it.
	// See the text p263 and chapter 11 slide p62
	ce.vmwriter.Add(PopCode("temp", 0))
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

	// Set a label for the starting point of the while loop
	ce.labelManager.StartWhile()
	ce.vmwriter.Add(LabelCode(ce.labelManager.WhileExpLabel()))

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

	// If the condition isn't met, jump to the end of the while.
	ce.vmwriter.Add("not")
	ce.vmwriter.Add(IfGotoCode(ce.labelManager.WhileEndLabel()))

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

	// Back to the head of the while.
	ce.vmwriter.Add(GotoCode(ce.labelManager.WhileExpLabel()))
	ce.vmwriter.Add(LabelCode(ce.labelManager.WhileEndLabel()))

	ce.labelManager.EndWhile()
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
		// Push 0 for void return
		// See p263.
		ce.vmwriter.Add(PushCode("constant", 0))
		ce.vmwriter.Add(ReturnCode())
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

	ce.vmwriter.Add(ReturnCode())
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

	ce.labelManager.StartIf()
	ce.vmwriter.Add(IfGotoCode(ce.labelManager.IfTrueLabel()))
	ce.vmwriter.Add(GotoCode(ce.labelManager.IfFalseLabel()))

	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "{")
	if err != nil {
		return nil, err
	}
	ifst.AddChild(ce.NewTokenElemCurrent())

	ce.vmwriter.Add(LabelCode(ce.labelManager.IfTrueLabel()))

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
		// Skip the else clause when the if condition met. This is needed only when else clause exists.
		ce.vmwriter.Add(GotoCode(ce.labelManager.IfEndLabel()))
		// Jumped here from the if clause when
		// - the condition didn't meet
		// - else clause exists
		ce.vmwriter.Add(LabelCode(ce.labelManager.IfFalseLabel()))

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

		ce.vmwriter.Add(LabelCode(ce.labelManager.IfEndLabel()))
	} else {
		// Jumped here from the if clause when
		// - the condition didn't meet
		// - else clause doesn't exists
		ce.vmwriter.Add(LabelCode(ce.labelManager.IfFalseLabel()))
	}
	ce.labelManager.EndIf()
	return ifst, nil
}

func (ce *CompilationEngine) compileExpression() (Elem, error) {
	expression := NewSyntaxElem("expression")

	nOps := 0
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
		nOps++

		ce.t.Advance()
		expression.AddChild(ce.NewTokenElemCurrent())
		op := ce.t.Current().String()
		switch op {
		case "+":
			ce.operatorStack.Push("add")
		case "-":
			ce.operatorStack.Push("sub")
		case "*":
			ce.operatorStack.Push(CallCode("Math.multiply", 2))
		case "/":
			ce.operatorStack.Push(CallCode("Math.divide", 2))
		case "=":
			ce.operatorStack.Push("eq")
		case ">":
			ce.operatorStack.Push("gt")
		case "<":
			ce.operatorStack.Push("lt")
		case "&":
			ce.operatorStack.Push("and")
		case "|":
			ce.operatorStack.Push("or")
		}

		ce.t.Advance()
	}
	// Shunting yard for operators
	for i := 0; i < nOps; i++ {
		op := ce.operatorStack.Pop()
		ce.vmwriter.Add(op)
	}
	return expression, nil
}

// Start: (
// End:   )
// It also returns the number of expressions
func (ce *CompilationEngine) compileExpressionList() (Elem, int, error) {
	expressionList := NewSyntaxElem("expressionList")
	nExpressions := 0
	if ce.t.Current().Type() == SYMBOL && ce.t.Current().String() == ")" {
		ce.t.Backward()
		return expressionList, nExpressions, nil
	}
	for {
		expression, err := ce.compileExpression()
		if err != nil {
			return nil, -1, fmt.Errorf("Failed to compile expression: %v", err)
		}
		expressionList.AddChild(expression)
		nExpressions++
		a, err := ce.t.LookAhead(1)
		if err != nil {
			return nil, -1, err
		}
		if !(a.Type() == SYMBOL && a.String() == ",") {
			break
		}
		// ,
		ce.t.Advance()
		expressionList.AddChild(ce.NewTokenElemCurrent())
		ce.t.Advance()
	}
	return expressionList, nExpressions, nil
}

func (ce *CompilationEngine) compileTerm() (Elem, error) {
	term := NewSyntaxElem("term")

	cur := ce.t.Current()
	if cur.Type() == INT_CONST {
		// integerConstant
		i, err := strconv.Atoi(cur.String())
		if err != nil {
			return nil, err
		}
		ce.vmwriter.Add(PushCode("constant", i))
		term.AddChild(ce.NewTokenElemCurrent())
	} else if cur.Type() == STR_CONST {
		// stringConstant
		term.AddChild(ce.NewTokenElemCurrent())
		strconst := ce.t.Current().String()
		ce.vmwriter.Add(PushCode("constant", len(strconst)))
		ce.vmwriter.Add(CallCode("String.new", 1))
		for _, r := range strconst {
			ce.vmwriter.Add(PushCode("constant", int(r)))
			ce.vmwriter.Add(CallCode("String.appendChar", 2))
		}

	} else if isKeywordConstant(ce.t.Current()) {
		// keywordConstant
		term.AddChild(ce.NewTokenElemCurrent())

		// Write keyword
		keyword := ce.t.Current().String()
		switch keyword {
		case "true":
			// true = -1 (0xFFFF)
			ce.vmwriter.Add(PushCode("constant", 0))
			ce.vmwriter.Add("not")
		case "false":
			// false = 0 (0x0000)
			ce.vmwriter.Add(PushCode("constant", 0))
		case "this":
			ce.vmwriter.Add(PushCode("pointer", 0))
		}
	} else if isUnaryOp(ce.t.Current()) {
		// UnaryOp term
		var op string
		switch ce.t.Current().String() {
		case "~":
			op = "not"
		case "-":
			op = "neg"
		}
		term.AddChild(ce.NewTokenElemCurrent())

		ce.t.Advance()
		term2, err := ce.compileTerm()
		if err != nil {
			return nil, fmt.Errorf("Failed to compile UnaryOp term: %v", err)
		}
		term.AddChild(term2)
		ce.vmwriter.Add(op)
	} else if cur.Type() == SYMBOL && cur.String() == "(" {
		// ( expression )
		term.AddChild(ce.NewTokenElemCurrent())

		ce.operatorStack.Push("(")

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

		// Shunting yard algorithm
		for op := ce.operatorStack.Pop(); op != "("; op = ce.operatorStack.Pop() {
			ce.vmwriter.Add(op)
		}

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
			expressionList, _, err := ce.compileExpressionList()
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
			// Array
			// varName [ expression ]
			varName := cur.String()
			var table *SymbolTable
			var ok bool
			if table, ok = ce.resolveVariableInSubroutine(varName); ok {
			} else {
				return nil, fmt.Errorf("Variable %s is not defined.", varName)
			}
			varIndex, _ := table.IndexOf(varName)
			varKind, _ := table.KindOf(varName)

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
			// array head + index
			ce.vmwriter.Add(PushCode(varKindToSegment(varKind), varIndex))
			ce.vmwriter.Add("add")
			ce.vmwriter.Add(PopCode("pointer", 1))
			ce.vmwriter.Add(PushCode("that", 0))
		} else if a.Type() == SYMBOL && a.String() == "." {
			// (className | varName).subroutineName(expressionList)
			ce.compileSubroutineCall(term)
		} else {
			// varName
			term.AddChild(ce.NewTokenElemCurrent())

			// Push the variable
			varName := ce.t.Current().String()
			table, ok := ce.resolveVariableInSubroutine(varName)
			if !ok {
				return nil, fmt.Errorf("Variable %s is undefined.", varName)
			}
			varKind, _ := table.KindOf(varName)
			varIndex, _ := table.IndexOf(varName)
			segment := varKindToSegment(varKind)
			ce.vmwriter.Add(PushCode(segment, varIndex))
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
	prefix := ce.t.Current().String()
	var subroutineName string
	e.AddChild(ce.NewTokenElemCurrent())

	a1, err := ce.t.LookAhead(1)
	if err != nil {
		return err
	}

	// Read the name and decide whether a method or a function/constuctor by its name. See p208.
	// TODO: Refactoring
	isMethod := false
	isReceiverVariable := false
	var varKind string
	var varType string
	var varIndex int
	if a1.String() == "." {
		// varName.foo() is a method.
		if table, ok := ce.resolveVariableInSubroutine(prefix); ok {
			isMethod = true
			isReceiverVariable = true
			varType, _ = table.TypeOf(prefix)
			varKind, _ = table.KindOf(prefix)
			prefix = varType
		}

		// className.foo() is a function/consturctor
		ce.t.Advance()
		err = ce.validateCurrent(SYMBOL, ".")
		if err != nil {
			return err
		}
		e.AddChild(ce.NewTokenElemCurrent())
		subroutineName = prefix + ce.t.Current().String()

		ce.t.Advance()
		err = ce.validateCurrentType(IDENTIFIER)
		if err != nil {
			return err
		}
		e.AddChild(ce.NewTokenElemCurrent())
		subroutineName += ce.t.Current().String()
	} else {
		// foo() is a method. Completing type name.
		isMethod = true
		subroutineName = fmt.Sprintf("%s.%s", ce.classTable().Name(), prefix)
	}
	ce.t.Advance()
	err = ce.validateCurrent(SYMBOL, "(")
	if err != nil {
		return err
	}
	e.AddChild(ce.NewTokenElemCurrent())

	ce.t.Advance()
	expressionList, nExpressions, err := ce.compileExpressionList()
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

	// Pass the instance to the method as 1st parameter
	if isMethod {
		nExpressions++
		if isReceiverVariable {
			segment := varKindToSegment(varKind)
			ce.vmwriter.Add(PushCode(segment, varIndex))
		} else {
			// The receiver is this
			ce.vmwriter.Add(PushCode("pointer", 0))
		}
	}
	ce.vmwriter.Add(CallCode(subroutineName, nExpressions))
	return nil
}

func varKindToSegment(varKind string) string {
	switch varKind {
	case "var":
		return "local"
	case "field":
		return "this"
	case "argument":
		return "argument"
		//TODO static
	}
	return ""
}
