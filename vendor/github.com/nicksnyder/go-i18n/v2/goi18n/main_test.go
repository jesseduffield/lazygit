package main

import (
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	testCases := []struct {
		args     []string
		exitCode int
	}{
		{
			args:     []string{"-help"},
			exitCode: 2,
		},
		{
			args:     []string{"extract"},
			exitCode: 0,
		},
		{
			args:     []string{"merge"},
			exitCode: 1,
		},
	}
	for _, testCase := range testCases {
		t.Run(strings.Join(testCase.args, " "), func(t *testing.T) {
			if code := testableMain(testCase.args); code != testCase.exitCode {
				t.Fatalf("expected exit code %d; got %d", testCase.exitCode, code)
			}
		})
	}
}
