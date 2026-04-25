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

func GetUserByName(username string) error {
	_, err := AppDB.Querie.GetUser(AppDB.Ctx, username)
	return err
}

func RemoveUserByName(username string) error {
	err := AppDB.Querie.RemoveUser(AppDB.Ctx, username)
	return err
}
