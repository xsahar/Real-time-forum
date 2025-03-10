// src/static/js/chat.js
import { loadMessages, sendMessage, initMessages, setActiveRecipient } from './messages.js';

export const initChat = () => {
    const chatSection = document.getElementById('chat-section');
    const chatToggle = document.getElementById('chat-toggle');
    const usersList = chatSection.querySelector('.users-list');
    const chatMessages = chatSection.querySelector('.chat-messages');
    const messagesContainer = chatSection.querySelector('.messages-container');
    const messageInput = chatSection.querySelector('.message-input input');
    const sendButton = chatSection.querySelector('.send-message');
    const backButton = chatSection.querySelector('.back-to-users');

    let currentChatUser = null;

    if (!chatSection || !chatToggle) {
        console.error('Chat elements not found:', { chatSection, chatToggle });
        return;
    }

    // Initialize messages module
    initMessages();

    // Toggle chat sidebar
    chatToggle.addEventListener('click', (e) => {
        e.preventDefault();
        e.stopPropagation();
        console.log('Chat toggle clicked');
        chatSection.classList.toggle('open');

        if (chatSection.classList.contains('open')) {
            loadOnlineUsers();
        }
    });

    // Close chat when clicking outside
    document.addEventListener('click', (e) => {
        if (!chatSection.contains(e.target) &&
            !chatToggle.contains(e.target) &&
            chatSection.classList.contains('open')) {
            chatSection.classList.remove('open');
        }
    });

    // Handle user selection
    usersList.addEventListener('click', (e) => {
        const userItem = e.target.closest('.user-item');
        if (userItem) {
            const userId = userItem.dataset.userId;
            currentChatUser = userId;
            setActiveRecipient(userId);
            showChatMessages(userId);
        }
    });

    // Handle message sending
    sendButton.addEventListener('click', () => sendChatMessage());
    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            sendChatMessage();
        }
    });

    // Add back button handler
    backButton.addEventListener('click', () => {
        chatMessages.style.display = 'none';
        usersList.style.display = 'block';
        backButton.style.display = 'none';
        currentChatUser = null;
        setActiveRecipient(null);
    });

    // Show chat toggle for logged-in users
    const token = localStorage.getItem('token');
    if (token) {
        console.log('User is logged in, showing chat toggle');
        chatToggle.classList.add('visible');
        loadOnlineUsers();
    } else {
        console.log('User is not logged in, hiding chat toggle');
        chatToggle.classList.remove('visible');
    }

    async function loadOnlineUsers() {
        try {
            const response = await fetch('/api/users/online', {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`
                }
            });

            if (!response.ok) throw new Error('Failed to load users');

            const users = await response.json();
            sortAndRenderUsers(users);
        } catch (error) {
            console.error('Error loading users:', error);
            usersList.innerHTML = '<div class="error">Failed to load users</div>';
        }
    }

    function sortAndRenderUsers(users) {
        // Sort users: first by having recent messages, then by online status, then by nickname
        users.sort((a, b) => {
            // First priority: has messages (this would need additional logic from your backend)
            if (a.hasMessages && !b.hasMessages) return -1;
            if (!a.hasMessages && b.hasMessages) return 1;

            // Second priority: online status
            if (a.online && !b.online) return -1;
            if (!a.online && b.online) return 1;

            // Last priority: nickname alphabetically
            return a.nickname.localeCompare(b.nickname);
        });

        usersList.innerHTML = users.map(user => `
            <div class="user-item" data-user-id="${user.id}">
                <span class="user-status ${user.online ? 'online' : 'offline'}"></span>
                <span class="user-name">${user.nickname}</span>
            </div>
        `).join('');
    }

    async function showChatMessages(userId) {
        usersList.style.display = 'none';
        chatMessages.style.display = 'flex';
        backButton.style.display = 'block';

        // Clear previous messages
        messagesContainer.innerHTML = '<div class="loading-indicator">Loading messages...</div>';

        await loadChatHistory(userId);
    }

    async function loadChatHistory(userId) {
        try {
            const messages = await loadMessages(userId);
            renderMessages(messages);
        } catch (error) {
            console.error('Error loading messages:', error);
            messagesContainer.innerHTML = '<div class="error">Failed to load messages</div>';
        }
    }

    function renderMessages(messages) {
        if (!messages || messages.length === 0) {
            messagesContainer.innerHTML = '<div class="no-messages">No messages yet. Say hello!</div>';
            return;
        }
    
        const currentUserId = JSON.parse(localStorage.getItem('user')).id;
        messagesContainer.innerHTML = '';
    
        // Make sure messages are sorted by date (oldest first)
        const sortedMessages = [...messages].sort((a, b) => {
            return new Date(a.createdAt) - new Date(b.createdAt);
        });
    
        // Render each message
        sortedMessages.forEach(message => {
            const isIncoming = message.senderId !== currentUserId;
            const messageElement = document.createElement('div');
            messageElement.className = `message ${isIncoming ? 'received' : 'sent'}`;
            messageElement.dataset.messageId = message.id;
    
            messageElement.innerHTML = `
                <div class="message-content">${message.content}</div>
                <div class="message-time">${new Date(message.createdAt).toLocaleTimeString()}</div>
            `;
            messagesContainer.appendChild(messageElement);
        });
    
        // Scroll to bottom
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    function sendChatMessage() {
        if (!currentChatUser || !messageInput.value.trim()) return;

        const content = messageInput.value.trim();

        if (sendMessage(currentChatUser, content)) {
            // Add message to UI immediately for better UX
            const messageElement = document.createElement('div');
            messageElement.className = 'message sent';

            const tempId = 'temp-' + Date.now();
            messageElement.dataset.messageId = tempId;

            messageElement.innerHTML = `
                <div class="message-content">${content}</div>
                <div class="message-time">${new Date().toLocaleTimeString()}</div>
            `;

            const noMessages = messagesContainer.querySelector('.no-messages');
            if (noMessages) {
                noMessages.remove();
            }

            messagesContainer.appendChild(messageElement);

            // Clear input and scroll to bottom
            messageInput.value = '';
            messagesContainer.scrollTop = messagesContainer.scrollHeight;

            // Don't reload all messages after sending - WebSocket will handle confirmation
            // Remove this line: loadChatHistory(currentChatUser);
        }
    }
};