# Companies Microservice

## Overview

Companies is a microservice for managing company information, built with Go using a layered architecture. The service provides a REST API for CRUD operations on companies, publishes change events to Apache Kafka, and supports JWT authentication for secured endpoints.

### Key Features
- 🏢 Full CRUD for company management
- 🔐 JWT authentication for secured operations
- 📨 Asynchronous event publishing to Kafka
- 🗄️ PostgreSQL for data storage
- 🐳 Docker containerization
- ✅ Unit and integration test coverage

## Quick Start

### Prerequisites
- Go 1.19+
- Docker and Docker Compose
- Make
- curl (for API testing)

### Local Setup

1. Clone the repository:
```bash
git clone https://github.com/ezhdanovskiy/companies.git
cd companies
```

2. Run the application with infrastructure:
```bash
make run/local
```
This command will:
- Start PostgreSQL and Kafka in Docker
- Create Kafka topic `companies-mutations`
- Apply database migrations
- Build and run the application

3. In a separate terminal, test the API:
```bash
make company/lifecycle
```
This will execute a complete CRUD cycle on a company.

4. To view Kafka events:
```bash
make kafka/topic/consume
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests
make test/int

# Full test cycle with docker-compose
make test/int/docker-compose
```

## Project Structure

```
.
├── cmd/
│   └── companies/          # Application entry point
│       └── main.go
├── internal/              # Internal application packages
│   ├── application/       # Initialization and orchestration
│   │   ├── application.go
│   │   └── logger.go
│   ├── auth/             # JWT authentication
│   │   └── jwt.go
│   ├── config/           # Application configuration
│   │   └── config.go
│   ├── http/             # HTTP layer (Gin)
│   │   ├── handlers.go
│   │   ├── server.go
│   │   ├── dependencies.go
│   │   ├── requests/     # Request DTOs
│   │   └── mocks/        # Test mocks
│   ├── kafka/            # Kafka producer
│   │   ├── producer.go
│   │   └── message.go
│   ├── middlewares/      # HTTP middlewares
│   │   └── auth.go
│   ├── models/           # Domain models
│   │   ├── company.go
│   │   └── errors.go
│   ├── repository/       # Database layer
│   │   ├── repository.go
│   │   ├── repository_test.go
│   │   └── entities.go
│   ├── service/          # Business logic
│   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── dependencies.go
│   │   └── mocks/
│   └── tests/            # Integration tests
│       └── integration_test.go
├── migrations/           # SQL migrations
├── docker-compose.yml    # Docker configuration
├── Dockerfile           # Application image
├── Makefile            # Development commands
├── go.mod              # Go module
└── CLAUDE.md           # Claude AI instructions
```

## Architecture

The application is built using a Layered Architecture pattern:

![Package Dependencies Diagram](docs/diagrams/package-dependencies.png)

### Application Layers

1. **HTTP Layer** (`internal/http/`)
   - Handles HTTP requests using Gin framework
   - Input validation
   - Routing and middleware

2. **Service Layer** (`internal/service/`)
   - Business logic implementation
   - Event publishing to Kafka
   - Coordination between repository and external services

3. **Repository Layer** (`internal/repository/`)
   - PostgreSQL operations using Bun ORM
   - Database logic encapsulation

4. **Application Layer** (`internal/application/`)
   - Component initialization
   - Application lifecycle management
   - Logging configuration

### External Dependencies

- **PostgreSQL** - Primary data storage
- **Apache Kafka** - Message broker for asynchronous events
- **Zookeeper** - Kafka coordinator

### API Endpoints

#### Public Endpoints
- `GET /api/v1/companies/:uuid` - Get company information

#### Secured Endpoints (require JWT token)
- `POST /api/v1/secured/companies` - Create new company
- `PATCH /api/v1/secured/companies/:uuid` - Update company
- `DELETE /api/v1/secured/companies/:uuid` - Delete company

### Data Model

```go
type Company struct {
    ID              uuid.UUID
    Name            string    // unique, max 15 characters
    Description     string    // max 3000 characters
    EmployeesAmount int
    Registered      bool
    Type            CompanyType
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// CompanyType - company types
type CompanyType string

const (
    Corporations       CompanyType = "Corporations"
    NonProfit         CompanyType = "NonProfit"
    Cooperative       CompanyType = "Cooperative"
    SoleProprietorship CompanyType = "Sole Proprietorship"
)
```

## Configuration

The application is configured through environment variables. All settings are loaded via Viper.

### Environment Variables

#### Database
- `DB_HOST` - PostgreSQL host (default: `localhost`)
- `DB_PORT` - PostgreSQL port (default: `5432`)
- `DB_USER` - Database user (default: `db`)
- `DB_PASSWORD` - Database password (default: `db`)
- `DB_NAME` - Database name (default: `db`)

#### Kafka
- `KAFKA_ADDR` - Kafka broker address (default: `localhost:9092`)
- `KAFKA_TOPIC` - Event topic (default: `companies-mutations`)

#### HTTP Server
- `HTTP_PORT` - HTTP server port (default: `8080`)

#### Authentication
- `JWT_KEY` - Secret key for JWT tokens

#### Logging
- `LOG_LEVEL` - Log level (debug, info, warn, error)
- `LOG_ENCODING` - Log format (json, console)

## Available Commands

### Build and Run
```bash
make build                # Build binary
make run                  # Run built binary
make run/local            # Full local run with infrastructure
```

### Testing
```bash
make test                 # Run unit tests
make test/int             # Run unit and integration tests
make test/int/docker-compose  # Full integration test cycle
```

### Code Quality
```bash
make lint                 # Run golangci-lint
make fmt                  # Format code
make generate             # Generate test mocks
```

### Docker and Infrastructure
```bash
make up                   # Start PostgreSQL and Kafka
make down                 # Stop containers
make kafka/topic/create   # Create Kafka topic
make kafka/topic/consume  # View topic messages
```

### Database Migrations
```bash
make migrate/up           # Apply migrations
make migrate/down         # Rollback last migration
```

### API Testing
```bash
make company/lifecycle    # Full CRUD cycle via curl
make company/create       # Create company
make company/get          # Get company
make company/patch        # Update company
make company/delete       # Delete company
```

### Diagrams
```bash
make diagrams             # Generate diagrams from DOT files
```

## Development

### Code Conventions
- Standard Go formatting (gofmt)
- Linter: golangci-lint with default settings
- Mocks generated using gomock

### Testing
- Unit tests are located alongside code (`*_test.go`)
- Integration tests in `internal/tests/`
- Integration tests use the `integration` build tag
- Tests require running PostgreSQL and Kafka

### Kafka Events
All mutating operations (CREATE, UPDATE, DELETE) publish events to the `companies-mutations` topic:

```json
{
  "type": "CREATE|UPDATE|DELETE",
  "companyId": "uuid",
  "data": {
    // company data
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### JWT Authentication
- Algorithm: HS256
- Token passed in header: `Authorization: Bearer <token>`
- Secured endpoints require valid token

## Development

### Adding New Features

1. Define models in `internal/models/`
2. Add repository methods in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Create HTTP handlers in `internal/http/`
5. Write tests for each layer
6. Update documentation

### Generating Test Mocks

```bash
make generate
```

This will create mocks for interfaces marked with comment:
```go
//go:generate mockgen -source=file.go -destination=mocks/file_mock.go
```

## Deployment

### Docker

To build Docker image:
```bash
docker build -t companies:latest .
```

### Docker Compose

Full deployment with infrastructure:
```bash
docker-compose up -d
```

## Monitoring and Logs

The application uses structured logging via Zap. Logs are output to stdout in JSON format (production) or console format (development).

View logs:
```bash
# Application logs
docker-compose logs -f companies

# All service logs
docker-compose logs -f
```

## Troubleshooting

### Database Connection Issues
1. Check PostgreSQL is running: `docker-compose ps`
2. Verify environment variables
3. Ensure migrations are applied: `make migrate/up`

### Kafka Issues
1. Check Kafka and Zookeeper are running
2. Ensure topic is created: `make kafka/topic/create`
3. Check Kafka logs: `docker-compose logs kafka`

### Test Issues
1. Integration tests require running infrastructure
2. Use `make test/int/docker-compose` for full cycle
3. Check ports 5432 (PostgreSQL) and 9092 (Kafka) are available

## Language Support

For Russian documentation, see [README_ru.md](README_ru.md).

## License

This project is distributed under the MIT License. See [LICENSE](LICENSE) file for details.