APP_NAME=app

run:
	docker compose -f docker-compose.local.yml up -d
	go run ./cmd/server
down:
	docker compose -f docker-compose.local.yml down
restart:
	docker compose -f docker-compose.local.yml down
	docker compose -f docker-compose.local.yml up -d
	go run ./cmd/server

build:
	go build -o bin/$(APP_NAME) ./cmd/server

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf bin
