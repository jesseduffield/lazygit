package a

import (
	. "github.com/alecthomas/chroma" // nolint
	"github.com/alecthomas/chroma/lexers/internal"
)

// Arduino lexer.
var Arduino = internal.Register(MustNewLazyLexer(
	&Config{
		Name:      "Arduino",
		Aliases:   []string{"arduino"},
		Filenames: []string{"*.ino"},
		MimeTypes: []string{"text/x-arduino"},
		EnsureNL:  true,
	},
	arduinoRules,
))

func arduinoRules() Rules {
	return Rules{
		"statements": {
			{Words(``, `\b`, `catch`, `const_cast`, `delete`, `dynamic_cast`, `explicit`, `export`, `friend`, `mutable`, `namespace`, `new`, `operator`, `private`, `protected`, `public`, `reinterpret_cast`, `restrict`, `static_cast`, `template`, `this`, `throw`, `throws`, `try`, `typeid`, `typename`, `using`, `virtual`, `constexpr`, `nullptr`, `decltype`, `thread_local`, `alignas`, `alignof`, `static_assert`, `noexcept`, `override`, `final`), Keyword, nil},
			{`char(16_t|32_t)\b`, KeywordType, nil},
			{`(class)\b`, ByGroups(Keyword, Text), Push("classname")},
			{`(R)(")([^\\()\s]{,16})(\()((?:.|\n)*?)(\)\3)(")`, ByGroups(LiteralStringAffix, LiteralString, LiteralStringDelimiter, LiteralStringDelimiter, LiteralString, LiteralStringDelimiter, LiteralString), nil},
			{`(u8|u|U)(")`, ByGroups(LiteralStringAffix, LiteralString), Push("string")},
			{`(L?)(")`, ByGroups(LiteralStringAffix, LiteralString), Push("string")},
			{`(L?)(')(\\.|\\[0-7]{1,3}|\\x[a-fA-F0-9]{1,2}|[^\\\'\n])(')`, ByGroups(LiteralStringAffix, LiteralStringChar, LiteralStringChar, LiteralStringChar), nil},
			{`(\d+\.\d*|\.\d+|\d+)[eE][+-]?\d+[LlUu]*`, LiteralNumberFloat, nil},
			{`(\d+\.\d*|\.\d+|\d+[fF])[fF]?`, LiteralNumberFloat, nil},
			{`0x[0-9a-fA-F]+[LlUu]*`, LiteralNumberHex, nil},
			{`0[0-7]+[LlUu]*`, LiteralNumberOct, nil},
			{`\d+[LlUu]*`, LiteralNumberInteger, nil},
			{`\*/`, Error, nil},
			{`[~!%^&*+=|?:<>/-]`, Operator, nil},
			{`[()\[\],.]`, Punctuation, nil},
			{Words(``, `\b`, `asm`, `auto`, `break`, `case`, `const`, `continue`, `default`, `do`, `else`, `enum`, `extern`, `for`, `goto`, `if`, `register`, `restricted`, `return`, `sizeof`, `static`, `struct`, `switch`, `typedef`, `union`, `volatile`, `while`), Keyword, nil},
			{`(_Bool|_Complex|_Imaginary|array|atomic_bool|atomic_char|atomic_int|atomic_llong|atomic_long|atomic_schar|atomic_short|atomic_uchar|atomic_uint|atomic_ullong|atomic_ulong|atomic_ushort|auto|bool|boolean|BooleanVariables|Byte|byte|Char|char|char16_t|char32_t|class|complex|Const|const|const_cast|delete|double|dynamic_cast|enum|explicit|extern|Float|float|friend|inline|Int|int|int16_t|int32_t|int64_t|int8_t|Long|long|new|NULL|null|operator|private|PROGMEM|protected|public|register|reinterpret_cast|short|signed|sizeof|Static|static|static_cast|String|struct|typedef|uint16_t|uint32_t|uint64_t|uint8_t|union|unsigned|virtual|Void|void|Volatile|volatile|word)\b`, KeywordType, nil},
			// Start of: Arduino-specific syntax
			{`(and|final|If|Loop|loop|not|or|override|setup|Setup|throw|try|xor)\b`, Keyword, nil}, // Addition to keywords already defined by C++
			{`(ANALOG_MESSAGE|BIN|CHANGE|DEC|DEFAULT|DIGITAL_MESSAGE|EXTERNAL|FALLING|FIRMATA_STRING|HALF_PI|HEX|HIGH|INPUT|INPUT_PULLUP|INTERNAL|INTERNAL1V1|INTERNAL1V1|INTERNAL2V56|INTERNAL2V56|LED_BUILTIN|LED_BUILTIN_RX|LED_BUILTIN_TX|LOW|LSBFIRST|MSBFIRST|OCT|OUTPUT|PI|REPORT_ANALOG|REPORT_DIGITAL|RISING|SET_PIN_MODE|SYSEX_START|SYSTEM_RESET|TWO_PI)\b`, KeywordConstant, nil},
			{`(boolean|const|byte|word|string|String|array)\b`, NameVariable, nil},
			{`(Keyboard|KeyboardController|MouseController|SoftwareSerial|EthernetServer|EthernetClient|LiquidCrystal|RobotControl|GSMVoiceCall|EthernetUDP|EsploraTFT|HttpClient|RobotMotor|WiFiClient|GSMScanner|FileSystem|Scheduler|GSMServer|YunClient|YunServer|IPAddress|GSMClient|GSMModem|Keyboard|Ethernet|Console|GSMBand|Esplora|Stepper|Process|WiFiUDP|GSM_SMS|Mailbox|USBHost|Firmata|PImage|Client|Server|GSMPIN|FileIO|Bridge|Serial|EEPROM|Stream|Mouse|Audio|Servo|File|Task|GPRS|WiFi|Wire|TFT|GSM|SPI|SD)\b`, NameClass, nil},
			{`(abs|Abs|accept|ACos|acos|acosf|addParameter|analogRead|AnalogRead|analogReadResolution|AnalogReadResolution|analogReference|AnalogReference|analogWrite|AnalogWrite|analogWriteResolution|AnalogWriteResolution|answerCall|asin|ASin|asinf|atan|ATan|atan2|ATan2|atan2f|atanf|attach|attached|attachGPRS|attachInterrupt|AttachInterrupt|autoscroll|available|availableForWrite|background|beep|begin|beginPacket|beginSD|beginSMS|beginSpeaker|beginTFT|beginTransmission|beginWrite|bit|Bit|BitClear|bitClear|bitRead|BitRead|bitSet|BitSet|BitWrite|bitWrite|blink|blinkVersion|BSSID|buffer|byte|cbrt|cbrtf|Ceil|ceil|ceilf|changePIN|char|charAt|checkPIN|checkPUK|checkReg|circle|cityNameRead|cityNameWrite|clear|clearScreen|click|close|compareTo|compassRead|concat|config|connect|connected|constrain|Constrain|copysign|copysignf|cos|Cos|cosf|cosh|coshf|countryNameRead|countryNameWrite|createChar|cursor|debugPrint|degrees|Delay|delay|DelayMicroseconds|delayMicroseconds|detach|DetachInterrupt|detachInterrupt|DigitalPinToInterrupt|digitalPinToInterrupt|DigitalRead|digitalRead|DigitalWrite|digitalWrite|disconnect|display|displayLogos|drawBMP|drawCompass|encryptionType|end|endPacket|endSMS|endsWith|endTransmission|endWrite|equals|equalsIgnoreCase|exists|exitValue|Exp|exp|expf|fabs|fabsf|fdim|fdimf|fill|find|findUntil|float|floor|Floor|floorf|flush|fma|fmaf|fmax|fmaxf|fmin|fminf|fmod|fmodf|gatewayIP|get|getAsynchronously|getBand|getButton|getBytes|getCurrentCarrier|getIMEI|getKey|getModifiers|getOemKey|getPINUsed|getResult|getSignalStrength|getSocket|getVoiceCallStatus|getXChange|getYChange|hangCall|height|highByte|HighByte|home|hypot|hypotf|image|indexOf|int|interrupts|IPAddress|IRread|isActionDone|isAlpha|isAlphaNumeric|isAscii|isControl|isDigit|isDirectory|isfinite|isGraph|isHexadecimalDigit|isinf|isListening|isLowerCase|isnan|isPIN|isPressed|isPrintable|isPunct|isSpace|isUpperCase|isValid|isWhitespace|keyboardRead|keyPressed|keyReleased|knobRead|lastIndexOf|ldexp|ldexpf|leftToRight|length|line|lineFollowConfig|listen|listenOnLocalhost|loadImage|localIP|log|Log|log10|log10f|logf|long|lowByte|LowByte|lrint|lrintf|lround|lroundf|macAddress|maintain|map|Map|Max|max|messageAvailable|Micros|micros|millis|Millis|Min|min|mkdir|motorsStop|motorsWrite|mouseDragged|mouseMoved|mousePressed|mouseReleased|move|noAutoscroll|noBlink|noBuffer|noCursor|noDisplay|noFill|noInterrupts|NoInterrupts|noListenOnLocalhost|noStroke|noTone|NoTone|onReceive|onRequest|open|openNextFile|overflow|parseCommand|parseFloat|parseInt|parsePacket|pauseMode|peek|PinMode|pinMode|playFile|playMelody|point|pointTo|position|Pow|pow|powf|prepare|press|print|printFirmwareVersion|println|printVersion|process|processInput|PulseIn|pulseIn|pulseInLong|PulseInLong|put|radians|random|Random|randomSeed|RandomSeed|read|readAccelerometer|readBlue|readButton|readBytes|readBytesUntil|readGreen|readJoystickButton|readJoystickSwitch|readJoystickX|readJoystickY|readLightSensor|readMessage|readMicrophone|readNetworks|readRed|readSlider|readString|readStringUntil|readTemperature|ready|rect|release|releaseAll|remoteIP|remoteNumber|remotePort|remove|replace|requestFrom|retrieveCallingNumber|rewindDirectory|rightToLeft|rmdir|robotNameRead|robotNameWrite|round|roundf|RSSI|run|runAsynchronously|running|runShellCommand|runShellCommandAsynchronously|scanNetworks|scrollDisplayLeft|scrollDisplayRight|seek|sendAnalog|sendDigitalPortPair|sendDigitalPorts|sendString|sendSysex|Serial_Available|Serial_Begin|Serial_End|Serial_Flush|Serial_Peek|Serial_Print|Serial_Println|Serial_Read|serialEvent|setBand|setBitOrder|setCharAt|setClockDivider|setCursor|setDataMode|setDNS|setFirmwareVersion|setMode|setPINUsed|setSpeed|setTextSize|setTimeout|ShiftIn|shiftIn|ShiftOut|shiftOut|shutdown|signbit|sin|Sin|sinf|sinh|sinhf|size|sizeof|Sq|sq|Sqrt|sqrt|sqrtf|SSID|startLoop|startsWith|step|stop|stroke|subnetMask|substring|switchPIN|tan|Tan|tanf|tanh|tanhf|tempoWrite|text|toCharArray|toInt|toLowerCase|tone|Tone|toUpperCase|transfer|trim|trunc|truncf|tuneWrite|turn|updateIR|userNameRead|userNameWrite|voiceCall|waitContinue|width|WiFiServer|word|write|writeBlue|writeGreen|writeJSON|writeMessage|writeMicroseconds|writeRed|writeRGB|yield|Yield)\b`, NameFunction, nil},
			// End of: Arduino-specific syntax
			{Words(``, `\b`, `inline`, `_inline`, `__inline`, `naked`, `restrict`, `thread`, `typename`), KeywordReserved, nil},
			{`(__m(128i|128d|128|64))\b`, KeywordReserved, nil},
			{Words(`__`, `\b`, `asm`, `int8`, `based`, `except`, `int16`, `stdcall`, `cdecl`, `fastcall`, `int32`, `declspec`, `finally`, `int64`, `try`, `leave`, `wchar_t`, `w64`, `unaligned`, `raise`, `noop`, `identifier`, `forceinline`, `assume`), KeywordReserved, nil},
			{`(true|false|NULL)\b`, NameBuiltin, nil},
			{`([a-zA-Z_]\w*)(\s*)(:)(?!:)`, ByGroups(NameLabel, Text, Punctuation), nil},
			{`[a-zA-Z_]\w*`, Name, nil},
		},
		"root": {
			Include("whitespace"),
			{`((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;{]*)(\{)`, ByGroups(UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), Push("function")},
			{`((?:[\w*\s])+?(?:\s|[*]))([a-zA-Z_]\w*)(\s*\([^;]*?\))([^;]*)(;)`, ByGroups(UsingSelf("root"), NameFunction, UsingSelf("root"), UsingSelf("root"), Punctuation), nil},
			Default(Push("statement")),
			{Words(`__`, `\b`, `virtual_inheritance`, `uuidof`, `super`, `single_inheritance`, `multiple_inheritance`, `interface`, `event`), KeywordReserved, nil},
			{`__(offload|blockingoffload|outer)\b`, KeywordPseudo, nil},
		},
		"classname": {
			{`[a-zA-Z_]\w*`, NameClass, Pop(1)},
			{`\s*(?=>)`, Text, Pop(1)},
		},
		"whitespace": {
			{`^#if\s+0`, CommentPreproc, Push("if0")},
			{`^#`, CommentPreproc, Push("macro")},
			{`^(\s*(?:/[*].*?[*]/\s*)?)(#if\s+0)`, ByGroups(UsingSelf("root"), CommentPreproc), Push("if0")},
			{`^(\s*(?:/[*].*?[*]/\s*)?)(#)`, ByGroups(UsingSelf("root"), CommentPreproc), Push("macro")},
			{`\n`, Text, nil},
			{`\s+`, Text, nil},
			{`\\\n`, Text, nil},
			{`//(\n|[\w\W]*?[^\\]\n)`, CommentSingle, nil},
			{`/(\\\n)?[*][\w\W]*?[*](\\\n)?/`, CommentMultiline, nil},
			{`/(\\\n)?[*][\w\W]*`, CommentMultiline, nil},
		},
		"statement": {
			Include("whitespace"),
			Include("statements"),
			{`[{}]`, Punctuation, nil},
			{`;`, Punctuation, Pop(1)},
		},
		"function": {
			Include("whitespace"),
			Include("statements"),
			{`;`, Punctuation, nil},
			{`\{`, Punctuation, Push()},
			{`\}`, Punctuation, Pop(1)},
		},
		"string": {
			{`"`, LiteralString, Pop(1)},
			{`\\([\\abfnrtv"\']|x[a-fA-F0-9]{2,4}|u[a-fA-F0-9]{4}|U[a-fA-F0-9]{8}|[0-7]{1,3})`, LiteralStringEscape, nil},
			{`[^\\"\n]+`, LiteralString, nil},
			{`\\\n`, LiteralString, nil},
			{`\\`, LiteralString, nil},
		},
		"macro": {
			{`(include)(\s*(?:/[*].*?[*]/\s*)?)([^\n]+)`, ByGroups(CommentPreproc, Text, CommentPreprocFile), nil},
			{`[^/\n]+`, CommentPreproc, nil},
			{`/[*](.|\n)*?[*]/`, CommentMultiline, nil},
			{`//.*?\n`, CommentSingle, Pop(1)},
			{`/`, CommentPreproc, nil},
			{`(?<=\\)\n`, CommentPreproc, nil},
			{`\n`, CommentPreproc, Pop(1)},
		},
		"if0": {
			{`^\s*#if.*?(?<!\\)\n`, CommentPreproc, Push()},
			{`^\s*#el(?:se|if).*\n`, CommentPreproc, Pop(1)},
			{`^\s*#endif.*?(?<!\\)\n`, CommentPreproc, Pop(1)},
			{`.*?\n`, Comment, nil},
		},
	}
}
