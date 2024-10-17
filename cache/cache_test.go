package cache

import (
	"context"
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache Suite")
}

var _ = Describe("Using the LinkedIn profile retrieval", Ordered, func() {
	var ctx context.Context
	var underTest *Cache

	BeforeAll(func() {
		ctx = context.Background()
		var err error
		underTest, err = NewCache(&ctx, "linkedin-profile-grabber")
		if err != nil {
			log.Panic(err)
		}
	})

	Describe("normal function", func() {
		Context("performs standard get/put operations", func() {
			It("should allow a put of a basic string value", func() {
				err := underTest.Put("test", "this is a string value")
				Ω(err).Should(BeNil())
			})

			It("should allow a read of a basic string value", func() {
				var value string
				err := underTest.Get("test", &value)
				Ω(err).Should(BeNil())
				Ω(value).Should(BeEquivalentTo("this is a string value"))
			})

			It("should allow the removal of an entry", func() {
				err := underTest.Remove("test")
				Ω(err).Should(BeNil())
			})
		})

		Context("error conditions", func() {
			It("attempted get of object that does not exist should produce proper error", func() {
				var value string
				err := underTest.Get("thisobjectdoesnotexist", &value)
				Ω(err).ShouldNot(BeNil())
				Ω(value).Should(BeEmpty())
			})
		})
	})

	AfterAll(func() {
		underTest.Close()
	})
})
