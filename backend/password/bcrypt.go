package password

import (
	log "github.com/Sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"github.com/VirrageS/chirp/backend/config"
)

type BcryptPasswordManager struct {
	randomPasswordLength int
}

func NewBcryptPasswordManager(config config.PasswordManagerConfig) Manager {
	randomPasswordLength := config.GetRandomPasswordLength()
	return &BcryptPasswordManager{
		randomPasswordLength: randomPasswordLength,
	}
}

func (pm *BcryptPasswordManager) HashPassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("Error generating hash from password.")
	}

	return string(passwordHash), err
}

func (pm *BcryptPasswordManager) GenerateRandomPassword() (string, error) {
	return generateRandomString(pm.randomPasswordLength)
}

func (pm *BcryptPasswordManager) ValidatePassword(password, hashedPassword string) bool {
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
