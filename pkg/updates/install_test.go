package updates

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/gui/types"
	"github.com/stretchr/testify/assert"
)

func TestAssetName(t *testing.T) {
	// The assets of the v0.66.2 release
	expectedNames := []string{
		"lazygit_0.66.2_darwin_arm64.tar.gz",
		"lazygit_0.66.2_darwin_x86_64.tar.gz",
		"lazygit_0.66.2_freebsd_32-bit.tar.gz",
		"lazygit_0.66.2_freebsd_arm64.tar.gz",
		"lazygit_0.66.2_freebsd_armv6.tar.gz",
		"lazygit_0.66.2_freebsd_x86_64.tar.gz",
		"lazygit_0.66.2_linux_32-bit.tar.gz",
		"lazygit_0.66.2_linux_arm64.tar.gz",
		"lazygit_0.66.2_linux_armv6.tar.gz",
		"lazygit_0.66.2_linux_x86_64.tar.gz",
		"lazygit_0.66.2_windows_32-bit.zip",
		"lazygit_0.66.2_windows_arm64.zip",
		"lazygit_0.66.2_windows_x86_64.zip",
	}

	version := &types.VersionNumber{Major: 0, Minor: 66, Patch: 2}
	names := []string{}
	for _, goos := range []string{"darwin", "freebsd", "linux", "windows"} {
		for _, goarch := range []string{"arm64", "amd64", "386", "arm"} {
			// goreleaser doesn't build these
			if (goos == "darwin" && (goarch == "386" || goarch == "arm")) || (goos == "windows" && goarch == "arm") {
				continue
			}
			names = append(names, assetName(version, goos, goarch))
		}
	}

	assert.ElementsMatch(t, expectedNames, names)
}

func TestAssetURL(t *testing.T) {
	version := &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}
	assert.Equal(t,
		"https://github.com/jesseduffield/lazygit/releases/download/v0.67.0/lazygit_0.67.0_linux_x86_64.tar.gz",
		assetURL("https://github.com/jesseduffield/lazygit", version, "linux", "amd64"))
}

func TestExtractBinary(t *testing.T) {
	files := map[string]string{
		"LICENSE":     "license",
		"README.md":   "readme",
		"lazygit":     "unix binary",
		"lazygit.exe": "windows binary",
	}

	binary, err := extractBinary(makeTarGz(t, files), "lazygit_0.67.0_linux_x86_64.tar.gz", "lazygit")
	assert.NoError(t, err)
	assert.Equal(t, "unix binary", string(binary))

	binary, err = extractBinary(makeZip(t, files), "lazygit_0.67.0_windows_x86_64.zip", "lazygit.exe")
	assert.NoError(t, err)
	assert.Equal(t, "windows binary", string(binary))

	_, err = extractBinary(makeTarGz(t, map[string]string{"LICENSE": "license"}), "archive.tar.gz", "lazygit")
	assert.EqualError(t, err, "archive.tar.gz doesn't contain lazygit")

	_, err = extractBinary([]byte("not an archive"), "archive.tar.gz", "lazygit")
	assert.Error(t, err)
}

func TestDownloadBinary(t *testing.T) {
	archive := makeTarGz(t, map[string]string{"LICENSE": "license", "lazygit": "new"})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/releases/download/v0.67.0/lazygit_0.67.0_linux_x86_64.tar.gz" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write(archive)
	}))
	defer server.Close()

	version := &types.VersionNumber{Major: 0, Minor: 67, Patch: 0}
	url := assetURL(server.URL, version, "linux", "amd64")
	binary, err := downloadBinary(context.Background(), server.Client(), url, "linux")
	assert.NoError(t, err)
	assert.Equal(t, "new", string(binary))

	url = assetURL(server.URL, version, "linux", "arm64")
	_, err = downloadBinary(context.Background(), server.Client(), url, "linux")
	assert.ErrorContains(t, err, "lazygit_0.67.0_linux_arm64.tar.gz returned 404 Not Found")
}

func makeTarGz(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buffer bytes.Buffer
	gzipWriter := gzip.NewWriter(&buffer)
	tarWriter := tar.NewWriter(gzipWriter)
	for name, content := range files {
		header := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}
		assert.NoError(t, tarWriter.WriteHeader(header))
		_, err := tarWriter.Write([]byte(content))
		assert.NoError(t, err)
	}
	assert.NoError(t, tarWriter.Close())
	assert.NoError(t, gzipWriter.Close())
	return buffer.Bytes()
}

func makeZip(t *testing.T, files map[string]string) []byte {
	t.Helper()

	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)
	for name, content := range files {
		writer, err := zipWriter.Create(name)
		assert.NoError(t, err)
		_, err = writer.Write([]byte(content))
		assert.NoError(t, err)
	}
	assert.NoError(t, zipWriter.Close())
	return buffer.Bytes()
}
