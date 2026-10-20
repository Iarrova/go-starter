.PHONY: build run test lint check generate openapi openapi-check docker-build docker-build-migrate migrate-diff migrate-validate migrate-apply

MIGRATIONS_DIR := file://migrations
DEV_URL := docker://postgres/18/dev?search_path=public

build:
	go build -o bin/ ./cmd/...

run:
	go run ./cmd/api

test:
	go test -race ./...

lint:
	golangci-lint run

# Everything CI should run.
check: lint test openapi-check migrate-validate

generate:
	go generate ./ent

openapi:
	mkdir -p api
	go run ./cmd/openapi > api/openapi.yaml.tmp && mv api/openapi.yaml.tmp api/openapi.yaml

# Fails if api/openapi.yaml is out of date with the code.
openapi-check:
	go run ./cmd/openapi | diff -u api/openapi.yaml -

docker-build:
	docker build -f build/package/Dockerfile -t go-starter .

docker-build-migrate:
	docker build -f build/package/Dockerfile.migrate -t go-starter-migrate .

migrate-diff:
	$(if $(name),,$(error usage: make migrate-diff name=<description>))
	atlas migrate diff $(name) \
		--dir "$(MIGRATIONS_DIR)" \
		--to "ent://ent/schema" \
		--dev-url "$(DEV_URL)"

migrate-validate:
	atlas migrate validate \
		--dir "$(MIGRATIONS_DIR)" \
		--dev-url "$(DEV_URL)"

migrate-apply:
	atlas migrate apply \
		--dir "$(MIGRATIONS_DIR)" \
		--url "$(DATABASE_URL)"
