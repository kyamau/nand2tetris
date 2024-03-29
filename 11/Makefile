SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command 

BINARIES := compiler.exe
GO_FILES := $(shell Get-ChildItem -r -Filter "*.go" -Name)
TOOL_DIR := ..\..\tools

build: $(BINARIES)
$(BINARIES): $(GO_FILES)
	go build -o $@ .

clean: 
	# See https://ir9ex.hatenablog.jp/entry/20121206/1354774247
	-Remove-Item ./compiler.exe > $$null
	-Remove-Item -Recurse test/*.out > $$null

# Tests for chapter 10 parsing
test_parser: build test_parser/ArrayTest test_parser/ExpressionLessSquare test_parser/Square

test_parser/ArrayTest: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/*.xml -Exclude *T.xml)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ($$file.FullName + ".out") ($$file).FullName}

test_parser/ExpressionLessSquare: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/*.xml -Exclude *T.xml)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ($$file.FullName + ".out") ($$file).FullName}

test_parser/Square: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/*.xml -Exclude *T.xml)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ($$file.FullName + ".out") ($$file).FullName}

# Test for chapter 11 code generation
test/Seven: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/ans/*.vm)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ("$@/"+$$file.Name+".out") ("$@/ans/"+$$file.Name)}
	
test/ConvertToBin: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/ans/*.vm)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ("$@/"+$$file.Name+".out") ("$@/ans/"+$$file.Name)}
	
test/Square: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/ans/*.vm)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ("$@/"+$$file.Name+".out") ("$@/ans/"+$$file.Name)}

test/Average: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/ans/*.vm)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ("$@/"+$$file.Name+".out") ("$@/ans/"+$$file.Name)}

test/Pong: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/ans/*.vm)) {echo $$file.Name; $(TOOL_DIR)/TextComparer.bat ("$@/"+$$file.Name+".out") ("$@/ans/"+$$file.Name)}

test/ComplexArrays: build
	./compiler.exe $@
	-foreach ($$file in (Get-ChildItem -Path $@/ans/*.vm)) {echo $$file.Name; $(TOOL_DIR)\TextComparer.bat ("$@/"+$$file.Name+".out") ("$@/ans/"+$$file.Name)}

test_compiler: test/Seven test/ConvertToBin test/Square test/Average test/Pong test/ComplexArrays
