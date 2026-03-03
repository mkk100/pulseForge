.PHONY: up down run test

db-create:
	docker compose up -d

db-destroy:
	docker compose down

backend:
	go run ./cmd/api

test:
	go test ./...