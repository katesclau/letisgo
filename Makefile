# Makefile for building, running, and hot reloading main.go with air

include .env.test
.PHONY: build run hot-reload

build:
	go build -o main backend/main.go

run:
	go run backend/main.go

dev:
	./bin/air --build.cmd "go build -o ./bin/api backend/main.go" --build.bin "./bin/api"

test:
	go test -v ./...
