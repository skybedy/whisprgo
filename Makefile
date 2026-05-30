APP_NAME := whisprgo
CMD_PATH := ./cmd/whisprgo
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

.PHONY: build install uninstall test fmt dev

build:
	go build -o $(APP_NAME) $(CMD_PATH)

install:
	install -d "$(BINDIR)"
	go build -o "$(BINDIR)/$(APP_NAME)" $(CMD_PATH)

uninstall:
	rm -f "$(BINDIR)/$(APP_NAME)"

fmt:
	go fmt ./...

test:
	go test ./...

dev:
	go mod tidy
	go fmt ./...
	go test ./...
