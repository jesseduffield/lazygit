// Code generated via go generate from gen_properties.go. DO NOT EDIT.

package uniseg

// graphemeCodePoints are taken from
// https://www.unicode.org/Public/15.0.0/ucd/auxiliary/GraphemeBreakProperty.txt
// and
// https://unicode.org/Public/15.0.0/ucd/emoji/emoji-data.txt
// ("Extended_Pictographic" only)
// on September 5, 2023. See https://www.unicode.org/license.html for the Unicode
// license agreement.
var graphemeCodePoints = [][3]int{
	{0x0000, 0x0009, prControl},                // Cc  [10] <control-0000>..<control-0009>
	{0x000A, 0x000A, prLF},                     // Cc       <control-000A>
	{0x000B, 0x000C, prControl},                // Cc   [2] <control-000B>..<control-000C>
	{0x000D, 0x000D, prCR},                     // Cc       <control-000D>
	{0x000E, 0x001F, prControl},                // Cc  [18] <control-000E>..<control-001F>
	{0x007F, 0x009F, prControl},                // Cc  [33] <control-007F>..<control-009F>
	{0x00A9, 0x00A9, prExtendedPictographic},   // E0.6   [1] (©️)       copyright
	{0x00AD, 0x00AD, prControl},                // Cf       SOFT HYPHEN
	{0x00AE, 0x00AE, prExtendedPictographic},   // E0.6   [1] (®️)       registered
	{0x0300, 0x036F, prExtend},                 // Mn [112] COMBINING GRAVE ACCENT..COMBINING LATIN SMALL LETTER X
	{0x0483, 0x0487, prExtend},                 // Mn   [5] COMBINING CYRILLIC TITLO..COMBINING CYRILLIC POKRYTIE
	{0x0488, 0x0489, prExtend},                 // Me   [2] COMBINING CYRILLIC HUNDRED THOUSANDS SIGN..COMBINING CYRILLIC MILLIONS SIGN
	{0x0591, 0x05BD, prExtend},                 // Mn  [45] HEBREW ACCENT ETNAHTA..HEBREW POINT METEG
	{0x05BF, 0x05BF, prExtend},                 // Mn       HEBREW POINT RAFE
	{0x05C1, 0x05C2, prExtend},                 // Mn   [2] HEBREW POINT SHIN DOT..HEBREW POINT SIN DOT
	{0x05C4, 0x05C5, prExtend},                 // Mn   [2] HEBREW MARK UPPER DOT..HEBREW MARK LOWER DOT
	{0x05C7, 0x05C7, prExtend},                 // Mn       HEBREW POINT QAMATS QATAN
	{0x0600, 0x0605, prPrepend},                // Cf   [6] ARABIC NUMBER SIGN..ARABIC NUMBER MARK ABOVE
	{0x0610, 0x061A, prExtend},                 // Mn  [11] ARABIC SIGN SALLALLAHOU ALAYHE WASSALLAM..ARABIC SMALL KASRA
	{0x061C, 0x061C, prControl},                // Cf       ARABIC LETTER MARK
	{0x064B, 0x065F, prExtend},                 // Mn  [21] ARABIC FATHATAN..ARABIC WAVY HAMZA BELOW
	{0x0670, 0x0670, prExtend},                 // Mn       ARABIC LETTER SUPERSCRIPT ALEF
	{0x06D6, 0x06DC, prExtend},                 // Mn   [7] ARABIC SMALL HIGH LIGATURE SAD WITH LAM WITH ALEF MAKSURA..ARABIC SMALL HIGH SEEN
	{0x06DD, 0x06DD, prPrepend},                // Cf       ARABIC END OF AYAH
	{0x06DF, 0x06E4, prExtend},                 // Mn   [6] ARABIC SMALL HIGH ROUNDED ZERO..ARABIC SMALL HIGH MADDA
	{0x06E7, 0x06E8, prExtend},                 // Mn   [2] ARABIC SMALL HIGH YEH..ARABIC SMALL HIGH NOON
	{0x06EA, 0x06ED, prExtend},                 // Mn   [4] ARABIC EMPTY CENTRE LOW STOP..ARABIC SMALL LOW MEEM
	{0x070F, 0x070F, prPrepend},                // Cf       SYRIAC ABBREVIATION MARK
	{0x0711, 0x0711, prExtend},                 // Mn       SYRIAC LETTER SUPERSCRIPT ALAPH
	{0x0730, 0x074A, prExtend},                 // Mn  [27] SYRIAC PTHAHA ABOVE..SYRIAC BARREKH
	{0x07A6, 0x07B0, prExtend},                 // Mn  [11] THAANA ABAFILI..THAANA SUKUN
	{0x07EB, 0x07F3, prExtend},                 // Mn   [9] NKO COMBINING SHORT HIGH TONE..NKO COMBINING DOUBLE DOT ABOVE
	{0x07FD, 0x07FD, prExtend},                 // Mn       NKO DANTAYALAN
	{0x0816, 0x0819, prExtend},                 // Mn   [4] SAMARITAN MARK IN..SAMARITAN MARK DAGESH
	{0x081B, 0x0823, prExtend},                 // Mn   [9] SAMARITAN MARK EPENTHETIC YUT..SAMARITAN VOWEL SIGN A
	{0x0825, 0x0827, prExtend},                 // Mn   [3] SAMARITAN VOWEL SIGN SHORT A..SAMARITAN VOWEL SIGN U
	{0x0829, 0x082D, prExtend},                 // Mn   [5] SAMARITAN VOWEL SIGN LONG I..SAMARITAN MARK NEQUDAA
	{0x0859, 0x085B, prExtend},                 // Mn   [3] MANDAIC AFFRICATION MARK..MANDAIC GEMINATION MARK
	{0x0890, 0x0891, prPrepend},                // Cf   [2] ARABIC POUND MARK ABOVE..ARABIC PIASTRE MARK ABOVE
	{0x0898, 0x089F, prExtend},                 // Mn   [8] ARABIC SMALL HIGH WORD AL-JUZ..ARABIC HALF MADDA OVER MADDA
	{0x08CA, 0x08E1, prExtend},                 // Mn  [24] ARABIC SMALL HIGH FARSI YEH..ARABIC SMALL HIGH SIGN SAFHA
	{0x08E2, 0x08E2, prPrepend},                // Cf       ARABIC DISPUTED END OF AYAH
	{0x08E3, 0x0902, prExtend},                 // Mn  [32] ARABIC TURNED DAMMA BELOW..DEVANAGARI SIGN ANUSVARA
	{0x0903, 0x0903, prSpacingMark},            // Mc       DEVANAGARI SIGN VISARGA
	{0x093A, 0x093A, prExtend},                 // Mn       DEVANAGARI VOWEL SIGN OE
	{0x093B, 0x093B, prSpacingMark},            // Mc       DEVANAGARI VOWEL SIGN OOE
	{0x093C, 0x093C, prExtend},                 // Mn       DEVANAGARI SIGN NUKTA
	{0x093E, 0x0940, prSpacingMark},            // Mc   [3] DEVANAGARI VOWEL SIGN AA..DEVANAGARI VOWEL SIGN II
	{0x0941, 0x0948, prExtend},                 // Mn   [8] DEVANAGARI VOWEL SIGN U..DEVANAGARI VOWEL SIGN AI
	{0x0949, 0x094C, prSpacingMark},            // Mc   [4] DEVANAGARI VOWEL SIGN CANDRA O..DEVANAGARI VOWEL SIGN AU
	{0x094D, 0x094D, prExtend},                 // Mn       DEVANAGARI SIGN VIRAMA
	{0x094E, 0x094F, prSpacingMark},            // Mc   [2] DEVANAGARI VOWEL SIGN PRISHTHAMATRA E..DEVANAGARI VOWEL SIGN AW
	{0x0951, 0x0957, prExtend},                 // Mn   [7] DEVANAGARI STRESS SIGN UDATTA..DEVANAGARI VOWEL SIGN UUE
	{0x0962, 0x0963, prExtend},                 // Mn   [2] DEVANAGARI VOWEL SIGN VOCALIC L..DEVANAGARI VOWEL SIGN VOCALIC LL
	{0x0981, 0x0981, prExtend},                 // Mn       BENGALI SIGN CANDRABINDU
	{0x0982, 0x0983, prSpacingMark},            // Mc   [2] BENGALI SIGN ANUSVARA..BENGALI SIGN VISARGA
	{0x09BC, 0x09BC, prExtend},                 // Mn       BENGALI SIGN NUKTA
	{0x09BE, 0x09BE, prExtend},                 // Mc       BENGALI VOWEL SIGN AA
	{0x09BF, 0x09C0, prSpacingMark},            // Mc   [2] BENGALI VOWEL SIGN I..BENGALI VOWEL SIGN II
	{0x09C1, 0x09C4, prExtend},                 // Mn   [4] BENGALI VOWEL SIGN U..BENGALI VOWEL SIGN VOCALIC RR
	{0x09C7, 0x09C8, prSpacingMark},            // Mc   [2] BENGALI VOWEL SIGN E..BENGALI VOWEL SIGN AI
	{0x09CB, 0x09CC, prSpacingMark},            // Mc   [2] BENGALI VOWEL SIGN O..BENGALI VOWEL SIGN AU
	{0x09CD, 0x09CD, prExtend},                 // Mn       BENGALI SIGN VIRAMA
	{0x09D7, 0x09D7, prExtend},                 // Mc       BENGALI AU LENGTH MARK
	{0x09E2, 0x09E3, prExtend},                 // Mn   [2] BENGALI VOWEL SIGN VOCALIC L..BENGALI VOWEL SIGN VOCALIC LL
	{0x09FE, 0x09FE, prExtend},                 // Mn       BENGALI SANDHI MARK
	{0x0A01, 0x0A02, prExtend},                 // Mn   [2] GURMUKHI SIGN ADAK BINDI..GURMUKHI SIGN BINDI
	{0x0A03, 0x0A03, prSpacingMark},            // Mc       GURMUKHI SIGN VISARGA
	{0x0A3C, 0x0A3C, prExtend},                 // Mn       GURMUKHI SIGN NUKTA
	{0x0A3E, 0x0A40, prSpacingMark},            // Mc   [3] GURMUKHI VOWEL SIGN AA..GURMUKHI VOWEL SIGN II
	{0x0A41, 0x0A42, prExtend},                 // Mn   [2] GURMUKHI VOWEL SIGN U..GURMUKHI VOWEL SIGN UU
	{0x0A47, 0x0A48, prExtend},                 // Mn   [2] GURMUKHI VOWEL SIGN EE..GURMUKHI VOWEL SIGN AI
	{0x0A4B, 0x0A4D, prExtend},                 // Mn   [3] GURMUKHI VOWEL SIGN OO..GURMUKHI SIGN VIRAMA
	{0x0A51, 0x0A51, prExtend},                 // Mn       GURMUKHI SIGN UDAAT
	{0x0A70, 0x0A71, prExtend},                 // Mn   [2] GURMUKHI TIPPI..GURMUKHI ADDAK
	{0x0A75, 0x0A75, prExtend},                 // Mn       GURMUKHI SIGN YAKASH
	{0x0A81, 0x0A82, prExtend},                 // Mn   [2] GUJARATI SIGN CANDRABINDU..GUJARATI SIGN ANUSVARA
	{0x0A83, 0x0A83, prSpacingMark},            // Mc       GUJARATI SIGN VISARGA
	{0x0ABC, 0x0ABC, prExtend},                 // Mn       GUJARATI SIGN NUKTA
	{0x0ABE, 0x0AC0, prSpacingMark},            // Mc   [3] GUJARATI VOWEL SIGN AA..GUJARATI VOWEL SIGN II
	{0x0AC1, 0x0AC5, prExtend},                 // Mn   [5] GUJARATI VOWEL SIGN U..GUJARATI VOWEL SIGN CANDRA E
	{0x0AC7, 0x0AC8, prExtend},                 // Mn   [2] GUJARATI VOWEL SIGN E..GUJARATI VOWEL SIGN AI
	{0x0AC9, 0x0AC9, prSpacingMark},            // Mc       GUJARATI VOWEL SIGN CANDRA O
	{0x0ACB, 0x0ACC, prSpacingMark},            // Mc   [2] GUJARATI VOWEL SIGN O..GUJARATI VOWEL SIGN AU
	{0x0ACD, 0x0ACD, prExtend},                 // Mn       GUJARATI SIGN VIRAMA
	{0x0AE2, 0x0AE3, prExtend},                 // Mn   [2] GUJARATI VOWEL SIGN VOCALIC L..GUJARATI VOWEL SIGN VOCALIC LL
	{0x0AFA, 0x0AFF, prExtend},                 // Mn   [6] GUJARATI SIGN SUKUN..GUJARATI SIGN TWO-CIRCLE NUKTA ABOVE
	{0x0B01, 0x0B01, prExtend},                 // Mn       ORIYA SIGN CANDRABINDU
	{0x0B02, 0x0B03, prSpacingMark},            // Mc   [2] ORIYA SIGN ANUSVARA..ORIYA SIGN VISARGA
	{0x0B3C, 0x0B3C, prExtend},                 // Mn       ORIYA SIGN NUKTA
	{0x0B3E, 0x0B3E, prExtend},                 // Mc       ORIYA VOWEL SIGN AA
	{0x0B3F, 0x0B3F, prExtend},                 // Mn       ORIYA VOWEL SIGN I
	{0x0B40, 0x0B40, prSpacingMark},            // Mc       ORIYA VOWEL SIGN II
	{0x0B41, 0x0B44, prExtend},                 // Mn   [4] ORIYA VOWEL SIGN U..ORIYA VOWEL SIGN VOCALIC RR
	{0x0B47, 0x0B48, prSpacingMark},            // Mc   [2] ORIYA VOWEL SIGN E..ORIYA VOWEL SIGN AI
	{0x0B4B, 0x0B4C, prSpacingMark},            // Mc   [2] ORIYA VOWEL SIGN O..ORIYA VOWEL SIGN AU
	{0x0B4D, 0x0B4D, prExtend},                 // Mn       ORIYA SIGN VIRAMA
	{0x0B55, 0x0B56, prExtend},                 // Mn   [2] ORIYA SIGN OVERLINE..ORIYA AI LENGTH MARK
	{0x0B57, 0x0B57, prExtend},                 // Mc       ORIYA AU LENGTH MARK
	{0x0B62, 0x0B63, prExtend},                 // Mn   [2] ORIYA VOWEL SIGN VOCALIC L..ORIYA VOWEL SIGN VOCALIC LL
	{0x0B82, 0x0B82, prExtend},                 // Mn       TAMIL SIGN ANUSVARA
	{0x0BBE, 0x0BBE, prExtend},                 // Mc       TAMIL VOWEL SIGN AA
	{0x0BBF, 0x0BBF, prSpacingMark},            // Mc       TAMIL VOWEL SIGN I
	{0x0BC0, 0x0BC0, prExtend},                 // Mn       TAMIL VOWEL SIGN II
	{0x0BC1, 0x0BC2, prSpacingMark},            // Mc   [2] TAMIL VOWEL SIGN U..TAMIL VOWEL SIGN UU
	{0x0BC6, 0x0BC8, prSpacingMark},            // Mc   [3] TAMIL VOWEL SIGN E..TAMIL VOWEL SIGN AI
	{0x0BCA, 0x0BCC, prSpacingMark},            // Mc   [3] TAMIL VOWEL SIGN O..TAMIL VOWEL SIGN AU
	{0x0BCD, 0x0BCD, prExtend},                 // Mn       TAMIL SIGN VIRAMA
	{0x0BD7, 0x0BD7, prExtend},                 // Mc       TAMIL AU LENGTH MARK
	{0x0C00, 0x0C00, prExtend},                 // Mn       TELUGU SIGN COMBINING CANDRABINDU ABOVE
	{0x0C01, 0x0C03, prSpacingMark},            // Mc   [3] TELUGU SIGN CANDRABINDU..TELUGU SIGN VISARGA
	{0x0C04, 0x0C04, prExtend},                 // Mn       TELUGU SIGN COMBINING ANUSVARA ABOVE
	{0x0C3C, 0x0C3C, prExtend},                 // Mn       TELUGU SIGN NUKTA
	{0x0C3E, 0x0C40, prExtend},                 // Mn   [3] TELUGU VOWEL SIGN AA..TELUGU VOWEL SIGN II
	{0x0C41, 0x0C44, prSpacingMark},            // Mc   [4] TELUGU VOWEL SIGN U..TELUGU VOWEL SIGN VOCALIC RR
	{0x0C46, 0x0C48, prExtend},                 // Mn   [3] TELUGU VOWEL SIGN E..TELUGU VOWEL SIGN AI
	{0x0C4A, 0x0C4D, prExtend},                 // Mn   [4] TELUGU VOWEL SIGN O..TELUGU SIGN VIRAMA
	{0x0C55, 0x0C56, prExtend},                 // Mn   [2] TELUGU LENGTH MARK..TELUGU AI LENGTH MARK
	{0x0C62, 0x0C63, prExtend},                 // Mn   [2] TELUGU VOWEL SIGN VOCALIC L..TELUGU VOWEL SIGN VOCALIC LL
	{0x0C81, 0x0C81, prExtend},                 // Mn       KANNADA SIGN CANDRABINDU
	{0x0C82, 0x0C83, prSpacingMark},            // Mc   [2] KANNADA SIGN ANUSVARA..KANNADA SIGN VISARGA
	{0x0CBC, 0x0CBC, prExtend},                 // Mn       KANNADA SIGN NUKTA
	{0x0CBE, 0x0CBE, prSpacingMark},            // Mc       KANNADA VOWEL SIGN AA
	{0x0CBF, 0x0CBF, prExtend},                 // Mn       KANNADA VOWEL SIGN I
	{0x0CC0, 0x0CC1, prSpacingMark},            // Mc   [2] KANNADA VOWEL SIGN II..KANNADA VOWEL SIGN U
	{0x0CC2, 0x0CC2, prExtend},                 // Mc       KANNADA VOWEL SIGN UU
	{0x0CC3, 0x0CC4, prSpacingMark},            // Mc   [2] KANNADA VOWEL SIGN VOCALIC R..KANNADA VOWEL SIGN VOCALIC RR
	{0x0CC6, 0x0CC6, prExtend},                 // Mn       KANNADA VOWEL SIGN E
	{0x0CC7, 0x0CC8, prSpacingMark},            // Mc   [2] KANNADA VOWEL SIGN EE..KANNADA VOWEL SIGN AI
	{0x0CCA, 0x0CCB, prSpacingMark},            // Mc   [2] KANNADA VOWEL SIGN O..KANNADA VOWEL SIGN OO
	{0x0CCC, 0x0CCD, prExtend},                 // Mn   [2] KANNADA VOWEL SIGN AU..KANNADA SIGN VIRAMA
	{0x0CD5, 0x0CD6, prExtend},                 // Mc   [2] KANNADA LENGTH MARK..KANNADA AI LENGTH MARK
	{0x0CE2, 0x0CE3, prExtend},                 // Mn   [2] KANNADA VOWEL SIGN VOCALIC L..KANNADA VOWEL SIGN VOCALIC LL
	{0x0CF3, 0x0CF3, prSpacingMark},            // Mc       KANNADA SIGN COMBINING ANUSVARA ABOVE RIGHT
	{0x0D00, 0x0D01, prExtend},                 // Mn   [2] MALAYALAM SIGN COMBINING ANUSVARA ABOVE..MALAYALAM SIGN CANDRABINDU
	{0x0D02, 0x0D03, prSpacingMark},            // Mc   [2] MALAYALAM SIGN ANUSVARA..MALAYALAM SIGN VISARGA
	{0x0D3B, 0x0D3C, prExtend},                 // Mn   [2] MALAYALAM SIGN VERTICAL BAR VIRAMA..MALAYALAM SIGN CIRCULAR VIRAMA
	{0x0D3E, 0x0D3E, prExtend},                 // Mc       MALAYALAM VOWEL SIGN AA
	{0x0D3F, 0x0D40, prSpacingMark},            // Mc   [2] MALAYALAM VOWEL SIGN I..MALAYALAM VOWEL SIGN II
	{0x0D41, 0x0D44, prExtend},                 // Mn   [4] MALAYALAM VOWEL SIGN U..MALAYALAM VOWEL SIGN VOCALIC RR
	{0x0D46, 0x0D48, prSpacingMark},            // Mc   [3] MALAYALAM VOWEL SIGN E..MALAYALAM VOWEL SIGN AI
	{0x0D4A, 0x0D4C, prSpacingMark},            // Mc   [3] MALAYALAM VOWEL SIGN O..MALAYALAM VOWEL SIGN AU
	{0x0D4D, 0x0D4D, prExtend},                 // Mn       MALAYALAM SIGN VIRAMA
	{0x0D4E, 0x0D4E, prPrepend},                // Lo       MALAYALAM LETTER DOT REPH
	{0x0D57, 0x0D57, prExtend},                 // Mc       MALAYALAM AU LENGTH MARK
	{0x0D62, 0x0D63, prExtend},                 // Mn   [2] MALAYALAM VOWEL SIGN VOCALIC L..MALAYALAM VOWEL SIGN VOCALIC LL
	{0x0D81, 0x0D81, prExtend},                 // Mn       SINHALA SIGN CANDRABINDU
	{0x0D82, 0x0D83, prSpacingMark},            // Mc   [2] SINHALA SIGN ANUSVARAYA..SINHALA SIGN VISARGAYA
	{0x0DCA, 0x0DCA, prExtend},                 // Mn       SINHALA SIGN AL-LAKUNA
	{0x0DCF, 0x0DCF, prExtend},                 // Mc       SINHALA VOWEL SIGN AELA-PILLA
	{0x0DD0, 0x0DD1, prSpacingMark},            // Mc   [2] SINHALA VOWEL SIGN KETTI AEDA-PILLA..SINHALA VOWEL SIGN DIGA AEDA-PILLA
	{0x0DD2, 0x0DD4, prExtend},                 // Mn   [3] SINHALA VOWEL SIGN KETTI IS-PILLA..SINHALA VOWEL SIGN KETTI PAA-PILLA
	{0x0DD6, 0x0DD6, prExtend},                 // Mn       SINHALA VOWEL SIGN DIGA PAA-PILLA
	{0x0DD8, 0x0DDE, prSpacingMark},            // Mc   [7] SINHALA VOWEL SIGN GAETTA-PILLA..SINHALA VOWEL SIGN KOMBUVA HAA GAYANUKITTA
	{0x0DDF, 0x0DDF, prExtend},                 // Mc       SINHALA VOWEL SIGN GAYANUKITTA
	{0x0DF2, 0x0DF3, prSpacingMark},            // Mc   [2] SINHALA VOWEL SIGN DIGA GAETTA-PILLA..SINHALA VOWEL SIGN DIGA GAYANUKITTA
	{0x0E31, 0x0E31, prExtend},                 // Mn       THAI CHARACTER MAI HAN-AKAT
	{0x0E33, 0x0E33, prSpacingMark},            // Lo       THAI CHARACTER SARA AM
	{0x0E34, 0x0E3A, prExtend},                 // Mn   [7] THAI CHARACTER SARA I..THAI CHARACTER PHINTHU
	{0x0E47, 0x0E4E, prExtend},                 // Mn   [8] THAI CHARACTER MAITAIKHU..THAI CHARACTER YAMAKKAN
	{0x0EB1, 0x0EB1, prExtend},                 // Mn       LAO VOWEL SIGN MAI KAN
	{0x0EB3, 0x0EB3, prSpacingMark},            // Lo       LAO VOWEL SIGN AM
	{0x0EB4, 0x0EBC, prExtend},                 // Mn   [9] LAO VOWEL SIGN I..LAO SEMIVOWEL SIGN LO
	{0x0EC8, 0x0ECE, prExtend},                 // Mn   [7] LAO TONE MAI EK..LAO YAMAKKAN
	{0x0F18, 0x0F19, prExtend},                 // Mn   [2] TIBETAN ASTROLOGICAL SIGN -KHYUD PA..TIBETAN ASTROLOGICAL SIGN SDONG TSHUGS
	{0x0F35, 0x0F35, prExtend},                 // Mn       TIBETAN MARK NGAS BZUNG NYI ZLA
	{0x0F37, 0x0F37, prExtend},                 // Mn       TIBETAN MARK NGAS BZUNG SGOR RTAGS
	{0x0F39, 0x0F39, prExtend},                 // Mn       TIBETAN MARK TSA -PHRU
	{0x0F3E, 0x0F3F, prSpacingMark},            // Mc   [2] TIBETAN SIGN YAR TSHES..TIBETAN SIGN MAR TSHES
	{0x0F71, 0x0F7E, prExtend},                 // Mn  [14] TIBETAN VOWEL SIGN AA..TIBETAN SIGN RJES SU NGA RO
	{0x0F7F, 0x0F7F, prSpacingMark},            // Mc       TIBETAN SIGN RNAM BCAD
	{0x0F80, 0x0F84, prExtend},                 // Mn   [5] TIBETAN VOWEL SIGN REVERSED I..TIBETAN MARK HALANTA
	{0x0F86, 0x0F87, prExtend},                 // Mn   [2] TIBETAN SIGN LCI RTAGS..TIBETAN SIGN YANG RTAGS
	{0x0F8D, 0x0F97, prExtend},                 // Mn  [11] TIBETAN SUBJOINED SIGN LCE TSA CAN..TIBETAN SUBJOINED LETTER JA
	{0x0F99, 0x0FBC, prExtend},                 // Mn  [36] TIBETAN SUBJOINED LETTER NYA..TIBETAN SUBJOINED LETTER FIXED-FORM RA
	{0x0FC6, 0x0FC6, prExtend},                 // Mn       TIBETAN SYMBOL PADMA GDAN
	{0x102D, 0x1030, prExtend},                 // Mn   [4] MYANMAR VOWEL SIGN I..MYANMAR VOWEL SIGN UU
	{0x1031, 0x1031, prSpacingMark},            // Mc       MYANMAR VOWEL SIGN E
	{0x1032, 0x1037, prExtend},                 // Mn   [6] MYANMAR VOWEL SIGN AI..MYANMAR SIGN DOT BELOW
	{0x1039, 0x103A, prExtend},                 // Mn   [2] MYANMAR SIGN VIRAMA..MYANMAR SIGN ASAT
	{0x103B, 0x103C, prSpacingMark},            // Mc   [2] MYANMAR CONSONANT SIGN MEDIAL YA..MYANMAR CONSONANT SIGN MEDIAL RA
	{0x103D, 0x103E, prExtend},                 // Mn   [2] MYANMAR CONSONANT SIGN MEDIAL WA..MYANMAR CONSONANT SIGN MEDIAL HA
	{0x1056, 0x1057, prSpacingMark},            // Mc   [2] MYANMAR VOWEL SIGN VOCALIC R..MYANMAR VOWEL SIGN VOCALIC RR
	{0x1058, 0x1059, prExtend},                 // Mn   [2] MYANMAR VOWEL SIGN VOCALIC L..MYANMAR VOWEL SIGN VOCALIC LL
	{0x105E, 0x1060, prExtend},                 // Mn   [3] MYANMAR CONSONANT SIGN MON MEDIAL NA..MYANMAR CONSONANT SIGN MON MEDIAL LA
	{0x1071, 0x1074, prExtend},                 // Mn   [4] MYANMAR VOWEL SIGN GEBA KAREN I..MYANMAR VOWEL SIGN KAYAH EE
	{0x1082, 0x1082, prExtend},                 // Mn       MYANMAR CONSONANT SIGN SHAN MEDIAL WA
	{0x1084, 0x1084, prSpacingMark},            // Mc       MYANMAR VOWEL SIGN SHAN E
	{0x1085, 0x1086, prExtend},                 // Mn   [2] MYANMAR VOWEL SIGN SHAN E ABOVE..MYANMAR VOWEL SIGN SHAN FINAL Y
	{0x108D, 0x108D, prExtend},                 // Mn       MYANMAR SIGN SHAN COUNCIL EMPHATIC TONE
	{0x109D, 0x109D, prExtend},                 // Mn       MYANMAR VOWEL SIGN AITON AI
	{0x1100, 0x115F, prL},                      // Lo  [96] HANGUL CHOSEONG KIYEOK..HANGUL CHOSEONG FILLER
	{0x1160, 0x11A7, prV},                      // Lo  [72] HANGUL JUNGSEONG FILLER..HANGUL JUNGSEONG O-YAE
	{0x11A8, 0x11FF, prT},                      // Lo  [88] HANGUL JONGSEONG KIYEOK..HANGUL JONGSEONG SSANGNIEUN
	{0x135D, 0x135F, prExtend},                 // Mn   [3] ETHIOPIC COMBINING GEMINATION AND VOWEL LENGTH MARK..ETHIOPIC COMBINING GEMINATION MARK
	{0x1712, 0x1714, prExtend},                 // Mn   [3] TAGALOG VOWEL SIGN I..TAGALOG SIGN VIRAMA
	{0x1715, 0x1715, prSpacingMark},            // Mc       TAGALOG SIGN PAMUDPOD
	{0x1732, 0x1733, prExtend},                 // Mn   [2] HANUNOO VOWEL SIGN I..HANUNOO VOWEL SIGN U
	{0x1734, 0x1734, prSpacingMark},            // Mc       HANUNOO SIGN PAMUDPOD
	{0x1752, 0x1753, prExtend},                 // Mn   [2] BUHID VOWEL SIGN I..BUHID VOWEL SIGN U
	{0x1772, 0x1773, prExtend},                 // Mn   [2] TAGBANWA VOWEL SIGN I..TAGBANWA VOWEL SIGN U
	{0x17B4, 0x17B5, prExtend},                 // Mn   [2] KHMER VOWEL INHERENT AQ..KHMER VOWEL INHERENT AA
	{0x17B6, 0x17B6, prSpacingMark},            // Mc       KHMER VOWEL SIGN AA
	{0x17B7, 0x17BD, prExtend},                 // Mn   [7] KHMER VOWEL SIGN I..KHMER VOWEL SIGN UA
	{0x17BE, 0x17C5, prSpacingMark},            // Mc   [8] KHMER VOWEL SIGN OE..KHMER VOWEL SIGN AU
	{0x17C6, 0x17C6, prExtend},                 // Mn       KHMER SIGN NIKAHIT
	{0x17C7, 0x17C8, prSpacingMark},            // Mc   [2] KHMER SIGN REAHMUK..KHMER SIGN YUUKALEAPINTU
	{0x17C9, 0x17D3, prExtend},                 // Mn  [11] KHMER SIGN MUUSIKATOAN..KHMER SIGN BATHAMASAT
	{0x17DD, 0x17DD, prExtend},                 // Mn       KHMER SIGN ATTHACAN
	{0x180B, 0x180D, prExtend},                 // Mn   [3] MONGOLIAN FREE VARIATION SELECTOR ONE..MONGOLIAN FREE VARIATION SELECTOR THREE
	{0x180E, 0x180E, prControl},                // Cf       MONGOLIAN VOWEL SEPARATOR
	{0x180F, 0x180F, prExtend},                 // Mn       MONGOLIAN FREE VARIATION SELECTOR FOUR
	{0x1885, 0x1886, prExtend},                 // Mn   [2] MONGOLIAN LETTER ALI GALI BALUDA..MONGOLIAN LETTER ALI GALI THREE BALUDA
	{0x18A9, 0x18A9, prExtend},                 // Mn       MONGOLIAN LETTER ALI GALI DAGALGA
	{0x1920, 0x1922, prExtend},                 // Mn   [3] LIMBU VOWEL SIGN A..LIMBU VOWEL SIGN U
	{0x1923, 0x1926, prSpacingMark},            // Mc   [4] LIMBU VOWEL SIGN EE..LIMBU VOWEL SIGN AU
	{0x1927, 0x1928, prExtend},                 // Mn   [2] LIMBU VOWEL SIGN E..LIMBU VOWEL SIGN O
	{0x1929, 0x192B, prSpacingMark},            // Mc   [3] LIMBU SUBJOINED LETTER YA..LIMBU SUBJOINED LETTER WA
	{0x1930, 0x1931, prSpacingMark},            // Mc   [2] LIMBU SMALL LETTER KA..LIMBU SMALL LETTER NGA
	{0x1932, 0x1932, prExtend},                 // Mn       LIMBU SMALL LETTER ANUSVARA
	{0x1933, 0x1938, prSpacingMark},            // Mc   [6] LIMBU SMALL LETTER TA..LIMBU SMALL LETTER LA
	{0x1939, 0x193B, prExtend},                 // Mn   [3] LIMBU SIGN MUKPHRENG..LIMBU SIGN SA-I
	{0x1A17, 0x1A18, prExtend},                 // Mn   [2] BUGINESE VOWEL SIGN I..BUGINESE VOWEL SIGN U
	{0x1A19, 0x1A1A, prSpacingMark},            // Mc   [2] BUGINESE VOWEL SIGN E..BUGINESE VOWEL SIGN O
	{0x1A1B, 0x1A1B, prExtend},                 // Mn       BUGINESE VOWEL SIGN AE
	{0x1A55, 0x1A55, prSpacingMark},            // Mc       TAI THAM CONSONANT SIGN MEDIAL RA
	{0x1A56, 0x1A56, prExtend},                 // Mn       TAI THAM CONSONANT SIGN MEDIAL LA
	{0x1A57, 0x1A57, prSpacingMark},            // Mc       TAI THAM CONSONANT SIGN LA TANG LAI
	{0x1A58, 0x1A5E, prExtend},                 // Mn   [7] TAI THAM SIGN MAI KANG LAI..TAI THAM CONSONANT SIGN SA
	{0x1A60, 0x1A60, prExtend},                 // Mn       TAI THAM SIGN SAKOT
	{0x1A62, 0x1A62, prExtend},                 // Mn       TAI THAM VOWEL SIGN MAI SAT
	{0x1A65, 0x1A6C, prExtend},                 // Mn   [8] TAI THAM VOWEL SIGN I..TAI THAM VOWEL SIGN OA BELOW
	{0x1A6D, 0x1A72, prSpacingMark},            // Mc   [6] TAI THAM VOWEL SIGN OY..TAI THAM VOWEL SIGN THAM AI
	{0x1A73, 0x1A7C, prExtend},                 // Mn  [10] TAI THAM VOWEL SIGN OA ABOVE..TAI THAM SIGN KHUEN-LUE KARAN
	{0x1A7F, 0x1A7F, prExtend},                 // Mn       TAI THAM COMBINING CRYPTOGRAMMIC DOT
	{0x1AB0, 0x1ABD, prExtend},                 // Mn  [14] COMBINING DOUBLED CIRCUMFLEX ACCENT..COMBINING PARENTHESES BELOW
	{0x1ABE, 0x1ABE, prExtend},                 // Me       COMBINING PARENTHESES OVERLAY
	{0x1ABF, 0x1ACE, prExtend},                 // Mn  [16] COMBINING LATIN SMALL LETTER W BELOW..COMBINING LATIN SMALL LETTER INSULAR T
	{0x1B00, 0x1B03, prExtend},                 // Mn   [4] BALINESE SIGN ULU RICEM..BALINESE SIGN SURANG
	{0x1B04, 0x1B04, prSpacingMark},            // Mc       BALINESE SIGN BISAH
	{0x1B34, 0x1B34, prExtend},                 // Mn       BALINESE SIGN REREKAN
	{0x1B35, 0x1B35, prExtend},                 // Mc       BALINESE VOWEL SIGN TEDUNG
	{0x1B36, 0x1B3A, prExtend},                 // Mn   [5] BALINESE VOWEL SIGN ULU..BALINESE VOWEL SIGN RA REPA
	{0x1B3B, 0x1B3B, prSpacingMark},            // Mc       BALINESE VOWEL SIGN RA REPA TEDUNG
	{0x1B3C, 0x1B3C, prExtend},                 // Mn       BALINESE VOWEL SIGN LA LENGA
	{0x1B3D, 0x1B41, prSpacingMark},            // Mc   [5] BALINESE VOWEL SIGN LA LENGA TEDUNG..BALINESE VOWEL SIGN TALING REPA TEDUNG
	{0x1B42, 0x1B42, prExtend},                 // Mn       BALINESE VOWEL SIGN PEPET
	{0x1B43, 0x1B44, prSpacingMark},            // Mc   [2] BALINESE VOWEL SIGN PEPET TEDUNG..BALINESE ADEG ADEG
	{0x1B6B, 0x1B73, prExtend},                 // Mn   [9] BALINESE MUSICAL SYMBOL COMBINING TEGEH..BALINESE MUSICAL SYMBOL COMBINING GONG
	{0x1B80, 0x1B81, prExtend},                 // Mn   [2] SUNDANESE SIGN PANYECEK..SUNDANESE SIGN PANGLAYAR
	{0x1B82, 0x1B82, prSpacingMark},            // Mc       SUNDANESE SIGN PANGWISAD
	{0x1BA1, 0x1BA1, prSpacingMark},            // Mc       SUNDANESE CONSONANT SIGN PAMINGKAL
	{0x1BA2, 0x1BA5, prExtend},                 // Mn   [4] SUNDANESE CONSONANT SIGN PANYAKRA..SUNDANESE VOWEL SIGN PANYUKU
	{0x1BA6, 0x1BA7, prSpacingMark},            // Mc   [2] SUNDANESE VOWEL SIGN PANAELAENG..SUNDANESE VOWEL SIGN PANOLONG
	{0x1BA8, 0x1BA9, prExtend},                 // Mn   [2] SUNDANESE VOWEL SIGN PAMEPET..SUNDANESE VOWEL SIGN PANEULEUNG
	{0x1BAA, 0x1BAA, prSpacingMark},            // Mc       SUNDANESE SIGN PAMAAEH
	{0x1BAB, 0x1BAD, prExtend},                 // Mn   [3] SUNDANESE SIGN VIRAMA..SUNDANESE CONSONANT SIGN PASANGAN WA
	{0x1BE6, 0x1BE6, prExtend},                 // Mn       BATAK SIGN TOMPI
	{0x1BE7, 0x1BE7, prSpacingMark},            // Mc       BATAK VOWEL SIGN E
	{0x1BE8, 0x1BE9, prExtend},                 // Mn   [2] BATAK VOWEL SIGN PAKPAK E..BATAK VOWEL SIGN EE
	{0x1BEA, 0x1BEC, prSpacingMark},            // Mc   [3] BATAK VOWEL SIGN I..BATAK VOWEL SIGN O
	{0x1BED, 0x1BED, prExtend},                 // Mn       BATAK VOWEL SIGN KARO O
	{0x1BEE, 0x1BEE, prSpacingMark},            // Mc       BATAK VOWEL SIGN U
	{0x1BEF, 0x1BF1, prExtend},                 // Mn   [3] BATAK VOWEL SIGN U FOR SIMALUNGUN SA..BATAK CONSONANT SIGN H
	{0x1BF2, 0x1BF3, prSpacingMark},            // Mc   [2] BATAK PANGOLAT..BATAK PANONGONAN
	{0x1C24, 0x1C2B, prSpacingMark},            // Mc   [8] LEPCHA SUBJOINED LETTER YA..LEPCHA VOWEL SIGN UU
	{0x1C2C, 0x1C33, prExtend},                 // Mn   [8] LEPCHA VOWEL SIGN E..LEPCHA CONSONANT SIGN T
	{0x1C34, 0x1C35, prSpacingMark},            // Mc   [2] LEPCHA CONSONANT SIGN NYIN-DO..LEPCHA CONSONANT SIGN KANG
	{0x1C36, 0x1C37, prExtend},                 // Mn   [2] LEPCHA SIGN RAN..LEPCHA SIGN NUKTA
	{0x1CD0, 0x1CD2, prExtend},                 // Mn   [3] VEDIC TONE KARSHANA..VEDIC TONE PRENKHA
	{0x1CD4, 0x1CE0, prExtend},                 // Mn  [13] VEDIC SIGN YAJURVEDIC MIDLINE SVARITA..VEDIC TONE RIGVEDIC KASHMIRI INDEPENDENT SVARITA
	{0x1CE1, 0x1CE1, prSpacingMark},            // Mc       VEDIC TONE ATHARVAVEDIC INDEPENDENT SVARITA
	{0x1CE2, 0x1CE8, prExtend},                 // Mn   [7] VEDIC SIGN VISARGA SVARITA..VEDIC SIGN VISARGA ANUDATTA WITH TAIL
	{0x1CED, 0x1CED, prExtend},                 // Mn       VEDIC SIGN TIRYAK
	{0x1CF4, 0x1CF4, prExtend},                 // Mn       VEDIC TONE CANDRA ABOVE
	{0x1CF7, 0x1CF7, prSpacingMark},            // Mc       VEDIC SIGN ATIKRAMA
	{0x1CF8, 0x1CF9, prExtend},                 // Mn   [2] VEDIC TONE RING ABOVE..VEDIC TONE DOUBLE RING ABOVE
	{0x1DC0, 0x1DFF, prExtend},                 // Mn  [64] COMBINING DOTTED GRAVE ACCENT..COMBINING RIGHT ARROWHEAD AND DOWN ARROWHEAD BELOW
	{0x200B, 0x200B, prControl},                // Cf       ZERO WIDTH SPACE
	{0x200C, 0x200C, prExtend},                 // Cf       ZERO WIDTH NON-JOINER
	{0x200D, 0x200D, prZWJ},                    // Cf       ZERO WIDTH JOINER
	{0x200E, 0x200F, prControl},                // Cf   [2] LEFT-TO-RIGHT MARK..RIGHT-TO-LEFT MARK
	{0x2028, 0x2028, prControl},                // Zl       LINE SEPARATOR
	{0x2029, 0x2029, prControl},                // Zp       PARAGRAPH SEPARATOR
	{0x202A, 0x202E, prControl},                // Cf   [5] LEFT-TO-RIGHT EMBEDDING..RIGHT-TO-LEFT OVERRIDE
	{0x203C, 0x203C, prExtendedPictographic},   // E0.6   [1] (‼️)       double exclamation mark
	{0x2049, 0x2049, prExtendedPictographic},   // E0.6   [1] (⁉️)       exclamation question mark
	{0x2060, 0x2064, prControl},                // Cf   [5] WORD JOINER..INVISIBLE PLUS
	{0x2065, 0x2065, prControl},                // Cn       <reserved-2065>
	{0x2066, 0x206F, prControl},                // Cf  [10] LEFT-TO-RIGHT ISOLATE..NOMINAL DIGIT SHAPES
	{0x20D0, 0x20DC, prExtend},                 // Mn  [13] COMBINING LEFT HARPOON ABOVE..COMBINING FOUR DOTS ABOVE
	{0x20DD, 0x20E0, prExtend},                 // Me   [4] COMBINING ENCLOSING CIRCLE..COMBINING ENCLOSING CIRCLE BACKSLASH
	{0x20E1, 0x20E1, prExtend},                 // Mn       COMBINING LEFT RIGHT ARROW ABOVE
	{0x20E2, 0x20E4, prExtend},                 // Me   [3] COMBINING ENCLOSING SCREEN..COMBINING ENCLOSING UPWARD POINTING TRIANGLE
	{0x20E5, 0x20F0, prExtend},                 // Mn  [12] COMBINING REVERSE SOLIDUS OVERLAY..COMBINING ASTERISK ABOVE
	{0x2122, 0x2122, prExtendedPictographic},   // E0.6   [1] (™️)       trade mark
	{0x2139, 0x2139, prExtendedPictographic},   // E0.6   [1] (ℹ️)       information
	{0x2194, 0x2199, prExtendedPictographic},   // E0.6   [6] (↔️..↙️)    left-right arrow..down-left arrow
	{0x21A9, 0x21AA, prExtendedPictographic},   // E0.6   [2] (↩️..↪️)    right arrow curving left..left arrow curving right
	{0x231A, 0x231B, prExtendedPictographic},   // E0.6   [2] (⌚..⌛)    watch..hourglass done
	{0x2328, 0x2328, prExtendedPictographic},   // E1.0   [1] (⌨️)       keyboard
	{0x2388, 0x2388, prExtendedPictographic},   // E0.0   [1] (⎈)       HELM SYMBOL
	{0x23CF, 0x23CF, prExtendedPictographic},   // E1.0   [1] (⏏️)       eject button
	{0x23E9, 0x23EC, prExtendedPictographic},   // E0.6   [4] (⏩..⏬)    fast-forward button..fast down button
	{0x23ED, 0x23EE, prExtendedPictographic},   // E0.7   [2] (⏭️..⏮️)    next track button..last track button
	{0x23EF, 0x23EF, prExtendedPictographic},   // E1.0   [1] (⏯️)       play or pause button
	{0x23F0, 0x23F0, prExtendedPictographic},   // E0.6   [1] (⏰)       alarm clock
	{0x23F1, 0x23F2, prExtendedPictographic},   // E1.0   [2] (⏱️..⏲️)    stopwatch..timer clock
	{0x23F3, 0x23F3, prExtendedPictographic},   // E0.6   [1] (⏳)       hourglass not done
	{0x23F8, 0x23FA, prExtendedPictographic},   // E0.7   [3] (⏸️..⏺️)    pause button..record button
	{0x24C2, 0x24C2, prExtendedPictographic},   // E0.6   [1] (Ⓜ️)       circled M
	{0x25AA, 0x25AB, prExtendedPictographic},   // E0.6   [2] (▪️..▫️)    black small square..white small square
	{0x25B6, 0x25B6, prExtendedPictographic},   // E0.6   [1] (▶️)       play button
	{0x25C0, 0x25C0, prExtendedPictographic},   // E0.6   [1] (◀️)       reverse button
	{0x25FB, 0x25FE, prExtendedPictographic},   // E0.6   [4] (◻️..◾)    white medium square..black medium-small square
	{0x2600, 0x2601, prExtendedPictographic},   // E0.6   [2] (☀️..☁️)    sun..cloud
	{0x2602, 0x2603, prExtendedPictographic},   // E0.7   [2] (☂️..☃️)    umbrella..snowman
	{0x2604, 0x2604, prExtendedPictographic},   // E1.0   [1] (☄️)       comet
	{0x2605, 0x2605, prExtendedPictographic},   // E0.0   [1] (★)       BLACK STAR
	{0x2607, 0x260D, prExtendedPictographic},   // E0.0   [7] (☇..☍)    LIGHTNING..OPPOSITION
	{0x260E, 0x260E, prExtendedPictographic},   // E0.6   [1] (☎️)       telephone
	{0x260F, 0x2610, prExtendedPictographic},   // E0.0   [2] (☏..☐)    WHITE TELEPHONE..BALLOT BOX
	{0x2611, 0x2611, prExtendedPictographic},   // E0.6   [1] (☑️)       check box with check
	{0x2612, 0x2612, prExtendedPictographic},   // E0.0   [1] (☒)       BALLOT BOX WITH X
	{0x2614, 0x2615, prExtendedPictographic},   // E0.6   [2] (☔..☕)    umbrella with rain drops..hot beverage
	{0x2616, 0x2617, prExtendedPictographic},   // E0.0   [2] (☖..☗)    WHITE SHOGI PIECE..BLACK SHOGI PIECE
	{0x2618, 0x2618, prExtendedPictographic},   // E1.0   [1] (☘️)       shamrock
	{0x2619, 0x261C, prExtendedPictographic},   // E0.0   [4] (☙..☜)    REVERSED ROTATED FLORAL HEART BULLET..WHITE LEFT POINTING INDEX
	{0x261D, 0x261D, prExtendedPictographic},   // E0.6   [1] (☝️)       index pointing up
	{0x261E, 0x261F, prExtendedPictographic},   // E0.0   [2] (☞..☟)    WHITE RIGHT POINTING INDEX..WHITE DOWN POINTING INDEX
	{0x2620, 0x2620, prExtendedPictographic},   // E1.0   [1] (☠️)       skull and crossbones
	{0x2621, 0x2621, prExtendedPictographic},   // E0.0   [1] (☡)       CAUTION SIGN
	{0x2622, 0x2623, prExtendedPictographic},   // E1.0   [2] (☢️..☣️)    radioactive..biohazard
	{0x2624, 0x2625, prExtendedPictographic},   // E0.0   [2] (☤..☥)    CADUCEUS..ANKH
	{0x2626, 0x2626, prExtendedPictographic},   // E1.0   [1] (☦️)       orthodox cross
	{0x2627, 0x2629, prExtendedPictographic},   // E0.0   [3] (☧..☩)    CHI RHO..CROSS OF JERUSALEM
	{0x262A, 0x262A, prExtendedPictographic},   // E0.7   [1] (☪️)       star and crescent
	{0x262B, 0x262D, prExtendedPictographic},   // E0.0   [3] (☫..☭)    FARSI SYMBOL..HAMMER AND SICKLE
	{0x262E, 0x262E, prExtendedPictographic},   // E1.0   [1] (☮️)       peace symbol
	{0x262F, 0x262F, prExtendedPictographic},   // E0.7   [1] (☯️)       yin yang
	{0x2630, 0x2637, prExtendedPictographic},   // E0.0   [8] (☰..☷)    TRIGRAM FOR HEAVEN..TRIGRAM FOR EARTH
	{0x2638, 0x2639, prExtendedPictographic},   // E0.7   [2] (☸️..☹️)    wheel of dharma..frowning face
	{0x263A, 0x263A, prExtendedPictographic},   // E0.6   [1] (☺️)       smiling face
	{0x263B, 0x263F, prExtendedPictographic},   // E0.0   [5] (☻..☿)    BLACK SMILING FACE..MERCURY
	{0x2640, 0x2640, prExtendedPictographic},   // E4.0   [1] (♀️)       female sign
	{0x2641, 0x2641, prExtendedPictographic},   // E0.0   [1] (♁)       EARTH
	{0x2642, 0x2642, prExtendedPictographic},   // E4.0   [1] (♂️)       male sign
	{0x2643, 0x2647, prExtendedPictographic},   // E0.0   [5] (♃..♇)    JUPITER..PLUTO
	{0x2648, 0x2653, prExtendedPictographic},   // E0.6  [12] (♈..♓)    Aries..Pisces
	{0x2654, 0x265E, prExtendedPictographic},   // E0.0  [11] (♔..♞)    WHITE CHESS KING..BLACK CHESS KNIGHT
	{0x265F, 0x265F, prExtendedPictographic},   // E11.0  [1] (♟️)       chess pawn
	{0x2660, 0x2660, prExtendedPictographic},   // E0.6   [1] (♠️)       spade suit
	{0x2661, 0x2662, prExtendedPictographic},   // E0.0   [2] (♡..♢)    WHITE HEART SUIT..WHITE DIAMOND SUIT
	{0x2663, 0x2663, prExtendedPictographic},   // E0.6   [1] (♣️)       club suit
	{0x2664, 0x2664, prExtendedPictographic},   // E0.0   [1] (♤)       WHITE SPADE SUIT
	{0x2665, 0x2666, prExtendedPictographic},   // E0.6   [2] (♥️..♦️)    heart suit..diamond suit
	{0x2667, 0x2667, prExtendedPictographic},   // E0.0   [1] (♧)       WHITE CLUB SUIT
	{0x2668, 0x2668, prExtendedPictographic},   // E0.6   [1] (♨️)       hot springs
	{0x2669, 0x267A, prExtendedPictographic},   // E0.0  [18] (♩..♺)    QUARTER NOTE..RECYCLING SYMBOL FOR GENERIC MATERIALS
	{0x267B, 0x267B, prExtendedPictographic},   // E0.6   [1] (♻️)       recycling symbol
	{0x267C, 0x267D, prExtendedPictographic},   // E0.0   [2] (♼..♽)    RECYCLED PAPER SYMBOL..PARTIALLY-RECYCLED PAPER SYMBOL
	{0x267E, 0x267E, prExtendedPictographic},   // E11.0  [1] (♾️)       infinity
	{0x267F, 0x267F, prExtendedPictographic},   // E0.6   [1] (♿)       wheelchair symbol
	{0x2680, 0x2685, prExtendedPictographic},   // E0.0   [6] (⚀..⚅)    DIE FACE-1..DIE FACE-6
	{0x2690, 0x2691, prExtendedPictographic},   // E0.0   [2] (⚐..⚑)    WHITE FLAG..BLACK FLAG
	{0x2692, 0x2692, prExtendedPictographic},   // E1.0   [1] (⚒️)       hammer and pick
	{0x2693, 0x2693, prExtendedPictographic},   // E0.6   [1] (⚓)       anchor
	{0x2694, 0x2694, prExtendedPictographic},   // E1.0   [1] (⚔️)       crossed swords
	{0x2695, 0x2695, prExtendedPictographic},   // E4.0   [1] (⚕️)       medical symbol
	{0x2696, 0x2697, prExtendedPictographic},   // E1.0   [2] (⚖️..⚗️)    balance scale..alembic
	{0x2698, 0x2698, prExtendedPictographic},   // E0.0   [1] (⚘)       FLOWER
	{0x2699, 0x2699, prExtendedPictographic},   // E1.0   [1] (⚙️)       gear
	{0x269A, 0x269A, prExtendedPictographic},   // E0.0   [1] (⚚)       STAFF OF HERMES
	{0x269B, 0x269C, prExtendedPictographic},   // E1.0   [2] (⚛️..⚜️)    atom symbol..fleur-de-lis
	{0x269D, 0x269F, prExtendedPictographic},   // E0.0   [3] (⚝..⚟)    OUTLINED WHITE STAR..THREE LINES CONVERGING LEFT
	{0x26A0, 0x26A1, prExtendedPictographic},   // E0.6   [2] (⚠️..⚡)    warning..high voltage
	{0x26A2, 0x26A6, prExtendedPictographic},   // E0.0   [5] (⚢..⚦)    DOUBLED FEMALE SIGN..MALE WITH STROKE SIGN
	{0x26A7, 0x26A7, prExtendedPictographic},   // E13.0  [1] (⚧️)       transgender symbol
	{0x26A8, 0x26A9, prExtendedPictographic},   // E0.0   [2] (⚨..⚩)    VERTICAL MALE WITH STROKE SIGN..HORIZONTAL MALE WITH STROKE SIGN
	{0x26AA, 0x26AB, prExtendedPictographic},   // E0.6   [2] (⚪..⚫)    white circle..black circle
	{0x26AC, 0x26AF, prExtendedPictographic},   // E0.0   [4] (⚬..⚯)    MEDIUM SMALL WHITE CIRCLE..UNMARRIED PARTNERSHIP SYMBOL
	{0x26B0, 0x26B1, prExtendedPictographic},   // E1.0   [2] (⚰️..⚱️)    coffin..funeral urn
	{0x26B2, 0x26BC, prExtendedPictographic},   // E0.0  [11] (⚲..⚼)    NEUTER..SESQUIQUADRATE
	{0x26BD, 0x26BE, prExtendedPictographic},   // E0.6   [2] (⚽..⚾)    soccer ball..baseball
	{0x26BF, 0x26C3, prExtendedPictographic},   // E0.0   [5] (⚿..⛃)    SQUARED KEY..BLACK DRAUGHTS KING
	{0x26C4, 0x26C5, prExtendedPictographic},   // E0.6   [2] (⛄..⛅)    snowman without snow..sun behind cloud
	{0x26C6, 0x26C7, prExtendedPictographic},   // E0.0   [2] (⛆..⛇)    RAIN..BLACK SNOWMAN
	{0x26C8, 0x26C8, prExtendedPictographic},   // E0.7   [1] (⛈️)       cloud with lightning and rain
	{0x26C9, 0x26CD, prExtendedPictographic},   // E0.0   [5] (⛉..⛍)    TURNED WHITE SHOGI PIECE..DISABLED CAR
	{0x26CE, 0x26CE, prExtendedPictographic},   // E0.6   [1] (⛎)       Ophiuchus
	{0x26CF, 0x26CF, prExtendedPictographic},   // E0.7   [1] (⛏️)       pick
	{0x26D0, 0x26D0, prExtendedPictographic},   // E0.0   [1] (⛐)       CAR SLIDING
	{0x26D1, 0x26D1, prExtendedPictographic},   // E0.7   [1] (⛑️)       rescue worker’s helmet
	{0x26D2, 0x26D2, prExtendedPictographic},   // E0.0   [1] (⛒)       CIRCLED CROSSING LANES
	{0x26D3, 0x26D3, prExtendedPictographic},   // E0.7   [1] (⛓️)       chains
	{0x26D4, 0x26D4, prExtendedPictographic},   // E0.6   [1] (⛔)       no entry
	{0x26D5, 0x26E8, prExtendedPictographic},   // E0.0  [20] (⛕..⛨)    ALTERNATE ONE-WAY LEFT WAY TRAFFIC..BLACK CROSS ON SHIELD
	{0x26E9, 0x26E9, prExtendedPictographic},   // E0.7   [1] (⛩️)       shinto shrine
	{0x26EA, 0x26EA, prExtendedPictographic},   // E0.6   [1] (⛪)       church
	{0x26EB, 0x26EF, prExtendedPictographic},   // E0.0   [5] (⛫..⛯)    CASTLE..MAP SYMBOL FOR LIGHTHOUSE
	{0x26F0, 0x26F1, prExtendedPictographic},   // E0.7   [2] (⛰️..⛱️)    mountain..umbrella on ground
	{0x26F2, 0x26F3, prExtendedPictographic},   // E0.6   [2] (⛲..⛳)    fountain..flag in hole
	{0x26F4, 0x26F4, prExtendedPictographic},   // E0.7   [1] (⛴️)       ferry
	{0x26F5, 0x26F5, prExtendedPictographic},   // E0.6   [1] (⛵)       sailboat
	{0x26F6, 0x26F6, prExtendedPictographic},   // E0.0   [1] (⛶)       SQUARE FOUR CORNERS
	{0x26F7, 0x26F9, prExtendedPictographic},   // E0.7   [3] (⛷️..⛹️)    skier..person bouncing ball
	{0x26FA, 0x26FA, prExtendedPictographic},   // E0.6   [1] (⛺)       tent
	{0x26FB, 0x26FC, prExtendedPictographic},   // E0.0   [2] (⛻..⛼)    JAPANESE BANK SYMBOL..HEADSTONE GRAVEYARD SYMBOL
	{0x26FD, 0x26FD, prExtendedPictographic},   // E0.6   [1] (⛽)       fuel pump
	{0x26FE, 0x2701, prExtendedPictographic},   // E0.0   [4] (⛾..✁)    CUP ON BLACK SQUARE..UPPER BLADE SCISSORS
	{0x2702, 0x2702, prExtendedPictographic},   // E0.6   [1] (✂️)       scissors
	{0x2703, 0x2704, prExtendedPictographic},   // E0.0   [2] (✃..✄)    LOWER BLADE SCISSORS..WHITE SCISSORS
	{0x2705, 0x2705, prExtendedPictographic},   // E0.6   [1] (✅)       check mark button
	{0x2708, 0x270C, prExtendedPictographic},   // E0.6   [5] (✈️..✌️)    airplane..victory hand
	{0x270D, 0x270D, prExtendedPictographic},   // E0.7   [1] (✍️)       writing hand
	{0x270E, 0x270E, prExtendedPictographic},   // E0.0   [1] (✎)       LOWER RIGHT PENCIL
	{0x270F, 0x270F, prExtendedPictographic},   // E0.6   [1] (✏️)       pencil
	{0x2710, 0x2711, prExtendedPictographic},   // E0.0   [2] (✐..✑)    UPPER RIGHT PENCIL..WHITE NIB
	{0x2712, 0x2712, prExtendedPictographic},   // E0.6   [1] (✒️)       black nib
	{0x2714, 0x2714, prExtendedPictographic},   // E0.6   [1] (✔️)       check mark
	{0x2716, 0x2716, prExtendedPictographic},   // E0.6   [1] (✖️)       multiply
	{0x271D, 0x271D, prExtendedPictographic},   // E0.7   [1] (✝️)       latin cross
	{0x2721, 0x2721, prExtendedPictographic},   // E0.7   [1] (✡️)       star of David
	{0x2728, 0x2728, prExtendedPictographic},   // E0.6   [1] (✨)       sparkles
	{0x2733, 0x2734, prExtendedPictographic},   // E0.6   [2] (✳️..✴️)    eight-spoked asterisk..eight-pointed star
	{0x2744, 0x2744, prExtendedPictographic},   // E0.6   [1] (❄️)       snowflake
	{0x2747, 0x2747, prExtendedPictographic},   // E0.6   [1] (❇️)       sparkle
	{0x274C, 0x274C, prExtendedPictographic},   // E0.6   [1] (❌)       cross mark
	{0x274E, 0x274E, prExtendedPictographic},   // E0.6   [1] (❎)       cross mark button
	{0x2753, 0x2755, prExtendedPictographic},   // E0.6   [3] (❓..❕)    red question mark..white exclamation mark
	{0x2757, 0x2757, prExtendedPictographic},   // E0.6   [1] (❗)       red exclamation mark
	{0x2763, 0x2763, prExtendedPictographic},   // E1.0   [1] (❣️)       heart exclamation
	{0x2764, 0x2764, prExtendedPictographic},   // E0.6   [1] (❤️)       red heart
	{0x2765, 0x2767, prExtendedPictographic},   // E0.0   [3] (❥..❧)    ROTATED HEAVY BLACK HEART BULLET..ROTATED FLORAL HEART BULLET
	{0x2795, 0x2797, prExtendedPictographic},   // E0.6   [3] (➕..➗)    plus..divide
	{0x27A1, 0x27A1, prExtendedPictographic},   // E0.6   [1] (➡️)       right arrow
	{0x27B0, 0x27B0, prExtendedPictographic},   // E0.6   [1] (➰)       curly loop
	{0x27BF, 0x27BF, prExtendedPictographic},   // E1.0   [1] (➿)       double curly loop
	{0x2934, 0x2935, prExtendedPictographic},   // E0.6   [2] (⤴️..⤵️)    right arrow curving up..right arrow curving down
	{0x2B05, 0x2B07, prExtendedPictographic},   // E0.6   [3] (⬅️..⬇️)    left arrow..down arrow
	{0x2B1B, 0x2B1C, prExtendedPictographic},   // E0.6   [2] (⬛..⬜)    black large square..white large square
	{0x2B50, 0x2B50, prExtendedPictographic},   // E0.6   [1] (⭐)       star
	{0x2B55, 0x2B55, prExtendedPictographic},   // E0.6   [1] (⭕)       hollow red circle
	{0x2CEF, 0x2CF1, prExtend},                 // Mn   [3] COPTIC COMBINING NI ABOVE..COPTIC COMBINING SPIRITUS LENIS
	{0x2D7F, 0x2D7F, prExtend},                 // Mn       TIFINAGH CONSONANT JOINER
	{0x2DE0, 0x2DFF, prExtend},                 // Mn  [32] COMBINING CYRILLIC LETTER BE..COMBINING CYRILLIC LETTER IOTIFIED BIG YUS
	{0x302A, 0x302D, prExtend},                 // Mn   [4] IDEOGRAPHIC LEVEL TONE MARK..IDEOGRAPHIC ENTERING TONE MARK
	{0x302E, 0x302F, prExtend},                 // Mc   [2] HANGUL SINGLE DOT TONE MARK..HANGUL DOUBLE DOT TONE MARK
	{0x3030, 0x3030, prExtendedPictographic},   // E0.6   [1] (〰️)       wavy dash
	{0x303D, 0x303D, prExtendedPictographic},   // E0.6   [1] (〽️)       part alternation mark
	{0x3099, 0x309A, prExtend},                 // Mn   [2] COMBINING KATAKANA-HIRAGANA VOICED SOUND MARK..COMBINING KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	{0x3297, 0x3297, prExtendedPictographic},   // E0.6   [1] (㊗️)       Japanese “congratulations” button
	{0x3299, 0x3299, prExtendedPictographic},   // E0.6   [1] (㊙️)       Japanese “secret” button
	{0xA66F, 0xA66F, prExtend},                 // Mn       COMBINING CYRILLIC VZMET
	{0xA670, 0xA672, prExtend},                 // Me   [3] COMBINING CYRILLIC TEN MILLIONS SIGN..COMBINING CYRILLIC THOUSAND MILLIONS SIGN
	{0xA674, 0xA67D, prExtend},                 // Mn  [10] COMBINING CYRILLIC LETTER UKRAINIAN IE..COMBINING CYRILLIC PAYEROK
	{0xA69E, 0xA69F, prExtend},                 // Mn   [2] COMBINING CYRILLIC LETTER EF..COMBINING CYRILLIC LETTER IOTIFIED E
	{0xA6F0, 0xA6F1, prExtend},                 // Mn   [2] BAMUM COMBINING MARK KOQNDON..BAMUM COMBINING MARK TUKWENTIS
	{0xA802, 0xA802, prExtend},                 // Mn       SYLOTI NAGRI SIGN DVISVARA
	{0xA806, 0xA806, prExtend},                 // Mn       SYLOTI NAGRI SIGN HASANTA
	{0xA80B, 0xA80B, prExtend},                 // Mn       SYLOTI NAGRI SIGN ANUSVARA
	{0xA823, 0xA824, prSpacingMark},            // Mc   [2] SYLOTI NAGRI VOWEL SIGN A..SYLOTI NAGRI VOWEL SIGN I
	{0xA825, 0xA826, prExtend},                 // Mn   [2] SYLOTI NAGRI VOWEL SIGN U..SYLOTI NAGRI VOWEL SIGN E
	{0xA827, 0xA827, prSpacingMark},            // Mc       SYLOTI NAGRI VOWEL SIGN OO
	{0xA82C, 0xA82C, prExtend},                 // Mn       SYLOTI NAGRI SIGN ALTERNATE HASANTA
	{0xA880, 0xA881, prSpacingMark},            // Mc   [2] SAURASHTRA SIGN ANUSVARA..SAURASHTRA SIGN VISARGA
	{0xA8B4, 0xA8C3, prSpacingMark},            // Mc  [16] SAURASHTRA CONSONANT SIGN HAARU..SAURASHTRA VOWEL SIGN AU
	{0xA8C4, 0xA8C5, prExtend},                 // Mn   [2] SAURASHTRA SIGN VIRAMA..SAURASHTRA SIGN CANDRABINDU
	{0xA8E0, 0xA8F1, prExtend},                 // Mn  [18] COMBINING DEVANAGARI DIGIT ZERO..COMBINING DEVANAGARI SIGN AVAGRAHA
	{0xA8FF, 0xA8FF, prExtend},                 // Mn       DEVANAGARI VOWEL SIGN AY
	{0xA926, 0xA92D, prExtend},                 // Mn   [8] KAYAH LI VOWEL UE..KAYAH LI TONE CALYA PLOPHU
	{0xA947, 0xA951, prExtend},                 // Mn  [11] REJANG VOWEL SIGN I..REJANG CONSONANT SIGN R
	{0xA952, 0xA953, prSpacingMark},            // Mc   [2] REJANG CONSONANT SIGN H..REJANG VIRAMA
	{0xA960, 0xA97C, prL},                      // Lo  [29] HANGUL CHOSEONG TIKEUT-MIEUM..HANGUL CHOSEONG SSANGYEORINHIEUH
	{0xA980, 0xA982, prExtend},                 // Mn   [3] JAVANESE SIGN PANYANGGA..JAVANESE SIGN LAYAR
	{0xA983, 0xA983, prSpacingMark},            // Mc       JAVANESE SIGN WIGNYAN
	{0xA9B3, 0xA9B3, prExtend},                 // Mn       JAVANESE SIGN CECAK TELU
	{0xA9B4, 0xA9B5, prSpacingMark},            // Mc   [2] JAVANESE VOWEL SIGN TARUNG..JAVANESE VOWEL SIGN TOLONG
	{0xA9B6, 0xA9B9, prExtend},                 // Mn   [4] JAVANESE VOWEL SIGN WULU..JAVANESE VOWEL SIGN SUKU MENDUT
	{0xA9BA, 0xA9BB, prSpacingMark},            // Mc   [2] JAVANESE VOWEL SIGN TALING..JAVANESE VOWEL SIGN DIRGA MURE
	{0xA9BC, 0xA9BD, prExtend},                 // Mn   [2] JAVANESE VOWEL SIGN PEPET..JAVANESE CONSONANT SIGN KERET
	{0xA9BE, 0xA9C0, prSpacingMark},            // Mc   [3] JAVANESE CONSONANT SIGN PENGKAL..JAVANESE PANGKON
	{0xA9E5, 0xA9E5, prExtend},                 // Mn       MYANMAR SIGN SHAN SAW
	{0xAA29, 0xAA2E, prExtend},                 // Mn   [6] CHAM VOWEL SIGN AA..CHAM VOWEL SIGN OE
	{0xAA2F, 0xAA30, prSpacingMark},            // Mc   [2] CHAM VOWEL SIGN O..CHAM VOWEL SIGN AI
	{0xAA31, 0xAA32, prExtend},                 // Mn   [2] CHAM VOWEL SIGN AU..CHAM VOWEL SIGN UE
	{0xAA33, 0xAA34, prSpacingMark},            // Mc   [2] CHAM CONSONANT SIGN YA..CHAM CONSONANT SIGN RA
	{0xAA35, 0xAA36, prExtend},                 // Mn   [2] CHAM CONSONANT SIGN LA..CHAM CONSONANT SIGN WA
	{0xAA43, 0xAA43, prExtend},                 // Mn       CHAM CONSONANT SIGN FINAL NG
	{0xAA4C, 0xAA4C, prExtend},                 // Mn       CHAM CONSONANT SIGN FINAL M
	{0xAA4D, 0xAA4D, prSpacingMark},            // Mc       CHAM CONSONANT SIGN FINAL H
	{0xAA7C, 0xAA7C, prExtend},                 // Mn       MYANMAR SIGN TAI LAING TONE-2
	{0xAAB0, 0xAAB0, prExtend},                 // Mn       TAI VIET MAI KANG
	{0xAAB2, 0xAAB4, prExtend},                 // Mn   [3] TAI VIET VOWEL I..TAI VIET VOWEL U
	{0xAAB7, 0xAAB8, prExtend},                 // Mn   [2] TAI VIET MAI KHIT..TAI VIET VOWEL IA
	{0xAABE, 0xAABF, prExtend},                 // Mn   [2] TAI VIET VOWEL AM..TAI VIET TONE MAI EK
	{0xAAC1, 0xAAC1, prExtend},                 // Mn       TAI VIET TONE MAI THO
	{0xAAEB, 0xAAEB, prSpacingMark},            // Mc       MEETEI MAYEK VOWEL SIGN II
	{0xAAEC, 0xAAED, prExtend},                 // Mn   [2] MEETEI MAYEK VOWEL SIGN UU..MEETEI MAYEK VOWEL SIGN AAI
	{0xAAEE, 0xAAEF, prSpacingMark},            // Mc   [2] MEETEI MAYEK VOWEL SIGN AU..MEETEI MAYEK VOWEL SIGN AAU
	{0xAAF5, 0xAAF5, prSpacingMark},            // Mc       MEETEI MAYEK VOWEL SIGN VISARGA
	{0xAAF6, 0xAAF6, prExtend},                 // Mn       MEETEI MAYEK VIRAMA
	{0xABE3, 0xABE4, prSpacingMark},            // Mc   [2] MEETEI MAYEK VOWEL SIGN ONAP..MEETEI MAYEK VOWEL SIGN INAP
	{0xABE5, 0xABE5, prExtend},                 // Mn       MEETEI MAYEK VOWEL SIGN ANAP
	{0xABE6, 0xABE7, prSpacingMark},            // Mc   [2] MEETEI MAYEK VOWEL SIGN YENAP..MEETEI MAYEK VOWEL SIGN SOUNAP
	{0xABE8, 0xABE8, prExtend},                 // Mn       MEETEI MAYEK VOWEL SIGN UNAP
	{0xABE9, 0xABEA, prSpacingMark},            // Mc   [2] MEETEI MAYEK VOWEL SIGN CHEINAP..MEETEI MAYEK VOWEL SIGN NUNG
	{0xABEC, 0xABEC, prSpacingMark},            // Mc       MEETEI MAYEK LUM IYEK
	{0xABED, 0xABED, prExtend},                 // Mn       MEETEI MAYEK APUN IYEK
	{0xAC00, 0xAC00, prLV},                     // Lo       HANGUL SYLLABLE GA
	{0xAC01, 0xAC1B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GAG..HANGUL SYLLABLE GAH
	{0xAC1C, 0xAC1C, prLV},                     // Lo       HANGUL SYLLABLE GAE
	{0xAC1D, 0xAC37, prLVT},                    // Lo  [27] HANGUL SYLLABLE GAEG..HANGUL SYLLABLE GAEH
	{0xAC38, 0xAC38, prLV},                     // Lo       HANGUL SYLLABLE GYA
	{0xAC39, 0xAC53, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYAG..HANGUL SYLLABLE GYAH
	{0xAC54, 0xAC54, prLV},                     // Lo       HANGUL SYLLABLE GYAE
	{0xAC55, 0xAC6F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYAEG..HANGUL SYLLABLE GYAEH
	{0xAC70, 0xAC70, prLV},                     // Lo       HANGUL SYLLABLE GEO
	{0xAC71, 0xAC8B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GEOG..HANGUL SYLLABLE GEOH
	{0xAC8C, 0xAC8C, prLV},                     // Lo       HANGUL SYLLABLE GE
	{0xAC8D, 0xACA7, prLVT},                    // Lo  [27] HANGUL SYLLABLE GEG..HANGUL SYLLABLE GEH
	{0xACA8, 0xACA8, prLV},                     // Lo       HANGUL SYLLABLE GYEO
	{0xACA9, 0xACC3, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYEOG..HANGUL SYLLABLE GYEOH
	{0xACC4, 0xACC4, prLV},                     // Lo       HANGUL SYLLABLE GYE
	{0xACC5, 0xACDF, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYEG..HANGUL SYLLABLE GYEH
	{0xACE0, 0xACE0, prLV},                     // Lo       HANGUL SYLLABLE GO
	{0xACE1, 0xACFB, prLVT},                    // Lo  [27] HANGUL SYLLABLE GOG..HANGUL SYLLABLE GOH
	{0xACFC, 0xACFC, prLV},                     // Lo       HANGUL SYLLABLE GWA
	{0xACFD, 0xAD17, prLVT},                    // Lo  [27] HANGUL SYLLABLE GWAG..HANGUL SYLLABLE GWAH
	{0xAD18, 0xAD18, prLV},                     // Lo       HANGUL SYLLABLE GWAE
	{0xAD19, 0xAD33, prLVT},                    // Lo  [27] HANGUL SYLLABLE GWAEG..HANGUL SYLLABLE GWAEH
	{0xAD34, 0xAD34, prLV},                     // Lo       HANGUL SYLLABLE GOE
	{0xAD35, 0xAD4F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GOEG..HANGUL SYLLABLE GOEH
	{0xAD50, 0xAD50, prLV},                     // Lo       HANGUL SYLLABLE GYO
	{0xAD51, 0xAD6B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYOG..HANGUL SYLLABLE GYOH
	{0xAD6C, 0xAD6C, prLV},                     // Lo       HANGUL SYLLABLE GU
	{0xAD6D, 0xAD87, prLVT},                    // Lo  [27] HANGUL SYLLABLE GUG..HANGUL SYLLABLE GUH
	{0xAD88, 0xAD88, prLV},                     // Lo       HANGUL SYLLABLE GWEO
	{0xAD89, 0xADA3, prLVT},                    // Lo  [27] HANGUL SYLLABLE GWEOG..HANGUL SYLLABLE GWEOH
	{0xADA4, 0xADA4, prLV},                     // Lo       HANGUL SYLLABLE GWE
	{0xADA5, 0xADBF, prLVT},                    // Lo  [27] HANGUL SYLLABLE GWEG..HANGUL SYLLABLE GWEH
	{0xADC0, 0xADC0, prLV},                     // Lo       HANGUL SYLLABLE GWI
	{0xADC1, 0xADDB, prLVT},                    // Lo  [27] HANGUL SYLLABLE GWIG..HANGUL SYLLABLE GWIH
	{0xADDC, 0xADDC, prLV},                     // Lo       HANGUL SYLLABLE GYU
	{0xADDD, 0xADF7, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYUG..HANGUL SYLLABLE GYUH
	{0xADF8, 0xADF8, prLV},                     // Lo       HANGUL SYLLABLE GEU
	{0xADF9, 0xAE13, prLVT},                    // Lo  [27] HANGUL SYLLABLE GEUG..HANGUL SYLLABLE GEUH
	{0xAE14, 0xAE14, prLV},                     // Lo       HANGUL SYLLABLE GYI
	{0xAE15, 0xAE2F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GYIG..HANGUL SYLLABLE GYIH
	{0xAE30, 0xAE30, prLV},                     // Lo       HANGUL SYLLABLE GI
	{0xAE31, 0xAE4B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GIG..HANGUL SYLLABLE GIH
	{0xAE4C, 0xAE4C, prLV},                     // Lo       HANGUL SYLLABLE GGA
	{0xAE4D, 0xAE67, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGAG..HANGUL SYLLABLE GGAH
	{0xAE68, 0xAE68, prLV},                     // Lo       HANGUL SYLLABLE GGAE
	{0xAE69, 0xAE83, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGAEG..HANGUL SYLLABLE GGAEH
	{0xAE84, 0xAE84, prLV},                     // Lo       HANGUL SYLLABLE GGYA
	{0xAE85, 0xAE9F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYAG..HANGUL SYLLABLE GGYAH
	{0xAEA0, 0xAEA0, prLV},                     // Lo       HANGUL SYLLABLE GGYAE
	{0xAEA1, 0xAEBB, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYAEG..HANGUL SYLLABLE GGYAEH
	{0xAEBC, 0xAEBC, prLV},                     // Lo       HANGUL SYLLABLE GGEO
	{0xAEBD, 0xAED7, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGEOG..HANGUL SYLLABLE GGEOH
	{0xAED8, 0xAED8, prLV},                     // Lo       HANGUL SYLLABLE GGE
	{0xAED9, 0xAEF3, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGEG..HANGUL SYLLABLE GGEH
	{0xAEF4, 0xAEF4, prLV},                     // Lo       HANGUL SYLLABLE GGYEO
	{0xAEF5, 0xAF0F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYEOG..HANGUL SYLLABLE GGYEOH
	{0xAF10, 0xAF10, prLV},                     // Lo       HANGUL SYLLABLE GGYE
	{0xAF11, 0xAF2B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYEG..HANGUL SYLLABLE GGYEH
	{0xAF2C, 0xAF2C, prLV},                     // Lo       HANGUL SYLLABLE GGO
	{0xAF2D, 0xAF47, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGOG..HANGUL SYLLABLE GGOH
	{0xAF48, 0xAF48, prLV},                     // Lo       HANGUL SYLLABLE GGWA
	{0xAF49, 0xAF63, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGWAG..HANGUL SYLLABLE GGWAH
	{0xAF64, 0xAF64, prLV},                     // Lo       HANGUL SYLLABLE GGWAE
	{0xAF65, 0xAF7F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGWAEG..HANGUL SYLLABLE GGWAEH
	{0xAF80, 0xAF80, prLV},                     // Lo       HANGUL SYLLABLE GGOE
	{0xAF81, 0xAF9B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGOEG..HANGUL SYLLABLE GGOEH
	{0xAF9C, 0xAF9C, prLV},                     // Lo       HANGUL SYLLABLE GGYO
	{0xAF9D, 0xAFB7, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYOG..HANGUL SYLLABLE GGYOH
	{0xAFB8, 0xAFB8, prLV},                     // Lo       HANGUL SYLLABLE GGU
	{0xAFB9, 0xAFD3, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGUG..HANGUL SYLLABLE GGUH
	{0xAFD4, 0xAFD4, prLV},                     // Lo       HANGUL SYLLABLE GGWEO
	{0xAFD5, 0xAFEF, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGWEOG..HANGUL SYLLABLE GGWEOH
	{0xAFF0, 0xAFF0, prLV},                     // Lo       HANGUL SYLLABLE GGWE
	{0xAFF1, 0xB00B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGWEG..HANGUL SYLLABLE GGWEH
	{0xB00C, 0xB00C, prLV},                     // Lo       HANGUL SYLLABLE GGWI
	{0xB00D, 0xB027, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGWIG..HANGUL SYLLABLE GGWIH
	{0xB028, 0xB028, prLV},                     // Lo       HANGUL SYLLABLE GGYU
	{0xB029, 0xB043, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYUG..HANGUL SYLLABLE GGYUH
	{0xB044, 0xB044, prLV},                     // Lo       HANGUL SYLLABLE GGEU
	{0xB045, 0xB05F, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGEUG..HANGUL SYLLABLE GGEUH
	{0xB060, 0xB060, prLV},                     // Lo       HANGUL SYLLABLE GGYI
	{0xB061, 0xB07B, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGYIG..HANGUL SYLLABLE GGYIH
	{0xB07C, 0xB07C, prLV},                     // Lo       HANGUL SYLLABLE GGI
	{0xB07D, 0xB097, prLVT},                    // Lo  [27] HANGUL SYLLABLE GGIG..HANGUL SYLLABLE GGIH
	{0xB098, 0xB098, prLV},                     // Lo       HANGUL SYLLABLE NA
	{0xB099, 0xB0B3, prLVT},                    // Lo  [27] HANGUL SYLLABLE NAG..HANGUL SYLLABLE NAH
	{0xB0B4, 0xB0B4, prLV},                     // Lo       HANGUL SYLLABLE NAE
	{0xB0B5, 0xB0CF, prLVT},                    // Lo  [27] HANGUL SYLLABLE NAEG..HANGUL SYLLABLE NAEH
	{0xB0D0, 0xB0D0, prLV},                     // Lo       HANGUL SYLLABLE NYA
	{0xB0D1, 0xB0EB, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYAG..HANGUL SYLLABLE NYAH
	{0xB0EC, 0xB0EC, prLV},                     // Lo       HANGUL SYLLABLE NYAE
	{0xB0ED, 0xB107, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYAEG..HANGUL SYLLABLE NYAEH
	{0xB108, 0xB108, prLV},                     // Lo       HANGUL SYLLABLE NEO
	{0xB109, 0xB123, prLVT},                    // Lo  [27] HANGUL SYLLABLE NEOG..HANGUL SYLLABLE NEOH
	{0xB124, 0xB124, prLV},                     // Lo       HANGUL SYLLABLE NE
	{0xB125, 0xB13F, prLVT},                    // Lo  [27] HANGUL SYLLABLE NEG..HANGUL SYLLABLE NEH
	{0xB140, 0xB140, prLV},                     // Lo       HANGUL SYLLABLE NYEO
	{0xB141, 0xB15B, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYEOG..HANGUL SYLLABLE NYEOH
	{0xB15C, 0xB15C, prLV},                     // Lo       HANGUL SYLLABLE NYE
	{0xB15D, 0xB177, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYEG..HANGUL SYLLABLE NYEH
	{0xB178, 0xB178, prLV},                     // Lo       HANGUL SYLLABLE NO
	{0xB179, 0xB193, prLVT},                    // Lo  [27] HANGUL SYLLABLE NOG..HANGUL SYLLABLE NOH
	{0xB194, 0xB194, prLV},                     // Lo       HANGUL SYLLABLE NWA
	{0xB195, 0xB1AF, prLVT},                    // Lo  [27] HANGUL SYLLABLE NWAG..HANGUL SYLLABLE NWAH
	{0xB1B0, 0xB1B0, prLV},                     // Lo       HANGUL SYLLABLE NWAE
	{0xB1B1, 0xB1CB, prLVT},                    // Lo  [27] HANGUL SYLLABLE NWAEG..HANGUL SYLLABLE NWAEH
	{0xB1CC, 0xB1CC, prLV},                     // Lo       HANGUL SYLLABLE NOE
	{0xB1CD, 0xB1E7, prLVT},                    // Lo  [27] HANGUL SYLLABLE NOEG..HANGUL SYLLABLE NOEH
	{0xB1E8, 0xB1E8, prLV},                     // Lo       HANGUL SYLLABLE NYO
	{0xB1E9, 0xB203, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYOG..HANGUL SYLLABLE NYOH
	{0xB204, 0xB204, prLV},                     // Lo       HANGUL SYLLABLE NU
	{0xB205, 0xB21F, prLVT},                    // Lo  [27] HANGUL SYLLABLE NUG..HANGUL SYLLABLE NUH
	{0xB220, 0xB220, prLV},                     // Lo       HANGUL SYLLABLE NWEO
	{0xB221, 0xB23B, prLVT},                    // Lo  [27] HANGUL SYLLABLE NWEOG..HANGUL SYLLABLE NWEOH
	{0xB23C, 0xB23C, prLV},                     // Lo       HANGUL SYLLABLE NWE
	{0xB23D, 0xB257, prLVT},                    // Lo  [27] HANGUL SYLLABLE NWEG..HANGUL SYLLABLE NWEH
	{0xB258, 0xB258, prLV},                     // Lo       HANGUL SYLLABLE NWI
	{0xB259, 0xB273, prLVT},                    // Lo  [27] HANGUL SYLLABLE NWIG..HANGUL SYLLABLE NWIH
	{0xB274, 0xB274, prLV},                     // Lo       HANGUL SYLLABLE NYU
	{0xB275, 0xB28F, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYUG..HANGUL SYLLABLE NYUH
	{0xB290, 0xB290, prLV},                     // Lo       HANGUL SYLLABLE NEU
	{0xB291, 0xB2AB, prLVT},                    // Lo  [27] HANGUL SYLLABLE NEUG..HANGUL SYLLABLE NEUH
	{0xB2AC, 0xB2AC, prLV},                     // Lo       HANGUL SYLLABLE NYI
	{0xB2AD, 0xB2C7, prLVT},                    // Lo  [27] HANGUL SYLLABLE NYIG..HANGUL SYLLABLE NYIH
	{0xB2C8, 0xB2C8, prLV},                     // Lo       HANGUL SYLLABLE NI
	{0xB2C9, 0xB2E3, prLVT},                    // Lo  [27] HANGUL SYLLABLE NIG..HANGUL SYLLABLE NIH
	{0xB2E4, 0xB2E4, prLV},                     // Lo       HANGUL SYLLABLE DA
	{0xB2E5, 0xB2FF, prLVT},                    // Lo  [27] HANGUL SYLLABLE DAG..HANGUL SYLLABLE DAH
	{0xB300, 0xB300, prLV},                     // Lo       HANGUL SYLLABLE DAE
	{0xB301, 0xB31B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DAEG..HANGUL SYLLABLE DAEH
	{0xB31C, 0xB31C, prLV},                     // Lo       HANGUL SYLLABLE DYA
	{0xB31D, 0xB337, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYAG..HANGUL SYLLABLE DYAH
	{0xB338, 0xB338, prLV},                     // Lo       HANGUL SYLLABLE DYAE
	{0xB339, 0xB353, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYAEG..HANGUL SYLLABLE DYAEH
	{0xB354, 0xB354, prLV},                     // Lo       HANGUL SYLLABLE DEO
	{0xB355, 0xB36F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DEOG..HANGUL SYLLABLE DEOH
	{0xB370, 0xB370, prLV},                     // Lo       HANGUL SYLLABLE DE
	{0xB371, 0xB38B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DEG..HANGUL SYLLABLE DEH
	{0xB38C, 0xB38C, prLV},                     // Lo       HANGUL SYLLABLE DYEO
	{0xB38D, 0xB3A7, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYEOG..HANGUL SYLLABLE DYEOH
	{0xB3A8, 0xB3A8, prLV},                     // Lo       HANGUL SYLLABLE DYE
	{0xB3A9, 0xB3C3, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYEG..HANGUL SYLLABLE DYEH
	{0xB3C4, 0xB3C4, prLV},                     // Lo       HANGUL SYLLABLE DO
	{0xB3C5, 0xB3DF, prLVT},                    // Lo  [27] HANGUL SYLLABLE DOG..HANGUL SYLLABLE DOH
	{0xB3E0, 0xB3E0, prLV},                     // Lo       HANGUL SYLLABLE DWA
	{0xB3E1, 0xB3FB, prLVT},                    // Lo  [27] HANGUL SYLLABLE DWAG..HANGUL SYLLABLE DWAH
	{0xB3FC, 0xB3FC, prLV},                     // Lo       HANGUL SYLLABLE DWAE
	{0xB3FD, 0xB417, prLVT},                    // Lo  [27] HANGUL SYLLABLE DWAEG..HANGUL SYLLABLE DWAEH
	{0xB418, 0xB418, prLV},                     // Lo       HANGUL SYLLABLE DOE
	{0xB419, 0xB433, prLVT},                    // Lo  [27] HANGUL SYLLABLE DOEG..HANGUL SYLLABLE DOEH
	{0xB434, 0xB434, prLV},                     // Lo       HANGUL SYLLABLE DYO
	{0xB435, 0xB44F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYOG..HANGUL SYLLABLE DYOH
	{0xB450, 0xB450, prLV},                     // Lo       HANGUL SYLLABLE DU
	{0xB451, 0xB46B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DUG..HANGUL SYLLABLE DUH
	{0xB46C, 0xB46C, prLV},                     // Lo       HANGUL SYLLABLE DWEO
	{0xB46D, 0xB487, prLVT},                    // Lo  [27] HANGUL SYLLABLE DWEOG..HANGUL SYLLABLE DWEOH
	{0xB488, 0xB488, prLV},                     // Lo       HANGUL SYLLABLE DWE
	{0xB489, 0xB4A3, prLVT},                    // Lo  [27] HANGUL SYLLABLE DWEG..HANGUL SYLLABLE DWEH
	{0xB4A4, 0xB4A4, prLV},                     // Lo       HANGUL SYLLABLE DWI
	{0xB4A5, 0xB4BF, prLVT},                    // Lo  [27] HANGUL SYLLABLE DWIG..HANGUL SYLLABLE DWIH
	{0xB4C0, 0xB4C0, prLV},                     // Lo       HANGUL SYLLABLE DYU
	{0xB4C1, 0xB4DB, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYUG..HANGUL SYLLABLE DYUH
	{0xB4DC, 0xB4DC, prLV},                     // Lo       HANGUL SYLLABLE DEU
	{0xB4DD, 0xB4F7, prLVT},                    // Lo  [27] HANGUL SYLLABLE DEUG..HANGUL SYLLABLE DEUH
	{0xB4F8, 0xB4F8, prLV},                     // Lo       HANGUL SYLLABLE DYI
	{0xB4F9, 0xB513, prLVT},                    // Lo  [27] HANGUL SYLLABLE DYIG..HANGUL SYLLABLE DYIH
	{0xB514, 0xB514, prLV},                     // Lo       HANGUL SYLLABLE DI
	{0xB515, 0xB52F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DIG..HANGUL SYLLABLE DIH
	{0xB530, 0xB530, prLV},                     // Lo       HANGUL SYLLABLE DDA
	{0xB531, 0xB54B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDAG..HANGUL SYLLABLE DDAH
	{0xB54C, 0xB54C, prLV},                     // Lo       HANGUL SYLLABLE DDAE
	{0xB54D, 0xB567, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDAEG..HANGUL SYLLABLE DDAEH
	{0xB568, 0xB568, prLV},                     // Lo       HANGUL SYLLABLE DDYA
	{0xB569, 0xB583, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYAG..HANGUL SYLLABLE DDYAH
	{0xB584, 0xB584, prLV},                     // Lo       HANGUL SYLLABLE DDYAE
	{0xB585, 0xB59F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYAEG..HANGUL SYLLABLE DDYAEH
	{0xB5A0, 0xB5A0, prLV},                     // Lo       HANGUL SYLLABLE DDEO
	{0xB5A1, 0xB5BB, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDEOG..HANGUL SYLLABLE DDEOH
	{0xB5BC, 0xB5BC, prLV},                     // Lo       HANGUL SYLLABLE DDE
	{0xB5BD, 0xB5D7, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDEG..HANGUL SYLLABLE DDEH
	{0xB5D8, 0xB5D8, prLV},                     // Lo       HANGUL SYLLABLE DDYEO
	{0xB5D9, 0xB5F3, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYEOG..HANGUL SYLLABLE DDYEOH
	{0xB5F4, 0xB5F4, prLV},                     // Lo       HANGUL SYLLABLE DDYE
	{0xB5F5, 0xB60F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYEG..HANGUL SYLLABLE DDYEH
	{0xB610, 0xB610, prLV},                     // Lo       HANGUL SYLLABLE DDO
	{0xB611, 0xB62B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDOG..HANGUL SYLLABLE DDOH
	{0xB62C, 0xB62C, prLV},                     // Lo       HANGUL SYLLABLE DDWA
	{0xB62D, 0xB647, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDWAG..HANGUL SYLLABLE DDWAH
	{0xB648, 0xB648, prLV},                     // Lo       HANGUL SYLLABLE DDWAE
	{0xB649, 0xB663, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDWAEG..HANGUL SYLLABLE DDWAEH
	{0xB664, 0xB664, prLV},                     // Lo       HANGUL SYLLABLE DDOE
	{0xB665, 0xB67F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDOEG..HANGUL SYLLABLE DDOEH
	{0xB680, 0xB680, prLV},                     // Lo       HANGUL SYLLABLE DDYO
	{0xB681, 0xB69B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYOG..HANGUL SYLLABLE DDYOH
	{0xB69C, 0xB69C, prLV},                     // Lo       HANGUL SYLLABLE DDU
	{0xB69D, 0xB6B7, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDUG..HANGUL SYLLABLE DDUH
	{0xB6B8, 0xB6B8, prLV},                     // Lo       HANGUL SYLLABLE DDWEO
	{0xB6B9, 0xB6D3, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDWEOG..HANGUL SYLLABLE DDWEOH
	{0xB6D4, 0xB6D4, prLV},                     // Lo       HANGUL SYLLABLE DDWE
	{0xB6D5, 0xB6EF, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDWEG..HANGUL SYLLABLE DDWEH
	{0xB6F0, 0xB6F0, prLV},                     // Lo       HANGUL SYLLABLE DDWI
	{0xB6F1, 0xB70B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDWIG..HANGUL SYLLABLE DDWIH
	{0xB70C, 0xB70C, prLV},                     // Lo       HANGUL SYLLABLE DDYU
	{0xB70D, 0xB727, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYUG..HANGUL SYLLABLE DDYUH
	{0xB728, 0xB728, prLV},                     // Lo       HANGUL SYLLABLE DDEU
	{0xB729, 0xB743, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDEUG..HANGUL SYLLABLE DDEUH
	{0xB744, 0xB744, prLV},                     // Lo       HANGUL SYLLABLE DDYI
	{0xB745, 0xB75F, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDYIG..HANGUL SYLLABLE DDYIH
	{0xB760, 0xB760, prLV},                     // Lo       HANGUL SYLLABLE DDI
	{0xB761, 0xB77B, prLVT},                    // Lo  [27] HANGUL SYLLABLE DDIG..HANGUL SYLLABLE DDIH
	{0xB77C, 0xB77C, prLV},                     // Lo       HANGUL SYLLABLE RA
	{0xB77D, 0xB797, prLVT},                    // Lo  [27] HANGUL SYLLABLE RAG..HANGUL SYLLABLE RAH
	{0xB798, 0xB798, prLV},                     // Lo       HANGUL SYLLABLE RAE
	{0xB799, 0xB7B3, prLVT},                    // Lo  [27] HANGUL SYLLABLE RAEG..HANGUL SYLLABLE RAEH
	{0xB7B4, 0xB7B4, prLV},                     // Lo       HANGUL SYLLABLE RYA
	{0xB7B5, 0xB7CF, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYAG..HANGUL SYLLABLE RYAH
	{0xB7D0, 0xB7D0, prLV},                     // Lo       HANGUL SYLLABLE RYAE
	{0xB7D1, 0xB7EB, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYAEG..HANGUL SYLLABLE RYAEH
	{0xB7EC, 0xB7EC, prLV},                     // Lo       HANGUL SYLLABLE REO
	{0xB7ED, 0xB807, prLVT},                    // Lo  [27] HANGUL SYLLABLE REOG..HANGUL SYLLABLE REOH
	{0xB808, 0xB808, prLV},                     // Lo       HANGUL SYLLABLE RE
	{0xB809, 0xB823, prLVT},                    // Lo  [27] HANGUL SYLLABLE REG..HANGUL SYLLABLE REH
	{0xB824, 0xB824, prLV},                     // Lo       HANGUL SYLLABLE RYEO
	{0xB825, 0xB83F, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYEOG..HANGUL SYLLABLE RYEOH
	{0xB840, 0xB840, prLV},                     // Lo       HANGUL SYLLABLE RYE
	{0xB841, 0xB85B, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYEG..HANGUL SYLLABLE RYEH
	{0xB85C, 0xB85C, prLV},                     // Lo       HANGUL SYLLABLE RO
	{0xB85D, 0xB877, prLVT},                    // Lo  [27] HANGUL SYLLABLE ROG..HANGUL SYLLABLE ROH
	{0xB878, 0xB878, prLV},                     // Lo       HANGUL SYLLABLE RWA
	{0xB879, 0xB893, prLVT},                    // Lo  [27] HANGUL SYLLABLE RWAG..HANGUL SYLLABLE RWAH
	{0xB894, 0xB894, prLV},                     // Lo       HANGUL SYLLABLE RWAE
	{0xB895, 0xB8AF, prLVT},                    // Lo  [27] HANGUL SYLLABLE RWAEG..HANGUL SYLLABLE RWAEH
	{0xB8B0, 0xB8B0, prLV},                     // Lo       HANGUL SYLLABLE ROE
	{0xB8B1, 0xB8CB, prLVT},                    // Lo  [27] HANGUL SYLLABLE ROEG..HANGUL SYLLABLE ROEH
	{0xB8CC, 0xB8CC, prLV},                     // Lo       HANGUL SYLLABLE RYO
	{0xB8CD, 0xB8E7, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYOG..HANGUL SYLLABLE RYOH
	{0xB8E8, 0xB8E8, prLV},                     // Lo       HANGUL SYLLABLE RU
	{0xB8E9, 0xB903, prLVT},                    // Lo  [27] HANGUL SYLLABLE RUG..HANGUL SYLLABLE RUH
	{0xB904, 0xB904, prLV},                     // Lo       HANGUL SYLLABLE RWEO
	{0xB905, 0xB91F, prLVT},                    // Lo  [27] HANGUL SYLLABLE RWEOG..HANGUL SYLLABLE RWEOH
	{0xB920, 0xB920, prLV},                     // Lo       HANGUL SYLLABLE RWE
	{0xB921, 0xB93B, prLVT},                    // Lo  [27] HANGUL SYLLABLE RWEG..HANGUL SYLLABLE RWEH
	{0xB93C, 0xB93C, prLV},                     // Lo       HANGUL SYLLABLE RWI
	{0xB93D, 0xB957, prLVT},                    // Lo  [27] HANGUL SYLLABLE RWIG..HANGUL SYLLABLE RWIH
	{0xB958, 0xB958, prLV},                     // Lo       HANGUL SYLLABLE RYU
	{0xB959, 0xB973, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYUG..HANGUL SYLLABLE RYUH
	{0xB974, 0xB974, prLV},                     // Lo       HANGUL SYLLABLE REU
	{0xB975, 0xB98F, prLVT},                    // Lo  [27] HANGUL SYLLABLE REUG..HANGUL SYLLABLE REUH
	{0xB990, 0xB990, prLV},                     // Lo       HANGUL SYLLABLE RYI
	{0xB991, 0xB9AB, prLVT},                    // Lo  [27] HANGUL SYLLABLE RYIG..HANGUL SYLLABLE RYIH
	{0xB9AC, 0xB9AC, prLV},                     // Lo       HANGUL SYLLABLE RI
	{0xB9AD, 0xB9C7, prLVT},                    // Lo  [27] HANGUL SYLLABLE RIG..HANGUL SYLLABLE RIH
	{0xB9C8, 0xB9C8, prLV},                     // Lo       HANGUL SYLLABLE MA
	{0xB9C9, 0xB9E3, prLVT},                    // Lo  [27] HANGUL SYLLABLE MAG..HANGUL SYLLABLE MAH
	{0xB9E4, 0xB9E4, prLV},                     // Lo       HANGUL SYLLABLE MAE
	{0xB9E5, 0xB9FF, prLVT},                    // Lo  [27] HANGUL SYLLABLE MAEG..HANGUL SYLLABLE MAEH
	{0xBA00, 0xBA00, prLV},                     // Lo       HANGUL SYLLABLE MYA
	{0xBA01, 0xBA1B, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYAG..HANGUL SYLLABLE MYAH
	{0xBA1C, 0xBA1C, prLV},                     // Lo       HANGUL SYLLABLE MYAE
	{0xBA1D, 0xBA37, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYAEG..HANGUL SYLLABLE MYAEH
	{0xBA38, 0xBA38, prLV},                     // Lo       HANGUL SYLLABLE MEO
	{0xBA39, 0xBA53, prLVT},                    // Lo  [27] HANGUL SYLLABLE MEOG..HANGUL SYLLABLE MEOH
	{0xBA54, 0xBA54, prLV},                     // Lo       HANGUL SYLLABLE ME
	{0xBA55, 0xBA6F, prLVT},                    // Lo  [27] HANGUL SYLLABLE MEG..HANGUL SYLLABLE MEH
	{0xBA70, 0xBA70, prLV},                     // Lo       HANGUL SYLLABLE MYEO
	{0xBA71, 0xBA8B, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYEOG..HANGUL SYLLABLE MYEOH
	{0xBA8C, 0xBA8C, prLV},                     // Lo       HANGUL SYLLABLE MYE
	{0xBA8D, 0xBAA7, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYEG..HANGUL SYLLABLE MYEH
	{0xBAA8, 0xBAA8, prLV},                     // Lo       HANGUL SYLLABLE MO
	{0xBAA9, 0xBAC3, prLVT},                    // Lo  [27] HANGUL SYLLABLE MOG..HANGUL SYLLABLE MOH
	{0xBAC4, 0xBAC4, prLV},                     // Lo       HANGUL SYLLABLE MWA
	{0xBAC5, 0xBADF, prLVT},                    // Lo  [27] HANGUL SYLLABLE MWAG..HANGUL SYLLABLE MWAH
	{0xBAE0, 0xBAE0, prLV},                     // Lo       HANGUL SYLLABLE MWAE
	{0xBAE1, 0xBAFB, prLVT},                    // Lo  [27] HANGUL SYLLABLE MWAEG..HANGUL SYLLABLE MWAEH
	{0xBAFC, 0xBAFC, prLV},                     // Lo       HANGUL SYLLABLE MOE
	{0xBAFD, 0xBB17, prLVT},                    // Lo  [27] HANGUL SYLLABLE MOEG..HANGUL SYLLABLE MOEH
	{0xBB18, 0xBB18, prLV},                     // Lo       HANGUL SYLLABLE MYO
	{0xBB19, 0xBB33, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYOG..HANGUL SYLLABLE MYOH
	{0xBB34, 0xBB34, prLV},                     // Lo       HANGUL SYLLABLE MU
	{0xBB35, 0xBB4F, prLVT},                    // Lo  [27] HANGUL SYLLABLE MUG..HANGUL SYLLABLE MUH
	{0xBB50, 0xBB50, prLV},                     // Lo       HANGUL SYLLABLE MWEO
	{0xBB51, 0xBB6B, prLVT},                    // Lo  [27] HANGUL SYLLABLE MWEOG..HANGUL SYLLABLE MWEOH
	{0xBB6C, 0xBB6C, prLV},                     // Lo       HANGUL SYLLABLE MWE
	{0xBB6D, 0xBB87, prLVT},                    // Lo  [27] HANGUL SYLLABLE MWEG..HANGUL SYLLABLE MWEH
	{0xBB88, 0xBB88, prLV},                     // Lo       HANGUL SYLLABLE MWI
	{0xBB89, 0xBBA3, prLVT},                    // Lo  [27] HANGUL SYLLABLE MWIG..HANGUL SYLLABLE MWIH
	{0xBBA4, 0xBBA4, prLV},                     // Lo       HANGUL SYLLABLE MYU
	{0xBBA5, 0xBBBF, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYUG..HANGUL SYLLABLE MYUH
	{0xBBC0, 0xBBC0, prLV},                     // Lo       HANGUL SYLLABLE MEU
	{0xBBC1, 0xBBDB, prLVT},                    // Lo  [27] HANGUL SYLLABLE MEUG..HANGUL SYLLABLE MEUH
	{0xBBDC, 0xBBDC, prLV},                     // Lo       HANGUL SYLLABLE MYI
	{0xBBDD, 0xBBF7, prLVT},                    // Lo  [27] HANGUL SYLLABLE MYIG..HANGUL SYLLABLE MYIH
	{0xBBF8, 0xBBF8, prLV},                     // Lo       HANGUL SYLLABLE MI
	{0xBBF9, 0xBC13, prLVT},                    // Lo  [27] HANGUL SYLLABLE MIG..HANGUL SYLLABLE MIH
	{0xBC14, 0xBC14, prLV},                     // Lo       HANGUL SYLLABLE BA
	{0xBC15, 0xBC2F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BAG..HANGUL SYLLABLE BAH
	{0xBC30, 0xBC30, prLV},                     // Lo       HANGUL SYLLABLE BAE
	{0xBC31, 0xBC4B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BAEG..HANGUL SYLLABLE BAEH
	{0xBC4C, 0xBC4C, prLV},                     // Lo       HANGUL SYLLABLE BYA
	{0xBC4D, 0xBC67, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYAG..HANGUL SYLLABLE BYAH
	{0xBC68, 0xBC68, prLV},                     // Lo       HANGUL SYLLABLE BYAE
	{0xBC69, 0xBC83, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYAEG..HANGUL SYLLABLE BYAEH
	{0xBC84, 0xBC84, prLV},                     // Lo       HANGUL SYLLABLE BEO
	{0xBC85, 0xBC9F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BEOG..HANGUL SYLLABLE BEOH
	{0xBCA0, 0xBCA0, prLV},                     // Lo       HANGUL SYLLABLE BE
	{0xBCA1, 0xBCBB, prLVT},                    // Lo  [27] HANGUL SYLLABLE BEG..HANGUL SYLLABLE BEH
	{0xBCBC, 0xBCBC, prLV},                     // Lo       HANGUL SYLLABLE BYEO
	{0xBCBD, 0xBCD7, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYEOG..HANGUL SYLLABLE BYEOH
	{0xBCD8, 0xBCD8, prLV},                     // Lo       HANGUL SYLLABLE BYE
	{0xBCD9, 0xBCF3, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYEG..HANGUL SYLLABLE BYEH
	{0xBCF4, 0xBCF4, prLV},                     // Lo       HANGUL SYLLABLE BO
	{0xBCF5, 0xBD0F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BOG..HANGUL SYLLABLE BOH
	{0xBD10, 0xBD10, prLV},                     // Lo       HANGUL SYLLABLE BWA
	{0xBD11, 0xBD2B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BWAG..HANGUL SYLLABLE BWAH
	{0xBD2C, 0xBD2C, prLV},                     // Lo       HANGUL SYLLABLE BWAE
	{0xBD2D, 0xBD47, prLVT},                    // Lo  [27] HANGUL SYLLABLE BWAEG..HANGUL SYLLABLE BWAEH
	{0xBD48, 0xBD48, prLV},                     // Lo       HANGUL SYLLABLE BOE
	{0xBD49, 0xBD63, prLVT},                    // Lo  [27] HANGUL SYLLABLE BOEG..HANGUL SYLLABLE BOEH
	{0xBD64, 0xBD64, prLV},                     // Lo       HANGUL SYLLABLE BYO
	{0xBD65, 0xBD7F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYOG..HANGUL SYLLABLE BYOH
	{0xBD80, 0xBD80, prLV},                     // Lo       HANGUL SYLLABLE BU
	{0xBD81, 0xBD9B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BUG..HANGUL SYLLABLE BUH
	{0xBD9C, 0xBD9C, prLV},                     // Lo       HANGUL SYLLABLE BWEO
	{0xBD9D, 0xBDB7, prLVT},                    // Lo  [27] HANGUL SYLLABLE BWEOG..HANGUL SYLLABLE BWEOH
	{0xBDB8, 0xBDB8, prLV},                     // Lo       HANGUL SYLLABLE BWE
	{0xBDB9, 0xBDD3, prLVT},                    // Lo  [27] HANGUL SYLLABLE BWEG..HANGUL SYLLABLE BWEH
	{0xBDD4, 0xBDD4, prLV},                     // Lo       HANGUL SYLLABLE BWI
	{0xBDD5, 0xBDEF, prLVT},                    // Lo  [27] HANGUL SYLLABLE BWIG..HANGUL SYLLABLE BWIH
	{0xBDF0, 0xBDF0, prLV},                     // Lo       HANGUL SYLLABLE BYU
	{0xBDF1, 0xBE0B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYUG..HANGUL SYLLABLE BYUH
	{0xBE0C, 0xBE0C, prLV},                     // Lo       HANGUL SYLLABLE BEU
	{0xBE0D, 0xBE27, prLVT},                    // Lo  [27] HANGUL SYLLABLE BEUG..HANGUL SYLLABLE BEUH
	{0xBE28, 0xBE28, prLV},                     // Lo       HANGUL SYLLABLE BYI
	{0xBE29, 0xBE43, prLVT},                    // Lo  [27] HANGUL SYLLABLE BYIG..HANGUL SYLLABLE BYIH
	{0xBE44, 0xBE44, prLV},                     // Lo       HANGUL SYLLABLE BI
	{0xBE45, 0xBE5F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BIG..HANGUL SYLLABLE BIH
	{0xBE60, 0xBE60, prLV},                     // Lo       HANGUL SYLLABLE BBA
	{0xBE61, 0xBE7B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBAG..HANGUL SYLLABLE BBAH
	{0xBE7C, 0xBE7C, prLV},                     // Lo       HANGUL SYLLABLE BBAE
	{0xBE7D, 0xBE97, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBAEG..HANGUL SYLLABLE BBAEH
	{0xBE98, 0xBE98, prLV},                     // Lo       HANGUL SYLLABLE BBYA
	{0xBE99, 0xBEB3, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYAG..HANGUL SYLLABLE BBYAH
	{0xBEB4, 0xBEB4, prLV},                     // Lo       HANGUL SYLLABLE BBYAE
	{0xBEB5, 0xBECF, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYAEG..HANGUL SYLLABLE BBYAEH
	{0xBED0, 0xBED0, prLV},                     // Lo       HANGUL SYLLABLE BBEO
	{0xBED1, 0xBEEB, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBEOG..HANGUL SYLLABLE BBEOH
	{0xBEEC, 0xBEEC, prLV},                     // Lo       HANGUL SYLLABLE BBE
	{0xBEED, 0xBF07, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBEG..HANGUL SYLLABLE BBEH
	{0xBF08, 0xBF08, prLV},                     // Lo       HANGUL SYLLABLE BBYEO
	{0xBF09, 0xBF23, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYEOG..HANGUL SYLLABLE BBYEOH
	{0xBF24, 0xBF24, prLV},                     // Lo       HANGUL SYLLABLE BBYE
	{0xBF25, 0xBF3F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYEG..HANGUL SYLLABLE BBYEH
	{0xBF40, 0xBF40, prLV},                     // Lo       HANGUL SYLLABLE BBO
	{0xBF41, 0xBF5B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBOG..HANGUL SYLLABLE BBOH
	{0xBF5C, 0xBF5C, prLV},                     // Lo       HANGUL SYLLABLE BBWA
	{0xBF5D, 0xBF77, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBWAG..HANGUL SYLLABLE BBWAH
	{0xBF78, 0xBF78, prLV},                     // Lo       HANGUL SYLLABLE BBWAE
	{0xBF79, 0xBF93, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBWAEG..HANGUL SYLLABLE BBWAEH
	{0xBF94, 0xBF94, prLV},                     // Lo       HANGUL SYLLABLE BBOE
	{0xBF95, 0xBFAF, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBOEG..HANGUL SYLLABLE BBOEH
	{0xBFB0, 0xBFB0, prLV},                     // Lo       HANGUL SYLLABLE BBYO
	{0xBFB1, 0xBFCB, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYOG..HANGUL SYLLABLE BBYOH
	{0xBFCC, 0xBFCC, prLV},                     // Lo       HANGUL SYLLABLE BBU
	{0xBFCD, 0xBFE7, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBUG..HANGUL SYLLABLE BBUH
	{0xBFE8, 0xBFE8, prLV},                     // Lo       HANGUL SYLLABLE BBWEO
	{0xBFE9, 0xC003, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBWEOG..HANGUL SYLLABLE BBWEOH
	{0xC004, 0xC004, prLV},                     // Lo       HANGUL SYLLABLE BBWE
	{0xC005, 0xC01F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBWEG..HANGUL SYLLABLE BBWEH
	{0xC020, 0xC020, prLV},                     // Lo       HANGUL SYLLABLE BBWI
	{0xC021, 0xC03B, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBWIG..HANGUL SYLLABLE BBWIH
	{0xC03C, 0xC03C, prLV},                     // Lo       HANGUL SYLLABLE BBYU
	{0xC03D, 0xC057, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYUG..HANGUL SYLLABLE BBYUH
	{0xC058, 0xC058, prLV},                     // Lo       HANGUL SYLLABLE BBEU
	{0xC059, 0xC073, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBEUG..HANGUL SYLLABLE BBEUH
	{0xC074, 0xC074, prLV},                     // Lo       HANGUL SYLLABLE BBYI
	{0xC075, 0xC08F, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBYIG..HANGUL SYLLABLE BBYIH
	{0xC090, 0xC090, prLV},                     // Lo       HANGUL SYLLABLE BBI
	{0xC091, 0xC0AB, prLVT},                    // Lo  [27] HANGUL SYLLABLE BBIG..HANGUL SYLLABLE BBIH
	{0xC0AC, 0xC0AC, prLV},                     // Lo       HANGUL SYLLABLE SA
	{0xC0AD, 0xC0C7, prLVT},                    // Lo  [27] HANGUL SYLLABLE SAG..HANGUL SYLLABLE SAH
	{0xC0C8, 0xC0C8, prLV},                     // Lo       HANGUL SYLLABLE SAE
	{0xC0C9, 0xC0E3, prLVT},                    // Lo  [27] HANGUL SYLLABLE SAEG..HANGUL SYLLABLE SAEH
	{0xC0E4, 0xC0E4, prLV},                     // Lo       HANGUL SYLLABLE SYA
	{0xC0E5, 0xC0FF, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYAG..HANGUL SYLLABLE SYAH
	{0xC100, 0xC100, prLV},                     // Lo       HANGUL SYLLABLE SYAE
	{0xC101, 0xC11B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYAEG..HANGUL SYLLABLE SYAEH
	{0xC11C, 0xC11C, prLV},                     // Lo       HANGUL SYLLABLE SEO
	{0xC11D, 0xC137, prLVT},                    // Lo  [27] HANGUL SYLLABLE SEOG..HANGUL SYLLABLE SEOH
	{0xC138, 0xC138, prLV},                     // Lo       HANGUL SYLLABLE SE
	{0xC139, 0xC153, prLVT},                    // Lo  [27] HANGUL SYLLABLE SEG..HANGUL SYLLABLE SEH
	{0xC154, 0xC154, prLV},                     // Lo       HANGUL SYLLABLE SYEO
	{0xC155, 0xC16F, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYEOG..HANGUL SYLLABLE SYEOH
	{0xC170, 0xC170, prLV},                     // Lo       HANGUL SYLLABLE SYE
	{0xC171, 0xC18B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYEG..HANGUL SYLLABLE SYEH
	{0xC18C, 0xC18C, prLV},                     // Lo       HANGUL SYLLABLE SO
	{0xC18D, 0xC1A7, prLVT},                    // Lo  [27] HANGUL SYLLABLE SOG..HANGUL SYLLABLE SOH
	{0xC1A8, 0xC1A8, prLV},                     // Lo       HANGUL SYLLABLE SWA
	{0xC1A9, 0xC1C3, prLVT},                    // Lo  [27] HANGUL SYLLABLE SWAG..HANGUL SYLLABLE SWAH
	{0xC1C4, 0xC1C4, prLV},                     // Lo       HANGUL SYLLABLE SWAE
	{0xC1C5, 0xC1DF, prLVT},                    // Lo  [27] HANGUL SYLLABLE SWAEG..HANGUL SYLLABLE SWAEH
	{0xC1E0, 0xC1E0, prLV},                     // Lo       HANGUL SYLLABLE SOE
	{0xC1E1, 0xC1FB, prLVT},                    // Lo  [27] HANGUL SYLLABLE SOEG..HANGUL SYLLABLE SOEH
	{0xC1FC, 0xC1FC, prLV},                     // Lo       HANGUL SYLLABLE SYO
	{0xC1FD, 0xC217, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYOG..HANGUL SYLLABLE SYOH
	{0xC218, 0xC218, prLV},                     // Lo       HANGUL SYLLABLE SU
	{0xC219, 0xC233, prLVT},                    // Lo  [27] HANGUL SYLLABLE SUG..HANGUL SYLLABLE SUH
	{0xC234, 0xC234, prLV},                     // Lo       HANGUL SYLLABLE SWEO
	{0xC235, 0xC24F, prLVT},                    // Lo  [27] HANGUL SYLLABLE SWEOG..HANGUL SYLLABLE SWEOH
	{0xC250, 0xC250, prLV},                     // Lo       HANGUL SYLLABLE SWE
	{0xC251, 0xC26B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SWEG..HANGUL SYLLABLE SWEH
	{0xC26C, 0xC26C, prLV},                     // Lo       HANGUL SYLLABLE SWI
	{0xC26D, 0xC287, prLVT},                    // Lo  [27] HANGUL SYLLABLE SWIG..HANGUL SYLLABLE SWIH
	{0xC288, 0xC288, prLV},                     // Lo       HANGUL SYLLABLE SYU
	{0xC289, 0xC2A3, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYUG..HANGUL SYLLABLE SYUH
	{0xC2A4, 0xC2A4, prLV},                     // Lo       HANGUL SYLLABLE SEU
	{0xC2A5, 0xC2BF, prLVT},                    // Lo  [27] HANGUL SYLLABLE SEUG..HANGUL SYLLABLE SEUH
	{0xC2C0, 0xC2C0, prLV},                     // Lo       HANGUL SYLLABLE SYI
	{0xC2C1, 0xC2DB, prLVT},                    // Lo  [27] HANGUL SYLLABLE SYIG..HANGUL SYLLABLE SYIH
	{0xC2DC, 0xC2DC, prLV},                     // Lo       HANGUL SYLLABLE SI
	{0xC2DD, 0xC2F7, prLVT},                    // Lo  [27] HANGUL SYLLABLE SIG..HANGUL SYLLABLE SIH
	{0xC2F8, 0xC2F8, prLV},                     // Lo       HANGUL SYLLABLE SSA
	{0xC2F9, 0xC313, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSAG..HANGUL SYLLABLE SSAH
	{0xC314, 0xC314, prLV},                     // Lo       HANGUL SYLLABLE SSAE
	{0xC315, 0xC32F, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSAEG..HANGUL SYLLABLE SSAEH
	{0xC330, 0xC330, prLV},                     // Lo       HANGUL SYLLABLE SSYA
	{0xC331, 0xC34B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYAG..HANGUL SYLLABLE SSYAH
	{0xC34C, 0xC34C, prLV},                     // Lo       HANGUL SYLLABLE SSYAE
	{0xC34D, 0xC367, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYAEG..HANGUL SYLLABLE SSYAEH
	{0xC368, 0xC368, prLV},                     // Lo       HANGUL SYLLABLE SSEO
	{0xC369, 0xC383, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSEOG..HANGUL SYLLABLE SSEOH
	{0xC384, 0xC384, prLV},                     // Lo       HANGUL SYLLABLE SSE
	{0xC385, 0xC39F, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSEG..HANGUL SYLLABLE SSEH
	{0xC3A0, 0xC3A0, prLV},                     // Lo       HANGUL SYLLABLE SSYEO
	{0xC3A1, 0xC3BB, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYEOG..HANGUL SYLLABLE SSYEOH
	{0xC3BC, 0xC3BC, prLV},                     // Lo       HANGUL SYLLABLE SSYE
	{0xC3BD, 0xC3D7, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYEG..HANGUL SYLLABLE SSYEH
	{0xC3D8, 0xC3D8, prLV},                     // Lo       HANGUL SYLLABLE SSO
	{0xC3D9, 0xC3F3, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSOG..HANGUL SYLLABLE SSOH
	{0xC3F4, 0xC3F4, prLV},                     // Lo       HANGUL SYLLABLE SSWA
	{0xC3F5, 0xC40F, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSWAG..HANGUL SYLLABLE SSWAH
	{0xC410, 0xC410, prLV},                     // Lo       HANGUL SYLLABLE SSWAE
	{0xC411, 0xC42B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSWAEG..HANGUL SYLLABLE SSWAEH
	{0xC42C, 0xC42C, prLV},                     // Lo       HANGUL SYLLABLE SSOE
	{0xC42D, 0xC447, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSOEG..HANGUL SYLLABLE SSOEH
	{0xC448, 0xC448, prLV},                     // Lo       HANGUL SYLLABLE SSYO
	{0xC449, 0xC463, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYOG..HANGUL SYLLABLE SSYOH
	{0xC464, 0xC464, prLV},                     // Lo       HANGUL SYLLABLE SSU
	{0xC465, 0xC47F, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSUG..HANGUL SYLLABLE SSUH
	{0xC480, 0xC480, prLV},                     // Lo       HANGUL SYLLABLE SSWEO
	{0xC481, 0xC49B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSWEOG..HANGUL SYLLABLE SSWEOH
	{0xC49C, 0xC49C, prLV},                     // Lo       HANGUL SYLLABLE SSWE
	{0xC49D, 0xC4B7, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSWEG..HANGUL SYLLABLE SSWEH
	{0xC4B8, 0xC4B8, prLV},                     // Lo       HANGUL SYLLABLE SSWI
	{0xC4B9, 0xC4D3, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSWIG..HANGUL SYLLABLE SSWIH
	{0xC4D4, 0xC4D4, prLV},                     // Lo       HANGUL SYLLABLE SSYU
	{0xC4D5, 0xC4EF, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYUG..HANGUL SYLLABLE SSYUH
	{0xC4F0, 0xC4F0, prLV},                     // Lo       HANGUL SYLLABLE SSEU
	{0xC4F1, 0xC50B, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSEUG..HANGUL SYLLABLE SSEUH
	{0xC50C, 0xC50C, prLV},                     // Lo       HANGUL SYLLABLE SSYI
	{0xC50D, 0xC527, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSYIG..HANGUL SYLLABLE SSYIH
	{0xC528, 0xC528, prLV},                     // Lo       HANGUL SYLLABLE SSI
	{0xC529, 0xC543, prLVT},                    // Lo  [27] HANGUL SYLLABLE SSIG..HANGUL SYLLABLE SSIH
	{0xC544, 0xC544, prLV},                     // Lo       HANGUL SYLLABLE A
	{0xC545, 0xC55F, prLVT},                    // Lo  [27] HANGUL SYLLABLE AG..HANGUL SYLLABLE AH
	{0xC560, 0xC560, prLV},                     // Lo       HANGUL SYLLABLE AE
	{0xC561, 0xC57B, prLVT},                    // Lo  [27] HANGUL SYLLABLE AEG..HANGUL SYLLABLE AEH
	{0xC57C, 0xC57C, prLV},                     // Lo       HANGUL SYLLABLE YA
	{0xC57D, 0xC597, prLVT},                    // Lo  [27] HANGUL SYLLABLE YAG..HANGUL SYLLABLE YAH
	{0xC598, 0xC598, prLV},                     // Lo       HANGUL SYLLABLE YAE
	{0xC599, 0xC5B3, prLVT},                    // Lo  [27] HANGUL SYLLABLE YAEG..HANGUL SYLLABLE YAEH
	{0xC5B4, 0xC5B4, prLV},                     // Lo       HANGUL SYLLABLE EO
	{0xC5B5, 0xC5CF, prLVT},                    // Lo  [27] HANGUL SYLLABLE EOG..HANGUL SYLLABLE EOH
	{0xC5D0, 0xC5D0, prLV},                     // Lo       HANGUL SYLLABLE E
	{0xC5D1, 0xC5EB, prLVT},                    // Lo  [27] HANGUL SYLLABLE EG..HANGUL SYLLABLE EH
	{0xC5EC, 0xC5EC, prLV},                     // Lo       HANGUL SYLLABLE YEO
	{0xC5ED, 0xC607, prLVT},                    // Lo  [27] HANGUL SYLLABLE YEOG..HANGUL SYLLABLE YEOH
	{0xC608, 0xC608, prLV},                     // Lo       HANGUL SYLLABLE YE
	{0xC609, 0xC623, prLVT},                    // Lo  [27] HANGUL SYLLABLE YEG..HANGUL SYLLABLE YEH
	{0xC624, 0xC624, prLV},                     // Lo       HANGUL SYLLABLE O
	{0xC625, 0xC63F, prLVT},                    // Lo  [27] HANGUL SYLLABLE OG..HANGUL SYLLABLE OH
	{0xC640, 0xC640, prLV},                     // Lo       HANGUL SYLLABLE WA
	{0xC641, 0xC65B, prLVT},                    // Lo  [27] HANGUL SYLLABLE WAG..HANGUL SYLLABLE WAH
	{0xC65C, 0xC65C, prLV},                     // Lo       HANGUL SYLLABLE WAE
	{0xC65D, 0xC677, prLVT},                    // Lo  [27] HANGUL SYLLABLE WAEG..HANGUL SYLLABLE WAEH
	{0xC678, 0xC678, prLV},                     // Lo       HANGUL SYLLABLE OE
	{0xC679, 0xC693, prLVT},                    // Lo  [27] HANGUL SYLLABLE OEG..HANGUL SYLLABLE OEH
	{0xC694, 0xC694, prLV},                     // Lo       HANGUL SYLLABLE YO
	{0xC695, 0xC6AF, prLVT},                    // Lo  [27] HANGUL SYLLABLE YOG..HANGUL SYLLABLE YOH
	{0xC6B0, 0xC6B0, prLV},                     // Lo       HANGUL SYLLABLE U
	{0xC6B1, 0xC6CB, prLVT},                    // Lo  [27] HANGUL SYLLABLE UG..HANGUL SYLLABLE UH
	{0xC6CC, 0xC6CC, prLV},                     // Lo       HANGUL SYLLABLE WEO
	{0xC6CD, 0xC6E7, prLVT},                    // Lo  [27] HANGUL SYLLABLE WEOG..HANGUL SYLLABLE WEOH
	{0xC6E8, 0xC6E8, prLV},                     // Lo       HANGUL SYLLABLE WE
	{0xC6E9, 0xC703, prLVT},                    // Lo  [27] HANGUL SYLLABLE WEG..HANGUL SYLLABLE WEH
	{0xC704, 0xC704, prLV},                     // Lo       HANGUL SYLLABLE WI
	{0xC705, 0xC71F, prLVT},                    // Lo  [27] HANGUL SYLLABLE WIG..HANGUL SYLLABLE WIH
	{0xC720, 0xC720, prLV},                     // Lo       HANGUL SYLLABLE YU
	{0xC721, 0xC73B, prLVT},                    // Lo  [27] HANGUL SYLLABLE YUG..HANGUL SYLLABLE YUH
	{0xC73C, 0xC73C, prLV},                     // Lo       HANGUL SYLLABLE EU
	{0xC73D, 0xC757, prLVT},                    // Lo  [27] HANGUL SYLLABLE EUG..HANGUL SYLLABLE EUH
	{0xC758, 0xC758, prLV},                     // Lo       HANGUL SYLLABLE YI
	{0xC759, 0xC773, prLVT},                    // Lo  [27] HANGUL SYLLABLE YIG..HANGUL SYLLABLE YIH
	{0xC774, 0xC774, prLV},                     // Lo       HANGUL SYLLABLE I
	{0xC775, 0xC78F, prLVT},                    // Lo  [27] HANGUL SYLLABLE IG..HANGUL SYLLABLE IH
	{0xC790, 0xC790, prLV},                     // Lo       HANGUL SYLLABLE JA
	{0xC791, 0xC7AB, prLVT},                    // Lo  [27] HANGUL SYLLABLE JAG..HANGUL SYLLABLE JAH
	{0xC7AC, 0xC7AC, prLV},                     // Lo       HANGUL SYLLABLE JAE
	{0xC7AD, 0xC7C7, prLVT},                    // Lo  [27] HANGUL SYLLABLE JAEG..HANGUL SYLLABLE JAEH
	{0xC7C8, 0xC7C8, prLV},                     // Lo       HANGUL SYLLABLE JYA
	{0xC7C9, 0xC7E3, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYAG..HANGUL SYLLABLE JYAH
	{0xC7E4, 0xC7E4, prLV},                     // Lo       HANGUL SYLLABLE JYAE
	{0xC7E5, 0xC7FF, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYAEG..HANGUL SYLLABLE JYAEH
	{0xC800, 0xC800, prLV},                     // Lo       HANGUL SYLLABLE JEO
	{0xC801, 0xC81B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JEOG..HANGUL SYLLABLE JEOH
	{0xC81C, 0xC81C, prLV},                     // Lo       HANGUL SYLLABLE JE
	{0xC81D, 0xC837, prLVT},                    // Lo  [27] HANGUL SYLLABLE JEG..HANGUL SYLLABLE JEH
	{0xC838, 0xC838, prLV},                     // Lo       HANGUL SYLLABLE JYEO
	{0xC839, 0xC853, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYEOG..HANGUL SYLLABLE JYEOH
	{0xC854, 0xC854, prLV},                     // Lo       HANGUL SYLLABLE JYE
	{0xC855, 0xC86F, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYEG..HANGUL SYLLABLE JYEH
	{0xC870, 0xC870, prLV},                     // Lo       HANGUL SYLLABLE JO
	{0xC871, 0xC88B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JOG..HANGUL SYLLABLE JOH
	{0xC88C, 0xC88C, prLV},                     // Lo       HANGUL SYLLABLE JWA
	{0xC88D, 0xC8A7, prLVT},                    // Lo  [27] HANGUL SYLLABLE JWAG..HANGUL SYLLABLE JWAH
	{0xC8A8, 0xC8A8, prLV},                     // Lo       HANGUL SYLLABLE JWAE
	{0xC8A9, 0xC8C3, prLVT},                    // Lo  [27] HANGUL SYLLABLE JWAEG..HANGUL SYLLABLE JWAEH
	{0xC8C4, 0xC8C4, prLV},                     // Lo       HANGUL SYLLABLE JOE
	{0xC8C5, 0xC8DF, prLVT},                    // Lo  [27] HANGUL SYLLABLE JOEG..HANGUL SYLLABLE JOEH
	{0xC8E0, 0xC8E0, prLV},                     // Lo       HANGUL SYLLABLE JYO
	{0xC8E1, 0xC8FB, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYOG..HANGUL SYLLABLE JYOH
	{0xC8FC, 0xC8FC, prLV},                     // Lo       HANGUL SYLLABLE JU
	{0xC8FD, 0xC917, prLVT},                    // Lo  [27] HANGUL SYLLABLE JUG..HANGUL SYLLABLE JUH
	{0xC918, 0xC918, prLV},                     // Lo       HANGUL SYLLABLE JWEO
	{0xC919, 0xC933, prLVT},                    // Lo  [27] HANGUL SYLLABLE JWEOG..HANGUL SYLLABLE JWEOH
	{0xC934, 0xC934, prLV},                     // Lo       HANGUL SYLLABLE JWE
	{0xC935, 0xC94F, prLVT},                    // Lo  [27] HANGUL SYLLABLE JWEG..HANGUL SYLLABLE JWEH
	{0xC950, 0xC950, prLV},                     // Lo       HANGUL SYLLABLE JWI
	{0xC951, 0xC96B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JWIG..HANGUL SYLLABLE JWIH
	{0xC96C, 0xC96C, prLV},                     // Lo       HANGUL SYLLABLE JYU
	{0xC96D, 0xC987, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYUG..HANGUL SYLLABLE JYUH
	{0xC988, 0xC988, prLV},                     // Lo       HANGUL SYLLABLE JEU
	{0xC989, 0xC9A3, prLVT},                    // Lo  [27] HANGUL SYLLABLE JEUG..HANGUL SYLLABLE JEUH
	{0xC9A4, 0xC9A4, prLV},                     // Lo       HANGUL SYLLABLE JYI
	{0xC9A5, 0xC9BF, prLVT},                    // Lo  [27] HANGUL SYLLABLE JYIG..HANGUL SYLLABLE JYIH
	{0xC9C0, 0xC9C0, prLV},                     // Lo       HANGUL SYLLABLE JI
	{0xC9C1, 0xC9DB, prLVT},                    // Lo  [27] HANGUL SYLLABLE JIG..HANGUL SYLLABLE JIH
	{0xC9DC, 0xC9DC, prLV},                     // Lo       HANGUL SYLLABLE JJA
	{0xC9DD, 0xC9F7, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJAG..HANGUL SYLLABLE JJAH
	{0xC9F8, 0xC9F8, prLV},                     // Lo       HANGUL SYLLABLE JJAE
	{0xC9F9, 0xCA13, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJAEG..HANGUL SYLLABLE JJAEH
	{0xCA14, 0xCA14, prLV},                     // Lo       HANGUL SYLLABLE JJYA
	{0xCA15, 0xCA2F, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYAG..HANGUL SYLLABLE JJYAH
	{0xCA30, 0xCA30, prLV},                     // Lo       HANGUL SYLLABLE JJYAE
	{0xCA31, 0xCA4B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYAEG..HANGUL SYLLABLE JJYAEH
	{0xCA4C, 0xCA4C, prLV},                     // Lo       HANGUL SYLLABLE JJEO
	{0xCA4D, 0xCA67, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJEOG..HANGUL SYLLABLE JJEOH
	{0xCA68, 0xCA68, prLV},                     // Lo       HANGUL SYLLABLE JJE
	{0xCA69, 0xCA83, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJEG..HANGUL SYLLABLE JJEH
	{0xCA84, 0xCA84, prLV},                     // Lo       HANGUL SYLLABLE JJYEO
	{0xCA85, 0xCA9F, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYEOG..HANGUL SYLLABLE JJYEOH
	{0xCAA0, 0xCAA0, prLV},                     // Lo       HANGUL SYLLABLE JJYE
	{0xCAA1, 0xCABB, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYEG..HANGUL SYLLABLE JJYEH
	{0xCABC, 0xCABC, prLV},                     // Lo       HANGUL SYLLABLE JJO
	{0xCABD, 0xCAD7, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJOG..HANGUL SYLLABLE JJOH
	{0xCAD8, 0xCAD8, prLV},                     // Lo       HANGUL SYLLABLE JJWA
	{0xCAD9, 0xCAF3, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJWAG..HANGUL SYLLABLE JJWAH
	{0xCAF4, 0xCAF4, prLV},                     // Lo       HANGUL SYLLABLE JJWAE
	{0xCAF5, 0xCB0F, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJWAEG..HANGUL SYLLABLE JJWAEH
	{0xCB10, 0xCB10, prLV},                     // Lo       HANGUL SYLLABLE JJOE
	{0xCB11, 0xCB2B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJOEG..HANGUL SYLLABLE JJOEH
	{0xCB2C, 0xCB2C, prLV},                     // Lo       HANGUL SYLLABLE JJYO
	{0xCB2D, 0xCB47, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYOG..HANGUL SYLLABLE JJYOH
	{0xCB48, 0xCB48, prLV},                     // Lo       HANGUL SYLLABLE JJU
	{0xCB49, 0xCB63, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJUG..HANGUL SYLLABLE JJUH
	{0xCB64, 0xCB64, prLV},                     // Lo       HANGUL SYLLABLE JJWEO
	{0xCB65, 0xCB7F, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJWEOG..HANGUL SYLLABLE JJWEOH
	{0xCB80, 0xCB80, prLV},                     // Lo       HANGUL SYLLABLE JJWE
	{0xCB81, 0xCB9B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJWEG..HANGUL SYLLABLE JJWEH
	{0xCB9C, 0xCB9C, prLV},                     // Lo       HANGUL SYLLABLE JJWI
	{0xCB9D, 0xCBB7, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJWIG..HANGUL SYLLABLE JJWIH
	{0xCBB8, 0xCBB8, prLV},                     // Lo       HANGUL SYLLABLE JJYU
	{0xCBB9, 0xCBD3, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYUG..HANGUL SYLLABLE JJYUH
	{0xCBD4, 0xCBD4, prLV},                     // Lo       HANGUL SYLLABLE JJEU
	{0xCBD5, 0xCBEF, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJEUG..HANGUL SYLLABLE JJEUH
	{0xCBF0, 0xCBF0, prLV},                     // Lo       HANGUL SYLLABLE JJYI
	{0xCBF1, 0xCC0B, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJYIG..HANGUL SYLLABLE JJYIH
	{0xCC0C, 0xCC0C, prLV},                     // Lo       HANGUL SYLLABLE JJI
	{0xCC0D, 0xCC27, prLVT},                    // Lo  [27] HANGUL SYLLABLE JJIG..HANGUL SYLLABLE JJIH
	{0xCC28, 0xCC28, prLV},                     // Lo       HANGUL SYLLABLE CA
	{0xCC29, 0xCC43, prLVT},                    // Lo  [27] HANGUL SYLLABLE CAG..HANGUL SYLLABLE CAH
	{0xCC44, 0xCC44, prLV},                     // Lo       HANGUL SYLLABLE CAE
	{0xCC45, 0xCC5F, prLVT},                    // Lo  [27] HANGUL SYLLABLE CAEG..HANGUL SYLLABLE CAEH
	{0xCC60, 0xCC60, prLV},                     // Lo       HANGUL SYLLABLE CYA
	{0xCC61, 0xCC7B, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYAG..HANGUL SYLLABLE CYAH
	{0xCC7C, 0xCC7C, prLV},                     // Lo       HANGUL SYLLABLE CYAE
	{0xCC7D, 0xCC97, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYAEG..HANGUL SYLLABLE CYAEH
	{0xCC98, 0xCC98, prLV},                     // Lo       HANGUL SYLLABLE CEO
	{0xCC99, 0xCCB3, prLVT},                    // Lo  [27] HANGUL SYLLABLE CEOG..HANGUL SYLLABLE CEOH
	{0xCCB4, 0xCCB4, prLV},                     // Lo       HANGUL SYLLABLE CE
	{0xCCB5, 0xCCCF, prLVT},                    // Lo  [27] HANGUL SYLLABLE CEG..HANGUL SYLLABLE CEH
	{0xCCD0, 0xCCD0, prLV},                     // Lo       HANGUL SYLLABLE CYEO
	{0xCCD1, 0xCCEB, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYEOG..HANGUL SYLLABLE CYEOH
	{0xCCEC, 0xCCEC, prLV},                     // Lo       HANGUL SYLLABLE CYE
	{0xCCED, 0xCD07, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYEG..HANGUL SYLLABLE CYEH
	{0xCD08, 0xCD08, prLV},                     // Lo       HANGUL SYLLABLE CO
	{0xCD09, 0xCD23, prLVT},                    // Lo  [27] HANGUL SYLLABLE COG..HANGUL SYLLABLE COH
	{0xCD24, 0xCD24, prLV},                     // Lo       HANGUL SYLLABLE CWA
	{0xCD25, 0xCD3F, prLVT},                    // Lo  [27] HANGUL SYLLABLE CWAG..HANGUL SYLLABLE CWAH
	{0xCD40, 0xCD40, prLV},                     // Lo       HANGUL SYLLABLE CWAE
	{0xCD41, 0xCD5B, prLVT},                    // Lo  [27] HANGUL SYLLABLE CWAEG..HANGUL SYLLABLE CWAEH
	{0xCD5C, 0xCD5C, prLV},                     // Lo       HANGUL SYLLABLE COE
	{0xCD5D, 0xCD77, prLVT},                    // Lo  [27] HANGUL SYLLABLE COEG..HANGUL SYLLABLE COEH
	{0xCD78, 0xCD78, prLV},                     // Lo       HANGUL SYLLABLE CYO
	{0xCD79, 0xCD93, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYOG..HANGUL SYLLABLE CYOH
	{0xCD94, 0xCD94, prLV},                     // Lo       HANGUL SYLLABLE CU
	{0xCD95, 0xCDAF, prLVT},                    // Lo  [27] HANGUL SYLLABLE CUG..HANGUL SYLLABLE CUH
	{0xCDB0, 0xCDB0, prLV},                     // Lo       HANGUL SYLLABLE CWEO
	{0xCDB1, 0xCDCB, prLVT},                    // Lo  [27] HANGUL SYLLABLE CWEOG..HANGUL SYLLABLE CWEOH
	{0xCDCC, 0xCDCC, prLV},                     // Lo       HANGUL SYLLABLE CWE
	{0xCDCD, 0xCDE7, prLVT},                    // Lo  [27] HANGUL SYLLABLE CWEG..HANGUL SYLLABLE CWEH
	{0xCDE8, 0xCDE8, prLV},                     // Lo       HANGUL SYLLABLE CWI
	{0xCDE9, 0xCE03, prLVT},                    // Lo  [27] HANGUL SYLLABLE CWIG..HANGUL SYLLABLE CWIH
	{0xCE04, 0xCE04, prLV},                     // Lo       HANGUL SYLLABLE CYU
	{0xCE05, 0xCE1F, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYUG..HANGUL SYLLABLE CYUH
	{0xCE20, 0xCE20, prLV},                     // Lo       HANGUL SYLLABLE CEU
	{0xCE21, 0xCE3B, prLVT},                    // Lo  [27] HANGUL SYLLABLE CEUG..HANGUL SYLLABLE CEUH
	{0xCE3C, 0xCE3C, prLV},                     // Lo       HANGUL SYLLABLE CYI
	{0xCE3D, 0xCE57, prLVT},                    // Lo  [27] HANGUL SYLLABLE CYIG..HANGUL SYLLABLE CYIH
	{0xCE58, 0xCE58, prLV},                     // Lo       HANGUL SYLLABLE CI
	{0xCE59, 0xCE73, prLVT},                    // Lo  [27] HANGUL SYLLABLE CIG..HANGUL SYLLABLE CIH
	{0xCE74, 0xCE74, prLV},                     // Lo       HANGUL SYLLABLE KA
	{0xCE75, 0xCE8F, prLVT},                    // Lo  [27] HANGUL SYLLABLE KAG..HANGUL SYLLABLE KAH
	{0xCE90, 0xCE90, prLV},                     // Lo       HANGUL SYLLABLE KAE
	{0xCE91, 0xCEAB, prLVT},                    // Lo  [27] HANGUL SYLLABLE KAEG..HANGUL SYLLABLE KAEH
	{0xCEAC, 0xCEAC, prLV},                     // Lo       HANGUL SYLLABLE KYA
	{0xCEAD, 0xCEC7, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYAG..HANGUL SYLLABLE KYAH
	{0xCEC8, 0xCEC8, prLV},                     // Lo       HANGUL SYLLABLE KYAE
	{0xCEC9, 0xCEE3, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYAEG..HANGUL SYLLABLE KYAEH
	{0xCEE4, 0xCEE4, prLV},                     // Lo       HANGUL SYLLABLE KEO
	{0xCEE5, 0xCEFF, prLVT},                    // Lo  [27] HANGUL SYLLABLE KEOG..HANGUL SYLLABLE KEOH
	{0xCF00, 0xCF00, prLV},                     // Lo       HANGUL SYLLABLE KE
	{0xCF01, 0xCF1B, prLVT},                    // Lo  [27] HANGUL SYLLABLE KEG..HANGUL SYLLABLE KEH
	{0xCF1C, 0xCF1C, prLV},                     // Lo       HANGUL SYLLABLE KYEO
	{0xCF1D, 0xCF37, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYEOG..HANGUL SYLLABLE KYEOH
	{0xCF38, 0xCF38, prLV},                     // Lo       HANGUL SYLLABLE KYE
	{0xCF39, 0xCF53, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYEG..HANGUL SYLLABLE KYEH
	{0xCF54, 0xCF54, prLV},                     // Lo       HANGUL SYLLABLE KO
	{0xCF55, 0xCF6F, prLVT},                    // Lo  [27] HANGUL SYLLABLE KOG..HANGUL SYLLABLE KOH
	{0xCF70, 0xCF70, prLV},                     // Lo       HANGUL SYLLABLE KWA
	{0xCF71, 0xCF8B, prLVT},                    // Lo  [27] HANGUL SYLLABLE KWAG..HANGUL SYLLABLE KWAH
	{0xCF8C, 0xCF8C, prLV},                     // Lo       HANGUL SYLLABLE KWAE
	{0xCF8D, 0xCFA7, prLVT},                    // Lo  [27] HANGUL SYLLABLE KWAEG..HANGUL SYLLABLE KWAEH
	{0xCFA8, 0xCFA8, prLV},                     // Lo       HANGUL SYLLABLE KOE
	{0xCFA9, 0xCFC3, prLVT},                    // Lo  [27] HANGUL SYLLABLE KOEG..HANGUL SYLLABLE KOEH
	{0xCFC4, 0xCFC4, prLV},                     // Lo       HANGUL SYLLABLE KYO
	{0xCFC5, 0xCFDF, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYOG..HANGUL SYLLABLE KYOH
	{0xCFE0, 0xCFE0, prLV},                     // Lo       HANGUL SYLLABLE KU
	{0xCFE1, 0xCFFB, prLVT},                    // Lo  [27] HANGUL SYLLABLE KUG..HANGUL SYLLABLE KUH
	{0xCFFC, 0xCFFC, prLV},                     // Lo       HANGUL SYLLABLE KWEO
	{0xCFFD, 0xD017, prLVT},                    // Lo  [27] HANGUL SYLLABLE KWEOG..HANGUL SYLLABLE KWEOH
	{0xD018, 0xD018, prLV},                     // Lo       HANGUL SYLLABLE KWE
	{0xD019, 0xD033, prLVT},                    // Lo  [27] HANGUL SYLLABLE KWEG..HANGUL SYLLABLE KWEH
	{0xD034, 0xD034, prLV},                     // Lo       HANGUL SYLLABLE KWI
	{0xD035, 0xD04F, prLVT},                    // Lo  [27] HANGUL SYLLABLE KWIG..HANGUL SYLLABLE KWIH
	{0xD050, 0xD050, prLV},                     // Lo       HANGUL SYLLABLE KYU
	{0xD051, 0xD06B, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYUG..HANGUL SYLLABLE KYUH
	{0xD06C, 0xD06C, prLV},                     // Lo       HANGUL SYLLABLE KEU
	{0xD06D, 0xD087, prLVT},                    // Lo  [27] HANGUL SYLLABLE KEUG..HANGUL SYLLABLE KEUH
	{0xD088, 0xD088, prLV},                     // Lo       HANGUL SYLLABLE KYI
	{0xD089, 0xD0A3, prLVT},                    // Lo  [27] HANGUL SYLLABLE KYIG..HANGUL SYLLABLE KYIH
	{0xD0A4, 0xD0A4, prLV},                     // Lo       HANGUL SYLLABLE KI
	{0xD0A5, 0xD0BF, prLVT},                    // Lo  [27] HANGUL SYLLABLE KIG..HANGUL SYLLABLE KIH
	{0xD0C0, 0xD0C0, prLV},                     // Lo       HANGUL SYLLABLE TA
	{0xD0C1, 0xD0DB, prLVT},                    // Lo  [27] HANGUL SYLLABLE TAG..HANGUL SYLLABLE TAH
	{0xD0DC, 0xD0DC, prLV},                     // Lo       HANGUL SYLLABLE TAE
	{0xD0DD, 0xD0F7, prLVT},                    // Lo  [27] HANGUL SYLLABLE TAEG..HANGUL SYLLABLE TAEH
	{0xD0F8, 0xD0F8, prLV},                     // Lo       HANGUL SYLLABLE TYA
	{0xD0F9, 0xD113, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYAG..HANGUL SYLLABLE TYAH
	{0xD114, 0xD114, prLV},                     // Lo       HANGUL SYLLABLE TYAE
	{0xD115, 0xD12F, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYAEG..HANGUL SYLLABLE TYAEH
	{0xD130, 0xD130, prLV},                     // Lo       HANGUL SYLLABLE TEO
	{0xD131, 0xD14B, prLVT},                    // Lo  [27] HANGUL SYLLABLE TEOG..HANGUL SYLLABLE TEOH
	{0xD14C, 0xD14C, prLV},                     // Lo       HANGUL SYLLABLE TE
	{0xD14D, 0xD167, prLVT},                    // Lo  [27] HANGUL SYLLABLE TEG..HANGUL SYLLABLE TEH
	{0xD168, 0xD168, prLV},                     // Lo       HANGUL SYLLABLE TYEO
	{0xD169, 0xD183, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYEOG..HANGUL SYLLABLE TYEOH
	{0xD184, 0xD184, prLV},                     // Lo       HANGUL SYLLABLE TYE
	{0xD185, 0xD19F, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYEG..HANGUL SYLLABLE TYEH
	{0xD1A0, 0xD1A0, prLV},                     // Lo       HANGUL SYLLABLE TO
	{0xD1A1, 0xD1BB, prLVT},                    // Lo  [27] HANGUL SYLLABLE TOG..HANGUL SYLLABLE TOH
	{0xD1BC, 0xD1BC, prLV},                     // Lo       HANGUL SYLLABLE TWA
	{0xD1BD, 0xD1D7, prLVT},                    // Lo  [27] HANGUL SYLLABLE TWAG..HANGUL SYLLABLE TWAH
	{0xD1D8, 0xD1D8, prLV},                     // Lo       HANGUL SYLLABLE TWAE
	{0xD1D9, 0xD1F3, prLVT},                    // Lo  [27] HANGUL SYLLABLE TWAEG..HANGUL SYLLABLE TWAEH
	{0xD1F4, 0xD1F4, prLV},                     // Lo       HANGUL SYLLABLE TOE
	{0xD1F5, 0xD20F, prLVT},                    // Lo  [27] HANGUL SYLLABLE TOEG..HANGUL SYLLABLE TOEH
	{0xD210, 0xD210, prLV},                     // Lo       HANGUL SYLLABLE TYO
	{0xD211, 0xD22B, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYOG..HANGUL SYLLABLE TYOH
	{0xD22C, 0xD22C, prLV},                     // Lo       HANGUL SYLLABLE TU
	{0xD22D, 0xD247, prLVT},                    // Lo  [27] HANGUL SYLLABLE TUG..HANGUL SYLLABLE TUH
	{0xD248, 0xD248, prLV},                     // Lo       HANGUL SYLLABLE TWEO
	{0xD249, 0xD263, prLVT},                    // Lo  [27] HANGUL SYLLABLE TWEOG..HANGUL SYLLABLE TWEOH
	{0xD264, 0xD264, prLV},                     // Lo       HANGUL SYLLABLE TWE
	{0xD265, 0xD27F, prLVT},                    // Lo  [27] HANGUL SYLLABLE TWEG..HANGUL SYLLABLE TWEH
	{0xD280, 0xD280, prLV},                     // Lo       HANGUL SYLLABLE TWI
	{0xD281, 0xD29B, prLVT},                    // Lo  [27] HANGUL SYLLABLE TWIG..HANGUL SYLLABLE TWIH
	{0xD29C, 0xD29C, prLV},                     // Lo       HANGUL SYLLABLE TYU
	{0xD29D, 0xD2B7, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYUG..HANGUL SYLLABLE TYUH
	{0xD2B8, 0xD2B8, prLV},                     // Lo       HANGUL SYLLABLE TEU
	{0xD2B9, 0xD2D3, prLVT},                    // Lo  [27] HANGUL SYLLABLE TEUG..HANGUL SYLLABLE TEUH
	{0xD2D4, 0xD2D4, prLV},                     // Lo       HANGUL SYLLABLE TYI
	{0xD2D5, 0xD2EF, prLVT},                    // Lo  [27] HANGUL SYLLABLE TYIG..HANGUL SYLLABLE TYIH
	{0xD2F0, 0xD2F0, prLV},                     // Lo       HANGUL SYLLABLE TI
	{0xD2F1, 0xD30B, prLVT},                    // Lo  [27] HANGUL SYLLABLE TIG..HANGUL SYLLABLE TIH
	{0xD30C, 0xD30C, prLV},                     // Lo       HANGUL SYLLABLE PA
	{0xD30D, 0xD327, prLVT},                    // Lo  [27] HANGUL SYLLABLE PAG..HANGUL SYLLABLE PAH
	{0xD328, 0xD328, prLV},                     // Lo       HANGUL SYLLABLE PAE
	{0xD329, 0xD343, prLVT},                    // Lo  [27] HANGUL SYLLABLE PAEG..HANGUL SYLLABLE PAEH
	{0xD344, 0xD344, prLV},                     // Lo       HANGUL SYLLABLE PYA
	{0xD345, 0xD35F, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYAG..HANGUL SYLLABLE PYAH
	{0xD360, 0xD360, prLV},                     // Lo       HANGUL SYLLABLE PYAE
	{0xD361, 0xD37B, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYAEG..HANGUL SYLLABLE PYAEH
	{0xD37C, 0xD37C, prLV},                     // Lo       HANGUL SYLLABLE PEO
	{0xD37D, 0xD397, prLVT},                    // Lo  [27] HANGUL SYLLABLE PEOG..HANGUL SYLLABLE PEOH
	{0xD398, 0xD398, prLV},                     // Lo       HANGUL SYLLABLE PE
	{0xD399, 0xD3B3, prLVT},                    // Lo  [27] HANGUL SYLLABLE PEG..HANGUL SYLLABLE PEH
	{0xD3B4, 0xD3B4, prLV},                     // Lo       HANGUL SYLLABLE PYEO
	{0xD3B5, 0xD3CF, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYEOG..HANGUL SYLLABLE PYEOH
	{0xD3D0, 0xD3D0, prLV},                     // Lo       HANGUL SYLLABLE PYE
	{0xD3D1, 0xD3EB, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYEG..HANGUL SYLLABLE PYEH
	{0xD3EC, 0xD3EC, prLV},                     // Lo       HANGUL SYLLABLE PO
	{0xD3ED, 0xD407, prLVT},                    // Lo  [27] HANGUL SYLLABLE POG..HANGUL SYLLABLE POH
	{0xD408, 0xD408, prLV},                     // Lo       HANGUL SYLLABLE PWA
	{0xD409, 0xD423, prLVT},                    // Lo  [27] HANGUL SYLLABLE PWAG..HANGUL SYLLABLE PWAH
	{0xD424, 0xD424, prLV},                     // Lo       HANGUL SYLLABLE PWAE
	{0xD425, 0xD43F, prLVT},                    // Lo  [27] HANGUL SYLLABLE PWAEG..HANGUL SYLLABLE PWAEH
	{0xD440, 0xD440, prLV},                     // Lo       HANGUL SYLLABLE POE
	{0xD441, 0xD45B, prLVT},                    // Lo  [27] HANGUL SYLLABLE POEG..HANGUL SYLLABLE POEH
	{0xD45C, 0xD45C, prLV},                     // Lo       HANGUL SYLLABLE PYO
	{0xD45D, 0xD477, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYOG..HANGUL SYLLABLE PYOH
	{0xD478, 0xD478, prLV},                     // Lo       HANGUL SYLLABLE PU
	{0xD479, 0xD493, prLVT},                    // Lo  [27] HANGUL SYLLABLE PUG..HANGUL SYLLABLE PUH
	{0xD494, 0xD494, prLV},                     // Lo       HANGUL SYLLABLE PWEO
	{0xD495, 0xD4AF, prLVT},                    // Lo  [27] HANGUL SYLLABLE PWEOG..HANGUL SYLLABLE PWEOH
	{0xD4B0, 0xD4B0, prLV},                     // Lo       HANGUL SYLLABLE PWE
	{0xD4B1, 0xD4CB, prLVT},                    // Lo  [27] HANGUL SYLLABLE PWEG..HANGUL SYLLABLE PWEH
	{0xD4CC, 0xD4CC, prLV},                     // Lo       HANGUL SYLLABLE PWI
	{0xD4CD, 0xD4E7, prLVT},                    // Lo  [27] HANGUL SYLLABLE PWIG..HANGUL SYLLABLE PWIH
	{0xD4E8, 0xD4E8, prLV},                     // Lo       HANGUL SYLLABLE PYU
	{0xD4E9, 0xD503, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYUG..HANGUL SYLLABLE PYUH
	{0xD504, 0xD504, prLV},                     // Lo       HANGUL SYLLABLE PEU
	{0xD505, 0xD51F, prLVT},                    // Lo  [27] HANGUL SYLLABLE PEUG..HANGUL SYLLABLE PEUH
	{0xD520, 0xD520, prLV},                     // Lo       HANGUL SYLLABLE PYI
	{0xD521, 0xD53B, prLVT},                    // Lo  [27] HANGUL SYLLABLE PYIG..HANGUL SYLLABLE PYIH
	{0xD53C, 0xD53C, prLV},                     // Lo       HANGUL SYLLABLE PI
	{0xD53D, 0xD557, prLVT},                    // Lo  [27] HANGUL SYLLABLE PIG..HANGUL SYLLABLE PIH
	{0xD558, 0xD558, prLV},                     // Lo       HANGUL SYLLABLE HA
	{0xD559, 0xD573, prLVT},                    // Lo  [27] HANGUL SYLLABLE HAG..HANGUL SYLLABLE HAH
	{0xD574, 0xD574, prLV},                     // Lo       HANGUL SYLLABLE HAE
	{0xD575, 0xD58F, prLVT},                    // Lo  [27] HANGUL SYLLABLE HAEG..HANGUL SYLLABLE HAEH
	{0xD590, 0xD590, prLV},                     // Lo       HANGUL SYLLABLE HYA
	{0xD591, 0xD5AB, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYAG..HANGUL SYLLABLE HYAH
	{0xD5AC, 0xD5AC, prLV},                     // Lo       HANGUL SYLLABLE HYAE
	{0xD5AD, 0xD5C7, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYAEG..HANGUL SYLLABLE HYAEH
	{0xD5C8, 0xD5C8, prLV},                     // Lo       HANGUL SYLLABLE HEO
	{0xD5C9, 0xD5E3, prLVT},                    // Lo  [27] HANGUL SYLLABLE HEOG..HANGUL SYLLABLE HEOH
	{0xD5E4, 0xD5E4, prLV},                     // Lo       HANGUL SYLLABLE HE
	{0xD5E5, 0xD5FF, prLVT},                    // Lo  [27] HANGUL SYLLABLE HEG..HANGUL SYLLABLE HEH
	{0xD600, 0xD600, prLV},                     // Lo       HANGUL SYLLABLE HYEO
	{0xD601, 0xD61B, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYEOG..HANGUL SYLLABLE HYEOH
	{0xD61C, 0xD61C, prLV},                     // Lo       HANGUL SYLLABLE HYE
	{0xD61D, 0xD637, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYEG..HANGUL SYLLABLE HYEH
	{0xD638, 0xD638, prLV},                     // Lo       HANGUL SYLLABLE HO
	{0xD639, 0xD653, prLVT},                    // Lo  [27] HANGUL SYLLABLE HOG..HANGUL SYLLABLE HOH
	{0xD654, 0xD654, prLV},                     // Lo       HANGUL SYLLABLE HWA
	{0xD655, 0xD66F, prLVT},                    // Lo  [27] HANGUL SYLLABLE HWAG..HANGUL SYLLABLE HWAH
	{0xD670, 0xD670, prLV},                     // Lo       HANGUL SYLLABLE HWAE
	{0xD671, 0xD68B, prLVT},                    // Lo  [27] HANGUL SYLLABLE HWAEG..HANGUL SYLLABLE HWAEH
	{0xD68C, 0xD68C, prLV},                     // Lo       HANGUL SYLLABLE HOE
	{0xD68D, 0xD6A7, prLVT},                    // Lo  [27] HANGUL SYLLABLE HOEG..HANGUL SYLLABLE HOEH
	{0xD6A8, 0xD6A8, prLV},                     // Lo       HANGUL SYLLABLE HYO
	{0xD6A9, 0xD6C3, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYOG..HANGUL SYLLABLE HYOH
	{0xD6C4, 0xD6C4, prLV},                     // Lo       HANGUL SYLLABLE HU
	{0xD6C5, 0xD6DF, prLVT},                    // Lo  [27] HANGUL SYLLABLE HUG..HANGUL SYLLABLE HUH
	{0xD6E0, 0xD6E0, prLV},                     // Lo       HANGUL SYLLABLE HWEO
	{0xD6E1, 0xD6FB, prLVT},                    // Lo  [27] HANGUL SYLLABLE HWEOG..HANGUL SYLLABLE HWEOH
	{0xD6FC, 0xD6FC, prLV},                     // Lo       HANGUL SYLLABLE HWE
	{0xD6FD, 0xD717, prLVT},                    // Lo  [27] HANGUL SYLLABLE HWEG..HANGUL SYLLABLE HWEH
	{0xD718, 0xD718, prLV},                     // Lo       HANGUL SYLLABLE HWI
	{0xD719, 0xD733, prLVT},                    // Lo  [27] HANGUL SYLLABLE HWIG..HANGUL SYLLABLE HWIH
	{0xD734, 0xD734, prLV},                     // Lo       HANGUL SYLLABLE HYU
	{0xD735, 0xD74F, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYUG..HANGUL SYLLABLE HYUH
	{0xD750, 0xD750, prLV},                     // Lo       HANGUL SYLLABLE HEU
	{0xD751, 0xD76B, prLVT},                    // Lo  [27] HANGUL SYLLABLE HEUG..HANGUL SYLLABLE HEUH
	{0xD76C, 0xD76C, prLV},                     // Lo       HANGUL SYLLABLE HYI
	{0xD76D, 0xD787, prLVT},                    // Lo  [27] HANGUL SYLLABLE HYIG..HANGUL SYLLABLE HYIH
	{0xD788, 0xD788, prLV},                     // Lo       HANGUL SYLLABLE HI
	{0xD789, 0xD7A3, prLVT},                    // Lo  [27] HANGUL SYLLABLE HIG..HANGUL SYLLABLE HIH
	{0xD7B0, 0xD7C6, prV},                      // Lo  [23] HANGUL JUNGSEONG O-YEO..HANGUL JUNGSEONG ARAEA-E
	{0xD7CB, 0xD7FB, prT},                      // Lo  [49] HANGUL JONGSEONG NIEUN-RIEUL..HANGUL JONGSEONG PHIEUPH-THIEUTH
	{0xFB1E, 0xFB1E, prExtend},                 // Mn       HEBREW POINT JUDEO-SPANISH VARIKA
	{0xFE00, 0xFE0F, prExtend},                 // Mn  [16] VARIATION SELECTOR-1..VARIATION SELECTOR-16
	{0xFE20, 0xFE2F, prExtend},                 // Mn  [16] COMBINING LIGATURE LEFT HALF..COMBINING CYRILLIC TITLO RIGHT HALF
	{0xFEFF, 0xFEFF, prControl},                // Cf       ZERO WIDTH NO-BREAK SPACE
	{0xFF9E, 0xFF9F, prExtend},                 // Lm   [2] HALFWIDTH KATAKANA VOICED SOUND MARK..HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK
	{0xFFF0, 0xFFF8, prControl},                // Cn   [9] <reserved-FFF0>..<reserved-FFF8>
	{0xFFF9, 0xFFFB, prControl},                // Cf   [3] INTERLINEAR ANNOTATION ANCHOR..INTERLINEAR ANNOTATION TERMINATOR
	{0x101FD, 0x101FD, prExtend},               // Mn       PHAISTOS DISC SIGN COMBINING OBLIQUE STROKE
	{0x102E0, 0x102E0, prExtend},               // Mn       COPTIC EPACT THOUSANDS MARK
	{0x10376, 0x1037A, prExtend},               // Mn   [5] COMBINING OLD PERMIC LETTER AN..COMBINING OLD PERMIC LETTER SII
	{0x10A01, 0x10A03, prExtend},               // Mn   [3] KHAROSHTHI VOWEL SIGN I..KHAROSHTHI VOWEL SIGN VOCALIC R
	{0x10A05, 0x10A06, prExtend},               // Mn   [2] KHAROSHTHI VOWEL SIGN E..KHAROSHTHI VOWEL SIGN O
	{0x10A0C, 0x10A0F, prExtend},               // Mn   [4] KHAROSHTHI VOWEL LENGTH MARK..KHAROSHTHI SIGN VISARGA
	{0x10A38, 0x10A3A, prExtend},               // Mn   [3] KHAROSHTHI SIGN BAR ABOVE..KHAROSHTHI SIGN DOT BELOW
	{0x10A3F, 0x10A3F, prExtend},               // Mn       KHAROSHTHI VIRAMA
	{0x10AE5, 0x10AE6, prExtend},               // Mn   [2] MANICHAEAN ABBREVIATION MARK ABOVE..MANICHAEAN ABBREVIATION MARK BELOW
	{0x10D24, 0x10D27, prExtend},               // Mn   [4] HANIFI ROHINGYA SIGN HARBAHAY..HANIFI ROHINGYA SIGN TASSI
	{0x10EAB, 0x10EAC, prExtend},               // Mn   [2] YEZIDI COMBINING HAMZA MARK..YEZIDI COMBINING MADDA MARK
	{0x10EFD, 0x10EFF, prExtend},               // Mn   [3] ARABIC SMALL LOW WORD SAKTA..ARABIC SMALL LOW WORD MADDA
	{0x10F46, 0x10F50, prExtend},               // Mn  [11] SOGDIAN COMBINING DOT BELOW..SOGDIAN COMBINING STROKE BELOW
	{0x10F82, 0x10F85, prExtend},               // Mn   [4] OLD UYGHUR COMBINING DOT ABOVE..OLD UYGHUR COMBINING TWO DOTS BELOW
	{0x11000, 0x11000, prSpacingMark},          // Mc       BRAHMI SIGN CANDRABINDU
	{0x11001, 0x11001, prExtend},               // Mn       BRAHMI SIGN ANUSVARA
	{0x11002, 0x11002, prSpacingMark},          // Mc       BRAHMI SIGN VISARGA
	{0x11038, 0x11046, prExtend},               // Mn  [15] BRAHMI VOWEL SIGN AA..BRAHMI VIRAMA
	{0x11070, 0x11070, prExtend},               // Mn       BRAHMI SIGN OLD TAMIL VIRAMA
	{0x11073, 0x11074, prExtend},               // Mn   [2] BRAHMI VOWEL SIGN OLD TAMIL SHORT E..BRAHMI VOWEL SIGN OLD TAMIL SHORT O
	{0x1107F, 0x11081, prExtend},               // Mn   [3] BRAHMI NUMBER JOINER..KAITHI SIGN ANUSVARA
	{0x11082, 0x11082, prSpacingMark},          // Mc       KAITHI SIGN VISARGA
	{0x110B0, 0x110B2, prSpacingMark},          // Mc   [3] KAITHI VOWEL SIGN AA..KAITHI VOWEL SIGN II
	{0x110B3, 0x110B6, prExtend},               // Mn   [4] KAITHI VOWEL SIGN U..KAITHI VOWEL SIGN AI
	{0x110B7, 0x110B8, prSpacingMark},          // Mc   [2] KAITHI VOWEL SIGN O..KAITHI VOWEL SIGN AU
	{0x110B9, 0x110BA, prExtend},               // Mn   [2] KAITHI SIGN VIRAMA..KAITHI SIGN NUKTA
	{0x110BD, 0x110BD, prPrepend},              // Cf       KAITHI NUMBER SIGN
	{0x110C2, 0x110C2, prExtend},               // Mn       KAITHI VOWEL SIGN VOCALIC R
	{0x110CD, 0x110CD, prPrepend},              // Cf       KAITHI NUMBER SIGN ABOVE
	{0x11100, 0x11102, prExtend},               // Mn   [3] CHAKMA SIGN CANDRABINDU..CHAKMA SIGN VISARGA
	{0x11127, 0x1112B, prExtend},               // Mn   [5] CHAKMA VOWEL SIGN A..CHAKMA VOWEL SIGN UU
	{0x1112C, 0x1112C, prSpacingMark},          // Mc       CHAKMA VOWEL SIGN E
	{0x1112D, 0x11134, prExtend},               // Mn   [8] CHAKMA VOWEL SIGN AI..CHAKMA MAAYYAA
	{0x11145, 0x11146, prSpacingMark},          // Mc   [2] CHAKMA VOWEL SIGN AA..CHAKMA VOWEL SIGN EI
	{0x11173, 0x11173, prExtend},               // Mn       MAHAJANI SIGN NUKTA
	{0x11180, 0x11181, prExtend},               // Mn   [2] SHARADA SIGN CANDRABINDU..SHARADA SIGN ANUSVARA
	{0x11182, 0x11182, prSpacingMark},          // Mc       SHARADA SIGN VISARGA
	{0x111B3, 0x111B5, prSpacingMark},          // Mc   [3] SHARADA VOWEL SIGN AA..SHARADA VOWEL SIGN II
	{0x111B6, 0x111BE, prExtend},               // Mn   [9] SHARADA VOWEL SIGN U..SHARADA VOWEL SIGN O
	{0x111BF, 0x111C0, prSpacingMark},          // Mc   [2] SHARADA VOWEL SIGN AU..SHARADA SIGN VIRAMA
	{0x111C2, 0x111C3, prPrepend},              // Lo   [2] SHARADA SIGN JIHVAMULIYA..SHARADA SIGN UPADHMANIYA
	{0x111C9, 0x111CC, prExtend},               // Mn   [4] SHARADA SANDHI MARK..SHARADA EXTRA SHORT VOWEL MARK
	{0x111CE, 0x111CE, prSpacingMark},          // Mc       SHARADA VOWEL SIGN PRISHTHAMATRA E
	{0x111CF, 0x111CF, prExtend},               // Mn       SHARADA SIGN INVERTED CANDRABINDU
	{0x1122C, 0x1122E, prSpacingMark},          // Mc   [3] KHOJKI VOWEL SIGN AA..KHOJKI VOWEL SIGN II
	{0x1122F, 0x11231, prExtend},               // Mn   [3] KHOJKI VOWEL SIGN U..KHOJKI VOWEL SIGN AI
	{0x11232, 0x11233, prSpacingMark},          // Mc   [2] KHOJKI VOWEL SIGN O..KHOJKI VOWEL SIGN AU
	{0x11234, 0x11234, prExtend},               // Mn       KHOJKI SIGN ANUSVARA
	{0x11235, 0x11235, prSpacingMark},          // Mc       KHOJKI SIGN VIRAMA
	{0x11236, 0x11237, prExtend},               // Mn   [2] KHOJKI SIGN NUKTA..KHOJKI SIGN SHADDA
	{0x1123E, 0x1123E, prExtend},               // Mn       KHOJKI SIGN SUKUN
	{0x11241, 0x11241, prExtend},               // Mn       KHOJKI VOWEL SIGN VOCALIC R
	{0x112DF, 0x112DF, prExtend},               // Mn       KHUDAWADI SIGN ANUSVARA
	{0x112E0, 0x112E2, prSpacingMark},          // Mc   [3] KHUDAWADI VOWEL SIGN AA..KHUDAWADI VOWEL SIGN II
	{0x112E3, 0x112EA, prExtend},               // Mn   [8] KHUDAWADI VOWEL SIGN U..KHUDAWADI SIGN VIRAMA
	{0x11300, 0x11301, prExtend},               // Mn   [2] GRANTHA SIGN COMBINING ANUSVARA ABOVE..GRANTHA SIGN CANDRABINDU
	{0x11302, 0x11303, prSpacingMark},          // Mc   [2] GRANTHA SIGN ANUSVARA..GRANTHA SIGN VISARGA
	{0x1133B, 0x1133C, prExtend},               // Mn   [2] COMBINING BINDU BELOW..GRANTHA SIGN NUKTA
	{0x1133E, 0x1133E, prExtend},               // Mc       GRANTHA VOWEL SIGN AA
	{0x1133F, 0x1133F, prSpacingMark},          // Mc       GRANTHA VOWEL SIGN I
	{0x11340, 0x11340, prExtend},               // Mn       GRANTHA VOWEL SIGN II
	{0x11341, 0x11344, prSpacingMark},          // Mc   [4] GRANTHA VOWEL SIGN U..GRANTHA VOWEL SIGN VOCALIC RR
	{0x11347, 0x11348, prSpacingMark},          // Mc   [2] GRANTHA VOWEL SIGN EE..GRANTHA VOWEL SIGN AI
	{0x1134B, 0x1134D, prSpacingMark},          // Mc   [3] GRANTHA VOWEL SIGN OO..GRANTHA SIGN VIRAMA
	{0x11357, 0x11357, prExtend},               // Mc       GRANTHA AU LENGTH MARK
	{0x11362, 0x11363, prSpacingMark},          // Mc   [2] GRANTHA VOWEL SIGN VOCALIC L..GRANTHA VOWEL SIGN VOCALIC LL
	{0x11366, 0x1136C, prExtend},               // Mn   [7] COMBINING GRANTHA DIGIT ZERO..COMBINING GRANTHA DIGIT SIX
	{0x11370, 0x11374, prExtend},               // Mn   [5] COMBINING GRANTHA LETTER A..COMBINING GRANTHA LETTER PA
	{0x11435, 0x11437, prSpacingMark},          // Mc   [3] NEWA VOWEL SIGN AA..NEWA VOWEL SIGN II
	{0x11438, 0x1143F, prExtend},               // Mn   [8] NEWA VOWEL SIGN U..NEWA VOWEL SIGN AI
	{0x11440, 0x11441, prSpacingMark},          // Mc   [2] NEWA VOWEL SIGN O..NEWA VOWEL SIGN AU
	{0x11442, 0x11444, prExtend},               // Mn   [3] NEWA SIGN VIRAMA..NEWA SIGN ANUSVARA
	{0x11445, 0x11445, prSpacingMark},          // Mc       NEWA SIGN VISARGA
	{0x11446, 0x11446, prExtend},               // Mn       NEWA SIGN NUKTA
	{0x1145E, 0x1145E, prExtend},               // Mn       NEWA SANDHI MARK
	{0x114B0, 0x114B0, prExtend},               // Mc       TIRHUTA VOWEL SIGN AA
	{0x114B1, 0x114B2, prSpacingMark},          // Mc   [2] TIRHUTA VOWEL SIGN I..TIRHUTA VOWEL SIGN II
	{0x114B3, 0x114B8, prExtend},               // Mn   [6] TIRHUTA VOWEL SIGN U..TIRHUTA VOWEL SIGN VOCALIC LL
	{0x114B9, 0x114B9, prSpacingMark},          // Mc       TIRHUTA VOWEL SIGN E
	{0x114BA, 0x114BA, prExtend},               // Mn       TIRHUTA VOWEL SIGN SHORT E
	{0x114BB, 0x114BC, prSpacingMark},          // Mc   [2] TIRHUTA VOWEL SIGN AI..TIRHUTA VOWEL SIGN O
	{0x114BD, 0x114BD, prExtend},               // Mc       TIRHUTA VOWEL SIGN SHORT O
	{0x114BE, 0x114BE, prSpacingMark},          // Mc       TIRHUTA VOWEL SIGN AU
	{0x114BF, 0x114C0, prExtend},               // Mn   [2] TIRHUTA SIGN CANDRABINDU..TIRHUTA SIGN ANUSVARA
	{0x114C1, 0x114C1, prSpacingMark},          // Mc       TIRHUTA SIGN VISARGA
	{0x114C2, 0x114C3, prExtend},               // Mn   [2] TIRHUTA SIGN VIRAMA..TIRHUTA SIGN NUKTA
	{0x115AF, 0x115AF, prExtend},               // Mc       SIDDHAM VOWEL SIGN AA
	{0x115B0, 0x115B1, prSpacingMark},          // Mc   [2] SIDDHAM VOWEL SIGN I..SIDDHAM VOWEL SIGN II
	{0x115B2, 0x115B5, prExtend},               // Mn   [4] SIDDHAM VOWEL SIGN U..SIDDHAM VOWEL SIGN VOCALIC RR
	{0x115B8, 0x115BB, prSpacingMark},          // Mc   [4] SIDDHAM VOWEL SIGN E..SIDDHAM VOWEL SIGN AU
	{0x115BC, 0x115BD, prExtend},               // Mn   [2] SIDDHAM SIGN CANDRABINDU..SIDDHAM SIGN ANUSVARA
	{0x115BE, 0x115BE, prSpacingMark},          // Mc       SIDDHAM SIGN VISARGA
	{0x115BF, 0x115C0, prExtend},               // Mn   [2] SIDDHAM SIGN VIRAMA..SIDDHAM SIGN NUKTA
	{0x115DC, 0x115DD, prExtend},               // Mn   [2] SIDDHAM VOWEL SIGN ALTERNATE U..SIDDHAM VOWEL SIGN ALTERNATE UU
	{0x11630, 0x11632, prSpacingMark},          // Mc   [3] MODI VOWEL SIGN AA..MODI VOWEL SIGN II
	{0x11633, 0x1163A, prExtend},               // Mn   [8] MODI VOWEL SIGN U..MODI VOWEL SIGN AI
	{0x1163B, 0x1163C, prSpacingMark},          // Mc   [2] MODI VOWEL SIGN O..MODI VOWEL SIGN AU
	{0x1163D, 0x1163D, prExtend},               // Mn       MODI SIGN ANUSVARA
	{0x1163E, 0x1163E, prSpacingMark},          // Mc       MODI SIGN VISARGA
	{0x1163F, 0x11640, prExtend},               // Mn   [2] MODI SIGN VIRAMA..MODI SIGN ARDHACANDRA
	{0x116AB, 0x116AB, prExtend},               // Mn       TAKRI SIGN ANUSVARA
	{0x116AC, 0x116AC, prSpacingMark},          // Mc       TAKRI SIGN VISARGA
	{0x116AD, 0x116AD, prExtend},               // Mn       TAKRI VOWEL SIGN AA
	{0x116AE, 0x116AF, prSpacingMark},          // Mc   [2] TAKRI VOWEL SIGN I..TAKRI VOWEL SIGN II
	{0x116B0, 0x116B5, prExtend},               // Mn   [6] TAKRI VOWEL SIGN U..TAKRI VOWEL SIGN AU
	{0x116B6, 0x116B6, prSpacingMark},          // Mc       TAKRI SIGN VIRAMA
	{0x116B7, 0x116B7, prExtend},               // Mn       TAKRI SIGN NUKTA
	{0x1171D, 0x1171F, prExtend},               // Mn   [3] AHOM CONSONANT SIGN MEDIAL LA..AHOM CONSONANT SIGN MEDIAL LIGATING RA
	{0x11722, 0x11725, prExtend},               // Mn   [4] AHOM VOWEL SIGN I..AHOM VOWEL SIGN UU
	{0x11726, 0x11726, prSpacingMark},          // Mc       AHOM VOWEL SIGN E
	{0x11727, 0x1172B, prExtend},               // Mn   [5] AHOM VOWEL SIGN AW..AHOM SIGN KILLER
	{0x1182C, 0x1182E, prSpacingMark},          // Mc   [3] DOGRA VOWEL SIGN AA..DOGRA VOWEL SIGN II
	{0x1182F, 0x11837, prExtend},               // Mn   [9] DOGRA VOWEL SIGN U..DOGRA SIGN ANUSVARA
	{0x11838, 0x11838, prSpacingMark},          // Mc       DOGRA SIGN VISARGA
	{0x11839, 0x1183A, prExtend},               // Mn   [2] DOGRA SIGN VIRAMA..DOGRA SIGN NUKTA
	{0x11930, 0x11930, prExtend},               // Mc       DIVES AKURU VOWEL SIGN AA
	{0x11931, 0x11935, prSpacingMark},          // Mc   [5] DIVES AKURU VOWEL SIGN I..DIVES AKURU VOWEL SIGN E
	{0x11937, 0x11938, prSpacingMark},          // Mc   [2] DIVES AKURU VOWEL SIGN AI..DIVES AKURU VOWEL SIGN O
	{0x1193B, 0x1193C, prExtend},               // Mn   [2] DIVES AKURU SIGN ANUSVARA..DIVES AKURU SIGN CANDRABINDU
	{0x1193D, 0x1193D, prSpacingMark},          // Mc       DIVES AKURU SIGN HALANTA
	{0x1193E, 0x1193E, prExtend},               // Mn       DIVES AKURU VIRAMA
	{0x1193F, 0x1193F, prPrepend},              // Lo       DIVES AKURU PREFIXED NASAL SIGN
	{0x11940, 0x11940, prSpacingMark},          // Mc       DIVES AKURU MEDIAL YA
	{0x11941, 0x11941, prPrepend},              // Lo       DIVES AKURU INITIAL RA
	{0x11942, 0x11942, prSpacingMark},          // Mc       DIVES AKURU MEDIAL RA
	{0x11943, 0x11943, prExtend},               // Mn       DIVES AKURU SIGN NUKTA
	{0x119D1, 0x119D3, prSpacingMark},          // Mc   [3] NANDINAGARI VOWEL SIGN AA..NANDINAGARI VOWEL SIGN II
	{0x119D4, 0x119D7, prExtend},               // Mn   [4] NANDINAGARI VOWEL SIGN U..NANDINAGARI VOWEL SIGN VOCALIC RR
	{0x119DA, 0x119DB, prExtend},               // Mn   [2] NANDINAGARI VOWEL SIGN E..NANDINAGARI VOWEL SIGN AI
	{0x119DC, 0x119DF, prSpacingMark},          // Mc   [4] NANDINAGARI VOWEL SIGN O..NANDINAGARI SIGN VISARGA
	{0x119E0, 0x119E0, prExtend},               // Mn       NANDINAGARI SIGN VIRAMA
	{0x119E4, 0x119E4, prSpacingMark},          // Mc       NANDINAGARI VOWEL SIGN PRISHTHAMATRA E
	{0x11A01, 0x11A0A, prExtend},               // Mn  [10] ZANABAZAR SQUARE VOWEL SIGN I..ZANABAZAR SQUARE VOWEL LENGTH MARK
	{0x11A33, 0x11A38, prExtend},               // Mn   [6] ZANABAZAR SQUARE FINAL CONSONANT MARK..ZANABAZAR SQUARE SIGN ANUSVARA
	{0x11A39, 0x11A39, prSpacingMark},          // Mc       ZANABAZAR SQUARE SIGN VISARGA
	{0x11A3A, 0x11A3A, prPrepend},              // Lo       ZANABAZAR SQUARE CLUSTER-INITIAL LETTER RA
	{0x11A3B, 0x11A3E, prExtend},               // Mn   [4] ZANABAZAR SQUARE CLUSTER-FINAL LETTER YA..ZANABAZAR SQUARE CLUSTER-FINAL LETTER VA
	{0x11A47, 0x11A47, prExtend},               // Mn       ZANABAZAR SQUARE SUBJOINER
	{0x11A51, 0x11A56, prExtend},               // Mn   [6] SOYOMBO VOWEL SIGN I..SOYOMBO VOWEL SIGN OE
	{0x11A57, 0x11A58, prSpacingMark},          // Mc   [2] SOYOMBO VOWEL SIGN AI..SOYOMBO VOWEL SIGN AU
	{0x11A59, 0x11A5B, prExtend},               // Mn   [3] SOYOMBO VOWEL SIGN VOCALIC R..SOYOMBO VOWEL LENGTH MARK
	{0x11A84, 0x11A89, prPrepend},              // Lo   [6] SOYOMBO SIGN JIHVAMULIYA..SOYOMBO CLUSTER-INITIAL LETTER SA
	{0x11A8A, 0x11A96, prExtend},               // Mn  [13] SOYOMBO FINAL CONSONANT SIGN G..SOYOMBO SIGN ANUSVARA
	{0x11A97, 0x11A97, prSpacingMark},          // Mc       SOYOMBO SIGN VISARGA
	{0x11A98, 0x11A99, prExtend},               // Mn   [2] SOYOMBO GEMINATION MARK..SOYOMBO SUBJOINER
	{0x11C2F, 0x11C2F, prSpacingMark},          // Mc       BHAIKSUKI VOWEL SIGN AA
	{0x11C30, 0x11C36, prExtend},               // Mn   [7] BHAIKSUKI VOWEL SIGN I..BHAIKSUKI VOWEL SIGN VOCALIC L
	{0x11C38, 0x11C3D, prExtend},               // Mn   [6] BHAIKSUKI VOWEL SIGN E..BHAIKSUKI SIGN ANUSVARA
	{0x11C3E, 0x11C3E, prSpacingMark},          // Mc       BHAIKSUKI SIGN VISARGA
	{0x11C3F, 0x11C3F, prExtend},               // Mn       BHAIKSUKI SIGN VIRAMA
	{0x11C92, 0x11CA7, prExtend},               // Mn  [22] MARCHEN SUBJOINED LETTER KA..MARCHEN SUBJOINED LETTER ZA
	{0x11CA9, 0x11CA9, prSpacingMark},          // Mc       MARCHEN SUBJOINED LETTER YA
	{0x11CAA, 0x11CB0, prExtend},               // Mn   [7] MARCHEN SUBJOINED LETTER RA..MARCHEN VOWEL SIGN AA
	{0x11CB1, 0x11CB1, prSpacingMark},          // Mc       MARCHEN VOWEL SIGN I
	{0x11CB2, 0x11CB3, prExtend},               // Mn   [2] MARCHEN VOWEL SIGN U..MARCHEN VOWEL SIGN E
	{0x11CB4, 0x11CB4, prSpacingMark},          // Mc       MARCHEN VOWEL SIGN O
	{0x11CB5, 0x11CB6, prExtend},               // Mn   [2] MARCHEN SIGN ANUSVARA..MARCHEN SIGN CANDRABINDU
	{0x11D31, 0x11D36, prExtend},               // Mn   [6] MASARAM GONDI VOWEL SIGN AA..MASARAM GONDI VOWEL SIGN VOCALIC R
	{0x11D3A, 0x11D3A, prExtend},               // Mn       MASARAM GONDI VOWEL SIGN E
	{0x11D3C, 0x11D3D, prExtend},               // Mn   [2] MASARAM GONDI VOWEL SIGN AI..MASARAM GONDI VOWEL SIGN O
	{0x11D3F, 0x11D45, prExtend},               // Mn   [7] MASARAM GONDI VOWEL SIGN AU..MASARAM GONDI VIRAMA
	{0x11D46, 0x11D46, prPrepend},              // Lo       MASARAM GONDI REPHA
	{0x11D47, 0x11D47, prExtend},               // Mn       MASARAM GONDI RA-KARA
	{0x11D8A, 0x11D8E, prSpacingMark},          // Mc   [5] GUNJALA GONDI VOWEL SIGN AA..GUNJALA GONDI VOWEL SIGN UU
	{0x11D90, 0x11D91, prExtend},               // Mn   [2] GUNJALA GONDI VOWEL SIGN EE..GUNJALA GONDI VOWEL SIGN AI
	{0x11D93, 0x11D94, prSpacingMark},          // Mc   [2] GUNJALA GONDI VOWEL SIGN OO..GUNJALA GONDI VOWEL SIGN AU
	{0x11D95, 0x11D95, prExtend},               // Mn       GUNJALA GONDI SIGN ANUSVARA
	{0x11D96, 0x11D96, prSpacingMark},          // Mc       GUNJALA GONDI SIGN VISARGA
	{0x11D97, 0x11D97, prExtend},               // Mn       GUNJALA GONDI VIRAMA
	{0x11EF3, 0x11EF4, prExtend},               // Mn   [2] MAKASAR VOWEL SIGN I..MAKASAR VOWEL SIGN U
	{0x11EF5, 0x11EF6, prSpacingMark},          // Mc   [2] MAKASAR VOWEL SIGN E..MAKASAR VOWEL SIGN O
	{0x11F00, 0x11F01, prExtend},               // Mn   [2] KAWI SIGN CANDRABINDU..KAWI SIGN ANUSVARA
	{0x11F02, 0x11F02, prPrepend},              // Lo       KAWI SIGN REPHA
	{0x11F03, 0x11F03, prSpacingMark},          // Mc       KAWI SIGN VISARGA
	{0x11F34, 0x11F35, prSpacingMark},          // Mc   [2] KAWI VOWEL SIGN AA..KAWI VOWEL SIGN ALTERNATE AA
	{0x11F36, 0x11F3A, prExtend},               // Mn   [5] KAWI VOWEL SIGN I..KAWI VOWEL SIGN VOCALIC R
	{0x11F3E, 0x11F3F, prSpacingMark},          // Mc   [2] KAWI VOWEL SIGN E..KAWI VOWEL SIGN AI
	{0x11F40, 0x11F40, prExtend},               // Mn       KAWI VOWEL SIGN EU
	{0x11F41, 0x11F41, prSpacingMark},          // Mc       KAWI SIGN KILLER
	{0x11F42, 0x11F42, prExtend},               // Mn       KAWI CONJOINER
	{0x13430, 0x1343F, prControl},              // Cf  [16] EGYPTIAN HIEROGLYPH VERTICAL JOINER..EGYPTIAN HIEROGLYPH END WALLED ENCLOSURE
	{0x13440, 0x13440, prExtend},               // Mn       EGYPTIAN HIEROGLYPH MIRROR HORIZONTALLY
	{0x13447, 0x13455, prExtend},               // Mn  [15] EGYPTIAN HIEROGLYPH MODIFIER DAMAGED AT TOP START..EGYPTIAN HIEROGLYPH MODIFIER DAMAGED
	{0x16AF0, 0x16AF4, prExtend},               // Mn   [5] BASSA VAH COMBINING HIGH TONE..BASSA VAH COMBINING HIGH-LOW TONE
	{0x16B30, 0x16B36, prExtend},               // Mn   [7] PAHAWH HMONG MARK CIM TUB..PAHAWH HMONG MARK CIM TAUM
	{0x16F4F, 0x16F4F, prExtend},               // Mn       MIAO SIGN CONSONANT MODIFIER BAR
	{0x16F51, 0x16F87, prSpacingMark},          // Mc  [55] MIAO SIGN ASPIRATION..MIAO VOWEL SIGN UI
	{0x16F8F, 0x16F92, prExtend},               // Mn   [4] MIAO TONE RIGHT..MIAO TONE BELOW
	{0x16FE4, 0x16FE4, prExtend},               // Mn       KHITAN SMALL SCRIPT FILLER
	{0x16FF0, 0x16FF1, prSpacingMark},          // Mc   [2] VIETNAMESE ALTERNATE READING MARK CA..VIETNAMESE ALTERNATE READING MARK NHAY
	{0x1BC9D, 0x1BC9E, prExtend},               // Mn   [2] DUPLOYAN THICK LETTER SELECTOR..DUPLOYAN DOUBLE MARK
	{0x1BCA0, 0x1BCA3, prControl},              // Cf   [4] SHORTHAND FORMAT LETTER OVERLAP..SHORTHAND FORMAT UP STEP
	{0x1CF00, 0x1CF2D, prExtend},               // Mn  [46] ZNAMENNY COMBINING MARK GORAZDO NIZKO S KRYZHEM ON LEFT..ZNAMENNY COMBINING MARK KRYZH ON LEFT
	{0x1CF30, 0x1CF46, prExtend},               // Mn  [23] ZNAMENNY COMBINING TONAL RANGE MARK MRACHNO..ZNAMENNY PRIZNAK MODIFIER ROG
	{0x1D165, 0x1D165, prExtend},               // Mc       MUSICAL SYMBOL COMBINING STEM
	{0x1D166, 0x1D166, prSpacingMark},          // Mc       MUSICAL SYMBOL COMBINING SPRECHGESANG STEM
	{0x1D167, 0x1D169, prExtend},               // Mn   [3] MUSICAL SYMBOL COMBINING TREMOLO-1..MUSICAL SYMBOL COMBINING TREMOLO-3
	{0x1D16D, 0x1D16D, prSpacingMark},          // Mc       MUSICAL SYMBOL COMBINING AUGMENTATION DOT
	{0x1D16E, 0x1D172, prExtend},               // Mc   [5] MUSICAL SYMBOL COMBINING FLAG-1..MUSICAL SYMBOL COMBINING FLAG-5
	{0x1D173, 0x1D17A, prControl},              // Cf   [8] MUSICAL SYMBOL BEGIN BEAM..MUSICAL SYMBOL END PHRASE
	{0x1D17B, 0x1D182, prExtend},               // Mn   [8] MUSICAL SYMBOL COMBINING ACCENT..MUSICAL SYMBOL COMBINING LOURE
	{0x1D185, 0x1D18B, prExtend},               // Mn   [7] MUSICAL SYMBOL COMBINING DOIT..MUSICAL SYMBOL COMBINING TRIPLE TONGUE
	{0x1D1AA, 0x1D1AD, prExtend},               // Mn   [4] MUSICAL SYMBOL COMBINING DOWN BOW..MUSICAL SYMBOL COMBINING SNAP PIZZICATO
	{0x1D242, 0x1D244, prExtend},               // Mn   [3] COMBINING GREEK MUSICAL TRISEME..COMBINING GREEK MUSICAL PENTASEME
	{0x1DA00, 0x1DA36, prExtend},               // Mn  [55] SIGNWRITING HEAD RIM..SIGNWRITING AIR SUCKING IN
	{0x1DA3B, 0x1DA6C, prExtend},               // Mn  [50] SIGNWRITING MOUTH CLOSED NEUTRAL..SIGNWRITING EXCITEMENT
	{0x1DA75, 0x1DA75, prExtend},               // Mn       SIGNWRITING UPPER BODY TILTING FROM HIP JOINTS
	{0x1DA84, 0x1DA84, prExtend},               // Mn       SIGNWRITING LOCATION HEAD NECK
	{0x1DA9B, 0x1DA9F, prExtend},               // Mn   [5] SIGNWRITING FILL MODIFIER-2..SIGNWRITING FILL MODIFIER-6
	{0x1DAA1, 0x1DAAF, prExtend},               // Mn  [15] SIGNWRITING ROTATION MODIFIER-2..SIGNWRITING ROTATION MODIFIER-16
	{0x1E000, 0x1E006, prExtend},               // Mn   [7] COMBINING GLAGOLITIC LETTER AZU..COMBINING GLAGOLITIC LETTER ZHIVETE
	{0x1E008, 0x1E018, prExtend},               // Mn  [17] COMBINING GLAGOLITIC LETTER ZEMLJA..COMBINING GLAGOLITIC LETTER HERU
	{0x1E01B, 0x1E021, prExtend},               // Mn   [7] COMBINING GLAGOLITIC LETTER SHTA..COMBINING GLAGOLITIC LETTER YATI
	{0x1E023, 0x1E024, prExtend},               // Mn   [2] COMBINING GLAGOLITIC LETTER YU..COMBINING GLAGOLITIC LETTER SMALL YUS
	{0x1E026, 0x1E02A, prExtend},               // Mn   [5] COMBINING GLAGOLITIC LETTER YO..COMBINING GLAGOLITIC LETTER FITA
	{0x1E08F, 0x1E08F, prExtend},               // Mn       COMBINING CYRILLIC SMALL LETTER BYELORUSSIAN-UKRAINIAN I
	{0x1E130, 0x1E136, prExtend},               // Mn   [7] NYIAKENG PUACHUE HMONG TONE-B..NYIAKENG PUACHUE HMONG TONE-D
	{0x1E2AE, 0x1E2AE, prExtend},               // Mn       TOTO SIGN RISING TONE
	{0x1E2EC, 0x1E2EF, prExtend},               // Mn   [4] WANCHO TONE TUP..WANCHO TONE KOINI
	{0x1E4EC, 0x1E4EF, prExtend},               // Mn   [4] NAG MUNDARI SIGN MUHOR..NAG MUNDARI SIGN SUTUH
	{0x1E8D0, 0x1E8D6, prExtend},               // Mn   [7] MENDE KIKAKUI COMBINING NUMBER TEENS..MENDE KIKAKUI COMBINING NUMBER MILLIONS
	{0x1E944, 0x1E94A, prExtend},               // Mn   [7] ADLAM ALIF LENGTHENER..ADLAM NUKTA
	{0x1F000, 0x1F003, prExtendedPictographic}, // E0.0   [4] (🀀..🀃)    MAHJONG TILE EAST WIND..MAHJONG TILE NORTH WIND
	{0x1F004, 0x1F004, prExtendedPictographic}, // E0.6   [1] (🀄)       mahjong red dragon
	{0x1F005, 0x1F0CE, prExtendedPictographic}, // E0.0 [202] (🀅..🃎)    MAHJONG TILE GREEN DRAGON..PLAYING CARD KING OF DIAMONDS
	{0x1F0CF, 0x1F0CF, prExtendedPictographic}, // E0.6   [1] (🃏)       joker
	{0x1F0D0, 0x1F0FF, prExtendedPictographic}, // E0.0  [48] (🃐..🃿)    <reserved-1F0D0>..<reserved-1F0FF>
	{0x1F10D, 0x1F10F, prExtendedPictographic}, // E0.0   [3] (🄍..🄏)    CIRCLED ZERO WITH SLASH..CIRCLED DOLLAR SIGN WITH OVERLAID BACKSLASH
	{0x1F12F, 0x1F12F, prExtendedPictographic}, // E0.0   [1] (🄯)       COPYLEFT SYMBOL
	{0x1F16C, 0x1F16F, prExtendedPictographic}, // E0.0   [4] (🅬..🅯)    RAISED MR SIGN..CIRCLED HUMAN FIGURE
	{0x1F170, 0x1F171, prExtendedPictographic}, // E0.6   [2] (🅰️..🅱️)    A button (blood type)..B button (blood type)
	{0x1F17E, 0x1F17F, prExtendedPictographic}, // E0.6   [2] (🅾️..🅿️)    O button (blood type)..P button
	{0x1F18E, 0x1F18E, prExtendedPictographic}, // E0.6   [1] (🆎)       AB button (blood type)
	{0x1F191, 0x1F19A, prExtendedPictographic}, // E0.6  [10] (🆑..🆚)    CL button..VS button
	{0x1F1AD, 0x1F1E5, prExtendedPictographic}, // E0.0  [57] (🆭..🇥)    MASK WORK SYMBOL..<reserved-1F1E5>
	{0x1F1E6, 0x1F1FF, prRegionalIndicator},    // So  [26] REGIONAL INDICATOR SYMBOL LETTER A..REGIONAL INDICATOR SYMBOL LETTER Z
	{0x1F201, 0x1F202, prExtendedPictographic}, // E0.6   [2] (🈁..🈂️)    Japanese “here” button..Japanese “service charge” button
	{0x1F203, 0x1F20F, prExtendedPictographic}, // E0.0  [13] (🈃..🈏)    <reserved-1F203>..<reserved-1F20F>
	{0x1F21A, 0x1F21A, prExtendedPictographic}, // E0.6   [1] (🈚)       Japanese “free of charge” button
	{0x1F22F, 0x1F22F, prExtendedPictographic}, // E0.6   [1] (🈯)       Japanese “reserved” button
	{0x1F232, 0x1F23A, prExtendedPictographic}, // E0.6   [9] (🈲..🈺)    Japanese “prohibited” button..Japanese “open for business” button
	{0x1F23C, 0x1F23F, prExtendedPictographic}, // E0.0   [4] (🈼..🈿)    <reserved-1F23C>..<reserved-1F23F>
	{0x1F249, 0x1F24F, prExtendedPictographic}, // E0.0   [7] (🉉..🉏)    <reserved-1F249>..<reserved-1F24F>
	{0x1F250, 0x1F251, prExtendedPictographic}, // E0.6   [2] (🉐..🉑)    Japanese “bargain” button..Japanese “acceptable” button
	{0x1F252, 0x1F2FF, prExtendedPictographic}, // E0.0 [174] (🉒..🋿)    <reserved-1F252>..<reserved-1F2FF>
	{0x1F300, 0x1F30C, prExtendedPictographic}, // E0.6  [13] (🌀..🌌)    cyclone..milky way
	{0x1F30D, 0x1F30E, prExtendedPictographic}, // E0.7   [2] (🌍..🌎)    globe showing Europe-Africa..globe showing Americas
	{0x1F30F, 0x1F30F, prExtendedPictographic}, // E0.6   [1] (🌏)       globe showing Asia-Australia
	{0x1F310, 0x1F310, prExtendedPictographic}, // E1.0   [1] (🌐)       globe with meridians
	{0x1F311, 0x1F311, prExtendedPictographic}, // E0.6   [1] (🌑)       new moon
	{0x1F312, 0x1F312, prExtendedPictographic}, // E1.0   [1] (🌒)       waxing crescent moon
	{0x1F313, 0x1F315, prExtendedPictographic}, // E0.6   [3] (🌓..🌕)    first quarter moon..full moon
	{0x1F316, 0x1F318, prExtendedPictographic}, // E1.0   [3] (🌖..🌘)    waning gibbous moon..waning crescent moon
	{0x1F319, 0x1F319, prExtendedPictographic}, // E0.6   [1] (🌙)       crescent moon
	{0x1F31A, 0x1F31A, prExtendedPictographic}, // E1.0   [1] (🌚)       new moon face
	{0x1F31B, 0x1F31B, prExtendedPictographic}, // E0.6   [1] (🌛)       first quarter moon face
	{0x1F31C, 0x1F31C, prExtendedPictographic}, // E0.7   [1] (🌜)       last quarter moon face
	{0x1F31D, 0x1F31E, prExtendedPictographic}, // E1.0   [2] (🌝..🌞)    full moon face..sun with face
	{0x1F31F, 0x1F320, prExtendedPictographic}, // E0.6   [2] (🌟..🌠)    glowing star..shooting star
	{0x1F321, 0x1F321, prExtendedPictographic}, // E0.7   [1] (🌡️)       thermometer
	{0x1F322, 0x1F323, prExtendedPictographic}, // E0.0   [2] (🌢..🌣)    BLACK DROPLET..WHITE SUN
	{0x1F324, 0x1F32C, prExtendedPictographic}, // E0.7   [9] (🌤️..🌬️)    sun behind small cloud..wind face
	{0x1F32D, 0x1F32F, prExtendedPictographic}, // E1.0   [3] (🌭..🌯)    hot dog..burrito
	{0x1F330, 0x1F331, prExtendedPictographic}, // E0.6   [2] (🌰..🌱)    chestnut..seedling
	{0x1F332, 0x1F333, prExtendedPictographic}, // E1.0   [2] (🌲..🌳)    evergreen tree..deciduous tree
	{0x1F334, 0x1F335, prExtendedPictographic}, // E0.6   [2] (🌴..🌵)    palm tree..cactus
	{0x1F336, 0x1F336, prExtendedPictographic}, // E0.7   [1] (🌶️)       hot pepper
	{0x1F337, 0x1F34A, prExtendedPictographic}, // E0.6  [20] (🌷..🍊)    tulip..tangerine
	{0x1F34B, 0x1F34B, prExtendedPictographic}, // E1.0   [1] (🍋)       lemon
	{0x1F34C, 0x1F34F, prExtendedPictographic}, // E0.6   [4] (🍌..🍏)    banana..green apple
	{0x1F350, 0x1F350, prExtendedPictographic}, // E1.0   [1] (🍐)       pear
	{0x1F351, 0x1F37B, prExtendedPictographic}, // E0.6  [43] (🍑..🍻)    peach..clinking beer mugs
	{0x1F37C, 0x1F37C, prExtendedPictographic}, // E1.0   [1] (🍼)       baby bottle
	{0x1F37D, 0x1F37D, prExtendedPictographic}, // E0.7   [1] (🍽️)       fork and knife with plate
	{0x1F37E, 0x1F37F, prExtendedPictographic}, // E1.0   [2] (🍾..🍿)    bottle with popping cork..popcorn
	{0x1F380, 0x1F393, prExtendedPictographic}, // E0.6  [20] (🎀..🎓)    ribbon..graduation cap
	{0x1F394, 0x1F395, prExtendedPictographic}, // E0.0   [2] (🎔..🎕)    HEART WITH TIP ON THE LEFT..BOUQUET OF FLOWERS
	{0x1F396, 0x1F397, prExtendedPictographic}, // E0.7   [2] (🎖️..🎗️)    military medal..reminder ribbon
	{0x1F398, 0x1F398, prExtendedPictographic}, // E0.0   [1] (🎘)       MUSICAL KEYBOARD WITH JACKS
	{0x1F399, 0x1F39B, prExtendedPictographic}, // E0.7   [3] (🎙️..🎛️)    studio microphone..control knobs
	{0x1F39C, 0x1F39D, prExtendedPictographic}, // E0.0   [2] (🎜..🎝)    BEAMED ASCENDING MUSICAL NOTES..BEAMED DESCENDING MUSICAL NOTES
	{0x1F39E, 0x1F39F, prExtendedPictographic}, // E0.7   [2] (🎞️..🎟️)    film frames..admission tickets
	{0x1F3A0, 0x1F3C4, prExtendedPictographic}, // E0.6  [37] (🎠..🏄)    carousel horse..person surfing
	{0x1F3C5, 0x1F3C5, prExtendedPictographic}, // E1.0   [1] (🏅)       sports medal
	{0x1F3C6, 0x1F3C6, prExtendedPictographic}, // E0.6   [1] (🏆)       trophy
	{0x1F3C7, 0x1F3C7, prExtendedPictographic}, // E1.0   [1] (🏇)       horse racing
	{0x1F3C8, 0x1F3C8, prExtendedPictographic}, // E0.6   [1] (🏈)       american football
	{0x1F3C9, 0x1F3C9, prExtendedPictographic}, // E1.0   [1] (🏉)       rugby football
	{0x1F3CA, 0x1F3CA, prExtendedPictographic}, // E0.6   [1] (🏊)       person swimming
	{0x1F3CB, 0x1F3CE, prExtendedPictographic}, // E0.7   [4] (🏋️..🏎️)    person lifting weights..racing car
	{0x1F3CF, 0x1F3D3, prExtendedPictographic}, // E1.0   [5] (🏏..🏓)    cricket game..ping pong
	{0x1F3D4, 0x1F3DF, prExtendedPictographic}, // E0.7  [12] (🏔️..🏟️)    snow-capped mountain..stadium
	{0x1F3E0, 0x1F3E3, prExtendedPictographic}, // E0.6   [4] (🏠..🏣)    house..Japanese post office
	{0x1F3E4, 0x1F3E4, prExtendedPictographic}, // E1.0   [1] (🏤)       post office
	{0x1F3E5, 0x1F3F0, prExtendedPictographic}, // E0.6  [12] (🏥..🏰)    hospital..castle
	{0x1F3F1, 0x1F3F2, prExtendedPictographic}, // E0.0   [2] (🏱..🏲)    WHITE PENNANT..BLACK PENNANT
	{0x1F3F3, 0x1F3F3, prExtendedPictographic}, // E0.7   [1] (🏳️)       white flag
	{0x1F3F4, 0x1F3F4, prExtendedPictographic}, // E1.0   [1] (🏴)       black flag
	{0x1F3F5, 0x1F3F5, prExtendedPictographic}, // E0.7   [1] (🏵️)       rosette
	{0x1F3F6, 0x1F3F6, prExtendedPictographic}, // E0.0   [1] (🏶)       BLACK ROSETTE
	{0x1F3F7, 0x1F3F7, prExtendedPictographic}, // E0.7   [1] (🏷️)       label
	{0x1F3F8, 0x1F3FA, prExtendedPictographic}, // E1.0   [3] (🏸..🏺)    badminton..amphora
	{0x1F3FB, 0x1F3FF, prExtend},               // Sk   [5] EMOJI MODIFIER FITZPATRICK TYPE-1-2..EMOJI MODIFIER FITZPATRICK TYPE-6
	{0x1F400, 0x1F407, prExtendedPictographic}, // E1.0   [8] (🐀..🐇)    rat..rabbit
	{0x1F408, 0x1F408, prExtendedPictographic}, // E0.7   [1] (🐈)       cat
	{0x1F409, 0x1F40B, prExtendedPictographic}, // E1.0   [3] (🐉..🐋)    dragon..whale
	{0x1F40C, 0x1F40E, prExtendedPictographic}, // E0.6   [3] (🐌..🐎)    snail..horse
	{0x1F40F, 0x1F410, prExtendedPictographic}, // E1.0   [2] (🐏..🐐)    ram..goat
	{0x1F411, 0x1F412, prExtendedPictographic}, // E0.6   [2] (🐑..🐒)    ewe..monkey
	{0x1F413, 0x1F413, prExtendedPictographic}, // E1.0   [1] (🐓)       rooster
	{0x1F414, 0x1F414, prExtendedPictographic}, // E0.6   [1] (🐔)       chicken
	{0x1F415, 0x1F415, prExtendedPictographic}, // E0.7   [1] (🐕)       dog
	{0x1F416, 0x1F416, prExtendedPictographic}, // E1.0   [1] (🐖)       pig
	{0x1F417, 0x1F429, prExtendedPictographic}, // E0.6  [19] (🐗..🐩)    boar..poodle
	{0x1F42A, 0x1F42A, prExtendedPictographic}, // E1.0   [1] (🐪)       camel
	{0x1F42B, 0x1F43E, prExtendedPictographic}, // E0.6  [20] (🐫..🐾)    two-hump camel..paw prints
	{0x1F43F, 0x1F43F, prExtendedPictographic}, // E0.7   [1] (🐿️)       chipmunk
	{0x1F440, 0x1F440, prExtendedPictographic}, // E0.6   [1] (👀)       eyes
	{0x1F441, 0x1F441, prExtendedPictographic}, // E0.7   [1] (👁️)       eye
	{0x1F442, 0x1F464, prExtendedPictographic}, // E0.6  [35] (👂..👤)    ear..bust in silhouette
	{0x1F465, 0x1F465, prExtendedPictographic}, // E1.0   [1] (👥)       busts in silhouette
	{0x1F466, 0x1F46B, prExtendedPictographic}, // E0.6   [6] (👦..👫)    boy..woman and man holding hands
	{0x1F46C, 0x1F46D, prExtendedPictographic}, // E1.0   [2] (👬..👭)    men holding hands..women holding hands
	{0x1F46E, 0x1F4AC, prExtendedPictographic}, // E0.6  [63] (👮..💬)    police officer..speech balloon
	{0x1F4AD, 0x1F4AD, prExtendedPictographic}, // E1.0   [1] (💭)       thought balloon
	{0x1F4AE, 0x1F4B5, prExtendedPictographic}, // E0.6   [8] (💮..💵)    white flower..dollar banknote
	{0x1F4B6, 0x1F4B7, prExtendedPictographic}, // E1.0   [2] (💶..💷)    euro banknote..pound banknote
	{0x1F4B8, 0x1F4EB, prExtendedPictographic}, // E0.6  [52] (💸..📫)    money with wings..closed mailbox with raised flag
	{0x1F4EC, 0x1F4ED, prExtendedPictographic}, // E0.7   [2] (📬..📭)    open mailbox with raised flag..open mailbox with lowered flag
	{0x1F4EE, 0x1F4EE, prExtendedPictographic}, // E0.6   [1] (📮)       postbox
	{0x1F4EF, 0x1F4EF, prExtendedPictographic}, // E1.0   [1] (📯)       postal horn
	{0x1F4F0, 0x1F4F4, prExtendedPictographic}, // E0.6   [5] (📰..📴)    newspaper..mobile phone off
	{0x1F4F5, 0x1F4F5, prExtendedPictographic}, // E1.0   [1] (📵)       no mobile phones
	{0x1F4F6, 0x1F4F7, prExtendedPictographic}, // E0.6   [2] (📶..📷)    antenna bars..camera
	{0x1F4F8, 0x1F4F8, prExtendedPictographic}, // E1.0   [1] (📸)       camera with flash
	{0x1F4F9, 0x1F4FC, prExtendedPictographic}, // E0.6   [4] (📹..📼)    video camera..videocassette
	{0x1F4FD, 0x1F4FD, prExtendedPictographic}, // E0.7   [1] (📽️)       film projector
	{0x1F4FE, 0x1F4FE, prExtendedPictographic}, // E0.0   [1] (📾)       PORTABLE STEREO
	{0x1F4FF, 0x1F502, prExtendedPictographic}, // E1.0   [4] (📿..🔂)    prayer beads..repeat single button
	{0x1F503, 0x1F503, prExtendedPictographic}, // E0.6   [1] (🔃)       clockwise vertical arrows
	{0x1F504, 0x1F507, prExtendedPictographic}, // E1.0   [4] (🔄..🔇)    counterclockwise arrows button..muted speaker
	{0x1F508, 0x1F508, prExtendedPictographic}, // E0.7   [1] (🔈)       speaker low volume
	{0x1F509, 0x1F509, prExtendedPictographic}, // E1.0   [1] (🔉)       speaker medium volume
	{0x1F50A, 0x1F514, prExtendedPictographic}, // E0.6  [11] (🔊..🔔)    speaker high volume..bell
	{0x1F515, 0x1F515, prExtendedPictographic}, // E1.0   [1] (🔕)       bell with slash
	{0x1F516, 0x1F52B, prExtendedPictographic}, // E0.6  [22] (🔖..🔫)    bookmark..water pistol
	{0x1F52C, 0x1F52D, prExtendedPictographic}, // E1.0   [2] (🔬..🔭)    microscope..telescope
	{0x1F52E, 0x1F53D, prExtendedPictographic}, // E0.6  [16] (🔮..🔽)    crystal ball..downwards button
	{0x1F546, 0x1F548, prExtendedPictographic}, // E0.0   [3] (🕆..🕈)    WHITE LATIN CROSS..CELTIC CROSS
	{0x1F549, 0x1F54A, prExtendedPictographic}, // E0.7   [2] (🕉️..🕊️)    om..dove
	{0x1F54B, 0x1F54E, prExtendedPictographic}, // E1.0   [4] (🕋..🕎)    kaaba..menorah
	{0x1F54F, 0x1F54F, prExtendedPictographic}, // E0.0   [1] (🕏)       BOWL OF HYGIEIA
	{0x1F550, 0x1F55B, prExtendedPictographic}, // E0.6  [12] (🕐..🕛)    one o’clock..twelve o’clock
	{0x1F55C, 0x1F567, prExtendedPictographic}, // E0.7  [12] (🕜..🕧)    one-thirty..twelve-thirty
	{0x1F568, 0x1F56E, prExtendedPictographic}, // E0.0   [7] (🕨..🕮)    RIGHT SPEAKER..BOOK
	{0x1F56F, 0x1F570, prExtendedPictographic}, // E0.7   [2] (🕯️..🕰️)    candle..mantelpiece clock
	{0x1F571, 0x1F572, prExtendedPictographic}, // E0.0   [2] (🕱..🕲)    BLACK SKULL AND CROSSBONES..NO PIRACY
	{0x1F573, 0x1F579, prExtendedPictographic}, // E0.7   [7] (🕳️..🕹️)    hole..joystick
	{0x1F57A, 0x1F57A, prExtendedPictographic}, // E3.0   [1] (🕺)       man dancing
	{0x1F57B, 0x1F586, prExtendedPictographic}, // E0.0  [12] (🕻..🖆)    LEFT HAND TELEPHONE RECEIVER..PEN OVER STAMPED ENVELOPE
	{0x1F587, 0x1F587, prExtendedPictographic}, // E0.7   [1] (🖇️)       linked paperclips
	{0x1F588, 0x1F589, prExtendedPictographic}, // E0.0   [2] (🖈..🖉)    BLACK PUSHPIN..LOWER LEFT PENCIL
	{0x1F58A, 0x1F58D, prExtendedPictographic}, // E0.7   [4] (🖊️..🖍️)    pen..crayon
	{0x1F58E, 0x1F58F, prExtendedPictographic}, // E0.0   [2] (🖎..🖏)    LEFT WRITING HAND..TURNED OK HAND SIGN
	{0x1F590, 0x1F590, prExtendedPictographic}, // E0.7   [1] (🖐️)       hand with fingers splayed
	{0x1F591, 0x1F594, prExtendedPictographic}, // E0.0   [4] (🖑..🖔)    REVERSED RAISED HAND WITH FINGERS SPLAYED..REVERSED VICTORY HAND
	{0x1F595, 0x1F596, prExtendedPictographic}, // E1.0   [2] (🖕..🖖)    middle finger..vulcan salute
	{0x1F597, 0x1F5A3, prExtendedPictographic}, // E0.0  [13] (🖗..🖣)    WHITE DOWN POINTING LEFT HAND INDEX..BLACK DOWN POINTING BACKHAND INDEX
	{0x1F5A4, 0x1F5A4, prExtendedPictographic}, // E3.0   [1] (🖤)       black heart
	{0x1F5A5, 0x1F5A5, prExtendedPictographic}, // E0.7   [1] (🖥️)       desktop computer
	{0x1F5A6, 0x1F5A7, prExtendedPictographic}, // E0.0   [2] (🖦..🖧)    KEYBOARD AND MOUSE..THREE NETWORKED COMPUTERS
	{0x1F5A8, 0x1F5A8, prExtendedPictographic}, // E0.7   [1] (🖨️)       printer
	{0x1F5A9, 0x1F5B0, prExtendedPictographic}, // E0.0   [8] (🖩..🖰)    POCKET CALCULATOR..TWO BUTTON MOUSE
	{0x1F5B1, 0x1F5B2, prExtendedPictographic}, // E0.7   [2] (🖱️..🖲️)    computer mouse..trackball
	{0x1F5B3, 0x1F5BB, prExtendedPictographic}, // E0.0   [9] (🖳..🖻)    OLD PERSONAL COMPUTER..DOCUMENT WITH PICTURE
	{0x1F5BC, 0x1F5BC, prExtendedPictographic}, // E0.7   [1] (🖼️)       framed picture
	{0x1F5BD, 0x1F5C1, prExtendedPictographic}, // E0.0   [5] (🖽..🗁)    FRAME WITH TILES..OPEN FOLDER
	{0x1F5C2, 0x1F5C4, prExtendedPictographic}, // E0.7   [3] (🗂️..🗄️)    card index dividers..file cabinet
	{0x1F5C5, 0x1F5D0, prExtendedPictographic}, // E0.0  [12] (🗅..🗐)    EMPTY NOTE..PAGES
	{0x1F5D1, 0x1F5D3, prExtendedPictographic}, // E0.7   [3] (🗑️..🗓️)    wastebasket..spiral calendar
	{0x1F5D4, 0x1F5DB, prExtendedPictographic}, // E0.0   [8] (🗔..🗛)    DESKTOP WINDOW..DECREASE FONT SIZE SYMBOL
	{0x1F5DC, 0x1F5DE, prExtendedPictographic}, // E0.7   [3] (🗜️..🗞️)    clamp..rolled-up newspaper
	{0x1F5DF, 0x1F5E0, prExtendedPictographic}, // E0.0   [2] (🗟..🗠)    PAGE WITH CIRCLED TEXT..STOCK CHART
	{0x1F5E1, 0x1F5E1, prExtendedPictographic}, // E0.7   [1] (🗡️)       dagger
	{0x1F5E2, 0x1F5E2, prExtendedPictographic}, // E0.0   [1] (🗢)       LIPS
	{0x1F5E3, 0x1F5E3, prExtendedPictographic}, // E0.7   [1] (🗣️)       speaking head
	{0x1F5E4, 0x1F5E7, prExtendedPictographic}, // E0.0   [4] (🗤..🗧)    THREE RAYS ABOVE..THREE RAYS RIGHT
	{0x1F5E8, 0x1F5E8, prExtendedPictographic}, // E2.0   [1] (🗨️)       left speech bubble
	{0x1F5E9, 0x1F5EE, prExtendedPictographic}, // E0.0   [6] (🗩..🗮)    RIGHT SPEECH BUBBLE..LEFT ANGER BUBBLE
	{0x1F5EF, 0x1F5EF, prExtendedPictographic}, // E0.7   [1] (🗯️)       right anger bubble
	{0x1F5F0, 0x1F5F2, prExtendedPictographic}, // E0.0   [3] (🗰..🗲)    MOOD BUBBLE..LIGHTNING MOOD
	{0x1F5F3, 0x1F5F3, prExtendedPictographic}, // E0.7   [1] (🗳️)       ballot box with ballot
	{0x1F5F4, 0x1F5F9, prExtendedPictographic}, // E0.0   [6] (🗴..🗹)    BALLOT SCRIPT X..BALLOT BOX WITH BOLD CHECK
	{0x1F5FA, 0x1F5FA, prExtendedPictographic}, // E0.7   [1] (🗺️)       world map
	{0x1F5FB, 0x1F5FF, prExtendedPictographic}, // E0.6   [5] (🗻..🗿)    mount fuji..moai
	{0x1F600, 0x1F600, prExtendedPictographic}, // E1.0   [1] (😀)       grinning face
	{0x1F601, 0x1F606, prExtendedPictographic}, // E0.6   [6] (😁..😆)    beaming face with smiling eyes..grinning squinting face
	{0x1F607, 0x1F608, prExtendedPictographic}, // E1.0   [2] (😇..😈)    smiling face with halo..smiling face with horns
	{0x1F609, 0x1F60D, prExtendedPictographic}, // E0.6   [5] (😉..😍)    winking face..smiling face with heart-eyes
	{0x1F60E, 0x1F60E, prExtendedPictographic}, // E1.0   [1] (😎)       smiling face with sunglasses
	{0x1F60F, 0x1F60F, prExtendedPictographic}, // E0.6   [1] (😏)       smirking face
	{0x1F610, 0x1F610, prExtendedPictographic}, // E0.7   [1] (😐)       neutral face
	{0x1F611, 0x1F611, prExtendedPictographic}, // E1.0   [1] (😑)       expressionless face
	{0x1F612, 0x1F614, prExtendedPictographic}, // E0.6   [3] (😒..😔)    unamused face..pensive face
	{0x1F615, 0x1F615, prExtendedPictographic}, // E1.0   [1] (😕)       confused face
	{0x1F616, 0x1F616, prExtendedPictographic}, // E0.6   [1] (😖)       confounded face
	{0x1F617, 0x1F617, prExtendedPictographic}, // E1.0   [1] (😗)       kissing face
	{0x1F618, 0x1F618, prExtendedPictographic}, // E0.6   [1] (😘)       face blowing a kiss
	{0x1F619, 0x1F619, prExtendedPictographic}, // E1.0   [1] (😙)       kissing face with smiling eyes
	{0x1F61A, 0x1F61A, prExtendedPictographic}, // E0.6   [1] (😚)       kissing face with closed eyes
	{0x1F61B, 0x1F61B, prExtendedPictographic}, // E1.0   [1] (😛)       face with tongue
	{0x1F61C, 0x1F61E, prExtendedPictographic}, // E0.6   [3] (😜..😞)    winking face with tongue..disappointed face
	{0x1F61F, 0x1F61F, prExtendedPictographic}, // E1.0   [1] (😟)       worried face
	{0x1F620, 0x1F625, prExtendedPictographic}, // E0.6   [6] (😠..😥)    angry face..sad but relieved face
	{0x1F626, 0x1F627, prExtendedPictographic}, // E1.0   [2] (😦..😧)    frowning face with open mouth..anguished face
	{0x1F628, 0x1F62B, prExtendedPictographic}, // E0.6   [4] (😨..😫)    fearful face..tired face
	{0x1F62C, 0x1F62C, prExtendedPictographic}, // E1.0   [1] (😬)       grimacing face
	{0x1F62D, 0x1F62D, prExtendedPictographic}, // E0.6   [1] (😭)       loudly crying face
	{0x1F62E, 0x1F62F, prExtendedPictographic}, // E1.0   [2] (😮..😯)    face with open mouth..hushed face
	{0x1F630, 0x1F633, prExtendedPictographic}, // E0.6   [4] (😰..😳)    anxious face with sweat..flushed face
	{0x1F634, 0x1F634, prExtendedPictographic}, // E1.0   [1] (😴)       sleeping face
	{0x1F635, 0x1F635, prExtendedPictographic}, // E0.6   [1] (😵)       face with crossed-out eyes
	{0x1F636, 0x1F636, prExtendedPictographic}, // E1.0   [1] (😶)       face without mouth
	{0x1F637, 0x1F640, prExtendedPictographic}, // E0.6  [10] (😷..🙀)    face with medical mask..weary cat
	{0x1F641, 0x1F644, prExtendedPictographic}, // E1.0   [4] (🙁..🙄)    slightly frowning face..face with rolling eyes
	{0x1F645, 0x1F64F, prExtendedPictographic}, // E0.6  [11] (🙅..🙏)    person gesturing NO..folded hands
	{0x1F680, 0x1F680, prExtendedPictographic}, // E0.6   [1] (🚀)       rocket
	{0x1F681, 0x1F682, prExtendedPictographic}, // E1.0   [2] (🚁..🚂)    helicopter..locomotive
	{0x1F683, 0x1F685, prExtendedPictographic}, // E0.6   [3] (🚃..🚅)    railway car..bullet train
	{0x1F686, 0x1F686, prExtendedPictographic}, // E1.0   [1] (🚆)       train
	{0x1F687, 0x1F687, prExtendedPictographic}, // E0.6   [1] (🚇)       metro
	{0x1F688, 0x1F688, prExtendedPictographic}, // E1.0   [1] (🚈)       light rail
	{0x1F689, 0x1F689, prExtendedPictographic}, // E0.6   [1] (🚉)       station
	{0x1F68A, 0x1F68B, prExtendedPictographic}, // E1.0   [2] (🚊..🚋)    tram..tram car
	{0x1F68C, 0x1F68C, prExtendedPictographic}, // E0.6   [1] (🚌)       bus
	{0x1F68D, 0x1F68D, prExtendedPictographic}, // E0.7   [1] (🚍)       oncoming bus
	{0x1F68E, 0x1F68E, prExtendedPictographic}, // E1.0   [1] (🚎)       trolleybus
	{0x1F68F, 0x1F68F, prExtendedPictographic}, // E0.6   [1] (🚏)       bus stop
	{0x1F690, 0x1F690, prExtendedPictographic}, // E1.0   [1] (🚐)       minibus
	{0x1F691, 0x1F693, prExtendedPictographic}, // E0.6   [3] (🚑..🚓)    ambulance..police car
	{0x1F694, 0x1F694, prExtendedPictographic}, // E0.7   [1] (🚔)       oncoming police car
	{0x1F695, 0x1F695, prExtendedPictographic}, // E0.6   [1] (🚕)       taxi
	{0x1F696, 0x1F696, prExtendedPictographic}, // E1.0   [1] (🚖)       oncoming taxi
	{0x1F697, 0x1F697, prExtendedPictographic}, // E0.6   [1] (🚗)       automobile
	{0x1F698, 0x1F698, prExtendedPictographic}, // E0.7   [1] (🚘)       oncoming automobile
	{0x1F699, 0x1F69A, prExtendedPictographic}, // E0.6   [2] (🚙..🚚)    sport utility vehicle..delivery truck
	{0x1F69B, 0x1F6A1, prExtendedPictographic}, // E1.0   [7] (🚛..🚡)    articulated lorry..aerial tramway
	{0x1F6A2, 0x1F6A2, prExtendedPictographic}, // E0.6   [1] (🚢)       ship
	{0x1F6A3, 0x1F6A3, prExtendedPictographic}, // E1.0   [1] (🚣)       person rowing boat
	{0x1F6A4, 0x1F6A5, prExtendedPictographic}, // E0.6   [2] (🚤..🚥)    speedboat..horizontal traffic light
	{0x1F6A6, 0x1F6A6, prExtendedPictographic}, // E1.0   [1] (🚦)       vertical traffic light
	{0x1F6A7, 0x1F6AD, prExtendedPictographic}, // E0.6   [7] (🚧..🚭)    construction..no smoking
	{0x1F6AE, 0x1F6B1, prExtendedPictographic}, // E1.0   [4] (🚮..🚱)    litter in bin sign..non-potable water
	{0x1F6B2, 0x1F6B2, prExtendedPictographic}, // E0.6   [1] (🚲)       bicycle
	{0x1F6B3, 0x1F6B5, prExtendedPictographic}, // E1.0   [3] (🚳..🚵)    no bicycles..person mountain biking
	{0x1F6B6, 0x1F6B6, prExtendedPictographic}, // E0.6   [1] (🚶)       person walking
	{0x1F6B7, 0x1F6B8, prExtendedPictographic}, // E1.0   [2] (🚷..🚸)    no pedestrians..children crossing
	{0x1F6B9, 0x1F6BE, prExtendedPictographic}, // E0.6   [6] (🚹..🚾)    men’s room..water closet
	{0x1F6BF, 0x1F6BF, prExtendedPictographic}, // E1.0   [1] (🚿)       shower
	{0x1F6C0, 0x1F6C0, prExtendedPictographic}, // E0.6   [1] (🛀)       person taking bath
	{0x1F6C1, 0x1F6C5, prExtendedPictographic}, // E1.0   [5] (🛁..🛅)    bathtub..left luggage
	{0x1F6C6, 0x1F6CA, prExtendedPictographic}, // E0.0   [5] (🛆..🛊)    TRIANGLE WITH ROUNDED CORNERS..GIRLS SYMBOL
	{0x1F6CB, 0x1F6CB, prExtendedPictographic}, // E0.7   [1] (🛋️)       couch and lamp
	{0x1F6CC, 0x1F6CC, prExtendedPictographic}, // E1.0   [1] (🛌)       person in bed
	{0x1F6CD, 0x1F6CF, prExtendedPictographic}, // E0.7   [3] (🛍️..🛏️)    shopping bags..bed
	{0x1F6D0, 0x1F6D0, prExtendedPictographic}, // E1.0   [1] (🛐)       place of worship
	{0x1F6D1, 0x1F6D2, prExtendedPictographic}, // E3.0   [2] (🛑..🛒)    stop sign..shopping cart
	{0x1F6D3, 0x1F6D4, prExtendedPictographic}, // E0.0   [2] (🛓..🛔)    STUPA..PAGODA
	{0x1F6D5, 0x1F6D5, prExtendedPictographic}, // E12.0  [1] (🛕)       hindu temple
	{0x1F6D6, 0x1F6D7, prExtendedPictographic}, // E13.0  [2] (🛖..🛗)    hut..elevator
	{0x1F6D8, 0x1F6DB, prExtendedPictographic}, // E0.0   [4] (🛘..🛛)    <reserved-1F6D8>..<reserved-1F6DB>
	{0x1F6DC, 0x1F6DC, prExtendedPictographic}, // E15.0  [1] (🛜)       wireless
	{0x1F6DD, 0x1F6DF, prExtendedPictographic}, // E14.0  [3] (🛝..🛟)    playground slide..ring buoy
	{0x1F6E0, 0x1F6E5, prExtendedPictographic}, // E0.7   [6] (🛠️..🛥️)    hammer and wrench..motor boat
	{0x1F6E6, 0x1F6E8, prExtendedPictographic}, // E0.0   [3] (🛦..🛨)    UP-POINTING MILITARY AIRPLANE..UP-POINTING SMALL AIRPLANE
	{0x1F6E9, 0x1F6E9, prExtendedPictographic}, // E0.7   [1] (🛩️)       small airplane
	{0x1F6EA, 0x1F6EA, prExtendedPictographic}, // E0.0   [1] (🛪)       NORTHEAST-POINTING AIRPLANE
	{0x1F6EB, 0x1F6EC, prExtendedPictographic}, // E1.0   [2] (🛫..🛬)    airplane departure..airplane arrival
	{0x1F6ED, 0x1F6EF, prExtendedPictographic}, // E0.0   [3] (🛭..🛯)    <reserved-1F6ED>..<reserved-1F6EF>
	{0x1F6F0, 0x1F6F0, prExtendedPictographic}, // E0.7   [1] (🛰️)       satellite
	{0x1F6F1, 0x1F6F2, prExtendedPictographic}, // E0.0   [2] (🛱..🛲)    ONCOMING FIRE ENGINE..DIESEL LOCOMOTIVE
	{0x1F6F3, 0x1F6F3, prExtendedPictographic}, // E0.7   [1] (🛳️)       passenger ship
	{0x1F6F4, 0x1F6F6, prExtendedPictographic}, // E3.0   [3] (🛴..🛶)    kick scooter..canoe
	{0x1F6F7, 0x1F6F8, prExtendedPictographic}, // E5.0   [2] (🛷..🛸)    sled..flying saucer
	{0x1F6F9, 0x1F6F9, prExtendedPictographic}, // E11.0  [1] (🛹)       skateboard
	{0x1F6FA, 0x1F6FA, prExtendedPictographic}, // E12.0  [1] (🛺)       auto rickshaw
	{0x1F6FB, 0x1F6FC, prExtendedPictographic}, // E13.0  [2] (🛻..🛼)    pickup truck..roller skate
	{0x1F6FD, 0x1F6FF, prExtendedPictographic}, // E0.0   [3] (🛽..🛿)    <reserved-1F6FD>..<reserved-1F6FF>
	{0x1F774, 0x1F77F, prExtendedPictographic}, // E0.0  [12] (🝴..🝿)    LOT OF FORTUNE..ORCUS
	{0x1F7D5, 0x1F7DF, prExtendedPictographic}, // E0.0  [11] (🟕..🟟)    CIRCLED TRIANGLE..<reserved-1F7DF>
	{0x1F7E0, 0x1F7EB, prExtendedPictographic}, // E12.0 [12] (🟠..🟫)    orange circle..brown square
	{0x1F7EC, 0x1F7EF, prExtendedPictographic}, // E0.0   [4] (🟬..🟯)    <reserved-1F7EC>..<reserved-1F7EF>
	{0x1F7F0, 0x1F7F0, prExtendedPictographic}, // E14.0  [1] (🟰)       heavy equals sign
	{0x1F7F1, 0x1F7FF, prExtendedPictographic}, // E0.0  [15] (🟱..🟿)    <reserved-1F7F1>..<reserved-1F7FF>
	{0x1F80C, 0x1F80F, prExtendedPictographic}, // E0.0   [4] (🠌..🠏)    <reserved-1F80C>..<reserved-1F80F>
	{0x1F848, 0x1F84F, prExtendedPictographic}, // E0.0   [8] (🡈..🡏)    <reserved-1F848>..<reserved-1F84F>
	{0x1F85A, 0x1F85F, prExtendedPictographic}, // E0.0   [6] (🡚..🡟)    <reserved-1F85A>..<reserved-1F85F>
	{0x1F888, 0x1F88F, prExtendedPictographic}, // E0.0   [8] (🢈..🢏)    <reserved-1F888>..<reserved-1F88F>
	{0x1F8AE, 0x1F8FF, prExtendedPictographic}, // E0.0  [82] (🢮..🣿)    <reserved-1F8AE>..<reserved-1F8FF>
	{0x1F90C, 0x1F90C, prExtendedPictographic}, // E13.0  [1] (🤌)       pinched fingers
	{0x1F90D, 0x1F90F, prExtendedPictographic}, // E12.0  [3] (🤍..🤏)    white heart..pinching hand
	{0x1F910, 0x1F918, prExtendedPictographic}, // E1.0   [9] (🤐..🤘)    zipper-mouth face..sign of the horns
	{0x1F919, 0x1F91E, prExtendedPictographic}, // E3.0   [6] (🤙..🤞)    call me hand..crossed fingers
	{0x1F91F, 0x1F91F, prExtendedPictographic}, // E5.0   [1] (🤟)       love-you gesture
	{0x1F920, 0x1F927, prExtendedPictographic}, // E3.0   [8] (🤠..🤧)    cowboy hat face..sneezing face
	{0x1F928, 0x1F92F, prExtendedPictographic}, // E5.0   [8] (🤨..🤯)    face with raised eyebrow..exploding head
	{0x1F930, 0x1F930, prExtendedPictographic}, // E3.0   [1] (🤰)       pregnant woman
	{0x1F931, 0x1F932, prExtendedPictographic}, // E5.0   [2] (🤱..🤲)    breast-feeding..palms up together
	{0x1F933, 0x1F93A, prExtendedPictographic}, // E3.0   [8] (🤳..🤺)    selfie..person fencing
	{0x1F93C, 0x1F93E, prExtendedPictographic}, // E3.0   [3] (🤼..🤾)    people wrestling..person playing handball
	{0x1F93F, 0x1F93F, prExtendedPictographic}, // E12.0  [1] (🤿)       diving mask
	{0x1F940, 0x1F945, prExtendedPictographic}, // E3.0   [6] (🥀..🥅)    wilted flower..goal net
	{0x1F947, 0x1F94B, prExtendedPictographic}, // E3.0   [5] (🥇..🥋)    1st place medal..martial arts uniform
	{0x1F94C, 0x1F94C, prExtendedPictographic}, // E5.0   [1] (🥌)       curling stone
	{0x1F94D, 0x1F94F, prExtendedPictographic}, // E11.0  [3] (🥍..🥏)    lacrosse..flying disc
	{0x1F950, 0x1F95E, prExtendedPictographic}, // E3.0  [15] (🥐..🥞)    croissant..pancakes
	{0x1F95F, 0x1F96B, prExtendedPictographic}, // E5.0  [13] (🥟..🥫)    dumpling..canned food
	{0x1F96C, 0x1F970, prExtendedPictographic}, // E11.0  [5] (🥬..🥰)    leafy green..smiling face with hearts
	{0x1F971, 0x1F971, prExtendedPictographic}, // E12.0  [1] (🥱)       yawning face
	{0x1F972, 0x1F972, prExtendedPictographic}, // E13.0  [1] (🥲)       smiling face with tear
	{0x1F973, 0x1F976, prExtendedPictographic}, // E11.0  [4] (🥳..🥶)    partying face..cold face
	{0x1F977, 0x1F978, prExtendedPictographic}, // E13.0  [2] (🥷..🥸)    ninja..disguised face
	{0x1F979, 0x1F979, prExtendedPictographic}, // E14.0  [1] (🥹)       face holding back tears
	{0x1F97A, 0x1F97A, prExtendedPictographic}, // E11.0  [1] (🥺)       pleading face
	{0x1F97B, 0x1F97B, prExtendedPictographic}, // E12.0  [1] (🥻)       sari
	{0x1F97C, 0x1F97F, prExtendedPictographic}, // E11.0  [4] (🥼..🥿)    lab coat..flat shoe
	{0x1F980, 0x1F984, prExtendedPictographic}, // E1.0   [5] (🦀..🦄)    crab..unicorn
	{0x1F985, 0x1F991, prExtendedPictographic}, // E3.0  [13] (🦅..🦑)    eagle..squid
	{0x1F992, 0x1F997, prExtendedPictographic}, // E5.0   [6] (🦒..🦗)    giraffe..cricket
	{0x1F998, 0x1F9A2, prExtendedPictographic}, // E11.0 [11] (🦘..🦢)    kangaroo..swan
	{0x1F9A3, 0x1F9A4, prExtendedPictographic}, // E13.0  [2] (🦣..🦤)    mammoth..dodo
	{0x1F9A5, 0x1F9AA, prExtendedPictographic}, // E12.0  [6] (🦥..🦪)    sloth..oyster
	{0x1F9AB, 0x1F9AD, prExtendedPictographic}, // E13.0  [3] (🦫..🦭)    beaver..seal
	{0x1F9AE, 0x1F9AF, prExtendedPictographic}, // E12.0  [2] (🦮..🦯)    guide dog..white cane
	{0x1F9B0, 0x1F9B9, prExtendedPictographic}, // E11.0 [10] (🦰..🦹)    red hair..supervillain
	{0x1F9BA, 0x1F9BF, prExtendedPictographic}, // E12.0  [6] (🦺..🦿)    safety vest..mechanical leg
	{0x1F9C0, 0x1F9C0, prExtendedPictographic}, // E1.0   [1] (🧀)       cheese wedge
	{0x1F9C1, 0x1F9C2, prExtendedPictographic}, // E11.0  [2] (🧁..🧂)    cupcake..salt
	{0x1F9C3, 0x1F9CA, prExtendedPictographic}, // E12.0  [8] (🧃..🧊)    beverage box..ice
	{0x1F9CB, 0x1F9CB, prExtendedPictographic}, // E13.0  [1] (🧋)       bubble tea
	{0x1F9CC, 0x1F9CC, prExtendedPictographic}, // E14.0  [1] (🧌)       troll
	{0x1F9CD, 0x1F9CF, prExtendedPictographic}, // E12.0  [3] (🧍..🧏)    person standing..deaf person
	{0x1F9D0, 0x1F9E6, prExtendedPictographic}, // E5.0  [23] (🧐..🧦)    face with monocle..socks
	{0x1F9E7, 0x1F9FF, prExtendedPictographic}, // E11.0 [25] (🧧..🧿)    red envelope..nazar amulet
	{0x1FA00, 0x1FA6F, prExtendedPictographic}, // E0.0 [112] (🨀..🩯)    NEUTRAL CHESS KING..<reserved-1FA6F>
	{0x1FA70, 0x1FA73, prExtendedPictographic}, // E12.0  [4] (🩰..🩳)    ballet shoes..shorts
	{0x1FA74, 0x1FA74, prExtendedPictographic}, // E13.0  [1] (🩴)       thong sandal
	{0x1FA75, 0x1FA77, prExtendedPictographic}, // E15.0  [3] (🩵..🩷)    light blue heart..pink heart
	{0x1FA78, 0x1FA7A, prExtendedPictographic}, // E12.0  [3] (🩸..🩺)    drop of blood..stethoscope
	{0x1FA7B, 0x1FA7C, prExtendedPictographic}, // E14.0  [2] (🩻..🩼)    x-ray..crutch
	{0x1FA7D, 0x1FA7F, prExtendedPictographic}, // E0.0   [3] (🩽..🩿)    <reserved-1FA7D>..<reserved-1FA7F>
	{0x1FA80, 0x1FA82, prExtendedPictographic}, // E12.0  [3] (🪀..🪂)    yo-yo..parachute
	{0x1FA83, 0x1FA86, prExtendedPictographic}, // E13.0  [4] (🪃..🪆)    boomerang..nesting dolls
	{0x1FA87, 0x1FA88, prExtendedPictographic}, // E15.0  [2] (🪇..🪈)    maracas..flute
	{0x1FA89, 0x1FA8F, prExtendedPictographic}, // E0.0   [7] (🪉..🪏)    <reserved-1FA89>..<reserved-1FA8F>
	{0x1FA90, 0x1FA95, prExtendedPictographic}, // E12.0  [6] (🪐..🪕)    ringed planet..banjo
	{0x1FA96, 0x1FAA8, prExtendedPictographic}, // E13.0 [19] (🪖..🪨)    military helmet..rock
	{0x1FAA9, 0x1FAAC, prExtendedPictographic}, // E14.0  [4] (🪩..🪬)    mirror ball..hamsa
	{0x1FAAD, 0x1FAAF, prExtendedPictographic}, // E15.0  [3] (🪭..🪯)    folding hand fan..khanda
	{0x1FAB0, 0x1FAB6, prExtendedPictographic}, // E13.0  [7] (🪰..🪶)    fly..feather
	{0x1FAB7, 0x1FABA, prExtendedPictographic}, // E14.0  [4] (🪷..🪺)    lotus..nest with eggs
	{0x1FABB, 0x1FABD, prExtendedPictographic}, // E15.0  [3] (🪻..🪽)    hyacinth..wing
	{0x1FABE, 0x1FABE, prExtendedPictographic}, // E0.0   [1] (🪾)       <reserved-1FABE>
	{0x1FABF, 0x1FABF, prExtendedPictographic}, // E15.0  [1] (🪿)       goose
	{0x1FAC0, 0x1FAC2, prExtendedPictographic}, // E13.0  [3] (🫀..🫂)    anatomical heart..people hugging
	{0x1FAC3, 0x1FAC5, prExtendedPictographic}, // E14.0  [3] (🫃..🫅)    pregnant man..person with crown
	{0x1FAC6, 0x1FACD, prExtendedPictographic}, // E0.0   [8] (🫆..🫍)    <reserved-1FAC6>..<reserved-1FACD>
	{0x1FACE, 0x1FACF, prExtendedPictographic}, // E15.0  [2] (🫎..🫏)    moose..donkey
	{0x1FAD0, 0x1FAD6, prExtendedPictographic}, // E13.0  [7] (🫐..🫖)    blueberries..teapot
	{0x1FAD7, 0x1FAD9, prExtendedPictographic}, // E14.0  [3] (🫗..🫙)    pouring liquid..jar
	{0x1FADA, 0x1FADB, prExtendedPictographic}, // E15.0  [2] (🫚..🫛)    ginger root..pea pod
	{0x1FADC, 0x1FADF, prExtendedPictographic}, // E0.0   [4] (🫜..🫟)    <reserved-1FADC>..<reserved-1FADF>
	{0x1FAE0, 0x1FAE7, prExtendedPictographic}, // E14.0  [8] (🫠..🫧)    melting face..bubbles
	{0x1FAE8, 0x1FAE8, prExtendedPictographic}, // E15.0  [1] (🫨)       shaking face
	{0x1FAE9, 0x1FAEF, prExtendedPictographic}, // E0.0   [7] (🫩..🫯)    <reserved-1FAE9>..<reserved-1FAEF>
	{0x1FAF0, 0x1FAF6, prExtendedPictographic}, // E14.0  [7] (🫰..🫶)    hand with index finger and thumb crossed..heart hands
	{0x1FAF7, 0x1FAF8, prExtendedPictographic}, // E15.0  [2] (🫷..🫸)    leftwards pushing hand..rightwards pushing hand
	{0x1FAF9, 0x1FAFF, prExtendedPictographic}, // E0.0   [7] (🫹..🫿)    <reserved-1FAF9>..<reserved-1FAFF>
	{0x1FC00, 0x1FFFD, prExtendedPictographic}, // E0.0[1022] (🰀..🿽)    <reserved-1FC00>..<reserved-1FFFD>
	{0xE0000, 0xE0000, prControl},              // Cn       <reserved-E0000>
	{0xE0001, 0xE0001, prControl},              // Cf       LANGUAGE TAG
	{0xE0002, 0xE001F, prControl},              // Cn  [30] <reserved-E0002>..<reserved-E001F>
	{0xE0020, 0xE007F, prExtend},               // Cf  [96] TAG SPACE..CANCEL TAG
	{0xE0080, 0xE00FF, prControl},              // Cn [128] <reserved-E0080>..<reserved-E00FF>
	{0xE0100, 0xE01EF, prExtend},               // Mn [240] VARIATION SELECTOR-17..VARIATION SELECTOR-256
	{0xE01F0, 0xE0FFF, prControl},              // Cn [3600] <reserved-E01F0>..<reserved-E0FFF>
}
