.PHONY: setup dev devFrontend devBackend test testFrontend testBackend lint build migrateUp migrateDown simulate

setup:
	cd frontend && npm ci
	cd backend && go mod download

dev:
	docker compose up --build

devFrontend:
	cd frontend && npm run dev

devBackend:
	cd backend && go run ./cmd/server

test: testFrontend testBackend

testFrontend:
	cd frontend && npm test

testBackend:
	cd backend && go test ./...

lint:
	cd frontend && npm run lint && npm run formatCheck && npm run typecheck
	cd backend && go vet ./...

build:
	cd frontend && npm run build
	cd backend && go build ./cmd/server ./cmd/migrate ./cmd/simulate

migrateUp:
	cd backend && go run ./cmd/migrate up

migrateDown:
	cd backend && go run ./cmd/migrate down

simulate:
	cd backend && go run ./cmd/simulate
