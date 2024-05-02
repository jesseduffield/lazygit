// Code generated via go generate from gen_properties.go. DO NOT EDIT.

package uniseg

// workBreakCodePoints are taken from
// https://www.unicode.org/Public/15.0.0/ucd/auxiliary/WordBreakProperty.txt
// and
// https://unicode.org/Public/15.0.0/ucd/emoji/emoji-data.txt
// ("Extended_Pictographic" only)
// on September 5, 2023. See https://www.unicode.org/license.html for the Unicode
// license agreement.
var workBreakCodePoints = [][3]int{
	{0x000A, 0x000A, prLF},                     // Cc       <control-000A>
	{0x000B, 0x000C, prNewline},                // Cc   [2] <control-000B>..<control-000C>
	{0x000D, 0x000D, prCR},                     // Cc       <control-000D>
	{0x0020, 0x0020, prWSegSpace},              // Zs       SPACE
	{0x0022, 0x0022, prDoubleQuote},            // Po       QUOTATION MARK
	{0x0027, 0x0027, prSingleQuote},            // Po       APOSTROPHE
	{0x002C, 0x002C, prMidNum},                 // Po       COMMA
	{0x002E, 0x002E, prMidNumLet},              // Po       FULL STOP
	{0x0030, 0x0039, prNumeric},                // Nd  [10] DIGIT ZERO..DIGIT NINE
	{0x003A, 0x003A, prMidLetter},              // Po       COLON
	{0x003B, 0x003B, prMidNum},                 // Po       SEMICOLON
	{0x0041, 0x005A, prALetter},                // L&  [26] LATIN CAPITAL LETTER A..LATIN CAPITAL LETTER Z
	{0x005F, 0x005F, prExtendNumLet},           // Pc       LOW LINE
	{0x0061, 0x007A, prALetter},                // L&  [26] LATIN SMALL LETTER A..LATIN SMALL LETTER Z
	{0x0085, 0x0085, prNewline},                // Cc       <control-0085>
	{0x00A9, 0x00A9, prExtendedPictographic},   // E0.6   [1] (©️)       copyright
	{0x00AA, 0x00AA, prALetter},                // Lo       FEMININE ORDINAL INDICATOR
	{0x00AD, 0x00AD, prFormat},                 // Cf       SOFT HYPHEN
	{0x00AE, 0x00AE, prExtendedPictographic},   // E0.6   [1] (®️)       registered
	{0x00B5, 0x00B5, prALetter},                // L&       MICRO SIGN
	{0x00B7, 0x00B7, prMidLetter},              // Po       MIDDLE DOT
	{0x00BA, 0x00BA, prALetter},                // Lo       MASCULINE ORDINAL INDICATOR
	{0x00C0, 0x00D6, prALetter},                // L&  [23] LATIN CAPITAL LETTER A WITH GRAVE..LATIN CAPITAL LETTER O WITH DIAERESIS
	{0x00D8, 0x00F6, prALetter},                // L&  [31] LATIN CAPITAL LETTER O WITH STROKE..LATIN SMALL LETTER O WITH DIAERESIS
	{0x00F8, 0x01BA, prALetter},                // L& [195] LATIN SMALL LETTER O WITH STROKE..LATIN SMALL LETTER EZH WITH TAIL
	{0x01BB, 0x01BB, prALetter},                // Lo       LATIN LETTER TWO WITH STROKE
	{0x01BC, 0x01BF, prALetter},                // L&   [4] LATIN CAPITAL LETTER TONE FIVE..LATIN LETTER WYNN
	{0x01C0, 0x01C3, prALetter},                // Lo   [4] LATIN LETTER DENTAL CLICK..LATIN LETTER RETROFLEX CLICK
	{0x01C4, 0x0293, prALetter},                // L& [208] LATIN CAPITAL LETTER DZ WITH CARON..LATIN SMALL LETTER EZH WITH CURL
	{0x0294, 0x0294, prALetter},                // Lo       LATIN LETTER GLOTTAL STOP
	{0x0295, 0x02AF, prALetter},                // L&  [27] LATIN LETTER PHARYNGEAL VOICED FRICATIVE..LATIN SMALL LETTER TURNED H WITH FISHHOOK AND TAIL
	{0x02B0, 0x02C1, prALetter},                // Lm  [18] MODIFIER LETTER SMALL H..MODIFIER LETTER REVERSED GLOTTAL STOP
	{0x02C2, 0x02C5, prALetter},                // Sk   [4] MODIFIER LETTER LEFT ARROWHEAD..MODIFIER LETTER DOWN ARROWHEAD
	{0x02C6, 0x02D1, prALetter},                // Lm  [12] MODIFIER LETTER CIRCUMFLEX ACCENT..MODIFIER LETTER HALF TRIANGULAR COLON
	{0x02D2, 0x02D7, prALetter},                // Sk   [6] MODIFIER LETTER CENTRED RIGHT HALF RING..MODIFIER LETTER MINUS SIGN
	{0x02DE, 0x02DF, prALetter},                // Sk   [2] MODIFIER LETTER RHOTIC HOOK..MODIFIER LETTER CROSS ACCENT
	{0x02E0, 0x02E4, prALetter},                // Lm   [5] MODIFIER LETTER SMALL GAMMA..MODIFIER LETTER SMALL REVERSED GLOTTAL STOP
	{0x02E5, 0x02EB, prALetter},                // Sk   [7] MODIFIER LETTER EXTRA-HIGH TONE BAR..MODIFIER LETTER YANG DEPARTING TONE MARK
	{0x02EC, 0x02EC, prALetter},                // Lm       MODIFIER LETTER VOICING
	{0x02ED, 0x02ED, prALetter},                // Sk       MODIFIER LETTER UNASPIRATED
	{0x02EE, 0x02EE, prALetter},                // Lm       MODIFIER LETTER DOUBLE APOSTROPHE
	{0x02EF, 0x02FF, prALetter},                // Sk  [17] MODIFIER LETTER LOW DOWN ARROWHEAD..MODIFIER LETTER LOW LEFT ARROW
	{0x0300, 0x036F, prExtend},                 // Mn [112] COMBINING GRAVE ACCENT..COMBINING LATIN SMALL LETTER X
	{0x0370, 0x0373, prALetter},                // L&   [4] GREEK CAPITAL LETTER HETA..GREEK SMALL LETTER ARCHAIC SAMPI
	{0x0374, 0x0374, prALetter},                // Lm       GREEK NUMERAL SIGN
	{0x0376, 0x0377, prALetter},                // L&   [2] GREEK CAPITAL LETTER PAMPHYLIAN DIGAMMA..GREEK SMALL LETTER PAMPHYLIAN DIGAMMA
	{0x037A, 0x037A, prALetter},                // Lm       GREEK YPOGEGRAMMENI
	{0x037B, 0x037D, prALetter},                // L&   [3] GREEK SMALL REVERSED LUNATE SIGMA SYMBOL..GREEK SMALL REVERSED DOTTED LUNATE SIGMA SYMBOL
	{0x037E, 0x037E, prMidNum},                 // Po       GREEK QUESTION MARK
	{0x037F, 0x037F, prALetter},                // L&       GREEK CAPITAL LETTER YOT
	{0x0386, 0x0386, prALetter},                // L&       GREEK CAPITAL LETTER ALPHA WITH TONOS
	{0x0387, 0x0387, prMidLetter},              // Po       GREEK ANO TELEIA
	{0x0388, 0x038A, prALetter},                // L&   [3] GREEK CAPITAL LETTER EPSILON WITH TONOS..GREEK CAPITAL LETTER IOTA WITH TONOS
	{0x038C, 0x038C, prALetter},                // L&       GREEK CAPITAL LETTER OMICRON WITH TONOS
	{0x038E, 0x03A1, prALetter},                // L&  [20] GREEK CAPITAL LETTER UPSILON WITH TONOS..GREEK CAPITAL LETTER RHO
	{0x03A3, 0x03F5, prALetter},                // L&  [83] GREEK CAPITAL LETTER SIGMA..GREEK LUNATE EPSILON SYMBOL
	{0x03F7, 0x0481, prALetter},                // L& [139] GREEK CAPITAL LETTER SHO..CYRILLIC SMALL LETTER KOPPA
	{0x0483, 0x0487, prExtend},                 // Mn   [5] COMBINING CYRILLIC TITLO..COMBINING CYRILLIC POKRYTIE
	{0x0488, 0x0489, prExtend},                 // Me   [2] COMBINING CYRILLIC HUNDRED THOUSANDS SIGN..COMBINING CYRILLIC MILLIONS SIGN
	{0x048A, 0x052F, prALetter},                // L& [166] CYRILLIC CAPITAL LETTER SHORT I WITH TAIL..CYRILLIC SMALL LETTER EL WITH DESCENDER
	{0x0531, 0x0556, prALetter},                // L&  [38] ARMENIAN CAPITAL LETTER AYB..ARMENIAN CAPITAL LETTER FEH
	{0x0559, 0x0559, prALetter},                // Lm       ARMENIAN MODIFIER LETTER LEFT HALF RING
	{0x055A, 0x055C, prALetter},                // Po   [3] ARMENIAN APOSTROPHE..ARMENIAN EXCLAMATION MARK
	{0x055E, 0x055E, prALetter},                // Po       ARMENIAN QUESTION MARK
	{0x055F, 0x055F, prMidLetter},              // Po       ARMENIAN ABBREVIATION MARK
	{0x0560, 0x0588, prALetter},                // L&  [41] ARMENIAN SMALL LETTER TURNED AYB..ARMENIAN SMALL LETTER YI WITH STROKE
	{0x0589, 0x0589, prMidNum},                 // Po       ARMENIAN FULL STOP
	{0x058A, 0x058A, prALetter},                // Pd       ARMENIAN HYPHEN
	{0x0591, 0x05BD, prExtend},                 // Mn  [45] HEBREW ACCENT ETNAHTA..HEBREW POINT METEG
	{0x05BF, 0x05BF, prExtend},                 // Mn       HEBREW POINT RAFE
	{0x05C1, 0x05C2, prExtend},                 // Mn   [2] HEBREW POINT SHIN DOT..HEBREW POINT SIN DOT
	{0x05C4, 0x05C5, prExtend},                 // Mn   [2] HEBREW MARK UPPER DOT..HEBREW MARK LOWER DOT
	{0x05C7, 0x05C7, prExtend},                 // Mn       HEBREW POINT QAMATS QATAN
	{0x05D0, 0x05EA, prHebrewLetter},           // Lo  [27] HEBREW LETTER ALEF..HEBREW LETTER TAV
	{0x05EF, 0x05F2, prHebrewLetter},           // Lo   [4] HEBREW YOD TRIANGLE..HEBREW LIGATURE YIDDISH DOUBLE YOD
	{0x05F3, 0x05F3, prALetter},                // Po       HEBREW PUNCTUATION GERESH
	{0x05F4, 0x05F4, prMidLetter},              // Po       HEBREW PUNCTUATION GERSHAYIM
	{0x0600, 0x0605, prFormat},                 // Cf   [6] ARABIC NUMBER SIGN..ARABIC NUMBER MARK ABOVE
	{0x060C, 0x060D, prMidNum},                 // Po   [2] ARABIC COMMA..ARABIC DATE SEPARATOR
	{0x0610, 0x061A, prExtend},                 // Mn  [11] ARABIC SIGN SALLALLAHOU ALAYHE WASSALLAM..ARABIC SMALL KASRA
	{0x061C, 0x061C, prFormat},                 // Cf       ARABIC LETTER MARK
	{0x0620, 0x063F, prALetter},                // Lo  [32] ARABIC LETTER KASHMIRI YEH..ARABIC LETTER FARSI YEH WITH THREE DOTS ABOVE
	{0x0640, 0x0640, prALetter},                // Lm       ARABIC TATWEEL
	{0x0641, 0x064A, prALetter},                // Lo  [10] ARABIC LETTER FEH..ARABIC LETTER YEH
	{0x064B, 0x065F, prExtend},                 // Mn  [21] ARABIC FATHATAN..ARABIC WAVY HAMZA BELOW
	{0x0660, 0x0669, prNumeric},                // Nd  [10] ARABIC-INDIC DIGIT ZERO..ARABIC-INDIC DIGIT NINE
	{0x066B, 0x066B, prNumeric},                // Po       ARABIC DECIMAL SEPARATOR
	{0x066C, 0x066C, prMidNum},                 // Po       ARABIC THOUSANDS SEPARATOR
	{0x066E, 0x066F, prALetter},                // Lo   [2] ARABIC LETTER DOTLESS BEH..ARABIC LETTER DOTLESS QAF
	{0x0670, 0x0670, prExtend},                 // Mn       ARABIC LETTER SUPERSCRIPT ALEF
	{0x0671, 0x06D3, prALetter},                // Lo  [99] ARABIC LETTER ALEF WASLA..ARABIC LETTER YEH BARREE WITH HAMZA ABOVE
	{0x06D5, 0x06D5, prALetter},                // Lo       ARABIC LETTER AE
	{0x06D6, 0x06DC, prExtend},                 // Mn   [7] ARABIC SMALL HIGH LIGATURE SAD WITH LAM WITH ALEF MAKSURA..ARABIC SMALL HIGH SEEN
	{0x06DD, 0x06DD, prFormat},                 // Cf       ARABIC END OF AYAH
	{0x06DF, 0x06E4, prExtend},                 // Mn   [6] ARABIC SMALL HIGH ROUNDED ZERO..ARABIC SMALL HIGH MADDA
	{0x06E5, 0x06E6, prALetter},                // Lm   [2] ARABIC SMALL WAW..ARABIC SMALL YEH
	{0x06E7, 0x06E8, prExtend},                 // Mn   [2] ARABIC SMALL HIGH YEH..ARABIC SMALL HIGH NOON
	{0x06EA, 0x06ED, prExtend},                 // Mn   [4] ARABIC EMPTY CENTRE LOW STOP..ARABIC SMALL LOW MEEM
	{0x06EE, 0x06EF, prALetter},                // Lo   [2] ARABIC LETTER DAL WITH INVERTED V..ARABIC LETTER REH WITH INVERTED V
	{0x06F0, 0x06F9, prNumeric},                // Nd  [10] EXTENDED ARABIC-INDIC DIGIT ZERO..EXTENDED ARABIC-INDIC DIGIT NINE
	{0x06FA, 0x06FC, prALetter},                // Lo   [3] ARABIC LETTER SHEEN WITH DOT BELOW..ARABIC LETTER GHAIN WITH DOT BELOW
	{0x06FF, 0x06FF, prALetter},                // Lo       ARABIC LETTER HEH WITH INVERTED V
	{0x070F, 0x070F, prFormat},                 // Cf       SYRIAC ABBREVIATION MARK
	{0x0710, 0x0710, prALetter},                // Lo       SYRIAC LETTER ALAPH
	{0x0711, 0x0711, prExtend},                 // Mn       SYRIAC LETTER SUPERSCRIPT ALAPH
	{0x0712, 0x072F, prALetter},                // Lo  [30] SYRIAC LETTER BETH..SYRIAC LETTER PERSIAN DHALATH
	{0x0730, 0x074A, prExtend},                 // Mn  [27] SYRIAC PTHAHA ABOVE..SYRIAC BARREKH
	{0x074D, 0x07A5, prALetter},                // Lo  [89] SYRIAC LETTER SOGDIAN ZHAIN..THAANA LETTER WAAVU
	{0x07A6, 0x07B0, prExtend},                 // Mn  [11] THAANA ABAFILI..THAANA SUKUN
	{0x07B1, 0x07B1, prALetter},                // Lo       THAANA LETTER NAA
	{0x07C0, 0x07C9, prNumeric},                // Nd  [10] NKO DIGIT ZERO..NKO DIGIT NINE
	{0x07CA, 0x07EA, prALetter},                // Lo  [33] NKO LETTER A..NKO LETTER JONA RA
	{0x07EB, 0x07F3, prExtend},                 // Mn   [9] NKO COMBINING SHORT HIGH TONE..NKO COMBINING DOUBLE DOT ABOVE
	{0x07F4, 0x07F5, prALetter},                // Lm   [2] NKO HIGH TONE APOSTROPHE..NKO LOW TONE APOSTROPHE
	{0x07F8, 0x07F8, prMidNum},                 // Po       NKO COMMA
	{0x07FA, 0x07FA, prALetter},                // Lm       NKO LAJANYALAN
	{0x07FD, 0x07FD, prExtend},                 // Mn       NKO DANTAYALAN
	{0x0800, 0x0815, prALetter},                // Lo  [22] SAMARITAN LETTER ALAF..SAMARITAN LETTER TAAF
	{0x0816, 0x0819, prExtend},                 // Mn   [4] SAMARITAN MARK IN..SAMARITAN MARK DAGESH
	{0x081A, 0x081A, prALetter},                // Lm       SAMARITAN MODIFIER LETTER EPENTHETIC YUT
	{0x081B, 0x0823, prExtend},                 // Mn   [9] SAMARITAN MARK EPENTHETIC YUT..SAMARITAN VOWEL SIGN A
	{0x0824, 0x0824, prALetter},                // Lm       SAMARITAN MODIFIER LETTER SHORT A
	{0x0825, 0x0827, prExtend},                 // Mn   [3] SAMARITAN VOWEL SIGN SHORT A..SAMARITAN VOWEL SIGN U
	{0x0828, 0x0828, prALetter},                // Lm       SAMARITAN MODIFIER LETTER I
	{0x0829, 0x082D, prExtend},                 // Mn   [5] SAMARITAN VOWEL SIGN LONG I..SAMARITAN MARK NEQUDAA
	{0x0840, 0x0858, prALetter},                // Lo  [25] MANDAIC LETTER HALQA..MANDAIC LETTER AIN
	{0x0859, 0x085B, prExtend},                 // Mn   [3] MANDAIC AFFRICATION MARK..MANDAIC GEMINATION MARK
	{0x0860, 0x086A, prALetter},                // Lo  [11] SYRIAC LETTER MALAYALAM NGA..SYRIAC LETTER MALAYALAM SSA
	{0x0870, 0x0887, prALetter},                // Lo  [24] ARABIC LETTER ALEF WITH ATTACHED FATHA..ARABIC BASELINE ROUND DOT
	{0x0889, 0x088E, prALetter},                // Lo   [6] ARABIC LETTER NOON WITH INVERTED SMALL V..ARABIC VERTICAL TAIL
	{0x0890, 0x0891, prFormat},                 // Cf   [2] ARABIC POUND MARK ABOVE..ARABIC PIASTRE MARK ABOVE
	{0x0898, 0x089F, prExtend},                 // Mn   [8] ARABIC SMALL HIGH WORD AL-JUZ..ARABIC HALF MADDA OVER MADDA
	{0x08A0, 0x08C8, prALetter},                // Lo  [41] ARABIC LETTER BEH WITH SMALL V BELOW..ARABIC LETTER GRAF
	{0x08C9, 0x08C9, prALetter},                // Lm       ARABIC SMALL FARSI YEH
	{0x08CA, 0x08E1, prExtend},                 // Mn  [24] ARABIC SMALL HIGH FARSI YEH..ARABIC SMALL HIGH SIGN SAFHA
	{0x08E2, 0x08E2, prFormat},                 // Cf       ARABIC DISPUTED END OF AYAH
	{0x08E3, 0x0902, prExtend},                 // Mn  [32] ARABIC TURNED DAMMA BELOW..DEVANAGARI SIGN ANUSVARA
	{0x0903, 0x0903, prExtend},                 // Mc       DEVANAGARI SIGN VISARGA
	{0x0904, 0x0939, prALetter},                // Lo  [54] DEVANAGARI LETTER SHORT A..DEVANAGARI LETTER HA
	{0x093A, 0x093A, prExtend},                 // Mn       DEVANAGARI VOWEL SIGN OE
	{0x093B, 0x093B, prExtend},                 // Mc       DEVANAGARI VOWEL SIGN OOE
	{0x093C, 0x093C, prExtend},                 // Mn       DEVANAGARI SIGN NUKTA
	{0x093D, 0x093D, prALetter},                // Lo       DEVANAGARI SIGN AVAGRAHA
	{0x093E, 0x0940, prExtend},                 // Mc   [3] DEVANAGARI VOWEL SIGN AA..DEVANAGARI VOWEL SIGN II
	{0x0941, 0x0948, prExtend},                 // Mn   [8] DEVANAGARI VOWEL SIGN U..DEVANAGARI VOWEL SIGN AI
	{0x0949, 0x094C, prExtend},                 // Mc   [4] DEVANAGARI VOWEL SIGN CANDRA O..DEVANAGARI VOWEL SIGN AU
	{0x094D, 0x094D, prExtend},                 // Mn       DEVANAGARI SIGN VIRAMA
	{0x094E, 0x094F, prExtend},                 // Mc   [2] DEVANAGARI VOWEL SIGN PRISHTHAMATRA E..DEVANAGARI VOWEL SIGN AW
	{0x0950, 0x0950, prALetter},                // Lo       DEVANAGARI OM
	{0x0951, 0x0957, prExtend},                 // Mn   [7] DEVANAGARI STRESS SIGN UDATTA..DEVANAGARI VOWEL SIGN UUE
	{0x0958, 0x0961, prALetter},                // Lo  [10] DEVANAGARI LETTER QA..DEVANAGARI LETTER VOCALIC LL
	{0x0962, 0x0963, prExtend},                 // Mn   [2] DEVANAGARI VOWEL SIGN VOCALIC L..DEVANAGARI VOWEL SIGN VOCALIC LL
	{0x0966, 0x096F, prNumeric},                // Nd  [10] DEVANAGARI DIGIT ZERO..DEVANAGARI DIGIT NINE
	{0x0971, 0x0971, prALetter},                // Lm       DEVANAGARI SIGN HIGH SPACING DOT
	{0x0972, 0x0980, prALetter},                // Lo  [15] DEVANAGARI LETTER CANDRA A..BENGALI ANJI
	{0x0981, 0x0981, prExtend},                 // Mn       BENGALI SIGN CANDRABINDU
	{0x0982, 0x0983, prExtend},                 // Mc   [2] BENGALI SIGN ANUSVARA..BENGALI SIGN VISARGA
	{0x0985, 0x098C, prALetter},                // Lo   [8] BENGALI LETTER A..BENGALI LETTER VOCALIC L
	{0x098F, 0x0990, prALetter},                // Lo   [2] BENGALI LETTER E..BENGALI LETTER AI
	{0x0993, 0x09A8, prALetter},                // Lo  [22] BENGALI LETTER O..BENGALI LETTER NA
	{0x09AA, 0x09B0, prALetter},                // Lo   [7] BENGALI LETTER PA..BENGALI LETTER RA
	{0x09B2, 0x09B2, prALetter},                // Lo       BENGALI LETTER LA
	{0x09B6, 0x09B9, prALetter},                // Lo   [4] BENGALI LETTER SHA..BENGALI LETTER HA
	{0x09BC, 0x09BC, prExtend},                 // Mn       BENGALI SIGN NUKTA
	{0x09BD, 0x09BD, prALetter},                // Lo       BENGALI SIGN AVAGRAHA
	{0x09BE, 0x09C0, prExtend},                 // Mc   [3] BENGALI VOWEL SIGN AA..BENGALI VOWEL SIGN II
	{0x09C1, 0x09C4, prExtend},                 // Mn   [4] BENGALI VOWEL SIGN U..BENGALI VOWEL SIGN VOCALIC RR
	{0x09C7, 0x09C8, prExtend},                 // Mc   [2] BENGALI VOWEL SIGN E..BENGALI VOWEL SIGN AI
	{0x09CB, 0x09CC, prExtend},                 // Mc   [2] BENGALI VOWEL SIGN O..BENGALI VOWEL SIGN AU
	{0x09CD, 0x09CD, prExtend},                 // Mn       BENGALI SIGN VIRAMA
	{0x09CE, 0x09CE, prALetter},                // Lo       BENGALI LETTER KHANDA TA
	{0x09D7, 0x09D7, prExtend},                 // Mc       BENGALI AU LENGTH MARK
	{0x09DC, 0x09DD, prALetter},                // Lo   [2] BENGALI LETTER RRA..BENGALI LETTER RHA
	{0x09DF, 0x09E1, prALetter},                // Lo   [3] BENGALI LETTER YYA..BENGALI LETTER VOCALIC LL
	{0x09E2, 0x09E3, prExtend},                 // Mn   [2] BENGALI VOWEL SIGN VOCALIC L..BENGALI VOWEL SIGN VOCALIC LL
	{0x09E6, 0x09EF, prNumeric},                // Nd  [10] BENGALI DIGIT ZERO..BENGALI DIGIT NINE
	{0x09F0, 0x09F1, prALetter},                // Lo   [2] BENGALI LETTER RA WITH MIDDLE DIAGONAL..BENGALI LETTER RA WITH LOWER DIAGONAL
	{0x09FC, 0x09FC, prALetter},                // Lo       BENGALI LETTER VEDIC ANUSVARA
	{0x09FE, 0x09FE, prExtend},                 // Mn       BENGALI SANDHI MARK
	{0x0A01, 0x0A02, prExtend},                 // Mn   [2] GURMUKHI SIGN ADAK BINDI..GURMUKHI SIGN BINDI
	{0x0A03, 0x0A03, prExtend},                 // Mc       GURMUKHI SIGN VISARGA
	{0x0A05, 0x0A0A, prALetter},                // Lo   [6] GURMUKHI LETTER A..GURMUKHI LETTER UU
	{0x0A0F, 0x0A10, prALetter},                // Lo   [2] GURMUKHI LETTER EE..GURMUKHI LETTER AI
	{0x0A13, 0x0A28, prALetter},                // Lo  [22] GURMUKHI LETTER OO..GURMUKHI LETTER NA
	{0x0A2A, 0x0A30, prALetter},                // Lo   [7] GURMUKHI LETTER PA..GURMUKHI LETTER RA
	{0x0A32, 0x0A33, prALetter},                // Lo   [2] GURMUKHI LETTER LA..GURMUKHI LETTER LLA
	{0x0A35, 0x0A36, prALetter},                // Lo   [2] GURMUKHI LETTER VA..GURMUKHI LETTER SHA
	{0x0A38, 0x0A39, prALetter},                // Lo   [2] GURMUKHI LETTER SA..GURMUKHI LETTER HA
	{0x0A3C, 0x0A3C, prExtend},                 // Mn       GURMUKHI SIGN NUKTA
	{0x0A3E, 0x0A40, prExtend},                 // Mc   [3] GURMUKHI VOWEL SIGN AA..GURMUKHI VOWEL SIGN II
	{0x0A41, 0x0A42, prExtend},                 // Mn   [2] GURMUKHI VOWEL SIGN U..GURMUKHI VOWEL SIGN UU
	{0x0A47, 0x0A48, prExtend},                 // Mn   [2] GURMUKHI VOWEL SIGN EE..GURMUKHI VOWEL SIGN AI
	{0x0A4B, 0x0A4D, prExtend},                 // Mn   [3] GURMUKHI VOWEL SIGN OO..GURMUKHI SIGN VIRAMA
	{0x0A51, 0x0A51, prExtend},                 // Mn       GURMUKHI SIGN UDAAT
	{0x0A59, 0x0A5C, prALetter},                // Lo   [4] GURMUKHI LETTER KHHA..GURMUKHI LETTER RRA
	{0x0A5E, 0x0A5E, prALetter},                // Lo       GURMUKHI LETTER FA
	{0x0A66, 0x0A6F, prNumeric},                // Nd  [10] GURMUKHI DIGIT ZERO..GURMUKHI DIGIT NINE
	{0x0A70, 0x0A71, prExtend},                 // Mn   [2] GURMUKHI TIPPI..GURMUKHI ADDAK
	{0x0A72, 0x0A74, prALetter},                // Lo   [3] GURMUKHI IRI..GURMUKHI EK ONKAR
	{0x0A75, 0x0A75, prExtend},                 // Mn       GURMUKHI SIGN YAKASH
	{0x0A81, 0x0A82, prExtend},                 // Mn   [2] GUJARATI SIGN CANDRABINDU..GUJARATI SIGN ANUSVARA
	{0x0A83, 0x0A83, prExtend},                 // Mc       GUJARATI SIGN VISARGA
	{0x0A85, 0x0A8D, prALetter},                // Lo   [9] GUJARATI LETTER A..GUJARATI VOWEL CANDRA E
	{0x0A8F, 0x0A91, prALetter},                // Lo   [3] GUJARATI LETTER E..GUJARATI VOWEL CANDRA O
	{0x0A93, 0x0AA8, prALetter},                // Lo  [22] GUJARATI LETTER O..GUJARATI LETTER NA
	{0x0AAA, 0x0AB0, prALetter},                // Lo   [7] GUJARATI LETTER PA..GUJARATI LETTER RA
	{0x0AB2, 0x0AB3, prALetter},                // Lo   [2] GUJARATI LETTER LA..GUJARATI LETTER LLA
	{0x0AB5, 0x0AB9, prALetter},                // Lo   [5] GUJARATI LETTER VA..GUJARATI LETTER HA
	{0x0ABC, 0x0ABC, prExtend},                 // Mn       GUJARATI SIGN NUKTA
	{0x0ABD, 0x0ABD, prALetter},                // Lo       GUJARATI SIGN AVAGRAHA
	{0x0ABE, 0x0AC0, prExtend},                 // Mc   [3] GUJARATI VOWEL SIGN AA..GUJARATI VOWEL SIGN II
	{0x0AC1, 0x0AC5, prExtend},                 // Mn   [5] GUJARATI VOWEL SIGN U..GUJARATI VOWEL SIGN CANDRA E
	{0x0AC7, 0x0AC8, prExtend},                 // Mn   [2] GUJARATI VOWEL SIGN E..GUJARATI VOWEL SIGN AI
	{0x0AC9, 0x0AC9, prExtend},                 // Mc       GUJARATI VOWEL SIGN CANDRA O
	{0x0ACB, 0x0ACC, prExtend},                 // Mc   [2] GUJARATI VOWEL SIGN O..GUJARATI VOWEL SIGN AU
	{0x0ACD, 0x0ACD, prExtend},                 // Mn       GUJARATI SIGN VIRAMA
	{0x0AD0, 0x0AD0, prALetter},                // Lo       GUJARATI OM
	{0x0AE0, 0x0AE1, prALetter},                // Lo   [2] GUJARATI LETTER VOCALIC RR..GUJARATI LETTER VOCALIC LL
	{0x0AE2, 0x0AE3, prExtend},                 // Mn   [2] GUJARATI VOWEL SIGN VOCALIC L..GUJARATI VOWEL SIGN VOCALIC LL
	{0x0AE6, 0x0AEF, prNumeric},                // Nd  [10] GUJARATI DIGIT ZERO..GUJARATI DIGIT NINE
	{0x0AF9, 0x0AF9, prALetter},                // Lo       GUJARATI LETTER ZHA
	{0x0AFA, 0x0AFF, prExtend},                 // Mn   [6] GUJARATI SIGN SUKUN..GUJARATI SIGN TWO-CIRCLE NUKTA ABOVE
	{0x0B01, 0x0B01, prExtend},                 // Mn       ORIYA SIGN CANDRABINDU
	{0x0B02, 0x0B03, prExtend},                 // Mc   [2] ORIYA SIGN ANUSVARA..ORIYA SIGN VISARGA
	{0x0B05, 0x0B0C, prALetter},                // Lo   [8] ORIYA LETTER A..ORIYA LETTER VOCALIC L
	{0x0B0F, 0x0B10, prALetter},                // Lo   [2] ORIYA LETTER E..ORIYA LETTER AI
	{0x0B13, 0x0B28, prALetter},                // Lo  [22] ORIYA LETTER O..ORIYA LETTER NA
	{0x0B2A, 0x0B30, prALetter},                // Lo   [7] ORIYA LETTER PA..ORIYA LETTER RA
	{0x0B32, 0x0B33, prALetter},                // Lo   [2] ORIYA LETTER LA..ORIYA LETTER LLA
	{0x0B35, 0x0B39, prALetter},                // Lo   [5] ORIYA LETTER VA..ORIYA LETTER HA
	{0x0B3C, 0x0B3C, prExtend},                 // Mn       ORIYA SIGN NUKTA
	{0x0B3D, 0x0B3D, prALetter},                // Lo       ORIYA SIGN AVAGRAHA
	{0x0B3E, 0x0B3E, prExtend},                 // Mc       ORIYA VOWEL SIGN AA
	{0x0B3F, 0x0B3F, prExtend},                 // Mn       ORIYA VOWEL SIGN I
	{0x0B40, 0x0B40, prExtend},                 // Mc       ORIYA VOWEL SIGN II
	{0x0B41, 0x0B44, prExtend},                 // Mn   [4] ORIYA VOWEL SIGN U..ORIYA VOWEL SIGN VOCALIC RR
	{0x0B47, 0x0B48, prExtend},                 // Mc   [2] ORIYA VOWEL SIGN E..ORIYA VOWEL SIGN AI
	{0x0B4B, 0x0B4C, prExtend},                 // Mc   [2] ORIYA VOWEL SIGN O..ORIYA VOWEL SIGN AU
	{0x0B4D, 0x0B4D, prExtend},                 // Mn       ORIYA SIGN VIRAMA
	{0x0B55, 0x0B56, prExtend},                 // Mn   [2] ORIYA SIGN OVERLINE..ORIYA AI LENGTH MARK
	{0x0B57, 0x0B57, prExtend},                 // Mc       ORIYA AU LENGTH MARK
	{0x0B5C, 0x0B5D, prALetter},                // Lo   [2] ORIYA LETTER RRA..ORIYA LETTER RHA
	{0x0B5F, 0x0B61, prALetter},                // Lo   [3] ORIYA LETTER YYA..ORIYA LETTER VOCALIC LL
	{0x0B62, 0x0B63, prExtend},                 // Mn   [2] ORIYA VOWEL SIGN VOCALIC L..ORIYA VOWEL SIGN VOCALIC LL
	{0x0B66, 0x0B6F, prNumeric},                // Nd  [10] ORIYA DIGIT ZERO..ORIYA DIGIT NINE
	{0x0B71, 0x0B71, prALetter},                // Lo       ORIYA LETTER WA
	{0x0B82, 0x0B82, prExtend},                 // Mn       TAMIL SIGN ANUSVARA
	{0x0B83, 0x0B83, prALetter},                // Lo       TAMIL SIGN VISARGA
	{0x0B85, 0x0B8A, prALetter},                // Lo   [6] TAMIL LETTER A..TAMIL LETTER UU
	{0x0B8E, 0x0B90, prALetter},                // Lo   [3] TAMIL LETTER E..TAMIL LETTER AI
	{0x0B92, 0x0B95, prALetter},                // Lo   [4] TAMIL LETTER O..TAMIL LETTER KA
	{0x0B99, 0x0B9A, prALetter},                // Lo   [2] TAMIL LETTER NGA..TAMIL LETTER CA
	{0x0B9C, 0x0B9C, prALetter},                // Lo       TAMIL LETTER JA
	{0x0B9E, 0x0B9F, prALetter},                // Lo   [2] TAMIL LETTER NYA..TAMIL LETTER TTA
	{0x0BA3, 0x0BA4, prALetter},                // Lo   [2] TAMIL LETTER NNA..TAMIL LETTER TA
	{0x0BA8, 0x0BAA, prALetter},                // Lo   [3] TAMIL LETTER NA..TAMIL LETTER PA
	{0x0BAE, 0x0BB9, prALetter},                // Lo  [12] TAMIL LETTER MA..TAMIL LETTER HA
	{0x0BBE, 0x0BBF, prExtend},                 // Mc   [2] TAMIL VOWEL SIGN AA..TAMIL VOWEL SIGN I
	{0x0BC0, 0x0BC0, prExtend},                 // Mn       TAMIL VOWEL SIGN II
	{0x0BC1, 0x0BC2, prExtend},                 // Mc   [2] TAMIL VOWEL SIGN U..TAMIL VOWEL SIGN UU
	{0x0BC6, 0x0BC8, prExtend},                 // Mc   [3] TAMIL VOWEL SIGN E..TAMIL VOWEL SIGN AI
	{0x0BCA, 0x0BCC, prExtend},                 // Mc   [3] TAMIL VOWEL SIGN O..TAMIL VOWEL SIGN AU
	{0x0BCD, 0x0BCD, prExtend},                 // Mn       TAMIL SIGN VIRAMA
	{0x0BD0, 0x0BD0, prALetter},                // Lo       TAMIL OM
	{0x0BD7, 0x0BD7, prExtend},                 // Mc       TAMIL AU LENGTH MARK
	{0x0BE6, 0x0BEF, prNumeric},                // Nd  [10] TAMIL DIGIT ZERO..TAMIL DIGIT NINE
	{0x0C00, 0x0C00, prExtend},                 // Mn       TELUGU SIGN COMBINING CANDRABINDU ABOVE
	{0x0C01, 0x0C03, prExtend},                 // Mc   [3] TELUGU SIGN CANDRABINDU..TELUGU SIGN VISARGA
	{0x0C04, 0x0C04, prExtend},                 // Mn       TELUGU SIGN COMBINING ANUSVARA ABOVE
	{0x0C05, 0x0C0C, prALetter},                // Lo   [8] TELUGU LETTER A..TELUGU LETTER VOCALIC L
	{0x0C0E, 0x0C10, prALetter},                // Lo   [3] TELUGU LETTER E..TELUGU LETTER AI
	{0x0C12, 0x0C28, prALetter},                // Lo  [23] TELUGU LETTER O..TELUGU LETTER NA
	{0x0C2A, 0x0C39, prALetter},                // Lo  [16] TELUGU LETTER PA..TELUGU LETTER HA
	{0x0C3C, 0x0C3C, prExtend},                 // Mn       TELUGU SIGN NUKTA
	{0x0C3D, 0x0C3D, prALetter},                // Lo       TELUGU SIGN AVAGRAHA
	{0x0C3E, 0x0C40, prExtend},                 // Mn   [3] TELUGU VOWEL SIGN AA..TELUGU VOWEL SIGN II
	{0x0C41, 0x0C44, prExtend},                 // Mc   [4] TELUGU VOWEL SIGN U..TELUGU VOWEL SIGN VOCALIC RR
	{0x0C46, 0x0C48, prExtend},                 // Mn   [3] TELUGU VOWEL SIGN E..TELUGU VOWEL SIGN AI
	{0x0C4A, 0x0C4D, prExtend},                 // Mn   [4] TELUGU VOWEL SIGN O..TELUGU SIGN VIRAMA
	{0x0C55, 0x0C56, prExtend},                 // Mn   [2] TELUGU LENGTH MARK..TELUGU AI LENGTH MARK
	{0x0C58, 0x0C5A, prALetter},                // Lo   [3] TELUGU LETTER TSA..TELUGU LETTER RRRA
	{0x0C5D, 0x0C5D, prALetter},                // Lo       TELUGU LETTER NAKAARA POLLU
	{0x0C60, 0x0C61, prALetter},                // Lo   [2] TELUGU LETTER VOCALIC RR..TELUGU LETTER VOCALIC LL
	{0x0C62, 0x0C63, prExtend},                 // Mn   [2] TELUGU VOWEL SIGN VOCALIC L..TELUGU VOWEL SIGN VOCALIC LL
	{0x0C66, 0x0C6F, prNumeric},                // Nd  [10] TELUGU DIGIT ZERO..TELUGU DIGIT NINE
	{0x0C80, 0x0C80, prALetter},                // Lo       KANNADA SIGN SPACING CANDRABINDU
	{0x0C81, 0x0C81, prExtend},                 // Mn       KANNADA SIGN CANDRABINDU
	{0x0C82, 0x0C83, prExtend},                 // Mc   [2] KANNADA SIGN ANUSVARA..KANNADA SIGN VISARGA
	{0x0C85, 0x0C8C, prALetter},                // Lo   [8] KANNADA LETTER A..KANNADA LETTER VOCALIC L
	{0x0C8E, 0x0C90, prALetter},                // Lo   [3] KANNADA LETTER E..KANNADA LETTER AI
	{0x0C92, 0x0CA8, prALetter},                // Lo  [23] KANNADA LETTER O..KANNADA LETTER NA
	{0x0CAA, 0x0CB3, prALetter},                // Lo  [10] KANNADA LETTER PA..KANNADA LETTER LLA
	{0x0CB5, 0x0CB9, prALetter},                // Lo   [5] KANNADA LETTER VA..KANNADA LETTER HA
	{0x0CBC, 0x0CBC, prExtend},                 // Mn       KANNADA SIGN NUKTA
	{0x0CBD, 0x0CBD, prALetter},                // Lo       KANNADA SIGN AVAGRAHA
	{0x0CBE, 0x0CBE, prExtend},                 // Mc       KANNADA VOWEL SIGN AA
	{0x0CBF, 0x0CBF, prExtend},                 // Mn       KANNADA VOWEL SIGN I
	{0x0CC0, 0x0CC4, prExtend},                 // Mc   [5] KANNADA VOWEL SIGN II..KANNADA VOWEL SIGN VOCALIC RR
	{0x0CC6, 0x0CC6, prExtend},                 // Mn       KANNADA VOWEL SIGN E
	{0x0CC7, 0x0CC8, prExtend},                 // Mc   [2] KANNADA VOWEL SIGN EE..KANNADA VOWEL SIGN AI
	{0x0CCA, 0x0CCB, prExtend},                 // Mc   [2] KANNADA VOWEL SIGN O..KANNADA VOWEL SIGN OO
	{0x0CCC, 0x0CCD, prExtend},                 // Mn   [2] KANNADA VOWEL SIGN AU..KANNADA SIGN VIRAMA
	{0x0CD5, 0x0CD6, prExtend},                 // Mc   [2] KANNADA LENGTH MARK..KANNADA AI LENGTH MARK
	{0x0CDD, 0x0CDE, prALetter},                // Lo   [2] KANNADA LETTER NAKAARA POLLU..KANNADA LETTER FA
	{0x0CE0, 0x0CE1, prALetter},                // Lo   [2] KANNADA LETTER VOCALIC RR..KANNADA LETTER VOCALIC LL
	{0x0CE2, 0x0CE3, prExtend},                 // Mn   [2] KANNADA VOWEL SIGN VOCALIC L..KANNADA VOWEL SIGN VOCALIC LL
	{0x0CE6, 0x0CEF, prNumeric},                // Nd  [10] KANNADA DIGIT ZERO..KANNADA DIGIT NINE
	{0x0CF1, 0x0CF2, prALetter},                // Lo   [2] KANNADA SIGN JIHVAMULIYA..KANNADA SIGN UPADHMANIYA
	{0x0CF3, 0x0CF3, prExtend},                 // Mc       KANNADA SIGN COMBINING ANUSVARA ABOVE RIGHT
	{0x0D00, 0x0D01, prExtend},                 // Mn   [2] MALAYALAM SIGN COMBINING ANUSVARA ABOVE..MALAYALAM SIGN CANDRABINDU
	{0x0D02, 0x0D03, prExtend},                 // Mc   [2] MALAYALAM SIGN ANUSVARA..MALAYALAM SIGN VISARGA
	{0x0D04, 0x0D0C, prALetter},                // Lo   [9] MALAYALAM LETTER VEDIC ANUSVARA..MALAYALAM LETTER VOCALIC L
	{0x0D0E, 0x0D10, prALetter},                // Lo   [3] MALAYALAM LETTER E..MALAYALAM LETTER AI
	{0x0D12, 0x0D3A, prALetter},                // Lo  [41] MALAYALAM LETTER O..MALAYALAM LETTER TTTA
	{0x0D3B, 0x0D3C, prExtend},                 // Mn   [2] MALAYALAM SIGN VERTICAL BAR VIRAMA..MALAYALAM SIGN CIRCULAR VIRAMA
	{0x0D3D, 0x0D3D, prALetter},                // Lo       MALAYALAM SIGN AVAGRAHA
	{0x0D3E, 0x0D40, prExtend},                 // Mc   [3] MALAYALAM VOWEL SIGN AA..MALAYALAM VOWEL SIGN II
	{0x0D41, 0x0D44, prExtend},                 // Mn   [4] MALAYALAM VOWEL SIGN U..MALAYALAM VOWEL SIGN VOCALIC RR
	{0x0D46, 0x0D48, prExtend},                 // Mc   [3] MALAYALAM VOWEL SIGN E..MALAYALAM VOWEL SIGN AI
	{0x0D4A, 0x0D4C, prExtend},                 // Mc   [3] MALAYALAM VOWEL SIGN O..MALAYALAM VOWEL SIGN AU
	{0x0D4D, 0x0D4D, prExtend},                 // Mn       MALAYALAM SIGN VIRAMA
	{0x0D4E, 0x0D4E, prALetter},                // Lo       MALAYALAM LETTER DOT REPH
	{0x0D54, 0x0D56, prALetter},                // Lo   [3] MALAYALAM LETTER CHILLU M..MALAYALAM LETTER CHILLU LLL
	{0x0D57, 0x0D57, prExtend},                 // Mc       MALAYALAM AU LENGTH MARK
	{0x0D5F, 0x0D61, prALetter},                // Lo   [3] MALAYALAM LETTER ARCHAIC II..MALAYALAM LETTER VOCALIC LL
	{0x0D62, 0x0D63, prExtend},                 // Mn   [2] MALAYALAM VOWEL SIGN VOCALIC L..MALAYALAM VOWEL SIGN VOCALIC LL
	{0x0D66, 0x0D6F, prNumeric},                // Nd  [10] MALAYALAM DIGIT ZERO..MALAYALAM DIGIT NINE
	{0x0D7A, 0x0D7F, prALetter},                // Lo   [6] MALAYALAM LETTER CHILLU NN..MALAYALAM LETTER CHILLU K
	{0x0D81, 0x0D81, prExtend},                 // Mn       SINHALA SIGN CANDRABINDU
	{0x0D82, 0x0D83, prExtend},                 // Mc   [2] SINHALA SIGN ANUSVARAYA..SINHALA SIGN VISARGAYA
	{0x0D85, 0x0D96, prALetter},                // Lo  [18] SINHALA LETTER AYANNA..SINHALA LETTER AUYANNA
	{0x0D9A, 0x0DB1, prALetter},                // Lo  [24] SINHALA LETTER ALPAPRAANA KAYANNA..SINHALA LETTER DANTAJA NAYANNA
	{0x0DB3, 0x0DBB, prALetter},                // Lo   [9] SINHALA LETTER SANYAKA DAYANNA..SINHALA LETTER RAYANNA
	{0x0DBD, 0x0DBD, prALetter},                // Lo       SINHALA LETTER DANTAJA LAYANNA
	{0x0DC0, 0x0DC6, prALetter},                // Lo   [7] SINHALA LETTER VAYANNA..SINHALA LETTER FAYANNA
	{0x0DCA, 0x0DCA, prExtend},                 // Mn       SINHALA SIGN AL-LAKUNA
	{0x0DCF, 0x0DD1, prExtend},                 // Mc   [3] SINHALA VOWEL SIGN AELA-PILLA..SINHALA VOWEL SIGN DIGA AEDA-PILLA
	{0x0DD2, 0x0DD4, prExtend},                 // Mn   [3] SINHALA VOWEL SIGN KETTI IS-PILLA..SINHALA VOWEL SIGN KETTI PAA-PILLA
	{0x0DD6, 0x0DD6, prExtend},                 // Mn       SINHALA VOWEL SIGN DIGA PAA-PILLA
	{0x0DD8, 0x0DDF, prExtend},                 // Mc   [8] SINHALA VOWEL SIGN GAETTA-PILLA..SINHALA VOWEL SIGN GAYANUKITTA
	{0x0DE6, 0x0DEF, prNumeric},                // Nd  [10] SINHALA LITH DIGIT ZERO..SINHALA LITH DIGIT NINE
	{0x0DF2, 0x0DF3, prExtend},                 // Mc   [2] SINHALA VOWEL SIGN DIGA GAETTA-PILLA..SINHALA VOWEL SIGN DIGA GAYANUKITTA
	{0x0E31, 0x0E31, prExtend},                 // Mn       THAI CHARACTER MAI HAN-AKAT
	{0x0E34, 0x0E3A, prExtend},                 // Mn   [7] THAI CHARACTER SARA I..THAI CHARACTER PHINTHU
	{0x0E47, 0x0E4E, prExtend},                 // Mn   [8] THAI CHARACTER MAITAIKHU..THAI CHARACTER YAMAKKAN
	{0x0E50, 0x0E59, prNumeric},                // Nd  [10] THAI DIGIT ZERO..THAI DIGIT NINE
	{0x0EB1, 0x0EB1, prExtend},                 // Mn       LAO VOWEL SIGN MAI KAN
	{0x0EB4, 0x0EBC, prExtend},                 // Mn   [9] LAO VOWEL SIGN I..LAO SEMIVOWEL SIGN LO
	{0x0EC8, 0x0ECE, prExtend},                 // Mn   [7] LAO TONE MAI EK..LAO YAMAKKAN
	{0x0ED0, 0x0ED9, prNumeric},                // Nd  [10] LAO DIGIT ZERO..LAO DIGIT NINE
	{0x0F00, 0x0F00, prALetter},                // Lo       TIBETAN SYLLABLE OM
	{0x0F18, 0x0F19, prExtend},                 // Mn   [2] TIBETAN ASTROLOGICAL SIGN -KHYUD PA..TIBETAN ASTROLOGICAL SIGN SDONG TSHUGS
	{0x0F20, 0x0F29, prNumeric},                // Nd  [10] TIBETAN DIGIT ZERO..TIBETAN DIGIT NINE
	{0x0F35, 0x0F35, prExtend},                 // Mn       TIBETAN MARK NGAS BZUNG NYI ZLA
	{0x0F37, 0x0F37, prExtend},                 // Mn       TIBETAN MARK NGAS BZUNG SGOR RTAGS
	{0x0F39, 0x0F39, prExtend},                 // Mn       TIBETAN MARK TSA -PHRU
	{0x0F3E, 0x0F3F, prExtend},                 // Mc   [2] TIBETAN SIGN YAR TSHES..TIBETAN SIGN MAR TSHES
	{0x0F40, 0x0F47, prALetter},                // Lo   [8] TIBETAN LETTER KA..TIBETAN LETTER JA
	{0x0F49, 0x0F6C, prALetter},                // Lo  [36] TIBETAN LETTER NYA..TIBETAN LETTER RRA
	{0x0F71, 0x0F7E, prExtend},                 // Mn  [14] TIBETAN VOWEL SIGN AA..TIBETAN SIGN RJES SU NGA RO
	{0x0F7F, 0x0F7F, prExtend},                 // Mc       TIBETAN SIGN RNAM BCAD
	{0x0F80, 0x0F84, prExtend},                 // Mn   [5] TIBETAN VOWEL SIGN REVERSED I..TIBETAN MARK HALANTA
	{0x0F86, 0x0F87, prExtend},                 // Mn   [2] TIBETAN SIGN LCI RTAGS..TIBETAN SIGN YANG RTAGS
	{0x0F88, 0x0F8C, prALetter},                // Lo   [5] TIBETAN SIGN LCE TSA CAN..TIBETAN SIGN INVERTED MCHU CAN
	{0x0F8D, 0x0F97, prExtend},                 // Mn  [11] TIBETAN SUBJOINED SIGN LCE TSA CAN..TIBETAN SUBJOINED LETTER JA
	{0x0F99, 0x0FBC, prExtend},                 // Mn  [36] TIBETAN SUBJOINED LETTER NYA..TIBETAN SUBJOINED LETTER FIXED-FORM RA
	{0x0FC6, 0x0FC6, prExtend},                 // Mn       TIBETAN SYMBOL PADMA GDAN
	{0x102B, 0x102C, prExtend},                 // Mc   [2] MYANMAR VOWEL SIGN TALL AA..MYANMAR VOWEL SIGN AA
	{0x102D, 0x1030, prExtend},                 // Mn   [4] MYANMAR VOWEL SIGN I..MYANMAR VOWEL SIGN UU
	{0x1031, 0x1031, prExtend},                 // Mc       MYANMAR VOWEL SIGN E
	{0x1032, 0x1037, prExtend},                 // Mn   [6] MYANMAR VOWEL SIGN AI..MYANMAR SIGN DOT BELOW
	{0x1038, 0x1038, prExtend},                 // Mc       MYANMAR SIGN VISARGA
	{0x1039, 0x103A, prExtend},                 // Mn   [2] MYANMAR SIGN VIRAMA..MYANMAR SIGN ASAT
	{0x103B, 0x103C, prExtend},                 // Mc   [2] MYANMAR CONSONANT SIGN MEDIAL YA..MYANMAR CONSONANT SIGN MEDIAL RA
	{0x103D, 0x103E, prExtend},                 // Mn   [2] MYANMAR CONSONANT SIGN MEDIAL WA..MYANMAR CONSONANT SIGN MEDIAL HA
	{0x1040, 0x1049, prNumeric},                // Nd  [10] MYANMAR DIGIT ZERO..MYANMAR DIGIT NINE
	{0x1056, 0x1057, prExtend},                 // Mc   [2] MYANMAR VOWEL SIGN VOCALIC R..MYANMAR VOWEL SIGN VOCALIC RR
	{0x1058, 0x1059, prExtend},                 // Mn   [2] MYANMAR VOWEL SIGN VOCALIC L..MYANMAR VOWEL SIGN VOCALIC LL
	{0x105E, 0x1060, prExtend},                 // Mn   [3] MYANMAR CONSONANT SIGN MON MEDIAL NA..MYANMAR CONSONANT SIGN MON MEDIAL LA
	{0x1062, 0x1064, prExtend},                 // Mc   [3] MYANMAR VOWEL SIGN SGAW KAREN EU..MYANMAR TONE MARK SGAW KAREN KE PHO
	{0x1067, 0x106D, prExtend},                 // Mc   [7] MYANMAR VOWEL SIGN WESTERN PWO KAREN EU..MYANMAR SIGN WESTERN PWO KAREN TONE-5
	{0x1071, 0x1074, prExtend},                 // Mn   [4] MYANMAR VOWEL SIGN GEBA KAREN I..MYANMAR VOWEL SIGN KAYAH EE
	{0x1082, 0x1082, prExtend},                 // Mn       MYANMAR CONSONANT SIGN SHAN MEDIAL WA
	{0x1083, 0x1084, prExtend},                 // Mc   [2] MYANMAR VOWEL SIGN SHAN AA..MYANMAR VOWEL SIGN SHAN E
	{0x1085, 0x1086, prExtend},                 // Mn   [2] MYANMAR VOWEL SIGN SHAN E ABOVE..MYANMAR VOWEL SIGN SHAN FINAL Y
	{0x1087, 0x108C, prExtend},                 // Mc   [6] MYANMAR SIGN SHAN TONE-2..MYANMAR SIGN SHAN COUNCIL TONE-3
	{0x108D, 0x108D, prExtend},                 // Mn       MYANMAR SIGN SHAN COUNCIL EMPHATIC TONE
	{0x108F, 0x108F, prExtend},                 // Mc       MYANMAR SIGN RUMAI PALAUNG TONE-5
	{0x1090, 0x1099, prNumeric},                // Nd  [10] MYANMAR SHAN DIGIT ZERO..MYANMAR SHAN DIGIT NINE
	{0x109A, 0x109C, prExtend},                 // Mc   [3] MYANMAR SIGN KHAMTI TONE-1..MYANMAR VOWEL SIGN AITON A
	{0x109D, 0x109D, prExtend},                 // Mn       MYANMAR VOWEL SIGN AITON AI
	{0x10A0, 0x10C5, prALetter},                // L&  [38] GEORGIAN CAPITAL LETTER AN..GEORGIAN CAPITAL LETTER HOE
	{0x10C7, 0x10C7, prALetter},                // L&       GEORGIAN CAPITAL LETTER YN
	{0x10CD, 0x10CD, prALetter},                // L&       GEORGIAN CAPITAL LETTER AEN
	{0x10D0, 0x10FA, prALetter},                // L&  [43] GEORGIAN LETTER AN..GEORGIAN LETTER AIN
	{0x10FC, 0x10FC, prALetter},                // Lm       MODIFIER LETTER GEORGIAN NAR
	{0x10FD, 0x10FF, prALetter},                // L&   [3] GEORGIAN LETTER AEN..GEORGIAN LETTER LABIAL SIGN
	{0x1100, 0x1248, prALetter},                // Lo [329] HANGUL CHOSEONG KIYEOK..ETHIOPIC SYLLABLE QWA
	{0x124A, 0x124D, prALetter},                // Lo   [4] ETHIOPIC SYLLABLE QWI..ETHIOPIC SYLLABLE QWE
	{0x1250, 0x1256, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE QHA..ETHIOPIC SYLLABLE QHO
	{0x1258, 0x1258, prALetter},                // Lo       ETHIOPIC SYLLABLE QHWA
	{0x125A, 0x125D, prALetter},                // Lo   [4] ETHIOPIC SYLLABLE QHWI..ETHIOPIC SYLLABLE QHWE
	{0x1260, 0x1288, prALetter},                // Lo  [41] ETHIOPIC SYLLABLE BA..ETHIOPIC SYLLABLE XWA
	{0x128A, 0x128D, prALetter},                // Lo   [4] ETHIOPIC SYLLABLE XWI..ETHIOPIC SYLLABLE XWE
	{0x1290, 0x12B0, prALetter},                // Lo  [33] ETHIOPIC SYLLABLE NA..ETHIOPIC SYLLABLE KWA
	{0x12B2, 0x12B5, prALetter},                // Lo   [4] ETHIOPIC SYLLABLE KWI..ETHIOPIC SYLLABLE KWE
	{0x12B8, 0x12BE, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE KXA..ETHIOPIC SYLLABLE KXO
	{0x12C0, 0x12C0, prALetter},                // Lo       ETHIOPIC SYLLABLE KXWA
	{0x12C2, 0x12C5, prALetter},                // Lo   [4] ETHIOPIC SYLLABLE KXWI..ETHIOPIC SYLLABLE KXWE
	{0x12C8, 0x12D6, prALetter},                // Lo  [15] ETHIOPIC SYLLABLE WA..ETHIOPIC SYLLABLE PHARYNGEAL O
	{0x12D8, 0x1310, prALetter},                // Lo  [57] ETHIOPIC SYLLABLE ZA..ETHIOPIC SYLLABLE GWA
	{0x1312, 0x1315, prALetter},                // Lo   [4] ETHIOPIC SYLLABLE GWI..ETHIOPIC SYLLABLE GWE
	{0x1318, 0x135A, prALetter},                // Lo  [67] ETHIOPIC SYLLABLE GGA..ETHIOPIC SYLLABLE FYA
	{0x135D, 0x135F, prExtend},                 // Mn   [3] ETHIOPIC COMBINING GEMINATION AND VOWEL LENGTH MARK..ETHIOPIC COMBINING GEMINATION MARK
	{0x1380, 0x138F, prALetter},                // Lo  [16] ETHIOPIC SYLLABLE SEBATBEIT MWA..ETHIOPIC SYLLABLE PWE
	{0x13A0, 0x13F5, prALetter},                // L&  [86] CHEROKEE LETTER A..CHEROKEE LETTER MV
	{0x13F8, 0x13FD, prALetter},                // L&   [6] CHEROKEE SMALL LETTER YE..CHEROKEE SMALL LETTER MV
	{0x1401, 0x166C, prALetter},                // Lo [620] CANADIAN SYLLABICS E..CANADIAN SYLLABICS CARRIER TTSA
	{0x166F, 0x167F, prALetter},                // Lo  [17] CANADIAN SYLLABICS QAI..CANADIAN SYLLABICS BLACKFOOT W
	{0x1680, 0x1680, prWSegSpace},              // Zs       OGHAM SPACE MARK
	{0x1681, 0x169A, prALetter},                // Lo  [26] OGHAM LETTER BEITH..OGHAM LETTER PEITH
	{0x16A0, 0x16EA, prALetter},                // Lo  [75] RUNIC LETTER FEHU FEOH FE F..RUNIC LETTER X
	{0x16EE, 0x16F0, prALetter},                // Nl   [3] RUNIC ARLAUG SYMBOL..RUNIC BELGTHOR SYMBOL
	{0x16F1, 0x16F8, prALetter},                // Lo   [8] RUNIC LETTER K..RUNIC LETTER FRANKS CASKET AESC
	{0x1700, 0x1711, prALetter},                // Lo  [18] TAGALOG LETTER A..TAGALOG LETTER HA
	{0x1712, 0x1714, prExtend},                 // Mn   [3] TAGALOG VOWEL SIGN I..TAGALOG SIGN VIRAMA
	{0x1715, 0x1715, prExtend},                 // Mc       TAGALOG SIGN PAMUDPOD
	{0x171F, 0x1731, prALetter},                // Lo  [19] TAGALOG LETTER ARCHAIC RA..HANUNOO LETTER HA
	{0x1732, 0x1733, prExtend},                 // Mn   [2] HANUNOO VOWEL SIGN I..HANUNOO VOWEL SIGN U
	{0x1734, 0x1734, prExtend},                 // Mc       HANUNOO SIGN PAMUDPOD
	{0x1740, 0x1751, prALetter},                // Lo  [18] BUHID LETTER A..BUHID LETTER HA
	{0x1752, 0x1753, prExtend},                 // Mn   [2] BUHID VOWEL SIGN I..BUHID VOWEL SIGN U
	{0x1760, 0x176C, prALetter},                // Lo  [13] TAGBANWA LETTER A..TAGBANWA LETTER YA
	{0x176E, 0x1770, prALetter},                // Lo   [3] TAGBANWA LETTER LA..TAGBANWA LETTER SA
	{0x1772, 0x1773, prExtend},                 // Mn   [2] TAGBANWA VOWEL SIGN I..TAGBANWA VOWEL SIGN U
	{0x17B4, 0x17B5, prExtend},                 // Mn   [2] KHMER VOWEL INHERENT AQ..KHMER VOWEL INHERENT AA
	{0x17B6, 0x17B6, prExtend},                 // Mc       KHMER VOWEL SIGN AA
	{0x17B7, 0x17BD, prExtend},                 // Mn   [7] KHMER VOWEL SIGN I..KHMER VOWEL SIGN UA
	{0x17BE, 0x17C5, prExtend},                 // Mc   [8] KHMER VOWEL SIGN OE..KHMER VOWEL SIGN AU
	{0x17C6, 0x17C6, prExtend},                 // Mn       KHMER SIGN NIKAHIT
	{0x17C7, 0x17C8, prExtend},                 // Mc   [2] KHMER SIGN REAHMUK..KHMER SIGN YUUKALEAPINTU
	{0x17C9, 0x17D3, prExtend},                 // Mn  [11] KHMER SIGN MUUSIKATOAN..KHMER SIGN BATHAMASAT
	{0x17DD, 0x17DD, prExtend},                 // Mn       KHMER SIGN ATTHACAN
	{0x17E0, 0x17E9, prNumeric},                // Nd  [10] KHMER DIGIT ZERO..KHMER DIGIT NINE
	{0x180B, 0x180D, prExtend},                 // Mn   [3] MONGOLIAN FREE VARIATION SELECTOR ONE..MONGOLIAN FREE VARIATION SELECTOR THREE
	{0x180E, 0x180E, prFormat},                 // Cf       MONGOLIAN VOWEL SEPARATOR
	{0x180F, 0x180F, prExtend},                 // Mn       MONGOLIAN FREE VARIATION SELECTOR FOUR
	{0x1810, 0x1819, prNumeric},                // Nd  [10] MONGOLIAN DIGIT ZERO..MONGOLIAN DIGIT NINE
	{0x1820, 0x1842, prALetter},                // Lo  [35] MONGOLIAN LETTER A..MONGOLIAN LETTER CHI
	{0x1843, 0x1843, prALetter},                // Lm       MONGOLIAN LETTER TODO LONG VOWEL SIGN
	{0x1844, 0x1878, prALetter},                // Lo  [53] MONGOLIAN LETTER TODO E..MONGOLIAN LETTER CHA WITH TWO DOTS
	{0x1880, 0x1884, prALetter},                // Lo   [5] MONGOLIAN LETTER ALI GALI ANUSVARA ONE..MONGOLIAN LETTER ALI GALI INVERTED UBADAMA
	{0x1885, 0x1886, prExtend},                 // Mn   [2] MONGOLIAN LETTER ALI GALI BALUDA..MONGOLIAN LETTER ALI GALI THREE BALUDA
	{0x1887, 0x18A8, prALetter},                // Lo  [34] MONGOLIAN LETTER ALI GALI A..MONGOLIAN LETTER MANCHU ALI GALI BHA
	{0x18A9, 0x18A9, prExtend},                 // Mn       MONGOLIAN LETTER ALI GALI DAGALGA
	{0x18AA, 0x18AA, prALetter},                // Lo       MONGOLIAN LETTER MANCHU ALI GALI LHA
	{0x18B0, 0x18F5, prALetter},                // Lo  [70] CANADIAN SYLLABICS OY..CANADIAN SYLLABICS CARRIER DENTAL S
	{0x1900, 0x191E, prALetter},                // Lo  [31] LIMBU VOWEL-CARRIER LETTER..LIMBU LETTER TRA
	{0x1920, 0x1922, prExtend},                 // Mn   [3] LIMBU VOWEL SIGN A..LIMBU VOWEL SIGN U
	{0x1923, 0x1926, prExtend},                 // Mc   [4] LIMBU VOWEL SIGN EE..LIMBU VOWEL SIGN AU
	{0x1927, 0x1928, prExtend},                 // Mn   [2] LIMBU VOWEL SIGN E..LIMBU VOWEL SIGN O
	{0x1929, 0x192B, prExtend},                 // Mc   [3] LIMBU SUBJOINED LETTER YA..LIMBU SUBJOINED LETTER WA
	{0x1930, 0x1931, prExtend},                 // Mc   [2] LIMBU SMALL LETTER KA..LIMBU SMALL LETTER NGA
	{0x1932, 0x1932, prExtend},                 // Mn       LIMBU SMALL LETTER ANUSVARA
	{0x1933, 0x1938, prExtend},                 // Mc   [6] LIMBU SMALL LETTER TA..LIMBU SMALL LETTER LA
	{0x1939, 0x193B, prExtend},                 // Mn   [3] LIMBU SIGN MUKPHRENG..LIMBU SIGN SA-I
	{0x1946, 0x194F, prNumeric},                // Nd  [10] LIMBU DIGIT ZERO..LIMBU DIGIT NINE
	{0x19D0, 0x19D9, prNumeric},                // Nd  [10] NEW TAI LUE DIGIT ZERO..NEW TAI LUE DIGIT NINE
	{0x1A00, 0x1A16, prALetter},                // Lo  [23] BUGINESE LETTER KA..BUGINESE LETTER HA
	{0x1A17, 0x1A18, prExtend},                 // Mn   [2] BUGINESE VOWEL SIGN I..BUGINESE VOWEL SIGN U
	{0x1A19, 0x1A1A, prExtend},                 // Mc   [2] BUGINESE VOWEL SIGN E..BUGINESE VOWEL SIGN O
	{0x1A1B, 0x1A1B, prExtend},                 // Mn       BUGINESE VOWEL SIGN AE
	{0x1A55, 0x1A55, prExtend},                 // Mc       TAI THAM CONSONANT SIGN MEDIAL RA
	{0x1A56, 0x1A56, prExtend},                 // Mn       TAI THAM CONSONANT SIGN MEDIAL LA
	{0x1A57, 0x1A57, prExtend},                 // Mc       TAI THAM CONSONANT SIGN LA TANG LAI
	{0x1A58, 0x1A5E, prExtend},                 // Mn   [7] TAI THAM SIGN MAI KANG LAI..TAI THAM CONSONANT SIGN SA
	{0x1A60, 0x1A60, prExtend},                 // Mn       TAI THAM SIGN SAKOT
	{0x1A61, 0x1A61, prExtend},                 // Mc       TAI THAM VOWEL SIGN A
	{0x1A62, 0x1A62, prExtend},                 // Mn       TAI THAM VOWEL SIGN MAI SAT
	{0x1A63, 0x1A64, prExtend},                 // Mc   [2] TAI THAM VOWEL SIGN AA..TAI THAM VOWEL SIGN TALL AA
	{0x1A65, 0x1A6C, prExtend},                 // Mn   [8] TAI THAM VOWEL SIGN I..TAI THAM VOWEL SIGN OA BELOW
	{0x1A6D, 0x1A72, prExtend},                 // Mc   [6] TAI THAM VOWEL SIGN OY..TAI THAM VOWEL SIGN THAM AI
	{0x1A73, 0x1A7C, prExtend},                 // Mn  [10] TAI THAM VOWEL SIGN OA ABOVE..TAI THAM SIGN KHUEN-LUE KARAN
	{0x1A7F, 0x1A7F, prExtend},                 // Mn       TAI THAM COMBINING CRYPTOGRAMMIC DOT
	{0x1A80, 0x1A89, prNumeric},                // Nd  [10] TAI THAM HORA DIGIT ZERO..TAI THAM HORA DIGIT NINE
	{0x1A90, 0x1A99, prNumeric},                // Nd  [10] TAI THAM THAM DIGIT ZERO..TAI THAM THAM DIGIT NINE
	{0x1AB0, 0x1ABD, prExtend},                 // Mn  [14] COMBINING DOUBLED CIRCUMFLEX ACCENT..COMBINING PARENTHESES BELOW
	{0x1ABE, 0x1ABE, prExtend},                 // Me       COMBINING PARENTHESES OVERLAY
	{0x1ABF, 0x1ACE, prExtend},                 // Mn  [16] COMBINING LATIN SMALL LETTER W BELOW..COMBINING LATIN SMALL LETTER INSULAR T
	{0x1B00, 0x1B03, prExtend},                 // Mn   [4] BALINESE SIGN ULU RICEM..BALINESE SIGN SURANG
	{0x1B04, 0x1B04, prExtend},                 // Mc       BALINESE SIGN BISAH
	{0x1B05, 0x1B33, prALetter},                // Lo  [47] BALINESE LETTER AKARA..BALINESE LETTER HA
	{0x1B34, 0x1B34, prExtend},                 // Mn       BALINESE SIGN REREKAN
	{0x1B35, 0x1B35, prExtend},                 // Mc       BALINESE VOWEL SIGN TEDUNG
	{0x1B36, 0x1B3A, prExtend},                 // Mn   [5] BALINESE VOWEL SIGN ULU..BALINESE VOWEL SIGN RA REPA
	{0x1B3B, 0x1B3B, prExtend},                 // Mc       BALINESE VOWEL SIGN RA REPA TEDUNG
	{0x1B3C, 0x1B3C, prExtend},                 // Mn       BALINESE VOWEL SIGN LA LENGA
	{0x1B3D, 0x1B41, prExtend},                 // Mc   [5] BALINESE VOWEL SIGN LA LENGA TEDUNG..BALINESE VOWEL SIGN TALING REPA TEDUNG
	{0x1B42, 0x1B42, prExtend},                 // Mn       BALINESE VOWEL SIGN PEPET
	{0x1B43, 0x1B44, prExtend},                 // Mc   [2] BALINESE VOWEL SIGN PEPET TEDUNG..BALINESE ADEG ADEG
	{0x1B45, 0x1B4C, prALetter},                // Lo   [8] BALINESE LETTER KAF SASAK..BALINESE LETTER ARCHAIC JNYA
	{0x1B50, 0x1B59, prNumeric},                // Nd  [10] BALINESE DIGIT ZERO..BALINESE DIGIT NINE
	{0x1B6B, 0x1B73, prExtend},                 // Mn   [9] BALINESE MUSICAL SYMBOL COMBINING TEGEH..BALINESE MUSICAL SYMBOL COMBINING GONG
	{0x1B80, 0x1B81, prExtend},                 // Mn   [2] SUNDANESE SIGN PANYECEK..SUNDANESE SIGN PANGLAYAR
	{0x1B82, 0x1B82, prExtend},                 // Mc       SUNDANESE SIGN PANGWISAD
	{0x1B83, 0x1BA0, prALetter},                // Lo  [30] SUNDANESE LETTER A..SUNDANESE LETTER HA
	{0x1BA1, 0x1BA1, prExtend},                 // Mc       SUNDANESE CONSONANT SIGN PAMINGKAL
	{0x1BA2, 0x1BA5, prExtend},                 // Mn   [4] SUNDANESE CONSONANT SIGN PANYAKRA..SUNDANESE VOWEL SIGN PANYUKU
	{0x1BA6, 0x1BA7, prExtend},                 // Mc   [2] SUNDANESE VOWEL SIGN PANAELAENG..SUNDANESE VOWEL SIGN PANOLONG
	{0x1BA8, 0x1BA9, prExtend},                 // Mn   [2] SUNDANESE VOWEL SIGN PAMEPET..SUNDANESE VOWEL SIGN PANEULEUNG
	{0x1BAA, 0x1BAA, prExtend},                 // Mc       SUNDANESE SIGN PAMAAEH
	{0x1BAB, 0x1BAD, prExtend},                 // Mn   [3] SUNDANESE SIGN VIRAMA..SUNDANESE CONSONANT SIGN PASANGAN WA
	{0x1BAE, 0x1BAF, prALetter},                // Lo   [2] SUNDANESE LETTER KHA..SUNDANESE LETTER SYA
	{0x1BB0, 0x1BB9, prNumeric},                // Nd  [10] SUNDANESE DIGIT ZERO..SUNDANESE DIGIT NINE
	{0x1BBA, 0x1BE5, prALetter},                // Lo  [44] SUNDANESE AVAGRAHA..BATAK LETTER U
	{0x1BE6, 0x1BE6, prExtend},                 // Mn       BATAK SIGN TOMPI
	{0x1BE7, 0x1BE7, prExtend},                 // Mc       BATAK VOWEL SIGN E
	{0x1BE8, 0x1BE9, prExtend},                 // Mn   [2] BATAK VOWEL SIGN PAKPAK E..BATAK VOWEL SIGN EE
	{0x1BEA, 0x1BEC, prExtend},                 // Mc   [3] BATAK VOWEL SIGN I..BATAK VOWEL SIGN O
	{0x1BED, 0x1BED, prExtend},                 // Mn       BATAK VOWEL SIGN KARO O
	{0x1BEE, 0x1BEE, prExtend},                 // Mc       BATAK VOWEL SIGN U
	{0x1BEF, 0x1BF1, prExtend},                 // Mn   [3] BATAK VOWEL SIGN U FOR SIMALUNGUN SA..BATAK CONSONANT SIGN H
	{0x1BF2, 0x1BF3, prExtend},                 // Mc   [2] BATAK PANGOLAT..BATAK PANONGONAN
	{0x1C00, 0x1C23, prALetter},                // Lo  [36] LEPCHA LETTER KA..LEPCHA LETTER A
	{0x1C24, 0x1C2B, prExtend},                 // Mc   [8] LEPCHA SUBJOINED LETTER YA..LEPCHA VOWEL SIGN UU
	{0x1C2C, 0x1C33, prExtend},                 // Mn   [8] LEPCHA VOWEL SIGN E..LEPCHA CONSONANT SIGN T
	{0x1C34, 0x1C35, prExtend},                 // Mc   [2] LEPCHA CONSONANT SIGN NYIN-DO..LEPCHA CONSONANT SIGN KANG
	{0x1C36, 0x1C37, prExtend},                 // Mn   [2] LEPCHA SIGN RAN..LEPCHA SIGN NUKTA
	{0x1C40, 0x1C49, prNumeric},                // Nd  [10] LEPCHA DIGIT ZERO..LEPCHA DIGIT NINE
	{0x1C4D, 0x1C4F, prALetter},                // Lo   [3] LEPCHA LETTER TTA..LEPCHA LETTER DDA
	{0x1C50, 0x1C59, prNumeric},                // Nd  [10] OL CHIKI DIGIT ZERO..OL CHIKI DIGIT NINE
	{0x1C5A, 0x1C77, prALetter},                // Lo  [30] OL CHIKI LETTER LA..OL CHIKI LETTER OH
	{0x1C78, 0x1C7D, prALetter},                // Lm   [6] OL CHIKI MU TTUDDAG..OL CHIKI AHAD
	{0x1C80, 0x1C88, prALetter},                // L&   [9] CYRILLIC SMALL LETTER ROUNDED VE..CYRILLIC SMALL LETTER UNBLENDED UK
	{0x1C90, 0x1CBA, prALetter},                // L&  [43] GEORGIAN MTAVRULI CAPITAL LETTER AN..GEORGIAN MTAVRULI CAPITAL LETTER AIN
	{0x1CBD, 0x1CBF, prALetter},                // L&   [3] GEORGIAN MTAVRULI CAPITAL LETTER AEN..GEORGIAN MTAVRULI CAPITAL LETTER LABIAL SIGN
	{0x1CD0, 0x1CD2, prExtend},                 // Mn   [3] VEDIC TONE KARSHANA..VEDIC TONE PRENKHA
	{0x1CD4, 0x1CE0, prExtend},                 // Mn  [13] VEDIC SIGN YAJURVEDIC MIDLINE SVARITA..VEDIC TONE RIGVEDIC KASHMIRI INDEPENDENT SVARITA
	{0x1CE1, 0x1CE1, prExtend},                 // Mc       VEDIC TONE ATHARVAVEDIC INDEPENDENT SVARITA
	{0x1CE2, 0x1CE8, prExtend},                 // Mn   [7] VEDIC SIGN VISARGA SVARITA..VEDIC SIGN VISARGA ANUDATTA WITH TAIL
	{0x1CE9, 0x1CEC, prALetter},                // Lo   [4] VEDIC SIGN ANUSVARA ANTARGOMUKHA..VEDIC SIGN ANUSVARA VAMAGOMUKHA WITH TAIL
	{0x1CED, 0x1CED, prExtend},                 // Mn       VEDIC SIGN TIRYAK
	{0x1CEE, 0x1CF3, prALetter},                // Lo   [6] VEDIC SIGN HEXIFORM LONG ANUSVARA..VEDIC SIGN ROTATED ARDHAVISARGA
	{0x1CF4, 0x1CF4, prExtend},                 // Mn       VEDIC TONE CANDRA ABOVE
	{0x1CF5, 0x1CF6, prALetter},                // Lo   [2] VEDIC SIGN JIHVAMULIYA..VEDIC SIGN UPADHMANIYA
	{0x1CF7, 0x1CF7, prExtend},                 // Mc       VEDIC SIGN ATIKRAMA
	{0x1CF8, 0x1CF9, prExtend},                 // Mn   [2] VEDIC TONE RING ABOVE..VEDIC TONE DOUBLE RING ABOVE
	{0x1CFA, 0x1CFA, prALetter},                // Lo       VEDIC SIGN DOUBLE ANUSVARA ANTARGOMUKHA
	{0x1D00, 0x1D2B, prALetter},                // L&  [44] LATIN LETTER SMALL CAPITAL A..CYRILLIC LETTER SMALL CAPITAL EL
	{0x1D2C, 0x1D6A, prALetter},                // Lm  [63] MODIFIER LETTER CAPITAL A..GREEK SUBSCRIPT SMALL LETTER CHI
	{0x1D6B, 0x1D77, prALetter},                // L&  [13] LATIN SMALL LETTER UE..LATIN SMALL LETTER TURNED G
	{0x1D78, 0x1D78, prALetter},                // Lm       MODIFIER LETTER CYRILLIC EN
	{0x1D79, 0x1D9A, prALetter},                // L&  [34] LATIN SMALL LETTER INSULAR G..LATIN SMALL LETTER EZH WITH RETROFLEX HOOK
	{0x1D9B, 0x1DBF, prALetter},                // Lm  [37] MODIFIER LETTER SMALL TURNED ALPHA..MODIFIER LETTER SMALL THETA
	{0x1DC0, 0x1DFF, prExtend},                 // Mn  [64] COMBINING DOTTED GRAVE ACCENT..COMBINING RIGHT ARROWHEAD AND DOWN ARROWHEAD BELOW
	{0x1E00, 0x1F15, prALetter},                // L& [278] LATIN CAPITAL LETTER A WITH RING BELOW..GREEK SMALL LETTER EPSILON WITH DASIA AND OXIA
	{0x1F18, 0x1F1D, prALetter},                // L&   [6] GREEK CAPITAL LETTER EPSILON WITH PSILI..GREEK CAPITAL LETTER EPSILON WITH DASIA AND OXIA
	{0x1F20, 0x1F45, prALetter},                // L&  [38] GREEK SMALL LETTER ETA WITH PSILI..GREEK SMALL LETTER OMICRON WITH DASIA AND OXIA
	{0x1F48, 0x1F4D, prALetter},                // L&   [6] GREEK CAPITAL LETTER OMICRON WITH PSILI..GREEK CAPITAL LETTER OMICRON WITH DASIA AND OXIA
	{0x1F50, 0x1F57, prALetter},                // L&   [8] GREEK SMALL LETTER UPSILON WITH PSILI..GREEK SMALL LETTER UPSILON WITH DASIA AND PERISPOMENI
	{0x1F59, 0x1F59, prALetter},                // L&       GREEK CAPITAL LETTER UPSILON WITH DASIA
	{0x1F5B, 0x1F5B, prALetter},                // L&       GREEK CAPITAL LETTER UPSILON WITH DASIA AND VARIA
	{0x1F5D, 0x1F5D, prALetter},                // L&       GREEK CAPITAL LETTER UPSILON WITH DASIA AND OXIA
	{0x1F5F, 0x1F7D, prALetter},                // L&  [31] GREEK CAPITAL LETTER UPSILON WITH DASIA AND PERISPOMENI..GREEK SMALL LETTER OMEGA WITH OXIA
	{0x1F80, 0x1FB4, prALetter},                // L&  [53] GREEK SMALL LETTER ALPHA WITH PSILI AND YPOGEGRAMMENI..GREEK SMALL LETTER ALPHA WITH OXIA AND YPOGEGRAMMENI
	{0x1FB6, 0x1FBC, prALetter},                // L&   [7] GREEK SMALL LETTER ALPHA WITH PERISPOMENI..GREEK CAPITAL LETTER ALPHA WITH PROSGEGRAMMENI
	{0x1FBE, 0x1FBE, prALetter},                // L&       GREEK PROSGEGRAMMENI
	{0x1FC2, 0x1FC4, prALetter},                // L&   [3] GREEK SMALL LETTER ETA WITH VARIA AND YPOGEGRAMMENI..GREEK SMALL LETTER ETA WITH OXIA AND YPOGEGRAMMENI
	{0x1FC6, 0x1FCC, prALetter},                // L&   [7] GREEK SMALL LETTER ETA WITH PERISPOMENI..GREEK CAPITAL LETTER ETA WITH PROSGEGRAMMENI
	{0x1FD0, 0x1FD3, prALetter},                // L&   [4] GREEK SMALL LETTER IOTA WITH VRACHY..GREEK SMALL LETTER IOTA WITH DIALYTIKA AND OXIA
	{0x1FD6, 0x1FDB, prALetter},                // L&   [6] GREEK SMALL LETTER IOTA WITH PERISPOMENI..GREEK CAPITAL LETTER IOTA WITH OXIA
	{0x1FE0, 0x1FEC, prALetter},                // L&  [13] GREEK SMALL LETTER UPSILON WITH VRACHY..GREEK CAPITAL LETTER RHO WITH DASIA
	{0x1FF2, 0x1FF4, prALetter},                // L&   [3] GREEK SMALL LETTER OMEGA WITH VARIA AND YPOGEGRAMMENI..GREEK SMALL LETTER OMEGA WITH OXIA AND YPOGEGRAMMENI
	{0x1FF6, 0x1FFC, prALetter},                // L&   [7] GREEK SMALL LETTER OMEGA WITH PERISPOMENI..GREEK CAPITAL LETTER OMEGA WITH PROSGEGRAMMENI
	{0x2000, 0x2006, prWSegSpace},              // Zs   [7] EN QUAD..SIX-PER-EM SPACE
	{0x2008, 0x200A, prWSegSpace},              // Zs   [3] PUNCTUATION SPACE..HAIR SPACE
	{0x200C, 0x200C, prExtend},                 // Cf       ZERO WIDTH NON-JOINER
	{0x200D, 0x200D, prZWJ},                    // Cf       ZERO WIDTH JOINER
	{0x200E, 0x200F, prFormat},                 // Cf   [2] LEFT-TO-RIGHT MARK..RIGHT-TO-LEFT MARK
	{0x2018, 0x2018, prMidNumLet},              // Pi       LEFT SINGLE QUOTATION MARK
	{0x2019, 0x2019, prMidNumLet},              // Pf       RIGHT SINGLE QUOTATION MARK
	{0x2024, 0x2024, prMidNumLet},              // Po       ONE DOT LEADER
	{0x2027, 0x2027, prMidLetter},              // Po       HYPHENATION POINT
	{0x2028, 0x2028, prNewline},                // Zl       LINE SEPARATOR
	{0x2029, 0x2029, prNewline},                // Zp       PARAGRAPH SEPARATOR
	{0x202A, 0x202E, prFormat},                 // Cf   [5] LEFT-TO-RIGHT EMBEDDING..RIGHT-TO-LEFT OVERRIDE
	{0x202F, 0x202F, prExtendNumLet},           // Zs       NARROW NO-BREAK SPACE
	{0x203C, 0x203C, prExtendedPictographic},   // E0.6   [1] (‼️)       double exclamation mark
	{0x203F, 0x2040, prExtendNumLet},           // Pc   [2] UNDERTIE..CHARACTER TIE
	{0x2044, 0x2044, prMidNum},                 // Sm       FRACTION SLASH
	{0x2049, 0x2049, prExtendedPictographic},   // E0.6   [1] (⁉️)       exclamation question mark
	{0x2054, 0x2054, prExtendNumLet},           // Pc       INVERTED UNDERTIE
	{0x205F, 0x205F, prWSegSpace},              // Zs       MEDIUM MATHEMATICAL SPACE
	{0x2060, 0x2064, prFormat},                 // Cf   [5] WORD JOINER..INVISIBLE PLUS
	{0x2066, 0x206F, prFormat},                 // Cf  [10] LEFT-TO-RIGHT ISOLATE..NOMINAL DIGIT SHAPES
	{0x2071, 0x2071, prALetter},                // Lm       SUPERSCRIPT LATIN SMALL LETTER I
	{0x207F, 0x207F, prALetter},                // Lm       SUPERSCRIPT LATIN SMALL LETTER N
	{0x2090, 0x209C, prALetter},                // Lm  [13] LATIN SUBSCRIPT SMALL LETTER A..LATIN SUBSCRIPT SMALL LETTER T
	{0x20D0, 0x20DC, prExtend},                 // Mn  [13] COMBINING LEFT HARPOON ABOVE..COMBINING FOUR DOTS ABOVE
	{0x20DD, 0x20E0, prExtend},                 // Me   [4] COMBINING ENCLOSING CIRCLE..COMBINING ENCLOSING CIRCLE BACKSLASH
	{0x20E1, 0x20E1, prExtend},                 // Mn       COMBINING LEFT RIGHT ARROW ABOVE
	{0x20E2, 0x20E4, prExtend},                 // Me   [3] COMBINING ENCLOSING SCREEN..COMBINING ENCLOSING UPWARD POINTING TRIANGLE
	{0x20E5, 0x20F0, prExtend},                 // Mn  [12] COMBINING REVERSE SOLIDUS OVERLAY..COMBINING ASTERISK ABOVE
	{0x2102, 0x2102, prALetter},                // L&       DOUBLE-STRUCK CAPITAL C
	{0x2107, 0x2107, prALetter},                // L&       EULER CONSTANT
	{0x210A, 0x2113, prALetter},                // L&  [10] SCRIPT SMALL G..SCRIPT SMALL L
	{0x2115, 0x2115, prALetter},                // L&       DOUBLE-STRUCK CAPITAL N
	{0x2119, 0x211D, prALetter},                // L&   [5] DOUBLE-STRUCK CAPITAL P..DOUBLE-STRUCK CAPITAL R
	{0x2122, 0x2122, prExtendedPictographic},   // E0.6   [1] (™️)       trade mark
	{0x2124, 0x2124, prALetter},                // L&       DOUBLE-STRUCK CAPITAL Z
	{0x2126, 0x2126, prALetter},                // L&       OHM SIGN
	{0x2128, 0x2128, prALetter},                // L&       BLACK-LETTER CAPITAL Z
	{0x212A, 0x212D, prALetter},                // L&   [4] KELVIN SIGN..BLACK-LETTER CAPITAL C
	{0x212F, 0x2134, prALetter},                // L&   [6] SCRIPT SMALL E..SCRIPT SMALL O
	{0x2135, 0x2138, prALetter},                // Lo   [4] ALEF SYMBOL..DALET SYMBOL
	{0x2139, 0x2139, prExtendedPictographic},   // E0.6   [1] (ℹ️)       information
	{0x2139, 0x2139, prALetter},                // L&       INFORMATION SOURCE
	{0x213C, 0x213F, prALetter},                // L&   [4] DOUBLE-STRUCK SMALL PI..DOUBLE-STRUCK CAPITAL PI
	{0x2145, 0x2149, prALetter},                // L&   [5] DOUBLE-STRUCK ITALIC CAPITAL D..DOUBLE-STRUCK ITALIC SMALL J
	{0x214E, 0x214E, prALetter},                // L&       TURNED SMALL F
	{0x2160, 0x2182, prALetter},                // Nl  [35] ROMAN NUMERAL ONE..ROMAN NUMERAL TEN THOUSAND
	{0x2183, 0x2184, prALetter},                // L&   [2] ROMAN NUMERAL REVERSED ONE HUNDRED..LATIN SMALL LETTER REVERSED C
	{0x2185, 0x2188, prALetter},                // Nl   [4] ROMAN NUMERAL SIX LATE FORM..ROMAN NUMERAL ONE HUNDRED THOUSAND
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
	{0x24B6, 0x24E9, prALetter},                // So  [52] CIRCLED LATIN CAPITAL LETTER A..CIRCLED LATIN SMALL LETTER Z
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
	{0x2C00, 0x2C7B, prALetter},                // L& [124] GLAGOLITIC CAPITAL LETTER AZU..LATIN LETTER SMALL CAPITAL TURNED E
	{0x2C7C, 0x2C7D, prALetter},                // Lm   [2] LATIN SUBSCRIPT SMALL LETTER J..MODIFIER LETTER CAPITAL V
	{0x2C7E, 0x2CE4, prALetter},                // L& [103] LATIN CAPITAL LETTER S WITH SWASH TAIL..COPTIC SYMBOL KAI
	{0x2CEB, 0x2CEE, prALetter},                // L&   [4] COPTIC CAPITAL LETTER CRYPTOGRAMMIC SHEI..COPTIC SMALL LETTER CRYPTOGRAMMIC GANGIA
	{0x2CEF, 0x2CF1, prExtend},                 // Mn   [3] COPTIC COMBINING NI ABOVE..COPTIC COMBINING SPIRITUS LENIS
	{0x2CF2, 0x2CF3, prALetter},                // L&   [2] COPTIC CAPITAL LETTER BOHAIRIC KHEI..COPTIC SMALL LETTER BOHAIRIC KHEI
	{0x2D00, 0x2D25, prALetter},                // L&  [38] GEORGIAN SMALL LETTER AN..GEORGIAN SMALL LETTER HOE
	{0x2D27, 0x2D27, prALetter},                // L&       GEORGIAN SMALL LETTER YN
	{0x2D2D, 0x2D2D, prALetter},                // L&       GEORGIAN SMALL LETTER AEN
	{0x2D30, 0x2D67, prALetter},                // Lo  [56] TIFINAGH LETTER YA..TIFINAGH LETTER YO
	{0x2D6F, 0x2D6F, prALetter},                // Lm       TIFINAGH MODIFIER LETTER LABIALIZATION MARK
	{0x2D7F, 0x2D7F, prExtend},                 // Mn       TIFINAGH CONSONANT JOINER
	{0x2D80, 0x2D96, prALetter},                // Lo  [23] ETHIOPIC SYLLABLE LOA..ETHIOPIC SYLLABLE GGWE
	{0x2DA0, 0x2DA6, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE SSA..ETHIOPIC SYLLABLE SSO
	{0x2DA8, 0x2DAE, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE CCA..ETHIOPIC SYLLABLE CCO
	{0x2DB0, 0x2DB6, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE ZZA..ETHIOPIC SYLLABLE ZZO
	{0x2DB8, 0x2DBE, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE CCHA..ETHIOPIC SYLLABLE CCHO
	{0x2DC0, 0x2DC6, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE QYA..ETHIOPIC SYLLABLE QYO
	{0x2DC8, 0x2DCE, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE KYA..ETHIOPIC SYLLABLE KYO
	{0x2DD0, 0x2DD6, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE XYA..ETHIOPIC SYLLABLE XYO
	{0x2DD8, 0x2DDE, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE GYA..ETHIOPIC SYLLABLE GYO
	{0x2DE0, 0x2DFF, prExtend},                 // Mn  [32] COMBINING CYRILLIC LETTER BE..COMBINING CYRILLIC LETTER IOTIFIED BIG YUS
	{0x2E2F, 0x2E2F, prALetter},                // Lm       VERTICAL TILDE
	{0x3000, 0x3000, prWSegSpace},              // Zs       IDEOGRAPHIC SPACE
	{0x3005, 0x3005, prALetter},                // Lm       IDEOGRAPHIC ITERATION MARK
	{0x302A, 0x302D, prExtend},                 // Mn   [4] IDEOGRAPHIC LEVEL TONE MARK..IDEOGRAPHIC ENTERING TONE MARK
	{0x302E, 0x302F, prExtend},                 // Mc   [2] HANGUL SINGLE DOT TONE MARK..HANGUL DOUBLE DOT TONE MARK
	{0x3030, 0x3030, prExtendedPictographic},   // E0.6   [1] (〰️)       wavy dash
	{0x3031, 0x3035, prKatakana},               // Lm   [5] VERTICAL KANA REPEAT MARK..VERTICAL KANA REPEAT MARK LOWER HALF
	{0x303B, 0x303B, prALetter},                // Lm       VERTICAL IDEOGRAPHIC ITERATION MARK
	{0x303C, 0x303C, prALetter},                // Lo       MASU MARK
	{0x303D, 0x303D, prExtendedPictographic},   // E0.6   [1] (〽️)       part alternation mark
	{0x3099, 0x309A, prExtend},                 // Mn   [2] COMBINING KATAKANA-HIRAGANA VOICED SOUND MARK..COMBINING KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	{0x309B, 0x309C, prKatakana},               // Sk   [2] KATAKANA-HIRAGANA VOICED SOUND MARK..KATAKANA-HIRAGANA SEMI-VOICED SOUND MARK
	{0x30A0, 0x30A0, prKatakana},               // Pd       KATAKANA-HIRAGANA DOUBLE HYPHEN
	{0x30A1, 0x30FA, prKatakana},               // Lo  [90] KATAKANA LETTER SMALL A..KATAKANA LETTER VO
	{0x30FC, 0x30FE, prKatakana},               // Lm   [3] KATAKANA-HIRAGANA PROLONGED SOUND MARK..KATAKANA VOICED ITERATION MARK
	{0x30FF, 0x30FF, prKatakana},               // Lo       KATAKANA DIGRAPH KOTO
	{0x3105, 0x312F, prALetter},                // Lo  [43] BOPOMOFO LETTER B..BOPOMOFO LETTER NN
	{0x3131, 0x318E, prALetter},                // Lo  [94] HANGUL LETTER KIYEOK..HANGUL LETTER ARAEAE
	{0x31A0, 0x31BF, prALetter},                // Lo  [32] BOPOMOFO LETTER BU..BOPOMOFO LETTER AH
	{0x31F0, 0x31FF, prKatakana},               // Lo  [16] KATAKANA LETTER SMALL KU..KATAKANA LETTER SMALL RO
	{0x3297, 0x3297, prExtendedPictographic},   // E0.6   [1] (㊗️)       Japanese “congratulations” button
	{0x3299, 0x3299, prExtendedPictographic},   // E0.6   [1] (㊙️)       Japanese “secret” button
	{0x32D0, 0x32FE, prKatakana},               // So  [47] CIRCLED KATAKANA A..CIRCLED KATAKANA WO
	{0x3300, 0x3357, prKatakana},               // So  [88] SQUARE APAATO..SQUARE WATTO
	{0xA000, 0xA014, prALetter},                // Lo  [21] YI SYLLABLE IT..YI SYLLABLE E
	{0xA015, 0xA015, prALetter},                // Lm       YI SYLLABLE WU
	{0xA016, 0xA48C, prALetter},                // Lo [1143] YI SYLLABLE BIT..YI SYLLABLE YYR
	{0xA4D0, 0xA4F7, prALetter},                // Lo  [40] LISU LETTER BA..LISU LETTER OE
	{0xA4F8, 0xA4FD, prALetter},                // Lm   [6] LISU LETTER TONE MYA TI..LISU LETTER TONE MYA JEU
	{0xA500, 0xA60B, prALetter},                // Lo [268] VAI SYLLABLE EE..VAI SYLLABLE NG
	{0xA60C, 0xA60C, prALetter},                // Lm       VAI SYLLABLE LENGTHENER
	{0xA610, 0xA61F, prALetter},                // Lo  [16] VAI SYLLABLE NDOLE FA..VAI SYMBOL JONG
	{0xA620, 0xA629, prNumeric},                // Nd  [10] VAI DIGIT ZERO..VAI DIGIT NINE
	{0xA62A, 0xA62B, prALetter},                // Lo   [2] VAI SYLLABLE NDOLE MA..VAI SYLLABLE NDOLE DO
	{0xA640, 0xA66D, prALetter},                // L&  [46] CYRILLIC CAPITAL LETTER ZEMLYA..CYRILLIC SMALL LETTER DOUBLE MONOCULAR O
	{0xA66E, 0xA66E, prALetter},                // Lo       CYRILLIC LETTER MULTIOCULAR O
	{0xA66F, 0xA66F, prExtend},                 // Mn       COMBINING CYRILLIC VZMET
	{0xA670, 0xA672, prExtend},                 // Me   [3] COMBINING CYRILLIC TEN MILLIONS SIGN..COMBINING CYRILLIC THOUSAND MILLIONS SIGN
	{0xA674, 0xA67D, prExtend},                 // Mn  [10] COMBINING CYRILLIC LETTER UKRAINIAN IE..COMBINING CYRILLIC PAYEROK
	{0xA67F, 0xA67F, prALetter},                // Lm       CYRILLIC PAYEROK
	{0xA680, 0xA69B, prALetter},                // L&  [28] CYRILLIC CAPITAL LETTER DWE..CYRILLIC SMALL LETTER CROSSED O
	{0xA69C, 0xA69D, prALetter},                // Lm   [2] MODIFIER LETTER CYRILLIC HARD SIGN..MODIFIER LETTER CYRILLIC SOFT SIGN
	{0xA69E, 0xA69F, prExtend},                 // Mn   [2] COMBINING CYRILLIC LETTER EF..COMBINING CYRILLIC LETTER IOTIFIED E
	{0xA6A0, 0xA6E5, prALetter},                // Lo  [70] BAMUM LETTER A..BAMUM LETTER KI
	{0xA6E6, 0xA6EF, prALetter},                // Nl  [10] BAMUM LETTER MO..BAMUM LETTER KOGHOM
	{0xA6F0, 0xA6F1, prExtend},                 // Mn   [2] BAMUM COMBINING MARK KOQNDON..BAMUM COMBINING MARK TUKWENTIS
	{0xA708, 0xA716, prALetter},                // Sk  [15] MODIFIER LETTER EXTRA-HIGH DOTTED TONE BAR..MODIFIER LETTER EXTRA-LOW LEFT-STEM TONE BAR
	{0xA717, 0xA71F, prALetter},                // Lm   [9] MODIFIER LETTER DOT VERTICAL BAR..MODIFIER LETTER LOW INVERTED EXCLAMATION MARK
	{0xA720, 0xA721, prALetter},                // Sk   [2] MODIFIER LETTER STRESS AND HIGH TONE..MODIFIER LETTER STRESS AND LOW TONE
	{0xA722, 0xA76F, prALetter},                // L&  [78] LATIN CAPITAL LETTER EGYPTOLOGICAL ALEF..LATIN SMALL LETTER CON
	{0xA770, 0xA770, prALetter},                // Lm       MODIFIER LETTER US
	{0xA771, 0xA787, prALetter},                // L&  [23] LATIN SMALL LETTER DUM..LATIN SMALL LETTER INSULAR T
	{0xA788, 0xA788, prALetter},                // Lm       MODIFIER LETTER LOW CIRCUMFLEX ACCENT
	{0xA789, 0xA78A, prALetter},                // Sk   [2] MODIFIER LETTER COLON..MODIFIER LETTER SHORT EQUALS SIGN
	{0xA78B, 0xA78E, prALetter},                // L&   [4] LATIN CAPITAL LETTER SALTILLO..LATIN SMALL LETTER L WITH RETROFLEX HOOK AND BELT
	{0xA78F, 0xA78F, prALetter},                // Lo       LATIN LETTER SINOLOGICAL DOT
	{0xA790, 0xA7CA, prALetter},                // L&  [59] LATIN CAPITAL LETTER N WITH DESCENDER..LATIN SMALL LETTER S WITH SHORT STROKE OVERLAY
	{0xA7D0, 0xA7D1, prALetter},                // L&   [2] LATIN CAPITAL LETTER CLOSED INSULAR G..LATIN SMALL LETTER CLOSED INSULAR G
	{0xA7D3, 0xA7D3, prALetter},                // L&       LATIN SMALL LETTER DOUBLE THORN
	{0xA7D5, 0xA7D9, prALetter},                // L&   [5] LATIN SMALL LETTER DOUBLE WYNN..LATIN SMALL LETTER SIGMOID S
	{0xA7F2, 0xA7F4, prALetter},                // Lm   [3] MODIFIER LETTER CAPITAL C..MODIFIER LETTER CAPITAL Q
	{0xA7F5, 0xA7F6, prALetter},                // L&   [2] LATIN CAPITAL LETTER REVERSED HALF H..LATIN SMALL LETTER REVERSED HALF H
	{0xA7F7, 0xA7F7, prALetter},                // Lo       LATIN EPIGRAPHIC LETTER SIDEWAYS I
	{0xA7F8, 0xA7F9, prALetter},                // Lm   [2] MODIFIER LETTER CAPITAL H WITH STROKE..MODIFIER LETTER SMALL LIGATURE OE
	{0xA7FA, 0xA7FA, prALetter},                // L&       LATIN LETTER SMALL CAPITAL TURNED M
	{0xA7FB, 0xA801, prALetter},                // Lo   [7] LATIN EPIGRAPHIC LETTER REVERSED F..SYLOTI NAGRI LETTER I
	{0xA802, 0xA802, prExtend},                 // Mn       SYLOTI NAGRI SIGN DVISVARA
	{0xA803, 0xA805, prALetter},                // Lo   [3] SYLOTI NAGRI LETTER U..SYLOTI NAGRI LETTER O
	{0xA806, 0xA806, prExtend},                 // Mn       SYLOTI NAGRI SIGN HASANTA
	{0xA807, 0xA80A, prALetter},                // Lo   [4] SYLOTI NAGRI LETTER KO..SYLOTI NAGRI LETTER GHO
	{0xA80B, 0xA80B, prExtend},                 // Mn       SYLOTI NAGRI SIGN ANUSVARA
	{0xA80C, 0xA822, prALetter},                // Lo  [23] SYLOTI NAGRI LETTER CO..SYLOTI NAGRI LETTER HO
	{0xA823, 0xA824, prExtend},                 // Mc   [2] SYLOTI NAGRI VOWEL SIGN A..SYLOTI NAGRI VOWEL SIGN I
	{0xA825, 0xA826, prExtend},                 // Mn   [2] SYLOTI NAGRI VOWEL SIGN U..SYLOTI NAGRI VOWEL SIGN E
	{0xA827, 0xA827, prExtend},                 // Mc       SYLOTI NAGRI VOWEL SIGN OO
	{0xA82C, 0xA82C, prExtend},                 // Mn       SYLOTI NAGRI SIGN ALTERNATE HASANTA
	{0xA840, 0xA873, prALetter},                // Lo  [52] PHAGS-PA LETTER KA..PHAGS-PA LETTER CANDRABINDU
	{0xA880, 0xA881, prExtend},                 // Mc   [2] SAURASHTRA SIGN ANUSVARA..SAURASHTRA SIGN VISARGA
	{0xA882, 0xA8B3, prALetter},                // Lo  [50] SAURASHTRA LETTER A..SAURASHTRA LETTER LLA
	{0xA8B4, 0xA8C3, prExtend},                 // Mc  [16] SAURASHTRA CONSONANT SIGN HAARU..SAURASHTRA VOWEL SIGN AU
	{0xA8C4, 0xA8C5, prExtend},                 // Mn   [2] SAURASHTRA SIGN VIRAMA..SAURASHTRA SIGN CANDRABINDU
	{0xA8D0, 0xA8D9, prNumeric},                // Nd  [10] SAURASHTRA DIGIT ZERO..SAURASHTRA DIGIT NINE
	{0xA8E0, 0xA8F1, prExtend},                 // Mn  [18] COMBINING DEVANAGARI DIGIT ZERO..COMBINING DEVANAGARI SIGN AVAGRAHA
	{0xA8F2, 0xA8F7, prALetter},                // Lo   [6] DEVANAGARI SIGN SPACING CANDRABINDU..DEVANAGARI SIGN CANDRABINDU AVAGRAHA
	{0xA8FB, 0xA8FB, prALetter},                // Lo       DEVANAGARI HEADSTROKE
	{0xA8FD, 0xA8FE, prALetter},                // Lo   [2] DEVANAGARI JAIN OM..DEVANAGARI LETTER AY
	{0xA8FF, 0xA8FF, prExtend},                 // Mn       DEVANAGARI VOWEL SIGN AY
	{0xA900, 0xA909, prNumeric},                // Nd  [10] KAYAH LI DIGIT ZERO..KAYAH LI DIGIT NINE
	{0xA90A, 0xA925, prALetter},                // Lo  [28] KAYAH LI LETTER KA..KAYAH LI LETTER OO
	{0xA926, 0xA92D, prExtend},                 // Mn   [8] KAYAH LI VOWEL UE..KAYAH LI TONE CALYA PLOPHU
	{0xA930, 0xA946, prALetter},                // Lo  [23] REJANG LETTER KA..REJANG LETTER A
	{0xA947, 0xA951, prExtend},                 // Mn  [11] REJANG VOWEL SIGN I..REJANG CONSONANT SIGN R
	{0xA952, 0xA953, prExtend},                 // Mc   [2] REJANG CONSONANT SIGN H..REJANG VIRAMA
	{0xA960, 0xA97C, prALetter},                // Lo  [29] HANGUL CHOSEONG TIKEUT-MIEUM..HANGUL CHOSEONG SSANGYEORINHIEUH
	{0xA980, 0xA982, prExtend},                 // Mn   [3] JAVANESE SIGN PANYANGGA..JAVANESE SIGN LAYAR
	{0xA983, 0xA983, prExtend},                 // Mc       JAVANESE SIGN WIGNYAN
	{0xA984, 0xA9B2, prALetter},                // Lo  [47] JAVANESE LETTER A..JAVANESE LETTER HA
	{0xA9B3, 0xA9B3, prExtend},                 // Mn       JAVANESE SIGN CECAK TELU
	{0xA9B4, 0xA9B5, prExtend},                 // Mc   [2] JAVANESE VOWEL SIGN TARUNG..JAVANESE VOWEL SIGN TOLONG
	{0xA9B6, 0xA9B9, prExtend},                 // Mn   [4] JAVANESE VOWEL SIGN WULU..JAVANESE VOWEL SIGN SUKU MENDUT
	{0xA9BA, 0xA9BB, prExtend},                 // Mc   [2] JAVANESE VOWEL SIGN TALING..JAVANESE VOWEL SIGN DIRGA MURE
	{0xA9BC, 0xA9BD, prExtend},                 // Mn   [2] JAVANESE VOWEL SIGN PEPET..JAVANESE CONSONANT SIGN KERET
	{0xA9BE, 0xA9C0, prExtend},                 // Mc   [3] JAVANESE CONSONANT SIGN PENGKAL..JAVANESE PANGKON
	{0xA9CF, 0xA9CF, prALetter},                // Lm       JAVANESE PANGRANGKEP
	{0xA9D0, 0xA9D9, prNumeric},                // Nd  [10] JAVANESE DIGIT ZERO..JAVANESE DIGIT NINE
	{0xA9E5, 0xA9E5, prExtend},                 // Mn       MYANMAR SIGN SHAN SAW
	{0xA9F0, 0xA9F9, prNumeric},                // Nd  [10] MYANMAR TAI LAING DIGIT ZERO..MYANMAR TAI LAING DIGIT NINE
	{0xAA00, 0xAA28, prALetter},                // Lo  [41] CHAM LETTER A..CHAM LETTER HA
	{0xAA29, 0xAA2E, prExtend},                 // Mn   [6] CHAM VOWEL SIGN AA..CHAM VOWEL SIGN OE
	{0xAA2F, 0xAA30, prExtend},                 // Mc   [2] CHAM VOWEL SIGN O..CHAM VOWEL SIGN AI
	{0xAA31, 0xAA32, prExtend},                 // Mn   [2] CHAM VOWEL SIGN AU..CHAM VOWEL SIGN UE
	{0xAA33, 0xAA34, prExtend},                 // Mc   [2] CHAM CONSONANT SIGN YA..CHAM CONSONANT SIGN RA
	{0xAA35, 0xAA36, prExtend},                 // Mn   [2] CHAM CONSONANT SIGN LA..CHAM CONSONANT SIGN WA
	{0xAA40, 0xAA42, prALetter},                // Lo   [3] CHAM LETTER FINAL K..CHAM LETTER FINAL NG
	{0xAA43, 0xAA43, prExtend},                 // Mn       CHAM CONSONANT SIGN FINAL NG
	{0xAA44, 0xAA4B, prALetter},                // Lo   [8] CHAM LETTER FINAL CH..CHAM LETTER FINAL SS
	{0xAA4C, 0xAA4C, prExtend},                 // Mn       CHAM CONSONANT SIGN FINAL M
	{0xAA4D, 0xAA4D, prExtend},                 // Mc       CHAM CONSONANT SIGN FINAL H
	{0xAA50, 0xAA59, prNumeric},                // Nd  [10] CHAM DIGIT ZERO..CHAM DIGIT NINE
	{0xAA7B, 0xAA7B, prExtend},                 // Mc       MYANMAR SIGN PAO KAREN TONE
	{0xAA7C, 0xAA7C, prExtend},                 // Mn       MYANMAR SIGN TAI LAING TONE-2
	{0xAA7D, 0xAA7D, prExtend},                 // Mc       MYANMAR SIGN TAI LAING TONE-5
	{0xAAB0, 0xAAB0, prExtend},                 // Mn       TAI VIET MAI KANG
	{0xAAB2, 0xAAB4, prExtend},                 // Mn   [3] TAI VIET VOWEL I..TAI VIET VOWEL U
	{0xAAB7, 0xAAB8, prExtend},                 // Mn   [2] TAI VIET MAI KHIT..TAI VIET VOWEL IA
	{0xAABE, 0xAABF, prExtend},                 // Mn   [2] TAI VIET VOWEL AM..TAI VIET TONE MAI EK
	{0xAAC1, 0xAAC1, prExtend},                 // Mn       TAI VIET TONE MAI THO
	{0xAAE0, 0xAAEA, prALetter},                // Lo  [11] MEETEI MAYEK LETTER E..MEETEI MAYEK LETTER SSA
	{0xAAEB, 0xAAEB, prExtend},                 // Mc       MEETEI MAYEK VOWEL SIGN II
	{0xAAEC, 0xAAED, prExtend},                 // Mn   [2] MEETEI MAYEK VOWEL SIGN UU..MEETEI MAYEK VOWEL SIGN AAI
	{0xAAEE, 0xAAEF, prExtend},                 // Mc   [2] MEETEI MAYEK VOWEL SIGN AU..MEETEI MAYEK VOWEL SIGN AAU
	{0xAAF2, 0xAAF2, prALetter},                // Lo       MEETEI MAYEK ANJI
	{0xAAF3, 0xAAF4, prALetter},                // Lm   [2] MEETEI MAYEK SYLLABLE REPETITION MARK..MEETEI MAYEK WORD REPETITION MARK
	{0xAAF5, 0xAAF5, prExtend},                 // Mc       MEETEI MAYEK VOWEL SIGN VISARGA
	{0xAAF6, 0xAAF6, prExtend},                 // Mn       MEETEI MAYEK VIRAMA
	{0xAB01, 0xAB06, prALetter},                // Lo   [6] ETHIOPIC SYLLABLE TTHU..ETHIOPIC SYLLABLE TTHO
	{0xAB09, 0xAB0E, prALetter},                // Lo   [6] ETHIOPIC SYLLABLE DDHU..ETHIOPIC SYLLABLE DDHO
	{0xAB11, 0xAB16, prALetter},                // Lo   [6] ETHIOPIC SYLLABLE DZU..ETHIOPIC SYLLABLE DZO
	{0xAB20, 0xAB26, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE CCHHA..ETHIOPIC SYLLABLE CCHHO
	{0xAB28, 0xAB2E, prALetter},                // Lo   [7] ETHIOPIC SYLLABLE BBA..ETHIOPIC SYLLABLE BBO
	{0xAB30, 0xAB5A, prALetter},                // L&  [43] LATIN SMALL LETTER BARRED ALPHA..LATIN SMALL LETTER Y WITH SHORT RIGHT LEG
	{0xAB5B, 0xAB5B, prALetter},                // Sk       MODIFIER BREVE WITH INVERTED BREVE
	{0xAB5C, 0xAB5F, prALetter},                // Lm   [4] MODIFIER LETTER SMALL HENG..MODIFIER LETTER SMALL U WITH LEFT HOOK
	{0xAB60, 0xAB68, prALetter},                // L&   [9] LATIN SMALL LETTER SAKHA YAT..LATIN SMALL LETTER TURNED R WITH MIDDLE TILDE
	{0xAB69, 0xAB69, prALetter},                // Lm       MODIFIER LETTER SMALL TURNED W
	{0xAB70, 0xABBF, prALetter},                // L&  [80] CHEROKEE SMALL LETTER A..CHEROKEE SMALL LETTER YA
	{0xABC0, 0xABE2, prALetter},                // Lo  [35] MEETEI MAYEK LETTER KOK..MEETEI MAYEK LETTER I LONSUM
	{0xABE3, 0xABE4, prExtend},                 // Mc   [2] MEETEI MAYEK VOWEL SIGN ONAP..MEETEI MAYEK VOWEL SIGN INAP
	{0xABE5, 0xABE5, prExtend},                 // Mn       MEETEI MAYEK VOWEL SIGN ANAP
	{0xABE6, 0xABE7, prExtend},                 // Mc   [2] MEETEI MAYEK VOWEL SIGN YENAP..MEETEI MAYEK VOWEL SIGN SOUNAP
	{0xABE8, 0xABE8, prExtend},                 // Mn       MEETEI MAYEK VOWEL SIGN UNAP
	{0xABE9, 0xABEA, prExtend},                 // Mc   [2] MEETEI MAYEK VOWEL SIGN CHEINAP..MEETEI MAYEK VOWEL SIGN NUNG
	{0xABEC, 0xABEC, prExtend},                 // Mc       MEETEI MAYEK LUM IYEK
	{0xABED, 0xABED, prExtend},                 // Mn       MEETEI MAYEK APUN IYEK
	{0xABF0, 0xABF9, prNumeric},                // Nd  [10] MEETEI MAYEK DIGIT ZERO..MEETEI MAYEK DIGIT NINE
	{0xAC00, 0xD7A3, prALetter},                // Lo [11172] HANGUL SYLLABLE GA..HANGUL SYLLABLE HIH
	{0xD7B0, 0xD7C6, prALetter},                // Lo  [23] HANGUL JUNGSEONG O-YEO..HANGUL JUNGSEONG ARAEA-E
	{0xD7CB, 0xD7FB, prALetter},                // Lo  [49] HANGUL JONGSEONG NIEUN-RIEUL..HANGUL JONGSEONG PHIEUPH-THIEUTH
	{0xFB00, 0xFB06, prALetter},                // L&   [7] LATIN SMALL LIGATURE FF..LATIN SMALL LIGATURE ST
	{0xFB13, 0xFB17, prALetter},                // L&   [5] ARMENIAN SMALL LIGATURE MEN NOW..ARMENIAN SMALL LIGATURE MEN XEH
	{0xFB1D, 0xFB1D, prHebrewLetter},           // Lo       HEBREW LETTER YOD WITH HIRIQ
	{0xFB1E, 0xFB1E, prExtend},                 // Mn       HEBREW POINT JUDEO-SPANISH VARIKA
	{0xFB1F, 0xFB28, prHebrewLetter},           // Lo  [10] HEBREW LIGATURE YIDDISH YOD YOD PATAH..HEBREW LETTER WIDE TAV
	{0xFB2A, 0xFB36, prHebrewLetter},           // Lo  [13] HEBREW LETTER SHIN WITH SHIN DOT..HEBREW LETTER ZAYIN WITH DAGESH
	{0xFB38, 0xFB3C, prHebrewLetter},           // Lo   [5] HEBREW LETTER TET WITH DAGESH..HEBREW LETTER LAMED WITH DAGESH
	{0xFB3E, 0xFB3E, prHebrewLetter},           // Lo       HEBREW LETTER MEM WITH DAGESH
	{0xFB40, 0xFB41, prHebrewLetter},           // Lo   [2] HEBREW LETTER NUN WITH DAGESH..HEBREW LETTER SAMEKH WITH DAGESH
	{0xFB43, 0xFB44, prHebrewLetter},           // Lo   [2] HEBREW LETTER FINAL PE WITH DAGESH..HEBREW LETTER PE WITH DAGESH
	{0xFB46, 0xFB4F, prHebrewLetter},           // Lo  [10] HEBREW LETTER TSADI WITH DAGESH..HEBREW LIGATURE ALEF LAMED
	{0xFB50, 0xFBB1, prALetter},                // Lo  [98] ARABIC LETTER ALEF WASLA ISOLATED FORM..ARABIC LETTER YEH BARREE WITH HAMZA ABOVE FINAL FORM
	{0xFBD3, 0xFD3D, prALetter},                // Lo [363] ARABIC LETTER NG ISOLATED FORM..ARABIC LIGATURE ALEF WITH FATHATAN ISOLATED FORM
	{0xFD50, 0xFD8F, prALetter},                // Lo  [64] ARABIC LIGATURE TEH WITH JEEM WITH MEEM INITIAL FORM..ARABIC LIGATURE MEEM WITH KHAH WITH MEEM INITIAL FORM
	{0xFD92, 0xFDC7, prALetter},                // Lo  [54] ARABIC LIGATURE MEEM WITH JEEM WITH KHAH INITIAL FORM..ARABIC LIGATURE NOON WITH JEEM WITH YEH FINAL FORM
	{0xFDF0, 0xFDFB, prALetter},                // Lo  [12] ARABIC LIGATURE SALLA USED AS KORANIC STOP SIGN ISOLATED FORM..ARABIC LIGATURE JALLAJALALOUHOU
	{0xFE00, 0xFE0F, prExtend},                 // Mn  [16] VARIATION SELECTOR-1..VARIATION SELECTOR-16
	{0xFE10, 0xFE10, prMidNum},                 // Po       PRESENTATION FORM FOR VERTICAL COMMA
	{0xFE13, 0xFE13, prMidLetter},              // Po       PRESENTATION FORM FOR VERTICAL COLON
	{0xFE14, 0xFE14, prMidNum},                 // Po       PRESENTATION FORM FOR VERTICAL SEMICOLON
	{0xFE20, 0xFE2F, prExtend},                 // Mn  [16] COMBINING LIGATURE LEFT HALF..COMBINING CYRILLIC TITLO RIGHT HALF
	{0xFE33, 0xFE34, prExtendNumLet},           // Pc   [2] PRESENTATION FORM FOR VERTICAL LOW LINE..PRESENTATION FORM FOR VERTICAL WAVY LOW LINE
	{0xFE4D, 0xFE4F, prExtendNumLet},           // Pc   [3] DASHED LOW LINE..WAVY LOW LINE
	{0xFE50, 0xFE50, prMidNum},                 // Po       SMALL COMMA
	{0xFE52, 0xFE52, prMidNumLet},              // Po       SMALL FULL STOP
	{0xFE54, 0xFE54, prMidNum},                 // Po       SMALL SEMICOLON
	{0xFE55, 0xFE55, prMidLetter},              // Po       SMALL COLON
	{0xFE70, 0xFE74, prALetter},                // Lo   [5] ARABIC FATHATAN ISOLATED FORM..ARABIC KASRATAN ISOLATED FORM
	{0xFE76, 0xFEFC, prALetter},                // Lo [135] ARABIC FATHA ISOLATED FORM..ARABIC LIGATURE LAM WITH ALEF FINAL FORM
	{0xFEFF, 0xFEFF, prFormat},                 // Cf       ZERO WIDTH NO-BREAK SPACE
	{0xFF07, 0xFF07, prMidNumLet},              // Po       FULLWIDTH APOSTROPHE
	{0xFF0C, 0xFF0C, prMidNum},                 // Po       FULLWIDTH COMMA
	{0xFF0E, 0xFF0E, prMidNumLet},              // Po       FULLWIDTH FULL STOP
	{0xFF10, 0xFF19, prNumeric},                // Nd  [10] FULLWIDTH DIGIT ZERO..FULLWIDTH DIGIT NINE
	{0xFF1A, 0xFF1A, prMidLetter},              // Po       FULLWIDTH COLON
	{0xFF1B, 0xFF1B, prMidNum},                 // Po       FULLWIDTH SEMICOLON
	{0xFF21, 0xFF3A, prALetter},                // L&  [26] FULLWIDTH LATIN CAPITAL LETTER A..FULLWIDTH LATIN CAPITAL LETTER Z
	{0xFF3F, 0xFF3F, prExtendNumLet},           // Pc       FULLWIDTH LOW LINE
	{0xFF41, 0xFF5A, prALetter},                // L&  [26] FULLWIDTH LATIN SMALL LETTER A..FULLWIDTH LATIN SMALL LETTER Z
	{0xFF66, 0xFF6F, prKatakana},               // Lo  [10] HALFWIDTH KATAKANA LETTER WO..HALFWIDTH KATAKANA LETTER SMALL TU
	{0xFF70, 0xFF70, prKatakana},               // Lm       HALFWIDTH KATAKANA-HIRAGANA PROLONGED SOUND MARK
	{0xFF71, 0xFF9D, prKatakana},               // Lo  [45] HALFWIDTH KATAKANA LETTER A..HALFWIDTH KATAKANA LETTER N
	{0xFF9E, 0xFF9F, prExtend},                 // Lm   [2] HALFWIDTH KATAKANA VOICED SOUND MARK..HALFWIDTH KATAKANA SEMI-VOICED SOUND MARK
	{0xFFA0, 0xFFBE, prALetter},                // Lo  [31] HALFWIDTH HANGUL FILLER..HALFWIDTH HANGUL LETTER HIEUH
	{0xFFC2, 0xFFC7, prALetter},                // Lo   [6] HALFWIDTH HANGUL LETTER A..HALFWIDTH HANGUL LETTER E
	{0xFFCA, 0xFFCF, prALetter},                // Lo   [6] HALFWIDTH HANGUL LETTER YEO..HALFWIDTH HANGUL LETTER OE
	{0xFFD2, 0xFFD7, prALetter},                // Lo   [6] HALFWIDTH HANGUL LETTER YO..HALFWIDTH HANGUL LETTER YU
	{0xFFDA, 0xFFDC, prALetter},                // Lo   [3] HALFWIDTH HANGUL LETTER EU..HALFWIDTH HANGUL LETTER I
	{0xFFF9, 0xFFFB, prFormat},                 // Cf   [3] INTERLINEAR ANNOTATION ANCHOR..INTERLINEAR ANNOTATION TERMINATOR
	{0x10000, 0x1000B, prALetter},              // Lo  [12] LINEAR B SYLLABLE B008 A..LINEAR B SYLLABLE B046 JE
	{0x1000D, 0x10026, prALetter},              // Lo  [26] LINEAR B SYLLABLE B036 JO..LINEAR B SYLLABLE B032 QO
	{0x10028, 0x1003A, prALetter},              // Lo  [19] LINEAR B SYLLABLE B060 RA..LINEAR B SYLLABLE B042 WO
	{0x1003C, 0x1003D, prALetter},              // Lo   [2] LINEAR B SYLLABLE B017 ZA..LINEAR B SYLLABLE B074 ZE
	{0x1003F, 0x1004D, prALetter},              // Lo  [15] LINEAR B SYLLABLE B020 ZO..LINEAR B SYLLABLE B091 TWO
	{0x10050, 0x1005D, prALetter},              // Lo  [14] LINEAR B SYMBOL B018..LINEAR B SYMBOL B089
	{0x10080, 0x100FA, prALetter},              // Lo [123] LINEAR B IDEOGRAM B100 MAN..LINEAR B IDEOGRAM VESSEL B305
	{0x10140, 0x10174, prALetter},              // Nl  [53] GREEK ACROPHONIC ATTIC ONE QUARTER..GREEK ACROPHONIC STRATIAN FIFTY MNAS
	{0x101FD, 0x101FD, prExtend},               // Mn       PHAISTOS DISC SIGN COMBINING OBLIQUE STROKE
	{0x10280, 0x1029C, prALetter},              // Lo  [29] LYCIAN LETTER A..LYCIAN LETTER X
	{0x102A0, 0x102D0, prALetter},              // Lo  [49] CARIAN LETTER A..CARIAN LETTER UUU3
	{0x102E0, 0x102E0, prExtend},               // Mn       COPTIC EPACT THOUSANDS MARK
	{0x10300, 0x1031F, prALetter},              // Lo  [32] OLD ITALIC LETTER A..OLD ITALIC LETTER ESS
	{0x1032D, 0x10340, prALetter},              // Lo  [20] OLD ITALIC LETTER YE..GOTHIC LETTER PAIRTHRA
	{0x10341, 0x10341, prALetter},              // Nl       GOTHIC LETTER NINETY
	{0x10342, 0x10349, prALetter},              // Lo   [8] GOTHIC LETTER RAIDA..GOTHIC LETTER OTHAL
	{0x1034A, 0x1034A, prALetter},              // Nl       GOTHIC LETTER NINE HUNDRED
	{0x10350, 0x10375, prALetter},              // Lo  [38] OLD PERMIC LETTER AN..OLD PERMIC LETTER IA
	{0x10376, 0x1037A, prExtend},               // Mn   [5] COMBINING OLD PERMIC LETTER AN..COMBINING OLD PERMIC LETTER SII
	{0x10380, 0x1039D, prALetter},              // Lo  [30] UGARITIC LETTER ALPA..UGARITIC LETTER SSU
	{0x103A0, 0x103C3, prALetter},              // Lo  [36] OLD PERSIAN SIGN A..OLD PERSIAN SIGN HA
	{0x103C8, 0x103CF, prALetter},              // Lo   [8] OLD PERSIAN SIGN AURAMAZDAA..OLD PERSIAN SIGN BUUMISH
	{0x103D1, 0x103D5, prALetter},              // Nl   [5] OLD PERSIAN NUMBER ONE..OLD PERSIAN NUMBER HUNDRED
	{0x10400, 0x1044F, prALetter},              // L&  [80] DESERET CAPITAL LETTER LONG I..DESERET SMALL LETTER EW
	{0x10450, 0x1049D, prALetter},              // Lo  [78] SHAVIAN LETTER PEEP..OSMANYA LETTER OO
	{0x104A0, 0x104A9, prNumeric},              // Nd  [10] OSMANYA DIGIT ZERO..OSMANYA DIGIT NINE
	{0x104B0, 0x104D3, prALetter},              // L&  [36] OSAGE CAPITAL LETTER A..OSAGE CAPITAL LETTER ZHA
	{0x104D8, 0x104FB, prALetter},              // L&  [36] OSAGE SMALL LETTER A..OSAGE SMALL LETTER ZHA
	{0x10500, 0x10527, prALetter},              // Lo  [40] ELBASAN LETTER A..ELBASAN LETTER KHE
	{0x10530, 0x10563, prALetter},              // Lo  [52] CAUCASIAN ALBANIAN LETTER ALT..CAUCASIAN ALBANIAN LETTER KIW
	{0x10570, 0x1057A, prALetter},              // L&  [11] VITHKUQI CAPITAL LETTER A..VITHKUQI CAPITAL LETTER GA
	{0x1057C, 0x1058A, prALetter},              // L&  [15] VITHKUQI CAPITAL LETTER HA..VITHKUQI CAPITAL LETTER RE
	{0x1058C, 0x10592, prALetter},              // L&   [7] VITHKUQI CAPITAL LETTER SE..VITHKUQI CAPITAL LETTER XE
	{0x10594, 0x10595, prALetter},              // L&   [2] VITHKUQI CAPITAL LETTER Y..VITHKUQI CAPITAL LETTER ZE
	{0x10597, 0x105A1, prALetter},              // L&  [11] VITHKUQI SMALL LETTER A..VITHKUQI SMALL LETTER GA
	{0x105A3, 0x105B1, prALetter},              // L&  [15] VITHKUQI SMALL LETTER HA..VITHKUQI SMALL LETTER RE
	{0x105B3, 0x105B9, prALetter},              // L&   [7] VITHKUQI SMALL LETTER SE..VITHKUQI SMALL LETTER XE
	{0x105BB, 0x105BC, prALetter},              // L&   [2] VITHKUQI SMALL LETTER Y..VITHKUQI SMALL LETTER ZE
	{0x10600, 0x10736, prALetter},              // Lo [311] LINEAR A SIGN AB001..LINEAR A SIGN A664
	{0x10740, 0x10755, prALetter},              // Lo  [22] LINEAR A SIGN A701 A..LINEAR A SIGN A732 JE
	{0x10760, 0x10767, prALetter},              // Lo   [8] LINEAR A SIGN A800..LINEAR A SIGN A807
	{0x10780, 0x10785, prALetter},              // Lm   [6] MODIFIER LETTER SMALL CAPITAL AA..MODIFIER LETTER SMALL B WITH HOOK
	{0x10787, 0x107B0, prALetter},              // Lm  [42] MODIFIER LETTER SMALL DZ DIGRAPH..MODIFIER LETTER SMALL V WITH RIGHT HOOK
	{0x107B2, 0x107BA, prALetter},              // Lm   [9] MODIFIER LETTER SMALL CAPITAL Y..MODIFIER LETTER SMALL S WITH CURL
	{0x10800, 0x10805, prALetter},              // Lo   [6] CYPRIOT SYLLABLE A..CYPRIOT SYLLABLE JA
	{0x10808, 0x10808, prALetter},              // Lo       CYPRIOT SYLLABLE JO
	{0x1080A, 0x10835, prALetter},              // Lo  [44] CYPRIOT SYLLABLE KA..CYPRIOT SYLLABLE WO
	{0x10837, 0x10838, prALetter},              // Lo   [2] CYPRIOT SYLLABLE XA..CYPRIOT SYLLABLE XE
	{0x1083C, 0x1083C, prALetter},              // Lo       CYPRIOT SYLLABLE ZA
	{0x1083F, 0x10855, prALetter},              // Lo  [23] CYPRIOT SYLLABLE ZO..IMPERIAL ARAMAIC LETTER TAW
	{0x10860, 0x10876, prALetter},              // Lo  [23] PALMYRENE LETTER ALEPH..PALMYRENE LETTER TAW
	{0x10880, 0x1089E, prALetter},              // Lo  [31] NABATAEAN LETTER FINAL ALEPH..NABATAEAN LETTER TAW
	{0x108E0, 0x108F2, prALetter},              // Lo  [19] HATRAN LETTER ALEPH..HATRAN LETTER QOPH
	{0x108F4, 0x108F5, prALetter},              // Lo   [2] HATRAN LETTER SHIN..HATRAN LETTER TAW
	{0x10900, 0x10915, prALetter},              // Lo  [22] PHOENICIAN LETTER ALF..PHOENICIAN LETTER TAU
	{0x10920, 0x10939, prALetter},              // Lo  [26] LYDIAN LETTER A..LYDIAN LETTER C
	{0x10980, 0x109B7, prALetter},              // Lo  [56] MEROITIC HIEROGLYPHIC LETTER A..MEROITIC CURSIVE LETTER DA
	{0x109BE, 0x109BF, prALetter},              // Lo   [2] MEROITIC CURSIVE LOGOGRAM RMT..MEROITIC CURSIVE LOGOGRAM IMN
	{0x10A00, 0x10A00, prALetter},              // Lo       KHAROSHTHI LETTER A
	{0x10A01, 0x10A03, prExtend},               // Mn   [3] KHAROSHTHI VOWEL SIGN I..KHAROSHTHI VOWEL SIGN VOCALIC R
	{0x10A05, 0x10A06, prExtend},               // Mn   [2] KHAROSHTHI VOWEL SIGN E..KHAROSHTHI VOWEL SIGN O
	{0x10A0C, 0x10A0F, prExtend},               // Mn   [4] KHAROSHTHI VOWEL LENGTH MARK..KHAROSHTHI SIGN VISARGA
	{0x10A10, 0x10A13, prALetter},              // Lo   [4] KHAROSHTHI LETTER KA..KHAROSHTHI LETTER GHA
	{0x10A15, 0x10A17, prALetter},              // Lo   [3] KHAROSHTHI LETTER CA..KHAROSHTHI LETTER JA
	{0x10A19, 0x10A35, prALetter},              // Lo  [29] KHAROSHTHI LETTER NYA..KHAROSHTHI LETTER VHA
	{0x10A38, 0x10A3A, prExtend},               // Mn   [3] KHAROSHTHI SIGN BAR ABOVE..KHAROSHTHI SIGN DOT BELOW
	{0x10A3F, 0x10A3F, prExtend},               // Mn       KHAROSHTHI VIRAMA
	{0x10A60, 0x10A7C, prALetter},              // Lo  [29] OLD SOUTH ARABIAN LETTER HE..OLD SOUTH ARABIAN LETTER THETH
	{0x10A80, 0x10A9C, prALetter},              // Lo  [29] OLD NORTH ARABIAN LETTER HEH..OLD NORTH ARABIAN LETTER ZAH
	{0x10AC0, 0x10AC7, prALetter},              // Lo   [8] MANICHAEAN LETTER ALEPH..MANICHAEAN LETTER WAW
	{0x10AC9, 0x10AE4, prALetter},              // Lo  [28] MANICHAEAN LETTER ZAYIN..MANICHAEAN LETTER TAW
	{0x10AE5, 0x10AE6, prExtend},               // Mn   [2] MANICHAEAN ABBREVIATION MARK ABOVE..MANICHAEAN ABBREVIATION MARK BELOW
	{0x10B00, 0x10B35, prALetter},              // Lo  [54] AVESTAN LETTER A..AVESTAN LETTER HE
	{0x10B40, 0x10B55, prALetter},              // Lo  [22] INSCRIPTIONAL PARTHIAN LETTER ALEPH..INSCRIPTIONAL PARTHIAN LETTER TAW
	{0x10B60, 0x10B72, prALetter},              // Lo  [19] INSCRIPTIONAL PAHLAVI LETTER ALEPH..INSCRIPTIONAL PAHLAVI LETTER TAW
	{0x10B80, 0x10B91, prALetter},              // Lo  [18] PSALTER PAHLAVI LETTER ALEPH..PSALTER PAHLAVI LETTER TAW
	{0x10C00, 0x10C48, prALetter},              // Lo  [73] OLD TURKIC LETTER ORKHON A..OLD TURKIC LETTER ORKHON BASH
	{0x10C80, 0x10CB2, prALetter},              // L&  [51] OLD HUNGARIAN CAPITAL LETTER A..OLD HUNGARIAN CAPITAL LETTER US
	{0x10CC0, 0x10CF2, prALetter},              // L&  [51] OLD HUNGARIAN SMALL LETTER A..OLD HUNGARIAN SMALL LETTER US
	{0x10D00, 0x10D23, prALetter},              // Lo  [36] HANIFI ROHINGYA LETTER A..HANIFI ROHINGYA MARK NA KHONNA
	{0x10D24, 0x10D27, prExtend},               // Mn   [4] HANIFI ROHINGYA SIGN HARBAHAY..HANIFI ROHINGYA SIGN TASSI
	{0x10D30, 0x10D39, prNumeric},              // Nd  [10] HANIFI ROHINGYA DIGIT ZERO..HANIFI ROHINGYA DIGIT NINE
	{0x10E80, 0x10EA9, prALetter},              // Lo  [42] YEZIDI LETTER ELIF..YEZIDI LETTER ET
	{0x10EAB, 0x10EAC, prExtend},               // Mn   [2] YEZIDI COMBINING HAMZA MARK..YEZIDI COMBINING MADDA MARK
	{0x10EB0, 0x10EB1, prALetter},              // Lo   [2] YEZIDI LETTER LAM WITH DOT ABOVE..YEZIDI LETTER YOT WITH CIRCUMFLEX ABOVE
	{0x10EFD, 0x10EFF, prExtend},               // Mn   [3] ARABIC SMALL LOW WORD SAKTA..ARABIC SMALL LOW WORD MADDA
	{0x10F00, 0x10F1C, prALetter},              // Lo  [29] OLD SOGDIAN LETTER ALEPH..OLD SOGDIAN LETTER FINAL TAW WITH VERTICAL TAIL
	{0x10F27, 0x10F27, prALetter},              // Lo       OLD SOGDIAN LIGATURE AYIN-DALETH
	{0x10F30, 0x10F45, prALetter},              // Lo  [22] SOGDIAN LETTER ALEPH..SOGDIAN INDEPENDENT SHIN
	{0x10F46, 0x10F50, prExtend},               // Mn  [11] SOGDIAN COMBINING DOT BELOW..SOGDIAN COMBINING STROKE BELOW
	{0x10F70, 0x10F81, prALetter},              // Lo  [18] OLD UYGHUR LETTER ALEPH..OLD UYGHUR LETTER LESH
	{0x10F82, 0x10F85, prExtend},               // Mn   [4] OLD UYGHUR COMBINING DOT ABOVE..OLD UYGHUR COMBINING TWO DOTS BELOW
	{0x10FB0, 0x10FC4, prALetter},              // Lo  [21] CHORASMIAN LETTER ALEPH..CHORASMIAN LETTER TAW
	{0x10FE0, 0x10FF6, prALetter},              // Lo  [23] ELYMAIC LETTER ALEPH..ELYMAIC LIGATURE ZAYIN-YODH
	{0x11000, 0x11000, prExtend},               // Mc       BRAHMI SIGN CANDRABINDU
	{0x11001, 0x11001, prExtend},               // Mn       BRAHMI SIGN ANUSVARA
	{0x11002, 0x11002, prExtend},               // Mc       BRAHMI SIGN VISARGA
	{0x11003, 0x11037, prALetter},              // Lo  [53] BRAHMI SIGN JIHVAMULIYA..BRAHMI LETTER OLD TAMIL NNNA
	{0x11038, 0x11046, prExtend},               // Mn  [15] BRAHMI VOWEL SIGN AA..BRAHMI VIRAMA
	{0x11066, 0x1106F, prNumeric},              // Nd  [10] BRAHMI DIGIT ZERO..BRAHMI DIGIT NINE
	{0x11070, 0x11070, prExtend},               // Mn       BRAHMI SIGN OLD TAMIL VIRAMA
	{0x11071, 0x11072, prALetter},              // Lo   [2] BRAHMI LETTER OLD TAMIL SHORT E..BRAHMI LETTER OLD TAMIL SHORT O
	{0x11073, 0x11074, prExtend},               // Mn   [2] BRAHMI VOWEL SIGN OLD TAMIL SHORT E..BRAHMI VOWEL SIGN OLD TAMIL SHORT O
	{0x11075, 0x11075, prALetter},              // Lo       BRAHMI LETTER OLD TAMIL LLA
	{0x1107F, 0x11081, prExtend},               // Mn   [3] BRAHMI NUMBER JOINER..KAITHI SIGN ANUSVARA
	{0x11082, 0x11082, prExtend},               // Mc       KAITHI SIGN VISARGA
	{0x11083, 0x110AF, prALetter},              // Lo  [45] KAITHI LETTER A..KAITHI LETTER HA
	{0x110B0, 0x110B2, prExtend},               // Mc   [3] KAITHI VOWEL SIGN AA..KAITHI VOWEL SIGN II
	{0x110B3, 0x110B6, prExtend},               // Mn   [4] KAITHI VOWEL SIGN U..KAITHI VOWEL SIGN AI
	{0x110B7, 0x110B8, prExtend},               // Mc   [2] KAITHI VOWEL SIGN O..KAITHI VOWEL SIGN AU
	{0x110B9, 0x110BA, prExtend},               // Mn   [2] KAITHI SIGN VIRAMA..KAITHI SIGN NUKTA
	{0x110BD, 0x110BD, prFormat},               // Cf       KAITHI NUMBER SIGN
	{0x110C2, 0x110C2, prExtend},               // Mn       KAITHI VOWEL SIGN VOCALIC R
	{0x110CD, 0x110CD, prFormat},               // Cf       KAITHI NUMBER SIGN ABOVE
	{0x110D0, 0x110E8, prALetter},              // Lo  [25] SORA SOMPENG LETTER SAH..SORA SOMPENG LETTER MAE
	{0x110F0, 0x110F9, prNumeric},              // Nd  [10] SORA SOMPENG DIGIT ZERO..SORA SOMPENG DIGIT NINE
	{0x11100, 0x11102, prExtend},               // Mn   [3] CHAKMA SIGN CANDRABINDU..CHAKMA SIGN VISARGA
	{0x11103, 0x11126, prALetter},              // Lo  [36] CHAKMA LETTER AA..CHAKMA LETTER HAA
	{0x11127, 0x1112B, prExtend},               // Mn   [5] CHAKMA VOWEL SIGN A..CHAKMA VOWEL SIGN UU
	{0x1112C, 0x1112C, prExtend},               // Mc       CHAKMA VOWEL SIGN E
	{0x1112D, 0x11134, prExtend},               // Mn   [8] CHAKMA VOWEL SIGN AI..CHAKMA MAAYYAA
	{0x11136, 0x1113F, prNumeric},              // Nd  [10] CHAKMA DIGIT ZERO..CHAKMA DIGIT NINE
	{0x11144, 0x11144, prALetter},              // Lo       CHAKMA LETTER LHAA
	{0x11145, 0x11146, prExtend},               // Mc   [2] CHAKMA VOWEL SIGN AA..CHAKMA VOWEL SIGN EI
	{0x11147, 0x11147, prALetter},              // Lo       CHAKMA LETTER VAA
	{0x11150, 0x11172, prALetter},              // Lo  [35] MAHAJANI LETTER A..MAHAJANI LETTER RRA
	{0x11173, 0x11173, prExtend},               // Mn       MAHAJANI SIGN NUKTA
	{0x11176, 0x11176, prALetter},              // Lo       MAHAJANI LIGATURE SHRI
	{0x11180, 0x11181, prExtend},               // Mn   [2] SHARADA SIGN CANDRABINDU..SHARADA SIGN ANUSVARA
	{0x11182, 0x11182, prExtend},               // Mc       SHARADA SIGN VISARGA
	{0x11183, 0x111B2, prALetter},              // Lo  [48] SHARADA LETTER A..SHARADA LETTER HA
	{0x111B3, 0x111B5, prExtend},               // Mc   [3] SHARADA VOWEL SIGN AA..SHARADA VOWEL SIGN II
	{0x111B6, 0x111BE, prExtend},               // Mn   [9] SHARADA VOWEL SIGN U..SHARADA VOWEL SIGN O
	{0x111BF, 0x111C0, prExtend},               // Mc   [2] SHARADA VOWEL SIGN AU..SHARADA SIGN VIRAMA
	{0x111C1, 0x111C4, prALetter},              // Lo   [4] SHARADA SIGN AVAGRAHA..SHARADA OM
	{0x111C9, 0x111CC, prExtend},               // Mn   [4] SHARADA SANDHI MARK..SHARADA EXTRA SHORT VOWEL MARK
	{0x111CE, 0x111CE, prExtend},               // Mc       SHARADA VOWEL SIGN PRISHTHAMATRA E
	{0x111CF, 0x111CF, prExtend},               // Mn       SHARADA SIGN INVERTED CANDRABINDU
	{0x111D0, 0x111D9, prNumeric},              // Nd  [10] SHARADA DIGIT ZERO..SHARADA DIGIT NINE
	{0x111DA, 0x111DA, prALetter},              // Lo       SHARADA EKAM
	{0x111DC, 0x111DC, prALetter},              // Lo       SHARADA HEADSTROKE
	{0x11200, 0x11211, prALetter},              // Lo  [18] KHOJKI LETTER A..KHOJKI LETTER JJA
	{0x11213, 0x1122B, prALetter},              // Lo  [25] KHOJKI LETTER NYA..KHOJKI LETTER LLA
	{0x1122C, 0x1122E, prExtend},               // Mc   [3] KHOJKI VOWEL SIGN AA..KHOJKI VOWEL SIGN II
	{0x1122F, 0x11231, prExtend},               // Mn   [3] KHOJKI VOWEL SIGN U..KHOJKI VOWEL SIGN AI
	{0x11232, 0x11233, prExtend},               // Mc   [2] KHOJKI VOWEL SIGN O..KHOJKI VOWEL SIGN AU
	{0x11234, 0x11234, prExtend},               // Mn       KHOJKI SIGN ANUSVARA
	{0x11235, 0x11235, prExtend},               // Mc       KHOJKI SIGN VIRAMA
	{0x11236, 0x11237, prExtend},               // Mn   [2] KHOJKI SIGN NUKTA..KHOJKI SIGN SHADDA
	{0x1123E, 0x1123E, prExtend},               // Mn       KHOJKI SIGN SUKUN
	{0x1123F, 0x11240, prALetter},              // Lo   [2] KHOJKI LETTER QA..KHOJKI LETTER SHORT I
	{0x11241, 0x11241, prExtend},               // Mn       KHOJKI VOWEL SIGN VOCALIC R
	{0x11280, 0x11286, prALetter},              // Lo   [7] MULTANI LETTER A..MULTANI LETTER GA
	{0x11288, 0x11288, prALetter},              // Lo       MULTANI LETTER GHA
	{0x1128A, 0x1128D, prALetter},              // Lo   [4] MULTANI LETTER CA..MULTANI LETTER JJA
	{0x1128F, 0x1129D, prALetter},              // Lo  [15] MULTANI LETTER NYA..MULTANI LETTER BA
	{0x1129F, 0x112A8, prALetter},              // Lo  [10] MULTANI LETTER BHA..MULTANI LETTER RHA
	{0x112B0, 0x112DE, prALetter},              // Lo  [47] KHUDAWADI LETTER A..KHUDAWADI LETTER HA
	{0x112DF, 0x112DF, prExtend},               // Mn       KHUDAWADI SIGN ANUSVARA
	{0x112E0, 0x112E2, prExtend},               // Mc   [3] KHUDAWADI VOWEL SIGN AA..KHUDAWADI VOWEL SIGN II
	{0x112E3, 0x112EA, prExtend},               // Mn   [8] KHUDAWADI VOWEL SIGN U..KHUDAWADI SIGN VIRAMA
	{0x112F0, 0x112F9, prNumeric},              // Nd  [10] KHUDAWADI DIGIT ZERO..KHUDAWADI DIGIT NINE
	{0x11300, 0x11301, prExtend},               // Mn   [2] GRANTHA SIGN COMBINING ANUSVARA ABOVE..GRANTHA SIGN CANDRABINDU
	{0x11302, 0x11303, prExtend},               // Mc   [2] GRANTHA SIGN ANUSVARA..GRANTHA SIGN VISARGA
	{0x11305, 0x1130C, prALetter},              // Lo   [8] GRANTHA LETTER A..GRANTHA LETTER VOCALIC L
	{0x1130F, 0x11310, prALetter},              // Lo   [2] GRANTHA LETTER EE..GRANTHA LETTER AI
	{0x11313, 0x11328, prALetter},              // Lo  [22] GRANTHA LETTER OO..GRANTHA LETTER NA
	{0x1132A, 0x11330, prALetter},              // Lo   [7] GRANTHA LETTER PA..GRANTHA LETTER RA
	{0x11332, 0x11333, prALetter},              // Lo   [2] GRANTHA LETTER LA..GRANTHA LETTER LLA
	{0x11335, 0x11339, prALetter},              // Lo   [5] GRANTHA LETTER VA..GRANTHA LETTER HA
	{0x1133B, 0x1133C, prExtend},               // Mn   [2] COMBINING BINDU BELOW..GRANTHA SIGN NUKTA
	{0x1133D, 0x1133D, prALetter},              // Lo       GRANTHA SIGN AVAGRAHA
	{0x1133E, 0x1133F, prExtend},               // Mc   [2] GRANTHA VOWEL SIGN AA..GRANTHA VOWEL SIGN I
	{0x11340, 0x11340, prExtend},               // Mn       GRANTHA VOWEL SIGN II
	{0x11341, 0x11344, prExtend},               // Mc   [4] GRANTHA VOWEL SIGN U..GRANTHA VOWEL SIGN VOCALIC RR
	{0x11347, 0x11348, prExtend},               // Mc   [2] GRANTHA VOWEL SIGN EE..GRANTHA VOWEL SIGN AI
	{0x1134B, 0x1134D, prExtend},               // Mc   [3] GRANTHA VOWEL SIGN OO..GRANTHA SIGN VIRAMA
	{0x11350, 0x11350, prALetter},              // Lo       GRANTHA OM
	{0x11357, 0x11357, prExtend},               // Mc       GRANTHA AU LENGTH MARK
	{0x1135D, 0x11361, prALetter},              // Lo   [5] GRANTHA SIGN PLUTA..GRANTHA LETTER VOCALIC LL
	{0x11362, 0x11363, prExtend},               // Mc   [2] GRANTHA VOWEL SIGN VOCALIC L..GRANTHA VOWEL SIGN VOCALIC LL
	{0x11366, 0x1136C, prExtend},               // Mn   [7] COMBINING GRANTHA DIGIT ZERO..COMBINING GRANTHA DIGIT SIX
	{0x11370, 0x11374, prExtend},               // Mn   [5] COMBINING GRANTHA LETTER A..COMBINING GRANTHA LETTER PA
	{0x11400, 0x11434, prALetter},              // Lo  [53] NEWA LETTER A..NEWA LETTER HA
	{0x11435, 0x11437, prExtend},               // Mc   [3] NEWA VOWEL SIGN AA..NEWA VOWEL SIGN II
	{0x11438, 0x1143F, prExtend},               // Mn   [8] NEWA VOWEL SIGN U..NEWA VOWEL SIGN AI
	{0x11440, 0x11441, prExtend},               // Mc   [2] NEWA VOWEL SIGN O..NEWA VOWEL SIGN AU
	{0x11442, 0x11444, prExtend},               // Mn   [3] NEWA SIGN VIRAMA..NEWA SIGN ANUSVARA
	{0x11445, 0x11445, prExtend},               // Mc       NEWA SIGN VISARGA
	{0x11446, 0x11446, prExtend},               // Mn       NEWA SIGN NUKTA
	{0x11447, 0x1144A, prALetter},              // Lo   [4] NEWA SIGN AVAGRAHA..NEWA SIDDHI
	{0x11450, 0x11459, prNumeric},              // Nd  [10] NEWA DIGIT ZERO..NEWA DIGIT NINE
	{0x1145E, 0x1145E, prExtend},               // Mn       NEWA SANDHI MARK
	{0x1145F, 0x11461, prALetter},              // Lo   [3] NEWA LETTER VEDIC ANUSVARA..NEWA SIGN UPADHMANIYA
	{0x11480, 0x114AF, prALetter},              // Lo  [48] TIRHUTA ANJI..TIRHUTA LETTER HA
	{0x114B0, 0x114B2, prExtend},               // Mc   [3] TIRHUTA VOWEL SIGN AA..TIRHUTA VOWEL SIGN II
	{0x114B3, 0x114B8, prExtend},               // Mn   [6] TIRHUTA VOWEL SIGN U..TIRHUTA VOWEL SIGN VOCALIC LL
	{0x114B9, 0x114B9, prExtend},               // Mc       TIRHUTA VOWEL SIGN E
	{0x114BA, 0x114BA, prExtend},               // Mn       TIRHUTA VOWEL SIGN SHORT E
	{0x114BB, 0x114BE, prExtend},               // Mc   [4] TIRHUTA VOWEL SIGN AI..TIRHUTA VOWEL SIGN AU
	{0x114BF, 0x114C0, prExtend},               // Mn   [2] TIRHUTA SIGN CANDRABINDU..TIRHUTA SIGN ANUSVARA
	{0x114C1, 0x114C1, prExtend},               // Mc       TIRHUTA SIGN VISARGA
	{0x114C2, 0x114C3, prExtend},               // Mn   [2] TIRHUTA SIGN VIRAMA..TIRHUTA SIGN NUKTA
	{0x114C4, 0x114C5, prALetter},              // Lo   [2] TIRHUTA SIGN AVAGRAHA..TIRHUTA GVANG
	{0x114C7, 0x114C7, prALetter},              // Lo       TIRHUTA OM
	{0x114D0, 0x114D9, prNumeric},              // Nd  [10] TIRHUTA DIGIT ZERO..TIRHUTA DIGIT NINE
	{0x11580, 0x115AE, prALetter},              // Lo  [47] SIDDHAM LETTER A..SIDDHAM LETTER HA
	{0x115AF, 0x115B1, prExtend},               // Mc   [3] SIDDHAM VOWEL SIGN AA..SIDDHAM VOWEL SIGN II
	{0x115B2, 0x115B5, prExtend},               // Mn   [4] SIDDHAM VOWEL SIGN U..SIDDHAM VOWEL SIGN VOCALIC RR
	{0x115B8, 0x115BB, prExtend},               // Mc   [4] SIDDHAM VOWEL SIGN E..SIDDHAM VOWEL SIGN AU
	{0x115BC, 0x115BD, prExtend},               // Mn   [2] SIDDHAM SIGN CANDRABINDU..SIDDHAM SIGN ANUSVARA
	{0x115BE, 0x115BE, prExtend},               // Mc       SIDDHAM SIGN VISARGA
	{0x115BF, 0x115C0, prExtend},               // Mn   [2] SIDDHAM SIGN VIRAMA..SIDDHAM SIGN NUKTA
	{0x115D8, 0x115DB, prALetter},              // Lo   [4] SIDDHAM LETTER THREE-CIRCLE ALTERNATE I..SIDDHAM LETTER ALTERNATE U
	{0x115DC, 0x115DD, prExtend},               // Mn   [2] SIDDHAM VOWEL SIGN ALTERNATE U..SIDDHAM VOWEL SIGN ALTERNATE UU
	{0x11600, 0x1162F, prALetter},              // Lo  [48] MODI LETTER A..MODI LETTER LLA
	{0x11630, 0x11632, prExtend},               // Mc   [3] MODI VOWEL SIGN AA..MODI VOWEL SIGN II
	{0x11633, 0x1163A, prExtend},               // Mn   [8] MODI VOWEL SIGN U..MODI VOWEL SIGN AI
	{0x1163B, 0x1163C, prExtend},               // Mc   [2] MODI VOWEL SIGN O..MODI VOWEL SIGN AU
	{0x1163D, 0x1163D, prExtend},               // Mn       MODI SIGN ANUSVARA
	{0x1163E, 0x1163E, prExtend},               // Mc       MODI SIGN VISARGA
	{0x1163F, 0x11640, prExtend},               // Mn   [2] MODI SIGN VIRAMA..MODI SIGN ARDHACANDRA
	{0x11644, 0x11644, prALetter},              // Lo       MODI SIGN HUVA
	{0x11650, 0x11659, prNumeric},              // Nd  [10] MODI DIGIT ZERO..MODI DIGIT NINE
	{0x11680, 0x116AA, prALetter},              // Lo  [43] TAKRI LETTER A..TAKRI LETTER RRA
	{0x116AB, 0x116AB, prExtend},               // Mn       TAKRI SIGN ANUSVARA
	{0x116AC, 0x116AC, prExtend},               // Mc       TAKRI SIGN VISARGA
	{0x116AD, 0x116AD, prExtend},               // Mn       TAKRI VOWEL SIGN AA
	{0x116AE, 0x116AF, prExtend},               // Mc   [2] TAKRI VOWEL SIGN I..TAKRI VOWEL SIGN II
	{0x116B0, 0x116B5, prExtend},               // Mn   [6] TAKRI VOWEL SIGN U..TAKRI VOWEL SIGN AU
	{0x116B6, 0x116B6, prExtend},               // Mc       TAKRI SIGN VIRAMA
	{0x116B7, 0x116B7, prExtend},               // Mn       TAKRI SIGN NUKTA
	{0x116B8, 0x116B8, prALetter},              // Lo       TAKRI LETTER ARCHAIC KHA
	{0x116C0, 0x116C9, prNumeric},              // Nd  [10] TAKRI DIGIT ZERO..TAKRI DIGIT NINE
	{0x1171D, 0x1171F, prExtend},               // Mn   [3] AHOM CONSONANT SIGN MEDIAL LA..AHOM CONSONANT SIGN MEDIAL LIGATING RA
	{0x11720, 0x11721, prExtend},               // Mc   [2] AHOM VOWEL SIGN A..AHOM VOWEL SIGN AA
	{0x11722, 0x11725, prExtend},               // Mn   [4] AHOM VOWEL SIGN I..AHOM VOWEL SIGN UU
	{0x11726, 0x11726, prExtend},               // Mc       AHOM VOWEL SIGN E
	{0x11727, 0x1172B, prExtend},               // Mn   [5] AHOM VOWEL SIGN AW..AHOM SIGN KILLER
	{0x11730, 0x11739, prNumeric},              // Nd  [10] AHOM DIGIT ZERO..AHOM DIGIT NINE
	{0x11800, 0x1182B, prALetter},              // Lo  [44] DOGRA LETTER A..DOGRA LETTER RRA
	{0x1182C, 0x1182E, prExtend},               // Mc   [3] DOGRA VOWEL SIGN AA..DOGRA VOWEL SIGN II
	{0x1182F, 0x11837, prExtend},               // Mn   [9] DOGRA VOWEL SIGN U..DOGRA SIGN ANUSVARA
	{0x11838, 0x11838, prExtend},               // Mc       DOGRA SIGN VISARGA
	{0x11839, 0x1183A, prExtend},               // Mn   [2] DOGRA SIGN VIRAMA..DOGRA SIGN NUKTA
	{0x118A0, 0x118DF, prALetter},              // L&  [64] WARANG CITI CAPITAL LETTER NGAA..WARANG CITI SMALL LETTER VIYO
	{0x118E0, 0x118E9, prNumeric},              // Nd  [10] WARANG CITI DIGIT ZERO..WARANG CITI DIGIT NINE
	{0x118FF, 0x11906, prALetter},              // Lo   [8] WARANG CITI OM..DIVES AKURU LETTER E
	{0x11909, 0x11909, prALetter},              // Lo       DIVES AKURU LETTER O
	{0x1190C, 0x11913, prALetter},              // Lo   [8] DIVES AKURU LETTER KA..DIVES AKURU LETTER JA
	{0x11915, 0x11916, prALetter},              // Lo   [2] DIVES AKURU LETTER NYA..DIVES AKURU LETTER TTA
	{0x11918, 0x1192F, prALetter},              // Lo  [24] DIVES AKURU LETTER DDA..DIVES AKURU LETTER ZA
	{0x11930, 0x11935, prExtend},               // Mc   [6] DIVES AKURU VOWEL SIGN AA..DIVES AKURU VOWEL SIGN E
	{0x11937, 0x11938, prExtend},               // Mc   [2] DIVES AKURU VOWEL SIGN AI..DIVES AKURU VOWEL SIGN O
	{0x1193B, 0x1193C, prExtend},               // Mn   [2] DIVES AKURU SIGN ANUSVARA..DIVES AKURU SIGN CANDRABINDU
	{0x1193D, 0x1193D, prExtend},               // Mc       DIVES AKURU SIGN HALANTA
	{0x1193E, 0x1193E, prExtend},               // Mn       DIVES AKURU VIRAMA
	{0x1193F, 0x1193F, prALetter},              // Lo       DIVES AKURU PREFIXED NASAL SIGN
	{0x11940, 0x11940, prExtend},               // Mc       DIVES AKURU MEDIAL YA
	{0x11941, 0x11941, prALetter},              // Lo       DIVES AKURU INITIAL RA
	{0x11942, 0x11942, prExtend},               // Mc       DIVES AKURU MEDIAL RA
	{0x11943, 0x11943, prExtend},               // Mn       DIVES AKURU SIGN NUKTA
	{0x11950, 0x11959, prNumeric},              // Nd  [10] DIVES AKURU DIGIT ZERO..DIVES AKURU DIGIT NINE
	{0x119A0, 0x119A7, prALetter},              // Lo   [8] NANDINAGARI LETTER A..NANDINAGARI LETTER VOCALIC RR
	{0x119AA, 0x119D0, prALetter},              // Lo  [39] NANDINAGARI LETTER E..NANDINAGARI LETTER RRA
	{0x119D1, 0x119D3, prExtend},               // Mc   [3] NANDINAGARI VOWEL SIGN AA..NANDINAGARI VOWEL SIGN II
	{0x119D4, 0x119D7, prExtend},               // Mn   [4] NANDINAGARI VOWEL SIGN U..NANDINAGARI VOWEL SIGN VOCALIC RR
	{0x119DA, 0x119DB, prExtend},               // Mn   [2] NANDINAGARI VOWEL SIGN E..NANDINAGARI VOWEL SIGN AI
	{0x119DC, 0x119DF, prExtend},               // Mc   [4] NANDINAGARI VOWEL SIGN O..NANDINAGARI SIGN VISARGA
	{0x119E0, 0x119E0, prExtend},               // Mn       NANDINAGARI SIGN VIRAMA
	{0x119E1, 0x119E1, prALetter},              // Lo       NANDINAGARI SIGN AVAGRAHA
	{0x119E3, 0x119E3, prALetter},              // Lo       NANDINAGARI HEADSTROKE
	{0x119E4, 0x119E4, prExtend},               // Mc       NANDINAGARI VOWEL SIGN PRISHTHAMATRA E
	{0x11A00, 0x11A00, prALetter},              // Lo       ZANABAZAR SQUARE LETTER A
	{0x11A01, 0x11A0A, prExtend},               // Mn  [10] ZANABAZAR SQUARE VOWEL SIGN I..ZANABAZAR SQUARE VOWEL LENGTH MARK
	{0x11A0B, 0x11A32, prALetter},              // Lo  [40] ZANABAZAR SQUARE LETTER KA..ZANABAZAR SQUARE LETTER KSSA
	{0x11A33, 0x11A38, prExtend},               // Mn   [6] ZANABAZAR SQUARE FINAL CONSONANT MARK..ZANABAZAR SQUARE SIGN ANUSVARA
	{0x11A39, 0x11A39, prExtend},               // Mc       ZANABAZAR SQUARE SIGN VISARGA
	{0x11A3A, 0x11A3A, prALetter},              // Lo       ZANABAZAR SQUARE CLUSTER-INITIAL LETTER RA
	{0x11A3B, 0x11A3E, prExtend},               // Mn   [4] ZANABAZAR SQUARE CLUSTER-FINAL LETTER YA..ZANABAZAR SQUARE CLUSTER-FINAL LETTER VA
	{0x11A47, 0x11A47, prExtend},               // Mn       ZANABAZAR SQUARE SUBJOINER
	{0x11A50, 0x11A50, prALetter},              // Lo       SOYOMBO LETTER A
	{0x11A51, 0x11A56, prExtend},               // Mn   [6] SOYOMBO VOWEL SIGN I..SOYOMBO VOWEL SIGN OE
	{0x11A57, 0x11A58, prExtend},               // Mc   [2] SOYOMBO VOWEL SIGN AI..SOYOMBO VOWEL SIGN AU
	{0x11A59, 0x11A5B, prExtend},               // Mn   [3] SOYOMBO VOWEL SIGN VOCALIC R..SOYOMBO VOWEL LENGTH MARK
	{0x11A5C, 0x11A89, prALetter},              // Lo  [46] SOYOMBO LETTER KA..SOYOMBO CLUSTER-INITIAL LETTER SA
	{0x11A8A, 0x11A96, prExtend},               // Mn  [13] SOYOMBO FINAL CONSONANT SIGN G..SOYOMBO SIGN ANUSVARA
	{0x11A97, 0x11A97, prExtend},               // Mc       SOYOMBO SIGN VISARGA
	{0x11A98, 0x11A99, prExtend},               // Mn   [2] SOYOMBO GEMINATION MARK..SOYOMBO SUBJOINER
	{0x11A9D, 0x11A9D, prALetter},              // Lo       SOYOMBO MARK PLUTA
	{0x11AB0, 0x11AF8, prALetter},              // Lo  [73] CANADIAN SYLLABICS NATTILIK HI..PAU CIN HAU GLOTTAL STOP FINAL
	{0x11C00, 0x11C08, prALetter},              // Lo   [9] BHAIKSUKI LETTER A..BHAIKSUKI LETTER VOCALIC L
	{0x11C0A, 0x11C2E, prALetter},              // Lo  [37] BHAIKSUKI LETTER E..BHAIKSUKI LETTER HA
	{0x11C2F, 0x11C2F, prExtend},               // Mc       BHAIKSUKI VOWEL SIGN AA
	{0x11C30, 0x11C36, prExtend},               // Mn   [7] BHAIKSUKI VOWEL SIGN I..BHAIKSUKI VOWEL SIGN VOCALIC L
	{0x11C38, 0x11C3D, prExtend},               // Mn   [6] BHAIKSUKI VOWEL SIGN E..BHAIKSUKI SIGN ANUSVARA
	{0x11C3E, 0x11C3E, prExtend},               // Mc       BHAIKSUKI SIGN VISARGA
	{0x11C3F, 0x11C3F, prExtend},               // Mn       BHAIKSUKI SIGN VIRAMA
	{0x11C40, 0x11C40, prALetter},              // Lo       BHAIKSUKI SIGN AVAGRAHA
	{0x11C50, 0x11C59, prNumeric},              // Nd  [10] BHAIKSUKI DIGIT ZERO..BHAIKSUKI DIGIT NINE
	{0x11C72, 0x11C8F, prALetter},              // Lo  [30] MARCHEN LETTER KA..MARCHEN LETTER A
	{0x11C92, 0x11CA7, prExtend},               // Mn  [22] MARCHEN SUBJOINED LETTER KA..MARCHEN SUBJOINED LETTER ZA
	{0x11CA9, 0x11CA9, prExtend},               // Mc       MARCHEN SUBJOINED LETTER YA
	{0x11CAA, 0x11CB0, prExtend},               // Mn   [7] MARCHEN SUBJOINED LETTER RA..MARCHEN VOWEL SIGN AA
	{0x11CB1, 0x11CB1, prExtend},               // Mc       MARCHEN VOWEL SIGN I
	{0x11CB2, 0x11CB3, prExtend},               // Mn   [2] MARCHEN VOWEL SIGN U..MARCHEN VOWEL SIGN E
	{0x11CB4, 0x11CB4, prExtend},               // Mc       MARCHEN VOWEL SIGN O
	{0x11CB5, 0x11CB6, prExtend},               // Mn   [2] MARCHEN SIGN ANUSVARA..MARCHEN SIGN CANDRABINDU
	{0x11D00, 0x11D06, prALetter},              // Lo   [7] MASARAM GONDI LETTER A..MASARAM GONDI LETTER E
	{0x11D08, 0x11D09, prALetter},              // Lo   [2] MASARAM GONDI LETTER AI..MASARAM GONDI LETTER O
	{0x11D0B, 0x11D30, prALetter},              // Lo  [38] MASARAM GONDI LETTER AU..MASARAM GONDI LETTER TRA
	{0x11D31, 0x11D36, prExtend},               // Mn   [6] MASARAM GONDI VOWEL SIGN AA..MASARAM GONDI VOWEL SIGN VOCALIC R
	{0x11D3A, 0x11D3A, prExtend},               // Mn       MASARAM GONDI VOWEL SIGN E
	{0x11D3C, 0x11D3D, prExtend},               // Mn   [2] MASARAM GONDI VOWEL SIGN AI..MASARAM GONDI VOWEL SIGN O
	{0x11D3F, 0x11D45, prExtend},               // Mn   [7] MASARAM GONDI VOWEL SIGN AU..MASARAM GONDI VIRAMA
	{0x11D46, 0x11D46, prALetter},              // Lo       MASARAM GONDI REPHA
	{0x11D47, 0x11D47, prExtend},               // Mn       MASARAM GONDI RA-KARA
	{0x11D50, 0x11D59, prNumeric},              // Nd  [10] MASARAM GONDI DIGIT ZERO..MASARAM GONDI DIGIT NINE
	{0x11D60, 0x11D65, prALetter},              // Lo   [6] GUNJALA GONDI LETTER A..GUNJALA GONDI LETTER UU
	{0x11D67, 0x11D68, prALetter},              // Lo   [2] GUNJALA GONDI LETTER EE..GUNJALA GONDI LETTER AI
	{0x11D6A, 0x11D89, prALetter},              // Lo  [32] GUNJALA GONDI LETTER OO..GUNJALA GONDI LETTER SA
	{0x11D8A, 0x11D8E, prExtend},               // Mc   [5] GUNJALA GONDI VOWEL SIGN AA..GUNJALA GONDI VOWEL SIGN UU
	{0x11D90, 0x11D91, prExtend},               // Mn   [2] GUNJALA GONDI VOWEL SIGN EE..GUNJALA GONDI VOWEL SIGN AI
	{0x11D93, 0x11D94, prExtend},               // Mc   [2] GUNJALA GONDI VOWEL SIGN OO..GUNJALA GONDI VOWEL SIGN AU
	{0x11D95, 0x11D95, prExtend},               // Mn       GUNJALA GONDI SIGN ANUSVARA
	{0x11D96, 0x11D96, prExtend},               // Mc       GUNJALA GONDI SIGN VISARGA
	{0x11D97, 0x11D97, prExtend},               // Mn       GUNJALA GONDI VIRAMA
	{0x11D98, 0x11D98, prALetter},              // Lo       GUNJALA GONDI OM
	{0x11DA0, 0x11DA9, prNumeric},              // Nd  [10] GUNJALA GONDI DIGIT ZERO..GUNJALA GONDI DIGIT NINE
	{0x11EE0, 0x11EF2, prALetter},              // Lo  [19] MAKASAR LETTER KA..MAKASAR ANGKA
	{0x11EF3, 0x11EF4, prExtend},               // Mn   [2] MAKASAR VOWEL SIGN I..MAKASAR VOWEL SIGN U
	{0x11EF5, 0x11EF6, prExtend},               // Mc   [2] MAKASAR VOWEL SIGN E..MAKASAR VOWEL SIGN O
	{0x11F00, 0x11F01, prExtend},               // Mn   [2] KAWI SIGN CANDRABINDU..KAWI SIGN ANUSVARA
	{0x11F02, 0x11F02, prALetter},              // Lo       KAWI SIGN REPHA
	{0x11F03, 0x11F03, prExtend},               // Mc       KAWI SIGN VISARGA
	{0x11F04, 0x11F10, prALetter},              // Lo  [13] KAWI LETTER A..KAWI LETTER O
	{0x11F12, 0x11F33, prALetter},              // Lo  [34] KAWI LETTER KA..KAWI LETTER JNYA
	{0x11F34, 0x11F35, prExtend},               // Mc   [2] KAWI VOWEL SIGN AA..KAWI VOWEL SIGN ALTERNATE AA
	{0x11F36, 0x11F3A, prExtend},               // Mn   [5] KAWI VOWEL SIGN I..KAWI VOWEL SIGN VOCALIC R
	{0x11F3E, 0x11F3F, prExtend},               // Mc   [2] KAWI VOWEL SIGN E..KAWI VOWEL SIGN AI
	{0x11F40, 0x11F40, prExtend},               // Mn       KAWI VOWEL SIGN EU
	{0x11F41, 0x11F41, prExtend},               // Mc       KAWI SIGN KILLER
	{0x11F42, 0x11F42, prExtend},               // Mn       KAWI CONJOINER
	{0x11F50, 0x11F59, prNumeric},              // Nd  [10] KAWI DIGIT ZERO..KAWI DIGIT NINE
	{0x11FB0, 0x11FB0, prALetter},              // Lo       LISU LETTER YHA
	{0x12000, 0x12399, prALetter},              // Lo [922] CUNEIFORM SIGN A..CUNEIFORM SIGN U U
	{0x12400, 0x1246E, prALetter},              // Nl [111] CUNEIFORM NUMERIC SIGN TWO ASH..CUNEIFORM NUMERIC SIGN NINE U VARIANT FORM
	{0x12480, 0x12543, prALetter},              // Lo [196] CUNEIFORM SIGN AB TIMES NUN TENU..CUNEIFORM SIGN ZU5 TIMES THREE DISH TENU
	{0x12F90, 0x12FF0, prALetter},              // Lo  [97] CYPRO-MINOAN SIGN CM001..CYPRO-MINOAN SIGN CM114
	{0x13000, 0x1342F, prALetter},              // Lo [1072] EGYPTIAN HIEROGLYPH A001..EGYPTIAN HIEROGLYPH V011D
	{0x13430, 0x1343F, prFormat},               // Cf  [16] EGYPTIAN HIEROGLYPH VERTICAL JOINER..EGYPTIAN HIEROGLYPH END WALLED ENCLOSURE
	{0x13440, 0x13440, prExtend},               // Mn       EGYPTIAN HIEROGLYPH MIRROR HORIZONTALLY
	{0x13441, 0x13446, prALetter},              // Lo   [6] EGYPTIAN HIEROGLYPH FULL BLANK..EGYPTIAN HIEROGLYPH WIDE LOST SIGN
	{0x13447, 0x13455, prExtend},               // Mn  [15] EGYPTIAN HIEROGLYPH MODIFIER DAMAGED AT TOP START..EGYPTIAN HIEROGLYPH MODIFIER DAMAGED
	{0x14400, 0x14646, prALetter},              // Lo [583] ANATOLIAN HIEROGLYPH A001..ANATOLIAN HIEROGLYPH A530
	{0x16800, 0x16A38, prALetter},              // Lo [569] BAMUM LETTER PHASE-A NGKUE MFON..BAMUM LETTER PHASE-F VUEQ
	{0x16A40, 0x16A5E, prALetter},              // Lo  [31] MRO LETTER TA..MRO LETTER TEK
	{0x16A60, 0x16A69, prNumeric},              // Nd  [10] MRO DIGIT ZERO..MRO DIGIT NINE
	{0x16A70, 0x16ABE, prALetter},              // Lo  [79] TANGSA LETTER OZ..TANGSA LETTER ZA
	{0x16AC0, 0x16AC9, prNumeric},              // Nd  [10] TANGSA DIGIT ZERO..TANGSA DIGIT NINE
	{0x16AD0, 0x16AED, prALetter},              // Lo  [30] BASSA VAH LETTER ENNI..BASSA VAH LETTER I
	{0x16AF0, 0x16AF4, prExtend},               // Mn   [5] BASSA VAH COMBINING HIGH TONE..BASSA VAH COMBINING HIGH-LOW TONE
	{0x16B00, 0x16B2F, prALetter},              // Lo  [48] PAHAWH HMONG VOWEL KEEB..PAHAWH HMONG CONSONANT CAU
	{0x16B30, 0x16B36, prExtend},               // Mn   [7] PAHAWH HMONG MARK CIM TUB..PAHAWH HMONG MARK CIM TAUM
	{0x16B40, 0x16B43, prALetter},              // Lm   [4] PAHAWH HMONG SIGN VOS SEEV..PAHAWH HMONG SIGN IB YAM
	{0x16B50, 0x16B59, prNumeric},              // Nd  [10] PAHAWH HMONG DIGIT ZERO..PAHAWH HMONG DIGIT NINE
	{0x16B63, 0x16B77, prALetter},              // Lo  [21] PAHAWH HMONG SIGN VOS LUB..PAHAWH HMONG SIGN CIM NRES TOS
	{0x16B7D, 0x16B8F, prALetter},              // Lo  [19] PAHAWH HMONG CLAN SIGN TSHEEJ..PAHAWH HMONG CLAN SIGN VWJ
	{0x16E40, 0x16E7F, prALetter},              // L&  [64] MEDEFAIDRIN CAPITAL LETTER M..MEDEFAIDRIN SMALL LETTER Y
	{0x16F00, 0x16F4A, prALetter},              // Lo  [75] MIAO LETTER PA..MIAO LETTER RTE
	{0x16F4F, 0x16F4F, prExtend},               // Mn       MIAO SIGN CONSONANT MODIFIER BAR
	{0x16F50, 0x16F50, prALetter},              // Lo       MIAO LETTER NASALIZATION
	{0x16F51, 0x16F87, prExtend},               // Mc  [55] MIAO SIGN ASPIRATION..MIAO VOWEL SIGN UI
	{0x16F8F, 0x16F92, prExtend},               // Mn   [4] MIAO TONE RIGHT..MIAO TONE BELOW
	{0x16F93, 0x16F9F, prALetter},              // Lm  [13] MIAO LETTER TONE-2..MIAO LETTER REFORMED TONE-8
	{0x16FE0, 0x16FE1, prALetter},              // Lm   [2] TANGUT ITERATION MARK..NUSHU ITERATION MARK
	{0x16FE3, 0x16FE3, prALetter},              // Lm       OLD CHINESE ITERATION MARK
	{0x16FE4, 0x16FE4, prExtend},               // Mn       KHITAN SMALL SCRIPT FILLER
	{0x16FF0, 0x16FF1, prExtend},               // Mc   [2] VIETNAMESE ALTERNATE READING MARK CA..VIETNAMESE ALTERNATE READING MARK NHAY
	{0x1AFF0, 0x1AFF3, prKatakana},             // Lm   [4] KATAKANA LETTER MINNAN TONE-2..KATAKANA LETTER MINNAN TONE-5
	{0x1AFF5, 0x1AFFB, prKatakana},             // Lm   [7] KATAKANA LETTER MINNAN TONE-7..KATAKANA LETTER MINNAN NASALIZED TONE-5
	{0x1AFFD, 0x1AFFE, prKatakana},             // Lm   [2] KATAKANA LETTER MINNAN NASALIZED TONE-7..KATAKANA LETTER MINNAN NASALIZED TONE-8
	{0x1B000, 0x1B000, prKatakana},             // Lo       KATAKANA LETTER ARCHAIC E
	{0x1B120, 0x1B122, prKatakana},             // Lo   [3] KATAKANA LETTER ARCHAIC YI..KATAKANA LETTER ARCHAIC WU
	{0x1B155, 0x1B155, prKatakana},             // Lo       KATAKANA LETTER SMALL KO
	{0x1B164, 0x1B167, prKatakana},             // Lo   [4] KATAKANA LETTER SMALL WI..KATAKANA LETTER SMALL N
	{0x1BC00, 0x1BC6A, prALetter},              // Lo [107] DUPLOYAN LETTER H..DUPLOYAN LETTER VOCALIC M
	{0x1BC70, 0x1BC7C, prALetter},              // Lo  [13] DUPLOYAN AFFIX LEFT HORIZONTAL SECANT..DUPLOYAN AFFIX ATTACHED TANGENT HOOK
	{0x1BC80, 0x1BC88, prALetter},              // Lo   [9] DUPLOYAN AFFIX HIGH ACUTE..DUPLOYAN AFFIX HIGH VERTICAL
	{0x1BC90, 0x1BC99, prALetter},              // Lo  [10] DUPLOYAN AFFIX LOW ACUTE..DUPLOYAN AFFIX LOW ARROW
	{0x1BC9D, 0x1BC9E, prExtend},               // Mn   [2] DUPLOYAN THICK LETTER SELECTOR..DUPLOYAN DOUBLE MARK
	{0x1BCA0, 0x1BCA3, prFormat},               // Cf   [4] SHORTHAND FORMAT LETTER OVERLAP..SHORTHAND FORMAT UP STEP
	{0x1CF00, 0x1CF2D, prExtend},               // Mn  [46] ZNAMENNY COMBINING MARK GORAZDO NIZKO S KRYZHEM ON LEFT..ZNAMENNY COMBINING MARK KRYZH ON LEFT
	{0x1CF30, 0x1CF46, prExtend},               // Mn  [23] ZNAMENNY COMBINING TONAL RANGE MARK MRACHNO..ZNAMENNY PRIZNAK MODIFIER ROG
	{0x1D165, 0x1D166, prExtend},               // Mc   [2] MUSICAL SYMBOL COMBINING STEM..MUSICAL SYMBOL COMBINING SPRECHGESANG STEM
	{0x1D167, 0x1D169, prExtend},               // Mn   [3] MUSICAL SYMBOL COMBINING TREMOLO-1..MUSICAL SYMBOL COMBINING TREMOLO-3
	{0x1D16D, 0x1D172, prExtend},               // Mc   [6] MUSICAL SYMBOL COMBINING AUGMENTATION DOT..MUSICAL SYMBOL COMBINING FLAG-5
	{0x1D173, 0x1D17A, prFormat},               // Cf   [8] MUSICAL SYMBOL BEGIN BEAM..MUSICAL SYMBOL END PHRASE
	{0x1D17B, 0x1D182, prExtend},               // Mn   [8] MUSICAL SYMBOL COMBINING ACCENT..MUSICAL SYMBOL COMBINING LOURE
	{0x1D185, 0x1D18B, prExtend},               // Mn   [7] MUSICAL SYMBOL COMBINING DOIT..MUSICAL SYMBOL COMBINING TRIPLE TONGUE
	{0x1D1AA, 0x1D1AD, prExtend},               // Mn   [4] MUSICAL SYMBOL COMBINING DOWN BOW..MUSICAL SYMBOL COMBINING SNAP PIZZICATO
	{0x1D242, 0x1D244, prExtend},               // Mn   [3] COMBINING GREEK MUSICAL TRISEME..COMBINING GREEK MUSICAL PENTASEME
	{0x1D400, 0x1D454, prALetter},              // L&  [85] MATHEMATICAL BOLD CAPITAL A..MATHEMATICAL ITALIC SMALL G
	{0x1D456, 0x1D49C, prALetter},              // L&  [71] MATHEMATICAL ITALIC SMALL I..MATHEMATICAL SCRIPT CAPITAL A
	{0x1D49E, 0x1D49F, prALetter},              // L&   [2] MATHEMATICAL SCRIPT CAPITAL C..MATHEMATICAL SCRIPT CAPITAL D
	{0x1D4A2, 0x1D4A2, prALetter},              // L&       MATHEMATICAL SCRIPT CAPITAL G
	{0x1D4A5, 0x1D4A6, prALetter},              // L&   [2] MATHEMATICAL SCRIPT CAPITAL J..MATHEMATICAL SCRIPT CAPITAL K
	{0x1D4A9, 0x1D4AC, prALetter},              // L&   [4] MATHEMATICAL SCRIPT CAPITAL N..MATHEMATICAL SCRIPT CAPITAL Q
	{0x1D4AE, 0x1D4B9, prALetter},              // L&  [12] MATHEMATICAL SCRIPT CAPITAL S..MATHEMATICAL SCRIPT SMALL D
	{0x1D4BB, 0x1D4BB, prALetter},              // L&       MATHEMATICAL SCRIPT SMALL F
	{0x1D4BD, 0x1D4C3, prALetter},              // L&   [7] MATHEMATICAL SCRIPT SMALL H..MATHEMATICAL SCRIPT SMALL N
	{0x1D4C5, 0x1D505, prALetter},              // L&  [65] MATHEMATICAL SCRIPT SMALL P..MATHEMATICAL FRAKTUR CAPITAL B
	{0x1D507, 0x1D50A, prALetter},              // L&   [4] MATHEMATICAL FRAKTUR CAPITAL D..MATHEMATICAL FRAKTUR CAPITAL G
	{0x1D50D, 0x1D514, prALetter},              // L&   [8] MATHEMATICAL FRAKTUR CAPITAL J..MATHEMATICAL FRAKTUR CAPITAL Q
	{0x1D516, 0x1D51C, prALetter},              // L&   [7] MATHEMATICAL FRAKTUR CAPITAL S..MATHEMATICAL FRAKTUR CAPITAL Y
	{0x1D51E, 0x1D539, prALetter},              // L&  [28] MATHEMATICAL FRAKTUR SMALL A..MATHEMATICAL DOUBLE-STRUCK CAPITAL B
	{0x1D53B, 0x1D53E, prALetter},              // L&   [4] MATHEMATICAL DOUBLE-STRUCK CAPITAL D..MATHEMATICAL DOUBLE-STRUCK CAPITAL G
	{0x1D540, 0x1D544, prALetter},              // L&   [5] MATHEMATICAL DOUBLE-STRUCK CAPITAL I..MATHEMATICAL DOUBLE-STRUCK CAPITAL M
	{0x1D546, 0x1D546, prALetter},              // L&       MATHEMATICAL DOUBLE-STRUCK CAPITAL O
	{0x1D54A, 0x1D550, prALetter},              // L&   [7] MATHEMATICAL DOUBLE-STRUCK CAPITAL S..MATHEMATICAL DOUBLE-STRUCK CAPITAL Y
	{0x1D552, 0x1D6A5, prALetter},              // L& [340] MATHEMATICAL DOUBLE-STRUCK SMALL A..MATHEMATICAL ITALIC SMALL DOTLESS J
	{0x1D6A8, 0x1D6C0, prALetter},              // L&  [25] MATHEMATICAL BOLD CAPITAL ALPHA..MATHEMATICAL BOLD CAPITAL OMEGA
	{0x1D6C2, 0x1D6DA, prALetter},              // L&  [25] MATHEMATICAL BOLD SMALL ALPHA..MATHEMATICAL BOLD SMALL OMEGA
	{0x1D6DC, 0x1D6FA, prALetter},              // L&  [31] MATHEMATICAL BOLD EPSILON SYMBOL..MATHEMATICAL ITALIC CAPITAL OMEGA
	{0x1D6FC, 0x1D714, prALetter},              // L&  [25] MATHEMATICAL ITALIC SMALL ALPHA..MATHEMATICAL ITALIC SMALL OMEGA
	{0x1D716, 0x1D734, prALetter},              // L&  [31] MATHEMATICAL ITALIC EPSILON SYMBOL..MATHEMATICAL BOLD ITALIC CAPITAL OMEGA
	{0x1D736, 0x1D74E, prALetter},              // L&  [25] MATHEMATICAL BOLD ITALIC SMALL ALPHA..MATHEMATICAL BOLD ITALIC SMALL OMEGA
	{0x1D750, 0x1D76E, prALetter},              // L&  [31] MATHEMATICAL BOLD ITALIC EPSILON SYMBOL..MATHEMATICAL SANS-SERIF BOLD CAPITAL OMEGA
	{0x1D770, 0x1D788, prALetter},              // L&  [25] MATHEMATICAL SANS-SERIF BOLD SMALL ALPHA..MATHEMATICAL SANS-SERIF BOLD SMALL OMEGA
	{0x1D78A, 0x1D7A8, prALetter},              // L&  [31] MATHEMATICAL SANS-SERIF BOLD EPSILON SYMBOL..MATHEMATICAL SANS-SERIF BOLD ITALIC CAPITAL OMEGA
	{0x1D7AA, 0x1D7C2, prALetter},              // L&  [25] MATHEMATICAL SANS-SERIF BOLD ITALIC SMALL ALPHA..MATHEMATICAL SANS-SERIF BOLD ITALIC SMALL OMEGA
	{0x1D7C4, 0x1D7CB, prALetter},              // L&   [8] MATHEMATICAL SANS-SERIF BOLD ITALIC EPSILON SYMBOL..MATHEMATICAL BOLD SMALL DIGAMMA
	{0x1D7CE, 0x1D7FF, prNumeric},              // Nd  [50] MATHEMATICAL BOLD DIGIT ZERO..MATHEMATICAL MONOSPACE DIGIT NINE
	{0x1DA00, 0x1DA36, prExtend},               // Mn  [55] SIGNWRITING HEAD RIM..SIGNWRITING AIR SUCKING IN
	{0x1DA3B, 0x1DA6C, prExtend},               // Mn  [50] SIGNWRITING MOUTH CLOSED NEUTRAL..SIGNWRITING EXCITEMENT
	{0x1DA75, 0x1DA75, prExtend},               // Mn       SIGNWRITING UPPER BODY TILTING FROM HIP JOINTS
	{0x1DA84, 0x1DA84, prExtend},               // Mn       SIGNWRITING LOCATION HEAD NECK
	{0x1DA9B, 0x1DA9F, prExtend},               // Mn   [5] SIGNWRITING FILL MODIFIER-2..SIGNWRITING FILL MODIFIER-6
	{0x1DAA1, 0x1DAAF, prExtend},               // Mn  [15] SIGNWRITING ROTATION MODIFIER-2..SIGNWRITING ROTATION MODIFIER-16
	{0x1DF00, 0x1DF09, prALetter},              // L&  [10] LATIN SMALL LETTER FENG DIGRAPH WITH TRILL..LATIN SMALL LETTER T WITH HOOK AND RETROFLEX HOOK
	{0x1DF0A, 0x1DF0A, prALetter},              // Lo       LATIN LETTER RETROFLEX CLICK WITH RETROFLEX HOOK
	{0x1DF0B, 0x1DF1E, prALetter},              // L&  [20] LATIN SMALL LETTER ESH WITH DOUBLE BAR..LATIN SMALL LETTER S WITH CURL
	{0x1DF25, 0x1DF2A, prALetter},              // L&   [6] LATIN SMALL LETTER D WITH MID-HEIGHT LEFT HOOK..LATIN SMALL LETTER T WITH MID-HEIGHT LEFT HOOK
	{0x1E000, 0x1E006, prExtend},               // Mn   [7] COMBINING GLAGOLITIC LETTER AZU..COMBINING GLAGOLITIC LETTER ZHIVETE
	{0x1E008, 0x1E018, prExtend},               // Mn  [17] COMBINING GLAGOLITIC LETTER ZEMLJA..COMBINING GLAGOLITIC LETTER HERU
	{0x1E01B, 0x1E021, prExtend},               // Mn   [7] COMBINING GLAGOLITIC LETTER SHTA..COMBINING GLAGOLITIC LETTER YATI
	{0x1E023, 0x1E024, prExtend},               // Mn   [2] COMBINING GLAGOLITIC LETTER YU..COMBINING GLAGOLITIC LETTER SMALL YUS
	{0x1E026, 0x1E02A, prExtend},               // Mn   [5] COMBINING GLAGOLITIC LETTER YO..COMBINING GLAGOLITIC LETTER FITA
	{0x1E030, 0x1E06D, prALetter},              // Lm  [62] MODIFIER LETTER CYRILLIC SMALL A..MODIFIER LETTER CYRILLIC SMALL STRAIGHT U WITH STROKE
	{0x1E08F, 0x1E08F, prExtend},               // Mn       COMBINING CYRILLIC SMALL LETTER BYELORUSSIAN-UKRAINIAN I
	{0x1E100, 0x1E12C, prALetter},              // Lo  [45] NYIAKENG PUACHUE HMONG LETTER MA..NYIAKENG PUACHUE HMONG LETTER W
	{0x1E130, 0x1E136, prExtend},               // Mn   [7] NYIAKENG PUACHUE HMONG TONE-B..NYIAKENG PUACHUE HMONG TONE-D
	{0x1E137, 0x1E13D, prALetter},              // Lm   [7] NYIAKENG PUACHUE HMONG SIGN FOR PERSON..NYIAKENG PUACHUE HMONG SYLLABLE LENGTHENER
	{0x1E140, 0x1E149, prNumeric},              // Nd  [10] NYIAKENG PUACHUE HMONG DIGIT ZERO..NYIAKENG PUACHUE HMONG DIGIT NINE
	{0x1E14E, 0x1E14E, prALetter},              // Lo       NYIAKENG PUACHUE HMONG LOGOGRAM NYAJ
	{0x1E290, 0x1E2AD, prALetter},              // Lo  [30] TOTO LETTER PA..TOTO LETTER A
	{0x1E2AE, 0x1E2AE, prExtend},               // Mn       TOTO SIGN RISING TONE
	{0x1E2C0, 0x1E2EB, prALetter},              // Lo  [44] WANCHO LETTER AA..WANCHO LETTER YIH
	{0x1E2EC, 0x1E2EF, prExtend},               // Mn   [4] WANCHO TONE TUP..WANCHO TONE KOINI
	{0x1E2F0, 0x1E2F9, prNumeric},              // Nd  [10] WANCHO DIGIT ZERO..WANCHO DIGIT NINE
	{0x1E4D0, 0x1E4EA, prALetter},              // Lo  [27] NAG MUNDARI LETTER O..NAG MUNDARI LETTER ELL
	{0x1E4EB, 0x1E4EB, prALetter},              // Lm       NAG MUNDARI SIGN OJOD
	{0x1E4EC, 0x1E4EF, prExtend},               // Mn   [4] NAG MUNDARI SIGN MUHOR..NAG MUNDARI SIGN SUTUH
	{0x1E4F0, 0x1E4F9, prNumeric},              // Nd  [10] NAG MUNDARI DIGIT ZERO..NAG MUNDARI DIGIT NINE
	{0x1E7E0, 0x1E7E6, prALetter},              // Lo   [7] ETHIOPIC SYLLABLE HHYA..ETHIOPIC SYLLABLE HHYO
	{0x1E7E8, 0x1E7EB, prALetter},              // Lo   [4] ETHIOPIC SYLLABLE GURAGE HHWA..ETHIOPIC SYLLABLE HHWE
	{0x1E7ED, 0x1E7EE, prALetter},              // Lo   [2] ETHIOPIC SYLLABLE GURAGE MWI..ETHIOPIC SYLLABLE GURAGE MWEE
	{0x1E7F0, 0x1E7FE, prALetter},              // Lo  [15] ETHIOPIC SYLLABLE GURAGE QWI..ETHIOPIC SYLLABLE GURAGE PWEE
	{0x1E800, 0x1E8C4, prALetter},              // Lo [197] MENDE KIKAKUI SYLLABLE M001 KI..MENDE KIKAKUI SYLLABLE M060 NYON
	{0x1E8D0, 0x1E8D6, prExtend},               // Mn   [7] MENDE KIKAKUI COMBINING NUMBER TEENS..MENDE KIKAKUI COMBINING NUMBER MILLIONS
	{0x1E900, 0x1E943, prALetter},              // L&  [68] ADLAM CAPITAL LETTER ALIF..ADLAM SMALL LETTER SHA
	{0x1E944, 0x1E94A, prExtend},               // Mn   [7] ADLAM ALIF LENGTHENER..ADLAM NUKTA
	{0x1E94B, 0x1E94B, prALetter},              // Lm       ADLAM NASALIZATION MARK
	{0x1E950, 0x1E959, prNumeric},              // Nd  [10] ADLAM DIGIT ZERO..ADLAM DIGIT NINE
	{0x1EE00, 0x1EE03, prALetter},              // Lo   [4] ARABIC MATHEMATICAL ALEF..ARABIC MATHEMATICAL DAL
	{0x1EE05, 0x1EE1F, prALetter},              // Lo  [27] ARABIC MATHEMATICAL WAW..ARABIC MATHEMATICAL DOTLESS QAF
	{0x1EE21, 0x1EE22, prALetter},              // Lo   [2] ARABIC MATHEMATICAL INITIAL BEH..ARABIC MATHEMATICAL INITIAL JEEM
	{0x1EE24, 0x1EE24, prALetter},              // Lo       ARABIC MATHEMATICAL INITIAL HEH
	{0x1EE27, 0x1EE27, prALetter},              // Lo       ARABIC MATHEMATICAL INITIAL HAH
	{0x1EE29, 0x1EE32, prALetter},              // Lo  [10] ARABIC MATHEMATICAL INITIAL YEH..ARABIC MATHEMATICAL INITIAL QAF
	{0x1EE34, 0x1EE37, prALetter},              // Lo   [4] ARABIC MATHEMATICAL INITIAL SHEEN..ARABIC MATHEMATICAL INITIAL KHAH
	{0x1EE39, 0x1EE39, prALetter},              // Lo       ARABIC MATHEMATICAL INITIAL DAD
	{0x1EE3B, 0x1EE3B, prALetter},              // Lo       ARABIC MATHEMATICAL INITIAL GHAIN
	{0x1EE42, 0x1EE42, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED JEEM
	{0x1EE47, 0x1EE47, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED HAH
	{0x1EE49, 0x1EE49, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED YEH
	{0x1EE4B, 0x1EE4B, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED LAM
	{0x1EE4D, 0x1EE4F, prALetter},              // Lo   [3] ARABIC MATHEMATICAL TAILED NOON..ARABIC MATHEMATICAL TAILED AIN
	{0x1EE51, 0x1EE52, prALetter},              // Lo   [2] ARABIC MATHEMATICAL TAILED SAD..ARABIC MATHEMATICAL TAILED QAF
	{0x1EE54, 0x1EE54, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED SHEEN
	{0x1EE57, 0x1EE57, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED KHAH
	{0x1EE59, 0x1EE59, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED DAD
	{0x1EE5B, 0x1EE5B, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED GHAIN
	{0x1EE5D, 0x1EE5D, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED DOTLESS NOON
	{0x1EE5F, 0x1EE5F, prALetter},              // Lo       ARABIC MATHEMATICAL TAILED DOTLESS QAF
	{0x1EE61, 0x1EE62, prALetter},              // Lo   [2] ARABIC MATHEMATICAL STRETCHED BEH..ARABIC MATHEMATICAL STRETCHED JEEM
	{0x1EE64, 0x1EE64, prALetter},              // Lo       ARABIC MATHEMATICAL STRETCHED HEH
	{0x1EE67, 0x1EE6A, prALetter},              // Lo   [4] ARABIC MATHEMATICAL STRETCHED HAH..ARABIC MATHEMATICAL STRETCHED KAF
	{0x1EE6C, 0x1EE72, prALetter},              // Lo   [7] ARABIC MATHEMATICAL STRETCHED MEEM..ARABIC MATHEMATICAL STRETCHED QAF
	{0x1EE74, 0x1EE77, prALetter},              // Lo   [4] ARABIC MATHEMATICAL STRETCHED SHEEN..ARABIC MATHEMATICAL STRETCHED KHAH
	{0x1EE79, 0x1EE7C, prALetter},              // Lo   [4] ARABIC MATHEMATICAL STRETCHED DAD..ARABIC MATHEMATICAL STRETCHED DOTLESS BEH
	{0x1EE7E, 0x1EE7E, prALetter},              // Lo       ARABIC MATHEMATICAL STRETCHED DOTLESS FEH
	{0x1EE80, 0x1EE89, prALetter},              // Lo  [10] ARABIC MATHEMATICAL LOOPED ALEF..ARABIC MATHEMATICAL LOOPED YEH
	{0x1EE8B, 0x1EE9B, prALetter},              // Lo  [17] ARABIC MATHEMATICAL LOOPED LAM..ARABIC MATHEMATICAL LOOPED GHAIN
	{0x1EEA1, 0x1EEA3, prALetter},              // Lo   [3] ARABIC MATHEMATICAL DOUBLE-STRUCK BEH..ARABIC MATHEMATICAL DOUBLE-STRUCK DAL
	{0x1EEA5, 0x1EEA9, prALetter},              // Lo   [5] ARABIC MATHEMATICAL DOUBLE-STRUCK WAW..ARABIC MATHEMATICAL DOUBLE-STRUCK YEH
	{0x1EEAB, 0x1EEBB, prALetter},              // Lo  [17] ARABIC MATHEMATICAL DOUBLE-STRUCK LAM..ARABIC MATHEMATICAL DOUBLE-STRUCK GHAIN
	{0x1F000, 0x1F003, prExtendedPictographic}, // E0.0   [4] (🀀..🀃)    MAHJONG TILE EAST WIND..MAHJONG TILE NORTH WIND
	{0x1F004, 0x1F004, prExtendedPictographic}, // E0.6   [1] (🀄)       mahjong red dragon
	{0x1F005, 0x1F0CE, prExtendedPictographic}, // E0.0 [202] (🀅..🃎)    MAHJONG TILE GREEN DRAGON..PLAYING CARD KING OF DIAMONDS
	{0x1F0CF, 0x1F0CF, prExtendedPictographic}, // E0.6   [1] (🃏)       joker
	{0x1F0D0, 0x1F0FF, prExtendedPictographic}, // E0.0  [48] (🃐..🃿)    <reserved-1F0D0>..<reserved-1F0FF>
	{0x1F10D, 0x1F10F, prExtendedPictographic}, // E0.0   [3] (🄍..🄏)    CIRCLED ZERO WITH SLASH..CIRCLED DOLLAR SIGN WITH OVERLAID BACKSLASH
	{0x1F12F, 0x1F12F, prExtendedPictographic}, // E0.0   [1] (🄯)       COPYLEFT SYMBOL
	{0x1F130, 0x1F149, prALetter},              // So  [26] SQUARED LATIN CAPITAL LETTER A..SQUARED LATIN CAPITAL LETTER Z
	{0x1F150, 0x1F169, prALetter},              // So  [26] NEGATIVE CIRCLED LATIN CAPITAL LETTER A..NEGATIVE CIRCLED LATIN CAPITAL LETTER Z
	{0x1F16C, 0x1F16F, prExtendedPictographic}, // E0.0   [4] (🅬..🅯)    RAISED MR SIGN..CIRCLED HUMAN FIGURE
	{0x1F170, 0x1F189, prALetter},              // So  [26] NEGATIVE SQUARED LATIN CAPITAL LETTER A..NEGATIVE SQUARED LATIN CAPITAL LETTER Z
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
	{0x1FBF0, 0x1FBF9, prNumeric},              // Nd  [10] SEGMENTED DIGIT ZERO..SEGMENTED DIGIT NINE
	{0x1FC00, 0x1FFFD, prExtendedPictographic}, // E0.0[1022] (🰀..🿽)    <reserved-1FC00>..<reserved-1FFFD>
	{0xE0001, 0xE0001, prFormat},               // Cf       LANGUAGE TAG
	{0xE0020, 0xE007F, prExtend},               // Cf  [96] TAG SPACE..CANCEL TAG
	{0xE0100, 0xE01EF, prExtend},               // Mn [240] VARIATION SELECTOR-17..VARIATION SELECTOR-256
}
