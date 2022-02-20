set ROOT=..\..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=SimpleAdd

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe .
%TOOL%CPUEmulator.bat %TARGET%.tst
