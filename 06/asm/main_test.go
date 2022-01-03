package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type testCase struct {
	name string
	args args
	want []uint16
}
type args struct {
	r io.Reader
}

func readTestcases(dirPath string) []testCase {
	cases := make([]testCase, 0)
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".asm" {
			_, filename := filepath.Split(path)
			if strings.Contains(filename, "L") {
				asmf, _ := os.Open(path)
				r := bufio.NewReader(asmf)
				want := readHack(path[:len(path)-4] + ".hack")
				cases = append(cases, testCase{name: filename, args: args{r: r}, want: want})
			}
		}
		return nil
	})
	return cases
}

func readHack(filePath string) []uint16 {
	f, _ := os.Open(filePath)
	b, _ := ioutil.ReadAll(f)
	s := string(b)
	s = strings.TrimRight(s, "\r\n")
	hack := make([]uint16, 0)
	for _, l := range strings.Split(s, "\r\n") {
		u, _ := strconv.ParseUint(l, 2, 16)
		hack = append(hack, uint16(u))
	}
	return hack
}

func TestCompile(t *testing.T) {

	easy := testCase{name: "easy case",
		args: args{r: strings.NewReader("//foo\r\n\r\n@1\r\nD=A-1;JGT\r\n")},
		want: []uint16{0b0000000000000001, 0b1110110010010001},
	}
	tests := readTestcases("./test")
	tests = append([]testCase{easy}, tests...)
	for _, tt := range tests {
		log.Printf("Test %v", tt.name)
		t.Run(tt.name, func(t *testing.T) {
			got := Compile(tt.args.r)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compile() = %b, want %b", got, tt.want)
			}
		})
	}
}
