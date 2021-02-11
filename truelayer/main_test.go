package main_test

import (
	"os/exec"
	"testing"

	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gexec"
)

func TestTruelayerTakeHomePokemonApi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TruelayerTakeHomePokemonApi Suite")
}

var (
	path string
)

var _ = BeforeSuite(func() {
	var err error
	path, err = gexec.Build("github.com/tlwr/truelayer-take-home-pokemon-api")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

var _ = Describe("main", func() {
	var (
		session *gexec.Session

		resp *http.Response
		body []byte

		pokemonName string
	)

	BeforeEach(func() {
		var err error
		command := exec.Command(path, "-stubs", "-bind", "127.0.0.1", "-port", "8080")
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() int {
			r, _ := http.Get("http://localhost:8080/")
			if r == nil {
				return 0
			}
			return r.StatusCode
		}).Should(Equal(404), "waiting for server to be up")
	})

	JustBeforeEach(func() {
		var err error

		resp, err = http.Get("http://localhost:8080/pokemon/" + pokemonName)
		Expect(err).NotTo(HaveOccurred())
		defer resp.Body.Close()

		Expect(resp.Header.Get("Content-Type")).To(Equal("application/json"))

		body, err = ioutil.ReadAll(resp.Body)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		session.Terminate().Wait()
	})

	Context("when getting a known pokemon", func() {
		BeforeEach(func() {
			pokemonName = "pikachu"
		})

		It("returns HTTP 200 with name and translation", func() {
			Expect(resp.StatusCode).To(Equal(200))
			Expect(body).To(MatchJSON(`{"name": "pikachu", "description": "you speak an infinite deal of nothing"}`))
		})
	})

	Context("when getting a non-existent pokemon", func() {
		BeforeEach(func() {
			pokemonName = "batman"
		})

		It("returns HTTP 404", func() {
			Expect(resp.StatusCode).To(Equal(404))
			Expect(body).To(MatchJSON(`{"message": "not found"}`))
		})
	})

	Context("when trying to get missingno", func() {
		BeforeEach(func() {
			pokemonName = "missingno"
		})

		It("returns HTTP 500", func() {
			Expect(resp.StatusCode).To(Equal(500))
			Expect(body).To(MatchJSON(`{"message": "missingno"}`))
		})
	})
})
