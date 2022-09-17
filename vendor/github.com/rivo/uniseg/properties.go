package uniseg

// The Unicode properties as used in the various parsers. Only the ones needed
// in the context of this package are included.
const (
	prXX      = 0    // Same as prAny.
	prAny     = iota // prAny must be 0.
	prPrepend        // Grapheme properties must come first, to reduce the number of bits stored in the state vector.
	prCR
	prLF
	prControl
	prExtend
	prRegionalIndicator
	prSpacingMark
	prL
	prV
	prT
	prLV
	prLVT
	prZWJ
	prExtendedPictographic
	prNewline
	prWSegSpace
	prDoubleQuote
	prSingleQuote
	prMidNumLet
	prNumeric
	prMidLetter
	prMidNum
	prExtendNumLet
	prALetter
	prFormat
	prHebrewLetter
	prKatakana
	prSp
	prSTerm
	prClose
	prSContinue
	prATerm
	prUpper
	prLower
	prSep
	prOLetter
	prCM
	prBA
	prBK
	prSP
	prEX
	prQU
	prAL
	prPR
	prPO
	prOP
	prCP
	prIS
	prHY
	prSY
	prNU
	prCL
	prNL
	prGL
	prAI
	prBB
	prHL
	prSA
	prJL
	prJV
	prJT
	prNS
	prZW
	prB2
	prIN
	prWJ
	prID
	prEB
	prCJ
	prH2
	prH3
	prSG
	prCB
	prRI
	prEM
	prN
	prNa
	prA
	prW
	prH
	prF
	prEmojiPresentation
)

// Unicode General Categories. Only the ones needed in the context of this
// package are included.
const (
	gcNone = iota // gcNone must be 0.
	gcCc
	gcZs
	gcPo
	gcSc
	gcPs
	gcPe
	gcSm
	gcPd
	gcNd
	gcLu
	gcSk
	gcPc
	gcLl
	gcSo
	gcLo
	gcPi
	gcCf
	gcNo
	gcPf
	gcLC
	gcLm
	gcMn
	gcMe
	gcMc
	gcNl
	gcZl
	gcZp
	gcCn
	gcCs
	gcCo
)

// Special code points.
const (
	vs15 = 0xfe0e // Variation Selector-15 (text presentation)
	vs16 = 0xfe0f // Variation Selector-16 (emoji presentation)
)

// propertySearch performs a binary search on a property slice and returns the
// entry whose range (start = first array element, end = second array element)
// includes r, or an array of 0's if no such entry was found.
func propertySearch[E interface{ [3]int | [4]int }](dictionary []E, r rune) (result E) {
	// Run a binary search.
	from := 0
	to := len(dictionary)
	for to > from {
		middle := (from + to) / 2
		cpRange := dictionary[middle]
		if int(r) < cpRange[0] {
			to = middle
			continue
		}
		if int(r) > cpRange[1] {
			from = middle + 1
			continue
		}
		return cpRange
	}
	return
}

// property returns the Unicode property value (see constants above) of the
// given code point.
func property(dictionary [][3]int, r rune) int {
	return propertySearch(dictionary, r)[2]
}

// propertyWithGenCat returns the Unicode property value and General Category
// (see constants above) of the given code point.
func propertyWithGenCat(dictionary [][4]int, r rune) (property, generalCategory int) {
	entry := propertySearch(dictionary, r)
	return entry[2], entry[3]
}
