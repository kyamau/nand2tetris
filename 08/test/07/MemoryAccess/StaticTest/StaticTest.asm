@111
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@333
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@888
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
@13 // Load poped value to R13
M=D
@8 // Set segment + index address to D
D=A
@16
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@SP // Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@3 // Set segment + index address to D
D=A
@16
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@SP // Pop to the address in D
M=M-1
A=M
D=M
@13 // Load poped value to R13
M=D
@1 // Set segment + index address to D
D=A
@16
@14 // Load segment + index address to R14
M=D
@13 // Write the value in R13 to the address in R14
D=M
@14
A=M
M=D
@3 // Set segment + index address to D
D=A
@16
A=D
D=M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@1 // Set segment + index address to D
D=A
@16
A=D
D=M
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
@14 // sub
D=M
@13
D=D-M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@8 // Set segment + index address to D
D=A
@16
A=D
D=M
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
