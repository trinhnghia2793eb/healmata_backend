include .env
export

export GOOSE_DBSTRING=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
# Goose commands
migrate-up:
		goose up

migrate-down:
		goose down

migrate-status:
		goose status

migrate-redo:
		goose redo
		
migrate-reset:
		goose reset

# Run database constraint integration tests (requires DB_* env vars)
test-db:
		go test ./internal/db/migrations/... -v -count=1