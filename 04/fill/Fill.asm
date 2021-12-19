// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.


// Constants

// End of the screen
// 512*256/16=8192

(KEY)
  @KBD // = 24576
  D=M

  @WHITE
  D; JEQ

  @BLACK
  D; JNE


  // Store color in @COLOR and jump to FILL.
  // -1:=black, 0:=white
  (WHITE)
    @COLOR
    M=0
    @FILL
    0; JMP

  (BLACK)
    @COLOR
    M=-1
    @FILL
    0; JMP


  (FILL)
  // The start adress of the screen = 16384 (=@SCREEN)
  // The number of addresses in the screen = 512 / 16 * 32 = 8192.
  // The end of address the screen = 16384 + 8191 = 24575.
  // We use a OFFSET moving from 8191 to 0.
  // Current position is @SCREEN + OFFSET

  // Initialize offset = 8191.
  @8191
  D=A
  @OFFSET
  M=D

  // for OFFSET = 8192; OFFSET >= 0; OFFSET--
  (LOOP)
    @OFFSET
    D=M

    @SCREEN // The start address of the screen
    D=D+A 
    @POS // = @SCREEN + OFFSET
    M=D

    @COLOR // Get color
    D=M

    @POS
    A=M
    M=D // Write color

    @OFFSET
    M=M-1 // OFFSET -= 1

    @OFFSET
    D=M  
    @LOOP
    D; JGE
@KEY
0;JMP
