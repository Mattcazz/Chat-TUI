package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"go.mod/resources/file"
	"go.mod/resources/msg"
	"go.mod/resources/user"
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

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	userStore := user.NewUserStore(a.db)
	contactStore := user.NewContactStore(a.db)
	userService := user.NewService(userStore, contactStore)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	fileStore := file.NewFileStore(a.db)
	fileService := file.NewService(fileStore)
	fileHandler := file.NewHandler(fileService)
	fileHandler.RegisterRoutes(r)

	msgStore := msg.NewMsgStore(a.db)
	msgService := msg.NewService(msgStore)
	msgHandler := msg.NewHandler(msgService)
	msgHandler.RegisterRoutes(r)

	log.Println("Listening on address ", a.addr)

	return http.ListenAndServe(a.addr, r)
}
