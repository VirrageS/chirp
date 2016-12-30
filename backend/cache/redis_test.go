package cache

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/VirrageS/chirp/backend/model"
)

type mockCacheConfigProvider struct{}

var cacheTime time.Duration = 10 * time.Millisecond

func (m *mockCacheConfigProvider) GetCacheExpirationTime() time.Duration {
	return cacheTime
}

func (m *mockCacheConfigProvider) GetPassword() string {
	return ""
}

func (m *mockCacheConfigProvider) GetHost() string {
	return "localhost"
}

func (m *mockCacheConfigProvider) GetPort() string {
	return "6379"
}

func (m *mockCacheConfigProvider) GetDB() int {
	return 0
}

var _ = Describe("RedisCache", func() {
	var (
		redisCache CacheProvider = NewRedisCache(&mockCacheConfigProvider{})

		objectTests []struct {
			in  interface{}
			out interface{}
		}

		fieldsTests []Fields
	)

	BeforeEach(func() {
		objectTests = []struct {
			in  interface{}
			out interface{}
		}{
			{
				&model.User{
					ID:       1,
					Username: "username",
				},
				&model.User{},
			},
			{
				&[]*model.User{
					{ID: 1, Username: "username1"},
					{ID: 2, Username: "username2"},
				},
				&[]*model.User{},
			},
			{
				&model.PublicUser{
					ID:        2,
					Username:  "username",
					Name:      "name",
					AvatarUrl: "avatars@web.com",
					Following: true,
				},
				&model.PublicUser{},
			},
		}

		fieldsTests = []Fields{
			{"key", 12, "hello", -1},
			{"key", 12},
			{-1, "key", 12},
			{0},
			{"key"},
		}
	})

	AfterEach(func() {
		redisCache.Flush()
	})

	Describe("tests without fields", func() {
		It("should set values without error", func() {
			for _, test := range objectTests {
				err := redisCache.Set("key", test.in)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should set integer value without error", func() {
			err := redisCache.SetInt("key", 0)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not find elements when trying to get without setting", func() {
			for _, test := range objectTests {
				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should not find integer element when trying to get without setting", func() {
			var integer int64
			exists, err := redisCache.GetInt("key", &integer)
			Expect(exists).Should(BeFalse())
			Expect(err).NotTo(HaveOccurred())
		})

		It("should find elements when trying to get after set", func() {
			for _, test := range objectTests {
				err := redisCache.Set("key", test.in)
				Expect(err).NotTo(HaveOccurred())

				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeTrue())
				Expect(err).NotTo(HaveOccurred())
				Expect(test.in).To(Equal(test.out))
			}
		})

		It("should find integer element when trying to get after set", func() {
			input := int64(10)

			err := redisCache.SetInt("key", input)
			Expect(err).NotTo(HaveOccurred())

			var output int64
			exists, err := redisCache.GetInt("key", &output)
			Expect(exists).Should(BeTrue())
			Expect(err).NotTo(HaveOccurred())
			Expect(input).To(Equal(output))
		})

		It("should set value of a non-existing key to 0 and then increment it", func() {
			err := redisCache.Increment("key")
			Expect(err).NotTo(HaveOccurred())

			var v int64
			exists, err := redisCache.GetInt("key", &v)
			Expect(exists).Should(BeTrue())
			Expect(int64(1)).To(Equal(v))
		})

		It("should increment a value after set", func() {
			input := int64(10)

			By("Setting key")

			err := redisCache.SetInt("key", input)
			Expect(err).NotTo(HaveOccurred())

			By("Incrementing key")

			err = redisCache.Increment("key")
			Expect(err).NotTo(HaveOccurred())

			By("Getting key")

			var v int64
			exists, err := redisCache.GetInt("key", &v)
			Expect(exists).Should(BeTrue())
			Expect((input + 1)).To(Equal(v))
		})

		It("should set value of a non-existing key to 0 and then decrement it", func() {
			err := redisCache.Decrement("key")
			Expect(err).NotTo(HaveOccurred())

			var v int64
			exists, err := redisCache.GetInt("key", &v)
			Expect(exists).Should(BeTrue())
			Expect(int64(-1)).To(Equal(v))
		})

		It("should decrement a value after set", func() {
			input := int64(10)

			By("Setting key")

			err := redisCache.SetInt("key", input)
			Expect(err).NotTo(HaveOccurred())

			By("Decrementing key")

			err = redisCache.Decrement("key")
			Expect(err).NotTo(HaveOccurred())

			By("Getting key")

			var v int64
			exists, err := redisCache.GetInt("key", &v)
			Expect(exists).Should(BeTrue())
			Expect((input - 1)).To(Equal(v))
		})

		It("should delete without errors when key does not exist", func() {
			err := redisCache.Delete("key")
			Expect(err).NotTo(HaveOccurred())
		})

		It("should delete keys without errors after set", func() {
			for _, test := range objectTests {
				By("Setting key")

				err := redisCache.Set("key", test.in)
				Expect(err).NotTo(HaveOccurred())

				By("Deleting key")

				err = redisCache.Delete("key")
				Expect(err).NotTo(HaveOccurred())

				By("Getting key")

				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should automatically delete keys after specified time in config", func() {
			for _, test := range objectTests {
				err := redisCache.Set("key", test.in)
				Expect(err).NotTo(HaveOccurred())

				time.Sleep(2 * cacheTime)

				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
		})
	})

	Describe("tests with fields", func() {
		It("should set values without error", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					err := redisCache.SetWithFields(fields, test.in)
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		It("should set integer values without error", func() {
			for _, fields := range fieldsTests {
				err := redisCache.SetIntWithFields(fields, int64(0))
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should not find elements when trying to get without setting", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		It("should not find integer elements when trying to get without setting", func() {
			var integer int64
			for _, fields := range fieldsTests {
				exists, err := redisCache.GetIntWithFields(fields, &integer)
				Expect(exists).Should(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should find elements when trying to get after set", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					err := redisCache.SetWithFields(fields, test.in)
					Expect(err).NotTo(HaveOccurred())

					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeTrue())
					Expect(err).NotTo(HaveOccurred())
					Expect(test.in).To(Equal(test.out))
				}
			}
		})

		It("should find integer elements when trying to get after set", func() {
			input := int64(10)
			var output int64

			for _, fields := range fieldsTests {
				err := redisCache.SetIntWithFields(fields, input)
				Expect(err).NotTo(HaveOccurred())

				exists, err := redisCache.GetIntWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect(err).NotTo(HaveOccurred())
				Expect(input).To(Equal(output))
			}
		})

		It("should set value of a non-existing key to 0 and then increment it", func() {
			for _, fields := range fieldsTests {
				err := redisCache.IncrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				var v int64
				exists, err := redisCache.GetIntWithFields(fields, &v)
				Expect(exists).Should(BeTrue())
				Expect(int64(1)).To(Equal(v))
			}
		})

		It("should increment a value after set", func() {
			for _, fields := range fieldsTests {
				input := int64(10)
				var output int64

				By("Setting key")

				err := redisCache.SetIntWithFields(fields, input)
				Expect(err).NotTo(HaveOccurred())

				By("Incrementing key")

				err = redisCache.IncrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				By("Getting key")

				exists, err := redisCache.GetIntWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect((input + 1)).To(Equal(output))
			}
		})

		It("should set value of a non-existing key to 0 and then decrement it", func() {
			for _, fields := range fieldsTests {
				err := redisCache.DecrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				var v int64
				exists, err := redisCache.GetIntWithFields(fields, &v)
				Expect(exists).Should(BeTrue())
				Expect(int64(-1)).To(Equal(v))
			}
		})

		It("should decrement a value after set", func() {
			for _, fields := range fieldsTests {
				input := int64(10)
				var output int64

				By("Setting key")

				err := redisCache.SetIntWithFields(fields, input)
				Expect(err).NotTo(HaveOccurred())

				By("Decrementing key")

				err = redisCache.DecrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				By("Getting key")

				exists, err := redisCache.GetIntWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect((input - 1)).To(Equal(output))
			}
		})

		It("should delete without errors when key does not exist", func() {
			for _, fields := range fieldsTests {
				err := redisCache.DeleteWithFields(fields)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should delete keys without errors after set", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					By("Setting key")

					err := redisCache.SetWithFields(fields, test.in)
					Expect(err).NotTo(HaveOccurred())

					By("Deleting key")

					err = redisCache.DeleteWithFields(fields)
					Expect(err).NotTo(HaveOccurred())

					By("Getting key")

					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		It("should automatically delete keys after specified time in config", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					err := redisCache.SetWithFields(fields, test.in)
					Expect(err).NotTo(HaveOccurred())

					time.Sleep(2 * cacheTime)

					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})
