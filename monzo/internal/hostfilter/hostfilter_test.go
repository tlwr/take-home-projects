package hostfilter_test

import (
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tlwr/take-home-projects/monzo/internal/hostfilter"
)

var _ = Describe("Hostfilter", func() {
	Context("happy path", func() {
		var (
			err error
			hf  hostfilter.HostFilter

			boe *url.URL

			monzo *url.URL
			toby  *url.URL

			wwwDOTmonzo        *url.URL
			wwwDOTtobyDOTcodes *url.URL
		)

		BeforeEach(func() {
			hf, err = hostfilter.New([]string{
				"monzo.com",
				"*.monzo.com",
				"WWW.TOBY.CODES",
			})
			Expect(err).NotTo(HaveOccurred())

			boe, _ = url.Parse("https://www.bankofengland.co.uk")

			monzo, _ = url.Parse("https://monzo.com")
			wwwDOTmonzo, _ = url.Parse("https://www.monzo.com")

			toby, _ = url.Parse("https://toby.codes")
			wwwDOTtobyDOTcodes, _ = url.Parse("https://www.toby.codes")
		})

		It("allows allowed host", func() {
			Expect(hf.IsAllowed(monzo)).To(Equal(true))
			Expect(hf.IsAllowed(wwwDOTmonzo)).To(Equal(true))
			Expect(hf.IsAllowed(wwwDOTtobyDOTcodes)).To(Equal(true))
		})

		It("does not allow an absent host", func() {
			Expect(hf.IsAllowed(boe)).To(Equal(false))
			Expect(hf.IsAllowed(toby)).To(Equal(false))
		})
	})

	Context("malconstruction", func() {
		It("will not construct an invalid glob", func() {
			_, err := hostfilter.New([]string{
				"[!a-cat",
			})

			Expect(err).To(MatchError(ContainSubstring("host glob ([!a-cat) is not a valid glob")))
		})
	})
})
