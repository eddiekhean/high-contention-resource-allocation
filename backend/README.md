# High Contention Resource Allocation Backend

This project demonstrates a high-performance backend system designed to handle high-contention resource allocation scenarios, alongside a specialized maze generation service. It serves as a proof-of-concept for scalable architecture using Go, Python, and Redis.

## Project Overview

The core problem this system solves is the race condition inherent in allocating limited resources (slots) under high concurrent load. It uses Redis atomic operations to ensure data consistency without heavy database locking. Additionally, it features a separate Python-based service for generating perfect mazes using the Recursive Backtracker algorithm, demonstrating a polyglot microservices approach.

## Features

- **High-Performance Resource Allocation**: Uses Redis atomic counters (`DECR`/`INCR`) to manage concurrent slot reservations efficiently.
- **Maze Generation Engine**: A dedicated Python service producing fully connected, "perfect" mazes.
- **Microservices Architecture**: Separation of concerns between the core API (Go) and computational tasks (Python).
- **Rate Limiting**: Built-in middleware to prevent abuse.
- **Dockerized Environment**: fully containerized setup for easy deployment.

## System Architecture

The system consists of three main components:

- **Core API (Go)**: Acts as the gateway and business logic handler. It manages resource allocation and routes requests.
- **Maze Service (Python)**: A specialized worker service running FastAPI to generate maze structures.
- **Redis**: The shared state layer used for:
    - Atomic resource counters (Semaphores)
    - High-speed caching/storage

## Tech Stack

- **Backend (Core)**: Go 1.22+, Gin Web Framework
- **Backend (Maze)**: Python 3.10+, FastAPI, Pydantic
- **Data Store**: Redis 7
- **Infrastructure**: Docker, Docker Compose
- **Logging**: Logrus (Go)

## Project Structure

```
.
├── cmd/server/          # Entry point for the Go API server
├── internal/
│   ├── client/          # Redis client wrappers
│   ├── config/          # Configuration loading
│   ├── handler/         # HTTP request handlers (Controllers)
│   ├── middleware/      # Rate limiting, CORS, Logging
│   ├── models/          # Data structures and domain models
│   ├── service/         # Business logic layer
│   ├── storage/         # Redis storage implementations (SlotStore)
│   └── utils/           # Helper functions
├── maze_generator/      # Python Maze Service (FastAPI)
├── docker-compose.yml   # Production orchestration
├── docker-compose.local.yml # Local development orchestration
└── Makefile             # Development automation commands
```

## Setup & Installation

### Prerequisites
- [Go](https://go.dev/) 1.22+
- [Docker](https://www.docker.com/) & Docker Compose
- [Make](https://www.gnu.org/software/make/) (optional, for using Makefile)

### Local Setup

1. **Clone the repository:**
   ```bash
   git clone https://github.com/eddiekhean/high-contention-resource-allocation-backend.git
   cd high-contention-resource-allocation-backend
   ```

2. **Run with Docker Compose (Recommended):**
   This starts Redis, the Maze Service, and the Go API.
   ```bash
   make run
   # OR directly:
   docker compose -f docker-compose.local.yml up -d maze-service redis
   go run cmd/server/main.go
   ```

3. **Verify Health:**
   Check if the server is running:
   ```bash
   curl http://localhost:8080/health
   ```

## Usage

### API Endpoints

The API is exposed at `http://localhost:8080/api/v1/public`.

#### 1. Simulate Resource Allocation
Simulate high-concurrency requests for limited slots.
- **POST** `/simulate`

#### 2. Generate Maze
Generate a new random maze.
- **POST** `/leetcode/maze/generate`
- **Body:**
  ```json
  {
      "rows": 20,
      "cols": 20
  }
  ```

#### 3. Submit Maze Solution
Submit a path for verification (if implemented).
- **POST** `/leetcode/maze/submit`

## Configuration

Configuration is managed via `config.yaml` (default) or environment variables.

key | Description | Default
--- | --- | ---
`redis.addr` | Redis connection string | `localhost:6379`
`server.port` | API listening port | `8080`
`cors.allowed_origins` | CORS whitelist | `*`
`maze_service_url` | URL of internal Python service | `http://localhost:8000`

## Testing

Run unit and integration tests for the Go application:

```bash
make test
```

## Performance & Scalability

- **Optimistic Concurrency**: The system uses Redis `DECR` to acquire slots. If the result is negative, it rolls back with `INCR`. This avoids heavy locking enabling high throughput for "ticket booking" style problems.
- **Stateless Services**: Both Go and Python services are stateless, allowing them to be scaled horizontally (simulated in `docker-compose`).

## Security Considerations

- **Rate Limiting**: Token bucket algorithm implemented in middleware to mitigate DoS attacks.
- **Input Validation**: Strict binding and validation on all JSON inputs.

## Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.
