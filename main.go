package main

import (
	"backend/controllers"
	"log"
)

func main() {
	engine, _, err := controllers.SetupServer("config.toml", "users.json")
	if err != nil {
		panic(err)
	}
	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}
}
