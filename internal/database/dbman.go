package database

import (
	"context"
	"database/sql"
	// "fmt"

	_ "github.com/lib/pq"
	cfg "github.com/lugumedeiros/Blog-Aggregator-BootDev/internal/config"
)

type appDB struct {
	Db     *sql.DB
	Querie *Queries
	Ctx    context.Context
}

var AppDB appDB

func InitDB() {
	// Open DB state
	config, err := cfg.Read()
	if err != nil {
		panic(err)
	}

	db, err_open := sql.Open(config.App, config.Db_url)
	if err_open != nil {
		panic(err_open)
	}
	AppDB.Db = db
	AppDB.Querie = New(db)
	AppDB.Ctx = context.Background()
}

func CreateUser(username string) error {
	var params CreateUserParams
	params.Name = username
	_, err := AppDB.Querie.CreateUser(AppDB.Ctx, params)
	return err
}

func GetUserByName(username string) (User, error) {
	return AppDB.Querie.GetUser(AppDB.Ctx, username)
	
}

func GetUserById(id int32) (User, error) {
	return AppDB.Querie.GetUserById(AppDB.Ctx, id)
}

func RemoveUserByName(username string) error {
	return AppDB.Querie.RemoveUser(AppDB.Ctx, username)
}

func ResetUsers() error {
	return AppDB.Querie.Reset(AppDB.Ctx)
}

func GetUsers() ([]User, error) {
	return AppDB.Querie.GetAllUsers(AppDB.Ctx)
}

func AddFeed(username string, title string, url string) error {
	var param CreateFeedParams
	user_param, err := AppDB.Querie.GetUser(AppDB.Ctx, username)
	if err != nil{
		return err
	}
	param.Name = title
	param.Url = url
	param.UserID = user_param.ID
	return AppDB.Querie.CreateFeed(AppDB.Ctx, param)
}

func GetAllFeeds() ([]Feed, error){
	return AppDB.Querie.GetFeeds(AppDB.Ctx)
}

func GetFeedByUrl(url string) (Feed, error) {
	return AppDB.Querie.GetFeedByUrl(AppDB.Ctx, url)
}

func Follow(url string, username string) error {
	feed, err := GetFeedByUrl(url)
	if err != nil {
		return err
	}
	
	user, err_u := GetUserByName(username)
	if err_u != nil {
		return err_u
	}

	var relation CreateFeedFollowParams
	relation.FeedID = feed.ID
	relation.UserID = user.ID
	return AppDB.Querie.CreateFeedFollow(AppDB.Ctx, relation)
}

type feedRelation struct {
	Feed Feed
	Users []User
}

func Following() ([]feedRelation, error){
	var feedCollection []feedRelation
	feeds, errf := GetAllFeeds()
	if errf != nil {
		return nil, errf
	}
	for _, feed := range feeds {
		user_ids, err_id := AppDB.Querie.GetUsersByFeed(AppDB.Ctx, feed.ID)
		if err_id != nil {
			return nil, err_id
		}
		
		var feedRel feedRelation
		feedRel.Feed = feed
		for _, user_id := range user_ids {
			//can cache but lazy
			user, err_idq := GetUserById(user_id)
			if err_idq != nil {
				return nil, err_idq
			}
			feedRel.Users = append(feedRel.Users, user)
		}
		feedCollection = append(feedCollection, feedRel)
	}
	return feedCollection, nil
}