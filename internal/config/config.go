package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const TESTING_ENV = "GATOR_CONFIG_PATH"
const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	//Env override
	path := os.Getenv("GATOR_CONFIG_PATH")
	if path != "" {
		return path, nil
	}

	//Fallback to working dir
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, CONFIG_FILE_NAME), nil
}

func Read() (Config, error) {
	// Get path
	var config Config
	config_file_path, read_err := getConfigFilePath()
	if read_err != nil {
		return config, read_err
	}

	// Read json
	byteValue, err_path := os.ReadFile(config_file_path)
	if err_path != nil {
		return config, err_path
	}

	// Decode json to struct
	err_unmarshal := json.Unmarshal(byteValue, &config)
	if err_unmarshal != nil {
		return config, err_unmarshal
	}

	return config, nil
}

func SetUser(new_name string) error {
	file_path, err_get_path := getConfigFilePath()
	if err_get_path != nil {
		return err_get_path
	}

	// get struct
	config, err := Read()
	if err != nil {
		return err
	}

	config.Current_user_name = new_name
	bytes, err_marsh := json.Marshal(config)
	if err_marsh != nil {
		return err
	}
	return os.WriteFile(file_path, bytes, 0644)
}
