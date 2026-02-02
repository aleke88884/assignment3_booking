# Assignment 4 - SmartBooking Core System Implementation

## Team Information
**Project**: SmartBooking - Online Booking Management System

## Project Overview
SmartBooking is an online booking management system for small businesses (hotels, coworking spaces, meeting rooms, sport facilities, etc.). This submission implements the core backend system with working endpoints, data persistence, and concurrency.

## How to Run

### Prerequisites
- Go 1.25.3 or higher
- No database required (uses in-memory storage)

### Start the Server
```bash
# Option 1: Direct run
go run .

# Option 2: Build and run
go build -o smartbooking .
./smartbooking
```

The server will start on `http://localhost:8080`

### Run Tests
```bash
# Make sure server is running first
./test_api.sh
```

## Implemented Features

### 1. Authentication & User Management
- ✅ User registration with password hashing (bcrypt)
- ✅ User login with credential verification
- ✅ List all users
- ✅ Get user by ID

### 2. Resource Management
- ✅ Create resources (rooms, facilities, etc.)
- ✅ List all resources
- ✅ Get resource by ID
- ✅ Delete resources

### 3. Booking Management
- ✅ Create bookings with validation
- ✅ **Double booking prevention** - checks for overlapping reservations
- ✅ Cancel bookings
- ✅ List all bookings
- ✅ List bookings by user

### 4. Background Processing
- ✅ Statistics worker running in goroutine
- ✅ Periodic logging of system statistics (every 30 seconds)

### 5. Data Persistence
- ✅ In-memory storage with thread-safe operations
- ✅ CRUD operations for all entities
- ✅ Concurrent request handling with sync.RWMutex

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login

### Users
- `GET /api/users` - List all users
- `GET /api/users/{id}` - Get user by ID
- `GET /api/users/{id}/bookings` - Get user's bookings

### Resources
- `GET /api/resources` - List all resources
- `POST /api/resources` - Create resource
- `GET /api/resources/{id}` - Get resource by ID
- `DELETE /api/resources/{id}` - Delete resource

### Bookings
- `GET /api/bookings` - List all bookings
- `POST /api/bookings` - Create booking
- `GET /api/bookings/{id}` - Get booking by ID
- `POST /api/bookings/{id}/cancel` - Cancel booking

### Health
- `GET /health` - Health check

## Example API Calls

### Register a User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com", "password": "password123"}'
```

### Create a Resource
```bash
curl -X POST http://localhost:8080/api/resources \
  -H "Content-Type: application/json" \
  -d '{"name": "Conference Room A", "description": "Large meeting room", "capacity": 20}'
```

### Create a Booking
```bash
curl -X POST http://localhost:8080/api/bookings \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "resource_id": 1, "start_time": "2026-02-10T10:00:00Z", "end_time": "2026-02-10T11:00:00Z"}'
```

### Test Double Booking Prevention
```bash
# Try to create overlapping booking (will fail)
curl -X POST http://localhost:8080/api/bookings \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "resource_id": 1, "start_time": "2026-02-10T10:30:00Z", "end_time": "2026-02-10T11:30:00Z"}'
```

## Technical Implementation

### Architecture
- **Layered Architecture**: Handlers → Services → Repositories
- **Monolithic**: Single application (as required)
- **RESTful API**: JSON request/response format

### Data Models
Following the ERD from Assignment 3:
- **User**: id, name, email, password (hashed), role, timestamps
- **Resource**: id, name, description, capacity, timestamps
- **Booking**: id, user_id, resource_id, start_time, end_time, status, timestamps

### Concurrency Implementation
1. **Thread-safe repositories**: Using `sync.RWMutex` for all data operations
2. **Background worker**: Statistics tracking goroutine
3. **Safe concurrent requests**: Multiple clients can access API simultaneously

### Key Features
- **Password Security**: bcrypt hashing with cost factor 10
- **Conflict Prevention**: Checks for overlapping bookings before creation
- **Status Management**: Pending → Confirmed/Cancelled transitions
- **Role-based Data**: User/Admin roles (foundation for future auth)

## Assignment 4 Requirements Checklist

### Backend Application (30%)
✅ HTTP server using net/http
✅ 10+ working endpoints
✅ JSON input and output

### Data Model & Storage (25%)
✅ Data structures match Assignment 3 ERD
✅ CRUD operations for all core entities
✅ Safe data access (thread-safe with mutexes)

### Concurrency (15%)
✅ Background worker goroutine
✅ Thread-safe repositories with sync.RWMutex

### Core Features (5 implemented)
✅ User registration
✅ User authentication
✅ Resource management
✅ Booking creation
✅ Booking cancellation with conflict prevention

## Project Structure
```
assignment3/
├── main.go                          # Server setup and background worker
├── go.mod                           # Dependencies
├── config/
│   └── config.go                   # Configuration management
├── internal/
│   ├── models/                     # Data models
│   │   ├── user.go
│   │   ├── resource.go
│   │   └── booking.go
│   ├── repository/                 # Data access layer (in-memory)
│   │   ├── user_repository.go
│   │   ├── resource_repository.go
│   │   └── booking_repository.go
│   ├── service/                    # Business logic
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── resource_service.go
│   │   └── booking_service.go
│   └── handler/                    # HTTP handlers
│       ├── auth_handler.go
│       ├── user_handler.go
│       ├── resource_handler.go
│       └── booking_handler.go
├── test_api.sh                     # API test script
├── ASSIGNMENT4_README.md           # This file
└── IMPLEMENTATION_NOTES.md         # Implementation details and originality

```

## Dependencies
- `golang.org/x/crypto` - For bcrypt password hashing

## Testing
All endpoints have been tested and verified working:
- User registration and login ✅
- Resource CRUD operations ✅
- Booking creation and cancellation ✅
- Double booking prevention ✅
- Concurrent request handling ✅
- Background worker statistics ✅

## Notes for Defense

### What's Implemented
- Full working backend with all core features
- In-memory storage (as allowed for milestone)
- Thread-safe concurrent operations
- Background goroutine for statistics
- Comprehensive error handling

### What's NOT Implemented (by design)
- Full authentication/authorization (JWT, sessions) - Not required yet
- Database persistence - In-memory sufficient for milestone
- Complete error handling for all edge cases - Will be addressed in final
- UI/Frontend - Not required for backend milestone

### Following Assignment 3
The implementation strictly follows the approved architecture from Assignment 3:
- Same entity structure (User, Resource, Booking)
- Same layered architecture
- Same business rules (conflict prevention, status management)

## Future Improvements (Final Project)
- Database integration (PostgreSQL)
- JWT-based authentication
- Full role-based authorization
- Advanced search and filtering
- Email notifications
- Frontend interface

## How to Demo

1. **Start the server**: `go run .`
2. **Show endpoints list**: Server logs show all available endpoints
3. **Run test script**: `./test_api.sh` demonstrates all features
4. **Show conflict prevention**: Try creating overlapping bookings
5. **Show background worker**: Server logs show periodic statistics
6. **Show concurrent safety**: Multiple curl requests work simultaneously

## Conclusion

This implementation meets all Assignment 4 requirements:
- ✅ Working backend application
- ✅ Core domain models matching ERD
- ✅ 5+ core features implemented
- ✅ In-memory persistence with thread safety
- ✅ Concurrency with goroutines
- ✅ Clean code following Go best practices

The system is ready for demonstration and provides a solid foundation for the final project submission.
