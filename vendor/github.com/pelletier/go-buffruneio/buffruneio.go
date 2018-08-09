// Package buffruneio is a wrapper around bufio to provide buffered runes access with unlimited unreads.
package buffruneio

import (
	"bufio"
	"container/list"
	"errors"
	"io"
)

// Rune to indicate end of file.
const (
	EOF = -(iota + 1)
)

// ErrNoRuneToUnread is returned by UnreadRune() when the read index is already at the beginning of the buffer.
var ErrNoRuneToUnread = errors.New("no rune to unwind")

// Reader implements runes buffering for an io.Reader object.
type Reader struct {
	buffer  *list.List
	current *list.Element
	input   *bufio.Reader
}

// NewReader returns a new Reader.
func NewReader(rd io.Reader) *Reader {
	return &Reader{
		buffer: list.New(),
		input:  bufio.NewReader(rd),
	}
}

type runeWithSize struct {
	r    rune
	size int
}

func (rd *Reader) feedBuffer() error {
	r, size, err := rd.input.ReadRune()

	if err != nil {
		if err != io.EOF {
			return err
		}
		r = EOF
	}

	newRuneWithSize := runeWithSize{r, size}

	rd.buffer.PushBack(newRuneWithSize)
	if rd.current == nil {
		rd.current = rd.buffer.Back()
	}
	return nil
}

// ReadRune reads the next rune from buffer, or from the underlying reader if needed.
func (rd *Reader) ReadRune() (rune, int, error) {
	if rd.current == rd.buffer.Back() || rd.current == nil {
		err := rd.feedBuffer()
		if err != nil {
			return EOF, 0, err
		}
	}

	runeWithSize := rd.current.Value.(runeWithSize)
	rd.current = rd.current.Next()
	return runeWithSize.r, runeWithSize.size, nil
}

// UnreadRune pushes back the previously read rune in the buffer, extending it if needed.
func (rd *Reader) UnreadRune() error {
	if rd.current == rd.buffer.Front() {
		return ErrNoRuneToUnread
	}
	if rd.current == nil {
		rd.current = rd.buffer.Back()
	} else {
		rd.current = rd.current.Prev()
	}
	return nil
}

// Forget removes runes stored before the current stream position index.
func (rd *Reader) Forget() {
	if rd.current == nil {
		rd.current = rd.buffer.Back()
	}
	for ; rd.current != rd.buffer.Front(); rd.buffer.Remove(rd.current.Prev()) {
	}
}

// PeekRune returns at most the next n runes, reading from the uderlying source if
// needed. Does not move the current index. It includes EOF if reached.
func (rd *Reader) PeekRunes(n int) []rune {
	res := make([]rune, 0, n)
	cursor := rd.current
	for i := 0; i < n; i++ {
		if cursor == nil {
			err := rd.feedBuffer()
			if err != nil {
				return res
			}
			cursor = rd.buffer.Back()
		}
		if cursor != nil {
			r := cursor.Value.(runeWithSize).r
			res = append(res, r)
			if r == EOF {
				return res
			}
			cursor = cursor.Next()
		}
	}
	return res
}
