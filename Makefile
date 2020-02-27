run: build
	docker-compose run app
build:
	docker build -t exlibris:latest .
run-local:
	docker-compose -f docker-compose.local.yml up -d
	go run main.go
