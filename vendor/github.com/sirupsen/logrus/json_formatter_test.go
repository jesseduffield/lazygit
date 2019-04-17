package logrus

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestErrorNotLost(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("error", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["error"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("omg", errors.New("wild walrus")))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("time", "right now!"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.time"] != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}

	if entry["time"] != "0001-01-01T00:00:00Z" {
		t.Fatal("time field not set to current time, was: ", entry["time"])
	}
}

func TestFieldClashWithMsg(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("msg", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.msg"] != "something" {
		t.Fatal("fields.msg not set to original msg field")
	}
}

func TestFieldClashWithLevel(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.level"] != "something" {
		t.Fatal("fields.level not set to original level field")
	}
}

func TestFieldClashWithRemappedFields(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime:  "@timestamp",
			FieldKeyLevel: "@level",
			FieldKeyMsg:   "@message",
		},
	}

	b, err := formatter.Format(WithFields(Fields{
		"@timestamp": "@timestamp",
		"@level":     "@level",
		"@message":   "@message",
		"timestamp":  "timestamp",
		"level":      "level",
		"msg":        "msg",
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	for _, field := range []string{"timestamp", "level", "msg"} {
		if entry[field] != field {
			t.Errorf("Expected field %v to be untouched; got %v", field, entry[field])
		}

		remappedKey := fmt.Sprintf("fields.%s", field)
		if remapped, ok := entry[remappedKey]; ok {
			t.Errorf("Expected %s to be empty; got %v", remappedKey, remapped)
		}
	}

	for _, field := range []string{"@timestamp", "@level", "@message"} {
		if entry[field] == field {
			t.Errorf("Expected field %v to be mapped to an Entry value", field)
		}

		remappedKey := fmt.Sprintf("fields.%s", field)
		if remapped, ok := entry[remappedKey]; ok {
			if remapped != field {
				t.Errorf("Expected field %v to be copied to %s; got %v", field, remappedKey, remapped)
			}
		} else {
			t.Errorf("Expected field %v to be copied to %s; was absent", field, remappedKey)
		}
	}
}

func TestFieldsInNestedDictionary(t *testing.T) {
	formatter := &JSONFormatter{
		DataKey: "args",
	}

	logEntry := WithFields(Fields{
		"level": "level",
		"test":  "test",
	})
	logEntry.Level = InfoLevel

	b, err := formatter.Format(logEntry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	args := entry["args"].(map[string]interface{})

	for _, field := range []string{"test", "level"} {
		if value, present := args[field]; !present || value != field {
			t.Errorf("Expected field %v to be present under 'args'; untouched", field)
		}
	}

	for _, field := range []string{"test", "fields.level"} {
		if _, present := entry[field]; present {
			t.Errorf("Expected field %v not to be present at top level", field)
		}
	}

	// with nested object, "level" shouldn't clash
	if entry["level"] != "info" {
		t.Errorf("Expected 'level' field to contain 'info'")
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}

func TestJSONMessageKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyMsg: "message",
		},
	}

	b, err := formatter.Format(&Entry{Message: "oh hai"})
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !(strings.Contains(s, "message") && strings.Contains(s, "oh hai")) {
		t.Fatal("Expected JSON to format message key")
	}
}

func TestJSONLevelKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyLevel: "somelevel",
		},
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "somelevel") {
		t.Fatal("Expected JSON to format level key")
	}
}

func TestJSONTimeKey(t *testing.T) {
	formatter := &JSONFormatter{
		FieldMap: FieldMap{
			FieldKeyTime: "timeywimey",
		},
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, "timeywimey") {
		t.Fatal("Expected JSON to format time key")
	}
}

func TestFieldDoesNotClashWithCaller(t *testing.T) {
	SetReportCaller(false)
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("func", "howdy pardner"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["func"] != "howdy pardner" {
		t.Fatal("func field replaced when ReportCaller=false")
	}
}

func TestFieldClashWithCaller(t *testing.T) {
	SetReportCaller(true)
	formatter := &JSONFormatter{}
	e := WithField("func", "howdy pardner")
	e.Caller = &runtime.Frame{Function: "somefunc"}
	b, err := formatter.Format(e)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.func"] != "howdy pardner" {
		t.Fatalf("fields.func not set to original func field when ReportCaller=true (got '%s')",
			entry["fields.func"])
	}

	if entry["func"] != "somefunc" {
		t.Fatalf("func not set as expected when ReportCaller=true (got '%s')",
			entry["func"])
	}

	SetReportCaller(false) // return to default value
}

func TestJSONDisableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{
		DisableTimestamp: true,
	}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if strings.Contains(s, FieldKeyTime) {
		t.Error("Did not prevent timestamp", s)
	}
}

func TestJSONEnableTimestamp(t *testing.T) {
	formatter := &JSONFormatter{}

	b, err := formatter.Format(WithField("level", "something"))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}
	s := string(b)
	if !strings.Contains(s, FieldKeyTime) {
		t.Error("Timestamp not present", s)
	}
}
