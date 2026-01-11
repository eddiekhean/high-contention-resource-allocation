APP_NAME=app
COMPOSE_LOCAL=docker-compose.local.yml
GO_MAIN=cmd/server/main.go

# ===== Go =====
run:
	docker compose -f $(COMPOSE_LOCAL) up -d maze-service redis
	go run $(GO_MAIN)

build:
	go build -o bin/$(APP_NAME) $(GO_MAIN)

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -rf bin

# ===== Docker helpers =====
compose-up:
	docker compose -f $(COMPOSE_LOCAL) up --build

compose-down:
	docker compose -f $(COMPOSE_LOCAL) down

compose-logs:
	docker compose -f $(COMPOSE_LOCAL) logs -f

compose-ps:
	docker compose -f $(COMPOSE_LOCAL) ps
