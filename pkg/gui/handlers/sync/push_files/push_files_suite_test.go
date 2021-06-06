package push_files_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPushFiles(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PushFiles Suite")
}
