# Makefile for building, running, and hot reloading main.go with air

include .env.test
.PHONY: build run hot-reload

build:
	go build -o main main.go

run:
	go run main.go

dev:
	./bin/air

test:
	go test -v ./...
