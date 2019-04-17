// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lzma

import (
	"bytes"
	"io"
	"math/rand"
	"strings"
	"testing"

	"github.com/ulikunitz/xz/internal/randtxt"
)

func TestBinTree_Find(t *testing.T) {
	bt, err := newBinTree(30)
	if err != nil {
		t.Fatal(err)
	}
	const s = "Klopp feiert mit Liverpool seinen hoechsten SiegSieg"
	n, err := io.WriteString(bt, s)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if n != len(s) {
		t.Fatalf("WriteString returned %d; want %d", n, len(s))
	}

	/* dump info writes the complete tree
	if err = bt.dump(os.Stdout); err != nil {
		t.Fatalf("bt.dump error %s", err)
	}
	*/

	tests := []string{"Sieg", "Sieb", "Simu"}
	for _, c := range tests {
		x := xval([]byte(c))
		a, b := bt.search(bt.root, x)
		t.Logf("%q: a, b == %d, %d", c, a, b)
	}
}

func TestBinTree_PredSucc(t *testing.T) {
	bt, err := newBinTree(30)
	if err != nil {
		t.Fatal(err)
	}
	const s = "Klopp feiert mit Liverpool seinen hoechsten Sieg."
	n, err := io.WriteString(bt, s)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if n != len(s) {
		t.Fatalf("WriteString returned %d; want %d", n, len(s))
	}
	for v := bt.min(bt.root); v != null; v = bt.succ(v) {
		t.Log(dumpX(bt.node[v].x))
	}
	t.Log("")
	for v := bt.max(bt.root); v != null; v = bt.pred(v) {
		t.Log(dumpX(bt.node[v].x))
	}
}

func TestBinTree_Cycle(t *testing.T) {
	buf := new(bytes.Buffer)
	w, err := Writer2Config{
		DictCap: 4096,
		Matcher: BinaryTree,
	}.NewWriter2(buf)
	if err != nil {
		t.Fatalf("NewWriter error %s", err)
	}
	// const txtlen = 1024
	const txtlen = 10000
	io.CopyN(buf, randtxt.NewReader(rand.NewSource(42)), txtlen)
	txt := buf.String()
	buf.Reset()
	n, err := io.Copy(w, strings.NewReader(txt))
	if err != nil {
		t.Fatalf("Compressing copy error %s", err)
	}
	if n != txtlen {
		t.Fatalf("Compressing data length %d; want %d", n, txtlen)
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	t.Logf("buf.Len() %d", buf.Len())
	r, err := Reader2Config{DictCap: 4096}.NewReader2(buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	out := new(bytes.Buffer)
	n, err = io.Copy(out, r)
	if err != nil {
		t.Fatalf("Decompressing copy error %s after %d bytes", err, n)
	}
	if n != txtlen {
		t.Fatalf("Decompression data length %d; want %d", n, txtlen)
	}
	if txt != out.String() {
		t.Fatal("decompressed data differs from original")
	}
}
