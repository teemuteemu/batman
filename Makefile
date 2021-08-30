BIN_DIR=bin
BIN_NAME=batman
VERSION=$(shell git describe --tags --abbrev=0)

all: test build

test:
	go test -v ./...
.PHONY: test

lint:
	golangci-lint run --allow-parallel-runners
.PHONY: lint

build:
	CGO_ENABLED=0 go build -ldflags "-X main.version=$(VERSION)" -o $(BIN_NAME) main.go
.PHONY: build

install:
	go install
.PHONY: install

uninstall:
	rm $(GOPATH)/bin/$(BIN_NAME)
.PHONY: uninstall

clean:
	rm $(BIN_NAME)
.PHONY: clean
