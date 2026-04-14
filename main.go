package main

import (
	"fmt"
	internal_cfg "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/config"
)

const name_default = "default"

// 1. Read config file
// 2. Set user to a "name"
// 3. Read config again and print contents to cli
func main() {
	// fmt.Printf(config.GetConfigFilePath())
	config_s, err_config := internal_cfg.Read()
	if err_config != nil {
		panic(err_config)
	}
	fmt.Printf("Stored username: %v\n", config_s.Current_user_name)
	internal_cfg.SetUser(name_default)
	config_s, err_config = internal_cfg.Read()
	if err_config != nil {
		panic(err_config)
	}
	fmt.Printf("New username: %v\n", config_s.Current_user_name)
	fmt.Printf("db_url: '%v', user: '%v'\n", config_s.Db_url, config_s.Current_user_name)
}
