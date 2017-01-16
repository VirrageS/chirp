package password

type Manager interface {
	HashPassword(password string) (string, error)
	GenerateRandomPassword() (string, error)
	ValidatePassword(password, hashedPassword string) bool
}
