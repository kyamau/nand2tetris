package main

import (
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestCompile(t *testing.T) {
	type args struct {
		r         io.Reader
		vmName    string
		bootstrap bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Simple Push Pop",
			args: args{r: strings.NewReader("push constant 3\r\npop local 3"),
				vmName: "simple_push_pop.vm", bootstrap: false},
		},
	}
	// Just kick Compile
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compile(tt.args.r, tt.args.vmName, tt.args.bootstrap)
			fmt.Println(got)
		})
	}
}
