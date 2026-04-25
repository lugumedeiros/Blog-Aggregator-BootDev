package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	cfg "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/config"
	db "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/database"
)

type commands struct {
	name        string
	description string
	callback    func([]string) error
}

var commands_cli map[string]commands

func init() {
	commands_cli = map[string]commands{
		"login": {
			name:        "login",
			description: "login {username} - Used to login with a username",
			callback:    login,
		},
		"exit": {
			name:        "exit",
			description: "exit - Used to end current application",
			callback:    exit,
		},
		"help": {
			name:        "help",
			description: "help - Used to show every available command",
			callback:    help,
		},
		"register": {
			name:        "register",
			description: "register {username} - Used to register a new login username",
			callback:    register,
		},
		"unregister": {
			name:        "unregister",
			description: "unregister {username} - Used to remove a login username from database",
			callback:    unregister,
		},
	}
}

func login(args []string) error {
	if len(args) != 1 {
		return getErrorArgsQntd(1, len(args))
	}
	username := args[0]

	// Check DB
	errdb := db.GetUserByName(username)
	if errdb != nil {
		if strings.Contains(errdb.Error(), "no rows in result") {
			return fmt.Errorf("User '%v' not registered.", username)
		}
		return errdb
	}
	cfg.SetUser(username)
	fmt.Printf("User set to '%v'\n", username)
	return nil
}

func exit(args []string) error {
	os.Exit(1)
	return nil
}

func help(args []string) error {
	if len(args) != 0 {
		return getErrorArgsQntd(0, len(args))
	}

	for _, command_s := range commands_cli {
		fmt.Printf("Command: %v\nDescription: %v\n\n", command_s.name, command_s.description)
	}
	return nil
}

func register(args []string) error {
	if len(args) != 1 {
		return getErrorArgsQntd(1, len(args))
	}
	username := args[0]

	err_login := db.GetUserByName(username)
	if err_login == nil {
		return errors.New("Unable to register new user")
	}

	query_err := db.CreateUser(username)
	if query_err == nil {
		fmt.Printf("User '%v' registered\n", username)
		return login([]string{username})
	}

	if strings.Contains(query_err.Error(), "duplicate key value violates unique") {
		return fmt.Errorf("User '%v' already registered.", username)
	}
	return query_err
}

func unregister(args []string) error {
	if len(args) != 1 {
		return getErrorArgsQntd(1, len(args))
	}
	username := args[0]

	err_login := db.GetUserByName(username)
	if err_login != nil {
		return errors.New("Unable to unregister user")
	}

	query_err := db.RemoveUserByName(username)
	if query_err == nil {
		fmt.Printf("User '%v' unregistered\n", username)
		return nil
	}
	return query_err
}

func Execute(command_arg string, args []string) error {
	command_struct, ok := commands_cli[strings.ToLower(command_arg)]
	if !ok {
		return errors.New("Nonexistent command. Try 'help'.\n")
	}

	execution_err := command_struct.callback(args)
	return execution_err
}

// -------------------------------- //

func getErrorArgsQntd(expected int, current int) error {
	return fmt.Errorf("Expected '%v' argument but found '%v'.\n", expected, current)
}
