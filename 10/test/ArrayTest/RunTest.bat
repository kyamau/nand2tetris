set ROOT=..\..\..\
set TOOL=%ROOT%..\..\tools\
set TARGET=ArrayTest

go build  -o %ROOT%compiler.exe %ROOT%
%ROOT%compiler.exe .
rem %TOOL%CPUEmulator.bat %TARGET%.tst
