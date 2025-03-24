import { formatMessageDate } from './chatUtils.js';

export async function loadOnlineUsers() {
    const usersList = document.querySelector('.users-list');
    const token = localStorage.getItem('token');
    if (!token) {
        console.error('No authentication token found');
        if (usersList) {
            usersList.innerHTML = '<div class="message-info">Please log in to see users</div>';
        }
        return;
    }
    try {
        const response = await fetch('/api/users/online', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) {
            throw new Error(`Failed to load users: ${response.status}`);
        }
        const users = await response.json();
        if (!users || users.length === 0) {
            if (usersList) {
                usersList.innerHTML = '<div class="message-info">No other users found</div>';
            }
            return;
        }
        sortAndRenderUsers(users);
    } catch (error) {
        console.error('Error loading users:', error);
        if (usersList) {
            usersList.innerHTML = '<div class="message-info">Failed to load users</div>';
        }
    }
}

function sortAndRenderUsers(users) {
    const usersList = document.querySelector('.users-list');
    let unreadMessages = JSON.parse(localStorage.getItem('unreadMessages') || '{}');

    users.sort((a, b) => {
        if (unreadMessages[a.id] && !unreadMessages[b.id]) return -1;
        if (!unreadMessages[a.id] && unreadMessages[b.id]) return 1;
        if (a.hasMessages && b.hasMessages) {
            return new Date(b.lastMessageTime) - new Date(a.lastMessageTime);
        }
        if (a.hasMessages && !b.hasMessages) return -1;
        if (!a.hasMessages && b.hasMessages) return 1;
        if (a.online && !b.online) return -1;
        if (!a.online && b.online) return 1;
        return a.nickname.localeCompare(b.nickname);
    });

    if (usersList) {
        usersList.innerHTML = users.map(user => `
            <div class="user-item" data-user-id="${user.id}">
                <span class="user-status ${user.online ? 'online' : 'offline'}"></span>
                <span class="user-name">${user.nickname}</span>
                ${user.hasMessages ? `<span class="last-message-time">${formatMessageDate(new Date(user.lastMessageTime))}</span>` : ''}
                ${unreadMessages[user.id] ? '<span class="notification-dot"></span>' : ''}
            </div>
        `).join('');
    }
}

export async function loadChatHistory(userId) {
    const token = localStorage.getItem('token');
    try {
        const response = await fetch(`/api/messages/${userId}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) {
            throw new Error(`Failed to load chat history: ${response.status}`);
        }
        const messages = await response.json();
        return messages;
    } catch (error) {
        throw error;
    }
}
