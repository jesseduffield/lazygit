package b

import (
	"strings"

	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Bicep lexer.
var Bicep = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Bicep",
		Aliases:   []string{"bicep"},
		Filenames: []string{"*.bicep"},
	},
	bicepRules,
))

func bicepRules() Rules {
	bicepFunctions := []string{
		"any",
		"array",
		"concat",
		"contains",
		"empty",
		"first",
		"intersection",
		"items",
		"last",
		"length",
		"min",
		"max",
		"range",
		"skip",
		"take",
		"union",
		"dateTimeAdd",
		"utcNow",
		"deployment",
		"environment",
		"loadFileAsBase64",
		"loadTextContent",
		"int",
		"json",
		"extensionResourceId",
		"getSecret",
		"list",
		"listKeys",
		"listKeyValue",
		"listAccountSas",
		"listSecrets",
		"pickZones",
		"reference",
		"resourceId",
		"subscriptionResourceId",
		"tenantResourceId",
		"managementGroup",
		"resourceGroup",
		"subscription",
		"tenant",
		"base64",
		"base64ToJson",
		"base64ToString",
		"dataUri",
		"dataUriToString",
		"endsWith",
		"format",
		"guid",
		"indexOf",
		"lastIndexOf",
		"length",
		"newGuid",
		"padLeft",
		"replace",
		"split",
		"startsWith",
		"string",
		"substring",
		"toLower",
		"toUpper",
		"trim",
		"uniqueString",
		"uri",
		"uriComponent",
		"uriComponentToString",
	}

	return Rules{
		"root": {
			{`//[^\n\r]+`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`([']?\w+[']?)(:)`, ByGroups(NameProperty, Punctuation), nil},
			{`\b('(resourceGroup|subscription|managementGroup|tenant)')\b`, KeywordNamespace, nil},
			{`'[\w\$\{\(\)\}\.]{1,}?'`, LiteralStringInterpol, nil},
			{`('''|').*?('''|')`, LiteralString, nil},
			{`\b(allowed|batchSize|description|maxLength|maxValue|metadata|minLength|minValue|secure)\b`, NameDecorator, nil},
			{`\b(az|sys)\.`, NameNamespace, nil},
			{`\b(` + strings.Join(bicepFunctions, "|") + `)\b`, NameFunction, nil},
			// https://docs.microsoft.com/en-us/azure/azure-resource-manager/bicep/bicep-functions-logical
			{`\b(bool)(\()`, ByGroups(NameFunction, Punctuation), nil},
			{`\b(for|if|in)\b`, Keyword, nil},
			{`\b(module|output|param|resource|var)\b`, KeywordDeclaration, nil},
			{`\b(array|bool|int|object|string)\b`, KeywordType, nil},
			// https://docs.microsoft.com/en-us/azure/azure-resource-manager/bicep/operators
			{`(>=|>|<=|<|==|!=|=~|!~|::|&&|\?\?|!|-|%|\*|\/|\+)`, Operator, nil},
			{`[\(\)\[\]\.:\?{}@=]`, Punctuation, nil},
			{`[\w_-]+`, Text, nil},
			{`\s+`, TextWhitespace, nil},
		},
	}
}
