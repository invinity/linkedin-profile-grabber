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
		timeout, _ := time.ParseDuration("180s")
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
		profile, err = linkedin.RetrieveProfile("matthew", "pitts", "mattpitts")
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
				Ω(profile.Summary).Should(Not(BeEmpty()))
				Ω(profile.Summary).Should(Not(ContainSubstring("Show less")))
				Ω(profile.Summary).Should(Not(ContainSubstring("Show more")))
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
					Ω(e.CompanyImage).Should(Not(BeEquivalentTo("")))
					Ω(e.Positions).Should(Not(BeEmpty()))
				}
			})

			It("each Experience Position should have populate fields", func() {
				exp := profile.Experience
				for _, e := range exp {
					positions := e.Positions
					for _, v := range positions {
						Ω(v.Title).Should(Not(BeEmpty()))
						Ω(v.Location).Should(Not(BeEmpty()))
						Ω(v.Description).Should(Not(BeEmpty()))
						Ω(v.StartDate).Should(Not(BeEmpty()))
						Ω(v.Description).Should(Not(ContainSubstring("Show less")))
						Ω(v.Description).Should(Not(ContainSubstring("Show more")))
					}
				}
			})
		})

		Context("loads Project data", func() {
			It("should load all project info", func() {
				Ω(profile.Projects).Should(HaveLen(2))
			})

			It("Projects should have fields populated", func() {
				project := profile.Projects[0]
				Ω(project.Title).Should(BeEquivalentTo("Cloud Proxy Log Copying Automation"))
				Ω(project.StartDate).Should(BeEquivalentTo("Jun 2023"))
				Ω(project.EndDate).Should(BeEquivalentTo("Dec 2023"))
				Ω(project.Description).Should(Not(BeEmpty()))
				Ω(project.Description).Should(Not(ContainSubstring("Show less")))
				Ω(project.Description).Should(Not(ContainSubstring("Show more")))
			})
		})

		Context("loads Education data", func() {
			It("should load all education info", func() {
				Ω(profile.Education).Should(HaveLen(1))
			})

			It("Education should have fields populated", func() {
				edu := profile.Education[0]
				Ω(edu.Title).Should(BeEquivalentTo("University of North Carolina at Charlotte"))
				Ω(edu.Subtitle).Should(BeEquivalentTo("BSEE Electrical Engineering"))
				Ω(edu.StartDate).Should(BeEquivalentTo("1998"))
				Ω(edu.EndDate).Should(BeEquivalentTo("2003"))
				Ω(edu.Description).Should(BeEquivalentTo("Dean's List (4 of 8 semesters)"))
			})
		})

		Context("loads Certification data", func() {
			It("should load all certification info", func() {
				Ω(profile.Certifications).Should(HaveLen(2))
			})

			It("Certifications should have fields populated", func() {
				cert := profile.Certifications[1]
				Ω(cert.Title).Should(BeEquivalentTo("GIAC Secure Software Programmer - Java"))
				Ω(cert.Institution).Should(BeEquivalentTo("GIAC Certifications"))
				Ω(cert.IssuedOn).Should(BeEquivalentTo("May 2017"))
				Ω(cert.ExpiresOn).Should(BeEquivalentTo("May 2021"))
			})
		})
	})

	AfterAll(func() {
		browser.MustClose()
	})
})
