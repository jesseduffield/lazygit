//go:build amd64 && !purego
// +build amd64,!purego

#include "textflag.h"

// Depends on circl/math/fp25519 package
#include "../../math/fp25519/fp_amd64.h"
#include "curve_amd64.h"

// CTE_A24 is (A+2)/4 from Curve25519
#define CTE_A24 121666

#define Size 32

// multiplyA24Leg multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R13, FLAGS
// Instr: x86_64, cmov
#define multiplyA24Leg(z,x) \
    MOVL $CTE_A24, AX; MULQ  0+x; MOVQ AX,  R8; MOVQ DX,  R9; \
    MOVL $CTE_A24, AX; MULQ  8+x; MOVQ AX, R12; MOVQ DX, R10; \
    MOVL $CTE_A24, AX; MULQ 16+x; MOVQ AX, R13; MOVQ DX, R11; \
    MOVL $CTE_A24, AX; MULQ 24+x; \
    ADDQ R12,  R9; \
    ADCQ R13, R10; \
    ADCQ  AX, R11; \
    ADCQ  $0,  DX; \
    MOVL $38,  AX; /* 2*C = 38 = 2^256 MOD 2^255-19*/ \
    IMULQ AX, DX; \
    ADDQ DX, R8; \
    ADCQ $0,  R9;  MOVQ  R9,  8+z; \
    ADCQ $0, R10;  MOVQ R10, 16+z; \
    ADCQ $0, R11;  MOVQ R11, 24+z; \
    MOVQ $0, DX; \
    CMOVQCS AX, DX; \
    ADDQ DX, R8;  MOVQ  R8,   0+z;

// multiplyA24Adx multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R12, FLAGS
// Instr: x86_64, cmov, bmi2
#define multiplyA24Adx(z,x) \
    MOVQ  $CTE_A24, DX; \
    MULXQ  0+x,  R8, R10; \
    MULXQ  8+x,  R9, R11;  ADDQ R10,  R9; \
    MULXQ 16+x, R10,  AX;  ADCQ R11, R10; \
    MULXQ 24+x, R11, R12;  ADCQ  AX, R11; \
    ;;;;;;;;;;;;;;;;;;;;;  ADCQ  $0, R12; \
    MOVL $38,  DX; /* 2*C = 38 = 2^256 MOD 2^255-19*/ \
    IMULQ DX, R12; \
    ADDQ R12, R8; \
    ADCQ $0,  R9;  MOVQ  R9,  8+z; \
    ADCQ $0, R10;  MOVQ R10, 16+z; \
    ADCQ $0, R11;  MOVQ R11, 24+z; \
    MOVQ $0, R12; \
    CMOVQCS DX, R12; \
    ADDQ R12, R8;  MOVQ  R8,  0+z;

#define mulA24Legacy \
    multiplyA24Leg(0(DI),0(SI))
#define mulA24Bmi2Adx \
    multiplyA24Adx(0(DI),0(SI))

// func mulA24Amd64(z, x *fp255.Elt)
TEXT ·mulA24Amd64(SB),NOSPLIT,$0-16
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    CHECK_BMI2ADX(LMA24, mulA24Legacy, mulA24Bmi2Adx)


// func ladderStepAmd64(w *[5]fp255.Elt, b uint)
// ladderStepAmd64 calculates a point addition and doubling as follows:
// (x2,z2) = 2*(x2,z2) and (x3,z3) = (x2,z2)+(x3,z3) using as a difference (x1,-).
//  work  = (x1,x2,z2,x3,z3) are five fp255.Elt of 32 bytes.
//  stack = (t0,t1) are two fp.Elt of fp.Size bytes, and
//          (b0,b1) are two-double precision fp.Elt of 2*fp.Size bytes.
TEXT ·ladderStepAmd64(SB),NOSPLIT,$192-16
    // Parameters
    #define regWork DI
    #define regMove SI
    #define x1 0*Size(regWork)
    #define x2 1*Size(regWork)
    #define z2 2*Size(regWork)
    #define x3 3*Size(regWork)
    #define z3 4*Size(regWork)
    // Local variables
    #define t0 0*Size(SP)
    #define t1 1*Size(SP)
    #define b0 2*Size(SP)
    #define b1 4*Size(SP)
    MOVQ w+0(FP), regWork
    MOVQ b+8(FP), regMove
    CHECK_BMI2ADX(LLADSTEP, ladderStepLeg, ladderStepBmi2Adx)
    #undef regWork
    #undef regMove
    #undef x1
    #undef x2
    #undef z2
    #undef x3
    #undef z3
    #undef t0
    #undef t1
    #undef b0
    #undef b1

// func diffAddAmd64(w *[5]fp255.Elt, b uint)
// diffAddAmd64 calculates a differential point addition using a precomputed point.
// (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
//    w    = (mu,x1,z1,x2,z2) are five fp.Elt, and
//   stack = (b0,b1) are two-double precision fp.Elt of 2*fp.Size bytes.
TEXT ·diffAddAmd64(SB),NOSPLIT,$128-16
    // Parameters
    #define regWork DI
    #define regSwap SI
    #define ui 0*Size(regWork)
    #define x1 1*Size(regWork)
    #define z1 2*Size(regWork)
    #define x2 3*Size(regWork)
    #define z2 4*Size(regWork)
    // Local variables
    #define b0 0*Size(SP)
    #define b1 2*Size(SP)
    MOVQ w+0(FP), regWork
    MOVQ b+8(FP), regSwap
    cswap(x1,x2,regSwap)
    cswap(z1,z2,regSwap)
    CHECK_BMI2ADX(LDIFADD, difAddLeg, difAddBmi2Adx)
    #undef regWork
    #undef regSwap
    #undef ui
    #undef x1
    #undef z1
    #undef x2
    #undef z2
    #undef b0
    #undef b1

// func doubleAmd64(x, z *fp255.Elt)
// doubleAmd64 calculates a point doubling (x1,z1) = 2*(x1,z1).
//  stack = (t0,t1) are two fp.Elt of fp.Size bytes, and
//          (b0,b1) are two-double precision fp.Elt of 2*fp.Size bytes.
TEXT ·doubleAmd64(SB),NOSPLIT,$192-16
    // Parameters
    #define x1 0(DI)
    #define z1 0(SI)
    // Local variables
    #define t0 0*Size(SP)
    #define t1 1*Size(SP)
    #define b0 2*Size(SP)
    #define b1 4*Size(SP)
    MOVQ x+0(FP), DI
    MOVQ z+8(FP), SI
    CHECK_BMI2ADX(LDOUB,doubleLeg,doubleBmi2Adx)
    #undef x1
    #undef z1
    #undef t0
    #undef t1
    #undef b0
    #undef b1
