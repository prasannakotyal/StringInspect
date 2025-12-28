# Makefile for StringInspect

BINARY = stringinspect
SRC = ./...

.PHONY: all build run clean test test-coverage fmt lint install uninstall

all: build

build:
	go build -o $(BINARY) .

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	go clean

test:
	go test -v $(SRC)

test-coverage:
	go test -v -cover -coverprofile=coverage.out $(SRC)
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt $(SRC)
	goimports -w .

lint:
	golangci-lint run

install: build
	cp $(BINARY) /usr/local/bin/

uninstall:
	rm -f /usr/local/bin/$(BINARY)
