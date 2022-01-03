package parser

import (
	"reflect"
	"strings"
	"testing"
)

func Test_removeIrrelevants(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "comment",
			args: args{[]string{"abc//comment"}},
			want: []string{"abc"},
		},
		{
			name: "space and tab",
			args: args{[]string{"abc\tdef ghi", "   \t", "//comment"}},
			want: []string{"abcdefghi"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeIrrelevants(tt.args.lines); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Symbol(t *testing.T) {
	p1, _ := NewParser(strings.NewReader("@azAz_09$:"))
	tests := []struct {
		name string
		p    *Parser
		want string
	}{
		{
			name: "normal",
			p:    p1,
			want: "azAz_09$:",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Advance()
			if got := tt.p.Symbol(); got != tt.want {
				t.Errorf("Parser.Symbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestParser_Dest(t *testing.T) {
	p1, _ := NewParser(strings.NewReader("DA=M"))
	p2, _ := NewParser(strings.NewReader(";JMP"))
	tests := []struct {
		name string
		p    *Parser
		want string
	}{
		{
			name: "normal",
			p:    p1,
			want: "DA",
		},
		{
			name: "null",
			p:    p2,
			want: "null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Advance()
			if got := tt.p.Dest(); got != tt.want {
				t.Errorf("Parser.Dest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Jump(t *testing.T) {
	p1, _ := NewParser(strings.NewReader("A=D;JGT"))
	p2, _ := NewParser(strings.NewReader("D=A"))
	tests := []struct {
		name string
		p    *Parser
		want string
	}{
		{
			name: "normal",
			p:    p1,
			want: "JGT",
		},
		{
			name: "null",
			p:    p2,
			want: "null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Advance()
			if got := tt.p.Jump(); got != tt.want {
				t.Errorf("Parser.Jump() = %v, want %v", got, tt.want)
			}
		})
	}
}
