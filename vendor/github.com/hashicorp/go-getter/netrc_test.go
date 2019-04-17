package getter

import (
	"net/url"
	"testing"
)

func TestAddAuthFromNetrc(t *testing.T) {
	defer tempEnv(t, "NETRC", "./test-fixtures/netrc/basic")()

	u, err := url.Parse("http://example.com")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if err := addAuthFromNetrc(u); err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := "http://foo:bar@example.com"
	actual := u.String()
	if expected != actual {
		t.Fatalf("Mismatch: %q != %q", actual, expected)
	}
}

func TestAddAuthFromNetrc_hasAuth(t *testing.T) {
	defer tempEnv(t, "NETRC", "./test-fixtures/netrc/basic")()

	u, err := url.Parse("http://username:password@example.com")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := u.String()
	if err := addAuthFromNetrc(u); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := u.String()
	if expected != actual {
		t.Fatalf("Mismatch: %q != %q", actual, expected)
	}
}

func TestAddAuthFromNetrc_hasUsername(t *testing.T) {
	defer tempEnv(t, "NETRC", "./test-fixtures/netrc/basic")()

	u, err := url.Parse("http://username@example.com")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := u.String()
	if err := addAuthFromNetrc(u); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual := u.String()
	if expected != actual {
		t.Fatalf("Mismatch: %q != %q", actual, expected)
	}
}
