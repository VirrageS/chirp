package password

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/VirrageS/chirp/backend/config"
)

type bcryptManager struct {
	randomPasswordLength int
}

func NewBcryptManager(config config.PasswordConfigProvider) Manager {
	randomPasswordLength := config.GetRandomPasswordLength()
	return &bcryptManager{
		randomPasswordLength: randomPasswordLength,
	}
}

func (m *bcryptManager) HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("Error generating hash from password.")
	}

	return string(passwordHash), err
}

func (m *bcryptManager) GenerateRandomPassword() (string, error) {
	return generateRandomString(m.randomPasswordLength)
}

func (m *bcryptManager) ValidatePassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	if err != nil {
		// this is a real error, not a wrong password error
		log.WithError(err).Error("Error validating password.")
		return false
	}

	return true
}
