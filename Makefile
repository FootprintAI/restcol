.PHONY: build test test-race test-short vet lint tidy run-local run-postgres gen-proto gen-swagger gen gen-check clean help

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

## gen-proto: Regenerate protobuf/gateway/OpenAPI output from the .proto (requires buf)
gen-proto:
	cd api && ./gen-proto-go.sh

## gen-swagger: Regenerate the go-swagger client from the OpenAPI spec (requires the pinned swagger)
gen-swagger:
	cd api && ./gen-swagger-go-client.sh

## gen: Regenerate everything under api/
gen: gen-proto gen-swagger

## gen-check: Fail if the committed generated output does not match the generators
#
# Scoped to the three generated directories, not to api/ as a whole: api/ also
# holds the .proto, the buf config and the generator scripts, and diffing those
# means editing a generator reports its own edit as "generated code is stale" -
# true but useless. These three are the outputs the check exists to protect.
GENERATED := api/pb api/openapiv2 api/go-openapiv2

gen-check:
	cd api && buf generate
	cd api && ./gen-swagger-go-client.sh
	@git diff --exit-code -- $(GENERATED) \
		|| (echo ""; \
		    echo "Generated code is stale: run 'make gen' and commit the result."; \
		    exit 1)

## clean: Remove build artifacts
clean:
	rm -f $(BINARY) coverage.out

## help: Show this help message
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/## //'
