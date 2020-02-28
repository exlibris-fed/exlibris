run: build
	docker-compose up
build:
	docker build -t exlibris:latest .
run-local:
	docker-compose down
	docker-compose -f docker-compose.local.yml up
	go run main.go
dev:
	docker-compose -f docker-compose.local.yml down
	docker-compose -f docker-compose.local.yml up --build