package jibber_jabber_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestJibberJabber(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jibber Jabber Suite")
}
