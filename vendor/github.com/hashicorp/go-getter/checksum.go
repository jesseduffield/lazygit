package getter

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	urlhelper "github.com/hashicorp/go-getter/helper/url"
)

// fileChecksum helps verifying the checksum for a file.
type fileChecksum struct {
	Type     string
	Hash     hash.Hash
	Value    []byte
	Filename string
}

// A ChecksumError is returned when a checksum differs
type ChecksumError struct {
	Hash     hash.Hash
	Actual   []byte
	Expected []byte
	File     string
}

func (cerr *ChecksumError) Error() string {
	if cerr == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"Checksums did not match for %s.\nExpected: %s\nGot: %s\n%T",
		cerr.File,
		hex.EncodeToString(cerr.Expected),
		hex.EncodeToString(cerr.Actual),
		cerr.Hash, // ex: *sha256.digest
	)
}

// checksum is a simple method to compute the checksum of a source file
// and compare it to the given expected value.
func (c *fileChecksum) checksum(source string) error {
	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("Failed to open file for checksum: %s", err)
	}
	defer f.Close()

	c.Hash.Reset()
	if _, err := io.Copy(c.Hash, f); err != nil {
		return fmt.Errorf("Failed to hash: %s", err)
	}

	if actual := c.Hash.Sum(nil); !bytes.Equal(actual, c.Value) {
		return &ChecksumError{
			Hash:     c.Hash,
			Actual:   actual,
			Expected: c.Value,
			File:     source,
		}
	}

	return nil
}

// extractChecksum will return a fileChecksum based on the 'checksum'
// parameter of u.
// ex:
//  http://hashicorp.com/terraform?checksum=<checksumValue>
//  http://hashicorp.com/terraform?checksum=<checksumType>:<checksumValue>
//  http://hashicorp.com/terraform?checksum=file:<checksum_url>
// when checksumming from a file, extractChecksum will go get checksum_url
// in a temporary directory, parse the content of the file then delete it.
// Content of files are expected to be BSD style or GNU style.
//
// BSD-style checksum:
//  MD5 (file1) = <checksum>
//  MD5 (file2) = <checksum>
//
// GNU-style:
//  <checksum>  file1
//  <checksum> *file2
//
// see parseChecksumLine for more detail on checksum file parsing
func (c *Client) extractChecksum(u *url.URL) (*fileChecksum, error) {
	q := u.Query()
	v := q.Get("checksum")

	if v == "" {
		return nil, nil
	}

	vs := strings.SplitN(v, ":", 2)
	switch len(vs) {
	case 2:
		break // good
	default:
		// here, we try to guess the checksum from it's length
		// if the type was not passed
		return newChecksumFromValue(v, filepath.Base(u.EscapedPath()))
	}

	checksumType, checksumValue := vs[0], vs[1]

	switch checksumType {
	case "file":
		return c.checksumFromFile(checksumValue, u)
	default:
		return newChecksumFromType(checksumType, checksumValue, filepath.Base(u.EscapedPath()))
	}
}

func newChecksum(checksumValue, filename string) (*fileChecksum, error) {
	c := &fileChecksum{
		Filename: filename,
	}
	var err error
	c.Value, err = hex.DecodeString(checksumValue)
	if err != nil {
		return nil, fmt.Errorf("invalid checksum: %s", err)
	}
	return c, nil
}

func newChecksumFromType(checksumType, checksumValue, filename string) (*fileChecksum, error) {
	c, err := newChecksum(checksumValue, filename)
	if err != nil {
		return nil, err
	}

	c.Type = strings.ToLower(checksumType)
	switch c.Type {
	case "md5":
		c.Hash = md5.New()
	case "sha1":
		c.Hash = sha1.New()
	case "sha256":
		c.Hash = sha256.New()
	case "sha512":
		c.Hash = sha512.New()
	default:
		return nil, fmt.Errorf(
			"unsupported checksum type: %s", checksumType)
	}

	return c, nil
}

func newChecksumFromValue(checksumValue, filename string) (*fileChecksum, error) {
	c, err := newChecksum(checksumValue, filename)
	if err != nil {
		return nil, err
	}

	switch len(c.Value) {
	case md5.Size:
		c.Hash = md5.New()
		c.Type = "md5"
	case sha1.Size:
		c.Hash = sha1.New()
		c.Type = "sha1"
	case sha256.Size:
		c.Hash = sha256.New()
		c.Type = "sha256"
	case sha512.Size:
		c.Hash = sha512.New()
		c.Type = "sha512"
	default:
		return nil, fmt.Errorf("Unknown type for checksum %s", checksumValue)
	}

	return c, nil
}

// checksumsFromFile will return all the fileChecksums found in file
//
// checksumsFromFile will try to guess the hashing algorithm based on content
// of checksum file
//
// checksumsFromFile will only return checksums for files that match file
// behind src
func (c *Client) checksumFromFile(checksumFile string, src *url.URL) (*fileChecksum, error) {
	checksumFileURL, err := urlhelper.Parse(checksumFile)
	if err != nil {
		return nil, err
	}

	tempfile, err := tmpFile("", filepath.Base(checksumFileURL.Path))
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempfile)

	c2 := &Client{
		Ctx:              c.Ctx,
		Getters:          c.Getters,
		Decompressors:    c.Decompressors,
		Detectors:        c.Detectors,
		Pwd:              c.Pwd,
		Dir:              false,
		Src:              checksumFile,
		Dst:              tempfile,
		ProgressListener: c.ProgressListener,
	}
	if err = c2.Get(); err != nil {
		return nil, fmt.Errorf(
			"Error downloading checksum file: %s", err)
	}

	filename := filepath.Base(src.Path)
	absPath, err := filepath.Abs(src.Path)
	if err != nil {
		return nil, err
	}
	checksumFileDir := filepath.Dir(checksumFileURL.Path)
	relpath, err := filepath.Rel(checksumFileDir, absPath)
	switch {
	case err == nil ||
		err.Error() == "Rel: can't make "+absPath+" relative to "+checksumFileDir:
		// ex: on windows C:\gopath\...\content.txt cannot be relative to \
		// which is okay, may be another expected path will work.
		break
	default:
		return nil, err
	}

	// possible file identifiers:
	options := []string{
		filename,       // ubuntu-14.04.1-server-amd64.iso
		"*" + filename, // *ubuntu-14.04.1-server-amd64.iso  Standard checksum
		"?" + filename, // ?ubuntu-14.04.1-server-amd64.iso  shasum -p
		relpath,        // dir/ubuntu-14.04.1-server-amd64.iso
		"./" + relpath, // ./dir/ubuntu-14.04.1-server-amd64.iso
		absPath,        // fullpath; set if local
	}

	f, err := os.Open(tempfile)
	if err != nil {
		return nil, fmt.Errorf(
			"Error opening downloaded file: %s", err)
	}
	defer f.Close()
	rd := bufio.NewReader(f)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return nil, fmt.Errorf(
					"Error reading checksum file: %s", err)
			}
			break
		}
		checksum, err := parseChecksumLine(line)
		if err != nil || checksum == nil {
			continue
		}
		if checksum.Filename == "" {
			// filename not sure, let's try
			return checksum, nil
		}
		// make sure the checksum is for the right file
		for _, option := range options {
			if option != "" && checksum.Filename == option {
				// any checksum will work so we return the first one
				return checksum, nil
			}
		}
	}
	return nil, fmt.Errorf("no checksum found in: %s", checksumFile)
}

// parseChecksumLine takes a line from a checksum file and returns
// checksumType, checksumValue and filename parseChecksumLine guesses the style
// of the checksum BSD vs GNU by splitting the line and by counting the parts.
// of a line.
// for BSD type sums parseChecksumLine guesses the hashing algorithm
// by checking the length of the checksum.
func parseChecksumLine(line string) (*fileChecksum, error) {
	parts := strings.Fields(line)

	switch len(parts) {
	case 4:
		// BSD-style checksum:
		//  MD5 (file1) = <checksum>
		//  MD5 (file2) = <checksum>
		if len(parts[1]) <= 2 ||
			parts[1][0] != '(' || parts[1][len(parts[1])-1] != ')' {
			return nil, fmt.Errorf(
				"Unexpected BSD-style-checksum filename format: %s", line)
		}
		filename := parts[1][1 : len(parts[1])-1]
		return newChecksumFromType(parts[0], parts[3], filename)
	case 2:
		// GNU-style:
		//  <checksum>  file1
		//  <checksum> *file2
		return newChecksumFromValue(parts[0], parts[1])
	case 0:
		return nil, nil // empty line
	default:
		return newChecksumFromValue(parts[0], "")
	}
}
