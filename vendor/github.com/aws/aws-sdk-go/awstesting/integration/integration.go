// +build integration

// Package integration performs initialization and validation for integration
// tests.
package integration

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Session is a shared session for all integration tests to use.
var Session = session.Must(session.NewSession())

func init() {
	logLevel := Session.Config.LogLevel
	if os.Getenv("DEBUG") != "" {
		logLevel = aws.LogLevel(aws.LogDebug)
	}
	if os.Getenv("DEBUG_SIGNING") != "" {
		logLevel = aws.LogLevel(aws.LogDebugWithSigning)
	}
	if os.Getenv("DEBUG_BODY") != "" {
		logLevel = aws.LogLevel(aws.LogDebugWithSigning | aws.LogDebugWithHTTPBody)
	}
	Session.Config.LogLevel = logLevel
}

// UniqueID returns a unique UUID-like identifier for use in generating
// resources for integration tests.
func UniqueID() string {
	uuid := make([]byte, 16)
	io.ReadFull(rand.Reader, uuid)
	return fmt.Sprintf("%x", uuid)
}

// SessionWithDefaultRegion returns a copy of the integration session with the
// region set if one was not already provided.
func SessionWithDefaultRegion(region string) *session.Session {
	sess := Session.Copy()
	if v := aws.StringValue(sess.Config.Region); len(v) == 0 {
		sess.Config.Region = aws.String(region)
	}

	return sess
}
