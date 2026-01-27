package main

import (
	"log"

	"go.mod/db"
)

func main() {

	addr := "localhost:8080" // momentary

	db := db.ConnectDB()
	api := NewApiServer(addr, db)

	if err := api.Run(); err != nil {
		log.Fatal("Error running server")
	}

	log.Println("Listening on address ", addr)
}
