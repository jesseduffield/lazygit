/*
Package uniseg implements Unicode Text Segmentation and Unicode Line Breaking.
Unicode Text Segmentation conforms to Unicode Standard Annex #29
(https://unicode.org/reports/tr29/) and Unicode Line Breaking conforms to
Unicode Standard Annex #14 (https://unicode.org/reports/tr14/).

In short, using this package, you can split a string into grapheme clusters
(what people would usually refer to as a "character"), into words, and into
sentences. Or, in its simplest case, this package allows you to count the number
of characters in a string, especially when it contains complex characters such
as emojis, combining characters, or characters from Asian, Arabic, Hebrew, or
other languages. Additionally, you can use it to implement line breaking (or
"word wrapping"), that is, to determine where text can be broken over to the
next line when the width of the line is not big enough to fit the entire text.

Grapheme Clusters

Consider the rainbow flag emoji: 🏳️‍🌈. On most modern systems, it appears as one
character. But its string representation actually has 14 bytes, so counting
bytes (or using len("🏳️‍🌈")) will not work as expected. Counting runes won't,
either: The flag has 4 Unicode code points, thus 4 runes. The stdlib function
utf8.RuneCountInString("🏳️‍🌈") and len([]rune("🏳️‍🌈")) will both return 4.

The uniseg.GraphemeClusterCount(str) function will return 1 for the rainbow flag
emoji. The Graphemes class and a variety of functions in this package will allow
you to split strings into its grapheme clusters.

Word Boundaries

Word boundaries are used in a number of different contexts. The most familiar
ones are selection (double-click mouse selection), cursor movement ("move to
next word" control-arrow keys), and the dialog option "Whole Word Search" for
search and replace. This package provides methods for determining word
boundaries.

Sentence Boundaries

Sentence boundaries are often used for triple-click or some other method of
selecting or iterating through blocks of text that are larger than single words.
They are also used to determine whether words occur within the same sentence in
database queries. This package provides methods for determining sentence
boundaries.

Line Breaking

Line breaking, also known as word wrapping, is the process of breaking a section
of text into lines such that it will fit in the available width of a page,
window or other display area. This package provides methods to determine the
positions in a string where a line must be broken, may be broken, or must not be
broken.

*/
package uniseg
