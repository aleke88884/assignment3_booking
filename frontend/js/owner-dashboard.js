// Owner Dashboard JavaScript

const API_BASE_URL = (window.location.port === '80' || window.location.port === '3000' || window.location.port === '')
    ? '/api'
    : 'http://localhost:8080/api';

// Check if user is logged in
function checkAuth() {
    const user = localStorage.getItem('user');
    if (!user) {
        window.location.href = '/auth.html';
        return null;
    }
    return JSON.parse(user);
}

// Format currency
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(amount);
}

// Format date
function formatDate(dateString) {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    }).format(date);
}

// Load owner statistics
async function loadStatistics(ownerId) {
    try {
        const response = await fetch(`${API_BASE_URL}/owners/${ownerId}/statistics`);
        if (!response.ok) throw new Error('Failed to load statistics');

        const stats = await response.json();

        document.getElementById('totalResources').textContent = stats.total_resources || 0;
        document.getElementById('totalBookings').textContent = stats.total_bookings || 0;
        document.getElementById('activeBookings').textContent = stats.active_bookings || 0;
        document.getElementById('totalRevenue').textContent = formatCurrency(stats.total_revenue || 0);
        document.getElementById('averageRating').textContent = (stats.average_rating || 0).toFixed(1);
        document.getElementById('totalReviews').textContent = stats.total_reviews || 0;

        document.getElementById('loadingStats').style.display = 'none';
        document.getElementById('statsContainer').style.display = 'block';
    } catch (error) {
        console.error('Error loading statistics:', error);
        document.getElementById('loadingStats').textContent = 'Failed to load statistics';
    }
}

// Load owner resources
async function loadResources(ownerId) {
    try {
        const response = await fetch(`${API_BASE_URL}/owners/${ownerId}/resources`);
        if (!response.ok) throw new Error('Failed to load resources');

        const resources = await response.json();

        document.getElementById('loadingResources').style.display = 'none';

        if (!resources || resources.length === 0) {
            document.getElementById('noResources').style.display = 'block';
            return;
        }

        const container = document.getElementById('resourcesContainer');
        container.innerHTML = '';

        resources.forEach(resource => {
            const card = document.createElement('div');
            card.className = 'resource-card';

            const rating = resource.rating ? resource.rating.toFixed(1) : 'N/A';
            const reviewCount = resource.reviews_count || 0;
            const pricePerHour = resource.price_per_hour ? formatCurrency(resource.price_per_hour) : 'N/A';
            const status = resource.is_active ? 'Active' : 'Inactive';

            card.innerHTML = `
                <h3>${escapeHtml(resource.name)}</h3>
                <p class="resource-info"><strong>Category:</strong> ${escapeHtml(resource.category_name || 'N/A')}</p>
                <p class="resource-info"><strong>Capacity:</strong> ${resource.capacity} people</p>
                <p class="resource-info"><strong>Price:</strong> ${pricePerHour}/hour</p>
                <p class="resource-info"><strong>Location:</strong> ${escapeHtml(resource.city || 'N/A')}</p>
                <p class="resource-info"><strong>Status:</strong> ${status}</p>
                <p class="resource-rating">★ ${rating} (${reviewCount} reviews)</p>
            `;

            container.appendChild(card);
        });

        container.style.display = 'grid';
    } catch (error) {
        console.error('Error loading resources:', error);
        document.getElementById('loadingResources').textContent = 'Failed to load resources';
    }
}

// Load owner bookings
async function loadBookings(ownerId) {
    try {
        const response = await fetch(`${API_BASE_URL}/owners/${ownerId}/bookings`);
        if (!response.ok) throw new Error('Failed to load bookings');

        const bookings = await response.json();

        document.getElementById('loadingBookings').style.display = 'none';

        if (!bookings || bookings.length === 0) {
            document.getElementById('noBookings').style.display = 'block';
            return;
        }

        const tbody = document.getElementById('bookingsTableBody');
        tbody.innerHTML = '';

        bookings.forEach(booking => {
            const row = document.createElement('tr');

            const statusClass = `status-${booking.status}`;
            const price = booking.total_price ? formatCurrency(booking.total_price) : 'N/A';

            row.innerHTML = `
                <td>#${booking.id}</td>
                <td>${escapeHtml(booking.resource_name || 'N/A')}</td>
                <td>
                    ${escapeHtml(booking.user_name || 'N/A')}<br>
                    <small>${escapeHtml(booking.user_email || '')}</small>
                </td>
                <td>${formatDate(booking.start_time)}</td>
                <td>${formatDate(booking.end_time)}</td>
                <td><span class="status-badge ${statusClass}">${booking.status}</span></td>
                <td>${price}</td>
            `;

            tbody.appendChild(row);
        });

        document.getElementById('bookingsContainer').style.display = 'block';
    } catch (error) {
        console.error('Error loading bookings:', error);
        document.getElementById('loadingBookings').textContent = 'Failed to load bookings';
    }
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Navigation functions
function goToHome() {
    window.location.href = '/index.html';
}

function createNewResource() {
    // This would open a modal or redirect to a resource creation page
    alert('Resource creation form would open here. This feature needs to be implemented.');
    // Future: window.location.href = '/create-resource.html';
}

// Initialize dashboard
document.addEventListener('DOMContentLoaded', () => {
    const user = checkAuth();
    if (!user) return;

    // Display owner name
    document.getElementById('ownerName').textContent = user.name;

    // Load all dashboard data
    const ownerId = user.id;
    loadStatistics(ownerId);
    loadResources(ownerId);
    loadBookings(ownerId);
});
