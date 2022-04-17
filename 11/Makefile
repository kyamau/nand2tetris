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
	-Remove-Item test/Square/*.out > $$null
	-Remove-Item test/ArrayTest/*.out > $$null

# Tests for tokinization
token_square: build
	./compiler.exe -tokenize test/Square/Square.jack
	./compiler.exe -tokenize test/Square/Main.jack
	./compiler.exe -tokenize test/Square/SquareGame.jack
	$(TOOL_DIR)\TextComparer.bat ./test/Square/SquareT.xml  ./test/Square/SquareT.xml.out
	$(TOOL_DIR)\TextComparer.bat ./test/Square/MainT.xml  ./test/Square/MainT.xml.out
	$(TOOL_DIR)\TextComparer.bat ./test/Square/SquareGameT.xml  ./test/Square/SquareGameT.xml.out

token_arraytest: build
	./compiler.exe -tokenize test/ArrayTest/Main.jack
	$(TOOL_DIR)\TextComparer.bat test/ArrayTest/MainT.xml test/ArrayTest/MainT.xml.out

parse_arraytest: build
	./compiler.exe -tokenize test/ArrayTest/Main.jack
	$(TOOL_DIR)\TextComparer.bat test/ArrayTest/Main.xml.out test/ArrayTest/Main.xml

parse_explesssquare: build 
	./compiler.exe test/ExpressionLessSquare/
	$(TOOL_DIR)\TextComparer.bat test/ExpressionLessSquare/Main.xml.out test/ExpressionLessSquare/Main.xml
	$(TOOL_DIR)\TextComparer.bat test/ExpressionLessSquare/Square.xml.out test/ExpressionLessSquare/Square.xml
	$(TOOL_DIR)\TextComparer.bat test/ExpressionLessSquare/SquareGame.xml.out test/ExpressionLessSquare/SquareGame.xml

parse_square: build 
	./compiler.exe test/Square/
	$(TOOL_DIR)\TextComparer.bat test/Square/Main.xml.out test/Square/Main.xml
	$(TOOL_DIR)\TextComparer.bat test/Square/Square.xml.out test/Square/Square.xml
	$(TOOL_DIR)\TextComparer.bat test/Square/SquareGame.xml.out test/Square/SquareGame.xml


test: build token_square token_arraytest parse_arraytest parse_explesssquare parse_square

