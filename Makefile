SHELL := /bin/bash
include .dev.env
export

.phony: build, run

build:
	cd cmds/swdocs && go build .;

run: build
	cd cmds/swdocs && ./swdocs serve;

test:
	go test -v ./...

run_postgres:
	docker run -it -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres
