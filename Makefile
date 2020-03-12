SHELL := /bin/bash

run:
	docker-compose up --build
build:
	docker build -t exlibris:latest .
run-local:
	docker-compose -f docker-compose.local.yml stop
	docker-compose -f docker-compose.local.yml up -d
	npm install
	go get github.com/githubnemo/CompileDaemon
	set -a && source app.env && CompileDaemon -build='go build -o exlibris' -directory=. -command=./exlibris &
	set -a && source .env && npm run serve
dev:
	docker-compose -f docker-compose.dev.yml stop
	docker-compose -f docker-compose.dev.yml up --build
