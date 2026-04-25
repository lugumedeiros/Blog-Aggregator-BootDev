package cli

import (
	"testing"
	db "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/database"
)

func TestRegister(t *testing.T) {
	t.Setenv("GATOR_CONFIG_PATH", "../../.gatorconfig.json")
	db.InitDB()

	user := []string{"test_user"}

	//TEST 1
	err := Execute("login", user)
	if err == nil {
		t.Errorf("ERR: Has logged to a invalid user.")
	}
	
	// TEST 2
	_ = Execute("register", user)
	err = Execute("login", user)
	if err != nil {
		t.Errorf("ERR: user was not registered.")
	}

	//TEST 3
	_ = Execute("unregister", user)
	err = Execute("login", user)
	if err == nil {
		t.Errorf("ERR: user was not unregistered.")
	}
}
