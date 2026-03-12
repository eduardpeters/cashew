.PHONY: build run test fmt vet

BINARY = cashew

build:
	go build -o $(BINARY) ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...