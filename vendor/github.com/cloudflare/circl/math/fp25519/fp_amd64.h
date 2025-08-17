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
    MOVQ 24+x, AX; MOVQ 24+y, DX; CMOVQNE DX, AX; MOVQ AX, 24+x;

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
    MOVQ 24+x, AX; MOVQ AX, R8; MOVQ 24+y, DX; CMOVQNE DX, AX; CMOVQNE R8, DX; MOVQ AX, 24+x; MOVQ DX, 24+y;

// additionLeg adds x and y and stores in z
// Uses: AX, DX, R8-R11, FLAGS
// Instr: x86_64, cmov
#define additionLeg(z,x,y) \
    MOVL $38, AX; \
    MOVL  $0, DX; \
    MOVQ  0+x,  R8;  ADDQ  0+y,  R8; \
    MOVQ  8+x,  R9;  ADCQ  8+y,  R9; \
    MOVQ 16+x, R10;  ADCQ 16+y, R10; \
    MOVQ 24+x, R11;  ADCQ 24+y, R11; \
    CMOVQCS AX, DX;    \
    ADDQ DX,  R8; \
    ADCQ $0,  R9;  MOVQ  R9,  8+z; \
    ADCQ $0, R10;  MOVQ R10, 16+z; \
    ADCQ $0, R11;  MOVQ R11, 24+z; \
    MOVL $0,  DX; \
    CMOVQCS AX, DX; \
    ADDQ DX,  R8;  MOVQ  R8,  0+z;

// additionAdx adds x and y and stores in z
// Uses: AX, DX, R8-R11, FLAGS
// Instr: x86_64, cmov, adx
#define additionAdx(z,x,y) \
    MOVL $38, AX; \
    XORL  DX, DX; \
    MOVQ  0+x,  R8;  ADCXQ  0+y,  R8; \
    MOVQ  8+x,  R9;  ADCXQ  8+y,  R9; \
    MOVQ 16+x, R10;  ADCXQ 16+y, R10; \
    MOVQ 24+x, R11;  ADCXQ 24+y, R11; \
    CMOVQCS AX, DX ; \
    XORL  AX,  AX; \
    ADCXQ DX,  R8; \
    ADCXQ AX,  R9;  MOVQ  R9,  8+z; \
    ADCXQ AX, R10;  MOVQ R10, 16+z; \
    ADCXQ AX, R11;  MOVQ R11, 24+z; \
    MOVL $38,  DX; \
    CMOVQCS DX, AX; \
    ADDQ  AX,  R8;  MOVQ  R8,  0+z;

// subtraction subtracts y from x and stores in z
// Uses: AX, DX, R8-R11, FLAGS
// Instr: x86_64, cmov
#define subtraction(z,x,y) \
    MOVL   $38,  AX; \
    MOVQ  0+x,  R8;  SUBQ  0+y,  R8; \
    MOVQ  8+x,  R9;  SBBQ  8+y,  R9; \
    MOVQ 16+x, R10;  SBBQ 16+y, R10; \
    MOVQ 24+x, R11;  SBBQ 24+y, R11; \
    MOVL $0, DX; \
    CMOVQCS AX, DX; \
    SUBQ  DX,  R8; \
    SBBQ  $0,  R9;  MOVQ  R9,  8+z; \
    SBBQ  $0, R10;  MOVQ R10, 16+z; \
    SBBQ  $0, R11;  MOVQ R11, 24+z; \
    MOVL  $0,  DX; \
    CMOVQCS AX, DX; \
    SUBQ  DX,  R8;  MOVQ  R8,  0+z;

// integerMulAdx multiplies x and y and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64, bmi2, adx
#define integerMulAdx(z,x,y) \
    MOVL    $0,R15; \
    MOVQ   0+y, DX;       XORL  AX,  AX; \
    MULXQ  0+x, AX,  R8;  MOVQ  AX, 0+z; \
    MULXQ  8+x, AX,  R9;  ADCXQ AX,  R8; \
    MULXQ 16+x, AX, R10;  ADCXQ AX,  R9; \
    MULXQ 24+x, AX, R11;  ADCXQ AX, R10; \
    MOVL $0, AX;;;;;;;;;  ADCXQ AX, R11; \
    MOVQ   8+y, DX;       XORL   AX,  AX; \
    MULXQ  0+x, AX, R12;  ADCXQ  R8,  AX;  MOVQ  AX,  8+z; \
    MULXQ  8+x, AX, R13;  ADCXQ  R9, R12;  ADOXQ AX, R12; \
    MULXQ 16+x, AX, R14;  ADCXQ R10, R13;  ADOXQ AX, R13; \
    MULXQ 24+x, AX, R15;  ADCXQ R11, R14;  ADOXQ AX, R14; \
    MOVL $0, AX;;;;;;;;;  ADCXQ  AX, R15;  ADOXQ AX, R15; \
    MOVQ  16+y, DX;       XORL   AX,  AX; \
    MULXQ  0+x, AX,  R8;  ADCXQ R12,  AX;  MOVQ  AX, 16+z; \
    MULXQ  8+x, AX,  R9;  ADCXQ R13,  R8;  ADOXQ AX,  R8; \
    MULXQ 16+x, AX, R10;  ADCXQ R14,  R9;  ADOXQ AX,  R9; \
    MULXQ 24+x, AX, R11;  ADCXQ R15, R10;  ADOXQ AX, R10; \
    MOVL $0, AX;;;;;;;;;  ADCXQ  AX, R11;  ADOXQ AX, R11; \
    MOVQ  24+y, DX;       XORL   AX,  AX; \
    MULXQ  0+x, AX, R12;  ADCXQ  R8,  AX;  MOVQ  AX, 24+z; \
    MULXQ  8+x, AX, R13;  ADCXQ  R9, R12;  ADOXQ AX, R12;  MOVQ R12, 32+z; \
    MULXQ 16+x, AX, R14;  ADCXQ R10, R13;  ADOXQ AX, R13;  MOVQ R13, 40+z; \
    MULXQ 24+x, AX, R15;  ADCXQ R11, R14;  ADOXQ AX, R14;  MOVQ R14, 48+z; \
    MOVL $0, AX;;;;;;;;;  ADCXQ  AX, R15;  ADOXQ AX, R15;  MOVQ R15, 56+z;

// integerMulLeg multiplies x and y and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64
#define integerMulLeg(z,x,y) \
    MOVQ  0+y, R8; \
    MOVQ  0+x, AX; MULQ R8; MOVQ AX, 0+z; MOVQ DX, R15; \
    MOVQ  8+x, AX; MULQ R8; MOVQ AX, R13; MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R8; MOVQ AX, R14; MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R8; \
    ADDQ R13, R15; \
    ADCQ R14, R10;  MOVQ R10, 16+z; \
    ADCQ  AX, R11;  MOVQ R11, 24+z; \
    ADCQ  $0,  DX;  MOVQ  DX, 32+z; \
    MOVQ  8+y, R8; \
    MOVQ  0+x, AX; MULQ R8; MOVQ AX, R12; MOVQ DX,  R9; \
    MOVQ  8+x, AX; MULQ R8; MOVQ AX, R13; MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R8; MOVQ AX, R14; MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R8; \
    ADDQ R12, R15; MOVQ R15,  8+z; \
    ADCQ R13,  R9; \
    ADCQ R14, R10; \
    ADCQ  AX, R11; \
    ADCQ  $0,  DX; \
    ADCQ 16+z,  R9;  MOVQ  R9,  R15; \
    ADCQ 24+z, R10;  MOVQ R10, 24+z; \
    ADCQ 32+z, R11;  MOVQ R11, 32+z; \
    ADCQ   $0,  DX;  MOVQ  DX, 40+z; \
    MOVQ 16+y, R8; \
    MOVQ  0+x, AX; MULQ R8; MOVQ AX, R12; MOVQ DX,  R9; \
    MOVQ  8+x, AX; MULQ R8; MOVQ AX, R13; MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R8; MOVQ AX, R14; MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R8; \
    ADDQ R12, R15;  MOVQ R15, 16+z; \
    ADCQ R13,  R9; \
    ADCQ R14, R10; \
    ADCQ  AX, R11; \
    ADCQ  $0,  DX; \
    ADCQ 24+z,  R9;  MOVQ  R9,  R15; \
    ADCQ 32+z, R10;  MOVQ R10, 32+z; \
    ADCQ 40+z, R11;  MOVQ R11, 40+z; \
    ADCQ   $0,  DX;  MOVQ  DX, 48+z; \
    MOVQ 24+y, R8; \
    MOVQ  0+x, AX; MULQ R8; MOVQ AX, R12; MOVQ DX,  R9; \
    MOVQ  8+x, AX; MULQ R8; MOVQ AX, R13; MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R8; MOVQ AX, R14; MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R8; \
    ADDQ R12, R15; MOVQ R15, 24+z; \
    ADCQ R13,  R9; \
    ADCQ R14, R10; \
    ADCQ  AX, R11; \
    ADCQ  $0,  DX; \
    ADCQ 32+z,  R9;  MOVQ  R9, 32+z; \
    ADCQ 40+z, R10;  MOVQ R10, 40+z; \
    ADCQ 48+z, R11;  MOVQ R11, 48+z; \
    ADCQ   $0,  DX;  MOVQ  DX, 56+z;

// integerSqrLeg squares x and stores in z
// Uses: AX, CX, DX, R8-R15, FLAGS
// Instr: x86_64
#define integerSqrLeg(z,x) \
    MOVQ  0+x, R8; \
    MOVQ  8+x, AX; MULQ R8; MOVQ AX,  R9; MOVQ DX, R10; /* A[0]*A[1] */ \
    MOVQ 16+x, AX; MULQ R8; MOVQ AX, R14; MOVQ DX, R11; /* A[0]*A[2] */ \
    MOVQ 24+x, AX; MULQ R8; MOVQ AX, R15; MOVQ DX, R12; /* A[0]*A[3] */ \
    MOVQ 24+x, R8; \
    MOVQ  8+x, AX; MULQ R8; MOVQ AX,  CX; MOVQ DX, R13; /* A[3]*A[1] */ \
    MOVQ 16+x, AX; MULQ R8; /* A[3]*A[2] */ \
    \
    ADDQ R14, R10;\
    ADCQ R15, R11; MOVL $0, R15;\
    ADCQ  CX, R12;\
    ADCQ  AX, R13;\
    ADCQ  $0,  DX; MOVQ DX, R14;\
    MOVQ 8+x, AX; MULQ 16+x;\
    \
    ADDQ AX, R11;\
    ADCQ DX, R12;\
    ADCQ $0, R13;\
    ADCQ $0, R14;\
    ADCQ $0, R15;\
    \
    SHLQ $1, R14, R15; MOVQ R15, 56+z;\
    SHLQ $1, R13, R14; MOVQ R14, 48+z;\
    SHLQ $1, R12, R13; MOVQ R13, 40+z;\
    SHLQ $1, R11, R12; MOVQ R12, 32+z;\
    SHLQ $1, R10, R11; MOVQ R11, 24+z;\
    SHLQ $1,  R9, R10; MOVQ R10, 16+z;\
    SHLQ $1,  R9;      MOVQ  R9,  8+z;\
    \
    MOVQ  0+x,AX; MULQ AX; MOVQ AX, 0+z; MOVQ DX,  R9;\
    MOVQ  8+x,AX; MULQ AX; MOVQ AX, R10; MOVQ DX, R11;\
    MOVQ 16+x,AX; MULQ AX; MOVQ AX, R12; MOVQ DX, R13;\
    MOVQ 24+x,AX; MULQ AX; MOVQ AX, R14; MOVQ DX, R15;\
    \
    ADDQ  8+z,  R9; MOVQ  R9,  8+z;\
    ADCQ 16+z, R10; MOVQ R10, 16+z;\
    ADCQ 24+z, R11; MOVQ R11, 24+z;\
    ADCQ 32+z, R12; MOVQ R12, 32+z;\
    ADCQ 40+z, R13; MOVQ R13, 40+z;\
    ADCQ 48+z, R14; MOVQ R14, 48+z;\
    ADCQ 56+z, R15; MOVQ R15, 56+z;

// integerSqrAdx squares x and stores in z
// Uses: AX, CX, DX, R8-R15, FLAGS
// Instr: x86_64, bmi2, adx
#define integerSqrAdx(z,x) \
    MOVQ   0+x,  DX; /* A[0] */ \
    MULXQ  8+x,  R8, R14; /* A[1]*A[0] */  XORL  R15, R15; \
    MULXQ 16+x,  R9, R10; /* A[2]*A[0] */  ADCXQ R14,  R9; \
    MULXQ 24+x,  AX,  CX; /* A[3]*A[0] */  ADCXQ  AX, R10; \
    MOVQ  24+x,  DX; /* A[3] */ \
    MULXQ  8+x, R11, R12; /* A[1]*A[3] */  ADCXQ  CX, R11; \
    MULXQ 16+x,  AX, R13; /* A[2]*A[3] */  ADCXQ  AX, R12; \
    MOVQ   8+x,  DX; /* A[1] */            ADCXQ R15, R13; \
    MULXQ 16+x,  AX,  CX; /* A[2]*A[1] */  MOVL   $0, R14; \
    ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;  ADCXQ R15, R14; \
    XORL  R15, R15; \
    ADOXQ  AX, R10;  ADCXQ  R8,  R8; \
    ADOXQ  CX, R11;  ADCXQ  R9,  R9; \
    ADOXQ R15, R12;  ADCXQ R10, R10; \
    ADOXQ R15, R13;  ADCXQ R11, R11; \
    ADOXQ R15, R14;  ADCXQ R12, R12; \
    ;;;;;;;;;;;;;;;  ADCXQ R13, R13; \
    ;;;;;;;;;;;;;;;  ADCXQ R14, R14; \
    MOVQ  0+x, DX;  MULXQ DX, AX, CX; /* A[0]^2 */ \
    ;;;;;;;;;;;;;;;  MOVQ  AX,  0+z; \
    ADDQ CX,  R8;    MOVQ  R8,  8+z; \
    MOVQ  8+x, DX;  MULXQ DX, AX, CX; /* A[1]^2 */ \
    ADCQ AX,  R9;    MOVQ  R9, 16+z; \
    ADCQ CX, R10;    MOVQ R10, 24+z; \
    MOVQ 16+x, DX;  MULXQ DX, AX, CX; /* A[2]^2 */ \
    ADCQ AX, R11;    MOVQ R11, 32+z; \
    ADCQ CX, R12;    MOVQ R12, 40+z; \
    MOVQ 24+x, DX;  MULXQ DX, AX, CX; /* A[3]^2 */ \
    ADCQ AX, R13;    MOVQ R13, 48+z; \
    ADCQ CX, R14;    MOVQ R14, 56+z;

// reduceFromDouble finds z congruent to x modulo p such that 0<z<2^256
// Uses: AX, DX, R8-R13, FLAGS
// Instr: x86_64
#define reduceFromDoubleLeg(z,x) \
    /* 2*C = 38 = 2^256 */ \
    MOVL $38, AX; MULQ 32+x; MOVQ AX,  R8; MOVQ DX,  R9; /* C*C[4] */ \
    MOVL $38, AX; MULQ 40+x; MOVQ AX, R12; MOVQ DX, R10; /* C*C[5] */ \
    MOVL $38, AX; MULQ 48+x; MOVQ AX, R13; MOVQ DX, R11; /* C*C[6] */ \
    MOVL $38, AX; MULQ 56+x; /* C*C[7] */ \
    ADDQ R12,  R9; \
    ADCQ R13, R10; \
    ADCQ  AX, R11; \
    ADCQ  $0,  DX; \
    ADDQ  0+x,  R8; \
    ADCQ  8+x,  R9; \
    ADCQ 16+x, R10; \
    ADCQ 24+x, R11; \
    ADCQ    $0, DX; \
    MOVL $38, AX; \
    IMULQ AX, DX; /* C*C[4], CF=0, OF=0 */ \
    ADDQ DX,  R8; \
    ADCQ $0,  R9; MOVQ  R9,  8+z; \
    ADCQ $0, R10; MOVQ R10, 16+z; \
    ADCQ $0, R11; MOVQ R11, 24+z; \
    MOVL $0,  DX; \
    CMOVQCS AX, DX; \
    ADDQ DX,  R8; MOVQ  R8,  0+z;

// reduceFromDoubleAdx finds z congruent to x modulo p such that 0<z<2^256
// Uses: AX, DX, R8-R13, FLAGS
// Instr: x86_64, bmi2, adx
#define reduceFromDoubleAdx(z,x) \
    MOVL    $38,  DX; /* 2*C = 38 = 2^256 */ \
    MULXQ 32+x,  R8, R10; /* C*C[4] */  XORL AX, AX;     ADOXQ  0+x,  R8; \
    MULXQ 40+x,  R9, R11; /* C*C[5] */  ADCXQ R10,  R9;  ADOXQ  8+x,  R9; \
    MULXQ 48+x, R10, R13; /* C*C[6] */  ADCXQ R11, R10;  ADOXQ 16+x, R10; \
    MULXQ 56+x, R11, R12; /* C*C[7] */  ADCXQ R13, R11;  ADOXQ 24+x, R11; \
    ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;  ADCXQ  AX, R12;  ADOXQ   AX, R12; \
    IMULQ  DX, R12; /* C*C[4], CF=0, OF=0 */ \
    ADCXQ R12, R8; \
    ADCXQ AX,  R9; MOVQ  R9,  8+z; \
    ADCXQ AX, R10; MOVQ R10, 16+z; \
    ADCXQ AX, R11; MOVQ R11, 24+z; \
    MOVL  $0, R12; \
    CMOVQCS DX, R12; \
    ADDQ R12,  R8; MOVQ  R8,  0+z;

// addSub calculates two operations: x,y = x+y,x-y
// Uses: AX, DX, R8-R15, FLAGS
#define addSub(x,y) \
    MOVL $38, AX; \
    XORL  DX, DX; \
    MOVQ  0+x,  R8;  MOVQ  R8, R12;  ADDQ  0+y,  R8; \
    MOVQ  8+x,  R9;  MOVQ  R9, R13;  ADCQ  8+y,  R9; \
    MOVQ 16+x, R10;  MOVQ R10, R14;  ADCQ 16+y, R10; \
    MOVQ 24+x, R11;  MOVQ R11, R15;  ADCQ 24+y, R11; \
    CMOVQCS AX, DX; \
    XORL AX,  AX; \
    ADDQ DX,  R8; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ $0, R11; \
    MOVL $38, DX; \
    CMOVQCS DX, AX; \
    ADDQ  AX, R8; \
    MOVL $38, AX; \
    SUBQ  0+y, R12; \
    SBBQ  8+y, R13; \
    SBBQ 16+y, R14; \
    SBBQ 24+y, R15; \
    MOVL $0, DX; \
    CMOVQCS AX, DX; \
    SUBQ DX, R12; \
    SBBQ $0, R13; \
    SBBQ $0, R14; \
    SBBQ $0, R15; \
    MOVL $0,  DX; \
    CMOVQCS AX, DX; \
    SUBQ DX, R12; \
    MOVQ  R8,  0+x; \
    MOVQ  R9,  8+x; \
    MOVQ R10, 16+x; \
    MOVQ R11, 24+x; \
    MOVQ R12,  0+y; \
    MOVQ R13,  8+y; \
    MOVQ R14, 16+y; \
    MOVQ R15, 24+y;
