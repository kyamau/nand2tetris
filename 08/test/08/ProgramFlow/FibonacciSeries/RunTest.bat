set ROOT=..\..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=FibonacciSeries

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe %TARGET%.vm
%TOOL%CPUEmulator.bat %TARGET%.tst
