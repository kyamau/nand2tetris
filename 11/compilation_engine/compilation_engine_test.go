package compilation_engine

import (
	. "compiler/tokenizer"
	"io/ioutil"
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

func setupCompilationEngine(src string) *CompilationEngine {
	return NewCompilationEngine(*setupTokenizer(src))
}

func readAsString(path string) string {
	file, _ := ioutil.ReadFile(path)
	return string(file)

}
func TestCompilationEngine_XML(t *testing.T) {
	simpleClassSrc := readAsString("./test/simple_class.jack")
	simpleClassAns := readAsString("./test/simple_class.xml")
	simpleSubroutineSrc := readAsString("./test/simple_subroutine.jack")
	simpleSubroutineAns := readAsString("./test/simple_subroutine.xml")
	tests := []struct {
		name string
		ce   *CompilationEngine
		want string
	}{
		{name: "simple_class", ce: setupCompilationEngine(simpleClassSrc), want: simpleClassAns},
		{name: "simple_subroutine", ce: setupCompilationEngine(simpleSubroutineSrc), want: simpleSubroutineAns},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ce.Compile()
			if err != nil {
				t.Error(err)
			}
			if got := tt.ce.XML(); got != tt.want {
				t.Errorf("Parser.XML() = %v, want %v, diff=%v", got, tt.want, cmp.Diff(tt.want, got))
			}
		})
	}
}
