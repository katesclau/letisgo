# Makefile for building, running, and hot reloading main.go with air

include .env.test
.PHONY: build run hot-reload

prepare:
	./scripts/install_deps.sh

templ:
	@templ generate --watch --proxy="http://localhost:8080" --open-browser=false

server:
	./bin/air

tailwind:
	npm run watch

infra:
	docker compose up -d --wait
	cd backend/terraform/service && tflocal init && tflocal apply -auto-approve && cd -

infra-down:
	docker compose down

dev: infra
	REDIS_HOST=localhost:6379 make -j5 tailwind templ server

build:
	go build -o main backend/main.go

run:
	go run backend/main.go

test:
	go test -v ./...
