set ROOT=..\..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=BasicLoop

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe .
%TOOL%CPUEmulator.bat %TARGET%.tst
