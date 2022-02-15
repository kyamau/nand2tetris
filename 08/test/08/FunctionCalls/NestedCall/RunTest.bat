set ROOT=..\..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=NestedCall

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe Sys.vm
move Sys.asm NestedCall.asm
%TOOL%CPUEmulator.bat %TARGET%.tst
