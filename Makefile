.PHONY: db-create db-destroy backend test db-insert db-reset

db-create:
	docker compose up -d

db-destroy:
	docker compose down

backend:
	go run ./cmd/api

test:
	go test ./...

db-insert:
	docker compose up -d db
	docker compose exec -T db psql -U pulseforge -d pulseforge -v ON_ERROR_STOP=1 < migrations/init.up.sql

db-reset:
	docker compose exec -T db psql -U pulseforge -d pulseforge -v ON_ERROR_STOP=1 < migrations/init.down.sql
