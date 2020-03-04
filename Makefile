run:
	docker-compose up --build
build:
	docker build -t exlibris:latest .
run-local:
	docker-compose -f docker-compose.local.yml down
	docker-compose -f docker-compose.local.yml up -d
	npm run build
	go run main.go
dev:
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml up --build
