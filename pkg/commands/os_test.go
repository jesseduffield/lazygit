package commands

import "testing"

func TestQuote(t *testing.T) {
	osCommand := &OSCommand{
		Log:      nil,
		Platform: getPlatform(),
	}
	test := "hello `test`"
	expected := osCommand.Platform.escapedQuote + "hello \\`test\\`" + osCommand.Platform.escapedQuote
	test = osCommand.Quote(test)
	if test != expected {
		t.Error("Expected " + expected + ", got " + test)
	}
}
