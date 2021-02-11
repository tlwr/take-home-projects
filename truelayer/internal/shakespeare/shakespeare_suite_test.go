package shakespeare_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestShakespeare(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Shakespeare Suite")
}
