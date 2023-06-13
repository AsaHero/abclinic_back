-include .env
export

CURRENT_DIR=$(shell pwd)
APP=abclinic
CMD_DIR=./cmd

.DEFAULT_GOAL = build

# build for current os
.PHONY: build
build:
	go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/main.go

# build for linux amd64
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOARCH="amd64" GOOS=linux go build -ldflags="-s -w" -o ./bin/${APP} ${CMD_DIR}/main.go

# run service
.PHONY: run
run:
	go run ${CMD_DIR}/app/main.go

# generate swagger
.PHONY: swagger-gen
swagger-gen:
	swag init --parseDependency --dir ./api -g router.go -o ./api/docs

# migrate
.PHONY: migrate
migrate:
	migrate -source file://migrations -database postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable up
