#define ladderStepLeg          \
    addSub(x2,z2)              \
    addSub(x3,z3)              \
    integerMulLeg(b0,x2,z3)    \
    integerMulLeg(b1,x3,z2)    \
    reduceFromDoubleLeg(t0,b0) \
    reduceFromDoubleLeg(t1,b1) \
    addSub(t0,t1)              \
    cselect(x2,x3,regMove)     \
    cselect(z2,z3,regMove)     \
    integerSqrLeg(b0,t0)       \
    integerSqrLeg(b1,t1)       \
    reduceFromDoubleLeg(x3,b0) \
    reduceFromDoubleLeg(z3,b1) \
    integerMulLeg(b0,x1,z3)    \
    reduceFromDoubleLeg(z3,b0) \
    integerSqrLeg(b0,x2)       \
    integerSqrLeg(b1,z2)       \
    reduceFromDoubleLeg(x2,b0) \
    reduceFromDoubleLeg(z2,b1) \
    subtraction(t0,x2,z2)      \
    multiplyA24Leg(t1,t0)      \
    additionLeg(t1,t1,z2)      \
    integerMulLeg(b0,x2,z2)    \
    integerMulLeg(b1,t0,t1)    \
    reduceFromDoubleLeg(x2,b0) \
    reduceFromDoubleLeg(z2,b1)

#define ladderStepBmi2Adx      \
    addSub(x2,z2)              \
    addSub(x3,z3)              \
    integerMulAdx(b0,x2,z3)    \
    integerMulAdx(b1,x3,z2)    \
    reduceFromDoubleAdx(t0,b0) \
    reduceFromDoubleAdx(t1,b1) \
    addSub(t0,t1)              \
    cselect(x2,x3,regMove)     \
    cselect(z2,z3,regMove)     \
    integerSqrAdx(b0,t0)       \
    integerSqrAdx(b1,t1)       \
    reduceFromDoubleAdx(x3,b0) \
    reduceFromDoubleAdx(z3,b1) \
    integerMulAdx(b0,x1,z3)    \
    reduceFromDoubleAdx(z3,b0) \
    integerSqrAdx(b0,x2)       \
    integerSqrAdx(b1,z2)       \
    reduceFromDoubleAdx(x2,b0) \
    reduceFromDoubleAdx(z2,b1) \
    subtraction(t0,x2,z2)      \
    multiplyA24Adx(t1,t0)      \
    additionAdx(t1,t1,z2)      \
    integerMulAdx(b0,x2,z2)    \
    integerMulAdx(b1,t0,t1)    \
    reduceFromDoubleAdx(x2,b0) \
    reduceFromDoubleAdx(z2,b1)

#define difAddLeg              \
    addSub(x1,z1)              \
    integerMulLeg(b0,z1,ui)    \
    reduceFromDoubleLeg(z1,b0) \
    addSub(x1,z1)              \
    integerSqrLeg(b0,x1)       \
    integerSqrLeg(b1,z1)       \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1) \
    integerMulLeg(b0,x1,z2)    \
    integerMulLeg(b1,z1,x2)    \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1)

#define difAddBmi2Adx          \
    addSub(x1,z1)              \
    integerMulAdx(b0,z1,ui)    \
    reduceFromDoubleAdx(z1,b0) \
    addSub(x1,z1)              \
    integerSqrAdx(b0,x1)       \
    integerSqrAdx(b1,z1)       \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1) \
    integerMulAdx(b0,x1,z2)    \
    integerMulAdx(b1,z1,x2)    \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1)

#define doubleLeg              \
    addSub(x1,z1)              \
    integerSqrLeg(b0,x1)       \
    integerSqrLeg(b1,z1)       \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1) \
    subtraction(t0,x1,z1)      \
    multiplyA24Leg(t1,t0)      \
    additionLeg(t1,t1,z1)      \
    integerMulLeg(b0,x1,z1)    \
    integerMulLeg(b1,t0,t1)    \
    reduceFromDoubleLeg(x1,b0) \
    reduceFromDoubleLeg(z1,b1)

#define doubleBmi2Adx          \
    addSub(x1,z1)              \
    integerSqrAdx(b0,x1)       \
    integerSqrAdx(b1,z1)       \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1) \
    subtraction(t0,x1,z1)      \
    multiplyA24Adx(t1,t0)      \
    additionAdx(t1,t1,z1)      \
    integerMulAdx(b0,x1,z1)    \
    integerMulAdx(b1,t0,t1)    \
    reduceFromDoubleAdx(x1,b0) \
    reduceFromDoubleAdx(z1,b1)
