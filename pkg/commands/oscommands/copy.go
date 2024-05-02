package oscommands

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

/* MIT License
 *
 * Copyright (c) 2017 Roland Singer [roland.singer@desertbit.com]
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return //nolint: nakedret
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return //nolint: nakedret
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return //nolint: nakedret
	}

	err = out.Sync()
	if err != nil {
		return //nolint: nakedret
	}

	si, err := os.Stat(src)
	if err != nil {
		return //nolint: nakedret
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return //nolint: nakedret
	}

	return //nolint: nakedret
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist. If destination already exists we'll clobber it.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return errors.New("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return //nolint: nakedret
	}
	if err == nil {
		// it exists so let's remove it
		if err := os.RemoveAll(dst); err != nil {
			return err
		}
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return //nolint: nakedret
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return //nolint: nakedret
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return //nolint: nakedret
			}
		} else {
			var info os.FileInfo
			info, err = entry.Info()
			if err != nil {
				return //nolint: nakedret
			}

			// Skip symlinks.
			if info.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return //nolint: nakedret
			}
		}
	}

	return //nolint: nakedret
}
