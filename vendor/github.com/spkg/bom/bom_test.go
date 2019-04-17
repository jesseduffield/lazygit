package bom_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spkg/bom"
	"github.com/stretchr/testify/assert"
)

var testCases = []struct {
	Input    []byte
	Expected []byte
}{
	{
		Input:    nil,
		Expected: nil,
	},
	{
		Input:    []byte{},
		Expected: []byte{},
	},
	{
		Input:    []byte{0xef},
		Expected: []byte{0xef},
	},
	{
		Input:    []byte{0xef, 0xbb},
		Expected: []byte{0xef, 0xbb},
	},
	{
		Input:    []byte{0xef, 0xbb, 0xbf},
		Expected: []byte{},
	},
	{
		Input:    []byte{0xef, 0xbb, 0xbf, 0x41, 0x42, 0x43},
		Expected: []byte{0x41, 0x42, 0x43},
	},
	{
		Input:    []byte{0xef, 0xbb, 0x41, 0x42, 0x43},
		Expected: []byte{0xef, 0xbb, 0x41, 0x42, 0x43},
	},
	{
		Input:    []byte{0xef, 0x41, 0x42, 0x43},
		Expected: []byte{0xef, 0x41, 0x42, 0x43},
	},
	{
		Input:    []byte{0x41, 0x42, 0x43},
		Expected: []byte{0x41, 0x42, 0x43},
	},
}

func TestClean(t *testing.T) {
	assert := assert.New(t)

	for _, tc := range testCases {
		output := bom.Clean(tc.Input)
		assert.Equal(tc.Expected, output)
	}
}

func TestReader(t *testing.T) {
	assert := assert.New(t)

	for _, tc := range testCases {
		// An input value of nil works differently to the Clean function.
		// In this case it results in an empty buffer, not nil.
		expected := tc.Expected
		if tc.Input == nil {
			expected = []byte{}
		}
		r1 := bytes.NewReader(tc.Input)
		r2 := bom.NewReader(r1)
		output, err := ioutil.ReadAll(r2)
		assert.NoError(err)
		assert.Equal(expected, output)
	}
}
