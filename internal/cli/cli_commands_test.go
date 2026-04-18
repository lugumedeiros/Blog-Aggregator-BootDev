package cli

import (
	"testing"
)

func TestLogin(t *testing.T) {
	err := Execute("login", []string{"test_user"})
	if err != nil {
		t.Errorf("Failed to Login: %v", "test_user")
	}

	err = Execute("login", []string{"test_user", "otherarg"})
	if err == nil {
		t.Errorf("Failed to raise Error with wrong arguments")
	}
}
