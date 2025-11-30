## up: start all containers in background without forcing build
up:
	@echo "Starting all images..."
	docker-compose up -d
	@echo "All images started."

## up-build: stop all (if running), build and start all containers in background
up-build: build-all

build-all: build-broker build-auth build-logger
	@echo "Building all images..."
	docker-compose down
	docker-compose up --build -d
	@echo "All images built and started."

## build-broker: Build broker service using Docker Compose (not Docker build directly)
build-broker:
	@echo "Building broker service using Docker Compose..."
	docker-compose build broker-service
	@echo "Broker service built using Docker Compose."

build-auth:
	@echo "Building auth service using Docker Compose..."
	docker-compose build auth-service
	@echo "Auth service built using Docker Compose."

build-logger:
	@echo "Building logger service using Docker Compose..."
	docker-compose build logger-service
	@echo "Logger service built using Docker Compose."

## down: Stop and remove all containers
down:
	@echo "Stopping all images..."
	docker-compose down
	@echo "All images stopped."
