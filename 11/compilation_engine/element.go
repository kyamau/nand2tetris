package compilation_engine

import (
	"bytes"
	. "compiler/tokenizer"
	"encoding/xml"
	"fmt"
)

type Elem interface {
	AddChild(c Elem)
	GetChildren() []Elem
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

func (e *BaseElem) GetChildren() []Elem {
	return e.children
}

// Make string of the element and children for debugging
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

func NewTokenElem(token Token) Elem {
	e := TokenElem{BaseElem: &BaseElem{elemName: token.Type(), children: make([]Elem, 0)}, token: token}
	return &e
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

func NewSyntaxElem(name string) Elem {
	e := SyntaxElem{BaseElem: &BaseElem{elemName: name, children: []Elem{}}}
	return &e
}
