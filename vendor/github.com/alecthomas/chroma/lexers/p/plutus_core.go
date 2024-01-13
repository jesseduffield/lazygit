package p

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// nolint

// Lexer for the Plutus Core Languages (version 2.1)
//
// including both Typed- and Untyped- versions
// based on “Formal Specification of the Plutus Core Language (version 2.1)”, published 6th April 2021:
// https://hydra.iohk.io/build/8205579/download/1/plutus-core-specification.pdf

var PlutusCoreLang = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Plutus Core",
		Aliases:   []string{"plutus-core", "plc"},
		Filenames: []string{"*.plc"},
		MimeTypes: []string{"text/x-plutus-core", "application/x-plutus-core"},
	},
	plutusCoreRules,
))

func plutusCoreRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`(\(|\))`, Punctuation, nil},
			{`(\[|\])`, Punctuation, nil},
			{`({|})`, Punctuation, nil},

			// Constants. Figure 1.
			// For version, see handling of (program ...) below.
			{`([+-]?\d+)`, LiteralNumberInteger, nil},
			{`(#([a-fA-F0-9][a-fA-F0-9])+)`, LiteralString, nil},
			{`(\(\))`, NameConstant, nil},
			{`(True|False)`, NameConstant, nil},

			// Keywords. Figures 2 and 15.
			// Special handling for program because it is followed by a version.
			{`(con |abs |iwrap |unwrap |lam |builtin |delay |force |error)`, Keyword, nil},
			{`(fun |all |ifix |lam |con )`, Keyword, nil},
			{`(type|fun )`, Keyword, nil},
			{`(program )(\S+)`, ByGroups(Keyword, LiteralString), nil},

			// Built-in Types. Figure 12.
			{`(unit|bool|integer|bytestring|string)`, KeywordType, nil},

			// Built-ins Functions. Figure 14 but, more importantly, implementation:
			// https://github.com/input-output-hk/plutus/blob/6d759c4/plutus-core/plutus-core/src/PlutusCore/Default/Builtins.hs#L42-L111
			{`(addInteger |subtractInteger |multiplyInteger |divideInteger |quotientInteger |remainderInteger |modInteger |equalsInteger |lessThanInteger |lessThanEqualsInteger )`, NameBuiltin, nil},
			{`(appendByteString |consByteString |sliceByteString |lengthOfByteString |indexByteString |equalsByteString |lessThanByteString |lessThanEqualsByteString )`, NameBuiltin, nil},
			{`(sha2_256 |sha3_256 |blake2b_256 |verifySignature )`, NameBuiltin, nil},
			{`(appendString |equalsString |encodeUtf8 |decodeUtf8 )`, NameBuiltin, nil},
			{`(ifThenElse )`, NameBuiltin, nil},
			{`(chooseUnit )`, NameBuiltin, nil},
			{`(trace )`, NameBuiltin, nil},
			{`(fstPair |sndPair )`, NameBuiltin, nil},
			{`(chooseList |mkCons |headList |tailList |nullList )`, NameBuiltin, nil},
			{`(chooseData |constrData |mapData |listData |iData |bData |unConstrData |unMapData |unListData |unIData |unBData |equalsData )`, NameBuiltin, nil},
			{`(mkPairData |mkNilData |mkNilPairData )`, NameBuiltin, nil},

			// Name. Figure 1.
			{`([a-zA-Z][a-zA-Z0-9_']*)`, Name, nil},

			// Unicode String. Not in the specification.
			{`"`, LiteralStringDouble, Push("string")},
		},
		"string": {
			{`[^\\"]+`, LiteralStringDouble, nil},
			{`"`, LiteralStringDouble, Pop(1)},
		},
	}
}
