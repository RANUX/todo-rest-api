.PHONY: build
build:
		go build -v ./cmd/apiserver

.PHONY: migrate
migrate:
	migrate -source file://internal/app/postgres/migrations \
		-database postgres://postgres:39dkj29d@127.0.0.1:5432/todos_db?sslmode=disable up

.PHONY: rollback
rollback:
	migrate -source file://internal/app/postgres/migrations \
		-database postgres://postgres:39dkj29d@127.0.0.1:5432/todos_db?sslmode=disable down

.PHONY: drop
drop:
	migrate -source file://internal/app/postgres/migrations \
		-database postgres://postgres:39dkj29d@127.0.0.1:5432/todos_db?sslmode=disable drop


.PHONY: migrations
migrations:
		@read -p "Enter migration name: " name; \
			migrate create -ext sql -dir ./internal/app/postgres/migrations $$name


.PHONY: sqlc
sqlc:
	sqlc generate

.DEFAULT_GOAL := build