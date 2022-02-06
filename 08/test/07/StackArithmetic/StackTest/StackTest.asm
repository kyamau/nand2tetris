@17
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@17
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
@14 // eq
D=M
@13
D=D-M
@TRUE0 // Set true or false to D
D;JEQ
@0 // False: set 0 to D
D=A
@TFEND0
0;JMP
(TRUE0)
@1 // True: set -1 to D
D=-A
(TFEND0)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@17
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@16
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
@14 // eq
D=M
@13
D=D-M
@TRUE1 // Set true or false to D
D;JEQ
@0 // False: set 0 to D
D=A
@TFEND1
0;JMP
(TRUE1)
@1 // True: set -1 to D
D=-A
(TFEND1)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@16
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@17
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
@14 // eq
D=M
@13
D=D-M
@TRUE2 // Set true or false to D
D;JEQ
@0 // False: set 0 to D
D=A
@TFEND2
0;JMP
(TRUE2)
@1 // True: set -1 to D
D=-A
(TFEND2)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@892
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@891
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
@14 // lt
D=M
@13
D=D-M
@TRUE3 // Set true or false to D
D;JLT
@0 // False: set 0 to D
D=A
@TFEND3
0;JMP
(TRUE3)
@1 // True: set -1 to D
D=-A
(TFEND3)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@891
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@892
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
@14 // lt
D=M
@13
D=D-M
@TRUE4 // Set true or false to D
D;JLT
@0 // False: set 0 to D
D=A
@TFEND4
0;JMP
(TRUE4)
@1 // True: set -1 to D
D=-A
(TFEND4)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@891
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@891
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
@14 // lt
D=M
@13
D=D-M
@TRUE5 // Set true or false to D
D;JLT
@0 // False: set 0 to D
D=A
@TFEND5
0;JMP
(TRUE5)
@1 // True: set -1 to D
D=-A
(TFEND5)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@32767
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@32766
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
@14 // gt
D=M
@13
D=D-M
@TRUE6 // Set true or false to D
D;JGT
@0 // False: set 0 to D
D=A
@TFEND6
0;JMP
(TRUE6)
@1 // True: set -1 to D
D=-A
(TFEND6)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@32766
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@32767
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
@14 // gt
D=M
@13
D=D-M
@TRUE7 // Set true or false to D
D;JGT
@0 // False: set 0 to D
D=A
@TFEND7
0;JMP
(TRUE7)
@1 // True: set -1 to D
D=-A
(TFEND7)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@32766
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@32766
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
@14 // gt
D=M
@13
D=D-M
@TRUE8 // Set true or false to D
D;JGT
@0 // False: set 0 to D
D=A
@TFEND8
0;JMP
(TRUE8)
@1 // True: set -1 to D
D=-A
(TFEND8)
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@57
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@31
D=A
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@53
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
@112
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
@14 // sub
D=M
@13
D=D-M
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
@13 // neg
D=-M
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
@14 // and
D=M
@13
D=D&M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
@82
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
@14 // or
D=M
@13
D=D|M
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
@13 // not
D=!M
@SP // Push the value at the address in D
A=M
M=D
@SP
M=M+1
