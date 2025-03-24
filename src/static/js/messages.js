import { WebSocketClient } from './websocket.js';
import { refreshOnlineUsers } from './chat.js';

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

// Further fix the chat notification logic in handleIncomingMessage

const handleIncomingMessage = (payload) => {
    console.log('Received message:', payload);
    
    const currentUser = JSON.parse(localStorage.getItem('user'));
    if (!currentUser) return;
    
    const isSender = payload.senderId === currentUser.id;
    const conversationPartnerId = isSender ? payload.receiverId : payload.senderId;
    
    // Get scroll position BEFORE any DOM changes
    const messagesContainer = document.querySelector('.messages-container');
    const wasAtBottom = messagesContainer ? 
        (messagesContainer.scrollHeight - messagesContainer.scrollTop <= messagesContainer.clientHeight + 20) : 
        false;
    
    // Check if this is the active conversation directly from activeRecipient
    // More reliable than DOM checking
    const isActiveConversation = activeRecipient === conversationPartnerId;
    
    // Check if chat UI is visible
    const chatMessages = document.querySelector('.chat-messages');
    const isChatVisible = chatMessages && 
        window.getComputedStyle(chatMessages).display === 'flex';
    
    // A chat is truly OPEN when it's the active conversation AND the chat UI is visible
    const isChatOpen = isChatVisible && isActiveConversation;
    
    // First, update the UI if this is the active chat
    if (isActiveConversation && messagesContainer) {
        // Update message display immediately
        updateMessageDisplay(messagesContainer, payload, isSender);
        
        // Handle scrolling
        if (wasAtBottom) {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }
    }
    
    // Handle notification logic
    if (isSender || isChatOpen) {
        // If we're the sender OR the chat is open, clear any notification
        const unreadMessages = JSON.parse(localStorage.getItem('unreadMessages') || '{}');
        if (unreadMessages[conversationPartnerId]) {
            delete unreadMessages[conversationPartnerId];
            localStorage.setItem('unreadMessages', JSON.stringify(unreadMessages));
            
            // Update UI immediately
            refreshOnlineUsers();
        }
    } else if (!isSender) {
        // We're not the sender and the chat isn't open - mark as unread
        const unreadMessages = JSON.parse(localStorage.getItem('unreadMessages') || '{}');
        unreadMessages[payload.senderId] = true;
        localStorage.setItem('unreadMessages', JSON.stringify(unreadMessages));
        
        // Update UI immediately
        refreshOnlineUsers();
    }
    
    // Clear message cache for this conversation
    for (const key of messageCache.keys()) {
        if (key.includes(conversationPartnerId)) {
            messageCache.delete(key);
        }
    }
};

// Extract message display logic to a separate function
function updateMessageDisplay(container, message, isSender) {
    // Check for duplicate messages
    const existingMessage = Array.from(container.querySelectorAll('.message')).find(el => {
        return el.dataset.messageId === message.id;
    });
    
    if (existingMessage) {
        console.log('Message already displayed:', message.id);
        return;
    }
    
    // Remove "no messages" placeholder if it exists
    const noMessagesEl = container.querySelector('.no-messages');
    if (noMessagesEl) {
        noMessagesEl.remove();
    }
    
    // If this is a message from us, remove any temporary message we might have added
    if (isSender) {
        const tempMessages = container.querySelectorAll('.message.sent[data-message-id^="temp-"]');
        tempMessages.forEach(tempMsg => {
            // Check if content matches
            const content = tempMsg.querySelector('.message-content')?.textContent;
            if (content === message.content) {
                tempMsg.remove();
            }
        });
    }
    
    // Create message element with correct styles
    const messageElement = document.createElement('div');
    messageElement.className = `message ${isSender ? 'sent' : 'received'}`;
    messageElement.dataset.messageId = message.id;
    
    messageElement.innerHTML = `
        <div class="message-content">${message.content}</div>
        <div class="message-time">${new Date(message.createdAt).toLocaleTimeString()}</div>
    `;
    
    container.appendChild(messageElement);
    container.scrollTop = container.scrollHeight;
}

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

        // Force scrollbar to always show
        messagesContainer.style.overflowY = 'scroll';
        
        // Setup scroll event listener once container is available
        messagesContainer.addEventListener('scroll', debounce(async () => {
            if (messagesContainer.scrollTop <= 50 && activeRecipient) {
                // Prevent multiple simultaneous loading operations
                if (messagesContainer.dataset.isLoading === 'true') {
                    return;
                }
                messagesContainer.dataset.isLoading = 'true';
                
                const currentMessages = messagesContainer.querySelectorAll('.message').length;
                
                // Check if we actually need to load more messages
                // If we have very few messages, don't try to load more
                if (currentMessages < 5) {
                    console.log('Not enough messages to warrant loading more');
                    messagesContainer.dataset.isLoading = 'false';
                    return;
                }
                
                // Get existing message IDs
                const existingMessageIds = new Set();
                messagesContainer.querySelectorAll('.message').forEach(el => {
                    if (el.dataset.messageId) {
                        existingMessageIds.add(el.dataset.messageId);
                    }
                });
                
                try {
                    const newMessages = await loadMessages(activeRecipient, 10, currentMessages);
                    
                    // Filter out any messages that already exist in the DOM
                    const uniqueMessages = newMessages.filter(msg => 
                        msg.id && !existingMessageIds.has(msg.id)
                    );
                    
                    if (uniqueMessages.length > 0) {
                        // Save current scroll position and height
                        const scrollPos = messagesContainer.scrollHeight - messagesContainer.scrollTop;

                        // Add older messages to the top
                        uniqueMessages.reverse().forEach(msg => {
                            prependMessage(messagesContainer, msg);
                        });

                        // Restore scroll position
                        messagesContainer.scrollTop = messagesContainer.scrollHeight - scrollPos;
                    }
                } catch (error) {
                    console.error('Error loading more messages:', error);
                } finally {
                    // Always clear the loading flag
                    messagesContainer.dataset.isLoading = 'false';
                }
            }
        }, 300));
    }, { once: true });
};

// Helper to prepend a message to the container
const prependMessage = (container, message) => {
    const currentUser = JSON.parse(localStorage.getItem('user'));
    if (!currentUser) return;
    
    // Skip if this message ID already exists in the container
    if (message.id && container.querySelector(`.message[data-message-id="${message.id}"]`)) {
        return;
    }

    const isSender = message.senderId === currentUser.id;

    const messageElement = document.createElement('div');
    messageElement.className = `message ${isSender ? 'sent' : 'received'}`;
    messageElement.dataset.messageId = message.id;

    messageElement.innerHTML = `
        <div class="message-content">${message.content}</div>
        <div class="message-time">${new Date(message.createdAt).toLocaleTimeString()}</div>
    `;

    container.insertBefore(messageElement, container.firstChild);
};

export const clearUnreadStatus = (userId) => {
    if (!userId) return;
    
    console.log(`Clearing unread status for ${userId}`);
    
    // Clear from localStorage
    const unreadMessages = JSON.parse(localStorage.getItem('unreadMessages') || '{}');
    if (unreadMessages[userId]) {
        delete unreadMessages[userId];
        localStorage.setItem('unreadMessages', JSON.stringify(unreadMessages));
        
        // Also try to clear the notification dot visually
        const usersList = document.querySelector('.users-list');
        if (usersList) {
            const userItem = usersList.querySelector(`[data-user-id="${userId}"]`);
            if (userItem) {
                const notificationDot = userItem.querySelector('.notification-dot');
                if (notificationDot) {
                    notificationDot.remove();
                }
            }
        }
    }
};

// Then modify the setActiveRecipient function:
export const setActiveRecipient = (userId) => {
    // Update the activeRecipient value
    activeRecipient = userId;
    console.log('Set active recipient:', userId);
    
    // Clear unread status when setting active recipient
    if (userId) {
        clearUnreadStatus(userId);
        
        // Also update the UI by removing the notification dot
        const userItem = document.querySelector(`.user-item[data-user-id="${userId}"]`);
        if (userItem) {
            // Mark this user as active
            document.querySelectorAll('.user-item.active').forEach(item => {
                item.classList.remove('active');
            });
            userItem.classList.add('active');
            
            // Remove any notification dot
            const notificationDot = userItem.querySelector('.notification-dot');
            if (notificationDot) {
                notificationDot.remove();
            }
        }
    }
};