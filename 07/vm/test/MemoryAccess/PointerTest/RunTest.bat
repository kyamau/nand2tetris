set ROOT=..\..\..\
set TOOL=..\..\..\..\..\..\tools\
set TARGET=PointerTest

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe %TARGET%.vm >  %TARGET%.asm
%TOOL%CPUEmulator.bat %TARGET%.tst
