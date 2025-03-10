// src/static/js/messages.js
import { WebSocketClient } from './websocket.js';

let wsClient = null;
let activeRecipient = null;
let messageCache = new Map();

export const initMessages = () => {
    wsClient = new WebSocketClient();
    wsClient.connect();

    wsClient.onMessage(handleIncomingMessage);
    wsClient.onStatus(handleStatusUpdate);

    // Setup debouncing for scroll event
    setupScrollHandling();

    console.log('Messages module initialized');
};

export const loadMessages = async (userId, limit = 10, offset = 0) => {
    try {
        const token = localStorage.getItem('token');
        if (!token) throw new Error('Authentication required');

        // Check cache first
        const cacheKey = `${userId}_${offset}_${limit}`;
        if (messageCache.has(cacheKey)) {
            return messageCache.get(cacheKey);
        }

        console.log(`Loading messages with ${userId}, limit: ${limit}, offset: ${offset}`);

        const response = await fetch(`/api/messages/${userId}?limit=${limit}&offset=${offset}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            const errorText = await response.text();
            console.error('Message loading error:', response.status, errorText);
            throw new Error('Failed to load messages');
        }

        const messages = await response.json();
        console.log(`Loaded ${messages.length} messages with ${userId}`);

        // Cache the results
        messageCache.set(cacheKey, messages);
        return messages;
    } catch (error) {
        console.error('Error loading messages:', error);
        throw error;
    }
};

export const sendMessage = (receiverId, content) => {
    console.log(`Attempting to send message to ${receiverId}`);
    if (!receiverId || !content || content.trim() === '') {
        console.error('Invalid message parameters');
        return false;
    }

    // Try sending via WebSocket first
    if (wsClient && wsClient.sendMessage(receiverId, content)) {
        console.log('Message sent via WebSocket');
        return true;
    }

    console.warn('WebSocket unavailable, message not sent');
    return false;
};

const handleIncomingMessage = (payload) => {
    console.log('Received message:', payload);
    
    const currentUser = JSON.parse(localStorage.getItem('user'));
    if (!currentUser) return;
    
    const isSender = payload.senderId === currentUser.id;
    
    // Update UI if this is a conversation we're currently viewing
    const messagesContainer = document.querySelector('.messages-container');
    
    if (messagesContainer) {
        const isActiveChat = (isSender && payload.receiverId === activeRecipient) || 
            (!isSender && payload.senderId === activeRecipient);
        
        if (isActiveChat) {
            // Check for duplicate messages
            const existingMessage = Array.from(messagesContainer.querySelectorAll('.message')).find(el => {
                return el.dataset.messageId === payload.id;
            });
            
            if (existingMessage) {
                console.log('Message already displayed:', payload.id);
                return;
            }
            
            // Remove "no messages" placeholder if it exists
            const noMessagesEl = messagesContainer.querySelector('.no-messages');
            if (noMessagesEl) {
                noMessagesEl.remove();
            }
            
            // If this is a message from us, remove any temporary message we might have added
            if (isSender) {
                const tempMessages = messagesContainer.querySelectorAll('.message.sent[data-message-id^="temp-"]');
                tempMessages.forEach(tempMsg => {
                    // Check if content matches
                    const content = tempMsg.querySelector('.message-content')?.textContent;
                    if (content === payload.content) {
                        tempMsg.remove();
                    }
                });
            }
            
            // Create message element with correct styles
            const messageElement = document.createElement('div');
            messageElement.className = `message ${isSender ? 'sent' : 'received'}`;
            messageElement.dataset.messageId = payload.id;
            
            messageElement.innerHTML = `
                <div class="message-content">${payload.content}</div>
                <div class="message-time">${new Date(payload.createdAt).toLocaleTimeString()}</div>
            `;
            
            messagesContainer.appendChild(messageElement);
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        } else if (!isSender) {
            // Show notification for incoming message not in active chat
            const usersList = document.querySelector('.users-list');
            const userItem = usersList?.querySelector(`[data-user-id="${payload.senderId}"]`);
            if (userItem && !userItem.querySelector('.notification-dot')) {
                const notificationDot = document.createElement('span');
                notificationDot.className = 'notification-dot';
                userItem.appendChild(notificationDot);
            }
        }
    }
    
    // Always invalidate cache for this conversation to ensure fresh data on next load
    const userId = isSender ? payload.receiverId : payload.senderId;
    
    // Remove any cached messages for this conversation
    for (const key of messageCache.keys()) {
        if (key.includes(userId)) {
            messageCache.delete(key);
        }
    }
};

const handleStatusUpdate = (payload) => {
    console.log('User status update:', payload);

    // Update the online status in the user list
    const usersList = document.querySelector('.users-list');
    if (usersList) {
        const userItem = usersList.querySelector(`[data-user-id="${payload.userId}"]`);
        if (userItem) {
            const statusIndicator = userItem.querySelector('.status-indicator');
            if (statusIndicator) {
                statusIndicator.classList.toggle('online', payload.online);
                statusIndicator.classList.toggle('offline', !payload.online);
            }
        }
    }
};

// Helper function to debounce scroll events
const debounce = (func, wait) => {
    let timeout;
    return function (...args) {
        clearTimeout(timeout);
        timeout = setTimeout(() => func.apply(this, args), wait);
    };
};

// Handle scrolling to load more messages
const setupScrollHandling = () => {
    document.addEventListener('click', event => {
        const messagesContainer = document.querySelector('.messages-container');
        if (!messagesContainer) return;

        // Setup scroll event listener once container is available
        messagesContainer.addEventListener('scroll', debounce(async () => {
            if (messagesContainer.scrollTop <= 50 && activeRecipient) {
                const currentMessages = messagesContainer.querySelectorAll('.message').length;
                try {
                    const newMessages = await loadMessages(activeRecipient, 10, currentMessages);
                    if (newMessages.length > 0) {
                        // Save current scroll position and height
                        const scrollPos = messagesContainer.scrollHeight - messagesContainer.scrollTop;

                        // Add older messages to the top
                        newMessages.reverse().forEach(msg => {
                            prependMessage(messagesContainer, msg);
                        });

                        // Restore scroll position
                        messagesContainer.scrollTop = messagesContainer.scrollHeight - scrollPos;
                    }
                } catch (error) {
                    console.error('Error loading more messages:', error);
                }
            }
        }, 300));
    }, { once: true });
};

// Helper to append a message to the container
const appendMessage = (container, message) => {
    const currentUser = JSON.parse(localStorage.getItem('user'));
    if (!currentUser) return;

    const isSender = message.senderId === currentUser.id;

    const messageElement = document.createElement('div');
    messageElement.classList.add('message', isSender ? 'outgoing' : 'incoming');

    messageElement.innerHTML = `
        <div class="message-content">${message.content}</div>
        <div class="message-meta">
            <span class="message-time">${new Date(message.createdAt).toLocaleTimeString()}</span>
        </div>
    `;

    container.appendChild(messageElement);
};

// Helper to prepend a message to the container
const prependMessage = (container, message) => {
    const currentUser = JSON.parse(localStorage.getItem('user'));
    if (!currentUser) return;

    const isSender = message.senderId === currentUser.id;

    const messageElement = document.createElement('div');
    messageElement.classList.add('message', isSender ? 'outgoing' : 'incoming');

    messageElement.innerHTML = `
        <div class="message-content">${message.content}</div>
        <div class="message-meta">
            <span class="message-time">${new Date(message.createdAt).toLocaleTimeString()}</span>
        </div>
    `;

    container.insertBefore(messageElement, container.firstChild);
};

// Helper to update the user list order
const updateUserListOrder = (userId) => {
    const usersList = document.querySelector('.users-list');
    if (!usersList) return;

    const userItem = usersList.querySelector(`[data-user-id="${userId}"]`);
    if (userItem && usersList.firstChild !== userItem) {
        usersList.insertBefore(userItem, usersList.firstChild);
    }
};

// Helper to show message notification
const showMessageNotification = (senderId) => {
    const usersList = document.querySelector('.users-list');
    if (!usersList) return;

    const userItem = usersList.querySelector(`[data-user-id="${senderId}"]`);
    if (userItem) {
        // Add notification indicator
        if (!userItem.querySelector('.notification-badge')) {
            const badge = document.createElement('span');
            badge.classList.add('notification-badge');
            badge.textContent = '!';
            userItem.appendChild(badge);
        }
    }
};

// Clear cache entries for a specific user
const invalidateCache = (userId) => {
    for (const key of messageCache.keys()) {
        if (key.startsWith(`${userId}_`)) {
            messageCache.delete(key);
        }
    }
};

export const setActiveRecipient = (userId) => {
    activeRecipient = userId;
    console.log('Set active recipient:', userId);
};