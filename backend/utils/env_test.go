package utils

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env", func() {
	It("should get value when environment variable is set and has value", func() {
		value := "normal-value"
		key := "superkey"

		err := os.Setenv(key, value)
		Expect(err).NotTo(HaveOccurred())

		v := GetenvOrDefault(key, "....")
		Expect(v).To(Equal(value))
	})

	It("should return default value when environment variable is not present", func() {
		dv := "default-value"

		v := GetenvOrDefault("super", dv)
		Expect(v).To(Equal(dv))
	})
})
