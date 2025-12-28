# {{.ProjectName}}

{{.ProjectName}} is a Go web application built with go-start.

## Project Structure

```
{{.ProjectName}}/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── controller/          # HTTP handlers
│   ├── service/             # Business logic
│   ├── repository/          # Data access layer
│   └── model/               # Data models
├── config/
│   ├── config.go            # Configuration loading
│   └── config.yaml          # Configuration file
└── pkg/                     # Private packages
    ├── cache/               # Redis cache wrapper
    ├── database/            # Database wrapper
    └── httpx/               # HTTP utilities
```

## Getting Started

### Prerequisites

- Go 1.25.4 or higher
- MySQL 5.7+ or PostgreSQL 12+
{{if .WithRedis}}
- Redis 6.0+
{{end}}

### Installation

1. Clone the repository
```bash
git clone <your-repo-url>
cd {{.ProjectName}}
```

2. Install dependencies
```bash
go mod download
```

3. Configure the application
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your settings
```

4. Run the application
```bash
go run cmd/server/main.go
```

### Build

```bash
go build -o bin/server cmd/server/main.go
./bin/server
```

## API Documentation

### Health Check

```
GET /health
```

### User APIs

```
POST   /api/v1/users        # Create user
GET    /api/v1/users        # List users
GET    /api/v1/users/:id    # Get user
PUT    /api/v1/users/:id    # Update user
DELETE /api/v1/users/:id    # Delete user
```

## Configuration

Edit `config.yaml` to configure:

- Server port
- Database connection
{{if .WithRedis}}
- Redis connection
{{end}}
{{if .WithSwagger}}
- Swagger documentation is available at `/swagger/index.html`
{{end}}

## Development

### Add New API

1. Define model in `internal/model/`
2. Create repository in `internal/repository/`
3. Create service in `internal/service/`
4. Create controller in `internal/controller/`
5. Register routes in `cmd/server/main.go`

## License

MIT
