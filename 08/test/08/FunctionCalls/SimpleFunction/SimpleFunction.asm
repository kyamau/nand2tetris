(SimpleFunction.test)
@0 // [Start:WritePushPop - push(1, constant, 0)]
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@0 // [Start:WritePushPop - push(1, constant, 0)]
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@0 // [Start:WritePushPop - push(1, local, 0)] Set segment + index address to D
D=A
@LCL
D=D+M
A=D
D=M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@1 // [Start:WritePushPop - push(1, local, 1)] Set segment + index address to D
D=A
@LCL
D=D+M
A=D
D=M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(0)]Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP // Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // add
D=M
@13
D=D+M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(8)]Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@13 // not
D=!M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@0 // [Start:WritePushPop - push(1, argument, 0)] Set segment + index address to D
D=A
@ARG
D=D+M
A=D
D=M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(0)]Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP // Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // add
D=M
@13
D=D+M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@1 // [Start:WritePushPop - push(1, argument, 1)] Set segment + index address to D
D=A
@ARG
D=D+M
A=D
D=M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // [Start:WriteArithmetic(1)]Pop to the address in D
M=M-1
A=M
D=M
@13 // Pop y to R13
M=D
@SP // Pop to the address in D
M=M-1
A=M
D=M
@14 // Pop x to R14
M=D
@14 // sub
D=M
@13
D=D-M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@LCL // [Start:WriteReturn] R13 = LCL
D=M
@R15
M=D
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
@ARG
D=M
@SP
M=D+1
@R15
M=M-1
A=M
D=M
@THAT
M=D
@R15
M=M-1
A=M
D=M
@THIS
M=D
@R15
M=M-1
A=M
D=M
@ARG
M=D
@R15
M=M-1
A=M
D=M
@LCL
M=D
@R15
M=M-1
A=M
0;JMP // [Start:WriteGotoA()
