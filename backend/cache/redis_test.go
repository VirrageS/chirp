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

var _ = Describe("RedisCache", func() {
	var (
		redisCache CacheProvider = NewRedisCache("6379", &mockCacheConfigProvider{})

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
