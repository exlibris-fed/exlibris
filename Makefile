run: build
	docker-compose run app
build:
	docker build .
run-local:
	docker-compose -f docker-compose.local.yml up -d
	go run main.go
