BIN_DIR=bin
BIN_NAME=batman

all: test build

test:
	go test -v ./...

build:
	go build -o $(BIN_DIR)/$(BIN_NAME) cmd/main.go

install:
	go build -o $(GOPATH)/bin/$(BIN_NAME) cmd/main.go

uninstall:
	rm $(GOPATH)/bin/$(BIN_NAME)

clean:
	rm -rf $(BIN_DIR)
