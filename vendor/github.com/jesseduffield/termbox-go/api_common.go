// termbox is a library for creating cross-platform text-based interfaces
package termbox

// public API, common OS agnostic part

type (
	InputMode  int
	OutputMode int
	EventType  uint8
	Modifier   uint8
	Key        uint16
	Attribute  uint16
)

// This type represents a termbox event. The 'Mod', 'Key' and 'Ch' fields are
// valid if 'Type' is EventKey. The 'Width' and 'Height' fields are valid if
// 'Type' is EventResize. The 'Err' field is valid if 'Type' is EventError.
type Event struct {
	Type   EventType // one of Event* constants
	Mod    Modifier  // one of Mod* constants or 0
	Key    Key       // one of Key* constants, invalid if 'Ch' is not 0
	Ch     rune      // a unicode character
	Width  int       // width of the screen
	Height int       // height of the screen
	Err    error     // error in case if input failed
	MouseX int       // x coord of mouse
	MouseY int       // y coord of mouse
	N      int       // number of bytes written when getting a raw event
}

// A cell, single conceptual entity on the screen. The screen is basically a 2d
// array of cells. 'Ch' is a unicode character, 'Fg' and 'Bg' are foreground
// and background attributes respectively.
type Cell struct {
	Ch rune
	Fg Attribute
	Bg Attribute
}

// To know if termbox has been initialized or not
var (
	IsInit bool = false
)

// Key constants, see Event.Key field.
const (
	KeyF1 Key = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
	key_min // see terminfo
	MouseLeft
	MouseMiddle
	MouseRight
	MouseRelease
	MouseWheelUp
	MouseWheelDown
)

const (
	KeyCtrlTilde      Key = 0x00
	KeyCtrl2          Key = 0x00
	KeyCtrlSpace      Key = 0x00
	KeyCtrlA          Key = 0x01
	KeyCtrlB          Key = 0x02
	KeyCtrlC          Key = 0x03
	KeyCtrlD          Key = 0x04
	KeyCtrlE          Key = 0x05
	KeyCtrlF          Key = 0x06
	KeyCtrlG          Key = 0x07
	KeyBackspace      Key = 0x08
	KeyCtrlH          Key = 0x08
	KeyTab            Key = 0x09
	KeyCtrlI          Key = 0x09
	KeyCtrlJ          Key = 0x0A
	KeyCtrlK          Key = 0x0B
	KeyCtrlL          Key = 0x0C
	KeyEnter          Key = 0x0D
	KeyCtrlM          Key = 0x0D
	KeyCtrlN          Key = 0x0E
	KeyCtrlO          Key = 0x0F
	KeyCtrlP          Key = 0x10
	KeyCtrlQ          Key = 0x11
	KeyCtrlR          Key = 0x12
	KeyCtrlS          Key = 0x13
	KeyCtrlT          Key = 0x14
	KeyCtrlU          Key = 0x15
	KeyCtrlV          Key = 0x16
	KeyCtrlW          Key = 0x17
	KeyCtrlX          Key = 0x18
	KeyCtrlY          Key = 0x19
	KeyCtrlZ          Key = 0x1A
	KeyEsc            Key = 0x1B
	KeyCtrlLsqBracket Key = 0x1B
	KeyCtrl3          Key = 0x1B
	KeyCtrl4          Key = 0x1C
	KeyCtrlBackslash  Key = 0x1C
	KeyCtrl5          Key = 0x1D
	KeyCtrlRsqBracket Key = 0x1D
	KeyCtrl6          Key = 0x1E
	KeyCtrl7          Key = 0x1F
	KeyCtrlSlash      Key = 0x1F
	KeyCtrlUnderscore Key = 0x1F
	KeySpace          Key = 0x20
	KeyBackspace2     Key = 0x7F
	KeyCtrl8          Key = 0x7F
)

// Alt modifier constant, see Event.Mod field and SetInputMode function.
const (
	ModAlt Modifier = 1 << iota
	ModMotion
)

// Cell colors, you can combine a color with multiple attributes using bitwise
// OR ('|').
const (
	ColorDefault Attribute = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

// Cell attributes, it is possible to use multiple attributes by combining them
// using bitwise OR ('|'). Although, colors cannot be combined. But you can
// combine attributes and a single color.
//
// It's worth mentioning that some platforms don't support certain attributes.
// For example windows console doesn't support AttrUnderline. And on some
// terminals applying AttrBold to background may result in blinking text. Use
// them with caution and test your code on various terminals.
const (
	AttrBold Attribute = 1 << (iota + 9)
	AttrUnderline
	AttrReverse
)

// Input mode. See SetInputMode function.
const (
	InputEsc InputMode = 1 << iota
	InputAlt
	InputMouse
	InputCurrent InputMode = 0
)

// Output mode. See SetOutputMode function.
const (
	OutputCurrent OutputMode = iota
	OutputNormal
	Output256
	Output216
	OutputGrayscale
)

// Event type. See Event.Type field.
const (
	EventKey EventType = iota
	EventResize
	EventMouse
	EventError
	EventInterrupt
	EventRaw
	EventNone
)
