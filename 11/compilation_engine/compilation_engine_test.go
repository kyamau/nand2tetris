package compilation_engine

import (
	. "compiler/tokenizer"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func setupTokenizer(content string) *Tokenizer {
	r := strings.NewReader(content)
	tokenizer, _ := NewTokenizer(r)
	tokenizer.Tokenize()
	return tokenizer
}

func setupParser(src string) *CompilationEngine {
	return &CompilationEngine{*setupTokenizer(src), nil}
}

func TestParser_XML(t *testing.T) {
	src1 := `
		class Foo {
		}`
	ans1 := `<class>
  <keyword> class </keyword>
  <identifier> Foo </identifier>
  <symbol> { </symbol>
  <symbol> } </symbol>
</class>`
	src2 := `
		class Square {
			constructor Square new(int Ax, int Ay, int Asize) {
			}
		}`
	ans2 := `<class>
  <keyword> class </keyword>
  <identifier> Square </identifier>
  <symbol> { </symbol>
  <subroutineDec>
    <keyword> constructor </keyword>
    <identifier> Square </identifier>
    <identifier> new </identifier>
    <symbol> ( </symbol>
    <parameterList>
      <keyword> int </keyword>
      <identifier> Ax </identifier>
      <symbol> , </symbol>
      <keyword> int </keyword>
      <identifier> Ay </identifier>
      <symbol> , </symbol>
      <keyword> int </keyword>
      <identifier> Asize </identifier>
    </parameterList>
    <symbol> ) </symbol>
    <subroutineBody>
      <symbol> { </symbol>
      <statements>
      </statements>
      <symbol> } </symbol>
    </subroutineBody>
  </subroutineDec>
  <symbol> } </symbol>
</class>`
	tests := []struct {
		name string
		ce   *CompilationEngine
		want string
	}{
		{name: "simple_class", ce: setupParser(src1), want: ans1},
		{name: "simple_subroutine", ce: setupParser(src2), want: ans2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ce.Parse()
			if err != nil {
				t.Error(err)
			}
			if got := tt.ce.XML(); got != tt.want {
				t.Errorf("Parser.XML() = %v, want %v, diff=%v", got, tt.want, cmp.Diff(tt.want, got))
			}
		})
	}
}
