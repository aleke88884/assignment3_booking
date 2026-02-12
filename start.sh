#!/bin/bash

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🚀 SmartBooking - Starting..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ ERROR: Docker is not running!"
    echo ""
    echo "Please start Docker Desktop and try again."
    exit 1
fi

echo "✅ Docker is running"
echo ""

# Stop and remove old containers
echo "🧹 Cleaning up old containers..."
docker-compose down -v 2>/dev/null
echo ""

# Build and start services
echo "🔨 Building and starting services..."
echo "   This may take 1-2 minutes on first run..."
echo ""
docker-compose up -d --build

# Wait for services to be ready
echo ""
echo "⏳ Waiting for services to start..."
sleep 15

# Check if containers are running
echo ""
echo "📊 Container status:"
docker-compose ps
echo ""

# Wait a bit more for migrations
echo "⏳ Waiting for database migrations..."
sleep 10

# Check database
DB_CHECK=$(docker-compose exec -T postgres psql -U postgres -d smartbooking -c "SELECT COUNT(*) FROM users;" 2>/dev/null | grep -o '[0-9]\+' | head -1)

if [ ! -z "$DB_CHECK" ] && [ "$DB_CHECK" -gt 0 ]; then
    echo "✅ Database ready - $DB_CHECK users found"
else
    echo "⚠️  Database might still be initializing..."
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ SmartBooking is ready!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🌐 Open in browser:"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8080"
echo "   API Docs: http://localhost:8080/swagger/"
echo ""
echo "👤 Test accounts (password: password123):"
echo "   Admin:  admin@smartbooking.com"
echo "   Owner:  owner1@smartbooking.com"
echo "   User:   john@example.com"
echo ""
echo "📝 Useful commands:"
echo "   View logs:    docker-compose logs -f"
echo "   View app:     docker-compose logs -f app"
echo "   Stop all:     docker-compose down"
echo "   Restart:      ./start.sh"
echo ""
