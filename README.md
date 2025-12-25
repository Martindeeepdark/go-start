# go-start

A powerful Go project scaffold tool that helps you quickly create web applications with MVC architecture.

## Features

- 🚀 **Quick Start** - Create a new project with a single command
- 🏗️ **MVC Architecture** - Clean separation of concerns (Controller → Service → Repository)
- 🎯 **Gin Framework** - High-performance HTTP web framework
- 💾 **Database Support** - GORM with MySQL/PostgreSQL
- ⚡ **Redis Cache** - Built-in Redis integration
- 🔧 **Production Ready** - Structured logging, error handling, graceful shutdown

## Installation

### From Source

```bash
git clone https://github.com/yourname/go-start.git
cd go-start
make install
```

### Using Go Install

```bash
go install github.com/yourname/go-start/cmd/go-start@latest
```

## Quick Start

Create a new project:

```bash
go-start create my-api
cd my-api
go mod tidy
cp config.yaml.example config.yaml
# Edit config.yaml with your database and Redis settings
go run cmd/server/main.go
```

## Usage

```bash
# Create a new project
go-start create <project-name> [flags]

# Flags:
#   -a, --arch string    Project architecture (default "mvc")
#   -m, --module string  Go module name

# Examples:
go-start create my-api                                    # Create with default settings
go-start create my-api --module=github.com/myname/my-api  # Custom module name

# Run project in development mode (with hot reload if air is installed)
go-start run

# Show version
go-start version
```

## Project Structure

```
my-api/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── controller/          # HTTP handlers (presentation layer)
│   ├── service/             # Business logic layer
│   ├── repository/          # Data access layer
│   └── model/               # Data models
├── config/
│   ├── config.go            # Configuration loading
│   └── config.yaml          # Configuration file
├── pkg/                     # Private utilities
│   ├── cache/               # Redis cache wrapper
│   ├── database/            # Database wrapper
│   └── httpx/               # HTTP utilities (middleware, response, router)
├── go.mod
└── README.md
```

## Configuration

Edit `config.yaml`:

```yaml
server:
  port: 8080

database:
  driver: mysql  # or postgres
  host: localhost
  port: 3306
  database: my_api
  username: root
  password: ""

redis:
  host: localhost
  port: 6379
  db: 0
```

## Development

### Prerequisites

- Go 1.25.4 or higher
- MySQL 5.7+ or PostgreSQL 12+
- Redis 6.0+

### Hot Reload

Install [air](https://github.com/cosmtrek/air) for hot reload:

```bash
go install github.com/cosmtrek/air@latest
```

Then run:

```bash
go-start run
```

### Build

```bash
make build
```

### Run Tests

```bash
make test
```

## License

MIT License
