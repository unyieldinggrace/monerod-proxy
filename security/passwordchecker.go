package security

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/ini.v1"
)

type IPasswordChecker interface {
	CheckAdminPassword(password string) bool
}

type IPasswordHashGenerator interface {
	GeneratePasswordHash(password string) string
}

type PasswordChecker struct {
	AdminPasswordHash string
}

func CreatePasswordCheckerFromConfig(cfg *ini.File) *PasswordChecker {
	adminPasswordHash := cfg.Section("").Key("admin_password_hash").Value()

	passwordChecker := &PasswordChecker{
		AdminPasswordHash: adminPasswordHash,
	}

	return passwordChecker
}

func (passwordChecker *PasswordChecker) CheckAdminPassword(password string) bool {
	return passwordChecker.passwordMatchesHash(password, passwordChecker.AdminPasswordHash)
}

func (passwordChecker *PasswordChecker) passwordMatchesHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (passwordChecker *PasswordChecker) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
