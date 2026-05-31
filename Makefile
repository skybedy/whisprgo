APP_NAME := whispergo
CMD_PATH := ./cmd/whispergo
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
DIST_DIR ?= dist
VERSION ?= $(shell cat VERSION 2>/dev/null || echo dev)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X 'whispergo/internal/buildinfo.Version=$(VERSION)' -X 'whispergo/internal/buildinfo.Commit=$(COMMIT)' -X 'whispergo/internal/buildinfo.Date=$(DATE)'

.PHONY: build install uninstall test fmt dev clean release release-linux-amd64 release-linux-arm64

build:
	go build -ldflags "$(LDFLAGS)" -o $(APP_NAME) $(CMD_PATH)

install:
	install -d "$(BINDIR)"
	go build -ldflags "$(LDFLAGS)" -o "$(BINDIR)/$(APP_NAME)" $(CMD_PATH)

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

clean:
	rm -rf $(DIST_DIR) $(APP_NAME)

release-linux-amd64:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 $(CMD_PATH)

release-linux-arm64:
	mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 $(CMD_PATH)

release: release-linux-amd64 release-linux-arm64
