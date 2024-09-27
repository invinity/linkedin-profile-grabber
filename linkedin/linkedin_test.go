package linkedin

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLinkedIn(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LinkedIn Suite")
}

var _ = Describe("Using the LinkedIn profile retrieval", Ordered, func() {
	var browser *rod.Browser
	var linkedin *LinkedIn

	BeforeAll(func() {
		timeout, _ := time.ParseDuration("60s")
		chromePath, _ := os.LookupEnv("CHROME")
		browser = rod.New().ControlURL(launcher.New().Leakless(false).NoSandbox(true).Bin(chromePath).MustLaunch()).Trace(true).Timeout(timeout).MustConnect()
		browser.EachEvent(func(e *proto.NetworkResponseReceived) {
			log.Println(e)
		})
		linkedin = New(browser)
		enc := json.NewEncoder(log.Writer())
		enc.SetIndent("  ", "  ")
		enc.Encode(linkedin.RetrieveProfile())
	})

	Describe("normal function", func() {
		Context("loads basic profile data", func() {
			It("should load basic user info", func() {
				Ω(linkedin.RetrieveProfile().Name).Should(BeEquivalentTo("Matthew Pitts"))
				Ω(linkedin.RetrieveProfile().Headline).Should(Not(BeEmpty()))
			})
		})

		Context("loads Profile Experience data", func() {
			It("should load all experience info", func() {
				Ω(linkedin.RetrieveProfile().Experience).Should(HaveLen(4))
			})

			It("each Experience should have a Company", func() {
				exp := linkedin.RetrieveProfile().Experience
				for _, e := range exp {
					Ω(e.Company).Should(Not(BeEmpty()))
				}
			})

			It("each Experience should have one or more Positions", func() {
				exp := linkedin.RetrieveProfile().Experience
				for _, e := range exp {
					Ω(e.Positions).Should(Not(BeEmpty()))
				}
			})
		})
	})

	AfterAll(func() {
		browser.MustClose()
	})
})
