package main

import (
	"fmt"
	"os"

	cli_commands "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/cli"
	db "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/database"
)

// const name_default = "default"

// 1. Read config file
// 2. Set user to a "name"
// 3. Read config again and print contents to cli
func main() {
	db.InitDB()
	if len(os.Args) < 2 {
		fmt.Printf("No command provided\n")
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	err := cli_commands.Execute(cmd, args)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
}
