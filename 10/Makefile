SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -Command 

BINARIES:=compiler.exe
GO_FILES:=$(shell Get-ChildItem -r -Filter "*.go" -Name)

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
	./compiler.exe test/Square/Square.jack > test/Square/SquareT.xml.out
	./compiler.exe test/Square/Main.jack > test/Square/MainT.xml.out
	./compiler.exe test/Square/SquareGame.jack > test/Square/SquareGameT.xml.out
	Compare-Object -ReferenceObject @(Get-Content -Path ./test/Square/SquareT.xml) -DifferenceObject @(Get-Content -Path ./test/Square/SquareT.xml.out)
	Compare-Object -ReferenceObject @(Get-Content -Path ./test/Square/MainT.xml) -DifferenceObject @(Get-Content -Path ./test/Square/MainT.xml.out)
	Compare-Object -ReferenceObject @(Get-Content -Path ./test/Square/SquareGameT.xml) -DifferenceObject @(Get-Content -Path ./test/Square/SquareGameT.xml.out)

token_arraytest: build
	./compiler.exe test/ArrayTest/Main.jack > test/ArrayTest/MainT.xml.out
	Compare-Object -ReferenceObject @(Get-Content -Path test/ArrayTest/MainT.xml) -DifferenceObject @(Get-Content -Path test/ArrayTest/MainT.xml.out)
