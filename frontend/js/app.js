const API_URL = 'http://localhost:8080/api';

function getToken() {
    return localStorage.getItem('token');
}

function setToken(token) {
    localStorage.setItem('token', token);
}

function getUser() {
    const userStr = localStorage.getItem('user');
    return userStr ? JSON.parse(userStr) : null;
}

function setUser(user) {
    localStorage.setItem('user', JSON.stringify(user));
}

function isLoggedIn() {
    return !!getToken() && !!getUser();
}

function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    window.location.href = 'index.html';
}

function formatDateTime(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleString('ru-RU');
}

async function register(name, email, password) {
    try {
        const response = await fetch(API_URL + '/auth/register', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({name, email, password})
        });

        const data = await response.json();

        if (response.ok) {
            setUser(data);
            setToken(data.id.toString());
            alert('Регистрация успешна!');
            window.location.href = 'index.html';
        } else {
            document.getElementById('register-error').textContent = data.error || 'Ошибка регистрации';
        }
    } catch (error) {
        document.getElementById('register-error').textContent = 'Ошибка соединения с сервером';
    }
}

async function login(email, password) {
    try {
        const response = await fetch(API_URL + '/auth/login', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({email, password})
        });

        const data = await response.json();

        if (response.ok) {
            setUser(data);
            setToken(data.id.toString());
            alert('Вход выполнен успешно!');
            window.location.href = 'index.html';
        } else {
            document.getElementById('login-error').textContent = data.error || 'Неверный email или пароль';
        }
    } catch (error) {
        document.getElementById('login-error').textContent = 'Ошибка соединения с сервером';
    }
}

async function loadCategories() {
    try {
        const response = await fetch(API_URL + '/resources');
        const resources = await response.json();
        
        const categoriesMap = {};
        resources.forEach(r => {
            if (r.category_name && !categoriesMap[r.category_name]) {
                categoriesMap[r.category_name] = {
                    name: r.category_name,
                    count: 1
                };
            } else if (r.category_name) {
                categoriesMap[r.category_name].count++;
            }
        });

        const container = document.getElementById('categories-grid');
        if (!container) return;

        const icons = {'Баня и сауна': '🧖', 'Бассейн': '🏊', 'Спортивная площадка': '⚽', 
                      'Конференц-зал': '🏢', 'Коворкинг': '💼', 'Студия': '🎨'};

        container.innerHTML = Object.values(categoriesMap).map(cat => `
            <div class="category-card">
                <div class="icon">${icons[cat.name] || '📍'}</div>
                <h3>${cat.name}</h3>
                <p>${cat.count} объектов</p>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading categories:', error);
    }
}

async function loadResources() {
    try {
        const response = await fetch(API_URL + '/resources');
        allResources = await response.json();
        displayResources(allResources);
    } catch (error) {
        console.error('Error loading resources:', error);
    }
}

function displayResources(resources) {
    const container = document.getElementById('resources-grid');
    if (!container) return;

    if (resources.length === 0) {
        container.innerHTML = '<p class="empty-state">Ресурсы не найдены</p>';
        return;
    }

    container.innerHTML = resources.map(resource => `
        <div class="resource-card">
            <img src="${resource.photos && resource.photos[0] ? resource.photos[0].url : 'https://via.placeholder.com/400x200'}" 
                 alt="${resource.name}">
            <h3>${resource.name}</h3>
            <p>${resource.description || ''}</p>
            <div class="meta">
                <span>👥 ${resource.capacity} чел.</span>
                <span>📍 ${resource.city || 'Алматы'}</span>
            </div>
            ${resource.price_per_hour ? `<div class="price">${resource.price_per_hour} ₸/час</div>` : ''}
            ${resource.amenities && resource.amenities.length > 0 ? `
                <div style="margin: 10px 0; font-size: 0.9rem; color: #666;">
                    ${resource.amenities.slice(0, 3).join(', ')}
                </div>
            ` : ''}
            <button class="btn btn-primary btn-block" onclick="openBookingModal(${resource.id})">
                Забронировать
            </button>
        </div>
    `).join('');
}

async function createBooking() {
    if (!currentResource) return;

    const user = getUser();
    if (!user) {
        alert('Необходимо войти в систему');
        return;
    }

    const startTime = document.getElementById('start-time').value;
    const endTime = document.getElementById('end-time').value;
    const notes = document.getElementById('notes').value;

    try {
        const response = await fetch(API_URL + '/bookings', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + getToken()
            },
            body: JSON.stringify({
                user_id: user.id,
                resource_id: currentResource.id,
                start_time: new Date(startTime).toISOString(),
                end_time: new Date(endTime).toISOString(),
                notes: notes
            })
        });

        if (response.ok) {
            alert('Бронирование успешно создано!');
            closeBookingModal();
            window.location.href = 'bookings.html';
        } else {
            const data = await response.json();
            document.getElementById('booking-error').textContent = data.error || 'Ошибка создания бронирования';
        }
    } catch (error) {
        document.getElementById('booking-error').textContent = 'Ошибка соединения с сервером';
    }
}

async function loadUserBookings() {
    const user = getUser();
    if (!user) return;

    try {
        const response = await fetch(API_URL + '/users/' + user.id + '/bookings');
        allBookings = await response.json();
        
        const container = document.getElementById('bookings-list');
        if (!container) return;

        if (allBookings.length === 0) {
            container.innerHTML = '<p class="empty-state">У вас пока нет бронирований</p>';
            return;
        }

        container.innerHTML = allBookings.map(booking => `
            <div class="booking-card">
                <div class="booking-header">
                    <h3>Ресурс #${booking.resource_id}</h3>
                    <span class="status-badge ${booking.status}">${getStatusText(booking.status)}</span>
                </div>
                <div class="booking-details">
                    <p><strong>Начало:</strong> ${formatDateTime(booking.start_time)}</p>
                    <p><strong>Окончание:</strong> ${formatDateTime(booking.end_time)}</p>
                    ${booking.notes ? `<p><strong>Комментарий:</strong> ${booking.notes}</p>` : ''}
                </div>
                ${booking.status !== 'cancelled' ? `
                    <button class="btn btn-danger" onclick="cancelBooking(${booking.id})">
                        Отменить бронь
                    </button>
                ` : ''}
            </div>
        `).join('');
    } catch (error) {
        console.error('Error loading bookings:', error);
    }
}

function getStatusText(status) {
    const statuses = {
        'confirmed': 'Подтверждено',
        'pending': 'В ожидании',
        'cancelled': 'Отменено'
    };
    return statuses[status] || status;
}

async function cancelBooking(bookingId) {
    if (!confirm('Вы уверены, что хотите отменить это бронирование?')) {
        return;
    }

    try {
        const response = await fetch(API_URL + '/bookings/' + bookingId + '/cancel', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + getToken()
            }
        });

        if (response.ok) {
            alert('Бронирование успешно отменено');
            loadUserBookings();
        } else {
            alert('Ошибка при отмене бронирования');
        }
    } catch (error) {
        alert('Ошибка соединения с сервером');
    }
}

document.addEventListener('DOMContentLoaded', function() {
    const user = getUser();
    const authLink = document.getElementById('auth-link');
    const profileLink = document.getElementById('profile-link');
    const logoutLink = document.getElementById('logout-link');
    const bookingsLink = document.getElementById('bookings-link');

    if (user && authLink && profileLink && logoutLink) {
        authLink.style.display = 'none';
        profileLink.style.display = 'block';
        logoutLink.style.display = 'block';
        if (bookingsLink) bookingsLink.style.display = 'block';
        
        const userName = document.getElementById('user-name');
        if (userName) userName.textContent = user.name;
        
        logoutLink.addEventListener('click', function(e) {
            e.preventDefault();
            logout();
        });
    }
});
