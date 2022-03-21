package tokenizer

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func newIntConstIgnoreErr(s string) *IntConst {
	ic, _ := NewIntConst(s)
	return ic
}

func TestTokenize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{"stringConstant", args{"\"azAZあclass{09\n\""}, []Token{NewStrConst("azAZあclass{09")}},
		{"keyword", args{"class"}, []Token{NewKeyword("class")}},
		{"identifier", args{"class_09"}, []Token{NewIdentifier("class_09")}},
		{"symbol", args{"{}()[].,;+-*/&|<>=~"}, []Token{NewSymbol("{"), NewSymbol("}"), NewSymbol("("), NewSymbol(")"), NewSymbol("["), NewSymbol("]"), NewSymbol("."), NewSymbol(","), NewSymbol(";"), NewSymbol("+"), NewSymbol("-"), NewSymbol("*"), NewSymbol("/"), NewSymbol("&"), NewSymbol("|"), NewSymbol("<"), NewSymbol(">"), NewSymbol("="), NewSymbol("~")}},
		{"integerConstant", args{"09"}, []Token{newIntConstIgnoreErr("09")}},
		{"combination", args{"\"azAZあclass{09\n\"class class_09{123"}, []Token{NewStrConst("azAZあclass{09"), NewKeyword("class"), NewIdentifier("class_09"), NewSymbol("{"), newIntConstIgnoreErr("123")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.args.s)
			if err != nil {
				t.Errorf("%v", err)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreprocess(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"single line comment", args{"code//comment\ncode"}, "code\ncode"},
		{"multi line comment", args{"code/** comment1\ncomment2. */code"}, "codecode"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := preprocess(tt.args.src); got != tt.want {
				t.Errorf("preprocess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestXML(t *testing.T) {
	r := strings.NewReader("\"test1\n\" class{123")
	tokenizer, err := NewTokenizer(r)
	if err != nil {
		t.Errorf("%v", err)
	}
	tokenizer.Tokenize()
	xml := tokenizer.XML()
	fmt.Println(xml)
}
