.PHONY: build test test-race test-short vet lint tidy run-local run-postgres gen-proto clean help

BINARY := restcol
GO      ?= go

## build: Compile the server binary
build:
	$(GO) build -o $(BINARY) ./

## test: Run tests (excludes tests that require postgres via -short)
test:
	$(GO) test ./... -short

## test-race: Run tests with the race detector
test-race:
	$(GO) test -race ./... -short

## test-full: Run all tests including integration tests that require postgres
test-full:
	$(GO) test -race ./... -coverprofile=coverage.out

## vet: Run go vet
vet:
	$(GO) vet ./...

## tidy: Tidy module dependencies
tidy:
	$(GO) mod tidy

## run-postgres: Start the local postgres container used by tests and run-local
run-postgres:
	./run_postgres.sh

## run-local: Build and run the server against local postgres
run-local:
	./run_local.sh

## gen-proto: Regenerate proto/OpenAPI clients (requires buf)
gen-proto:
	cd api && ./gen-proto-go.sh

## clean: Remove build artifacts
clean:
	rm -f $(BINARY) coverage.out

## help: Show this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## //'
