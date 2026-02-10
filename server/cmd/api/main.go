package main

import (
	"log"

	"github.com/Mattcazz/Chat-TUI/server/db"
)

func main() {

	addr := ":8080" // momentary

	db := db.ConnectDB()
	api := NewApiServer(addr, db)

	if err := api.Run(); err != nil {
		log.Fatal("Error running server")
	}

}
