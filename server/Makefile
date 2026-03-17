run: 
	@go run cmd/api/*.go
	
migration: 
	@migrate create -ext sql -dir cmd/migrate/migrations

migrate-up: 
	@go run cmd/migrate/main.go up

migrate-down: 
	@go run cmd/migrate/main.go down

test:
	@echo "Running server tests..."
	@go test ./... -race 2>&1 | grep -E "^(ok|FAIL|--- FAIL)"

test-v:
	@echo "Running server tests with verbose output..."
	@go test ./... -v -rave
