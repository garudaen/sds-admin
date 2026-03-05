# SDS Admin - DNS Management System

A modern DNS management system built with Golang (Gin) backend and Vue.js 3 frontend.

## Features

- **Domain Management**: Add, edit, delete, enable/disable domains
- **DNS Record Management**: Support for A, AAAA, CNAME, MX, TXT record types
- **Smart DNS Resolution**: Client IP CIDR-based routing for A/AAAA/CNAME records
- **RESTful API**: Well-documented API with Swagger support
- **Modern UI**: Responsive web interface with mobile support

## Tech Stack

### Backend
- **Framework**: Gin (Golang)
- **ORM**: GORM
- **Database**: MySQL
- **Logger**: Logrus with Lumberjack rotation
- **API Documentation**: Swagger/OpenAPI

### Frontend
- **Framework**: Vue.js 3 (Composition API)
- **Build Tool**: Vite
- **HTTP Client**: Axios
- **Static Files**: Embedded with Statik

## Project Structure

```
sds-admin/
├── cmd/
│   └── sds-admin/          # Application entry point
│       └── main.go
├── configs/
│   ├── config-sample.yaml  # Sample configuration
│   └── config.yaml         # Actual configuration (gitignored)
├── internal/
│   ├── config/             # Configuration loading
│   ├── database/           # Database connection and migrations
│   ├── dto/                # Data Transfer Objects
│   ├── handler/            # HTTP handlers
│   ├── logger/             # Logging configuration
│   ├── models/             # Database models
│   ├── router/             # Route definitions
│   ├── service/            # Business logic
│   └── static/             # Embedded static files (generated)
├── pkg/                    # Public libraries
├── fe/                     # Frontend source
│   ├── src/
│   │   ├── App.vue
│   │   └── components/
│   │       ├── AdminLayout.vue
│   │       ├── DomainManagement.vue
│   │       └── RecordManagement.vue
│   └── package.json
├── docs/                   # Swagger documentation (generated)
├── logs/                   # Log files (gitignored)
├── pub/                    # Built frontend files (gitignored)
├── bin/                    # Compiled binaries (gitignored)
└── Makefile
```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+
- MySQL 8.0+

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd sds-admin
   ```

2. **Configure the application**
   ```bash
   cp configs/config-sample.yaml configs/config.yaml
   # Edit configs/config.yaml with your database credentials
   ```

3. **Create database**
   ```sql
   CREATE DATABASE sds_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

4. **Build and run**
   ```bash
   make build-all
   ./bin/sds-admin --config configs/config.yaml
   ```

   Or run directly:
   ```bash
   make run
   ```

5. **Access the application**
   - Web UI: http://localhost:8080
   - Swagger API Docs: http://localhost:8080/swagger/index.html

## API Endpoints

### Domains

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/domains` | List all domains |
| GET | `/api/v1/domains/:id` | Get domain by ID |
| POST | `/api/v1/domains` | Create domain |
| PUT | `/api/v1/domains/:id` | Update domain |
| DELETE | `/api/v1/domains/:id` | Delete domain |
| POST | `/api/v1/domains/:id/disable` | Disable domain |
| POST | `/api/v1/domains/:id/enable` | Enable domain |

### Records

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/domains/:id/records` | List records for a domain |
| GET | `/api/v1/domains/:id/records/:recordId` | Get record by ID |
| POST | `/api/v1/domains/:id/records` | Create record |
| PUT | `/api/v1/domains/:id/records/:recordId` | Update record |
| DELETE | `/api/v1/domains/:id/records/:recordId` | Delete record |
| POST | `/api/v1/domains/:id/records/:recordId/disable` | Disable record |
| POST | `/api/v1/domains/:id/records/:recordId/enable` | Enable record |

### Record Types

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/record-types` | List all record types |
| GET | `/api/v1/record-types/:id` | Get record type by ID |

## DNS Record Types

### A Record
- Maps a domain to an IPv4 address
- Supports multiple values with smart resolution via CIDR
- Requires at least one default value

### AAAA Record
- Maps a domain to an IPv6 address
- Supports multiple values with smart resolution via CIDR
- Requires at least one default value

### CNAME Record
- Creates an alias to another domain
- Cannot coexist with other record types on the same host
- Automatically appends trailing dot if missing

### MX Record
- Mail exchange record
- Each value has its own priority (0-65535)
- Lower priority number = higher priority

### TXT Record
- Text record for various purposes
- Simple value storage

## Development

### Available Make Commands

```bash
make help          # Show all available commands
make deps          # Install Go dependencies
make build         # Build backend only
make build-all     # Build frontend and backend
make run           # Run the application
make clean         # Clean build artifacts
make test          # Run tests
make fmt           # Format code
make lint          # Run linter
make swagger       # Generate Swagger documentation
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

### Frontend Development

```bash
cd fe
npm install
npm run dev     # Development server with hot reload
npm run build   # Production build
```

### Database Migrations

Tables are automatically created on first run:
- `domains` - Domain names
- `record_types` - DNS record types (A, AAAA, CNAME, MX, TXT)
- `records` - DNS records
- `record_values` - Record values with CIDR support

## Configuration

See `configs/config-sample.yaml` for all available options:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  mode: release          # debug, release, test

database:
  driver: mysql
  host: localhost
  port: 3306
  username: root
  password: your_password
  database: sds_admin

log:
  level: info            # debug, info, warn, error
  format: json           # json, text
  output: logs/sds-admin.log
  max_size: 100          # MB
  max_backups: 3
  max_age: 7             # days
  compress: true

swagger:
  enabled: true
```

## License

**WTFPL (Do What The Fuck You Want To Public License)**

This project is entirely generated by AI (Qoder). You can do whatever you want with it - no restrictions, no warranties.

```
            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                    Version 2, December 2004

 Copyright (C) 2024 AI Generated

 Everyone is permitted to copy and distribute verbatim or modified
 copies of this license document, and changing it is allowed as long
 as the name is changed.

            DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
   TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

  0. You just DO WHAT THE FUCK YOU WANT TO.
```
