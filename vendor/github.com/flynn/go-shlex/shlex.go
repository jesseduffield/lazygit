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

/*
Package shlex implements a simple lexer which splits input in to tokens using
shell-style rules for quoting and commenting.
*/
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

/*
A TokenType is a top-level token; a word, space, comment, unknown.
*/
type TokenType int

/*
A RuneTokenType is the type of a UTF-8 character; a character, quote, space, escape.
*/
type RuneTokenType int

type lexerState int

type Token struct {
	tokenType TokenType
	value     string
}

/*
Two tokens are equal if both their types and values are equal. A nil token can
never equal another token.
*/
func (a *Token) Equal(b *Token) bool {
	if a == nil || b == nil {
		return false
	}
	if a.tokenType != b.tokenType {
		return false
	}
	return a.value == b.value
}

const (
	RUNE_CHAR              string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._-,/@$*()+=><:;&^%~|!?[]{}"
	RUNE_SPACE             string = " \t\r\n"
	RUNE_ESCAPING_QUOTE    string = "\""
	RUNE_NONESCAPING_QUOTE string = "'"
	RUNE_ESCAPE                   = "\\"
	RUNE_COMMENT                  = "#"

	RUNETOKEN_UNKNOWN           RuneTokenType = 0
	RUNETOKEN_CHAR              RuneTokenType = 1
	RUNETOKEN_SPACE             RuneTokenType = 2
	RUNETOKEN_ESCAPING_QUOTE    RuneTokenType = 3
	RUNETOKEN_NONESCAPING_QUOTE RuneTokenType = 4
	RUNETOKEN_ESCAPE            RuneTokenType = 5
	RUNETOKEN_COMMENT           RuneTokenType = 6
	RUNETOKEN_EOF               RuneTokenType = 7

	TOKEN_UNKNOWN TokenType = 0
	TOKEN_WORD    TokenType = 1
	TOKEN_SPACE   TokenType = 2
	TOKEN_COMMENT TokenType = 3

	STATE_START           lexerState = 0
	STATE_INWORD          lexerState = 1
	STATE_ESCAPING        lexerState = 2
	STATE_ESCAPING_QUOTED lexerState = 3
	STATE_QUOTED_ESCAPING lexerState = 4
	STATE_QUOTED          lexerState = 5
	STATE_COMMENT         lexerState = 6

	INITIAL_TOKEN_CAPACITY int = 100
)

/*
A type for classifying characters. This allows for different sorts of
classifiers - those accepting extended non-ascii chars, or strict posix
compatibility, for example.
*/
type TokenClassifier struct {
	typeMap map[int32]RuneTokenType
}

func addRuneClass(typeMap *map[int32]RuneTokenType, runes string, tokenType RuneTokenType) {
	for _, rune := range runes {
		(*typeMap)[int32(rune)] = tokenType
	}
}

/*
Create a new classifier for basic ASCII characters.
*/
func NewDefaultClassifier() *TokenClassifier {
	typeMap := map[int32]RuneTokenType{}
	addRuneClass(&typeMap, RUNE_CHAR, RUNETOKEN_CHAR)
	addRuneClass(&typeMap, RUNE_SPACE, RUNETOKEN_SPACE)
	addRuneClass(&typeMap, RUNE_ESCAPING_QUOTE, RUNETOKEN_ESCAPING_QUOTE)
	addRuneClass(&typeMap, RUNE_NONESCAPING_QUOTE, RUNETOKEN_NONESCAPING_QUOTE)
	addRuneClass(&typeMap, RUNE_ESCAPE, RUNETOKEN_ESCAPE)
	addRuneClass(&typeMap, RUNE_COMMENT, RUNETOKEN_COMMENT)
	return &TokenClassifier{
		typeMap: typeMap}
}

func (classifier *TokenClassifier) ClassifyRune(rune int32) RuneTokenType {
	return classifier.typeMap[rune]
}

/*
A type for turning an input stream in to a sequence of strings. Whitespace and
comments are skipped.
*/
type Lexer struct {
	tokenizer *Tokenizer
}

/*
Create a new lexer.
*/
func NewLexer(r io.Reader) (*Lexer, error) {

	tokenizer, err := NewTokenizer(r)
	if err != nil {
		return nil, err
	}
	lexer := &Lexer{tokenizer: tokenizer}
	return lexer, nil
}

/*
Return the next word, and an error value. If there are no more words, the error
will be io.EOF.
*/
func (l *Lexer) NextWord() (string, error) {
	var token *Token
	var err error
	for {
		token, err = l.tokenizer.NextToken()
		if err != nil {
			return "", err
		}
		switch token.tokenType {
		case TOKEN_WORD:
			{
				return token.value, nil
			}
		case TOKEN_COMMENT:
			{
				// skip comments
			}
		default:
			{
				panic(fmt.Sprintf("Unknown token type: %v", token.tokenType))
			}
		}
	}
	return "", io.EOF
}

/*
A type for turning an input stream in to a sequence of typed tokens.
*/
type Tokenizer struct {
	input      *bufio.Reader
	classifier *TokenClassifier
}

/*
Create a new tokenizer.
*/
func NewTokenizer(r io.Reader) (*Tokenizer, error) {
	input := bufio.NewReader(r)
	classifier := NewDefaultClassifier()
	tokenizer := &Tokenizer{
		input:      input,
		classifier: classifier}
	return tokenizer, nil
}

/*
Scan the stream for the next token.

This uses an internal state machine. It will panic if it encounters a character
which it does not know how to handle.
*/
func (t *Tokenizer) scanStream() (*Token, error) {
	state := STATE_START
	var tokenType TokenType
	value := make([]int32, 0, INITIAL_TOKEN_CAPACITY)
	var (
		nextRune     int32
		nextRuneType RuneTokenType
		err          error
	)
SCAN:
	for {
		nextRune, _, err = t.input.ReadRune()
		nextRuneType = t.classifier.ClassifyRune(nextRune)
		if err != nil {
			if err == io.EOF {
				nextRuneType = RUNETOKEN_EOF
				err = nil
			} else {
				return nil, err
			}
		}
		switch state {
		case STATE_START: // no runes read yet
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						return nil, io.EOF
					}
				case RUNETOKEN_CHAR:
					{
						tokenType = TOKEN_WORD
						value = append(value, nextRune)
						state = STATE_INWORD
					}
				case RUNETOKEN_SPACE:
					{
					}
				case RUNETOKEN_ESCAPING_QUOTE:
					{
						tokenType = TOKEN_WORD
						state = STATE_QUOTED_ESCAPING
					}
				case RUNETOKEN_NONESCAPING_QUOTE:
					{
						tokenType = TOKEN_WORD
						state = STATE_QUOTED
					}
				case RUNETOKEN_ESCAPE:
					{
						tokenType = TOKEN_WORD
						state = STATE_ESCAPING
					}
				case RUNETOKEN_COMMENT:
					{
						tokenType = TOKEN_COMMENT
						state = STATE_COMMENT
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Unknown rune: %v", nextRune))
					}
				}
			}
		case STATE_INWORD: // in a regular word
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						break SCAN
					}
				case RUNETOKEN_CHAR, RUNETOKEN_COMMENT:
					{
						value = append(value, nextRune)
					}
				case RUNETOKEN_SPACE:
					{
						t.input.UnreadRune()
						break SCAN
					}
				case RUNETOKEN_ESCAPING_QUOTE:
					{
						state = STATE_QUOTED_ESCAPING
					}
				case RUNETOKEN_NONESCAPING_QUOTE:
					{
						state = STATE_QUOTED
					}
				case RUNETOKEN_ESCAPE:
					{
						state = STATE_ESCAPING
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_ESCAPING: // the next rune after an escape character
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						err = errors.New("EOF found after escape character")
						break SCAN
					}
				case RUNETOKEN_CHAR, RUNETOKEN_SPACE, RUNETOKEN_ESCAPING_QUOTE, RUNETOKEN_NONESCAPING_QUOTE, RUNETOKEN_ESCAPE, RUNETOKEN_COMMENT:
					{
						state = STATE_INWORD
						value = append(value, nextRune)
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_ESCAPING_QUOTED: // the next rune after an escape character, in double quotes
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						err = errors.New("EOF found after escape character")
						break SCAN
					}
				case RUNETOKEN_CHAR, RUNETOKEN_SPACE, RUNETOKEN_ESCAPING_QUOTE, RUNETOKEN_NONESCAPING_QUOTE, RUNETOKEN_ESCAPE, RUNETOKEN_COMMENT:
					{
						state = STATE_QUOTED_ESCAPING
						value = append(value, nextRune)
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_QUOTED_ESCAPING: // in escaping double quotes
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						err = errors.New("EOF found when expecting closing quote.")
						break SCAN
					}
				case RUNETOKEN_CHAR, RUNETOKEN_UNKNOWN, RUNETOKEN_SPACE, RUNETOKEN_NONESCAPING_QUOTE, RUNETOKEN_COMMENT:
					{
						value = append(value, nextRune)
					}
				case RUNETOKEN_ESCAPING_QUOTE:
					{
						state = STATE_INWORD
					}
				case RUNETOKEN_ESCAPE:
					{
						state = STATE_ESCAPING_QUOTED
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_QUOTED: // in non-escaping single quotes
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						err = errors.New("EOF found when expecting closing quote.")
						break SCAN
					}
				case RUNETOKEN_CHAR, RUNETOKEN_UNKNOWN, RUNETOKEN_SPACE, RUNETOKEN_ESCAPING_QUOTE, RUNETOKEN_ESCAPE, RUNETOKEN_COMMENT:
					{
						value = append(value, nextRune)
					}
				case RUNETOKEN_NONESCAPING_QUOTE:
					{
						state = STATE_INWORD
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		case STATE_COMMENT:
			{
				switch nextRuneType {
				case RUNETOKEN_EOF:
					{
						break SCAN
					}
				case RUNETOKEN_CHAR, RUNETOKEN_UNKNOWN, RUNETOKEN_ESCAPING_QUOTE, RUNETOKEN_ESCAPE, RUNETOKEN_COMMENT, RUNETOKEN_NONESCAPING_QUOTE:
					{
						value = append(value, nextRune)
					}
				case RUNETOKEN_SPACE:
					{
						if nextRune == '\n' {
							state = STATE_START
							break SCAN
						} else {
							value = append(value, nextRune)
						}
					}
				default:
					{
						return nil, errors.New(fmt.Sprintf("Uknown rune: %v", nextRune))
					}
				}
			}
		default:
			{
				panic(fmt.Sprintf("Unexpected state: %v", state))
			}
		}
	}
	token := &Token{
		tokenType: tokenType,
		value:     string(value)}
	return token, err
}

/*
Return the next token in the stream, and an error value. If there are no more
tokens available, the error value will be io.EOF.
*/
func (t *Tokenizer) NextToken() (*Token, error) {
	return t.scanStream()
}

/*
Split a string in to a slice of strings, based upon shell-style rules for
quoting, escaping, and spaces.
*/
func Split(s string) ([]string, error) {
	l, err := NewLexer(strings.NewReader(s))
	if err != nil {
		return nil, err
	}
	subStrings := []string{}
	for {
		word, err := l.NextWord()
		if err != nil {
			if err == io.EOF {
				return subStrings, nil
			}
			return subStrings, err
		}
		subStrings = append(subStrings, word)
	}
	return subStrings, nil
}
