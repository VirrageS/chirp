package database

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/VirrageS/chirp/backend/config"
)

var _ = Describe("Follows", func() {
	var (
		conf       *config.Configuration = config.New()
		db                               = NewPostgresDatabase(conf.Postgres)
		followsDAO                       = NewFollowsDAO(db)
	)

	BeforeEach(func() {
		// HACK: this is hack since TRUNCATE can execute up to 1s... whereas this ~5ms
		db.Exec(`DELETE FROM users; DELETE FROM follows;`)
	})

	AfterEach(func() {})

	It("should not follow user when followee and follower does not exists", func() {
		followed, err := followsDAO.FollowUser(1, 2)
		Expect(err).To(HaveOccurred())
		Expect(followed).To(BeFalse())
	})

	It("should not follow user when followee does not exists", func() {
		// TODO create follower
		followed, err := followsDAO.FollowUser(1, 2)
		Expect(err).To(HaveOccurred())
		Expect(followed).To(BeFalse())
	})

	It("should not follow user when followeer does not exists", func() {
		// TODO create followee
		followed, err := followsDAO.FollowUser(1, 2)
		Expect(err).To(HaveOccurred())
		Expect(followed).To(BeFalse())
	})
})
