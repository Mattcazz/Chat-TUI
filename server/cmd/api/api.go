package main

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.mod/service/file"
	"go.mod/service/msg"
	"go.mod/service/user"
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

	userStore := user.NewStore(a.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(r)

	fileStore := file.NewStore(a.db)
	fileHandler := file.NewHandler(fileStore)
	fileHandler.RegisterRoutes(r)

	msgStore := msg.NewStore(a.db)
	msgHandler := msg.NewHandler(msgStore)
	msgHandler.RegisterRoutes(r)

	return http.ListenAndServe(a.addr, r)
}
