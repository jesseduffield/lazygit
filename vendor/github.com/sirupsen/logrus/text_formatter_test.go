package logrus

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatting(t *testing.T) {
	tf := &TextFormatter{DisableColors: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`foo`, "time=\"0001-01-01T00:00:00Z\" level=panic test=foo\n"},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))

		if string(b) != tc.expected {
			t.Errorf("formatting expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestQuoting(t *testing.T) {
	tf := &TextFormatter{DisableColors: true}

	checkQuoting := func(q bool, value interface{}) {
		b, _ := tf.Format(WithField("test", value))
		idx := bytes.Index(b, ([]byte)("test="))
		cont := bytes.Contains(b[idx+5:], []byte("\""))
		if cont != q {
			if q {
				t.Errorf("quoting expected for: %#v", value)
			} else {
				t.Errorf("quoting not expected for: %#v", value)
			}
		}
	}

	checkQuoting(false, "")
	checkQuoting(false, "abcd")
	checkQuoting(false, "v1.0")
	checkQuoting(false, "1234567890")
	checkQuoting(false, "/foobar")
	checkQuoting(false, "foo_bar")
	checkQuoting(false, "foo@bar")
	checkQuoting(false, "foobar^")
	checkQuoting(false, "+/-_^@f.oobar")
	checkQuoting(true, "foobar$")
	checkQuoting(true, "&foobar")
	checkQuoting(true, "x y")
	checkQuoting(true, "x,y")
	checkQuoting(false, errors.New("invalid"))
	checkQuoting(true, errors.New("invalid argument"))

	// Test for quoting empty fields.
	tf.QuoteEmptyFields = true
	checkQuoting(true, "")
	checkQuoting(false, "abcd")
	checkQuoting(true, errors.New("invalid argument"))
}

func TestEscaping(t *testing.T) {
	tf := &TextFormatter{DisableColors: true}

	testCases := []struct {
		value    string
		expected string
	}{
		{`ba"r`, `ba\"r`},
		{`ba'r`, `ba'r`},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestEscaping_Interface(t *testing.T) {
	tf := &TextFormatter{DisableColors: true}

	ts := time.Now()

	testCases := []struct {
		value    interface{}
		expected string
	}{
		{ts, fmt.Sprintf("\"%s\"", ts.String())},
		{errors.New("error: something went wrong"), "\"error: something went wrong\""},
	}

	for _, tc := range testCases {
		b, _ := tf.Format(WithField("test", tc.value))
		if !bytes.Contains(b, []byte(tc.expected)) {
			t.Errorf("escaping expected for %q (result was %q instead of %q)", tc.value, string(b), tc.expected)
		}
	}
}

func TestTimestampFormat(t *testing.T) {
	checkTimeStr := func(format string) {
		customFormatter := &TextFormatter{DisableColors: true, TimestampFormat: format}
		customStr, _ := customFormatter.Format(WithField("test", "test"))
		timeStart := bytes.Index(customStr, ([]byte)("time="))
		timeEnd := bytes.Index(customStr, ([]byte)("level="))
		timeStr := customStr[timeStart+5+len("\"") : timeEnd-1-len("\"")]
		if format == "" {
			format = time.RFC3339
		}
		_, e := time.Parse(format, (string)(timeStr))
		if e != nil {
			t.Errorf("time string \"%s\" did not match provided time format \"%s\": %s", timeStr, format, e)
		}
	}

	checkTimeStr("2006-01-02T15:04:05.000000000Z07:00")
	checkTimeStr("Mon Jan _2 15:04:05 2006")
	checkTimeStr("")
}

func TestDisableLevelTruncation(t *testing.T) {
	entry := &Entry{
		Time:    time.Now(),
		Message: "testing",
	}
	keys := []string{}
	timestampFormat := "Mon Jan 2 15:04:05 -0700 MST 2006"
	checkDisableTruncation := func(disabled bool, level Level) {
		tf := &TextFormatter{DisableLevelTruncation: disabled}
		var b bytes.Buffer
		entry.Level = level
		tf.printColored(&b, entry, keys, nil, timestampFormat)
		logLine := (&b).String()
		if disabled {
			expected := strings.ToUpper(level.String())
			if !strings.Contains(logLine, expected) {
				t.Errorf("level string expected to be %s when truncation disabled", expected)
			}
		} else {
			expected := strings.ToUpper(level.String())
			if len(level.String()) > 4 {
				if strings.Contains(logLine, expected) {
					t.Errorf("level string %s expected to be truncated to %s when truncation is enabled", expected, expected[0:4])
				}
			} else {
				if !strings.Contains(logLine, expected) {
					t.Errorf("level string expected to be %s when truncation is enabled and level string is below truncation threshold", expected)
				}
			}
		}
	}

	checkDisableTruncation(true, DebugLevel)
	checkDisableTruncation(true, InfoLevel)
	checkDisableTruncation(false, ErrorLevel)
	checkDisableTruncation(false, InfoLevel)
}

func TestDisableTimestampWithColoredOutput(t *testing.T) {
	tf := &TextFormatter{DisableTimestamp: true, ForceColors: true}

	b, _ := tf.Format(WithField("test", "test"))
	if strings.Contains(string(b), "[0000]") {
		t.Error("timestamp not expected when DisableTimestamp is true")
	}
}

func TestNewlineBehavior(t *testing.T) {
	tf := &TextFormatter{ForceColors: true}

	// Ensure a single new line is removed as per stdlib log
	e := NewEntry(StandardLogger())
	e.Message = "test message\n"
	b, _ := tf.Format(e)
	if bytes.Contains(b, []byte("test message\n")) {
		t.Error("first newline at end of Entry.Message resulted in unexpected 2 newlines in output. Expected newline to be removed.")
	}

	// Ensure a double new line is reduced to a single new line
	e = NewEntry(StandardLogger())
	e.Message = "test message\n\n"
	b, _ = tf.Format(e)
	if bytes.Contains(b, []byte("test message\n\n")) {
		t.Error("Double newline at end of Entry.Message resulted in unexpected 2 newlines in output. Expected single newline")
	}
	if !bytes.Contains(b, []byte("test message\n")) {
		t.Error("Double newline at end of Entry.Message did not result in a single newline after formatting")
	}
}

func TestTextFormatterFieldMap(t *testing.T) {
	formatter := &TextFormatter{
		DisableColors: true,
		FieldMap: FieldMap{
			FieldKeyMsg:   "message",
			FieldKeyLevel: "somelevel",
			FieldKeyTime:  "timeywimey",
		},
	}

	entry := &Entry{
		Message: "oh hi",
		Level:   WarnLevel,
		Time:    time.Date(1981, time.February, 24, 4, 28, 3, 100, time.UTC),
		Data: Fields{
			"field1":     "f1",
			"message":    "messagefield",
			"somelevel":  "levelfield",
			"timeywimey": "timeywimeyfield",
		},
	}

	b, err := formatter.Format(entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	assert.Equal(t,
		`timeywimey="1981-02-24T04:28:03Z" `+
			`somelevel=warning `+
			`message="oh hi" `+
			`field1=f1 `+
			`fields.message=messagefield `+
			`fields.somelevel=levelfield `+
			`fields.timeywimey=timeywimeyfield`+"\n",
		string(b),
		"Formatted output doesn't respect FieldMap")
}

func TestTextFormatterIsColored(t *testing.T) {
	params := []struct {
		name               string
		expectedResult     bool
		isTerminal         bool
		disableColor       bool
		forceColor         bool
		envColor           bool
		clicolorIsSet      bool
		clicolorForceIsSet bool
		clicolorVal        string
		clicolorForceVal   string
	}{
		// Default values
		{
			name:               "testcase1",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal
		{
			name:               "testcase2",
			expectedResult:     true,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal with color disabled
		{
			name:               "testcase3",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       true,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output not on terminal with color disabled
		{
			name:               "testcase4",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       true,
			forceColor:         false,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output not on terminal with color forced
		{
			name:               "testcase5",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           false,
			clicolorIsSet:      false,
			clicolorForceIsSet: false,
		},
		// Output on terminal with clicolor set to "0"
		{
			name:               "testcase6",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output on terminal with clicolor set to "1"
		{
			name:               "testcase7",
			expectedResult:     true,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "0"
		{
			name:               "testcase8",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output not on terminal with clicolor set to "1"
		{
			name:               "testcase9",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "1" and force color
		{
			name:               "testcase10",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "1",
		},
		// Output not on terminal with clicolor set to "0" and force color
		{
			name:               "testcase11",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         true,
			envColor:           true,
			clicolorIsSet:      true,
			clicolorForceIsSet: false,
			clicolorVal:        "0",
		},
		// Output not on terminal with clicolor_force set to "1"
		{
			name:               "testcase12",
			expectedResult:     true,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "1",
		},
		// Output not on terminal with clicolor_force set to "0"
		{
			name:               "testcase13",
			expectedResult:     false,
			isTerminal:         false,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "0",
		},
		// Output on terminal with clicolor_force set to "0"
		{
			name:               "testcase14",
			expectedResult:     false,
			isTerminal:         true,
			disableColor:       false,
			forceColor:         false,
			envColor:           true,
			clicolorIsSet:      false,
			clicolorForceIsSet: true,
			clicolorForceVal:   "0",
		},
	}

	cleanenv := func() {
		os.Unsetenv("CLICOLOR")
		os.Unsetenv("CLICOLOR_FORCE")
	}

	defer cleanenv()

	for _, val := range params {
		t.Run("textformatter_"+val.name, func(subT *testing.T) {
			tf := TextFormatter{
				isTerminal:                val.isTerminal,
				DisableColors:             val.disableColor,
				ForceColors:               val.forceColor,
				EnvironmentOverrideColors: val.envColor,
			}
			cleanenv()
			if val.clicolorIsSet {
				os.Setenv("CLICOLOR", val.clicolorVal)
			}
			if val.clicolorForceIsSet {
				os.Setenv("CLICOLOR_FORCE", val.clicolorForceVal)
			}
			res := tf.isColored()
			if runtime.GOOS == "windows" && !tf.ForceColors && !val.clicolorForceIsSet {
				assert.Equal(subT, false, res)
			} else {
				assert.Equal(subT, val.expectedResult, res)
			}
		})
	}
}

func TestCustomSorting(t *testing.T) {
	formatter := &TextFormatter{
		DisableColors: true,
		SortingFunc: func(keys []string) {
			sort.Slice(keys, func(i, j int) bool {
				if keys[j] == "prefix" {
					return false
				}
				if keys[i] == "prefix" {
					return true
				}
				return strings.Compare(keys[i], keys[j]) == -1
			})
		},
	}

	entry := &Entry{
		Message: "Testing custom sort function",
		Time:    time.Now(),
		Level:   InfoLevel,
		Data: Fields{
			"test":      "testvalue",
			"prefix":    "the application prefix",
			"blablabla": "blablabla",
		},
	}
	b, err := formatter.Format(entry)
	require.NoError(t, err)
	require.True(t, strings.HasPrefix(string(b), "prefix="), "format output is %q", string(b))
}
