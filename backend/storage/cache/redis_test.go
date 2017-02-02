package cache

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/VirrageS/chirp/backend/config"
	"github.com/VirrageS/chirp/backend/model"
)

var _ = Describe("RedisCache", func() {
	var (
		conf       *config.Configuration = config.New()
		redisCache Accessor              = NewRedisCache(conf.Redis)

		valuesTests []struct {
			in  interface{}
			out interface{}
		}

		key Key
	)

	BeforeEach(func() {
		valuesTests = []struct {
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

		key = Key{"key", 12, "hello", -1}
	})

	AfterEach(func() {
		redisCache.Flush()
	})

	Describe("test standard functions", func() {
		var (
			entriesInTests  []Entry
			entriesOutTests []Entry
			keysTests       []Key
		)

		BeforeEach(func() {
			// TODO: is there better way to clean variables?
			entriesInTests = nil
			entriesOutTests = nil
			keysTests = nil

			for i, test := range valuesTests {
				key := Key{"somekey", -1, i}
				entriesInTests = append(entriesInTests, Entry{key, test.in})
				entriesOutTests = append(entriesOutTests, Entry{key, test.out})
				keysTests = append(keysTests, key)
			}
		})

		It("should set value", func() {
			err := redisCache.Set(entriesInTests...)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get single value", func() {
			err := redisCache.Set(entriesInTests...)
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < len(entriesInTests); i++ {
				exists, err := redisCache.GetSingle(entriesOutTests[i].Key, entriesOutTests[i].Value)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).Should(BeTrue())
				Expect(entriesInTests[i].Value).To(Equal(entriesOutTests[i].Value))
			}
		})

		It("should not get single value when key is not set", func() {
			for i := 0; i < len(entriesInTests); i++ {
				exists, err := redisCache.GetSingle(entriesOutTests[i].Key, entriesOutTests[i].Value)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).Should(BeFalse())
				Expect(entriesInTests[i].Value).NotTo(Equal(entriesOutTests[i].Value))
			}
		})

		It("should get values", func() {
			err := redisCache.Set(entriesInTests...)
			Expect(err).NotTo(HaveOccurred())

			existing, err := redisCache.Get(entriesOutTests...)
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < len(entriesInTests); i++ {
				Expect(existing[i]).To(BeTrue())
				Expect(entriesInTests[i].Value).To(Equal(entriesOutTests[i].Value))
			}
		})

		It("should not get values when key is not set", func() {
			existing, err := redisCache.Get(entriesOutTests...)
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < len(entriesInTests); i++ {
				Expect(existing[i]).To(BeFalse())
				Expect(entriesInTests[i].Value).NotTo(Equal(entriesOutTests[i].Value))
			}
		})

		It("should delete keys", func() {
			err := redisCache.Set(entriesInTests...)
			Expect(err).NotTo(HaveOccurred())

			err = redisCache.Delete(keysTests...)
			Expect(err).NotTo(HaveOccurred())

			existing, err := redisCache.Get(entriesOutTests...)
			Expect(err).NotTo(HaveOccurred())

			for i := 0; i < len(entriesInTests); i++ {
				Expect(existing[i]).To(BeFalse())
				Expect(entriesInTests[i].Value).NotTo(Equal(entriesOutTests[i].Value))
			}
		})
	})

	Describe("test integer functions", func() {
		var (
			integerTests []struct {
				key   Key
				value int64
			}
			keysTests []Key
		)

		BeforeEach(func() {
			integerTests = []struct {
				key   Key
				value int64
			}{
				{Key{"lol", -1, "najs"}, -1},
				{Key{"lol", -1, "najs1"}, 1},
				{Key{123, "9", -1, "najs"}, 1234124},
				{Key{"lol", "sdfsdf", "nsdfajs"}, -1124124141441},
				{Key{"losdfsl", -1, "nssfsajs"}, 0},
			}

			keysTests = []Key{}
			for _, test := range integerTests {
				keysTests = append(keysTests, test.key)
			}
		})

		Describe("tests with existsing keys", func() {
			BeforeEach(func() {
				for _, test := range integerTests {
					err := redisCache.Set(Entry{test.key, test.value})
					Expect(err).NotTo(HaveOccurred())
				}
			})

			It("should increase", func() {
				var value int64

				err := redisCache.Incr(keysTests...)
				Expect(err).NotTo(HaveOccurred())

				for _, test := range integerTests {
					exists, err := redisCache.GetSingle(test.key, &value)
					Expect(err).NotTo(HaveOccurred())
					Expect(exists).Should(BeTrue())
					Expect(value).To(Equal(test.value + 1))
				}
			})

			It("should decrease", func() {
				var value int64

				err := redisCache.Decr(keysTests...)
				Expect(err).NotTo(HaveOccurred())

				for _, test := range integerTests {
					exists, err := redisCache.GetSingle(test.key, &value)
					Expect(err).NotTo(HaveOccurred())
					Expect(exists).Should(BeTrue())
					Expect(value).To(Equal(test.value - 1))
				}
			})
		})

		Describe("tests with non-existing keys", func() {
			It("should increase non-existing keys setting their value to 1", func() {
				var value int64

				err := redisCache.Incr(keysTests...)
				Expect(err).NotTo(HaveOccurred())

				for _, key := range keysTests {
					exists, err := redisCache.GetSingle(key, &value)
					Expect(err).NotTo(HaveOccurred())
					Expect(exists).Should(BeTrue())
					Expect(value).To(Equal(int64(1)))
				}
			})

			It("should decrease non-existing keys setting their value to -1", func() {
				var value int64

				err := redisCache.Decr(keysTests...)
				Expect(err).NotTo(HaveOccurred())

				for _, key := range keysTests {
					exists, err := redisCache.GetSingle(key, &value)
					Expect(err).NotTo(HaveOccurred())
					Expect(exists).Should(BeTrue())
					Expect(value).To(Equal(int64(-1)))
				}
			})
		})
	})

	Describe("test set functions", func() {
		var (
			key            Key
			values         Values
			expectedValues []int64
		)

		BeforeEach(func() {
			key = Key{"test", "suppoer", -1, "function"}
			values = []int64{-1, 24525, 1231, 0, 12312, 1, 1, 1, 1, 1, 1}
			expectedValues = []int64{24525, 1231, 1, 12312, 0, -1} // set of values
		})

		It("should add members to set", func() {
			err := redisCache.SAdd(key, values)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should get set members", func() {
			var members []int64

			err := redisCache.SAdd(key, values)
			Expect(err).NotTo(HaveOccurred())

			exists, err := redisCache.SMembers(key, &members)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
			Expect(members).To(ConsistOf(expectedValues))
		})

		It("should not get set members when key is not set", func() {
			var members []int64

			exists, err := redisCache.SMembers(key, &members)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		It("shoud remove members from set", func() {
			var members []int64

			err := redisCache.SAdd(key, values)
			Expect(err).NotTo(HaveOccurred())

			err = redisCache.SRemove(key, values)
			Expect(err).NotTo(HaveOccurred())

			exists, err := redisCache.SMembers(key, &members)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})
})
