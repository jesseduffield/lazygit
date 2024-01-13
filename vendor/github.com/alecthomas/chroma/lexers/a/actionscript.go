package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Actionscript lexer.
var Actionscript = internal.Register(MustNewLazyLexer(
	&Config{
		Name:         "ActionScript",
		Aliases:      []string{"as", "actionscript"},
		Filenames:    []string{"*.as"},
		MimeTypes:    []string{"application/x-actionscript", "text/x-actionscript", "text/actionscript"},
		NotMultiline: true,
		DotAll:       true,
	},
	actionscriptRules,
))

func actionscriptRules() Rules {
	return Rules{
		"root": {
			{`\s+`, Text, nil},
			{`//.*?\n`, CommentSingle, nil},
			{`/\*.*?\*/`, CommentMultiline, nil},
			{`/(\\\\|\\/|[^/\n])*/[gim]*`, LiteralStringRegex, nil},
			{`[~^*!%&<>|+=:;,/?\\-]+`, Operator, nil},
			{`[{}\[\]();.]+`, Punctuation, nil},
			{Words(``, `\b`, `case`, `default`, `for`, `each`, `in`, `while`, `do`, `break`, `return`, `continue`, `if`, `else`, `throw`, `try`, `catch`, `var`, `with`, `new`, `typeof`, `arguments`, `instanceof`, `this`, `switch`), Keyword, nil},
			{Words(``, `\b`, `class`, `public`, `final`, `internal`, `native`, `override`, `private`, `protected`, `static`, `import`, `extends`, `implements`, `interface`, `intrinsic`, `return`, `super`, `dynamic`, `function`, `const`, `get`, `namespace`, `package`, `set`), KeywordDeclaration, nil},
			{`(true|false|null|NaN|Infinity|-Infinity|undefined|Void)\b`, KeywordConstant, nil},
			{Words(``, `\b`, `Accessibility`, `AccessibilityProperties`, `ActionScriptVersion`, `ActivityEvent`, `AntiAliasType`, `ApplicationDomain`, `AsBroadcaster`, `Array`, `AsyncErrorEvent`, `AVM1Movie`, `BevelFilter`, `Bitmap`, `BitmapData`, `BitmapDataChannel`, `BitmapFilter`, `BitmapFilterQuality`, `BitmapFilterType`, `BlendMode`, `BlurFilter`, `Boolean`, `ByteArray`, `Camera`, `Capabilities`, `CapsStyle`, `Class`, `Color`, `ColorMatrixFilter`, `ColorTransform`, `ContextMenu`, `ContextMenuBuiltInItems`, `ContextMenuEvent`, `ContextMenuItem`, `ConvultionFilter`, `CSMSettings`, `DataEvent`, `Date`, `DefinitionError`, `DeleteObjectSample`, `Dictionary`, `DisplacmentMapFilter`, `DisplayObject`, `DisplacmentMapFilterMode`, `DisplayObjectContainer`, `DropShadowFilter`, `Endian`, `EOFError`, `Error`, `ErrorEvent`, `EvalError`, `Event`, `EventDispatcher`, `EventPhase`, `ExternalInterface`, `FileFilter`, `FileReference`, `FileReferenceList`, `FocusDirection`, `FocusEvent`, `Font`, `FontStyle`, `FontType`, `FrameLabel`, `FullScreenEvent`, `Function`, `GlowFilter`, `GradientBevelFilter`, `GradientGlowFilter`, `GradientType`, `Graphics`, `GridFitType`, `HTTPStatusEvent`, `IBitmapDrawable`, `ID3Info`, `IDataInput`, `IDataOutput`, `IDynamicPropertyOutputIDynamicPropertyWriter`, `IEventDispatcher`, `IExternalizable`, `IllegalOperationError`, `IME`, `IMEConversionMode`, `IMEEvent`, `int`, `InteractiveObject`, `InterpolationMethod`, `InvalidSWFError`, `InvokeEvent`, `IOError`, `IOErrorEvent`, `JointStyle`, `Key`, `Keyboard`, `KeyboardEvent`, `KeyLocation`, `LineScaleMode`, `Loader`, `LoaderContext`, `LoaderInfo`, `LoadVars`, `LocalConnection`, `Locale`, `Math`, `Matrix`, `MemoryError`, `Microphone`, `MorphShape`, `Mouse`, `MouseEvent`, `MovieClip`, `MovieClipLoader`, `Namespace`, `NetConnection`, `NetStatusEvent`, `NetStream`, `NewObjectSample`, `Number`, `Object`, `ObjectEncoding`, `PixelSnapping`, `Point`, `PrintJob`, `PrintJobOptions`, `PrintJobOrientation`, `ProgressEvent`, `Proxy`, `QName`, `RangeError`, `Rectangle`, `ReferenceError`, `RegExp`, `Responder`, `Sample`, `Scene`, `ScriptTimeoutError`, `Security`, `SecurityDomain`, `SecurityError`, `SecurityErrorEvent`, `SecurityPanel`, `Selection`, `Shape`, `SharedObject`, `SharedObjectFlushStatus`, `SimpleButton`, `Socket`, `Sound`, `SoundChannel`, `SoundLoaderContext`, `SoundMixer`, `SoundTransform`, `SpreadMethod`, `Sprite`, `StackFrame`, `StackOverflowError`, `Stage`, `StageAlign`, `StageDisplayState`, `StageQuality`, `StageScaleMode`, `StaticText`, `StatusEvent`, `String`, `StyleSheet`, `SWFVersion`, `SyncEvent`, `SyntaxError`, `System`, `TextColorType`, `TextField`, `TextFieldAutoSize`, `TextFieldType`, `TextFormat`, `TextFormatAlign`, `TextLineMetrics`, `TextRenderer`, `TextSnapshot`, `Timer`, `TimerEvent`, `Transform`, `TypeError`, `uint`, `URIError`, `URLLoader`, `URLLoaderDataFormat`, `URLRequest`, `URLRequestHeader`, `URLRequestMethod`, `URLStream`, `URLVariabeles`, `VerifyError`, `Video`, `XML`, `XMLDocument`, `XMLList`, `XMLNode`, `XMLNodeType`, `XMLSocket`, `XMLUI`), NameBuiltin, nil},
			{Words(``, `\b`, `decodeURI`, `decodeURIComponent`, `encodeURI`, `escape`, `eval`, `isFinite`, `isNaN`, `isXMLName`, `clearInterval`, `fscommand`, `getTimer`, `getURL`, `getVersion`, `parseFloat`, `parseInt`, `setInterval`, `trace`, `updateAfterEvent`, `unescape`), NameFunction, nil},
			{`[$a-zA-Z_]\w*`, NameOther, nil},
			{`[0-9][0-9]*\.[0-9]+([eE][0-9]+)?[fd]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-f]+`, LiteralNumberHex, nil},
			{`[0-9]+`, LiteralNumberInteger, nil},
			{`"(\\\\|\\"|[^"])*"`, LiteralStringDouble, nil},
			{`'(\\\\|\\'|[^'])*'`, LiteralStringSingle, nil},
		},
	}
}
