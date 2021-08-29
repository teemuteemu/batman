BIN_DIR=bin
BIN_NAME=batman

all: test build

test:
	go test -v ./...
.PHONY: test

lint:
	go mod download
	golangci-lint run --allow-parallel-runners
.PHONY: lint

build:
	go build -o $(BIN_DIR)/$(BIN_NAME) cmd/main.go
.PHONY: build

install:
	go build -o $(GOPATH)/bin/$(BIN_NAME) cmd/main.go
.PHONY: install

uninstall:
	rm $(GOPATH)/bin/$(BIN_NAME)
.PHONY: uninstall

clean:
	rm -rf $(BIN_DIR)
.PHONY: clean
