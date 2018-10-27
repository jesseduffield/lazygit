/*
Copyright 2012 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shlex

import (
	"strings"
	"testing"
)

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func TestClassifier(t *testing.T) {
	classifier := NewDefaultClassifier()
	runeTests := map[int32]RuneTokenType{
		'a':  RUNETOKEN_CHAR,
		' ':  RUNETOKEN_SPACE,
		'"':  RUNETOKEN_ESCAPING_QUOTE,
		'\'': RUNETOKEN_NONESCAPING_QUOTE,
		'#':  RUNETOKEN_COMMENT}
	for rune, expectedType := range runeTests {
		foundType := classifier.ClassifyRune(rune)
		if foundType != expectedType {
			t.Logf("Expected type: %v for rune '%c'(%v). Found type: %v.", expectedType, rune, rune, foundType)
			t.Fail()
		}
	}
}

func TestTokenizer(t *testing.T) {
	testInput := strings.NewReader("one two \"three four\" \"five \\\"six\\\"\" seven#eight # nine # ten\n eleven")
	expectedTokens := []*Token{
		&Token{
			tokenType: TOKEN_WORD,
			value:     "one"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "two"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "three four"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "five \"six\""},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "seven#eight"},
		&Token{
			tokenType: TOKEN_COMMENT,
			value:     " nine # ten"},
		&Token{
			tokenType: TOKEN_WORD,
			value:     "eleven"}}

	tokenizer, err := NewTokenizer(testInput)
	checkError(err, t)
	for _, expectedToken := range expectedTokens {
		foundToken, err := tokenizer.NextToken()
		checkError(err, t)
		if !foundToken.Equal(expectedToken) {
			t.Error("Expected token:", expectedToken, ". Found:", foundToken)
		}
	}
}

func TestLexer(t *testing.T) {
	testInput := strings.NewReader("one")
	expectedWord := "one"
	lexer, err := NewLexer(testInput)
	checkError(err, t)
	foundWord, err := lexer.NextWord()
	checkError(err, t)
	if expectedWord != foundWord {
		t.Error("Expected word:", expectedWord, ". Found:", foundWord)
	}
}

func TestSplitSimple(t *testing.T) {
	testInput := "one two three"
	expectedOutput := []string{"one", "two", "three"}
	foundOutput, err := Split(testInput)
	if err != nil {
		t.Error("Split returned error:", err)
	}
	if len(expectedOutput) != len(foundOutput) {
		t.Error("Split expected:", len(expectedOutput), "results. Found:", len(foundOutput), "results")
	}
	for i := range foundOutput {
		if foundOutput[i] != expectedOutput[i] {
			t.Error("Item:", i, "(", foundOutput[i], ") differs from the expected value:", expectedOutput[i])
		}
	}
}

func TestSplitEscapingQuotes(t *testing.T) {
	testInput := "one \"два ${three}\" four"
	expectedOutput := []string{"one", "два ${three}", "four"}
	foundOutput, err := Split(testInput)
	if err != nil {
		t.Error("Split returned error:", err)
	}
	if len(expectedOutput) != len(foundOutput) {
		t.Error("Split expected:", len(expectedOutput), "results. Found:", len(foundOutput), "results")
	}
	for i := range foundOutput {
		if foundOutput[i] != expectedOutput[i] {
			t.Error("Item:", i, "(", foundOutput[i], ") differs from the expected value:", expectedOutput[i])
		}
	}
}

func TestGlobbingExpressions(t *testing.T) {
	testInput := "onefile *file one?ile onefil[de]"
	expectedOutput := []string{"onefile", "*file", "one?ile", "onefil[de]"}
	foundOutput, err := Split(testInput)
	if err != nil {
		t.Error("Split returned error", err)
	}
	if len(expectedOutput) != len(foundOutput) {
		t.Error("Split expected:", len(expectedOutput), "results. Found:", len(foundOutput), "results")
	}
	for i := range foundOutput {
		if foundOutput[i] != expectedOutput[i] {
			t.Error("Item:", i, "(", foundOutput[i], ") differs from the expected value:", expectedOutput[i])
		}
	}

}

func TestSplitNonEscapingQuotes(t *testing.T) {
	testInput := "one 'два ${three}' four"
	expectedOutput := []string{"one", "два ${three}", "four"}
	foundOutput, err := Split(testInput)
	if err != nil {
		t.Error("Split returned error:", err)
	}
	if len(expectedOutput) != len(foundOutput) {
		t.Error("Split expected:", len(expectedOutput), "results. Found:", len(foundOutput), "results")
	}
	for i := range foundOutput {
		if foundOutput[i] != expectedOutput[i] {
			t.Error("Item:", i, "(", foundOutput[i], ") differs from the expected value:", expectedOutput[i])
		}
	}
}
