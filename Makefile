include .env
export

.PHONY: dev test lint migrate-up migrate-down collector-run api-run frontend-run

dev:
	docker compose -f docker/docker-compose.yml up --build

test:
	cd backend && go test ./...

lint:
	cd backend && go vet ./...

migrate-up:
	docker compose -f docker/docker-compose.yml exec -T database psql "$$DATABASE_URL" -f /migrations/0001_init.up.sql

migrate-down:
	docker compose -f docker/docker-compose.yml exec -T database psql "$$DATABASE_URL" -f /migrations/0001_init.down.sql

collector-run:
	cd backend && go run ./cmd/collector

api-run:
	cd backend && go run ./main.go

frontend-run:
	cd frontend && npm run dev
