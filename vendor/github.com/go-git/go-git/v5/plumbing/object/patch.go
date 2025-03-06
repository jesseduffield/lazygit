package object

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	fdiff "github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/utils/diff"

	dmp "github.com/sergi/go-diff/diffmatchpatch"
)

var (
	ErrCanceled = errors.New("operation canceled")
)

func getPatch(message string, changes ...*Change) (*Patch, error) {
	ctx := context.Background()
	return getPatchContext(ctx, message, changes...)
}

func getPatchContext(ctx context.Context, message string, changes ...*Change) (*Patch, error) {
	var filePatches []fdiff.FilePatch
	for _, c := range changes {
		select {
		case <-ctx.Done():
			return nil, ErrCanceled
		default:
		}

		fp, err := filePatchWithContext(ctx, c)
		if err != nil {
			return nil, err
		}

		filePatches = append(filePatches, fp)
	}

	return &Patch{message, filePatches}, nil
}

func filePatchWithContext(ctx context.Context, c *Change) (fdiff.FilePatch, error) {
	from, to, err := c.Files()
	if err != nil {
		return nil, err
	}
	fromContent, fIsBinary, err := fileContent(from)
	if err != nil {
		return nil, err
	}

	toContent, tIsBinary, err := fileContent(to)
	if err != nil {
		return nil, err
	}

	if fIsBinary || tIsBinary {
		return &textFilePatch{from: c.From, to: c.To}, nil
	}

	diffs := diff.Do(fromContent, toContent)

	var chunks []fdiff.Chunk
	for _, d := range diffs {
		select {
		case <-ctx.Done():
			return nil, ErrCanceled
		default:
		}

		var op fdiff.Operation
		switch d.Type {
		case dmp.DiffEqual:
			op = fdiff.Equal
		case dmp.DiffDelete:
			op = fdiff.Delete
		case dmp.DiffInsert:
			op = fdiff.Add
		}

		chunks = append(chunks, &textChunk{d.Text, op})
	}

	return &textFilePatch{
		chunks: chunks,
		from:   c.From,
		to:     c.To,
	}, nil

}

func fileContent(f *File) (content string, isBinary bool, err error) {
	if f == nil {
		return
	}

	isBinary, err = f.IsBinary()
	if err != nil || isBinary {
		return
	}

	content, err = f.Contents()

	return
}

// Patch is an implementation of fdiff.Patch interface
type Patch struct {
	message     string
	filePatches []fdiff.FilePatch
}

func (p *Patch) FilePatches() []fdiff.FilePatch {
	return p.filePatches
}

func (p *Patch) Message() string {
	return p.message
}

func (p *Patch) Encode(w io.Writer) error {
	ue := fdiff.NewUnifiedEncoder(w, fdiff.DefaultContextLines)

	return ue.Encode(p)
}

func (p *Patch) Stats() FileStats {
	return getFileStatsFromFilePatches(p.FilePatches())
}

func (p *Patch) String() string {
	buf := bytes.NewBuffer(nil)
	err := p.Encode(buf)
	if err != nil {
		return fmt.Sprintf("malformed patch: %s", err.Error())
	}

	return buf.String()
}

// changeEntryWrapper is an implementation of fdiff.File interface
type changeEntryWrapper struct {
	ce ChangeEntry
}

func (f *changeEntryWrapper) Hash() plumbing.Hash {
	if !f.ce.TreeEntry.Mode.IsFile() {
		return plumbing.ZeroHash
	}

	return f.ce.TreeEntry.Hash
}

func (f *changeEntryWrapper) Mode() filemode.FileMode {
	return f.ce.TreeEntry.Mode
}
func (f *changeEntryWrapper) Path() string {
	if !f.ce.TreeEntry.Mode.IsFile() {
		return ""
	}

	return f.ce.Name
}

func (f *changeEntryWrapper) Empty() bool {
	return !f.ce.TreeEntry.Mode.IsFile()
}

// textFilePatch is an implementation of fdiff.FilePatch interface
type textFilePatch struct {
	chunks   []fdiff.Chunk
	from, to ChangeEntry
}

func (tf *textFilePatch) Files() (from fdiff.File, to fdiff.File) {
	f := &changeEntryWrapper{tf.from}
	t := &changeEntryWrapper{tf.to}

	if !f.Empty() {
		from = f
	}

	if !t.Empty() {
		to = t
	}

	return
}

func (tf *textFilePatch) IsBinary() bool {
	return len(tf.chunks) == 0
}

func (tf *textFilePatch) Chunks() []fdiff.Chunk {
	return tf.chunks
}

// textChunk is an implementation of fdiff.Chunk interface
type textChunk struct {
	content string
	op      fdiff.Operation
}

func (t *textChunk) Content() string {
	return t.content
}

func (t *textChunk) Type() fdiff.Operation {
	return t.op
}

// FileStat stores the status of changes in content of a file.
type FileStat struct {
	Name     string
	Addition int
	Deletion int
}

func (fs FileStat) String() string {
	return printStat([]FileStat{fs})
}

// FileStats is a collection of FileStat.
type FileStats []FileStat

func (fileStats FileStats) String() string {
	return printStat(fileStats)
}

// printStat prints the stats of changes in content of files.
// Original implementation: https://github.com/git/git/blob/1a87c842ece327d03d08096395969aca5e0a6996/diff.c#L2615
// Parts of the output:
// <pad><filename><pad>|<pad><changeNumber><pad><+++/---><newline>
// example: " main.go | 10 +++++++--- "
func printStat(fileStats []FileStat) string {
	maxGraphWidth := uint(53)
	maxNameLen := 0
	maxChangeLen := 0

	scaleLinear := func(it, width, max uint) uint {
		if it == 0 || max == 0 {
			return 0
		}

		return 1 + (it * (width - 1) / max)
	}

	for _, fs := range fileStats {
		if len(fs.Name) > maxNameLen {
			maxNameLen = len(fs.Name)
		}

		changes := strconv.Itoa(fs.Addition + fs.Deletion)
		if len(changes) > maxChangeLen {
			maxChangeLen = len(changes)
		}
	}

	result := ""
	for _, fs := range fileStats {
		add := uint(fs.Addition)
		del := uint(fs.Deletion)
		np := maxNameLen - len(fs.Name)
		cp := maxChangeLen - len(strconv.Itoa(fs.Addition+fs.Deletion))

		total := add + del
		if total > maxGraphWidth {
			add = scaleLinear(add, maxGraphWidth, total)
			del = scaleLinear(del, maxGraphWidth, total)
		}

		adds := strings.Repeat("+", int(add))
		dels := strings.Repeat("-", int(del))
		namePad := strings.Repeat(" ", np)
		changePad := strings.Repeat(" ", cp)

		result += fmt.Sprintf(" %s%s | %s%d %s%s\n", fs.Name, namePad, changePad, total, adds, dels)
	}
	return result
}

func getFileStatsFromFilePatches(filePatches []fdiff.FilePatch) FileStats {
	var fileStats FileStats

	for _, fp := range filePatches {
		// ignore empty patches (binary files, submodule refs updates)
		if len(fp.Chunks()) == 0 {
			continue
		}

		cs := FileStat{}
		from, to := fp.Files()
		if from == nil {
			// New File is created.
			cs.Name = to.Path()
		} else if to == nil {
			// File is deleted.
			cs.Name = from.Path()
		} else if from.Path() != to.Path() {
			// File is renamed.
			cs.Name = fmt.Sprintf("%s => %s", from.Path(), to.Path())
		} else {
			cs.Name = from.Path()
		}

		for _, chunk := range fp.Chunks() {
			s := chunk.Content()
			if len(s) == 0 {
				continue
			}

			switch chunk.Type() {
			case fdiff.Add:
				cs.Addition += strings.Count(s, "\n")
				if s[len(s)-1] != '\n' {
					cs.Addition++
				}
			case fdiff.Delete:
				cs.Deletion += strings.Count(s, "\n")
				if s[len(s)-1] != '\n' {
					cs.Deletion++
				}
			}
		}

		fileStats = append(fileStats, cs)
	}

	return fileStats
}
