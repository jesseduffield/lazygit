package color

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/mattn/go-colorable"
)

// Testing colors is kinda different. First we test for given colors and their
// escaped formatted results. Next we create some visual tests to be tested.
// Each visual test includes the color name to be compared.
func TestColor(t *testing.T) {
	rb := new(bytes.Buffer)
	Output = rb

	NoColor = false

	testColors := []struct {
		text string
		code Attribute
	}{
		{text: "black", code: FgBlack},
		{text: "red", code: FgRed},
		{text: "green", code: FgGreen},
		{text: "yellow", code: FgYellow},
		{text: "blue", code: FgBlue},
		{text: "magent", code: FgMagenta},
		{text: "cyan", code: FgCyan},
		{text: "white", code: FgWhite},
		{text: "hblack", code: FgHiBlack},
		{text: "hred", code: FgHiRed},
		{text: "hgreen", code: FgHiGreen},
		{text: "hyellow", code: FgHiYellow},
		{text: "hblue", code: FgHiBlue},
		{text: "hmagent", code: FgHiMagenta},
		{text: "hcyan", code: FgHiCyan},
		{text: "hwhite", code: FgHiWhite},
	}

	for _, c := range testColors {
		New(c.code).Print(c.text)

		line, _ := rb.ReadString('\n')
		scannedLine := fmt.Sprintf("%q", line)
		colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", c.code, c.text)
		escapedForm := fmt.Sprintf("%q", colored)

		fmt.Printf("%s\t: %s\n", c.text, line)

		if scannedLine != escapedForm {
			t.Errorf("Expecting %s, got '%s'\n", escapedForm, scannedLine)
		}
	}

	for _, c := range testColors {
		line := New(c.code).Sprintf("%s", c.text)
		scannedLine := fmt.Sprintf("%q", line)
		colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", c.code, c.text)
		escapedForm := fmt.Sprintf("%q", colored)

		fmt.Printf("%s\t: %s\n", c.text, line)

		if scannedLine != escapedForm {
			t.Errorf("Expecting %s, got '%s'\n", escapedForm, scannedLine)
		}
	}
}

func TestColorEquals(t *testing.T) {
	fgblack1 := New(FgBlack)
	fgblack2 := New(FgBlack)
	bgblack := New(BgBlack)
	fgbgblack := New(FgBlack, BgBlack)
	fgblackbgred := New(FgBlack, BgRed)
	fgred := New(FgRed)
	bgred := New(BgRed)

	if !fgblack1.Equals(fgblack2) {
		t.Error("Two black colors are not equal")
	}

	if fgblack1.Equals(bgblack) {
		t.Error("Fg and bg black colors are equal")
	}

	if fgblack1.Equals(fgbgblack) {
		t.Error("Fg black equals fg/bg black color")
	}

	if fgblack1.Equals(fgred) {
		t.Error("Fg black equals Fg red")
	}

	if fgblack1.Equals(bgred) {
		t.Error("Fg black equals Bg red")
	}

	if fgblack1.Equals(fgblackbgred) {
		t.Error("Fg black equals fg black bg red")
	}
}

func TestNoColor(t *testing.T) {
	rb := new(bytes.Buffer)
	Output = rb

	testColors := []struct {
		text string
		code Attribute
	}{
		{text: "black", code: FgBlack},
		{text: "red", code: FgRed},
		{text: "green", code: FgGreen},
		{text: "yellow", code: FgYellow},
		{text: "blue", code: FgBlue},
		{text: "magent", code: FgMagenta},
		{text: "cyan", code: FgCyan},
		{text: "white", code: FgWhite},
		{text: "hblack", code: FgHiBlack},
		{text: "hred", code: FgHiRed},
		{text: "hgreen", code: FgHiGreen},
		{text: "hyellow", code: FgHiYellow},
		{text: "hblue", code: FgHiBlue},
		{text: "hmagent", code: FgHiMagenta},
		{text: "hcyan", code: FgHiCyan},
		{text: "hwhite", code: FgHiWhite},
	}

	for _, c := range testColors {
		p := New(c.code)
		p.DisableColor()
		p.Print(c.text)

		line, _ := rb.ReadString('\n')
		if line != c.text {
			t.Errorf("Expecting %s, got '%s'\n", c.text, line)
		}
	}

	// global check
	NoColor = true
	defer func() {
		NoColor = false
	}()
	for _, c := range testColors {
		p := New(c.code)
		p.Print(c.text)

		line, _ := rb.ReadString('\n')
		if line != c.text {
			t.Errorf("Expecting %s, got '%s'\n", c.text, line)
		}
	}

}

func TestColorVisual(t *testing.T) {
	// First Visual Test
	Output = colorable.NewColorableStdout()

	New(FgRed).Printf("red\t")
	New(BgRed).Print("         ")
	New(FgRed, Bold).Println(" red")

	New(FgGreen).Printf("green\t")
	New(BgGreen).Print("         ")
	New(FgGreen, Bold).Println(" green")

	New(FgYellow).Printf("yellow\t")
	New(BgYellow).Print("         ")
	New(FgYellow, Bold).Println(" yellow")

	New(FgBlue).Printf("blue\t")
	New(BgBlue).Print("         ")
	New(FgBlue, Bold).Println(" blue")

	New(FgMagenta).Printf("magenta\t")
	New(BgMagenta).Print("         ")
	New(FgMagenta, Bold).Println(" magenta")

	New(FgCyan).Printf("cyan\t")
	New(BgCyan).Print("         ")
	New(FgCyan, Bold).Println(" cyan")

	New(FgWhite).Printf("white\t")
	New(BgWhite).Print("         ")
	New(FgWhite, Bold).Println(" white")
	fmt.Println("")

	// Second Visual test
	Black("black")
	Red("red")
	Green("green")
	Yellow("yellow")
	Blue("blue")
	Magenta("magenta")
	Cyan("cyan")
	White("white")
	HiBlack("hblack")
	HiRed("hred")
	HiGreen("hgreen")
	HiYellow("hyellow")
	HiBlue("hblue")
	HiMagenta("hmagenta")
	HiCyan("hcyan")
	HiWhite("hwhite")

	// Third visual test
	fmt.Println()
	Set(FgBlue)
	fmt.Println("is this blue?")
	Unset()

	Set(FgMagenta)
	fmt.Println("and this magenta?")
	Unset()

	// Fourth Visual test
	fmt.Println()
	blue := New(FgBlue).PrintlnFunc()
	blue("blue text with custom print func")

	red := New(FgRed).PrintfFunc()
	red("red text with a printf func: %d\n", 123)

	put := New(FgYellow).SprintFunc()
	warn := New(FgRed).SprintFunc()

	fmt.Fprintf(Output, "this is a %s and this is %s.\n", put("warning"), warn("error"))

	info := New(FgWhite, BgGreen).SprintFunc()
	fmt.Fprintf(Output, "this %s rocks!\n", info("package"))

	notice := New(FgBlue).FprintFunc()
	notice(os.Stderr, "just a blue notice to stderr")

	// Fifth Visual Test
	fmt.Println()

	fmt.Fprintln(Output, BlackString("black"))
	fmt.Fprintln(Output, RedString("red"))
	fmt.Fprintln(Output, GreenString("green"))
	fmt.Fprintln(Output, YellowString("yellow"))
	fmt.Fprintln(Output, BlueString("blue"))
	fmt.Fprintln(Output, MagentaString("magenta"))
	fmt.Fprintln(Output, CyanString("cyan"))
	fmt.Fprintln(Output, WhiteString("white"))
	fmt.Fprintln(Output, HiBlackString("hblack"))
	fmt.Fprintln(Output, HiRedString("hred"))
	fmt.Fprintln(Output, HiGreenString("hgreen"))
	fmt.Fprintln(Output, HiYellowString("hyellow"))
	fmt.Fprintln(Output, HiBlueString("hblue"))
	fmt.Fprintln(Output, HiMagentaString("hmagenta"))
	fmt.Fprintln(Output, HiCyanString("hcyan"))
	fmt.Fprintln(Output, HiWhiteString("hwhite"))
}

func TestNoFormat(t *testing.T) {
	fmt.Printf("%s   %%s = ", BlackString("Black"))
	Black("%s")

	fmt.Printf("%s     %%s = ", RedString("Red"))
	Red("%s")

	fmt.Printf("%s   %%s = ", GreenString("Green"))
	Green("%s")

	fmt.Printf("%s  %%s = ", YellowString("Yellow"))
	Yellow("%s")

	fmt.Printf("%s    %%s = ", BlueString("Blue"))
	Blue("%s")

	fmt.Printf("%s %%s = ", MagentaString("Magenta"))
	Magenta("%s")

	fmt.Printf("%s    %%s = ", CyanString("Cyan"))
	Cyan("%s")

	fmt.Printf("%s   %%s = ", WhiteString("White"))
	White("%s")

	fmt.Printf("%s   %%s = ", HiBlackString("HiBlack"))
	HiBlack("%s")

	fmt.Printf("%s     %%s = ", HiRedString("HiRed"))
	HiRed("%s")

	fmt.Printf("%s   %%s = ", HiGreenString("HiGreen"))
	HiGreen("%s")

	fmt.Printf("%s  %%s = ", HiYellowString("HiYellow"))
	HiYellow("%s")

	fmt.Printf("%s    %%s = ", HiBlueString("HiBlue"))
	HiBlue("%s")

	fmt.Printf("%s %%s = ", HiMagentaString("HiMagenta"))
	HiMagenta("%s")

	fmt.Printf("%s    %%s = ", HiCyanString("HiCyan"))
	HiCyan("%s")

	fmt.Printf("%s   %%s = ", HiWhiteString("HiWhite"))
	HiWhite("%s")
}

func TestNoFormatString(t *testing.T) {
	tests := []struct {
		f      func(string, ...interface{}) string
		format string
		args   []interface{}
		want   string
	}{
		{BlackString, "%s", nil, "\x1b[30m%s\x1b[0m"},
		{RedString, "%s", nil, "\x1b[31m%s\x1b[0m"},
		{GreenString, "%s", nil, "\x1b[32m%s\x1b[0m"},
		{YellowString, "%s", nil, "\x1b[33m%s\x1b[0m"},
		{BlueString, "%s", nil, "\x1b[34m%s\x1b[0m"},
		{MagentaString, "%s", nil, "\x1b[35m%s\x1b[0m"},
		{CyanString, "%s", nil, "\x1b[36m%s\x1b[0m"},
		{WhiteString, "%s", nil, "\x1b[37m%s\x1b[0m"},
		{HiBlackString, "%s", nil, "\x1b[90m%s\x1b[0m"},
		{HiRedString, "%s", nil, "\x1b[91m%s\x1b[0m"},
		{HiGreenString, "%s", nil, "\x1b[92m%s\x1b[0m"},
		{HiYellowString, "%s", nil, "\x1b[93m%s\x1b[0m"},
		{HiBlueString, "%s", nil, "\x1b[94m%s\x1b[0m"},
		{HiMagentaString, "%s", nil, "\x1b[95m%s\x1b[0m"},
		{HiCyanString, "%s", nil, "\x1b[96m%s\x1b[0m"},
		{HiWhiteString, "%s", nil, "\x1b[97m%s\x1b[0m"},
	}

	for i, test := range tests {
		s := fmt.Sprintf("%s", test.f(test.format, test.args...))
		if s != test.want {
			t.Errorf("[%d] want: %q, got: %q", i, test.want, s)
		}
	}
}
