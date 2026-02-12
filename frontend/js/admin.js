// Admin Panel JavaScript

const API_BASE_URL = (window.location.port === '80' || window.location.port === '3000' || window.location.port === '')
    ? '/api'
    : 'http://localhost:8080/api';

let allBookings = [];
let allResources = [];
let allUsers = [];
let allCategories = [];
document.addEventListener('DOMContentLoaded', function () {
    initCharts();
});

function initCharts() {
    // 1. Line Chart - Bookings over time
    const lineCtx = document.getElementById('lineChart').getContext('2d');
    new Chart(lineCtx, {
        type: 'line',
        data: {
            labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
            datasets: [{
                label: 'Bookings',
                data: [12, 19, 15, 25, 32, 45, 40],
                borderColor: '#667eea',
                backgroundColor: 'rgba(102, 126, 234, 0.1)',
                fill: true,
                tension: 0.4,
                borderWidth: 3,
                pointBackgroundColor: '#fff',
                pointBorderColor: '#667eea',
                pointRadius: 5
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: { display: false }
            },
            scales: {
                y: { beginAtZero: true, grid: { display: false } },
                x: { grid: { display: false } }
            }
        }
    });

    // 2. Doughnut Chart - Category distribution
    const doughnutCtx = document.getElementById('doughnutChart').getContext('2d');
    new Chart(doughnutCtx, {
        type: 'doughnut',
        data: {
            labels: ['Offices', 'Venues', 'Equipment', 'Cars'],
            datasets: [{
                data: [40, 25, 20, 15],
                backgroundColor: [
                    '#667eea',
                    '#48bb78',
                    '#ed8936',
                    '#4299e1'
                ],
                borderWidth: 0,
                hoverOffset: 10
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom',
                    labels: { usePointStyle: true, padding: 20 }
                }
            },
            cutout: '70%'
        }
    });
}

// Simple Tab Switcher
function switchTab(tabName) {
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.remove('active');
    });
    document.querySelectorAll('.tab').forEach(tab => {
        tab.classList.remove('active');
    });

    document.getElementById(tabName + '-tab').classList.add('active');
    event.currentTarget.classList.add('active');
}
// Check if user is admin
function checkAdminAuth() {
    const user = localStorage.getItem('user');
    if (!user) {
        window.location.href = '/auth.html';
        return null;
    }

    const userData = JSON.parse(user);
    if (userData.role !== 'admin') {
        alert('Access denied. Admin privileges required.');
        window.location.href = '/index.html';
        return null;
    }

    return userData;
}

// Format currency
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(amount || 0);
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

// Escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Tab switching
function switchTab(tabName) {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => {
        tab.classList.remove('active');
    });
    document.querySelectorAll('.tab').forEach(tab => {
        tab.classList.remove('active');
    });

    // Show selected tab
    document.getElementById(`${tabName}-tab`).classList.add('active');
    event.target.classList.add('active');

    // Load data for the tab if not loaded
    if (tabName === 'bookings' && allBookings.length === 0) {
        loadBookings();
    } else if (tabName === 'resources' && allResources.length === 0) {
        loadResources();
    } else if (tabName === 'users' && allUsers.length === 0) {
        loadUsers();
    } else if (tabName === 'categories' && allCategories.length === 0) {
        loadCategories();
    }
}

// Load Overview
async function loadOverview() {
    try {
        const [users, resources, bookings] = await Promise.all([
            fetch(`${API_BASE_URL}/users`).then(r => r.json()),
            fetch(`${API_BASE_URL}/resources`).then(r => r.json()),
            fetch(`${API_BASE_URL}/bookings`).then(r => r.json())
        ]);

        const activeBookings = bookings.filter(b => b.status === 'confirmed' || b.status === 'pending').length;

        document.getElementById('totalUsers').textContent = users.length;
        document.getElementById('totalResources').textContent = resources.length;
        document.getElementById('totalBookings').textContent = bookings.length;
        document.getElementById('activeBookings').textContent = activeBookings;

        document.getElementById('loadingOverview').style.display = 'none';
        document.getElementById('overviewContent').style.display = 'block';
    } catch (error) {
        console.error('Error loading overview:', error);
        document.getElementById('loadingOverview').textContent = 'Failed to load overview';
    }
}

// Load Bookings
async function loadBookings() {
    try {
        const response = await fetch(`${API_BASE_URL}/bookings`);
        if (!response.ok) throw new Error('Failed to load bookings');

        allBookings = await response.json();
        displayBookings(allBookings);

        document.getElementById('loadingBookings').style.display = 'none';
        document.getElementById('bookingsContent').style.display = 'block';
    } catch (error) {
        console.error('Error loading bookings:', error);
        document.getElementById('loadingBookings').textContent = 'Failed to load bookings';
    }
}

function displayBookings(bookings) {
    const tbody = document.getElementById('bookingsTableBody');
    tbody.innerHTML = '';

    if (!bookings || bookings.length === 0) {
        tbody.innerHTML = '<tr><td colspan="8" class="no-data">No bookings found</td></tr>';
        return;
    }

    bookings.forEach(booking => {
        const row = document.createElement('tr');
        const statusClass = `status-${booking.status}`;

        row.innerHTML = `
            <td>#${booking.id}</td>
            <td>User #${booking.user_id}</td>
            <td>Resource #${booking.resource_id}</td>
            <td>${formatDate(booking.start_time)}</td>
            <td>${formatDate(booking.end_time)}</td>
            <td><span class="status-badge ${statusClass}">${booking.status}</span></td>
            <td>${formatCurrency(booking.total_price)}</td>
            <td>
                <button class="action-btn edit" onclick="viewBooking(${booking.id})">View</button>
                <button class="action-btn delete" onclick="cancelBooking(${booking.id})">Cancel</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Load Resources
async function loadResources() {
    try {
        const response = await fetch(`${API_BASE_URL}/resources`);
        if (!response.ok) throw new Error('Failed to load resources');

        allResources = await response.json();
        displayResources(allResources);

        document.getElementById('loadingResources').style.display = 'none';
        document.getElementById('resourcesContent').style.display = 'block';
    } catch (error) {
        console.error('Error loading resources:', error);
        document.getElementById('loadingResources').textContent = 'Failed to load resources';
    }
}

function displayResources(resources) {
    const tbody = document.getElementById('resourcesTableBody');
    tbody.innerHTML = '';

    if (!resources || resources.length === 0) {
        tbody.innerHTML = '<tr><td colspan="9" class="no-data">No resources found</td></tr>';
        return;
    }

    resources.forEach(resource => {
        const row = document.createElement('tr');
        const statusClass = resource.is_active ? 'status-active' : 'status-inactive';
        const status = resource.is_active ? 'Active' : 'Inactive';

        row.innerHTML = `
            <td>#${resource.id}</td>
            <td>${escapeHtml(resource.name)}</td>
            <td>${escapeHtml(resource.owner_name || 'N/A')}</td>
            <td>${escapeHtml(resource.category_name || 'N/A')}</td>
            <td>${resource.capacity}</td>
            <td>${escapeHtml(resource.city || 'N/A')}</td>
            <td>${formatCurrency(resource.price_per_hour)}</td>
            <td><span class="status-badge ${statusClass}">${status}</span></td>
            <td>
                <button class="action-btn edit" onclick="editResource(${resource.id})">Edit</button>
                <button class="action-btn delete" onclick="deleteResource(${resource.id})">Delete</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Load Users
async function loadUsers() {
    try {
        const response = await fetch(`${API_BASE_URL}/users`);
        if (!response.ok) throw new Error('Failed to load users');

        allUsers = await response.json();
        displayUsers(allUsers);

        document.getElementById('loadingUsers').style.display = 'none';
        document.getElementById('usersContent').style.display = 'block';
    } catch (error) {
        console.error('Error loading users:', error);
        document.getElementById('loadingUsers').textContent = 'Failed to load users';
    }
}

function displayUsers(users) {
    const tbody = document.getElementById('usersTableBody');
    tbody.innerHTML = '';

    if (!users || users.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="no-data">No users found</td></tr>';
        return;
    }

    users.forEach(user => {
        const row = document.createElement('tr');

        row.innerHTML = `
            <td>#${user.id}</td>
            <td>${escapeHtml(user.name)}</td>
            <td>${escapeHtml(user.email)}</td>
            <td><span class="status-badge ${user.role === 'admin' ? 'status-confirmed' : 'status-pending'}">${user.role}</span></td>
            <td>${formatDate(user.created_at)}</td>
            <td>
                <button class="action-btn edit" onclick="viewUser(${user.id})">View</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Load Categories
async function loadCategories() {
    try {
        const response = await fetch(`${API_BASE_URL}/categories`);
        if (!response.ok) throw new Error('Failed to load categories');

        allCategories = await response.json();
        displayCategories(allCategories);

        document.getElementById('loadingCategories').style.display = 'none';
        document.getElementById('categoriesContent').style.display = 'block';
    } catch (error) {
        console.error('Error loading categories:', error);
        document.getElementById('loadingCategories').textContent = 'Failed to load categories';
    }
}

function displayCategories(categories) {
    const tbody = document.getElementById('categoriesTableBody');
    tbody.innerHTML = '';

    if (!categories || categories.length === 0) {
        tbody.innerHTML = '<tr><td colspan="4" class="no-data">No categories found</td></tr>';
        return;
    }

    categories.forEach(category => {
        const row = document.createElement('tr');

        row.innerHTML = `
            <td>#${category.id}</td>
            <td>${escapeHtml(category.name)}</td>
            <td>${escapeHtml(category.description || 'N/A')}</td>
            <td>
                <button class="action-btn edit" onclick="editCategory(${category.id})">Edit</button>
                <button class="action-btn delete" onclick="deleteCategory(${category.id})">Delete</button>
            </td>
        `;
        tbody.appendChild(row);
    });
}

// Search functions
function searchBookings() {
    const query = document.getElementById('bookingSearch').value.toLowerCase();
    const filtered = allBookings.filter(b =>
        b.id.toString().includes(query) ||
        b.user_id.toString().includes(query) ||
        b.resource_id.toString().includes(query)
    );
    displayBookings(filtered);
}

function searchResources() {
    const query = document.getElementById('resourceSearch').value.toLowerCase();
    const filtered = allResources.filter(r =>
        r.name.toLowerCase().includes(query) ||
        (r.city && r.city.toLowerCase().includes(query)) ||
        (r.category_name && r.category_name.toLowerCase().includes(query))
    );
    displayResources(filtered);
}

function searchUsers() {
    const query = document.getElementById('userSearch').value.toLowerCase();
    const filtered = allUsers.filter(u =>
        u.name.toLowerCase().includes(query) ||
        u.email.toLowerCase().includes(query)
    );
    displayUsers(filtered);
}

// Action handlers
function viewBooking(id) {
    alert(`View booking #${id} - Feature to be implemented`);
}

async function cancelBooking(id) {
    if (!confirm(`Are you sure you want to cancel booking #${id}?`)) return;

    try {
        const response = await fetch(`${API_BASE_URL}/bookings/${id}/cancel`, {
            method: 'POST'
        });

        if (!response.ok) throw new Error('Failed to cancel booking');

        alert('Booking cancelled successfully');
        loadBookings();
    } catch (error) {
        console.error('Error cancelling booking:', error);
        alert('Failed to cancel booking');
    }
}

function editResource(id) {
    alert(`Edit resource #${id} - Feature to be implemented`);
}

async function deleteResource(id) {
    if (!confirm(`Are you sure you want to delete resource #${id}?`)) return;

    try {
        const response = await fetch(`${API_BASE_URL}/resources/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete resource');

        alert('Resource deleted successfully');
        loadResources();
    } catch (error) {
        console.error('Error deleting resource:', error);
        alert('Failed to delete resource');
    }
}

function viewUser(id) {
    alert(`View user #${id} - Feature to be implemented`);
}

function createCategory() {
    alert('Create category - Feature to be implemented');
}

function editCategory(id) {
    alert(`Edit category #${id} - Feature to be implemented`);
}

async function deleteCategory(id) {
    if (!confirm(`Are you sure you want to delete category #${id}?`)) return;

    try {
        const response = await fetch(`${API_BASE_URL}/categories/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) throw new Error('Failed to delete category');

        alert('Category deleted successfully');
        loadCategories();
    } catch (error) {
        console.error('Error deleting category:', error);
        alert('Failed to delete category');
    }
}

function goToHome() {
    window.location.href = '/index.html';
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
    const user = checkAdminAuth();
    if (!user) return;

    document.getElementById('adminName').textContent = user.name;
    loadOverview();
});
