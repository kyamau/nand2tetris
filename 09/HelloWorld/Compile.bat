set ROOT=..\
set TOOL=%ROOT%..\..\tools\

echo Copy OS
copy %TOOL%OS\* .
call %TOOL%JackCompiler.bat .
if not %ERRORLEVEL% == 0 (
  exit /b 1
)
