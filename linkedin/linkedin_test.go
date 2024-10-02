package linkedin

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
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
	var profile *LinkedInProfile

	BeforeAll(func() {
		timeout, _ := time.ParseDuration("60s")
		chromePath, present := os.LookupEnv("CHROME")
		nosandbox := false
		if present {
			nosandbox = true
		} else {
			chromePath, present = launcher.LookPath()
		}
		if present {
			log.Println("Using chrome from ENV: ", chromePath)
		} else {
			log.Fatal("Unable to find path to chrome")
		}
		browser = rod.New().ControlURL(launcher.New().Leakless(false).NoSandbox(nosandbox).Headless(true).Bin(chromePath).MustLaunch()).Trace(true).Timeout(timeout)
		// browser.EachEvent(func(e *proto.NetworkResponseReceived) {
		// 	log.Println(e)
		// })
		linkedin = New(browser.MustConnect())
		var err error
		profile, err = linkedin.RetrieveProfile()
		if err != nil {
			log.Fatal(err)
		}
		enc := json.NewEncoder(log.Writer())
		enc.SetIndent("  ", "  ")
		enc.Encode(profile)
	})

	Describe("normal function", func() {
		Context("loads basic profile data", func() {
			It("should load basic user info", func() {
				Ω(profile.Name).Should(BeEquivalentTo("Matthew Pitts"))
				Ω(profile.Headline).Should(Not(BeEmpty()))
			})
		})

		Context("loads Profile Experience data", func() {
			It("should load all experience info", func() {
				Ω(profile.Experience).Should(HaveLen(4))
			})

			It("each Experience should have a Company", func() {
				exp := profile.Experience
				for _, e := range exp {
					Ω(e.Company).Should(Not(BeEmpty()))
				}
			})

			It("each Experience should have one or more Positions", func() {
				exp := profile.Experience
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
