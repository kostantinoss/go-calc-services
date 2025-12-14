# Go Calculator Microservices

A microservices-based calculator API built with Go, featuring an API gateway, multiple calculator server instances, and Docker containerization.

## Architecture

```
Client
  ↓
API Gateway (port 8080)
  ├─→ Server1 → Calculator Operations (/add, /sub)
  └─→ Server2 → Calculator Operations (/multi, /div)
```

## Project Structure

```
services/
├── gateway/              # API Gateway service
│   ├── internal/
│   │   ├── config.go    # Configuration management
│   │   ├── gateway.go   # Gateway logic & request forwarding
│   │   └── routing.go   # Route matching
│   ├── config.yaml      # Gateway configuration
│   ├── Dockerfile
│   └── main.go
├── server/              # Calculator server
│   ├── Dockerfile
│   └── main.go         # Calculator endpoints
├── loadbalancer/        # Load balancer (planned)
├── docker-compose.yml   # Service orchestration
└── go.work             # Go workspace file
```

## Quick Start

### Prerequisites

- Go 1.25.3 or later
- Docker & Docker Compose

### Running with Docker Compose

```bash
# Build and start all services in background
docker compose up --build -d

# Run in background
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

## API Endpoints

All endpoints accept JSON payloads with two numbers:

### Request Format
```json
{
  "a": 10,
  "b": 5
}
```

### Response Format
```json
{
  "result": 15
}
```

### Available Operations

| Endpoint | Method | Operation      | Example Result |
|----------|--------|----------------|----------------|
| `/add`   | POST   | Addition       | `a + b = 15`   |
| `/sub`   | POST   | Subtraction    | `a - b = 5`    |
| `/multi` | POST   | Multiplication | `a * b = 50`   |
| `/div`   | POST   | Division       | `a / b = 2`    |

### Example Requests

```bash
# Addition
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5}'

# Subtraction
curl -X POST http://localhost:8080/sub \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5}'

# Multiplication
curl -X POST http://localhost:8080/multi \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5}'

# Division
curl -X POST http://localhost:8080/div \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5}'
```

## Configuration

### Gateway Configuration

The gateway uses `config.yaml` for routing configuration:

```yaml
api_gateway:
  gateway_server:
    port: "8080"
    
  target_servers:
    - name: server1
      url: http://server1:8080
    - name: server2
      url: http://server2:8080

  routing:
    - path: /add
      methods: [POST]
      server: server1
    - path: /sub
      methods: [POST]
      server: server1
    - path: /multi
      methods: [POST]
      server: server2
    - path: /div
      methods: [POST]
      server: server2
```

### Environment Variables

| Variable      | Description                    | Default                |
|---------------|--------------------------------|------------------------|
| `CONFIG_PATH` | Path to gateway config file    | `/app/config.yaml`     |
| `PORT`        | Server port                    | `8080`                 |
| `SERVER1`     | Server1 URL                    | `http://server1:8080`  |
| `SERVER2`     | Server2 URL                    | `http://server2:8080`  |

## Docker

### Building Images

```bash
# Build all images
docker compose build

# Build specific service
docker compose build gateway
docker compose build server
```


## Services

### API Gateway
- Routes incoming requests to appropriate backend services
- YAML-based configuration
- Request/response forwarding
- Structured logging

### Calculator Server
- REST API for basic arithmetic operations
- JSON request/response handling

### Load Balancer (Planned)
- Distributes traffic across multiple server instances
- Health checking
- Multiple load balancing strategies

## Technology Stack

- **Language**: Go 1.25.3
- **Configuration**: YAML (gopkg.in/yaml.v2)
- **Containerization**: Docker & Docker Compose
- **Networking**: Docker networking (bridge network)
- **HTTP**: Standard library (`net/http`)

## Logging

All services use structured logging with the following format:

Example:
```
[GATEWAY] 2025/12/13 23:27:56 gateway.go:113: Router initialized
[GATEWAY] 2025/12/13 23:27:56 gateway.go:126: Route added: /add -> http://server1:8080
```


## TODOs

- [ ] Load balancer implementation
- [ ] Health check endpoints
- [ ] Request rate limiting
- [ ] Metrics & monitoring (Prometheus)
- [ ] Unit & integration tests
- [ ] CI/CD pipeline


## Author

Konstantinos Chondralis

---

**Note**: This is a learning project demonstrating microservices patterns in Go. It showcases:
- Service-to-service communication
- API Gateway pattern
- Docker containerization
- Configuration management

