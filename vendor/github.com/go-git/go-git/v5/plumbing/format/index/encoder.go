package index

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/plumbing/hash"
	"github.com/go-git/go-git/v5/utils/binary"
)

var (
	// EncodeVersionSupported is the range of supported index versions
	EncodeVersionSupported uint32 = 4

	// ErrInvalidTimestamp is returned by Encode if a Index with a Entry with
	// negative timestamp values
	ErrInvalidTimestamp = errors.New("negative timestamps are not allowed")
)

// An Encoder writes an Index to an output stream.
type Encoder struct {
	w         io.Writer
	hash      hash.Hash
	lastEntry *Entry
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	h := hash.New(hash.CryptoType)
	mw := io.MultiWriter(w, h)
	return &Encoder{mw, h, nil}
}

// Encode writes the Index to the stream of the encoder.
func (e *Encoder) Encode(idx *Index) error {
	return e.encode(idx, true)
}

func (e *Encoder) encode(idx *Index, footer bool) error {

	// TODO: support extensions
	if idx.Version > EncodeVersionSupported {
		return ErrUnsupportedVersion
	}

	if err := e.encodeHeader(idx); err != nil {
		return err
	}

	if err := e.encodeEntries(idx); err != nil {
		return err
	}

	if footer {
		return e.encodeFooter()
	}
	return nil
}

func (e *Encoder) encodeHeader(idx *Index) error {
	return binary.Write(e.w,
		indexSignature,
		idx.Version,
		uint32(len(idx.Entries)),
	)
}

func (e *Encoder) encodeEntries(idx *Index) error {
	sort.Sort(byName(idx.Entries))

	for _, entry := range idx.Entries {
		if err := e.encodeEntry(idx, entry); err != nil {
			return err
		}
		entryLength := entryHeaderLength
		if entry.IntentToAdd || entry.SkipWorktree {
			entryLength += 2
		}

		wrote := entryLength + len(entry.Name)
		if err := e.padEntry(idx, wrote); err != nil {
			return err
		}
	}

	return nil
}

func (e *Encoder) encodeEntry(idx *Index, entry *Entry) error {
	sec, nsec, err := e.timeToUint32(&entry.CreatedAt)
	if err != nil {
		return err
	}

	msec, mnsec, err := e.timeToUint32(&entry.ModifiedAt)
	if err != nil {
		return err
	}

	flags := uint16(entry.Stage&0x3) << 12
	if l := len(entry.Name); l < nameMask {
		flags |= uint16(l)
	} else {
		flags |= nameMask
	}

	flow := []interface{}{
		sec, nsec,
		msec, mnsec,
		entry.Dev,
		entry.Inode,
		entry.Mode,
		entry.UID,
		entry.GID,
		entry.Size,
		entry.Hash[:],
	}

	flagsFlow := []interface{}{flags}

	if entry.IntentToAdd || entry.SkipWorktree {
		var extendedFlags uint16

		if entry.IntentToAdd {
			extendedFlags |= intentToAddMask
		}
		if entry.SkipWorktree {
			extendedFlags |= skipWorkTreeMask
		}

		flagsFlow = []interface{}{flags | entryExtended, extendedFlags}
	}

	flow = append(flow, flagsFlow...)

	if err := binary.Write(e.w, flow...); err != nil {
		return err
	}

	switch idx.Version {
	case 2, 3:
		err = e.encodeEntryName(entry)
	case 4:
		err = e.encodeEntryNameV4(entry)
	default:
		err = ErrUnsupportedVersion
	}

	return err
}

func (e *Encoder) encodeEntryName(entry *Entry) error {
	return binary.Write(e.w, []byte(entry.Name))
}

func (e *Encoder) encodeEntryNameV4(entry *Entry) error {
	name := entry.Name
	l := 0
	if e.lastEntry != nil {
		dir := path.Dir(e.lastEntry.Name) + "/"
		if strings.HasPrefix(entry.Name, dir) {
			l = len(e.lastEntry.Name) - len(dir)
			name = strings.TrimPrefix(entry.Name, dir)
		} else {
			l = len(e.lastEntry.Name)
		}
	}

	e.lastEntry = entry

	err := binary.WriteVariableWidthInt(e.w, int64(l))
	if err != nil {
		return err
	}

	return binary.Write(e.w, []byte(name+string('\x00')))
}

func (e *Encoder) encodeRawExtension(signature string, data []byte) error {
	if len(signature) != 4 {
		return fmt.Errorf("invalid signature length")
	}

	_, err := e.w.Write([]byte(signature))
	if err != nil {
		return err
	}

	err = binary.WriteUint32(e.w, uint32(len(data)))
	if err != nil {
		return err
	}

	_, err = e.w.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) timeToUint32(t *time.Time) (uint32, uint32, error) {
	if t.IsZero() {
		return 0, 0, nil
	}

	if t.Unix() < 0 || t.UnixNano() < 0 {
		return 0, 0, ErrInvalidTimestamp
	}

	return uint32(t.Unix()), uint32(t.Nanosecond()), nil
}

func (e *Encoder) padEntry(idx *Index, wrote int) error {
	if idx.Version == 4 {
		return nil
	}

	padLen := 8 - wrote%8

	_, err := e.w.Write(bytes.Repeat([]byte{'\x00'}, padLen))
	return err
}

func (e *Encoder) encodeFooter() error {
	return binary.Write(e.w, e.hash.Sum(nil))
}

type byName []*Entry

func (l byName) Len() int           { return len(l) }
func (l byName) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byName) Less(i, j int) bool { return l[i].Name < l[j].Name }
