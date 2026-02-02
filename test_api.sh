#!/bin/bash

# SmartBooking API Test Script
# This script demonstrates all the core features implemented for Assignment 4

BASE_URL="http://localhost:8080"

echo "====================================="
echo "SmartBooking API Test Script"
echo "Assignment 4 - Core System Implementation"
echo "====================================="
echo ""

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s $BASE_URL/health | jq .
echo ""

# Test 2: User Registration
echo "2. Testing User Registration..."
echo "   Creating user: Alice Johnson"
curl -s -X POST $BASE_URL/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice Johnson", "email": "alice@example.com", "password": "alice123"}' | jq .
echo ""

# Test 3: User Login
echo "3. Testing User Login..."
echo "   Logging in as Alice"
curl -s -X POST $BASE_URL/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "alice@example.com", "password": "alice123"}' | jq .
echo ""

# Test 4: Create Resources
echo "4. Testing Resource Creation..."
echo "   Creating resource: Sauna"
curl -s -X POST $BASE_URL/api/resources \
  -H "Content-Type: application/json" \
  -d '{"name": "Sauna", "description": "Relaxing sauna with capacity for 8 people", "capacity": 8}' | jq .
echo ""

# Test 5: List Resources
echo "5. Testing Resource Listing..."
curl -s $BASE_URL/api/resources | jq .
echo ""

# Test 6: Create Booking
echo "6. Testing Booking Creation..."
echo "   Creating booking for Sauna"
curl -s -X POST $BASE_URL/api/bookings \
  -H "Content-Type: application/json" \
  -d '{"user_id": 3, "resource_id": 4, "start_time": "2026-02-20T18:00:00Z", "end_time": "2026-02-20T19:00:00Z"}' | jq .
echo ""

# Test 7: Test Double Booking Prevention
echo "7. Testing Double Booking Prevention..."
echo "   Attempting to create overlapping booking (should fail)"
curl -s -X POST $BASE_URL/api/bookings \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "resource_id": 4, "start_time": "2026-02-20T18:30:00Z", "end_time": "2026-02-20T19:30:00Z"}'
echo ""
echo ""

# Test 8: List All Bookings
echo "8. Testing Booking Listing..."
curl -s $BASE_URL/api/bookings | jq .
echo ""

# Test 9: List Users
echo "9. Testing User Listing..."
curl -s $BASE_URL/api/users | jq .
echo ""

# Test 10: Get User's Bookings
echo "10. Testing User's Bookings..."
echo "    Getting bookings for user ID 2"
curl -s $BASE_URL/api/users/2/bookings | jq .
echo ""

echo "====================================="
echo "All tests completed!"
echo "====================================="
