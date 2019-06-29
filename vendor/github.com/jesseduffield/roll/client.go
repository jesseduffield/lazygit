package roll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	// By default, all Rollbar API requests are sent to this endpoint.
	endpoint = "https://api.rollbar.com/api/1/item/"

	// Identify this Rollbar client library to the Rollbar API.
	clientName     = "go-roll"
	clientVersion  = "0.2.0"
	clientLanguage = "go"
)

var (
	// Endpoint is the default HTTP(S) endpoint that all Rollbar API requests
	// will be sent to. By default, this is Rollbar's "Items" API endpoint. If
	// this is blank, no items will be sent to Rollbar.
	Endpoint = endpoint

	// Rollbar access token for the global client. If this is blank, no items
	// will be sent to Rollbar.
	Token = ""

	// Environment for all items reported with the global client.
	Environment = "development"
)

type rollbarSuccess struct {
	Result map[string]string `json:"result"`
}

// Client reports items to a single Rollbar project.
type Client interface {
	Critical(err error, custom map[string]string) (uuid string, e error)
	CriticalStack(err error, ptrs []uintptr, custom map[string]string) (uuid string, e error)
	Error(err error, custom map[string]string) (uuid string, e error)
	ErrorStack(err error, ptrs []uintptr, custom map[string]string) (uuid string, e error)
	Warning(err error, custom map[string]string) (uuid string, e error)
	WarningStack(err error, ptrs []uintptr, custom map[string]string) (uuid string, e error)
	Info(msg string, custom map[string]string) (uuid string, e error)
	Debug(msg string, custom map[string]string) (uuid string, e error)
}

type rollbarClient struct {
	token string
	env   string
}

// New creates a new Rollbar client that reports items to the given project
// token and with the given environment (eg. "production", "development", etc).
func New(token, env string) Client {
	return &rollbarClient{token, env}
}

func Critical(err error, custom map[string]string) (uuid string, e error) {
	return CriticalStack(err, getCallers(2), custom)
}

func CriticalStack(err error, ptrs []uintptr, custom map[string]string) (uuid string, e error) {
	return New(Token, Environment).CriticalStack(err, ptrs, custom)
}

func Error(err error, custom map[string]string) (uuid string, e error) {
	return ErrorStack(err, getCallers(2), custom)
}

func ErrorStack(err error, ptrs []uintptr, custom map[string]string) (uuid string, e error) {
	return New(Token, Environment).ErrorStack(err, ptrs, custom)
}

func Warning(err error, custom map[string]string) (uuid string, e error) {
	return WarningStack(err, getCallers(2), custom)
}

func WarningStack(err error, ptrs []uintptr, custom map[string]string) (uuid string, e error) {
	return New(Token, Environment).WarningStack(err, ptrs, custom)
}

func Info(msg string, custom map[string]string) (uuid string, e error) {
	return New(Token, Environment).Info(msg, custom)
}

func Debug(msg string, custom map[string]string) (uuid string, e error) {
	return New(Token, Environment).Debug(msg, custom)
}

func (c *rollbarClient) Critical(err error, custom map[string]string) (uuid string, e error) {
	return c.CriticalStack(err, getCallers(2), custom)
}

func (c *rollbarClient) CriticalStack(err error, callers []uintptr, custom map[string]string) (uuid string, e error) {
	item := c.buildTraceItem("critical", err, callers, custom)
	return c.send(item)
}

func (c *rollbarClient) Error(err error, custom map[string]string) (uuid string, e error) {
	return c.ErrorStack(err, getCallers(2), custom)
}

func (c *rollbarClient) ErrorStack(err error, callers []uintptr, custom map[string]string) (uuid string, e error) {
	item := c.buildTraceItem("error", err, callers, custom)
	return c.send(item)
}

func (c *rollbarClient) Warning(err error, custom map[string]string) (uuid string, e error) {
	return c.WarningStack(err, getCallers(2), custom)
}

func (c *rollbarClient) WarningStack(err error, callers []uintptr, custom map[string]string) (uuid string, e error) {
	item := c.buildTraceItem("warning", err, callers, custom)
	return c.send(item)
}

func (c *rollbarClient) Info(msg string, custom map[string]string) (uuid string, e error) {
	item := c.buildMessageItem("info", msg, custom)
	return c.send(item)
}

func (c *rollbarClient) Debug(msg string, custom map[string]string) (uuid string, e error) {
	item := c.buildMessageItem("debug", msg, custom)
	return c.send(item)
}

func (c *rollbarClient) buildTraceItem(level string, err error, callers []uintptr, custom map[string]string) (item map[string]interface{}) {
	stack := buildRollbarFrames(callers)
	item = c.buildItem(level, err.Error(), custom)
	itemData := item["data"].(map[string]interface{})
	itemData["fingerprint"] = stack.fingerprint()
	itemData["body"] = map[string]interface{}{
		"trace": map[string]interface{}{
			"frames": stack,
			"exception": map[string]interface{}{
				"class":   errorClass(err),
				"message": err.Error(),
			},
		},
	}

	return item
}

func (c *rollbarClient) buildMessageItem(level string, msg string, custom map[string]string) (item map[string]interface{}) {
	item = c.buildItem(level, msg, custom)
	itemData := item["data"].(map[string]interface{})
	itemData["body"] = map[string]interface{}{
		"message": map[string]interface{}{
			"body": msg,
		},
	}

	return item
}

func (c *rollbarClient) buildItem(level, title string, custom map[string]string) map[string]interface{} {
	hostname, _ := os.Hostname()

	return map[string]interface{}{
		"access_token": c.token,
		"data": map[string]interface{}{
			"environment": c.env,
			"title":       title,
			"level":       level,
			"timestamp":   time.Now().Unix(),
			"platform":    runtime.GOOS,
			"language":    clientLanguage,
			"server": map[string]interface{}{
				"host": hostname,
			},
			"notifier": map[string]interface{}{
				"name":    clientName,
				"version": clientVersion,
			},
			"custom": custom,
		},
	}
}

// send reports the given item to Rollbar and returns either a UUID for the
// reported item or an error.
func (c *rollbarClient) send(item map[string]interface{}) (uuid string, err error) {
	if len(c.token) == 0 || len(Endpoint) == 0 {
		return "", nil
	}

	jsonBody, err := json.Marshal(item)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(Endpoint, "application/json", bytes.NewReader(jsonBody))
	if err != nil {
		// If something goes wrong it really does not matter
		return "", nil
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		// If something goes wrong it really does not matter
		return "", nil
	}

	// Extract UUID from JSON response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}
	success := rollbarSuccess{}
	json.Unmarshal(body, &success)

	return success.Result["uuid"], nil
}

// errorClass returns a class name for an error (eg.  "ErrUnexpectedEOF"). For
// string errors, it returns an Adler-32 checksum of the error string.
func errorClass(err error) string {
	class := reflect.TypeOf(err).String()
	if class == "" {
		return "panic"
	} else if class == "*errors.errorString" {
		checksum := adler32.Checksum([]byte(err.Error()))
		return fmt.Sprintf("{%x}", checksum)
	} else {
		return strings.TrimPrefix(class, "*")
	}
}
