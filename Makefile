SOURCE  = ./cmd/server/
OUTPUT  = ./bin/server

default: fmt build

setup:
	@echo "Downloading dependencies..."
	go mod download
	@echo "Done."

build:
	@echo "Building..."
	@mkdir -p bin
	go build -o $(OUTPUT) $(SOURCE)
	@echo "Build complete."

clean:
	rm -rf ./bin

fmt:
	go fmt ./...

test:
	go test ./... -v

run: build
	./$(OUTPUT) -config-path config/dev/config.yaml

.PHONY: setup build clean fmt test run
