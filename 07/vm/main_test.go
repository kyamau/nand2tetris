package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	type args struct {
		r       io.Reader
		outPath string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Simple Push Pop",
			args: args{r: strings.NewReader("push constant 3\r\npop local 3"),
				outPath: "./test.vm"},
		},
	}
	// Just kick Compile
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compile(tt.args.r)
			ioutil.WriteFile(tt.args.outPath, []byte(got), 644)
			fmt.Println(got)
		})
	}
}
