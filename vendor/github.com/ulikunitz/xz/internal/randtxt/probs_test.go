// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package randtxt

import (
	"bufio"
	"io"
	"math/rand"
	"testing"
)

func TestReader(t *testing.T) {
	lr := io.LimitReader(NewReader(rand.NewSource(13)), 195)
	pretty := NewGroupReader(lr)
	scanner := bufio.NewScanner(pretty)
	for scanner.Scan() {
		t.Log(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scanner error %s", err)
	}
}

func TestComap(t *testing.T) {
	prs := cmap["TH"]
	for _, p := range prs[3:6] {
		t.Logf("%v", p)
	}
	p := 0.2
	x := cmap.trigram("TH", p)
	if x != "THE" {
		t.Fatalf("cmap.trigram(%q, %.1f) returned %q; want %q",
			"TH", p, x, "THE")
	}
}
