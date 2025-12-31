APP_NAME=app

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
