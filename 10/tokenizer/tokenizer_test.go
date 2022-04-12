package tokenizer

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func newIntConstIgnoreErr(s string) *IntConst {
	ic, _ := NewIntConst(s, []int{0, 0})
	return ic
}

func TestTokenize(t *testing.T) {
	pos := []int{0, 0}
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{"stringConstant", args{"\"azAZあclass{09\n\""}, []Token{NewStrConst("azAZあclass{09", pos)}},
		{"keyword", args{"class"}, []Token{NewKeyword("class", pos)}},
		{"identifier", args{"class_09"}, []Token{NewIdentifier("class_09", pos)}},
		{"symbol", args{"{}()[].,;+-*/&|<>=~"}, []Token{NewSymbol("{", pos), NewSymbol("}", pos), NewSymbol("(", pos), NewSymbol(")", pos), NewSymbol("[", pos), NewSymbol("]", pos), NewSymbol(".", pos), NewSymbol(",", pos), NewSymbol(";", pos), NewSymbol("+", pos), NewSymbol("-", pos), NewSymbol("*", pos), NewSymbol("/", pos), NewSymbol("&", pos), NewSymbol("|", pos), NewSymbol("<", pos), NewSymbol(">", pos), NewSymbol("=", pos), NewSymbol("~", pos)}},
		{"integerConstant", args{"09"}, []Token{newIntConstIgnoreErr("09")}},
		{"combination", args{"\"azAZあclass{09\n\"class class_09{123"}, []Token{NewStrConst("azAZあclass{09", pos), NewKeyword("class", pos), NewIdentifier("class_09", pos), NewSymbol("{", pos), newIntConstIgnoreErr("123")}},
		{"single line comment", args{"code//comment\ncode"}, []Token{NewIdentifier("code", pos), NewIdentifier("code", pos)}},
		{"multi line comment", args{"code/** comment1\ncomment2. */code"}, []Token{NewIdentifier("code", pos), NewIdentifier("code", pos)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.args.s)
			if err != nil {
				t.Errorf("%v", err)
			} else if !isTokenEqualExceptPos(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenize_Pos(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{"simple", args{"foo\nbar bar2\nbuzz"}, []Token{NewIdentifier("foo", []int{1, 1}), NewIdentifier("bar", []int{2, 1}), NewIdentifier("bar2", []int{2, 5}), NewIdentifier("buzz", []int{3, 1})}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenize(tt.args.s)
			if err != nil {
				t.Errorf("%v", err)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%v, want %v", got, tt.want)
			}
		})
	}
}

func isTokenEqualExceptPos(tokens1 []Token, tokens2 []Token) bool {
	l := len(tokens1)
	if l2 := len(tokens2); l2 > l {
		l = l2
	}
	for i := 0; i < l; i++ {
		if tokens1[i].String() != tokens2[i].String() || tokens1[i].Type() != tokens2[i].Type() {
			fmt.Printf("want %v\n", tokens1[i])
			fmt.Printf("got %v\nn", tokens2[i])
			return false
		}
	}
	return true
}

// func TestPreprocess(t *testing.T) {
// 	type args struct {
// 		src string
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want string
// 	}{
// 		{"single line comment", args{"code//comment\ncode"}, "code\ncode"},
// 		{"multi line comment", args{"code/** comment1\ncomment2. */code"}, "codecode"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := preprocess(tt.args.src); got != tt.want {
// 				t.Errorf("preprocess() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

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
