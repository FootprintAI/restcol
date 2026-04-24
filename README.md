# Restcol

**One RESTful API for collaborative, schema-free document storage.**

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GitHub Issues](https://img.shields.io/github/issues/FootprintAI/restcol)](https://github.com/FootprintAI/restcol/issues)
[![GitHub Stars](https://img.shields.io/github/stars/FootprintAI/restcol)](https://github.com/FootprintAI/restcol/stargazers)

## Overview

Restcol organises data into **projects**, **collections**, and **documents**:

- **Project** — a tenant boundary. Every request is scoped to one project.
- **Collection** — a set of documents with similar shape. Collections track schema evolution over time.
- **Document** — the actual payload (JSON today; CSV/XML/media on the roadmap). Schemas are inferred on write; no up-front definition required.

Restcol speaks gRPC natively and exposes the same service over HTTP/JSON via `grpc-gateway`. Swagger/OpenAPI is auto-generated from the proto.

## Architecture

```
┌────────────┐    HTTP/JSON     ┌───────────────┐
│   client   │ ───────────────▶ │ grpc-gateway  │
└────────────┘                  │  (port 50091) │
                                └───────┬───────┘
                                        │ gRPC (internal)
                                        ▼
┌────────────┐       gRPC       ┌───────────────┐     GORM      ┌───────────┐
│ gRPC client│ ───────────────▶ │ restcol gRPC  │ ─────────────▶│ Postgres  │
└────────────┘   (port 50090)   │    server     │               └───────────┘
                                └───────────────┘
```

Server entrypoint: `main.go` → `pkg/server/app` wires auth middleware, storage (`pkg/storage/...`), schema inference (`pkg/schema`), and the RestColService handlers (`pkg/app`).

## Quick start

### 1. Start Postgres
```bash
make run-postgres   # or: ./run_postgres.sh
```
Spins up `library/postgres:16-alpine3.18` on `:5432` with user `postgres`, password `password`, database `unittest`.

### 2. Build and run the server
```bash
make run-local      # or: ./run_local.sh
```
Serves gRPC on `:50090` and HTTP/JSON on `:50091`. On first boot the default project (`1001`) is seeded so anonymous requests have a tenant.

### 3. Try it out
Swagger UI: <http://localhost:50091/swaggerui/>
Per-project API doc: <http://localhost:50091/v1/projects/1001/apidoc>

## API examples

All examples use the default project `1001`. Replace with your own project ID as needed.

### Create a collection
```bash
curl -X POST http://localhost:50091/v1/projects/1001/collections \
  -H 'Content-Type: application/json' \
  -d '{"description": "user events"}'
```

### Create a document (auto-inferred schema)
```bash
curl -X POST http://localhost:50091/v1/projects/1001/collections/<COLLECTION_ID>:newdoc \
  -H 'Content-Type: application/json' \
  -d '{"data": {"user": "alice", "event": "login", "ts": 1714000000}}'
```
Omit `collections/<id>:newdoc` to let the server auto-provision a collection:
```bash
curl -X POST http://localhost:50091/v1/projects/1001/newdoc \
  -H 'Content-Type: application/json' \
  -d '{"data": {"hello": "world"}}'
```

### Get a document
```bash
curl http://localhost:50091/v1/projects/1001/collections/<COLLECTION_ID>/docs/<DOC_ID>
```
Scope mismatch (wrong project or collection for the doc) returns `404 Not Found`, not an empty body.

### Query documents
```bash
curl 'http://localhost:50091/v1/projects/1001/collections/<COLLECTION_ID>/docs?limitCount=10'
```

### Delete a document
```bash
curl -X DELETE http://localhost:50091/v1/projects/1001/collections/<COLLECTION_ID>/docs/<DOC_ID>
```

### Delete a collection
By default, deleting a non-empty collection returns `409 Conflict`:
```bash
curl -X DELETE http://localhost:50091/v1/projects/1001/collections/<COLLECTION_ID>
# {"code":2,"message":"collection ... contains N documents; pass force=true to cascade-delete"}
```
Pass `force=true` to cascade-delete all documents first:
```bash
curl -X DELETE 'http://localhost:50091/v1/projects/1001/collections/<COLLECTION_ID>?force=true'
```

## Authentication

The server ships with `AnnonymousClaimParser` + `AllowEveryOne` authorisation — every request is accepted and mapped to the default project. This is fine for local development and demos; **configure a real JWT claim parser before exposing the service publicly**. See `pkg/server/app/app.go` for the middleware wiring.

## Development

### Common commands
```bash
make build           # compile server binary
make test            # short tests only (no postgres required)
make test-race       # short tests with -race
make test-full       # full suite (requires run-postgres first)
make vet             # go vet ./...
make tidy            # go mod tidy
make gen-proto       # regenerate api/pb/* from api/restcol.proto (requires buf)
make clean           # remove build artifacts
make help            # list all targets
```

### Repository layout

| Path                    | Purpose                                                |
|-------------------------|--------------------------------------------------------|
| `main.go`               | entrypoint; parses flags and starts the server         |
| `api/`                  | proto definitions + generated gRPC / OpenAPI clients   |
| `pkg/app/`              | gRPC service handlers (collections, documents)         |
| `pkg/server/app/`       | server assembly: middleware + storage + handlers       |
| `pkg/server/`           | swagger / OpenAPI route registration                   |
| `pkg/storage/`          | GORM-backed CRUD for projects, collections, documents  |
| `pkg/models/`           | domain models (what the storage layer reads/writes)    |
| `pkg/schema/`           | schema inference and field-path building               |
| `pkg/bootstrap/`        | seeds the default project used for anonymous auth      |
| `pkg/encoding/`         | JSON/CSV/XML payload decoders                          |
| `pkg/runtime/js/`       | goja-based JavaScript runtime for evaluating swagger   |
| `pkg/authn/`, `authz/`  | anonymous auth + allow-all authorisation (dev default) |
| `integrationtest/`      | end-to-end tests against a live server                 |

### Configuration

The server reads Postgres connection settings from flags prefixed with `--restcol_`. See `run_local.sh` for a working example. Flags:

| Flag                          | Default  |
|-------------------------------|----------|
| `--grpc_port`                 | `50090`  |
| `--http_port`                 | `50091`  |
| `--restcol_db_endpoint`       | —        |
| `--restcol_db_name`           | —        |
| `--restcol_db_user`           | —        |
| `--restcol_db_password`       | —        |
| `--restcol_auto_migrate`      | `false`  |

## Contributing

1. Fork and create a feature branch.
2. `make test-race` before pushing.
3. Open a pull request against `main`.

## License

Apache License 2.0 — see [LICENSE](LICENSE).

## Contact

Issues: <https://github.com/FootprintAI/restcol/issues>
