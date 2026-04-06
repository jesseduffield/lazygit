/*
Package uniseg implements Unicode Text Segmentation, Unicode Line Breaking, and
string width calculation for monospace fonts. Unicode Text Segmentation conforms
to Unicode Standard Annex #29 (https://unicode.org/reports/tr29/) and Unicode
Line Breaking conforms to Unicode Standard Annex #14
(https://unicode.org/reports/tr14/).

In short, using this package, you can split a string into grapheme clusters
(what people would usually refer to as a "character"), into words, and into
sentences. Or, in its simplest case, this package allows you to count the number
of characters in a string, especially when it contains complex characters such
as emojis, combining characters, or characters from Asian, Arabic, Hebrew, or
other languages. Additionally, you can use it to implement line breaking (or
"word wrapping"), that is, to determine where text can be broken over to the
next line when the width of the line is not big enough to fit the entire text.
Finally, you can use it to calculate the display width of a string for monospace
fonts.

# Getting Started

If you just want to count the number of characters in a string, you can use
[GraphemeClusterCount]. If you want to determine the display width of a string,
you can use [StringWidth]. If you want to iterate over a string, you can use
[Step], [StepString], or the [Graphemes] class (more convenient but less
performant). This will provide you with all information: grapheme clusters,
word boundaries, sentence boundaries, line breaks, and monospace character
widths. The specialized functions [FirstGraphemeCluster],
[FirstGraphemeClusterInString], [FirstWord], [FirstWordInString],
[FirstSentence], and [FirstSentenceInString] can be used if only one type of
information is needed.

# Grapheme Clusters

Consider the rainbow flag emoji: 🏳️‍🌈. On most modern systems, it appears as one
character. But its string representation actually has 14 bytes, so counting
bytes (or using len("🏳️‍🌈")) will not work as expected. Counting runes won't,
either: The flag has 4 Unicode code points, thus 4 runes. The stdlib function
utf8.RuneCountInString("🏳️‍🌈") and len([]rune("🏳️‍🌈")) will both return 4.

The [GraphemeClusterCount] function will return 1 for the rainbow flag emoji.
The Graphemes class and a variety of functions in this package will allow you to
split strings into its grapheme clusters.

# Word Boundaries

Word boundaries are used in a number of different contexts. The most familiar
ones are selection (double-click mouse selection), cursor movement ("move to
next word" control-arrow keys), and the dialog option "Whole Word Search" for
search and replace. This package provides methods for determining word
boundaries.

# Sentence Boundaries

Sentence boundaries are often used for triple-click or some other method of
selecting or iterating through blocks of text that are larger than single words.
They are also used to determine whether words occur within the same sentence in
database queries. This package provides methods for determining sentence
boundaries.

# Line Breaking

Line breaking, also known as word wrapping, is the process of breaking a section
of text into lines such that it will fit in the available width of a page,
window or other display area. This package provides methods to determine the
positions in a string where a line must be broken, may be broken, or must not be
broken.

# Monospace Width

Monospace width, as referred to in this package, is the width of a string in a
monospace font. This is commonly used in terminal user interfaces or text
displays or editors that don't support proportional fonts. A width of 1
corresponds to a single character cell. The C function [wcswidth()] and its
implementation in other programming languages is in widespread use for the same
purpose. However, there is no standard for the calculation of such widths, and
this package differs from wcswidth() in a number of ways, presumably to generate
more visually pleasing results.

To start, we assume that every code point has a width of 1, with the following
exceptions:

  - Code points with grapheme cluster break properties Control, CR, LF, Extend,
    and ZWJ have a width of 0.
  - U+2E3A, Two-Em Dash, has a width of 3.
  - U+2E3B, Three-Em Dash, has a width of 4.
  - Characters with the East-Asian Width properties "Fullwidth" (F) and "Wide"
    (W) have a width of 2. (Properties "Ambiguous" (A) and "Neutral" (N) both
    have a width of 1.)
  - Code points with grapheme cluster break property Regional Indicator have a
    width of 2.
  - Code points with grapheme cluster break property Extended Pictographic have
    a width of 2, unless their Emoji Presentation flag is "No", in which case
    the width is 1.

For Hangul grapheme clusters composed of conjoining Jamo and for Regional
Indicators (flags), all code points except the first one have a width of 0. For
grapheme clusters starting with an Extended Pictographic, any additional code
point will force a total width of 2, except if the Variation Selector-15
(U+FE0E) is included, in which case the total width is always 1. Grapheme
clusters ending with Variation Selector-16 (U+FE0F) have a width of 2.

Note that whether these widths appear correct depends on your application's
render engine, to which extent it conforms to the Unicode Standard, and its
choice of font.

[wcswidth()]: https://man7.org/linux/man-pages/man3/wcswidth.3.html
*/
package uniseg
