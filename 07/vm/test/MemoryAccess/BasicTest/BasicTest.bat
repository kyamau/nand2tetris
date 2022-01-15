set ROOT=..\..\..\

go build  -o %ROOT%vm.exe %ROOT%
%ROOT%vm.exe BasicTest.vm >  BasicTest.asm