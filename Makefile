# Makefile for building, running, and hot reloading main.go with air

include .env.test
.PHONY: build run hot-reload

prepare:
	{ \
        curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s;    \
        curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64; \
        chmod +x tailwindcss-macos-arm64; \
        mv tailwindcss-macos-arm64 ./bin/tailwindcss; \
    }

build:
	go build -o main backend/main.go

run:
	go run backend/main.go

dev:
	./bin/air

test:
	go test -v ./...
