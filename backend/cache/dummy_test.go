package cache

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type object struct {
	Str string
	Num int
}

var dummyCache CacheProvider = NewDummyCache()

var _ = Describe("DummyCache", func() {
	var (
		dummyCache CacheProvider = NewDummyCache()

		objectTest *object
		fieldTest  Fields
	)

	BeforeEach(func() {
		objectTest = &object{
			Str: "wtf",
			Num: 12,
		}

		fieldTest = Fields{"key", "super", 1}
	})

	Describe("tests without fields", func() {
		It("should set values without error", func() {
			err := dummyCache.Set("key", objectTest)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not find elements when trying to get without setting", func() {
			exists, err := dummyCache.Get("key", &objectTest)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not find elements when trying to get after set", func() {
			err := dummyCache.Set("key", objectTest)
			Expect(err).NotTo(HaveOccurred())

			var obj object
			exists, err := dummyCache.Get("key", &obj)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete without errors when key does not exist", func() {
			err := dummyCache.Delete("key")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete keys without errors after set", func() {
			err := dummyCache.Set("key", objectTest)
			Expect(err).NotTo(HaveOccurred())

			err = dummyCache.Delete("key")
			Expect(err).NotTo(HaveOccurred())

			var obj object
			exists, err := dummyCache.Get("key", &obj)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("tests with fields", func() {
		It("should set values without error", func() {
			err := dummyCache.SetWithFields(fieldTest, objectTest)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not find elements when trying to get without setting", func() {
			exists, err := dummyCache.GetWithFields(fieldTest, objectTest)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not find elements when trying to get after set", func() {
			err := dummyCache.SetWithFields(fieldTest, objectTest)
			Expect(err).NotTo(HaveOccurred())

			var obj object
			exists, err := dummyCache.GetWithFields(fieldTest, &obj)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete without errors when key does not exist", func() {
			err := dummyCache.DeleteWithFields(fieldTest)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete keys without errors after set", func() {
			err := dummyCache.SetWithFields(fieldTest, objectTest)
			Expect(err).NotTo(HaveOccurred())

			err = dummyCache.DeleteWithFields(fieldTest)
			Expect(err).NotTo(HaveOccurred())

			var obj object
			exists, err := dummyCache.GetWithFields(fieldTest, &obj)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
