#BUILDER STAGE
FROM golang:1.25.1-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/server

#RUNTIME STAGE
FROM alpine:latest
WORKDIR /app
COPY --from=builder ./app/app .
COPY config.yaml .
EXPOSE 8080
CMD ["./app"]