package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	internal_cfg "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/config"
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
	}
}

func login(args []string) error {
	if len(args) != 1 {
		return getErrorArgsQntd(1, len(args))
	}
	user := args[0]
	internal_cfg.SetUser(user)
	fmt.Printf("User set to '%v'\n", user)
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
