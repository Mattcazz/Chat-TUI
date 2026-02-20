package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/Mattcazz/Chat-TUI/server/db"
	"github.com/Mattcazz/Chat-TUI/server/resources/file"

	"github.com/Mattcazz/Chat-TUI/server/resources/chat"
	"github.com/Mattcazz/Chat-TUI/server/resources/user"
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

	txManager := db.NewTxManager(a.db)

	userStore := user.NewUserStore(a.db)
	contactStore := user.NewContactStore(a.db)
	challengeStore := user.NewChallengeStore(a.db)
	userService := user.NewService(userStore, contactStore, challengeStore, txManager)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(r)

	fileStore := file.NewFileStore(a.db)
	fileService := file.NewService(fileStore, txManager)
	fileHandler := file.NewHandler(fileService)
	fileHandler.RegisterRoutes(r)

	conversationStore := chat.NewConversationStore(a.db)
	conversationService := chat.NewService(conversationStore, txManager)
	conversationHandler := chat.NewHandler(conversationService)
	conversationHandler.RegisterRoutes(r)

	log.Println("Listening on address ", a.addr)

	return http.ListenAndServe(a.addr, r)
}
