package database

import (
	"context"
	"database/sql"

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