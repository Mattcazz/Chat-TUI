package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.mod/service/contact"
	"go.mod/service/file"
	"go.mod/service/msg"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (a *APIServer) Run() error {
	r := chi.NewRouter()

	contactStore := contact.NewStore(a.db)
	contactHandler := contact.NewHandler(contactStore)
	contactHandler.RegisterRoutes(r)

	fileStore := file.NewStore(a.db)
	fileHandler := file.NewHandler(fileStore)
	fileHandler.RegisterRoutes(r)

	msgStore := msg.NewStore(a.db)
	msgHandler := msg.NewHandler(msgStore)
	msgHandler.RegisterRoutes(r)

	log.Println("Listening on port ", a.addr)
	return http.ListenAndServe(a.addr, r)
}
