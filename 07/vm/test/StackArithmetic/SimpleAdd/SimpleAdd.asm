@7
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@8
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@SP // Pop to the address in D
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
