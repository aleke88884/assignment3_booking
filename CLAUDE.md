# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

SmartBooking is a production-ready booking system for small and medium businesses (saunas, pools, sports facilities, conference rooms). Built with Go backend, PostgreSQL, MinIO (S3-compatible storage), and vanilla JavaScript frontend.

## Core Architecture

### Three-Layer Architecture Pattern

The codebase follows a clean three-layer architecture:

1. **Handler Layer** (`internal/handler/`) - HTTP request/response handling
   - Parses HTTP requests and validates input
   - Calls service layer for business logic
   - Formats responses and handles errors
   - Example: `auth_handler.go`, `resource_handler.go`, `booking_handler.go`

2. **Service Layer** (`internal/service/`) - Business logic
   - Implements business rules and validation
   - Orchestrates multiple repository calls if needed
   - Example: `booking_service.go` checks availability before creating bookings

3. **Repository Layer** (`internal/repository/`) - Data access
   - Direct database operations using `database/sql`
   - SQL queries and result mapping
   - Returns domain models defined in `internal/models/`

**Important**: When modifying existing functionality, maintain this layering:
- Handlers should NOT contain business logic
- Services should NOT know about HTTP
- Repositories should NOT contain business validation

### Domain Models

All data structures are in `internal/models/`:
- `User`, `Resource`, `Booking`, `Photo`, `Review`, `Category`
- Each model includes request/response types (e.g., `ResourceCreateRequest`, `ResourceUpdateRequest`)
- Database nullable fields use pointers (e.g., `*int64`, `*float64`)

### Storage System

Photo storage is abstracted via `internal/storage/storage.go`:
- Interface supports local filesystem or S3-compatible storage (MinIO)
- Configured via `STORAGE_TYPE` environment variable
- Default: MinIO at `http://minio:9000` with bucket `smartbooking`
- Photo URLs are stored in database, actual files in storage

## Common Development Commands

### Docker-based Development (Recommended)

```bash
# Start all services (PostgreSQL, MinIO, Backend, Nginx, pgAdmin)
docker-compose up -d

# View application logs
docker-compose logs -f app

# View database logs
docker-compose logs -f postgres

# Stop all services
docker-compose down

# Clean slate (remove all data)
docker-compose down -v

# Rebuild and restart
docker-compose up -d --build

# Access PostgreSQL shell
docker-compose exec postgres psql -U postgres -d smartbooking

# Access application shell
docker-compose exec app sh
```

### Makefile Commands

```bash
make up          # Start all services
make down        # Stop all services
make logs        # View app logs
make logs-db     # View database logs
make clean       # Remove containers and volumes
make db-shell    # Connect to PostgreSQL
make app-shell   # Connect to app container
make rebuild     # Clean rebuild
```

### Local Development (Backend only)

```bash
# Start dependencies only
docker-compose up -d postgres minio

# Build and run backend locally
go build -o smartbooking .
./smartbooking

# Or run directly
go run main.go

# Serve frontend (choose one):
cd frontend && python3 -m http.server 3000
# OR use VS Code Live Server extension on frontend/index.html
```

### Testing

```bash
# Run all tests
go test ./...

# Test specific package
go test ./internal/service

# Health check
curl http://localhost:8080/health

# Test authentication
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@smartbooking.com","password":"password123"}'
```

## Database Migrations

Migrations are in `migrations/` directory and run automatically on PostgreSQL container startup (via `docker-entrypoint-initdb.d`).

Migration files (in order):
1. `001_create_tables.sql` - Core schema
2. `002_add_sessions_and_audit.sql` - Sessions and audit tables
3. `003_seed_data.sql` - Test users and initial data
4. `004_add_photos_and_categories.sql` - Photos and categories schema
5. `005_seed_categories_and_prices.sql` - Category data
6. `006_add_owner_to_resources.sql` - Owner relationships
7. `007_add_test_users.sql` - Additional test users

**Important**: Migrations only run once when database is empty. If you need to re-run:
```bash
docker-compose down -v  # Delete volumes
docker-compose up -d    # Recreate and run migrations
```

## Configuration

Configuration is managed via environment variables in `config/config.go`:

### Database
- `DB_HOST` (default: `localhost`)
- `DB_PORT` (default: `5432`)
- `DB_USER` (default: `postgres`)
- `DB_PASSWORD` (default: `postgres`)
- `DB_NAME` (default: `smartbooking`)

### Server
- `SERVER_HOST` (default: `0.0.0.0`)
- `SERVER_PORT` (default: `8080`)

### Storage (MinIO/S3)
- `STORAGE_TYPE` (default: `minio`, can be `local` or `s3`)
- `STORAGE_ENDPOINT` (default: `http://minio:9000`)
- `STORAGE_ACCESS_KEY` (default: `minioadmin`)
- `STORAGE_SECRET_KEY` (default: `minioadmin`)
- `STORAGE_BUCKET` (default: `smartbooking`)
- `STORAGE_REGION` (default: `us-east-1`)
- `STORAGE_USE_SSL` (default: `false`)
- `STORAGE_PUBLIC_URL` (default: `http://localhost:9000`)

## Service URLs

When running with `docker-compose up`:
- **Frontend**: http://localhost (nginx)
- **Backend API**: http://localhost:8080
- **Swagger Docs**: http://localhost:8080/swagger/
- **MinIO Console**: http://localhost:9001 (minioadmin/minioadmin)
- **pgAdmin**: http://localhost:5050 (admin@smartbooking.com/admin)

## Test Credentials

Default accounts created by migrations:
- **Admin**: `admin@smartbooking.com` / `password123`
- **User**: `john@example.com` / `password123`

## API Structure

35 RESTful endpoints organized by domain:

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login (returns user object)

### Resources
- `GET /api/resources` - List all resources
- `POST /api/resources` - Create resource
- `GET /api/resources/{id}` - Get resource details
- `DELETE /api/resources/{id}` - Delete resource
- `GET /api/resources/{id}/photos` - Get resource photos
- `GET /api/resources/{id}/reviews` - Get resource reviews
- `GET /api/resources/{id}/rating` - Get average rating

### Bookings
- `GET /api/bookings` - List all bookings
- `POST /api/bookings` - Create booking (checks availability)
- `GET /api/bookings/{id}` - Get booking by ID
- `POST /api/bookings/{id}/cancel` - Cancel booking

### Photos
- `POST /api/photos/upload` - Upload photo to MinIO
- `DELETE /api/photos/{id}` - Delete photo
- `PUT /api/photos/{id}/primary` - Set as primary photo

### Reviews
- `GET /api/reviews` - List reviews
- `POST /api/reviews` - Create review
- `GET /api/reviews/{id}` - Get review
- `PUT /api/reviews/{id}` - Update review
- `DELETE /api/reviews/{id}` - Delete review

### Categories
- `GET /api/categories` - List categories
- `POST /api/categories` - Create category
- `GET /api/categories/{id}` - Get category
- `PUT /api/categories/{id}` - Update category
- `DELETE /api/categories/{id}` - Delete category

### Owners (Dashboard)
- `GET /api/owners/{id}/resources` - Get owner's resources
- `GET /api/owners/{id}/bookings` - Get owner's bookings
- `GET /api/owners/{id}/statistics` - Get owner statistics

### System
- `GET /health` - Health check
- `GET /swagger/` - API documentation

## Frontend Structure

Vanilla JavaScript frontend in `frontend/`:
- `index.html` - Homepage with categories
- `auth.html` - Login/Registration
- `resources.html` - Resource listing with search/filters
- `bookings.html` - User's bookings management
- `owner-dashboard.html` - Owner dashboard
- `admin.html` - Admin panel
- `js/app.js` - Main JavaScript (API calls, auth, navigation)
- `css/style.css` - Styling

**Authentication**: Simple user object stored in localStorage (not production-ready, uses user ID as token)

## Important Implementation Notes

### Adding New Features

When adding new endpoints/features:

1. **Define the model** in `internal/models/` (if new entity)
2. **Create repository interface and implementation** in `internal/repository/`
3. **Implement service layer** in `internal/service/` with business logic
4. **Add handler** in `internal/handler/` for HTTP
5. **Register route** in `main.go` (around line 130-170)
6. **Add Swagger annotations** to handler for documentation

### Database Operations

- Use `context.Context` for all DB operations
- Use parameterized queries (e.g., `$1, $2`) to prevent SQL injection
- Handle `sql.Null*` types with helper functions in models
- PostgreSQL arrays (amenities) use custom parsing in `models/resource.go`

### Error Handling

- Repository layer returns domain-specific errors (e.g., `ErrResourceNotFound`)
- Service layer adds business logic validation
- Handler layer converts to appropriate HTTP status codes
- Use `internal/logger` for logging

### CORS

CORS middleware in `main.go` allows all origins (`*`) - adjust for production:
```go
w.Header().Set("Access-Control-Allow-Origin", "*")
```

### Background Worker

`startStatisticsWorker()` in `main.go` logs statistics every 30 seconds (lines 212-235). Example of background goroutine pattern.

## Swagger Documentation

API documentation generated using swaggo:
- Annotations in handler files (e.g., `// @Summary`, `// @Param`)
- Generated docs in `docs/` directory
- Regenerate with: `swag init` (requires swag CLI)
- Access at: http://localhost:8080/swagger/

## Troubleshooting

### Port conflicts
```bash
lsof -i :8080  # Check what's using port 8080
lsof -i :5432  # Check PostgreSQL port
```

### Database issues
```bash
# Verify tables exist
docker-compose exec postgres psql -U postgres -d smartbooking -c "\dt"

# Check database stats
docker-compose logs app | grep "БД статистика"
```

### MinIO/Storage issues
- Access console at http://localhost:9001
- Ensure bucket `smartbooking` exists and has public read policy
- Check `STORAGE_PUBLIC_URL` matches how frontend accesses it

### Migration issues
Migrations only run once on database initialization. To force re-run:
```bash
docker-compose down -v && docker-compose up -d
```
