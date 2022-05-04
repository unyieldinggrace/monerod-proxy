package security

import (
	"testing"
)

func TestWhenAdminPasswordIsCorrectThenCheckAdminPasswordReturnsTrue(t *testing.T) {
	passwordChecker := getDefaultPasswordChecker(t, "somepassword")

	result := passwordChecker.CheckAdminPassword("somepassword")
	if !result {
		t.Errorf("Password checker should return true for the same password that was used to generate the admin hash.")
	}
}

func TestWhenAdminPasswordIsIncorrectThenCheckAdminPasswordReturnsFalse(t *testing.T) {
	passwordChecker := getDefaultPasswordChecker(t, "somepassword")

	result := passwordChecker.CheckAdminPassword("otherpassword")
	if result {
		t.Errorf("Password checker should return false for a different password than the one that was used to generate the admin hash.")
	}
}

func getDefaultPasswordChecker(t *testing.T, password string) *PasswordChecker {
	passwordChecker := &PasswordChecker{}

	adminHash, err := passwordChecker.GeneratePasswordHash("somepassword")
	t.Log("Admin Hash: ", adminHash)
	if err != nil {
		t.Errorf("Error generating test password")
	}

	passwordChecker.AdminPasswordHash = adminHash
	return passwordChecker
}
