# SmartBooking Frontend-Backend Integration Testing

## Prerequisites
- Docker and Docker Compose installed
- Ports 80, 5432, 8080, 9000, 9001 available

## Starting the Application

### Option 1: Full Stack (Recommended)
```bash
docker-compose up -d
```
Access the application at: http://localhost

### Option 2: Development Mode (Backend only)
```bash
# Start dependencies
docker-compose up -d postgres minio

# Run backend locally
go run main.go
```
Access backend at: http://localhost:8080
Open frontend files directly in browser (use Live Server extension)

## Testing Checklist

### 1. Homepage (index.html)
- [ ] Navigate to http://localhost (or open index.html)
- [ ] Categories should load automatically
- [ ] Navigation bar should show "Войти" button
- [ ] Clicking "Посмотреть ресурсы" redirects to resources page

### 2. Authentication (auth.html)
- [ ] Navigate to http://localhost/auth.html
- [ ] Try login with test account:
  - Email: admin@smartbooking.com
  - Password: password123
- [ ] Should redirect to index.html on success
- [ ] Navigation should show user name and "Выйти" button
- [ ] Try invalid credentials - should show error message

### 3. Registration
- [ ] Switch to "Регистрация" tab
- [ ] Register new user:
  - Name: Test User
  - Email: test@example.com
  - Password: test123456
- [ ] Should redirect to index.html on success
- [ ] Should show error if email already exists

### 4. Resources Page (resources.html)
- [ ] Navigate to http://localhost/resources.html
- [ ] Resources should load in grid
- [ ] Each resource should show:
  - Name, description, capacity
  - Photo (or placeholder)
  - Price per hour
  - "Забронировать" button
- [ ] Test search filter
- [ ] Test category filter
- [ ] Test city filter

### 5. Booking Modal
- [ ] Click "Забронировать" on any resource (must be logged in)
- [ ] Modal should open showing:
  - Resource name
  - Resource photos (if available)
  - Resource details (price, capacity, description)
- [ ] Click on photo thumbnail to view larger version
- [ ] Fill in booking details:
  - Start time: tomorrow at 10:00
  - End time: tomorrow at 12:00
  - Notes: "Test booking"
- [ ] Submit booking
- [ ] Should redirect to bookings page

### 6. Photo Gallery
- [ ] In booking modal, if resource has multiple photos
- [ ] Click on any photo thumbnail
- [ ] Photo modal should open with:
  - Large photo display
  - Thumbnail navigation at bottom
- [ ] Click different thumbnails to switch photos
- [ ] Close modal by clicking X or outside

### 7. My Bookings Page (bookings.html)
- [ ] Navigate to http://localhost/bookings.html (must be logged in)
- [ ] Should show list of user's bookings
- [ ] Each booking should display:
  - Resource ID
  - Start and end time (formatted)
  - Status badge (Подтверждено/В ожидании/Отменено)
  - Notes (if any)
  - "Отменить бронь" button (if not cancelled)
- [ ] Click "Отменить бронь"
- [ ] Confirm cancellation
- [ ] Booking status should update to "Отменено"

### 8. API Endpoints Testing

#### Using curl:
```bash
# Health check
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test2@test.com","password":"test123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@smartbooking.com","password":"password123"}'

# Get resources
curl http://localhost:8080/api/resources

# Get categories
curl http://localhost:8080/api/categories

# Get bookings (replace USER_ID)
curl http://localhost:8080/api/users/1/bookings
```

#### Swagger UI:
Navigate to: http://localhost:8080/swagger/

## Known Issues & Solutions

### Issue 1: CORS Errors
**Symptom:** Browser console shows CORS policy errors
**Solution:** Make sure the backend is running and CORS middleware is enabled

### Issue 2: Cannot Connect to API
**Symptom:** "Ошибка соединения с сервером" messages
**Solution:**
- Check if backend is running: `docker ps` or `curl http://localhost:8080/health`
- Check browser console for actual error
- Verify API_URL in app.js matches your setup

### Issue 3: Photos Not Loading
**Symptom:** Images show placeholder
**Solution:**
- Check MinIO is running: http://localhost:9001 (admin/admin)
- Verify bucket "smartbooking" exists
- Check STORAGE_PUBLIC_URL in docker-compose.yml

### Issue 4: Database Connection Failed
**Symptom:** "Failed to connect to database" in logs
**Solution:**
- Wait for postgres to be fully ready: `docker-compose logs postgres`
- Check migrations ran: `docker-compose exec postgres psql -U postgres -d smartbooking -c "\dt"`

### Issue 5: Empty Categories/Resources
**Symptom:** Pages load but no data
**Solution:** Database needs seeding:
```bash
docker-compose exec postgres psql -U postgres -d smartbooking -f /docker-entrypoint-initdb.d/003_seed_data.sql
```

## Integration Points Summary

### Frontend → Backend Communication
| Frontend Action | API Endpoint | Method | Authentication |
|----------------|--------------|--------|----------------|
| Register | `/api/auth/register` | POST | No |
| Login | `/api/auth/login` | POST | No |
| Load resources | `/api/resources` | GET | No |
| Load categories | `/api/categories` | GET | No |
| Create booking | `/api/bookings` | POST | Bearer Token |
| List user bookings | `/api/users/{id}/bookings` | GET | Bearer Token |
| Cancel booking | `/api/bookings/{id}/cancel` | POST | Bearer Token |
| Get resource photos | `/api/resources/{id}/photos` | GET | No |

### Data Flow
1. User registers/logs in → Backend returns user object
2. Frontend stores user in localStorage
3. Frontend uses user.id as "token" (simple auth)
4. All authenticated requests include user.id
5. Backend validates and processes requests

### Environment Differences
- **Docker (Production):** Frontend served by nginx on port 80, proxies to backend
- **Local (Development):** Frontend accesses backend directly at localhost:8080

## Success Criteria
✅ All pages load without errors
✅ Can register and login successfully
✅ Resources display with filters working
✅ Can create and view bookings
✅ Can cancel bookings
✅ Photo galleries work properly
✅ Error messages display correctly
✅ Navigation between pages works
✅ User session persists across page refreshes
