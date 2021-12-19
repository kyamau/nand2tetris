// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Multiplies R0 and R1 and stores the result in R2.
// (R0, R1, R2 refer to RAM[0], RAM[1], and RAM[2], respectively.)
//
// This program only needs to handle arguments that satisfy
// R0 >= 0, R1 >= 0, and R0*R1 < 32768.

//R2 = 0
//for R3 = R1 ;R3 > 0 ;R3-- {
//  R2 += R0
//}
//
  
// R2 = 0
// Note: @0 means A = R0
@R2
M=0

// if R1 == 0 goto @END
(LOOP) 
@R1
D=M
@END
D; JEQ

// R1 -= 1
@R1
M=M-1

// R2 = R2 + R0
@R0
D=M
@R2
M=M+D

// goto @LOOP
@LOOP
0; JMP

@END
