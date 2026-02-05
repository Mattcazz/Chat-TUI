module github.com/Mattcazz/Chat-TUI/server

go 1.25.1

require (
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/chi/v5 v5.2.4
	github.com/golang-migrate/migrate/v4 v4.19.1
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.11.1
)

require (
	github.com/stretchr/testify v1.11.1 // indirect
	golang.org/x/sys v0.38.0 // indirect
)

require github.com/golang-jwt/jwt/v5 v5.3.1

require (
	github.com/Mattcazz/Chat-TUI/pkg v0.0.0
	golang.org/x/crypto v0.45.0
)

replace github.com/Mattcazz/Chat-TUI/pkg => ../pkg
