PROJECT_NAME = go-monolite

## init: used to initialize the Go project, tidy, docker, migration, build and deploy
.PHONY: init
init:
	find ./ -type f -not -path '*/.git/*' -exec sed -i 's/ms_delivery/$(PROJECT_NAME)/g' {} +
	go mod init gitlab.toledo24.ru/web/$(PROJECT_NAME) || true
	go mod tidy
	docker compose up -d
	go run cmd/migration/main.go
	go build -o build/package/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/main.go

## init-docker: used to initialize the Go project docker, tidy, docker, migration, build and deploy
.PHONY: init-docker
init-docker:
	go run cmd/migration/main.go
	go build -o build/package/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/main.go
	build/package/$(PROJECT_NAME)

## deploy: executing the deployment command
.PHONY: swagger-init
swagger-init:
	swag init -g cmd/go-monolite/main.go -o docs

## fast-start: quick launch of ms_delivery
.PHONY: fast-start
fast-start:
	go run cmd/migration/main.go -action=up
	go run cmd/$(PROJECT_NAME)/main.go

## start: build start of $(PROJECT_NAME)
.PHONY: start
start:
	go build -o build/package/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/main.go
	build/package/$(PROJECT_NAME)

## migration-up: start the migration stage with the database
.PHONY: migration-up
migration-up:
	go run cmd/migration/main.go -action=up

## migration-down: down the migration with the database
.PHONY: migration-down
migration-down:
	go run cmd/migration/main.go -action=down

## build: build a project
.PHONY: build
build:
	go build -o build/package/$(PROJECT_NAME) cmd/$(PROJECT_NAME)/main.go

## lint: format and golangci-lint the project
.PHONY: lint
lint:
	gofmt -s -w .
	golangci-lint run

## test: start test
.PHONY: test
test:
	go run cmd/migration/main.go -action=up -env=.env.test
	go test -v ./internal/module/...
