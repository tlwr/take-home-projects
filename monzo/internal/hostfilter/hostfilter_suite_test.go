package hostfilter_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHostfilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hostfilter Suite")
}
