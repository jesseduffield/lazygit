// Copyright 2021 Peter Mattis.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License. See the AUTHORS file
// for names of contributors.

// Assembly to mimic runtime.getg.

//go:build (386 || amd64 || amd64p32 || arm || arm64 || s390x) && gc && go1.5
// +build 386 amd64 amd64p32 arm arm64 s390x
// +build gc
// +build go1.5

#include "textflag.h"

// func getg() *g
TEXT ·getg(SB),NOSPLIT,$0-8
#ifdef GOARCH_386
	MOVL (TLS), AX
	MOVL AX, ret+0(FP)
#endif
#ifdef GOARCH_amd64
	MOVQ (TLS), AX
	MOVQ AX, ret+0(FP)
#endif
#ifdef GOARCH_arm
	MOVW g, ret+0(FP)
#endif
#ifdef GOARCH_arm64
	MOVD g, ret+0(FP)
#endif
#ifdef GOARCH_s390x
	MOVD g, ret+0(FP)
#endif
	RET
