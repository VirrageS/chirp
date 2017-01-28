package cache

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/model"
)

var _ = Describe("RedisCache", func() {
	var (
		conf       *config.Configuration = config.New()
		redisCache CacheProvider         = NewRedisCache(conf.Redis)

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
			// TODO: make this work
			// {
			// 	&model.Tweet{
			// 		ID: 2,
			// 		Author: &model.PublicUser {
			// 			ID: 1,
			// 			Username: "wtf",
			// 		},
			// 		LikeCount: 10,
			// 		RetweetCount: 200,
			// 		Retweeted: true,
			// 	},
			// 	&model.Tweet{},
			// },
			// {
			// 	&[]*model.Tweet{
			// 		{
			// 			ID: 2,
			// 			Author: &model.PublicUser {
			// 				ID: 1,
			// 				Username: "wtf",
			// 			},
			// 			LikeCount: 10,
			// 			RetweetCount: 200,
			// 			Retweeted: true,
			// 		},
			// 		{
			// 			ID: 10,
			// 			Author: &model.PublicUser {
			// 				ID: 10,
			// 				Username: "lolek",
			// 			},
			// 			LikeCount: 0,
			// 			RetweetCount: 1010101010,
			// 			Retweeted: false,
			// 		},
			// 	},
			// 	&[]*model.Tweet{},
			// },
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

		It("should set values without expiration without error", func() {
			for _, test := range objectTests {
				err := redisCache.SetWithoutExpiration("key", test.in)
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should not find elements when trying to get without setting", func() {
			for _, test := range objectTests {
				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
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

		It("should find elements when trying to get after set without expiration", func() {
			for _, test := range objectTests {
				err := redisCache.SetWithoutExpiration("key", test.in)
				Expect(err).NotTo(HaveOccurred())

				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeTrue())
				Expect(err).NotTo(HaveOccurred())
				Expect(test.in).To(Equal(test.out))
			}
		})

		It("should set value of a non-existing key to 0 and then increment it", func() {
			err := redisCache.Increment("key")
			Expect(err).NotTo(HaveOccurred())

			var v int64
			exists, err := redisCache.Get("key", &v)
			Expect(exists).Should(BeTrue())
			Expect(v).To(Equal(int64(1)))
		})

		It("should increment a value after set", func() {
			input := int64(10)

			By("Setting key")

			err := redisCache.Set("key", input)
			Expect(err).NotTo(HaveOccurred())

			By("Incrementing key")

			err = redisCache.Increment("key")
			Expect(err).NotTo(HaveOccurred())

			By("Getting key")

			var output int64
			exists, err := redisCache.Get("key", &output)
			Expect(exists).Should(BeTrue())
			Expect(output).To(Equal(input + 1))
		})

		It("should set value of a non-existing key to 0 and then decrement it", func() {
			err := redisCache.Decrement("key")
			Expect(err).NotTo(HaveOccurred())

			var output int64
			exists, err := redisCache.Get("key", &output)
			Expect(exists).Should(BeTrue())
			Expect(output).To(Equal(int64(-1)))
		})

		It("should decrement a value after set", func() {
			input := int64(10)

			By("Setting key")

			err := redisCache.Set("key", input)
			Expect(err).NotTo(HaveOccurred())

			By("Decrementing key")

			err = redisCache.Decrement("key")
			Expect(err).NotTo(HaveOccurred())

			By("Getting key")

			var output int64
			exists, err := redisCache.Get("key", &output)
			Expect(exists).Should(BeTrue())
			Expect(output).To(Equal((input - 1)))
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

				time.Sleep(2 * conf.Redis.GetExpirationTime())

				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeFalse())
				Expect(err).NotTo(HaveOccurred())
			}
		})

		It("should not delete keys after specified time in config if set without expiration", func() {
			for _, test := range objectTests {
				err := redisCache.SetWithoutExpiration("key", test.in)
				Expect(err).NotTo(HaveOccurred())

				time.Sleep(2 * conf.Redis.GetExpirationTime())

				exists, err := redisCache.Get("key", test.out)
				Expect(exists).Should(BeTrue())
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

		It("should set values without expiration without error", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					err := redisCache.SetWithFieldsWithoutExpiration(fields, test.in)
					Expect(err).NotTo(HaveOccurred())
				}
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

		It("should find elements when trying to get after set without expiration", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					err := redisCache.SetWithFieldsWithoutExpiration(fields, test.in)
					Expect(err).NotTo(HaveOccurred())

					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeTrue())
					Expect(err).NotTo(HaveOccurred())
					Expect(test.in).To(Equal(test.out))
				}
			}
		})

		It("should set value of a non-existing key to 0 and then increment it", func() {
			for _, fields := range fieldsTests {
				err := redisCache.IncrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				var output int64
				exists, err := redisCache.GetWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect(output).To(Equal(int64(1)))
			}
		})

		It("should increment a value after set", func() {
			for _, fields := range fieldsTests {
				input := int64(10)

				By("Setting key")

				err := redisCache.SetWithFields(fields, input)
				Expect(err).NotTo(HaveOccurred())

				By("Incrementing key")

				err = redisCache.IncrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				By("Getting key")

				var output int64
				exists, err := redisCache.GetWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect((output)).To(Equal(input + 1))
			}
		})

		It("should set value of a non-existing key to 0 and then decrement it", func() {
			for _, fields := range fieldsTests {
				err := redisCache.DecrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				var output int64
				exists, err := redisCache.GetWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect(output).To(Equal(int64(-1)))
			}
		})

		It("should decrement a value after set", func() {
			for _, fields := range fieldsTests {
				input := int64(10)

				By("Setting key")

				err := redisCache.SetWithFields(fields, input)
				Expect(err).NotTo(HaveOccurred())

				By("Decrementing key")

				err = redisCache.DecrementWithFields(fields)
				Expect(err).NotTo(HaveOccurred())

				By("Getting key")

				var output int64
				exists, err := redisCache.GetWithFields(fields, &output)
				Expect(exists).Should(BeTrue())
				Expect(output).To(Equal(input - 1))
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

					time.Sleep(2 * conf.Redis.GetExpirationTime())

					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeFalse())
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})

		It("should not delete keys after specified time in config if set without expiration", func() {
			for _, test := range objectTests {
				for _, fields := range fieldsTests {
					err := redisCache.SetWithFieldsWithoutExpiration(fields, test.in)
					Expect(err).NotTo(HaveOccurred())

					time.Sleep(2 * conf.Redis.GetExpirationTime())

					exists, err := redisCache.GetWithFields(fields, test.out)
					Expect(exists).Should(BeTrue())
					Expect(err).NotTo(HaveOccurred())
				}
			}
		})
	})
})
