# Makefile for building, running, and hot reloading main.go with air

include .env.test
.PHONY: build run hot-reload

prepare:
	./scripts/install_deps.sh

templ:
	@templ generate --watch --proxy="http://localhost:8080"

server:
	./bin/air

tailwind:
	./bin/tailwindcss -i ./frontend/static/css/input.css -o ./frontend/static/css/output.css --watch

dev:
	make -j5 tailwind server templ

build:
	go build -o main backend/main.go

run:
	go run backend/main.go

test:
	go test -v ./...
