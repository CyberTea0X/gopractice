package controllers

import (
	"backend/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const BasePath = "/api"

func SetupServer(configPath string, usersJson string) (*gin.Engine, *sql.DB, error) {
	config, err := models.ParseConfig(configPath)
	if err != nil {
		return nil, nil, err
	}
	db, err := models.SetupDatabase(&config.Database)
	if err != nil {
		return nil, nil, err
	}
	for err := db.Ping(); err != nil; err = db.Ping() {
		log.Println(err)
		time.Sleep(time.Second)
		log.Println("Database ping failed, retrying..")
	}
	log.Println("Successfully connected to the database")
	err = models.MigrateDatabase(db)
	if err != nil {
		panic(err)
	}

	if err := models.MigrateDatabase(db); err != nil {
		return nil, nil, errors.Join(errors.New("failed to automigrate database"), err)
	}
	users, err := models.ParseUsersFromJson(usersJson)
	config.SmsCodes.Lifespan()
	if err != nil {
		return nil, nil, err
	}
	userIds := make([]int64, len(users))
	for i, user := range users {
		id, err := models.CreateOrUpdateUser(db, user.Phone, user.Roles)
		userIds[i] = id
		if err != nil {
			return nil, nil, err
		}
	}
	fmt.Println("UserIds:", userIds)
	engine := gin.New()
	api := engine.Group(BasePath)
	api.GET("/login", Login(
		db,
		config.Tokens.Refresh.Lifespan(),
		config.Tokens.Access.Lifespan(),
		config.SmsCodes.Lifespan(),
		config.Tokens.Refresh.Secret,
		config.Tokens.Access.Secret,
	))
	api.GET("/refresh", Refresh(
		db,
		config.Tokens.Refresh.Lifespan(),
		config.Tokens.Access.Lifespan(),
		config.Tokens.Refresh.Secret,
		config.Tokens.Access.Secret,
	))
	api.GET("/auth", Auth(config.Tokens.Access.Secret))
	api.GET("/sendcode", SendCode(db))
	return engine, db, nil
}
