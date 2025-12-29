# ----------------------
# Project Variables
# ----------------------
DB_DSN=postgres://postgres:secret@localhost:5432/rssagg?sslmode=disable
MIGRATIONS_DIR=./internal/adapters/database/migrations/
BINARY=bin/rssagg
CMD_DIR=./cmd

# ----------------------
# Build / Run / Test
# ----------------------
build:
	@go build -mod=vendor -o $(BINARY) $(CMD_DIR)

run: build
	@./$(BINARY)

test:
	@go test -mod=vendor -v ./...

# ----------------------
# Database Migrations
# ----------------------
.PHONY: migrate rollback status down-to

migrate:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

rollback:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

status:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

down-to:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down-to $(VERSION)

# ----------------------
# Optional: fresh (drop all + migrate)
# ----------------------
.PHONY: fresh
fresh:
	@goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset
	@make migrate
