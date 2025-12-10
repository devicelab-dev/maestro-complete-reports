BINARY_NAME=maestro-complete-reports
CMD_PATH=./cmd/maestro-complete-reports
BIN_DIR=bin

.PHONY: all clean build build-all darwin linux windows

all: build-all

build:
	go build -o $(BIN_DIR)/$(BINARY_NAME) $(CMD_PATH)

build-all: darwin linux windows

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_PATH)
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-darwin-arm64 $(CMD_PATH)

linux:
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_PATH)
	GOOS=linux GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-linux-arm64 $(CMD_PATH)

windows:
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_PATH)
	GOOS=windows GOARCH=arm64 go build -o $(BIN_DIR)/$(BINARY_NAME)-windows-arm64.exe $(CMD_PATH)

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./...
