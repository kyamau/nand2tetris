set ROOT=..\..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=StackTest

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe -bootstrap=false .
%TOOL%CPUEmulator.bat %TARGET%.tst
