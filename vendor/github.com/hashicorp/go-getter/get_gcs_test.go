package getter

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

// initGCPCredentials writes a temporary GCS credentials file if necessary and
// returns the path and a function to clean it up. allAuthenticatedUsers can
// access go-getter-test with read only access.
func initGCPCredentials(t *testing.T) func() {
	if gc := os.Getenv("GOOGLE_CREDENTIALS"); gc != "" &&
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") == "" {
		file, cleanup := tempFileContents(t, gc)
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", file)
		return func() {
			os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "")
			cleanup()
		}
	}
	return func() {}
}

func TestGCSGetter_impl(t *testing.T) {
	var _ Getter = new(GCSGetter)
}

func TestGCSGetter(t *testing.T) {
	defer initGCPCredentials(t)()

	g := new(GCSGetter)
	dst := tempDir(t)

	// With a dir that doesn't exist
	err := g.Get(
		dst, testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/folder"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	mainPath := filepath.Join(dst, "main.tf")
	if _, err := os.Stat(mainPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGCSGetter_subdir(t *testing.T) {
	defer initGCPCredentials(t)()

	g := new(GCSGetter)
	dst := tempDir(t)

	// With a dir that doesn't exist
	err := g.Get(
		dst, testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/folder/subfolder"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the sub file exists
	subPath := filepath.Join(dst, "sub.tf")
	if _, err := os.Stat(subPath); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGCSGetter_GetFile(t *testing.T) {
	defer initGCPCredentials(t)()

	g := new(GCSGetter)
	dst := tempTestFile(t)
	defer os.RemoveAll(filepath.Dir(dst))

	// Download
	err := g.GetFile(
		dst, testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/folder/main.tf"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the main file exists
	if _, err := os.Stat(dst); err != nil {
		t.Fatalf("err: %s", err)
	}
	assertContents(t, dst, "# Main\n")
}

func TestGCSGetter_GetFile_notfound(t *testing.T) {
	g := new(GCSGetter)
	dst := tempTestFile(t)
	defer os.RemoveAll(filepath.Dir(dst))

	// Download
	err := g.GetFile(
		dst, testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/folder/404.tf"))
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGCSGetter_ClientMode_dir(t *testing.T) {
	defer initGCPCredentials(t)()

	g := new(GCSGetter)

	// Check client mode on a key prefix with only a single key.
	mode, err := g.ClientMode(
		testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/folder/subfolder"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if mode != ClientModeDir {
		t.Fatal("expect ClientModeDir")
	}
}

func TestGCSGetter_ClientMode_file(t *testing.T) {
	defer initGCPCredentials(t)()

	g := new(GCSGetter)

	// Check client mode on a key prefix which contains sub-keys.
	mode, err := g.ClientMode(
		testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/folder/subfolder/sub.tf"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if mode != ClientModeFile {
		t.Fatal("expect ClientModeFile")
	}
}

func TestGCSGetter_ClientMode_notfound(t *testing.T) {
	defer initGCPCredentials(t)()

	g := new(GCSGetter)

	// Check the client mode when a non-existent key is looked up. This does not
	// return an error, but rather should just return the file mode.
	mode, err := g.ClientMode(
		testURL("https://www.googleapis.com/storage/v1/go-getter-test/go-getter/foobar"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if mode != ClientModeFile {
		t.Fatal("expect ClientModeFile")
	}
}

func TestGCSGetter_Url(t *testing.T) {
	defer initGCPCredentials(t)()

	var gcstests = []struct {
		name   string
		url    string
		bucket string
		path   string
	}{
		{
			name:   "test1",
			url:    "gcs::https://www.googleapis.com/storage/v1/go-getter-test/go-getter/foo/null.zip",
			bucket: "go-getter-test",
			path:   "go-getter/foo/null.zip",
		},
	}

	for i, pt := range gcstests {
		t.Run(pt.name, func(t *testing.T) {
			g := new(GCSGetter)
			forced, src := getForcedGetter(pt.url)
			u, err := url.Parse(src)

			if err != nil {
				t.Errorf("test %d: unexpected error: %s", i, err)
			}
			if forced != "gcs" {
				t.Fatalf("expected forced protocol to be gcs")
			}

			bucket, path, err := g.parseURL(u)

			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if bucket != pt.bucket {
				t.Fatalf("expected %s, got %s", pt.bucket, bucket)
			}
			if path != pt.path {
				t.Fatalf("expected %s, got %s", pt.path, path)
			}
		})
	}
}
