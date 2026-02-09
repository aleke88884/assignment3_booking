# SmartBooking - Quick Start Guide

## Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)
- Make (optional)

## Quick Start (Recommended)

### 1. Start Everything with Docker
```bash
# Start all services (PostgreSQL, MinIO, Backend, Nginx)
docker-compose up -d

# Check if all services are running
docker-compose ps

# View logs
docker-compose logs -f
```

**Access the application:**
- Frontend: http://localhost
- Backend API: http://localhost:8080
- Swagger Docs: http://localhost:8080/swagger/
- MinIO Console: http://localhost:9001 (admin/minioadmin)

### 2. Stop Everything
```bash
docker-compose down

# Stop and remove volumes (clean slate)
docker-compose down -v
```

---

## Local Development (Backend Only)

### 1. Start Dependencies
```bash
# Start only PostgreSQL and MinIO
docker-compose up -d postgres minio
```

### 2. Build and Run Backend
```bash
# Build the binary
go build -o smartbooking .

# Run the backend
./smartbooking
```

Or use Go directly:
```bash
go run main.go
```

### 3. Open Frontend
Option A: Use a local web server
```bash
cd frontend
python3 -m http.server 3000
# Visit http://localhost:3000
```

Option B: Use VS Code Live Server extension
- Right-click on `frontend/index.html` → "Open with Live Server"

---

## Building

### Build Go Binary
```bash
# Development build
go build -o smartbooking .

# Production build (optimized)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smartbooking .

# Run the binary
./smartbooking
```

### Build Docker Image
```bash
# Build the image
docker build -t smartbooking:latest .

# Run the container
docker run -p 8080:8080 smartbooking:latest
```

### Build with Make (if you have Make)
```bash
# Build binary
make build

# Run with Docker Compose
make up

# Stop Docker Compose
make down

# View logs
make logs

# Clean everything
make clean
```

---

## Testing

### 1. Test Backend Health
```bash
curl http://localhost:8080/health
# Expected: {"status": "ok"}
```

### 2. Test Registration
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "test123456"
  }'
```

### 3. Test Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@smartbooking.com",
    "password": "password123"
  }'
```

### 4. Test Get Resources
```bash
curl http://localhost:8080/api/resources
```

### 5. Run Go Tests
```bash
go test ./...
```

---

## Default Test Accounts

After migrations run, these accounts are available:

**Admin Account:**
- Email: `admin@smartbooking.com`
- Password: `password123`

**User Account:**
- Email: `john@example.com`
- Password: `password123`

---

## Environment Variables

Create a `.env` file in the project root (optional):

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=smartbooking

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Storage (MinIO)
STORAGE_TYPE=minio
STORAGE_ENDPOINT=http://localhost:9000
STORAGE_ACCESS_KEY=minioadmin
STORAGE_SECRET_KEY=minioadmin
STORAGE_BUCKET=smartbooking
STORAGE_REGION=us-east-1
STORAGE_USE_SSL=false
STORAGE_PUBLIC_URL=http://localhost:9000
```

---

## Troubleshooting

### Backend won't start
```bash
# Check if ports are in use
lsof -i :8080
lsof -i :5432

# Check database connection
docker-compose logs postgres

# Rebuild everything
docker-compose down -v
docker-compose up --build
```

### Database is empty
```bash
# Check if migrations ran
docker-compose exec postgres psql -U postgres -d smartbooking -c "\dt"

# Manually run migrations
docker-compose exec postgres psql -U postgres -d smartbooking < migrations/001_create_tables.sql
docker-compose exec postgres psql -U postgres -d smartbooking < migrations/003_seed_data.sql
```

### Frontend can't connect to backend
- Check if backend is running: `curl http://localhost:8080/health`
- Check browser console for errors
- Verify API_URL in `frontend/js/app.js`

### MinIO not working
```bash
# Access MinIO console
http://localhost:9001
# Login: minioadmin / minioadmin

# Check if bucket exists
# Create bucket named "smartbooking" if missing
```

---

## Project Structure

```
assignment3/
├── main.go                 # Entry point
├── config/                 # Configuration
├── internal/
│   ├── handler/           # HTTP handlers (35 endpoints)
│   ├── service/           # Business logic
│   ├── repository/        # Database layer
│   ├── models/            # Data structures
│   ├── storage/           # File storage (S3/MinIO)
│   └── database/          # DB connection & migrations
├── migrations/            # SQL migration files
├── frontend/              # Frontend files
│   ├── index.html
│   ├── auth.html
│   ├── resources.html
│   ├── bookings.html
│   ├── js/app.js
│   └── css/style.css
├── nginx/                 # Nginx config
├── Dockerfile             # Backend container
├── docker-compose.yml     # Full stack setup
└── go.mod                 # Go dependencies
```

---

## API Endpoints (35 total)

| Method | Endpoint | Description |
|--------|----------|-------------|
| **Auth** |
| POST | `/api/auth/register` | Register new user |
| POST | `/api/auth/login` | User login |
| **Users** |
| GET | `/api/users` | List all users |
| GET | `/api/users/{id}` | Get user by ID |
| GET | `/api/users/{id}/bookings` | Get user's bookings |
| GET | `/api/users/{user_id}/reviews` | Get user's reviews |
| **Resources** |
| GET | `/api/resources` | List all resources |
| POST | `/api/resources` | Create resource |
| GET | `/api/resources/{id}` | Get resource by ID |
| DELETE | `/api/resources/{id}` | Delete resource |
| GET | `/api/resources/{id}/photos` | Get resource photos |
| GET | `/api/resources/{id}/reviews` | Get resource reviews |
| GET | `/api/resources/{id}/rating` | Get average rating |
| **Bookings** |
| GET | `/api/bookings` | List all bookings |
| POST | `/api/bookings` | Create booking |
| GET | `/api/bookings/{id}` | Get booking by ID |
| POST | `/api/bookings/{id}/cancel` | Cancel booking |
| **Photos** |
| POST | `/api/photos/upload` | Upload photo |
| DELETE | `/api/photos/{id}` | Delete photo |
| PUT | `/api/photos/{id}/primary` | Set primary photo |
| **Reviews** |
| GET | `/api/reviews` | List reviews |
| POST | `/api/reviews` | Create review |
| GET | `/api/reviews/{id}` | Get review by ID |
| PUT | `/api/reviews/{id}` | Update review |
| DELETE | `/api/reviews/{id}` | Delete review |
| **Categories** |
| GET | `/api/categories` | List categories |
| POST | `/api/categories` | Create category |
| GET | `/api/categories/{id}` | Get category by ID |
| PUT | `/api/categories/{id}` | Update category |
| DELETE | `/api/categories/{id}` | Delete category |
| **System** |
| GET | `/health` | Health check |
| GET | `/swagger/` | API documentation |

---

## Quick Commands Cheatsheet

```bash
# Start everything
docker-compose up -d

# Rebuild and start
docker-compose up -d --build

# View logs
docker-compose logs -f app

# Stop everything
docker-compose down

# Clean slate (remove volumes)
docker-compose down -v

# Check service status
docker-compose ps

# Enter postgres container
docker-compose exec postgres psql -U postgres -d smartbooking

# Restart single service
docker-compose restart app

# Build Go binary
go build -o smartbooking .

# Run Go binary
./smartbooking

# Test API
curl http://localhost:8080/health
```

That's it! Your SmartBooking application should now be running at http://localhost
