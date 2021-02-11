package ye_olde_pokemon_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestYeOldePokemon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "YeOldePokemon Suite")
}
