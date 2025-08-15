// This code was imported from https://github.com/armfazh/rfc7748_precomputed

// CHECK_BMI2ADX triggers bmi2adx if supported,
// otherwise it fallbacks to legacy code.
#define CHECK_BMI2ADX(label, legacy, bmi2adx) \
    CMPB ·hasBmi2Adx(SB), $0  \
    JE label                  \
    bmi2adx                   \
    RET                       \
    label:                    \
    legacy                    \
    RET

// cselect is a conditional move
// if b=1: it copies y into x;
// if b=0: x remains with the same value;
// if b<> 0,1: undefined.
// Uses: AX, DX, FLAGS
// Instr: x86_64, cmov
#define cselect(x,y,b) \
    TESTQ b, b \
    MOVQ  0+x, AX; MOVQ  0+y, DX; CMOVQNE DX, AX; MOVQ AX,  0+x; \
    MOVQ  8+x, AX; MOVQ  8+y, DX; CMOVQNE DX, AX; MOVQ AX,  8+x; \
    MOVQ 16+x, AX; MOVQ 16+y, DX; CMOVQNE DX, AX; MOVQ AX, 16+x; \
    MOVQ 24+x, AX; MOVQ 24+y, DX; CMOVQNE DX, AX; MOVQ AX, 24+x; \
    MOVQ 32+x, AX; MOVQ 32+y, DX; CMOVQNE DX, AX; MOVQ AX, 32+x; \
    MOVQ 40+x, AX; MOVQ 40+y, DX; CMOVQNE DX, AX; MOVQ AX, 40+x; \
    MOVQ 48+x, AX; MOVQ 48+y, DX; CMOVQNE DX, AX; MOVQ AX, 48+x;

// cswap is a conditional swap
// if b=1: x,y <- y,x;
// if b=0: x,y remain with the same values;
// if b<> 0,1: undefined.
// Uses: AX, DX, R8, FLAGS
// Instr: x86_64, cmov
#define cswap(x,y,b) \
    TESTQ b, b \
    MOVQ  0+x, AX; MOVQ AX, R8; MOVQ  0+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX,  0+x; MOVQ DX,  0+y; \
    MOVQ  8+x, AX; MOVQ AX, R8; MOVQ  8+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX,  8+x; MOVQ DX,  8+y; \
    MOVQ 16+x, AX; MOVQ AX, R8; MOVQ 16+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX, 16+x; MOVQ DX, 16+y; \
    MOVQ 24+x, AX; MOVQ AX, R8; MOVQ 24+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX, 24+x; MOVQ DX, 24+y; \
    MOVQ 32+x, AX; MOVQ AX, R8; MOVQ 32+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX, 32+x; MOVQ DX, 32+y; \
    MOVQ 40+x, AX; MOVQ AX, R8; MOVQ 40+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX, 40+x; MOVQ DX, 40+y; \
    MOVQ 48+x, AX; MOVQ AX, R8; MOVQ 48+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX, 48+x; MOVQ DX, 48+y;

// additionLeg adds x and y and stores in z
// Uses: AX, DX, R8-R14, FLAGS
// Instr: x86_64
#define additionLeg(z,x,y) \
    MOVQ  0+x,  R8;  ADDQ  0+y,  R8; \
    MOVQ  8+x,  R9;  ADCQ  8+y,  R9; \
    MOVQ 16+x, R10;  ADCQ 16+y, R10; \
    MOVQ 24+x, R11;  ADCQ 24+y, R11; \
    MOVQ 32+x, R12;  ADCQ 32+y, R12; \
    MOVQ 40+x, R13;  ADCQ 40+y, R13; \
    MOVQ 48+x, R14;  ADCQ 48+y, R14; \
    MOVQ   $0,  AX;  ADCQ   $0,  AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    ADDQ AX,  R8; MOVQ  $0, AX; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ DX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    ADCQ $0,  AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    ADDQ AX,  R8;  MOVQ  R8,  0+z; \
    ADCQ $0,  R9;  MOVQ  R9,  8+z; \
    ADCQ $0, R10;  MOVQ R10, 16+z; \
    ADCQ DX, R11;  MOVQ R11, 24+z; \
    ADCQ $0, R12;  MOVQ R12, 32+z; \
    ADCQ $0, R13;  MOVQ R13, 40+z; \
    ADCQ $0, R14;  MOVQ R14, 48+z;


// additionAdx adds x and y and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64, adx
#define additionAdx(z,x,y) \
    MOVL $32, R15; \
    XORL DX, DX; \
    MOVQ  0+x,  R8;  ADCXQ  0+y,  R8; \
    MOVQ  8+x,  R9;  ADCXQ  8+y,  R9; \
    MOVQ 16+x, R10;  ADCXQ 16+y, R10; \
    MOVQ 24+x, R11;  ADCXQ 24+y, R11; \
    MOVQ 32+x, R12;  ADCXQ 32+y, R12; \
    MOVQ 40+x, R13;  ADCXQ 40+y, R13; \
    MOVQ 48+x, R14;  ADCXQ 48+y, R14; \
    ;;;;;;;;;;;;;;;  ADCXQ   DX,  DX; \
    XORL AX, AX; \
    ADCXQ DX,  R8; SHLXQ R15, DX, DX; \
    ADCXQ AX,  R9; \
    ADCXQ AX, R10; \
    ADCXQ DX, R11; \
    ADCXQ AX, R12; \
    ADCXQ AX, R13; \
    ADCXQ AX, R14; \
    ADCXQ AX,  AX; \
    XORL  DX,  DX; \
    ADCXQ AX,  R8;  MOVQ  R8,  0+z; SHLXQ R15, AX, AX; \
    ADCXQ DX,  R9;  MOVQ  R9,  8+z; \
    ADCXQ DX, R10;  MOVQ R10, 16+z; \
    ADCXQ AX, R11;  MOVQ R11, 24+z; \
    ADCXQ DX, R12;  MOVQ R12, 32+z; \
    ADCXQ DX, R13;  MOVQ R13, 40+z; \
    ADCXQ DX, R14;  MOVQ R14, 48+z;

// subtraction subtracts y from x and stores in z
// Uses: AX, DX, R8-R14, FLAGS
// Instr: x86_64
#define subtraction(z,x,y) \
    MOVQ  0+x,  R8;  SUBQ  0+y,  R8; \
    MOVQ  8+x,  R9;  SBBQ  8+y,  R9; \
    MOVQ 16+x, R10;  SBBQ 16+y, R10; \
    MOVQ 24+x, R11;  SBBQ 24+y, R11; \
    MOVQ 32+x, R12;  SBBQ 32+y, R12; \
    MOVQ 40+x, R13;  SBBQ 40+y, R13; \
    MOVQ 48+x, R14;  SBBQ 48+y, R14; \
    MOVQ   $0,  AX;  SETCS AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    SUBQ AX,  R8; MOVQ  $0, AX; \
    SBBQ $0,  R9; \
    SBBQ $0, R10; \
    SBBQ DX, R11; \
    SBBQ $0, R12; \
    SBBQ $0, R13; \
    SBBQ $0, R14; \
    SETCS AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    SUBQ AX,  R8;  MOVQ  R8,  0+z; \
    SBBQ $0,  R9;  MOVQ  R9,  8+z; \
    SBBQ $0, R10;  MOVQ R10, 16+z; \
    SBBQ DX, R11;  MOVQ R11, 24+z; \
    SBBQ $0, R12;  MOVQ R12, 32+z; \
    SBBQ $0, R13;  MOVQ R13, 40+z; \
    SBBQ $0, R14;  MOVQ R14, 48+z;

// maddBmi2Adx multiplies x and y and accumulates in z
// Uses: AX, DX, R15, FLAGS
// Instr: x86_64, bmi2, adx
#define maddBmi2Adx(z,x,y,i,r0,r1,r2,r3,r4,r5,r6) \
    MOVQ   i+y, DX; XORL AX, AX; \
    MULXQ  0+x, AX, R8;  ADOXQ AX, r0;  ADCXQ R8, r1; MOVQ r0,i+z; \
    MULXQ  8+x, AX, r0;  ADOXQ AX, r1;  ADCXQ r0, r2; MOVQ $0, R8; \
    MULXQ 16+x, AX, r0;  ADOXQ AX, r2;  ADCXQ r0, r3; \
    MULXQ 24+x, AX, r0;  ADOXQ AX, r3;  ADCXQ r0, r4; \
    MULXQ 32+x, AX, r0;  ADOXQ AX, r4;  ADCXQ r0, r5; \
    MULXQ 40+x, AX, r0;  ADOXQ AX, r5;  ADCXQ r0, r6; \
    MULXQ 48+x, AX, r0;  ADOXQ AX, r6;  ADCXQ R8, r0; \
    ;;;;;;;;;;;;;;;;;;;  ADOXQ R8, r0;

// integerMulAdx multiplies x and y and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64, bmi2, adx
#define integerMulAdx(z,x,y) \
    MOVL    $0,R15; \
    MOVQ   0+y, DX;  XORL AX, AX;  MOVQ $0, R8; \
    MULXQ  0+x, AX,  R9;  MOVQ  AX, 0+z; \
    MULXQ  8+x, AX, R10;  ADCXQ AX,  R9; \
    MULXQ 16+x, AX, R11;  ADCXQ AX, R10; \
    MULXQ 24+x, AX, R12;  ADCXQ AX, R11; \
    MULXQ 32+x, AX, R13;  ADCXQ AX, R12; \
    MULXQ 40+x, AX, R14;  ADCXQ AX, R13; \
    MULXQ 48+x, AX, R15;  ADCXQ AX, R14; \
    ;;;;;;;;;;;;;;;;;;;;  ADCXQ R8, R15; \
    maddBmi2Adx(z,x,y, 8, R9,R10,R11,R12,R13,R14,R15) \
    maddBmi2Adx(z,x,y,16,R10,R11,R12,R13,R14,R15, R9) \
    maddBmi2Adx(z,x,y,24,R11,R12,R13,R14,R15, R9,R10) \
    maddBmi2Adx(z,x,y,32,R12,R13,R14,R15, R9,R10,R11) \
    maddBmi2Adx(z,x,y,40,R13,R14,R15, R9,R10,R11,R12) \
    maddBmi2Adx(z,x,y,48,R14,R15, R9,R10,R11,R12,R13) \
    MOVQ R15,  56+z; \
    MOVQ  R9,  64+z; \
    MOVQ R10,  72+z; \
    MOVQ R11,  80+z; \
    MOVQ R12,  88+z; \
    MOVQ R13,  96+z; \
    MOVQ R14, 104+z;

// maddLegacy multiplies x and y and accumulates in z
// Uses: AX, DX, R15, FLAGS
// Instr: x86_64
#define maddLegacy(z,x,y,i) \
    MOVQ  i+y, R15; \
    MOVQ  0+x, AX; MULQ R15; MOVQ AX,  R8; ;;;;;;;;;;;; MOVQ DX,  R9; \
    MOVQ  8+x, AX; MULQ R15; ADDQ AX,  R9; ADCQ $0, DX; MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R15; ADDQ AX, R10; ADCQ $0, DX; MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R15; ADDQ AX, R11; ADCQ $0, DX; MOVQ DX, R12; \
    MOVQ 32+x, AX; MULQ R15; ADDQ AX, R12; ADCQ $0, DX; MOVQ DX, R13; \
    MOVQ 40+x, AX; MULQ R15; ADDQ AX, R13; ADCQ $0, DX; MOVQ DX, R14; \
    MOVQ 48+x, AX; MULQ R15; ADDQ AX, R14; ADCQ $0, DX; \
    ADDQ  0+i+z,  R8; MOVQ  R8,  0+i+z; \
    ADCQ  8+i+z,  R9; MOVQ  R9,  8+i+z; \
    ADCQ 16+i+z, R10; MOVQ R10, 16+i+z; \
    ADCQ 24+i+z, R11; MOVQ R11, 24+i+z; \
    ADCQ 32+i+z, R12; MOVQ R12, 32+i+z; \
    ADCQ 40+i+z, R13; MOVQ R13, 40+i+z; \
    ADCQ 48+i+z, R14; MOVQ R14, 48+i+z; \
    ADCQ     $0,  DX; MOVQ  DX, 56+i+z;

// integerMulLeg multiplies x and y and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64
#define integerMulLeg(z,x,y) \
    MOVQ  0+y, R15; \
    MOVQ  0+x, AX; MULQ R15; MOVQ AX, 0+z; ;;;;;;;;;;;; MOVQ DX,  R8; \
    MOVQ  8+x, AX; MULQ R15; ADDQ AX,  R8; ADCQ $0, DX; MOVQ DX,  R9; MOVQ  R8,  8+z; \
    MOVQ 16+x, AX; MULQ R15; ADDQ AX,  R9; ADCQ $0, DX; MOVQ DX, R10; MOVQ  R9, 16+z; \
    MOVQ 24+x, AX; MULQ R15; ADDQ AX, R10; ADCQ $0, DX; MOVQ DX, R11; MOVQ R10, 24+z; \
    MOVQ 32+x, AX; MULQ R15; ADDQ AX, R11; ADCQ $0, DX; MOVQ DX, R12; MOVQ R11, 32+z; \
    MOVQ 40+x, AX; MULQ R15; ADDQ AX, R12; ADCQ $0, DX; MOVQ DX, R13; MOVQ R12, 40+z; \
    MOVQ 48+x, AX; MULQ R15; ADDQ AX, R13; ADCQ $0, DX; MOVQ DX,56+z; MOVQ R13, 48+z; \
    maddLegacy(z,x,y, 8) \
    maddLegacy(z,x,y,16) \
    maddLegacy(z,x,y,24) \
    maddLegacy(z,x,y,32) \
    maddLegacy(z,x,y,40) \
    maddLegacy(z,x,y,48)

// integerSqrLeg squares x and stores in z
// Uses: AX, CX, DX, R8-R15, FLAGS
// Instr: x86_64
#define integerSqrLeg(z,x) \
    XORL R15, R15; \
    MOVQ  0+x, CX; \
    MOVQ   CX, AX; MULQ CX; MOVQ AX, 0+z; MOVQ DX, R8; \
    ADDQ   CX, CX; ADCQ $0, R15; \
    MOVQ  8+x, AX; MULQ CX; ADDQ AX,  R8; ADCQ $0, DX; MOVQ DX,  R9; MOVQ R8, 8+z; \
    MOVQ 16+x, AX; MULQ CX; ADDQ AX,  R9; ADCQ $0, DX; MOVQ DX, R10; \
    MOVQ 24+x, AX; MULQ CX; ADDQ AX, R10; ADCQ $0, DX; MOVQ DX, R11; \
    MOVQ 32+x, AX; MULQ CX; ADDQ AX, R11; ADCQ $0, DX; MOVQ DX, R12; \
    MOVQ 40+x, AX; MULQ CX; ADDQ AX, R12; ADCQ $0, DX; MOVQ DX, R13; \
    MOVQ 48+x, AX; MULQ CX; ADDQ AX, R13; ADCQ $0, DX; MOVQ DX, R14; \
    \
    MOVQ  8+x, CX; \
    MOVQ   CX, AX; ADDQ R15, CX; MOVQ $0, R15; ADCQ $0, R15; \
    ;;;;;;;;;;;;;; MULQ CX; ADDQ  AX, R9; ADCQ $0, DX; MOVQ R9,16+z; \
    MOVQ  R15, AX; NEGQ AX; ANDQ 8+x, AX; ADDQ AX, DX; ADCQ $0, R11; MOVQ DX, R8; \
    ADDQ  8+x, CX; ADCQ $0, R15; \
    MOVQ 16+x, AX; MULQ CX; ADDQ AX, R10; ADCQ $0, DX; ADDQ R8, R10; ADCQ $0, DX; MOVQ DX, R8; MOVQ R10, 24+z; \
    MOVQ 24+x, AX; MULQ CX; ADDQ AX, R11; ADCQ $0, DX; ADDQ R8, R11; ADCQ $0, DX; MOVQ DX, R8; \
    MOVQ 32+x, AX; MULQ CX; ADDQ AX, R12; ADCQ $0, DX; ADDQ R8, R12; ADCQ $0, DX; MOVQ DX, R8; \
    MOVQ 40+x, AX; MULQ CX; ADDQ AX, R13; ADCQ $0, DX; ADDQ R8, R13; ADCQ $0, DX; MOVQ DX, R8; \
    MOVQ 48+x, AX; MULQ CX; ADDQ AX, R14; ADCQ $0, DX; ADDQ R8, R14; ADCQ $0, DX; MOVQ DX, R9; \
    \
    MOVQ 16+x, CX; \
    MOVQ   CX, AX; ADDQ R15, CX; MOVQ $0, R15; ADCQ $0, R15; \
    ;;;;;;;;;;;;;; MULQ CX; ADDQ AX, R11; ADCQ $0, DX; MOVQ R11, 32+z; \
    MOVQ  R15, AX; NEGQ AX; ANDQ 16+x,AX; ADDQ AX, DX; ADCQ $0, R13; MOVQ DX, R8; \
    ADDQ 16+x, CX; ADCQ $0, R15; \
    MOVQ 24+x, AX; MULQ CX; ADDQ AX, R12; ADCQ $0, DX; ADDQ R8, R12; ADCQ $0, DX; MOVQ DX, R8; MOVQ R12, 40+z; \
    MOVQ 32+x, AX; MULQ CX; ADDQ AX, R13; ADCQ $0, DX; ADDQ R8, R13; ADCQ $0, DX; MOVQ DX, R8; \
    MOVQ 40+x, AX; MULQ CX; ADDQ AX, R14; ADCQ $0, DX; ADDQ R8, R14; ADCQ $0, DX; MOVQ DX, R8; \
    MOVQ 48+x, AX; MULQ CX; ADDQ AX,  R9; ADCQ $0, DX; ADDQ R8,  R9; ADCQ $0, DX; MOVQ DX,R10; \
    \
    MOVQ 24+x, CX; \
    MOVQ   CX, AX; ADDQ R15, CX; MOVQ $0, R15; ADCQ $0, R15; \
    ;;;;;;;;;;;;;; MULQ CX; ADDQ AX, R13; ADCQ $0, DX; MOVQ R13, 48+z; \
    MOVQ  R15, AX; NEGQ AX; ANDQ 24+x,AX; ADDQ AX, DX; ADCQ $0,  R9; MOVQ DX, R8; \
    ADDQ 24+x, CX; ADCQ $0, R15; \
    MOVQ 32+x, AX; MULQ CX; ADDQ AX, R14; ADCQ $0, DX; ADDQ R8, R14; ADCQ $0, DX; MOVQ DX, R8; MOVQ R14, 56+z; \
    MOVQ 40+x, AX; MULQ CX; ADDQ AX,  R9; ADCQ $0, DX; ADDQ R8,  R9; ADCQ $0, DX; MOVQ DX, R8; \
    MOVQ 48+x, AX; MULQ CX; ADDQ AX, R10; ADCQ $0, DX; ADDQ R8, R10; ADCQ $0, DX; MOVQ DX,R11; \
    \
    MOVQ 32+x, CX; \
    MOVQ   CX, AX; ADDQ R15, CX; MOVQ $0, R15; ADCQ $0, R15; \
    ;;;;;;;;;;;;;; MULQ CX; ADDQ AX,  R9; ADCQ $0, DX; MOVQ R9, 64+z; \
    MOVQ  R15, AX; NEGQ AX; ANDQ 32+x,AX; ADDQ AX, DX; ADCQ $0, R11; MOVQ DX, R8; \
    ADDQ 32+x, CX; ADCQ $0, R15; \
    MOVQ 40+x, AX; MULQ CX; ADDQ AX, R10; ADCQ $0, DX; ADDQ R8, R10; ADCQ $0, DX; MOVQ DX, R8; MOVQ R10, 72+z; \
    MOVQ 48+x, AX; MULQ CX; ADDQ AX, R11; ADCQ $0, DX; ADDQ R8, R11; ADCQ $0, DX; MOVQ DX,R12; \
    \
    XORL R13, R13; \
    XORL R14, R14; \
    MOVQ 40+x, CX; \
    MOVQ   CX, AX; ADDQ R15, CX; MOVQ $0, R15; ADCQ $0, R15; \
    ;;;;;;;;;;;;;; MULQ CX; ADDQ AX, R11; ADCQ $0, DX; MOVQ R11, 80+z; \
    MOVQ  R15, AX; NEGQ AX; ANDQ 40+x,AX; ADDQ AX, DX; ADCQ $0, R13; MOVQ DX, R8; \
    ADDQ 40+x, CX; ADCQ $0, R15; \
    MOVQ 48+x, AX; MULQ CX; ADDQ AX, R12; ADCQ $0, DX; ADDQ R8, R12; ADCQ $0, DX; MOVQ DX, R8; MOVQ R12, 88+z; \
    ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;; ADDQ R8, R13; ADCQ $0,R14; \
    \
    XORL   R9, R9; \
    MOVQ 48+x, CX; \
    MOVQ   CX, AX; ADDQ R15, CX; MOVQ $0, R15; ADCQ $0, R15; \
    ;;;;;;;;;;;;;; MULQ CX; ADDQ AX, R13; ADCQ $0, DX; MOVQ R13, 96+z; \
    MOVQ  R15, AX; NEGQ AX; ANDQ 48+x,AX; ADDQ AX, DX; ADCQ $0, R9; MOVQ DX, R8; \
    ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;; ADDQ R8,R14; ADCQ $0, R9; MOVQ R14, 104+z;


// integerSqrAdx squares x and stores in z
// Uses: AX, CX, DX, R8-R15, FLAGS
// Instr: x86_64, bmi2, adx
#define integerSqrAdx(z,x) \
    XORL R15, R15; \
    MOVQ  0+x, DX; \
    ;;;;;;;;;;;;;; MULXQ DX, AX, R8; MOVQ AX, 0+z; \
    ADDQ   DX, DX; ADCQ $0, R15; CLC; \
    MULXQ  8+x, AX,  R9; ADCXQ AX,  R8; MOVQ R8, 8+z; \
    MULXQ 16+x, AX, R10; ADCXQ AX,  R9; MOVQ $0, R8;\
    MULXQ 24+x, AX, R11; ADCXQ AX, R10; \
    MULXQ 32+x, AX, R12; ADCXQ AX, R11; \
    MULXQ 40+x, AX, R13; ADCXQ AX, R12; \
    MULXQ 48+x, AX, R14; ADCXQ AX, R13; \
    ;;;;;;;;;;;;;;;;;;;; ADCXQ R8, R14; \
    \
    MOVQ  8+x, DX; \
    MOVQ   DX, AX; ADDQ R15, DX; MOVQ $0, R15; ADCQ  $0, R15; \
    MULXQ AX,  AX, CX; \
    MOVQ R15,  R8; NEGQ R8; ANDQ 8+x, R8; \
    ADDQ AX,  R9; MOVQ R9, 16+z; \
    ADCQ CX,  R8; \
    ADCQ $0, R11; \
    ADDQ  8+x,  DX; \
    ADCQ   $0, R15; \
    XORL R9, R9; ;;;;;;;;;;;;;;;;;;;;; ADOXQ R8, R10; \
    MULXQ 16+x, AX, CX; ADCXQ AX, R10; ADOXQ CX, R11; MOVQ R10, 24+z; \
    MULXQ 24+x, AX, CX; ADCXQ AX, R11; ADOXQ CX, R12; MOVQ  $0, R10; \
    MULXQ 32+x, AX, CX; ADCXQ AX, R12; ADOXQ CX, R13; \
    MULXQ 40+x, AX, CX; ADCXQ AX, R13; ADOXQ CX, R14; \
    MULXQ 48+x, AX, CX; ADCXQ AX, R14; ADOXQ CX,  R9; \
    ;;;;;;;;;;;;;;;;;;; ADCXQ R10, R9; \
    \
    MOVQ 16+x, DX; \
    MOVQ   DX, AX; ADDQ R15, DX; MOVQ $0, R15; ADCQ  $0, R15; \
    MULXQ AX,  AX, CX; \
    MOVQ R15,  R8; NEGQ R8; ANDQ 16+x, R8; \
    ADDQ AX, R11; MOVQ R11, 32+z; \
    ADCQ CX,  R8; \
    ADCQ $0, R13; \
    ADDQ 16+x,  DX; \
    ADCQ   $0, R15; \
    XORL R11, R11; ;;;;;;;;;;;;;;;;;;; ADOXQ R8, R12; \
    MULXQ 24+x, AX, CX; ADCXQ AX, R12; ADOXQ CX, R13; MOVQ R12, 40+z; \
    MULXQ 32+x, AX, CX; ADCXQ AX, R13; ADOXQ CX, R14; MOVQ  $0, R12; \
    MULXQ 40+x, AX, CX; ADCXQ AX, R14; ADOXQ CX,  R9; \
    MULXQ 48+x, AX, CX; ADCXQ AX,  R9; ADOXQ CX, R10; \
    ;;;;;;;;;;;;;;;;;;; ADCXQ R11,R10; \
    \
    MOVQ 24+x, DX; \
    MOVQ   DX, AX; ADDQ R15, DX; MOVQ $0, R15; ADCQ  $0, R15; \
    MULXQ AX,  AX, CX; \
    MOVQ R15,  R8; NEGQ R8; ANDQ 24+x, R8; \
    ADDQ AX, R13; MOVQ R13, 48+z; \
    ADCQ CX,  R8; \
    ADCQ $0,  R9; \
    ADDQ 24+x,  DX; \
    ADCQ   $0, R15; \
    XORL R13, R13; ;;;;;;;;;;;;;;;;;;; ADOXQ R8, R14; \
    MULXQ 32+x, AX, CX; ADCXQ AX, R14; ADOXQ CX,  R9; MOVQ R14, 56+z; \
    MULXQ 40+x, AX, CX; ADCXQ AX,  R9; ADOXQ CX, R10; MOVQ  $0, R14; \
    MULXQ 48+x, AX, CX; ADCXQ AX, R10; ADOXQ CX, R11; \
    ;;;;;;;;;;;;;;;;;;; ADCXQ R12,R11; \
    \
    MOVQ 32+x, DX; \
    MOVQ   DX, AX; ADDQ R15, DX; MOVQ $0, R15; ADCQ  $0, R15; \
    MULXQ AX,  AX, CX; \
    MOVQ R15,  R8; NEGQ R8; ANDQ 32+x, R8; \
    ADDQ AX,  R9; MOVQ R9, 64+z; \
    ADCQ CX,  R8; \
    ADCQ $0, R11; \
    ADDQ 32+x,  DX; \
    ADCQ   $0, R15; \
    XORL R9, R9; ;;;;;;;;;;;;;;;;;;;;; ADOXQ R8, R10; \
    MULXQ 40+x, AX, CX; ADCXQ AX, R10; ADOXQ CX, R11; MOVQ R10, 72+z; \
    MULXQ 48+x, AX, CX; ADCXQ AX, R11; ADOXQ CX, R12; \
    ;;;;;;;;;;;;;;;;;;; ADCXQ R13,R12; \
    \
    MOVQ 40+x, DX; \
    MOVQ   DX, AX; ADDQ R15, DX; MOVQ $0, R15; ADCQ  $0, R15; \
    MULXQ AX,  AX, CX; \
    MOVQ R15,  R8; NEGQ R8; ANDQ 40+x, R8; \
    ADDQ AX, R11; MOVQ R11, 80+z; \
    ADCQ CX,  R8; \
    ADCQ $0, R13; \
    ADDQ 40+x,  DX; \
    ADCQ   $0, R15; \
    XORL R11, R11; ;;;;;;;;;;;;;;;;;;; ADOXQ R8, R12; \
    MULXQ 48+x, AX, CX; ADCXQ AX, R12; ADOXQ CX, R13; MOVQ R12, 88+z; \
    ;;;;;;;;;;;;;;;;;;; ADCXQ R14,R13; \
    \
    MOVQ 48+x, DX; \
    MOVQ   DX, AX; ADDQ R15, DX; MOVQ $0, R15; ADCQ  $0, R15; \
    MULXQ AX,  AX, CX; \
    MOVQ R15,  R8; NEGQ R8; ANDQ 48+x, R8; \
    XORL R10, R10; ;;;;;;;;;;;;;; ADOXQ CX, R14; \
    ;;;;;;;;;;;;;; ADCXQ AX, R13; ;;;;;;;;;;;;;; MOVQ R13, 96+z; \
    ;;;;;;;;;;;;;; ADCXQ R8, R14; MOVQ R14, 104+z;

// reduceFromDoubleLeg finds a z=x modulo p such that z<2^448 and stores in z
// Uses: AX, R8-R15, FLAGS
// Instr: x86_64
#define reduceFromDoubleLeg(z,x) \
    /* (   ,2C13,2C12,2C11,2C10|C10,C9,C8, C7) + (C6,...,C0) */ \
    /* (r14, r13, r12, r11,     r10,r9,r8,r15) */ \
    MOVQ 80+x,AX; MOVQ AX,R10; \
    MOVQ $0xFFFFFFFF00000000, R8; \
    ANDQ R8,R10; \
    \
    MOVQ $0,R14; \
    MOVQ 104+x,R13; SHLQ $1,R13,R14; \
    MOVQ  96+x,R12; SHLQ $1,R12,R13; \
    MOVQ  88+x,R11; SHLQ $1,R11,R12; \
    MOVQ  72+x, R9; SHLQ $1,R10,R11; \
    MOVQ  64+x, R8; SHLQ $1,R10; \
    MOVQ $0xFFFFFFFF,R15; ANDQ R15,AX; ORQ AX,R10; \
    MOVQ  56+x,R15; \
    \
    ADDQ  0+x,R15; MOVQ R15, 0+z; MOVQ  56+x,R15; \
    ADCQ  8+x, R8; MOVQ  R8, 8+z; MOVQ  64+x, R8; \
    ADCQ 16+x, R9; MOVQ  R9,16+z; MOVQ  72+x, R9; \
    ADCQ 24+x,R10; MOVQ R10,24+z; MOVQ  80+x,R10; \
    ADCQ 32+x,R11; MOVQ R11,32+z; MOVQ  88+x,R11; \
    ADCQ 40+x,R12; MOVQ R12,40+z; MOVQ  96+x,R12; \
    ADCQ 48+x,R13; MOVQ R13,48+z; MOVQ 104+x,R13; \
    ADCQ   $0,R14; \
    /* (c10c9,c9c8,c8c7,c7c13,c13c12,c12c11,c11c10) + (c6,...,c0) */ \
    /* (   r9,  r8, r15,  r13,   r12,   r11,   r10) */ \
    MOVQ R10, AX; \
    SHRQ $32,R11,R10; \
    SHRQ $32,R12,R11; \
    SHRQ $32,R13,R12; \
    SHRQ $32,R15,R13; \
    SHRQ $32, R8,R15; \
    SHRQ $32, R9, R8; \
    SHRQ $32, AX, R9; \
    \
    ADDQ  0+z,R10; \
    ADCQ  8+z,R11; \
    ADCQ 16+z,R12; \
    ADCQ 24+z,R13; \
    ADCQ 32+z,R15; \
    ADCQ 40+z, R8; \
    ADCQ 48+z, R9; \
    ADCQ   $0,R14; \
    /* ( c7) + (c6,...,c0) */ \
    /* (r14) */ \
    MOVQ R14, AX; SHLQ $32, AX; \
    ADDQ R14,R10; MOVQ  $0,R14; \
    ADCQ  $0,R11; \
    ADCQ  $0,R12; \
    ADCQ  AX,R13; \
    ADCQ  $0,R15; \
    ADCQ  $0, R8; \
    ADCQ  $0, R9; \
    ADCQ  $0,R14; \
    /* ( c7) + (c6,...,c0) */ \
    /* (r14) */ \
    MOVQ R14, AX; SHLQ $32,AX; \
    ADDQ R14,R10; MOVQ R10, 0+z; \
    ADCQ  $0,R11; MOVQ R11, 8+z; \
    ADCQ  $0,R12; MOVQ R12,16+z; \
    ADCQ  AX,R13; MOVQ R13,24+z; \
    ADCQ  $0,R15; MOVQ R15,32+z; \
    ADCQ  $0, R8; MOVQ  R8,40+z; \
    ADCQ  $0, R9; MOVQ  R9,48+z;

// reduceFromDoubleAdx finds a z=x modulo p such that z<2^448 and stores in z
// Uses: AX, R8-R15, FLAGS
// Instr: x86_64, adx
#define reduceFromDoubleAdx(z,x) \
    /* (   ,2C13,2C12,2C11,2C10|C10,C9,C8, C7) + (C6,...,C0) */ \
    /* (r14, r13, r12, r11,     r10,r9,r8,r15) */ \
    MOVQ 80+x,AX; MOVQ AX,R10; \
    MOVQ $0xFFFFFFFF00000000, R8; \
    ANDQ R8,R10; \
    \
    MOVQ $0,R14; \
    MOVQ 104+x,R13; SHLQ $1,R13,R14; \
    MOVQ  96+x,R12; SHLQ $1,R12,R13; \
    MOVQ  88+x,R11; SHLQ $1,R11,R12; \
    MOVQ  72+x, R9; SHLQ $1,R10,R11; \
    MOVQ  64+x, R8; SHLQ $1,R10; \
    MOVQ $0xFFFFFFFF,R15; ANDQ R15,AX; ORQ AX,R10; \
    MOVQ  56+x,R15; \
    \
    XORL AX,AX; \
    ADCXQ  0+x,R15; MOVQ R15, 0+z; MOVQ  56+x,R15; \
    ADCXQ  8+x, R8; MOVQ  R8, 8+z; MOVQ  64+x, R8; \
    ADCXQ 16+x, R9; MOVQ  R9,16+z; MOVQ  72+x, R9; \
    ADCXQ 24+x,R10; MOVQ R10,24+z; MOVQ  80+x,R10; \
    ADCXQ 32+x,R11; MOVQ R11,32+z; MOVQ  88+x,R11; \
    ADCXQ 40+x,R12; MOVQ R12,40+z; MOVQ  96+x,R12; \
    ADCXQ 48+x,R13; MOVQ R13,48+z; MOVQ 104+x,R13; \
    ADCXQ   AX,R14; \
    /* (c10c9,c9c8,c8c7,c7c13,c13c12,c12c11,c11c10) + (c6,...,c0) */ \
    /* (   r9,  r8, r15,  r13,   r12,   r11,   r10) */ \
    MOVQ R10, AX; \
    SHRQ $32,R11,R10; \
    SHRQ $32,R12,R11; \
    SHRQ $32,R13,R12; \
    SHRQ $32,R15,R13; \
    SHRQ $32, R8,R15; \
    SHRQ $32, R9, R8; \
    SHRQ $32, AX, R9; \
    \
    XORL AX,AX; \
    ADCXQ  0+z,R10; \
    ADCXQ  8+z,R11; \
    ADCXQ 16+z,R12; \
    ADCXQ 24+z,R13; \
    ADCXQ 32+z,R15; \
    ADCXQ 40+z, R8; \
    ADCXQ 48+z, R9; \
    ADCXQ   AX,R14; \
    /* ( c7) + (c6,...,c0) */ \
    /* (r14) */ \
    MOVQ R14, AX; SHLQ $32, AX; \
    CLC; \
    ADCXQ R14,R10; MOVQ $0,R14; \
    ADCXQ R14,R11; \
    ADCXQ R14,R12; \
    ADCXQ  AX,R13; \
    ADCXQ R14,R15; \
    ADCXQ R14, R8; \
    ADCXQ R14, R9; \
    ADCXQ R14,R14; \
    /* ( c7) + (c6,...,c0) */ \
    /* (r14) */ \
    MOVQ R14, AX; SHLQ $32, AX; \
    CLC; \
    ADCXQ R14,R10; MOVQ R10, 0+z; MOVQ $0,R14; \
    ADCXQ R14,R11; MOVQ R11, 8+z; \
    ADCXQ R14,R12; MOVQ R12,16+z; \
    ADCXQ  AX,R13; MOVQ R13,24+z; \
    ADCXQ R14,R15; MOVQ R15,32+z; \
    ADCXQ R14, R8; MOVQ  R8,40+z; \
    ADCXQ R14, R9; MOVQ  R9,48+z;

// addSub calculates two operations: x,y = x+y,x-y
// Uses: AX, DX, R8-R15, FLAGS
#define addSub(x,y) \
    MOVQ  0+x,  R8;  ADDQ  0+y,  R8; \
    MOVQ  8+x,  R9;  ADCQ  8+y,  R9; \
    MOVQ 16+x, R10;  ADCQ 16+y, R10; \
    MOVQ 24+x, R11;  ADCQ 24+y, R11; \
    MOVQ 32+x, R12;  ADCQ 32+y, R12; \
    MOVQ 40+x, R13;  ADCQ 40+y, R13; \
    MOVQ 48+x, R14;  ADCQ 48+y, R14; \
    MOVQ   $0,  AX;  ADCQ   $0,  AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    ADDQ AX,  R8; MOVQ  $0, AX; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ DX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    ADCQ $0,  AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    ADDQ AX,  R8;  MOVQ  0+x,AX; MOVQ  R8,  0+x; MOVQ AX,  R8; \
    ADCQ $0,  R9;  MOVQ  8+x,AX; MOVQ  R9,  8+x; MOVQ AX,  R9; \
    ADCQ $0, R10;  MOVQ 16+x,AX; MOVQ R10, 16+x; MOVQ AX, R10; \
    ADCQ DX, R11;  MOVQ 24+x,AX; MOVQ R11, 24+x; MOVQ AX, R11; \
    ADCQ $0, R12;  MOVQ 32+x,AX; MOVQ R12, 32+x; MOVQ AX, R12; \
    ADCQ $0, R13;  MOVQ 40+x,AX; MOVQ R13, 40+x; MOVQ AX, R13; \
    ADCQ $0, R14;  MOVQ 48+x,AX; MOVQ R14, 48+x; MOVQ AX, R14; \
    SUBQ  0+y,  R8; \
    SBBQ  8+y,  R9; \
    SBBQ 16+y, R10; \
    SBBQ 24+y, R11; \
    SBBQ 32+y, R12; \
    SBBQ 40+y, R13; \
    SBBQ 48+y, R14; \
    MOVQ   $0,  AX;  SETCS AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    SUBQ AX,  R8; MOVQ  $0, AX; \
    SBBQ $0,  R9; \
    SBBQ $0, R10; \
    SBBQ DX, R11; \
    SBBQ $0, R12; \
    SBBQ $0, R13; \
    SBBQ $0, R14; \
    SETCS AX; \
    MOVQ AX,  DX; \
    SHLQ $32, DX; \
    SUBQ AX,  R8;  MOVQ  R8,  0+y; \
    SBBQ $0,  R9;  MOVQ  R9,  8+y; \
    SBBQ $0, R10;  MOVQ R10, 16+y; \
    SBBQ DX, R11;  MOVQ R11, 24+y; \
    SBBQ $0, R12;  MOVQ R12, 32+y; \
    SBBQ $0, R13;  MOVQ R13, 40+y; \
    SBBQ $0, R14;  MOVQ R14, 48+y;
