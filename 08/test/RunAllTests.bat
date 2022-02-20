@echo off
for /r . %%f in (*.bat) do (
  if not %%~nf == %~n0 ( 
    echo RUN %%f
    cd %%~dpf
    call %%f
  )
)
