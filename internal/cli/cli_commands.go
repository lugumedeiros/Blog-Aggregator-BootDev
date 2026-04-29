package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	cfg "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/config"
	db "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/database"
	rss "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/rss"
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
		"reset" : {
			name: "reset",
			description: "reset - Used to clear the entire database",
			callback: reset,
		},
		"users" : {
			name: "users",
			description: "users - Get all users registered",
			callback: users,
		},
		"current" : {
			name: "current",
			description: "current - Get logged user",
			callback: current,
		},
		"agg" : {
			name: "agg",
			description: "agg {url} - IDK",
			callback: agg,
		},
		"addfeed" : {
			name: "addfeed",
			description: "addfeed {title} {url} - Add a feed to current user",
			callback: addfeed,
		},
		"feeds" : {
			name: "feeds",
			description: "feeds - Get all feeds from database",
			callback: feeds,
		},
		"follow" : {
			name: "follow",
			description: "follow {url} - Makes current user follow a feed",
			callback: follow,
		},
		"following" : {
			name: "following",
			description: "following - List all feeds followed by current user",
			callback: following,
		},
		"followingAll": {
			name: "followingAll",
			description: "followingAll - List all feeds with it's followers",
			callback: followingAll,
		},
	}
}

func login(args []string) error {
	if len(args) != 1 {
		return getErrorArgsQntd(1, len(args))
	}
	username := args[0]

	// Check DB
	_, errdb := db.GetUserByName(username)
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

	_, err_login := db.GetUserByName(username)
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

	_, err_login := db.GetUserByName(username)
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

func reset(args []string) error {
	if len(args) != 0 {
		return getErrorArgsQntd(0, len(args))
	}

	query_err := db.ResetUsers()
	if query_err == nil {
		fmt.Printf("Database cleared\n")
		return nil
	}

	return query_err
}

func current(args []string) error{
	username := getCurrentUsername()
	fmt.Printf("Current user: '%v'\n", username)
	return nil
}

func users(args []string) error {
	logged_user := getCurrentUsername()
	users, err_db := db.GetUsers()
	if err_db != nil {
		return err_db
	}
	
	for _, user := range users{
		if user.Name == logged_user {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func agg(args []string) error {
	// if len(args) <= 1 {
	// 	return getErrorArgsQntd(0, len(args))
	// }
	// url := args[0]

	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(url)	
	if err != nil {
		return err
	}
	fmt.Printf("TITLE: %v\nLINK: %v\nDESCRIPTION: %v\n", feed.Channel.Title, feed.Channel.Link, feed.Channel.Description)
	for idx, feed_item := range feed.Channel.Item {
		fmt.Printf("%v. %v\nlink: %v\nDescription: %v\nPubDate: %v\n\n", idx+1, feed_item.Title, feed_item.Link, feed_item.Description, feed_item.PubDate)
	}
	return nil
}

func addfeed(args []string) error {
	if len(args) != 2 {
		return getErrorArgsQntd(2, len(args))
	}
	title := args[0]
	url := args[1]
	user := getCurrentUsername()
	err := db.AddFeed(user, title, url)
	if err != nil {
		return err
	}
	return follow([]string{url})
}

func feeds(args []string) error {
	if len(args) != 0 {
		return getErrorArgsQntd(0, len(args))
	}
	feeds_s, err := db.GetAllFeeds()
	if err != nil {
		return err
	}

	for idx, feed := range feeds_s {
		user, err_us := db.GetUserById(feed.UserID)
		if err_us != nil {
			return err_us
		}
		title := feed.Name
		url := feed.Url
		fmt.Printf("%v. %v\nUSER: %v\nURL: %v\n\n", idx, title, user.Name, url)
	}
	return nil
}

func follow(args []string) error {
	if len(args) != 1 {
		return getErrorArgsQntd(1, len(args))
	}
	url := args[0]
	username := getCurrentUsername()
	return db.Follow(url, username)
}

func followingAll(args []string) error {
	if len(args) != 0 {
		return getErrorArgsQntd(0, len(args))
	}
	relations, err := db.Following()
	if err != nil {
		return err
	}
	for _, relation := range relations{
		fmt.Printf("FEED: %v\nURL: %v\nFOLLOWERS:\n", relation.Feed.Name, relation.Feed.Url)
		for _, user := range relation.Users{
			fmt.Printf(" - %v\n", user.Name)
		}
		fmt.Print("\n")
	}
	return nil
}

func following(args []string) error {
	if len(args) != 0 {
		return getErrorArgsQntd(0, len(args))
	}
	relations, err := db.Following()
	if err != nil {
		return err
	}
	username := getCurrentUsername()

	fmt.Printf("Feeds followed by %v:\n", username)
	for _, relation := range relations{
		for _, user := range relation.Users {
			if user.Name == username {
				fmt.Printf(" - %v\n", relation.Feed.Name)
			}
		}
	}
	return nil
}


// --------------------------------- //

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
	return fmt.Errorf("Expected '%v' argument(s) but found '%v'.\n", expected, current)
}

func getCurrentUsername() string {
	gator, err := cfg.Read()
	if err != nil {
		return ""
	}
	return gator.Current_user_name
}
