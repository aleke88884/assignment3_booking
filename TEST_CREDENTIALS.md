# Test Credentials for SmartBooking

## Quick Reference

All test accounts use the password: **`password`**

For simplified login, you can use just the username/email field with these values:

---

## Admin Accounts

### Admin User
- **Email/Username**: `admin`
- **Password**: `password`
- **Role**: admin
- **Access**: Full system access, admin panel

### Super Admin
- **Email/Username**: `superadmin`
- **Password**: `password`
- **Role**: admin
- **Access**: Full system access, admin panel

---

## Owner Accounts (Resource Owners)

### Owner 1
- **Email/Username**: `owner1`
- **Password**: `password`
- **Role**: user
- **Resources**: 3 resources
  - Downtown Office Suite
  - Rooftop Conference Room
  - Creative Workshop Space
- **Access**: Owner dashboard, can view their bookings and statistics

### Owner 2
- **Email/Username**: `owner2`
- **Password**: `password`
- **Role**: user
- **Resources**: 3 resources
  - Beachside Villa
  - Mountain Cabin Retreat
  - Urban Loft Apartment
- **Access**: Owner dashboard, can view their bookings and statistics

### Owner 3
- **Email/Username**: `owner3`
- **Password**: `password`
- **Role**: user
- **Resources**: 4 resources
  - Professional Recording Studio
  - Photography Studio
  - Event Hall
  - Small Meeting Pod
- **Access**: Owner dashboard, can view their bookings and statistics

---

## Regular User Accounts

### User 1
- **Email/Username**: `user1`
- **Password**: `password`
- **Role**: user
- **Bookings**: 2 bookings on owner1's resources

### User 2
- **Email/Username**: `user2`
- **Password**: `password`
- **Role**: user
- **Bookings**: 2 bookings on owner1's resources

### User 3
- **Email/Username**: `user3`
- **Password**: `password`
- **Role**: user
- **Bookings**: 2 bookings on owner2's resources

### User 4
- **Email/Username**: `user4`
- **Password**: `password`
- **Role**: user
- **Bookings**: None yet

### User 5
- **Email/Username**: `user5`
- **Password**: `password`
- **Role**: user
- **Bookings**: None yet

---

## How to Use

### Logging In

1. Go to the login page (or use the API)
2. Enter email: `admin` (or any username above)
3. Enter password: `password`
4. Click Login

### Testing Different Features

**To test Admin Panel:**
```
Email: admin
Password: password
Visit: http://localhost/admin.html
```

**To test Owner Dashboard:**
```
Email: owner1 (or owner2, owner3)
Password: password
Visit: http://localhost/owner-dashboard.html
```

**To test User Bookings:**
```
Email: user1 (or user2, user3)
Password: password
Visit: http://localhost/bookings.html
```

---

## Test Data Summary

### Resources Created
- **10 test resources** across 3 owners
- **Various categories**: Offices, Apartments, Studios, Event Spaces
- **Different price ranges**: $25/hr to $200/hr
- **Multiple locations**: New York, Miami, Denver, Chicago, LA, Las Vegas, San Francisco

### Bookings Created
- **6 test bookings** with different statuses (confirmed, pending)
- **Price range**: $150 to $7,200
- **Various durations**: 2 hours to 48 hours

### Reviews Created
- **3 test reviews** with ratings 4-5 stars
- **Detailed comments** for realism

---

## Database Admin Access

### pgAdmin
- **URL**: http://localhost:5050
- **Email**: `admin@smartbooking.com`
- **Password**: `admin`

**To connect to database:**
- Host: `postgres`
- Port: `5432`
- Database: `smartbooking`
- Username: `postgres`
- Password: `postgres`

### MinIO (File Storage)
- **URL**: http://localhost:9001
- **Username**: `minioadmin`
- **Password**: `minioadmin`

---

## API Testing

You can test the API endpoints with these credentials using curl:

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin","password":"password"}'
```

### Get Owner Resources (as owner1)
```bash
# First login to get user ID, then:
curl http://localhost:8080/api/owners/1/resources
```

### Get Owner Statistics (as owner1)
```bash
curl http://localhost:8080/api/owners/1/statistics
```

---

## Quick Start Guide

1. **Start the application:**
   ```bash
   docker-compose up -d
   ```

2. **Wait for migrations to run** (about 10-30 seconds)

3. **Test admin access:**
   - Go to: http://localhost/admin.html
   - Login: `admin` / `password`
   - Explore the admin panel

4. **Test owner dashboard:**
   - Go to: http://localhost/owner-dashboard.html
   - Login: `owner1` / `password`
   - View your 3 resources and their bookings

5. **Test regular user:**
   - Go to: http://localhost
   - Login: `user1` / `password`
   - Browse and book resources

---

## Resetting Test Data

If you need to reset the test data:

```bash
# Stop and remove containers
docker-compose down -v

# Start fresh
docker-compose up -d
```

The migrations will automatically re-create all test data.

---

## Security Note

⚠️ **IMPORTANT**: These credentials are for **development and testing only**!

- All test accounts use the same simple password: `password`
- Never use these credentials in production
- Always use strong, unique passwords in production
- The bcrypt hash used is: `$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy`

---

## Summary Table

| Email | Password | Role | Resources | Bookings | Purpose |
|-------|----------|------|-----------|----------|---------|
| admin | password | admin | - | - | Full admin access |
| superadmin | password | admin | - | - | Full admin access |
| owner1 | password | user | 3 | 2 (incoming) | Test owner dashboard |
| owner2 | password | user | 3 | 1 (incoming) | Test owner dashboard |
| owner3 | password | user | 4 | 0 | Test owner dashboard |
| user1 | password | user | - | 2 | Test user bookings |
| user2 | password | user | - | 2 | Test user bookings |
| user3 | password | user | - | 2 | Test user bookings |
| user4 | password | user | - | 0 | Test new user |
| user5 | password | user | - | 0 | Test new user |

---

## Troubleshooting

**Can't login?**
- Make sure migrations have run: `docker-compose logs postgres`
- Check if users exist: `docker-compose exec postgres psql -U postgres -d smartbooking -c "SELECT email, role FROM users;"`

**No resources showing?**
- Check if migrations completed: `docker-compose logs app`
- Verify resources exist: Visit http://localhost:8080/api/resources

**Owner dashboard empty?**
- Make sure you're logged in as owner1, owner2, or owner3
- Check if resources have owner_id set: Use pgAdmin to verify

---

## Need More Test Data?

You can easily add more test data by:
1. Using the admin panel to create resources
2. Using the frontend to make bookings
3. Adding more users through the registration page
4. Manually inserting data via pgAdmin

Enjoy testing! 🚀
