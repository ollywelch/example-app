generate:
	oapi-codegen -package main openapi.yml > server.gen.go

build:
	@go build -o ./bin/app

run: build
	@./bin/app
