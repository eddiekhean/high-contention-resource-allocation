APP_NAME=app
COMPOSE_LOCAL=docker-compose.local.yml

# ===== Go local =====
run:
	go run main.go

build:
	go build -o bin/$(APP_NAME) .

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf bin

# ===== Docker Compose (Local Dev) =====
up:
	docker compose -f $(COMPOSE_LOCAL) up --build

down:
	docker compose -f $(COMPOSE_LOCAL) down

restart:
	docker compose -f $(COMPOSE_LOCAL) down && docker compose -f $(COMPOSE_LOCAL) up --build

logs:
	docker compose -f $(COMPOSE_LOCAL) logs -f

ps:
	docker compose -f $(COMPOSE_LOCAL) ps
