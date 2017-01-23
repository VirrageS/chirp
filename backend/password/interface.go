package password

type PasswordManager interface {
	HashPassword(password string) (string, error)
	GenerateRandomPassword() (string, error)
	ValidatePassword(password, hashedPassword string) bool
}
