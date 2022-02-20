set ROOT=..\..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=PointerTest

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe -bootstrap=false .
%TOOL%CPUEmulator.bat %TARGET%.tst
