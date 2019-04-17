package object

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing"

	. "gopkg.in/check.v1"
	"gopkg.in/src-d/go-git-fixtures.v3"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type SuiteCommit struct {
	BaseObjectsSuite
	Commit *Commit
}

var _ = Suite(&SuiteCommit{})

func (s *SuiteCommit) SetUpSuite(c *C) {
	s.BaseObjectsSuite.SetUpSuite(c)

	hash := plumbing.NewHash("1669dce138d9b841a518c64b10914d88f5e488ea")

	s.Commit = s.commit(c, hash)
}

func (s *SuiteCommit) TestDecodeNonCommit(c *C) {
	hash := plumbing.NewHash("9a48f23120e880dfbe41f7c9b7b708e9ee62a492")
	blob, err := s.Storer.EncodedObject(plumbing.AnyObject, hash)
	c.Assert(err, IsNil)

	commit := &Commit{}
	err = commit.Decode(blob)
	c.Assert(err, Equals, ErrUnsupportedObject)
}

func (s *SuiteCommit) TestType(c *C) {
	c.Assert(s.Commit.Type(), Equals, plumbing.CommitObject)
}

func (s *SuiteCommit) TestTree(c *C) {
	tree, err := s.Commit.Tree()
	c.Assert(err, IsNil)
	c.Assert(tree.ID().String(), Equals, "eba74343e2f15d62adedfd8c883ee0262b5c8021")
}

func (s *SuiteCommit) TestParents(c *C) {
	expected := []string{
		"35e85108805c84807bc66a02d91535e1e24b38b9",
		"a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69",
	}

	var output []string
	i := s.Commit.Parents()
	err := i.ForEach(func(commit *Commit) error {
		output = append(output, commit.ID().String())
		return nil
	})

	c.Assert(err, IsNil)
	c.Assert(output, DeepEquals, expected)

	i.Close()
}

func (s *SuiteCommit) TestParent(c *C) {
	commit, err := s.Commit.Parent(1)
	c.Assert(err, IsNil)
	c.Assert(commit.Hash.String(), Equals, "a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69")
}

func (s *SuiteCommit) TestParentNotFound(c *C) {
	commit, err := s.Commit.Parent(42)
	c.Assert(err, Equals, ErrParentNotFound)
	c.Assert(commit, IsNil)
}

func (s *SuiteCommit) TestPatch(c *C) {
	from := s.commit(c, plumbing.NewHash("918c48b83bd081e863dbe1b80f8998f058cd8294"))
	to := s.commit(c, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))

	patch, err := from.Patch(to)
	c.Assert(err, IsNil)

	buf := bytes.NewBuffer(nil)
	err = patch.Encode(buf)
	c.Assert(err, IsNil)

	c.Assert(buf.String(), Equals, `diff --git a/vendor/foo.go b/vendor/foo.go
new file mode 100644
index 0000000000000000000000000000000000000000..9dea2395f5403188298c1dabe8bdafe562c491e3
--- /dev/null
+++ b/vendor/foo.go
@@ -0,0 +1,7 @@
+package main
+
+import "fmt"
+
+func main() {
+	fmt.Println("Hello, playground")
+}
`)
	c.Assert(buf.String(), Equals, patch.String())

	from = s.commit(c, plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"))
	to = s.commit(c, plumbing.NewHash("35e85108805c84807bc66a02d91535e1e24b38b9"))

	patch, err = from.Patch(to)
	c.Assert(err, IsNil)

	buf.Reset()
	err = patch.Encode(buf)
	c.Assert(err, IsNil)

	c.Assert(buf.String(), Equals, `diff --git a/CHANGELOG b/CHANGELOG
deleted file mode 100644
index d3ff53e0564a9f87d8e84b6e28e5060e517008aa..0000000000000000000000000000000000000000
--- a/CHANGELOG
+++ /dev/null
@@ -1 +0,0 @@
-Initial changelog
diff --git a/binary.jpg b/binary.jpg
new file mode 100644
index 0000000000000000000000000000000000000000..d5c0f4ab811897cadf03aec358ae60d21f91c50d
Binary files /dev/null and b/binary.jpg differ
`)

	c.Assert(buf.String(), Equals, patch.String())
}

func (s *SuiteCommit) TestPatchContext(c *C) {
	from := s.commit(c, plumbing.NewHash("918c48b83bd081e863dbe1b80f8998f058cd8294"))
	to := s.commit(c, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))

	patch, err := from.PatchContext(context.Background(), to)
	c.Assert(err, IsNil)

	buf := bytes.NewBuffer(nil)
	err = patch.Encode(buf)
	c.Assert(err, IsNil)

	c.Assert(buf.String(), Equals, `diff --git a/vendor/foo.go b/vendor/foo.go
new file mode 100644
index 0000000000000000000000000000000000000000..9dea2395f5403188298c1dabe8bdafe562c491e3
--- /dev/null
+++ b/vendor/foo.go
@@ -0,0 +1,7 @@
+package main
+
+import "fmt"
+
+func main() {
+	fmt.Println("Hello, playground")
+}
`)
	c.Assert(buf.String(), Equals, patch.String())

	from = s.commit(c, plumbing.NewHash("b8e471f58bcbca63b07bda20e428190409c2db47"))
	to = s.commit(c, plumbing.NewHash("35e85108805c84807bc66a02d91535e1e24b38b9"))

	patch, err = from.PatchContext(context.Background(), to)
	c.Assert(err, IsNil)

	buf.Reset()
	err = patch.Encode(buf)
	c.Assert(err, IsNil)

	c.Assert(buf.String(), Equals, `diff --git a/CHANGELOG b/CHANGELOG
deleted file mode 100644
index d3ff53e0564a9f87d8e84b6e28e5060e517008aa..0000000000000000000000000000000000000000
--- a/CHANGELOG
+++ /dev/null
@@ -1 +0,0 @@
-Initial changelog
diff --git a/binary.jpg b/binary.jpg
new file mode 100644
index 0000000000000000000000000000000000000000..d5c0f4ab811897cadf03aec358ae60d21f91c50d
Binary files /dev/null and b/binary.jpg differ
`)

	c.Assert(buf.String(), Equals, patch.String())
}

func (s *SuiteCommit) TestCommitEncodeDecodeIdempotent(c *C) {
	ts, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05-07:00")
	c.Assert(err, IsNil)
	commits := []*Commit{
		{
			Author:       Signature{Name: "Foo", Email: "foo@example.local", When: ts},
			Committer:    Signature{Name: "Bar", Email: "bar@example.local", When: ts},
			Message:      "Message\n\nFoo\nBar\nWith trailing blank lines\n\n",
			TreeHash:     plumbing.NewHash("f000000000000000000000000000000000000001"),
			ParentHashes: []plumbing.Hash{plumbing.NewHash("f000000000000000000000000000000000000002")},
		},
		{
			Author:    Signature{Name: "Foo", Email: "foo@example.local", When: ts},
			Committer: Signature{Name: "Bar", Email: "bar@example.local", When: ts},
			Message:   "Message\n\nFoo\nBar\nWith no trailing blank lines",
			TreeHash:  plumbing.NewHash("0000000000000000000000000000000000000003"),
			ParentHashes: []plumbing.Hash{
				plumbing.NewHash("f000000000000000000000000000000000000004"),
				plumbing.NewHash("f000000000000000000000000000000000000005"),
				plumbing.NewHash("f000000000000000000000000000000000000006"),
				plumbing.NewHash("f000000000000000000000000000000000000007"),
			},
		},
	}
	for _, commit := range commits {
		obj := &plumbing.MemoryObject{}
		err = commit.Encode(obj)
		c.Assert(err, IsNil)
		newCommit := &Commit{}
		err = newCommit.Decode(obj)
		c.Assert(err, IsNil)
		commit.Hash = obj.Hash()
		c.Assert(newCommit, DeepEquals, commit)
	}
}

func (s *SuiteCommit) TestFile(c *C) {
	file, err := s.Commit.File("CHANGELOG")
	c.Assert(err, IsNil)
	c.Assert(file.Name, Equals, "CHANGELOG")
}

func (s *SuiteCommit) TestNumParents(c *C) {
	c.Assert(s.Commit.NumParents(), Equals, 2)
}

func (s *SuiteCommit) TestString(c *C) {
	c.Assert(s.Commit.String(), Equals, ""+
		"commit 1669dce138d9b841a518c64b10914d88f5e488ea\n"+
		"Author: Máximo Cuadros Ortiz <mcuadros@gmail.com>\n"+
		"Date:   Tue Mar 31 13:48:14 2015 +0200\n"+
		"\n"+
		"    Merge branch 'master' of github.com:tyba/git-fixture\n"+
		"\n",
	)
}

func (s *SuiteCommit) TestStringMultiLine(c *C) {
	hash := plumbing.NewHash("e7d896db87294e33ca3202e536d4d9bb16023db3")

	f := fixtures.ByURL("https://github.com/src-d/go-git.git").One()
	sto, err := filesystem.NewStorage(f.DotGit())
	c.Assert(err, IsNil)

	o, err := sto.EncodedObject(plumbing.CommitObject, hash)
	c.Assert(err, IsNil)
	commit, err := DecodeCommit(sto, o)
	c.Assert(err, IsNil)

	c.Assert(commit.String(), Equals, ""+
		"commit e7d896db87294e33ca3202e536d4d9bb16023db3\n"+
		"Author: Alberto Cortés <alberto@sourced.tech>\n"+
		"Date:   Wed Jan 27 11:13:49 2016 +0100\n"+
		"\n"+
		"    fix zlib invalid header error\n"+
		"\n"+
		"    The return value of reads to the packfile were being ignored, so zlib\n"+
		"    was getting invalid data on it read buffers.\n"+
		"\n",
	)
}

func (s *SuiteCommit) TestCommitIterNext(c *C) {
	i := s.Commit.Parents()

	commit, err := i.Next()
	c.Assert(err, IsNil)
	c.Assert(commit.ID().String(), Equals, "35e85108805c84807bc66a02d91535e1e24b38b9")

	commit, err = i.Next()
	c.Assert(err, IsNil)
	c.Assert(commit.ID().String(), Equals, "a5b8b09e2f8fcb0bb99d3ccb0958157b40890d69")

	commit, err = i.Next()
	c.Assert(err, Equals, io.EOF)
	c.Assert(commit, IsNil)
}

func (s *SuiteCommit) TestLongCommitMessageSerialization(c *C) {
	encoded := &plumbing.MemoryObject{}
	decoded := &Commit{}
	commit := *s.Commit

	longMessage := "my message: message\n\n" + strings.Repeat("test", 4096) + "\nOK"
	commit.Message = longMessage

	err := commit.Encode(encoded)
	c.Assert(err, IsNil)

	err = decoded.Decode(encoded)
	c.Assert(err, IsNil)
	c.Assert(decoded.Message, Equals, longMessage)
}

func (s *SuiteCommit) TestPGPSignatureSerialization(c *C) {
	encoded := &plumbing.MemoryObject{}
	decoded := &Commit{}
	commit := *s.Commit

	pgpsignature := `-----BEGIN PGP SIGNATURE-----

iQEcBAABAgAGBQJTZbQlAAoJEF0+sviABDDrZbQH/09PfE51KPVPlanr6q1v4/Ut
LQxfojUWiLQdg2ESJItkcuweYg+kc3HCyFejeDIBw9dpXt00rY26p05qrpnG+85b
hM1/PswpPLuBSr+oCIDj5GMC2r2iEKsfv2fJbNW8iWAXVLoWZRF8B0MfqX/YTMbm
ecorc4iXzQu7tupRihslbNkfvfciMnSDeSvzCpWAHl7h8Wj6hhqePmLm9lAYqnKp
8S5B/1SSQuEAjRZgI4IexpZoeKGVDptPHxLLS38fozsyi0QyDyzEgJxcJQVMXxVi
RUysgqjcpT8+iQM1PblGfHR4XAhuOqN5Fx06PSaFZhqvWFezJ28/CLyX5q+oIVk=
=EFTF
-----END PGP SIGNATURE-----
`
	commit.PGPSignature = pgpsignature

	err := commit.Encode(encoded)
	c.Assert(err, IsNil)

	err = decoded.Decode(encoded)
	c.Assert(err, IsNil)
	c.Assert(decoded.PGPSignature, Equals, pgpsignature)

	// signature in author name

	commit.PGPSignature = ""
	commit.Author.Name = beginpgp
	encoded = &plumbing.MemoryObject{}
	decoded = &Commit{}

	err = commit.Encode(encoded)
	c.Assert(err, IsNil)

	err = decoded.Decode(encoded)
	c.Assert(err, IsNil)
	c.Assert(decoded.PGPSignature, Equals, "")
	c.Assert(decoded.Author.Name, Equals, beginpgp)

	// broken signature

	commit.PGPSignature = beginpgp + "\n" +
		"some\n" +
		"trash\n" +
		endpgp +
		"text\n"
	encoded = &plumbing.MemoryObject{}
	decoded = &Commit{}

	err = commit.Encode(encoded)
	c.Assert(err, IsNil)

	err = decoded.Decode(encoded)
	c.Assert(err, IsNil)
	c.Assert(decoded.PGPSignature, Equals, commit.PGPSignature)
}

func (s *SuiteCommit) TestStat(c *C) {
	aCommit := s.commit(c, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))
	fileStats, err := aCommit.Stats()
	c.Assert(err, IsNil)

	c.Assert(fileStats[0].Name, Equals, "vendor/foo.go")
	c.Assert(fileStats[0].Addition, Equals, 7)
	c.Assert(fileStats[0].Deletion, Equals, 0)
	c.Assert(fileStats[0].String(), Equals, " vendor/foo.go | 7 +++++++\n")

	// Stats for another commit.
	aCommit = s.commit(c, plumbing.NewHash("918c48b83bd081e863dbe1b80f8998f058cd8294"))
	fileStats, err = aCommit.Stats()
	c.Assert(err, IsNil)

	c.Assert(fileStats[0].Name, Equals, "go/example.go")
	c.Assert(fileStats[0].Addition, Equals, 142)
	c.Assert(fileStats[0].Deletion, Equals, 0)
	c.Assert(fileStats[0].String(), Equals, " go/example.go | 142 +++++++++++++++++++++++++++++++++++++++++++++++++++++\n")

	c.Assert(fileStats[1].Name, Equals, "php/crappy.php")
	c.Assert(fileStats[1].Addition, Equals, 259)
	c.Assert(fileStats[1].Deletion, Equals, 0)
	c.Assert(fileStats[1].String(), Equals, " php/crappy.php | 259 ++++++++++++++++++++++++++++++++++++++++++++++++++++\n")
}

func (s *SuiteCommit) TestVerify(c *C) {
	ts := time.Unix(1511197315, 0)
	loc, _ := time.LoadLocation("Asia/Kolkata")
	commit := &Commit{
		Hash:      plumbing.NewHash("8a9cea36fe052711fbc42b86e1f99a4fa0065deb"),
		Author:    Signature{Name: "Sunny", Email: "me@darkowlzz.space", When: ts.In(loc)},
		Committer: Signature{Name: "Sunny", Email: "me@darkowlzz.space", When: ts.In(loc)},
		Message: `status: simplify template command selection
`,
		TreeHash:     plumbing.NewHash("6572ba6df4f1fb323c8aaa24ce07bca0648b161e"),
		ParentHashes: []plumbing.Hash{plumbing.NewHash("ede5f57ea1280a0065beec96d3e1a3453d010dbd")},
		PGPSignature: `
-----BEGIN PGP SIGNATURE-----

iQFHBAABCAAxFiEEoRt6IzxHaZkkUslhQyLeMqcmyU4FAloTCrsTHG1lQGRhcmtv
d2x6ei5zcGFjZQAKCRBDIt4ypybJTul5CADmVxB4kqlqRZ9fAcSU5LKva3GRXx0+
leX6vbzoyQztSWYgl7zALh4kB3a3t2C9EnnM6uehlgaORNigyMArCSY1ivWVviCT
BvldSVi8f8OvnqwbWX0I/5a8KmItthDf5WqZRFjhcRlY1AK5Bo2hUGVRq71euf8F
rE6wNhDoyBCEpftXuXbq8duD7D6qJ7QiOS4m5+ej1UCssS2WQ60yta7q57odduHY
+txqTKI8MQUpBgoTqh+V4lOkwQQxLiz7hIQ/ZYLUcnp6fan7/kY/G7YoLt9pOG1Y
vLzAWdidLH2P+EUOqlNMuVScHYWD1FZB0/L5LJ8no5pTowQd2Z+Nggxl
=0uC8
-----END PGP SIGNATURE-----
`,
	}

	armoredKeyRing := `
-----BEGIN PGP PUBLIC KEY BLOCK-----

mQENBFmtHgABCADnfThM7q8D4pgUub9jMppSpgFh3ev84g3Csc3yQUlszEOVgXmu
YiSWP1oAiWFQ8ahCydh3LT8TnEB2QvoRNiExUI5XlXFwVfKW3cpDu8gdhtufs90Q
NvpaHOgTqRf/texGEKwXi6fvS47fpyaQ9BKNdN52LeaaHzDDZkVsAFmroE+7MMvj
P4Mq8qDn2WcWnX9zheQKYrX6Cs48Tx80eehHor4f/XnuaP8DLmPQx7URdJ0Igckh
N+i91Qv2ujin8zxUwhkfus66EZS9lQ4qR9iVHs4WHOs3j7whsejd4VhajonilVHj
uqTtqHmpN/4njbIKb8q8uQkS26VQYoSYm2UvABEBAAG0GlN1bm55IDxtZUBkYXJr
b3dsenouc3BhY2U+iQFUBBMBCAA+FiEEoRt6IzxHaZkkUslhQyLeMqcmyU4FAlmt
HgACGwMFCQPCZwAFCwkIBwIGFQgJCgsCBBYCAwECHgECF4AACgkQQyLeMqcmyU7V
nAf+J5BYu26B2i+iwctOzDRFcPwCLka9cBwe5wcDvoF2qL8QRo8NPWBBH4zWHa/k
BthtGo1b89a53I2hnTwTQ0NOtAUNV+Vvu6nOHJd9Segsx3E1nM43bd2bUfGJ1eeO
jDOlOvtP4ozuV6Ej+0Ln2ouMOc87yAwbAzTfQ9axU6CKUbqy0/t2dW1jdKntGH+t
VPeFxJHL2gXjP89skCSPYA7yKqqyJRPFvC+7rde1OLdCmZi4VwghUiNbh3s1+xM3
gfr2ahsRDTN2SQzwuHu4y1EgZgPtuWfRxzHqduoRoSgfOfFr9H9Il3UMHf2Etleu
rif40YZJhge6STwsIycGh4wOiLkBDQRZrR4AAQgArpUvPdGC/W9X4AuZXrXEShvx
TqM4K2Jk9n0j+ABx87k9fm48qgtae7+TayMbb0i7kcbgnjltKbauTbyRbju/EJvN
CdIw76IPpjy6jUM37wG2QGLFo6Ku3x8/ZpNGGOZ8KMU258/EBqDlJQ/4g4kJ8D+m
9yOH0r6/Xpe/jOY2V8Jo9pdFTm+8eAsSyZF0Cl7drz603Pymq1IS2wrwQbdxQA/w
B75pQ5es7X34Ac7/9UZCwCPmZDAldnjHyw5dZgZe8XLrG84BIfbG0Hj8PjrFdF1D
Czt9bk+PbYAnLORW2oX1oedxVrNFo5UrbWgBSjA1ppbGFjwSDHFlyjuEuxqyFwAR
AQABiQE8BBgBCAAmFiEEoRt6IzxHaZkkUslhQyLeMqcmyU4FAlmtHgACGwwFCQPC
ZwAACgkQQyLeMqcmyU7ZBggArzc8UUVSjde987Vqnu/S5Cv8Qhz+UB7gAFyTW2iF
VYvB86r30H/NnfjvjCVkBE6FHCNHoxWVyDWmuxKviB7nkReHuwqniQHPgdJDcTKC
tBboeX2IYBLJbEvEJuz5NSvnvFuYkIpZHqySFaqdl/qu9XcmoPL5AmIzIFOeiNty
qT0ldkf3ru6yQQDDqBDpkfz4AzkpFnLYL59z6IbJDK2Hz7aKeSEeVOGiZLCjIZZV
uISZThYqh5zUkvF346OHLDqfDdgQ4RZriqd/DTtRJPlz2uL0QcEIjJuYCkG0UWgl
sYyf9RfOnw/KUFAQbdtvLx3ikODQC+D3KBtuKI9ISHQfgw==
=FPev
-----END PGP PUBLIC KEY BLOCK-----
`

	e, err := commit.Verify(armoredKeyRing)
	c.Assert(err, IsNil)

	_, ok := e.Identities["Sunny <me@darkowlzz.space>"]
	c.Assert(ok, Equals, true)
}

func (s *SuiteCommit) TestPatchCancel(c *C) {
	from := s.commit(c, plumbing.NewHash("918c48b83bd081e863dbe1b80f8998f058cd8294"))
	to := s.commit(c, plumbing.NewHash("6ecf0ef2c2dffb796033e5a02219af86ec6584e5"))

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	patch, err := from.PatchContext(ctx, to)
	c.Assert(patch, IsNil)
	c.Assert(err, ErrorMatches, "operation canceled")

}
