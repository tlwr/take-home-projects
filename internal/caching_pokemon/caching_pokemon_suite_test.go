package caching_pokemon_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCachingPokemon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CachingPokemon Suite")
}
