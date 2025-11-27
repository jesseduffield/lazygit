// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.25

package blake2b

import "hash"

var _ hash.XOF = (*xof)(nil)
