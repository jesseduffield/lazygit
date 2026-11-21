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
	"path"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
)

// Our release archives are about 10 MB. Anything much bigger than that is not
// what we expect.
const maxArchiveSize = 200 * 1024 * 1024

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
