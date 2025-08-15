package color

import "fmt"

/*************************************************************
 * colored message Printer
 *************************************************************/

// PrinterFace interface
type PrinterFace interface {
	fmt.Stringer
	Sprint(a ...interface{}) string
	Sprintf(format string, a ...interface{}) string
	Print(a ...interface{})
	Printf(format string, a ...interface{})
	Println(a ...interface{})
}

// Printer a generic color message printer.
//
// Usage:
// 	p := &Printer{Code: "32;45;3"}
// 	p.Print("message")
type Printer struct {
	// NoColor disable color.
	NoColor bool
	// Code color code string. eg "32;45;3"
	Code string
}

// NewPrinter instance
func NewPrinter(colorCode string) *Printer {
	return &Printer{Code: colorCode}
}

// String returns color code string. eg: "32;45;3"
func (p *Printer) String() string {
	// panic("implement me")
	return p.Code
}

// Sprint returns rendering colored messages
func (p *Printer) Sprint(a ...interface{}) string {
	return RenderCode(p.String(), a...)
}

// Sprintf returns format and rendering colored messages
func (p *Printer) Sprintf(format string, a ...interface{}) string {
	return RenderString(p.String(), fmt.Sprintf(format, a...))
}

// Print rendering colored messages
func (p *Printer) Print(a ...interface{}) {
	doPrintV2(p.String(), fmt.Sprint(a...))
}

// Printf format and rendering colored messages
func (p *Printer) Printf(format string, a ...interface{}) {
	doPrintV2(p.String(), fmt.Sprintf(format, a...))
}

// Println rendering colored messages with newline
func (p *Printer) Println(a ...interface{}) {
	doPrintlnV2(p.Code, a)
}

// IsEmpty color code
func (p *Printer) IsEmpty() bool {
	return p.Code == ""
}

/*************************************************************
 * SimplePrinter struct
 *************************************************************/

// SimplePrinter use for quick use color print on inject to struct
type SimplePrinter struct{}

// Print message
func (s *SimplePrinter) Print(v ...interface{}) {
	Print(v...)
}

// Printf message
func (s *SimplePrinter) Printf(format string, v ...interface{}) {
	Printf(format, v...)
}

// Println message
func (s *SimplePrinter) Println(v ...interface{}) {
	Println(v...)
}

// Infof message
func (s *SimplePrinter) Infof(format string, a ...interface{}) {
	Info.Printf(format, a...)
}

// Infoln message
func (s *SimplePrinter) Infoln(a ...interface{}) {
	Info.Println(a...)
}

// Warnf message
func (s *SimplePrinter) Warnf(format string, a ...interface{}) {
	Warn.Printf(format, a...)
}

// Warnln message
func (s *SimplePrinter) Warnln(a ...interface{}) {
	Warn.Println(a...)
}

// Errorf message
func (s *SimplePrinter) Errorf(format string, a ...interface{}) {
	Error.Printf(format, a...)
}

// Errorln message
func (s *SimplePrinter) Errorln(a ...interface{}) {
	Error.Println(a...)
}
