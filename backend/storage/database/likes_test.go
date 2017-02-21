package database

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/VirrageS/chirp/backend/config"
)

var _ = Describe("Likes", func() {
	var (
		conf *config.Configuration = config.New()
		db                         = NewPostgresDatabase(conf.Postgres)
	)

	BeforeEach(func() {})

	AfterEach(func() {
		// HACK: this is hack since TRUNCATE can execute up to 1s... whereas this ~5ms
		db.Exec(`DELETE FROM users; DELETE FROM likes;`)
	})

	It("should do something", func() {
		Expect(true).To(BeTrue())
	})
})
