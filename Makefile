clean:
	go clean -testcache

build:
	go build -v ./...

test:
	go test -tags test -v ./...

lint:
	golangci-lint run ./...

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0

all: clean build test lint

