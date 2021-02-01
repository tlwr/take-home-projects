package integration_test

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

func TestMonzoTakeHomeCrawler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MonzoTakeHomeCrawler Integration Suite")
}

var (
	path string
)

var _ = BeforeSuite(func() {
	var err error
	path, err = Build("github.com/tlwr/monzo-take-home-crawler")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	CleanupBuildArtifacts()
})

var _ = Describe("Usage", func() {
	It("displays a help message", func() {
		command := exec.Command(path, "-help")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(Exit(0))
		Eventually(session.Err).Should(Say("Usage of "))
		Eventually(session.Err).Should(Say("monzo-take-home-crawler"))
	})
})

var _ = Describe("-parallel", func() {
	Context("when zero", func() {
		It("displays a helpful error message", func() {
			command := exec.Command(path, "-parallel", "1024")
			session, err := Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(Exit(1))
			Eventually(session.Err).Should(Say("parallel flag should be between 1 and 256"))
		})
	})
})

var _ = Describe("Scraping", func() {
	It("scrapes my personal website", func() {
		command := exec.Command(path, "-url", "https://www.toby.codes", "-host", "toby.codes", "-host", "www.toby.codes")
		session, err := Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		Eventually(session).Should(Exit(0))

		lines := []string{
			"will include URLs within", "toby.codes",
			"queueing first url", "https://www.toby.codes",
			"results for page", "https://www.toby.codes",
			"https://github.com/tlwr",
			"results for page:", "https://www.toby.codes/posts",
			"https://www.toby.codes/posts/FOSDEM-2020",
		}
		for _, line := range lines {
			Eventually(session.Err).Should(Say(line))
		}
	})
})
