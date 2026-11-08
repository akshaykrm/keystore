.PHONY: sqlc tags play

build: 
	go run ./cmd/server/main.go
sqlc:
	sqlc generate
tags:
	ctags -R .
play:
	go run ./cmd/playground/main.go
check:
	@go build ./cmd/server/main.go
