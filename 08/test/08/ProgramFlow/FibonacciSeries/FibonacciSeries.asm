@1 // [Start:WritePushPop - push(1, argument, 1)] Set segment + index address to D
D=A
@ARG
D=D+M
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, pointer, 1)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@1 //  Set segment + index address to D
D=A
@3
D=D+A
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@0 // [Start:WritePushPop - push(1, constant, 0)]
D=A
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, that, 0)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@0 //  Set segment + index address to D
D=A
@THAT
D=D+M
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@1 // [Start:WritePushPop - push(1, constant, 1)]
D=A
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, that, 1)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@1 //  Set segment + index address to D
D=A
@THAT
D=D+M
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@0 // [Start:WritePushPop - push(1, argument, 0)] Set segment + index address to D
D=A
@ARG
D=D+M
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@2 // [Start:WritePushPop - push(1, constant, 2)]
D=A
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(1)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP //  Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // sub
D=M
@13
D=D-M
@SP // [End:WriteArithmetic] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, argument, 0)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@0 //  Set segment + index address to D
D=A
@ARG
D=D+M
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
(MAIN_LOOP_START)
@0 // [Start:WritePushPop - push(1, argument, 0)] Set segment + index address to D
D=A
@ARG
D=D+M
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP //  Pop to the address in D
M=M-1
A=M
D=M
@COMPUTE_ELEMENT // [Start:WriteIf(COMPUTE_ELEMENT)]
D;JNE // [End:WriteIf(COMPUTE_ELEMENT)]
@END_PROGRAM // [Start:WriteGoto(END_PROGRAM)]
0;JMP // [End:WriteGoto(END_PROGRAM)]
(COMPUTE_ELEMENT)
@0 // [Start:WritePushPop - push(1, that, 0)] Set segment + index address to D
D=A
@THAT
D=D+M
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@1 // [Start:WritePushPop - push(1, that, 1)] Set segment + index address to D
D=A
@THAT
D=D+M
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(0)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP //  Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // add
D=M
@13
D=D+M
@SP // [End:WriteArithmetic] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, that, 2)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@2 //  Set segment + index address to D
D=A
@THAT
D=D+M
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@1 // [Start:WritePushPop - push(1, pointer, 1)] Set segment + index address to D
D=A
@3
D=D+A
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@1 // [Start:WritePushPop - push(1, constant, 1)]
D=A
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(0)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP //  Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // add
D=M
@13
D=D+M
@SP // [End:WriteArithmetic] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, pointer, 1)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@1 //  Set segment + index address to D
D=A
@3
D=D+A
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@0 // [Start:WritePushPop - push(1, argument, 0)] Set segment + index address to D
D=A
@ARG
D=D+M
A=D
D=M // 
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@1 // [Start:WritePushPop - push(1, constant, 1)]
D=A
@SP // [End:WritePushPop] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(1)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP //  Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // sub
D=M
@13
D=D-M
@SP // [End:WriteArithmetic] Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WritePushPop - pop(2, argument, 0)] Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@0 //  Set segment + index address to D
D=A
@ARG
D=D+M
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@MAIN_LOOP_START // [Start:WriteGoto(MAIN_LOOP_START)]
0;JMP // [End:WriteGoto(MAIN_LOOP_START)]
(END_PROGRAM)
