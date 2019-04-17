// Copyright 2014-2017 Ulrich Kunitz. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package randtxt supports the generation of random text using a
// trigram model for the English language.
package randtxt

import (
	"math"
	"math/rand"
	"sort"
)

// ngram stores an entry from the language model.
type ngram struct {
	s   string
	lgP float64
	lgQ float64
}

// ngrams represents a slice of ngram values and is used to represent a
// language model.
type ngrams []ngram

func (s ngrams) Len() int           { return len(s) }
func (s ngrams) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ngrams) Less(i, j int) bool { return s[i].s < s[j].s }

// Sorts the language model in the sequence of their ngrams.
func (s ngrams) Sort() { sort.Sort(s) }

// Search is looking for an ngram or the position where it would be
// inserted.
func (s ngrams) Search(g string) int {
	return sort.Search(len(s), func(k int) bool { return s[k].s >= g })
}

// prob represents a string, usually an ngram, and a probability value.
type prob struct {
	s string
	p float64
}

// probs is a slice of prob values that can be sorted and searched.
type probs []prob

func (s probs) Len() int           { return len(s) }
func (s probs) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s probs) Less(i, j int) bool { return s[i].s < s[j].s }

// SortByNgram sorts the probs slice by ngram, field s.
func (s probs) SortByNgram() { sort.Sort(s) }

// SortsByProb sorts the probs slice by probability, field p.
func (s probs) SortByProb() { sort.Sort(byProb{s}) }

// SearchNgram searches for an ngram or the position where it would be
// inserted.
func (s probs) SearchNgram(g string) int {
	return sort.Search(len(s), func(k int) bool { return s[k].s >= g })
}

// SearchProb searches ngrams for a specific probability or where it
// would be inserted.
func (s probs) SearchProb(p float64) int {
	return sort.Search(len(s), func(k int) bool { return s[k].p >= p })
}

// byProb is used to sort probs slice by probability, field p.
type byProb struct {
	probs
}

func (s byProb) Less(i, j int) bool {
	return s.probs[i].p < s.probs[j].p
}

// cdf can be used to setup a cumulative distribution function
// represented by a probs slice. We should have returned an actual
// function.
func cdf(n int, p func(i int) prob) probs {
	prs := make(probs, n)
	sum := 0.0
	for i := range prs {
		pr := p(i)
		sum += pr.p
		prs[i] = pr
	}
	q := 1.0 / sum
	x := 0.0
	for i, pr := range prs {
		x += pr.p * q
		if x > 1.0 {
			x = 1.0
		}
		prs[i].p = x
	}
	if !sort.IsSorted(byProb{prs}) {
		panic("cdf not sorted")
	}
	return prs
}

// pCDFOfLM converts a language model into a cumulative distribution
// function represented by probs.
func pCDFOfLM(lm ngrams) probs {
	return cdf(len(lm), func(i int) prob {
		return prob{lm[i].s, math.Exp2(lm[i].lgP)}
	})
}

// cCDF converts a ngrams slice into a cumulative distribution function
// using the conditional probability lgQ.
func cCDF(s ngrams) probs {
	return cdf(len(s), func(i int) prob {
		return prob{s[i].s, math.Exp2(s[i].lgQ)}
	})
}

// comap contains a map of conditional distribution function for the
// last character.
type comap map[string]probs

// comapOfLM converts a language model in a map of conditional
// distribution functions.
func comapOfLM(lm ngrams) comap {
	if !sort.IsSorted(lm) {
		panic("lm is not sorted")
	}
	m := make(comap, 26*26)
	for i := 0; i < len(lm); {
		j := i
		g := lm[i].s
		g2 := g[:2]
		z := g2 + "Z"
		i = lm.Search(z)
		if i >= len(lm) || lm[i].s != z {
			panic("unexpected search result")
		}
		i++
		m[g2] = cCDF(lm[j:i])
	}
	return m
}

// trigram returns the trigram with prefix g2 using a probability value
// in the range [0.0,1.0).
func (c comap) trigram(g2 string, p float64) string {
	prs := c[g2]
	i := prs.SearchProb(p)
	return prs[i].s
}

var (
	// CDF for normal probabilities
	pcdf = pCDFOfLM(englm3)
	// map of two letter conditionals
	cmap = comapOfLM(englm3)
)

// Reader generates a stream of text of uppercase letters with trigrams
// distributed according to a language model of the English language.
type Reader struct {
	rnd *rand.Rand
	g3  string
}

// NewReader creates a new reader. The argument src must create a uniformly
// distributed stream of random values.
func NewReader(src rand.Source) *Reader {
	rnd := rand.New(src)
	i := pcdf.SearchProb(rnd.Float64())
	return &Reader{rnd, pcdf[i].s}
}

// Read reads random text. The Read function will always return len(p)
// bytes and will never return an error.
func (r *Reader) Read(p []byte) (n int, err error) {
	for i := range p {
		r.g3 = cmap.trigram(r.g3[1:], r.rnd.Float64())
		p[i] = r.g3[2]
	}
	return len(p), nil
}
