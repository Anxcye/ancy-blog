# File: Makefile
# Purpose: Provide repeatable developer commands for build, test, and migrations.
# Module: backend, local automation layer.
# Related: cmd/server and cmd/migrate entrypoints.

.PHONY: run test fmt migrate-up migrate-down migrate-version

run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

migrate-up:
	go run ./cmd/migrate -action up

migrate-down:
	go run ./cmd/migrate -action down

migrate-version:
	go run ./cmd/migrate -action version
