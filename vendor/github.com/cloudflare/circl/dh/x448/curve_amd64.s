//go:build amd64 && !purego
// +build amd64,!purego

#include "textflag.h"

// Depends on circl/math/fp448 package
#include "../../math/fp448/fp_amd64.h"
#include "curve_amd64.h"

// CTE_A24 is (A+2)/4 from Curve448
#define CTE_A24 39082

#define Size 56

// multiplyA24Leg multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R15, FLAGS
// Instr: x86_64, cmov, adx
#define multiplyA24Leg(z,x) \
    MOVQ $CTE_A24, R15; \
    MOVQ  0+x, AX; MULQ R15; MOVQ AX,  R8; ;;;;;;;;;;;;  MOVQ DX,  R9; \
    MOVQ  8+x, AX; MULQ R15; ADDQ AX,  R9; ADCQ $0, DX;  MOVQ DX, R10; \
    MOVQ 16+x, AX; MULQ R15; ADDQ AX, R10; ADCQ $0, DX;  MOVQ DX, R11; \
    MOVQ 24+x, AX; MULQ R15; ADDQ AX, R11; ADCQ $0, DX;  MOVQ DX, R12; \
    MOVQ 32+x, AX; MULQ R15; ADDQ AX, R12; ADCQ $0, DX;  MOVQ DX, R13; \
    MOVQ 40+x, AX; MULQ R15; ADDQ AX, R13; ADCQ $0, DX;  MOVQ DX, R14; \
    MOVQ 48+x, AX; MULQ R15; ADDQ AX, R14; ADCQ $0, DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; MOVQ $0, DX; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    ADCQ $0,  DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    MOVQ  R8,  0+z; \
    MOVQ  R9,  8+z; \
    MOVQ R10, 16+z; \
    MOVQ R11, 24+z; \
    MOVQ R12, 32+z; \
    MOVQ R13, 40+z; \
    MOVQ R14, 48+z;

// multiplyA24Adx multiplies x times CTE_A24 and stores in z
// Uses: AX, DX, R8-R14, FLAGS
// Instr: x86_64, bmi2
#define multiplyA24Adx(z,x) \
    MOVQ $CTE_A24, DX; \
    MULXQ  0+x, R8,  R9; \
    MULXQ  8+x, AX, R10;  ADDQ AX,  R9; \
    MULXQ 16+x, AX, R11;  ADCQ AX, R10; \
    MULXQ 24+x, AX, R12;  ADCQ AX, R11; \
    MULXQ 32+x, AX, R13;  ADCQ AX, R12; \
    MULXQ 40+x, AX, R14;  ADCQ AX, R13; \
    MULXQ 48+x, AX,  DX;  ADCQ AX, R14; \
    ;;;;;;;;;;;;;;;;;;;;  ADCQ $0,  DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; MOVQ $0, DX; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    ADCQ $0,  DX; \
    MOVQ DX,  AX; \
    SHLQ $32, AX; \
    ADDQ DX,  R8; \
    ADCQ $0,  R9; \
    ADCQ $0, R10; \
    ADCQ AX, R11; \
    ADCQ $0, R12; \
    ADCQ $0, R13; \
    ADCQ $0, R14; \
    MOVQ  R8,  0+z; \
    MOVQ  R9,  8+z; \
    MOVQ R10, 16+z; \
    MOVQ R11, 24+z; \
    MOVQ R12, 32+z; \
    MOVQ R13, 40+z; \
    MOVQ R14, 48+z;

#define mulA24Legacy \
    multiplyA24Leg(0(DI),0(SI))
#define mulA24Bmi2Adx \
    multiplyA24Adx(0(DI),0(SI))

// func mulA24Amd64(z, x *fp448.Elt)
TEXT ·mulA24Amd64(SB),NOSPLIT,$0-16
    MOVQ z+0(FP), DI
    MOVQ x+8(FP), SI
    CHECK_BMI2ADX(LMA24, mulA24Legacy, mulA24Bmi2Adx)

// func ladderStepAmd64(w *[5]fp448.Elt, b uint)
// ladderStepAmd64 calculates a point addition and doubling as follows:
// (x2,z2) = 2*(x2,z2) and (x3,z3) = (x2,z2)+(x3,z3) using as a difference (x1,-).
//    w    = {x1,x2,z2,x3,z4} are five fp255.Elt of 56 bytes.
//  stack  = (t0,t1) are two fp.Elt of fp.Size bytes, and
//           (b0,b1) are two-double precision fp.Elt of 2*fp.Size bytes.
TEXT ·ladderStepAmd64(SB),NOSPLIT,$336-16
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

// func diffAddAmd64(work *[5]fp.Elt, swap uint)
// diffAddAmd64 calculates a differential point addition using a precomputed point.
// (x1,z1) = (x1,z1)+(mu) using a difference point (x2,z2)
//    work = {mu,x1,z1,x2,z2} are five fp448.Elt of 56 bytes, and
//   stack = (b0,b1) are two-double precision fp.Elt of 2*fp.Size bytes.
// This is Equation 7 at https://eprint.iacr.org/2017/264.
TEXT ·diffAddAmd64(SB),NOSPLIT,$224-16
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

// func doubleAmd64(x, z *fp448.Elt)
// doubleAmd64 calculates a point doubling (x1,z1) = 2*(x1,z1).
//  stack = (t0,t1) are two fp.Elt of fp.Size bytes, and
//          (b0,b1) are two-double precision fp.Elt of 2*fp.Size bytes.
TEXT ·doubleAmd64(SB),NOSPLIT,$336-16
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
