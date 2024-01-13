package s

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Swift lexer.
var Swift = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Swift",
		Aliases:   []string{"swift"},
		Filenames: []string{"*.swift"},
		MimeTypes: []string{"text/x-swift"},
	},
	swiftRules,
))

func swiftRules() Rules {
	return Rules{
		"root": {
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`//`, CommentSingle, Push("comment-single")},
			{`/\*`, CommentMultiline, Push("comment-multi")},
			{`#(if|elseif|else|endif|available)\b`, CommentPreproc, Push("preproc")},
			Include("keywords"),
			{Words(``, `\b`, `Array`, `AutoreleasingUnsafeMutablePointer`, `BidirectionalReverseView`, `Bit`, `Bool`, `CFunctionPointer`, `COpaquePointer`, `CVaListPointer`, `Character`, `ClosedInterval`, `CollectionOfOne`, `ContiguousArray`, `Dictionary`, `DictionaryGenerator`, `DictionaryIndex`, `Double`, `EmptyCollection`, `EmptyGenerator`, `EnumerateGenerator`, `EnumerateSequence`, `FilterCollectionView`, `FilterCollectionViewIndex`, `FilterGenerator`, `FilterSequenceView`, `Float`, `Float80`, `FloatingPointClassification`, `GeneratorOf`, `GeneratorOfOne`, `GeneratorSequence`, `HalfOpenInterval`, `HeapBuffer`, `HeapBufferStorage`, `ImplicitlyUnwrappedOptional`, `IndexingGenerator`, `Int`, `Int16`, `Int32`, `Int64`, `Int8`, `LazyBidirectionalCollection`, `LazyForwardCollection`, `LazyRandomAccessCollection`, `LazySequence`, `MapCollectionView`, `MapSequenceGenerator`, `MapSequenceView`, `MirrorDisposition`, `ObjectIdentifier`, `OnHeap`, `Optional`, `PermutationGenerator`, `QuickLookObject`, `RandomAccessReverseView`, `Range`, `RangeGenerator`, `RawByte`, `Repeat`, `ReverseBidirectionalIndex`, `ReverseRandomAccessIndex`, `SequenceOf`, `SinkOf`, `Slice`, `StaticString`, `StrideThrough`, `StrideThroughGenerator`, `StrideTo`, `StrideToGenerator`, `String`, `UInt`, `UInt16`, `UInt32`, `UInt64`, `UInt8`, `UTF16`, `UTF32`, `UTF8`, `UnicodeDecodingResult`, `UnicodeScalar`, `Unmanaged`, `UnsafeBufferPointer`, `UnsafeBufferPointerGenerator`, `UnsafeMutableBufferPointer`, `UnsafeMutablePointer`, `UnsafePointer`, `Zip2`, `ZipGenerator2`, `AbsoluteValuable`, `AnyObject`, `ArrayLiteralConvertible`, `BidirectionalIndexType`, `BitwiseOperationsType`, `BooleanLiteralConvertible`, `BooleanType`, `CVarArgType`, `CollectionType`, `Comparable`, `DebugPrintable`, `DictionaryLiteralConvertible`, `Equatable`, `ExtendedGraphemeClusterLiteralConvertible`, `ExtensibleCollectionType`, `FloatLiteralConvertible`, `FloatingPointType`, `ForwardIndexType`, `GeneratorType`, `Hashable`, `IntegerArithmeticType`, `IntegerLiteralConvertible`, `IntegerType`, `IntervalType`, `MirrorType`, `MutableCollectionType`, `MutableSliceable`, `NilLiteralConvertible`, `OutputStreamType`, `Printable`, `RandomAccessIndexType`, `RangeReplaceableCollectionType`, `RawOptionSetType`, `RawRepresentable`, `Reflectable`, `SequenceType`, `SignedIntegerType`, `SignedNumberType`, `SinkType`, `Sliceable`, `Streamable`, `Strideable`, `StringInterpolationConvertible`, `StringLiteralConvertible`, `UnicodeCodecType`, `UnicodeScalarLiteralConvertible`, `UnsignedIntegerType`, `_ArrayBufferType`, `_BidirectionalIndexType`, `_CocoaStringType`, `_CollectionType`, `_Comparable`, `_ExtensibleCollectionType`, `_ForwardIndexType`, `_Incrementable`, `_IntegerArithmeticType`, `_IntegerType`, `_ObjectiveCBridgeable`, `_RandomAccessIndexType`, `_RawOptionSetType`, `_SequenceType`, `_Sequence_Type`, `_SignedIntegerType`, `_SignedNumberType`, `_Sliceable`, `_Strideable`, `_SwiftNSArrayRequiredOverridesType`, `_SwiftNSArrayType`, `_SwiftNSCopyingType`, `_SwiftNSDictionaryRequiredOverridesType`, `_SwiftNSDictionaryType`, `_SwiftNSEnumeratorType`, `_SwiftNSFastEnumerationType`, `_SwiftNSStringRequiredOverridesType`, `_SwiftNSStringType`, `_UnsignedIntegerType`, `C_ARGC`, `C_ARGV`, `Process`, `Any`, `AnyClass`, `BooleanLiteralType`, `CBool`, `CChar`, `CChar16`, `CChar32`, `CDouble`, `CFloat`, `CInt`, `CLong`, `CLongLong`, `CShort`, `CSignedChar`, `CUnsignedInt`, `CUnsignedLong`, `CUnsignedShort`, `CWideChar`, `ExtendedGraphemeClusterType`, `Float32`, `Float64`, `FloatLiteralType`, `IntMax`, `IntegerLiteralType`, `StringLiteralType`, `UIntMax`, `UWord`, `UnicodeScalarType`, `Void`, `Word`, `NSErrorPointer`, `NSObjectProtocol`, `Selector`), NameBuiltin, nil},
			{Words(``, `\b`, `abs`, `advance`, `alignof`, `alignofValue`, `assert`, `assertionFailure`, `contains`, `count`, `countElements`, `debugPrint`, `debugPrintln`, `distance`, `dropFirst`, `dropLast`, `dump`, `enumerate`, `equal`, `extend`, `fatalError`, `filter`, `find`, `first`, `getVaList`, `indices`, `insert`, `isEmpty`, `join`, `last`, `lazy`, `lexicographicalCompare`, `map`, `max`, `maxElement`, `min`, `minElement`, `numericCast`, `overlaps`, `partition`, `precondition`, `preconditionFailure`, `prefix`, `print`, `println`, `reduce`, `reflect`, `removeAll`, `removeAtIndex`, `removeLast`, `removeRange`, `reverse`, `sizeof`, `sizeofValue`, `sort`, `sorted`, `splice`, `split`, `startsWith`, `stride`, `strideof`, `strideofValue`, `suffix`, `swap`, `toDebugString`, `toString`, `transcode`, `underestimateCount`, `unsafeAddressOf`, `unsafeBitCast`, `unsafeDowncast`, `withExtendedLifetime`, `withUnsafeMutablePointer`, `withUnsafeMutablePointers`, `withUnsafePointer`, `withUnsafePointers`, `withVaList`), NameBuiltinPseudo, nil},
			{`\$\d+`, NameVariable, nil},
			{`0b[01_]+`, LiteralNumberBin, nil},
			{`0o[0-7_]+`, LiteralNumberOct, nil},
			{`0x[0-9a-fA-F_]+`, LiteralNumberHex, nil},
			{`[0-9][0-9_]*(\.[0-9_]+[eE][+\-]?[0-9_]+|\.[0-9_]*|[eE][+\-]?[0-9_]+)`, LiteralNumberFloat, nil},
			{`[0-9][0-9_]*`, LiteralNumberInteger, nil},
			{`"`, LiteralString, Push("string")},
			{"[(){}\\[\\].,:;=@#`?]|->|[<&?](?=\\w)|(?<=\\w)[>!?]", Punctuation, nil},
			{`[/=\-+!*%<>&|^?~]+`, Operator, nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"keywords": {
			{Words(``, `\b`, `as`, `break`, `case`, `catch`, `continue`, `default`, `defer`, `do`, `else`, `fallthrough`, `for`, `guard`, `if`, `in`, `is`, `repeat`, `return`, `#selector`, `switch`, `throw`, `try`, `where`, `while`), Keyword, nil},
			{`@availability\([^)]+\)`, KeywordReserved, nil},
			{Words(``, `\b`, `associativity`, `convenience`, `dynamic`, `didSet`, `final`, `get`, `indirect`, `infix`, `inout`, `lazy`, `left`, `mutating`, `none`, `nonmutating`, `optional`, `override`, `postfix`, `precedence`, `prefix`, `Protocol`, `required`, `rethrows`, `right`, `set`, `throws`, `Type`, `unowned`, `weak`, `willSet`, `@availability`, `@autoclosure`, `@noreturn`, `@NSApplicationMain`, `@NSCopying`, `@NSManaged`, `@objc`, `@UIApplicationMain`, `@IBAction`, `@IBDesignable`, `@IBInspectable`, `@IBOutlet`), KeywordReserved, nil},
			{`(as|dynamicType|false|is|nil|self|Self|super|true|__COLUMN__|__FILE__|__FUNCTION__|__LINE__|_|#(?:file|line|column|function))\b`, KeywordConstant, nil},
			{`import\b`, KeywordDeclaration, Push("module")},
			{`(class|enum|extension|struct|protocol)(\s+)([a-zA-Z_]\w*)`, ByGroups(KeywordDeclaration, Text, NameClass), nil},
			{`(func)(\s+)([a-zA-Z_]\w*)`, ByGroups(KeywordDeclaration, Text, NameFunction), nil},
			{`(var|let)(\s+)([a-zA-Z_]\w*)`, ByGroups(KeywordDeclaration, Text, NameVariable), nil},
			{Words(``, `\b`, `class`, `deinit`, `enum`, `extension`, `func`, `import`, `init`, `internal`, `let`, `operator`, `private`, `protocol`, `public`, `static`, `struct`, `subscript`, `typealias`, `var`), KeywordDeclaration, nil},
		},
		"comment": {
			{`:param: [a-zA-Z_]\w*|:returns?:|(FIXME|MARK|TODO):`, CommentSpecial, nil},
		},
		"comment-single": {
			{`\n`, Text, Pop(1)},
			Include("comment"),
			{`[^\n]`, CommentSingle, nil},
		},
		"comment-multi": {
			Include("comment"),
			{`[^*/]`, CommentMultiline, nil},
			{`/\*`, CommentMultiline, Push()},
			{`\*/`, CommentMultiline, Pop(1)},
			{`[*/]`, CommentMultiline, nil},
		},
		"module": {
			{`\n`, Text, Pop(1)},
			{`[a-zA-Z_]\w*`, NameClass, nil},
			Include("root"),
		},
		"preproc": {
			{`\n`, Text, Pop(1)},
			Include("keywords"),
			{`[A-Za-z]\w*`, CommentPreproc, nil},
			Include("root"),
		},
		"string": {
			{`\\\(`, LiteralStringInterpol, Push("string-intp")},
			{`"`, LiteralString, Pop(1)},
			{`\\['"\\nrt]|\\x[0-9a-fA-F]{2}|\\[0-7]{1,3}|\\u[0-9a-fA-F]{4}|\\U[0-9a-fA-F]{8}`, LiteralStringEscape, nil},
			{`[^\\"]+`, LiteralString, nil},
			{`\\`, LiteralString, nil},
		},
		"string-intp": {
			{`\(`, LiteralStringInterpol, Push()},
			{`\)`, LiteralStringInterpol, Pop(1)},
			Include("root"),
		},
	}
}
