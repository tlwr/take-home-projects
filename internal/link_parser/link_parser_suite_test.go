package link_parser_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLinkParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkParser Suite")
}
