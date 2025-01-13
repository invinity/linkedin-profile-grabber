package controller

import (
	"errors"
	"testing"
	"time"

	"github.com/invinity/linkedin-profile-grabber/cache"
	"github.com/invinity/linkedin-profile-grabber/linkedin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type DummyProfileRetriever struct {
	profile *linkedin.LinkedInProfile
	err     error
}

func (r DummyProfileRetriever) Get() (*linkedin.LinkedInProfile, error) {
	return r.profile, r.err
}

func NewTestRetriever(profile *linkedin.LinkedInProfile, err error) LinkedinProfileRetriever {
	return DummyProfileRetriever{profile: profile, err: err}
}

func TestProfileRetriever(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ProfileRetriever Suite")
}

var _ = Describe("Using the profile retrieval", Ordered, func() {

	BeforeAll(func() {

	})

	Describe("normal function", func() {
		Context("cached copy is present but stale and linked in call fails", func() {
			cachedProfile := linkedin.LinkedInProfile{}
			testCache := cache.NewMemoryCache()
			testCache.Put("thing", cachedProfile)
			testRetriever := NewTestRetriever(nil, errors.New("unable to call linked in"))
			underTest := NewCacheHandlingRetriever(testCache, testRetriever)
			retrievedProfile, err := underTest.Get()
			It("should return the cached copy", func() {
				Ω(retrievedProfile).Should(BeEquivalentTo(&cachedProfile))
			})

			It("should return a nil error", func() {
				Ω(err).Should(BeNil())
			})
		})

		Context("cached copy is not present and linked in call succeeds", func() {
			freshProfile := linkedin.LinkedInProfile{}
			testCache := cache.NewMemoryCache()
			testRetriever := NewTestRetriever(&freshProfile, nil)
			underTest := NewCacheHandlingRetriever(testCache, testRetriever)
			retrievedProfile, err := underTest.Get()
			It("should return the fresh copy", func() {
				Ω(retrievedProfile).Should(BeEquivalentTo(&freshProfile))
			})

			It("should return a nil error", func() {
				Ω(err).Should(BeNil())
			})
		})

		Context("cached copy is present but stale and linked in call succeeds", func() {
			cachedProfile := linkedin.LinkedInProfile{GeneratedAt: time.Now().Add(-6 * time.Hour)}
			freshProfile := linkedin.LinkedInProfile{GeneratedAt: time.Now()}
			testCache := cache.NewMemoryCache()
			testCache.Put("thing", cachedProfile)
			testRetriever := NewTestRetriever(&freshProfile, nil)
			underTest := NewCacheHandlingRetriever(testCache, testRetriever)
			retrievedProfile, err := underTest.Get()
			It("should return the fresh copy", func() {
				Ω(retrievedProfile).Should(BeEquivalentTo(&freshProfile))
			})

			It("should return a nil error", func() {
				Ω(err).Should(BeNil())
			})
		})

		Context("cached copy is not present and linkedin call fails", func() {
			testCache := cache.NewMemoryCache()
			testRetriever := NewTestRetriever(nil, errors.New("unable to call linked in"))
			underTest := NewCacheHandlingRetriever(testCache, testRetriever)
			retrievedProfile, err := underTest.Get()
			It("should return a nil profile", func() {
				Ω(retrievedProfile).Should(BeNil())
			})

			It("should return an error", func() {
				Ω(err).ShouldNot(BeNil())
			})
		})
	})

	AfterAll(func() {
	})
})
