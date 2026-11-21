package updates

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Our release archives are about 10 MB. Anything much bigger than that is not
// what we expect.
const maxArchiveSize = 200 * 1024 * 1024

// The prefix of the temporary file that we write the new binary to. It is in
// the folder of the executable, so that we can rename it over the executable.
const tempFilePrefix = ".lazygit-update-"

// assetName returns the name of the release archive for the given platform,
// following the name_template in .goreleaser.yml.
func assetName(version *types.VersionNumber, goos string, goarch string) string {
	arch := goarch
	switch goarch {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "32-bit"
	case "arm":
		arch = "armv6"
	}

	extension := "tar.gz"
	if goos == "windows" {
		extension = "zip"
	}

	return fmt.Sprintf("lazygit_%d.%d.%d_%s_%s.%s", version.Major, version.Minor, version.Patch, goos, arch, extension)
}

func assetURL(repoURL string, version *types.VersionNumber, goos string, goarch string) string {
	return fmt.Sprintf("%s/releases/download/%s/%s", repoURL, version.String(), assetName(version, goos, goarch))
}

func binaryName(goos string) string {
	if goos == "windows" {
		return "lazygit.exe"
	}
	return "lazygit"
}

// checkCanReplaceExecutable returns an error if we can't create files in the
// folder of the executable. In that case we can't install an update.
func checkCanReplaceExecutable(executablePath string) error {
	file, err := os.CreateTemp(filepath.Dir(executablePath), tempFilePrefix+"*")
	if err != nil {
		return err
	}
	_ = file.Close()
	return os.Remove(file.Name())
}

// downloadBinary downloads the release archive from url, and returns the
// lazygit binary in it.
func downloadBinary(ctx context.Context, client *http.Client, url string, goos string) ([]byte, error) {
	archive, err := download(ctx, client, url)
	if err != nil {
		return nil, err
	}

	return extractBinary(archive, url, binaryName(goos))
}

func download(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s returned %s", url, resp.Status)
	}

	content, err := io.ReadAll(io.LimitReader(resp.Body, maxArchiveSize+1))
	if err != nil {
		return nil, err
	}
	if len(content) > maxArchiveSize {
		return nil, fmt.Errorf("%s is larger than %d bytes", url, maxArchiveSize)
	}
	return content, nil
}

func extractBinary(archive []byte, archiveName string, name string) ([]byte, error) {
	var binary []byte
	var err error
	if strings.HasSuffix(archiveName, ".zip") {
		binary, err = extractFromZip(archive, name)
	} else {
		binary, err = extractFromTarGz(archive, name)
	}
	if err != nil {
		return nil, err
	}
	if binary == nil {
		return nil, fmt.Errorf("%s doesn't contain %s", archiveName, name)
	}
	return binary, nil
}

func extractFromTarGz(archive []byte, name string) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		if header.Typeflag == tar.TypeReg && path.Base(header.Name) == name {
			return io.ReadAll(tarReader)
		}
	}
}

func extractFromZip(archive []byte, name string) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		return nil, err
	}
	for _, file := range zipReader.File {
		if !file.FileInfo().IsDir() && path.Base(file.Name) == name {
			reader, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer reader.Close()
			return io.ReadAll(reader)
		}
	}
	return nil, nil
}

// writeNewExecutable writes the new binary to a temporary file next to the
// executable, with the same permissions, and returns the path of that file.
// replaceExecutable can then rename it over the executable.
func writeNewExecutable(executablePath string, binary []byte) (string, error) {
	info, err := os.Stat(executablePath)
	if err != nil {
		return "", err
	}

	return writeTempFile(filepath.Dir(executablePath), binary, info.Mode().Perm())
}

// replaceExecutable renames the file that writeNewExecutable wrote over the
// executable. Renaming makes the change atomic, and gives the new binary a new
// inode. The latter matters on macOS, where the system kills a signed binary if
// it was overwritten in place. If renaming fails, it removes the new file.
func replaceExecutable(executablePath string, newPath string) error {
	err := os.Rename(newPath, executablePath)
	if err != nil {
		_ = os.Remove(newPath)
	}
	return err
}

func writeTempFile(dir string, content []byte, mode os.FileMode) (string, error) {
	file, err := os.CreateTemp(dir, tempFilePrefix+"*")
	if err != nil {
		return "", err
	}

	_, err = file.Write(content)
	if err == nil {
		err = file.Sync()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		err = os.Chmod(file.Name(), mode)
	}
	if err != nil {
		_ = os.Remove(file.Name())
		return "", err
	}
	return file.Name(), nil
}
