package roll

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// -- Test helpers

type CustomError struct {
	s string
}

func (e *CustomError) Error() string {
	return e.s
}

func setup() func() {
	Token = os.Getenv("TOKEN")
	Environment = "test"

	if Token == "" {
		Token = "test"
		originalEndpoint := Endpoint
		server := httptest.NewServer(http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(`{"result": {"uuid": "01234567890123456789012345678901"}}`))
			},
		))

		Endpoint = server.URL

		return func() {
			Endpoint = originalEndpoint
			Token = ""
			server.Close()
		}
	}

	// Assume Token was provided and we want integration tests.
	return func() {}
}

// -- Tests

func TestErrorClass(t *testing.T) {
	errors := map[string]error{
		"{508e076d}":       fmt.Errorf("Something is broken!"),
		"roll.CustomError": &CustomError{"Terrible mistakes were made."},
	}

	for expected, err := range errors {
		if errorClass(err) != expected {
			t.Error("Got:", errorClass(err), "Expected:", expected)
		}
	}
}

func TestCritical(t *testing.T) {
	teardown := setup()
	defer teardown()

	uuid, err := Critical(errors.New("global critical"), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestError(t *testing.T) {
	teardown := setup()
	defer teardown()

	uuid, err := Error(errors.New("global error"), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestWarning(t *testing.T) {
	teardown := setup()
	defer teardown()

	uuid, err := Warning(errors.New("global warning"), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestInfo(t *testing.T) {
	teardown := setup()
	defer teardown()

	uuid, err := Info("global info", map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestDebug(t *testing.T) {
	teardown := setup()
	defer teardown()

	uuid, err := Debug("global debug", map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientCritical(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.Critical(errors.New("new client critical"), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientError(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.Error(errors.New("new client error"), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientWarning(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.Warning(errors.New("new client warning"), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientInfo(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.Info("new client info", map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientDebug(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.Debug("new client debug", map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientCriticalStack(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.CriticalStack(errors.New("new client critical"), getCallers(0), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientErrorStack(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.ErrorStack(errors.New("new client error"), getCallers(0), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}

func TestRollbarClientWarningStack(t *testing.T) {
	teardown := setup()
	defer teardown()

	client := New(Token, "test")

	uuid, err := client.WarningStack(errors.New("new client warning"), getCallers(0), map[string]string{"extras": "true"})
	if err != nil {
		t.Error(err)
	}
	if len(uuid) != 32 {
		t.Errorf("expected UUID, got: %#v", uuid)
	}
}
