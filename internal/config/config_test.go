package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateUser(t *testing.T) {
	test_name := "Testing_name"
	og_name := "Default_name"

	// New json config file for testing
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.json")
	os.Setenv(TESTING_ENV, path)
	content := fmt.Sprintf(`{"db_url":"None", "current_user_name":"%v"}`, og_name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Errorf("Failed to Create Test json: %v", err)
	}

	// TEST 1
	config, err_read := Read()
	if err_read != nil {
		t.Errorf("Failed to READ: %v", err_read)
	}
	if og_name != config.Current_user_name {
		t.Errorf("1. Expected name '%v' but got '%v' instead", og_name, config.Current_user_name)
	}

	// TEST 2
	SetUser(test_name)
	config, err_read = Read()
	if err_read != nil {
		t.Errorf("Failed to READ: %v", err_read)
	}
	if test_name != config.Current_user_name {
		t.Errorf("2. Expected name '%v' but got '%v' instead", test_name, config.Current_user_name)
	}
}
