import { loadOnlineUsers} from './chatAPI.js';
import { formatMessageDate } from './chatUtils.js';
import { loadMessages, sendMessage, initMessages, setActiveRecipient } from './messages.js';

let usersList; // DOM element reference
let currentChatUser = null;
let unreadMessages = {}; // Unread message status

export const initChat = () => {
    const chatSection = document.getElementById('chat-section');
    const chatToggle = document.getElementById('chat-toggle');
    usersList = chatSection.querySelector('.users-list');
    const chatMessages = chatSection.querySelector('.chat-messages');
    const messagesContainer = chatSection.querySelector('.messages-container');
    const messageInput = chatSection.querySelector('.message-input input');
    const sendButton = chatSection.querySelector('.send-message');
    const backButton = chatSection.querySelector('.back-to-users');

    // Load unread messages from localStorage
    unreadMessages = JSON.parse(localStorage.getItem('unreadMessages') || '{}');

    if (!chatSection || !chatToggle) {
        console.error('Chat elements not found:', { chatSection, chatToggle });
        return;
    }

    // Initialize messages module
    initMessages();

    const isMobile = window.innerWidth <= 768;
    if (isMobile) {
        chatToggle.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
            console.log('Chat toggle clicked (mobile)');
            chatSection.classList.toggle('open');
            if (chatSection.classList.contains('open')) {
                loadOnlineUsers();
            }
        });
        document.addEventListener('click', (e) => {
            if (!chatSection.contains(e.target) && !chatToggle.contains(e.target) && chatSection.classList.contains('open')) {
                chatSection.classList.remove('open');
            }
        });
    } else {
        chatSection.classList.add('open');
        chatToggle.style.display = 'none';
    }

    // Handle user selection from the users list
    usersList.addEventListener('click', (e) => {
        const userItem = e.target.closest('.user-item');
        if (userItem) {
            const userId = userItem.dataset.userId;
            currentChatUser = userId;
            setActiveRecipient(userId);
            showChatMessages(userId);
            const unread = JSON.parse(localStorage.getItem('unreadMessages') || '{}');
            if (unread[userId]) {
                delete unread[userId];
                localStorage.setItem('unreadMessages', JSON.stringify(unread));
            }
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

    // Back button handler
    backButton.addEventListener('click', () => {
        chatMessages.style.display = 'none';
        usersList.style.display = 'block';
        backButton.style.display = 'none';
        
        // Clear the active recipient and remove active class from all user items
        currentChatUser = null;
        setActiveRecipient(null);
        
        // IMPORTANT: Remove the currentChatOpen flag when closing chat
        localStorage.removeItem('currentChatOpen');
        
        document.querySelectorAll('.user-item.active').forEach(item => {
            item.classList.remove('active');
        });
        
        // Refresh users list to ensure correct status display
        loadOnlineUsers();
    });

    // Show chat if user is logged in
    const token = localStorage.getItem('token');
    if (token) {
        console.log('User is logged in, showing chat');
        if (isMobile) {
            chatToggle.classList.add('visible');
        } else {
            loadOnlineUsers();
        }
    } else {
        console.log('User is not logged in, hiding chat');
        chatSection.classList.remove('open');
        chatToggle.classList.remove('visible');
    }

    // Reload page on viewport change (mobile/desktop)
    window.addEventListener('resize', () => {
        const newIsMobile = window.innerWidth <= 768;
        if (newIsMobile !== isMobile) {
            location.reload();
        }
    });

    // Show messages for a selected user
    async function showChatMessages(userId) {
        usersList.style.display = 'none';
        chatMessages.style.display = 'flex';
        backButton.style.display = 'block';
        messagesContainer.innerHTML = '<div class="loading-indicator"><div class="spinner"></div>Loading messages...</div>';
        
        // Set active recipient
        setActiveRecipient(userId);
        
        // Mark the chat as viewed in localStorage
        localStorage.setItem('currentChatOpen', userId);
        
        // Remove any notification dot
        const unread = JSON.parse(localStorage.getItem('unreadMessages') || '{}');
        if (unread[userId]) {
            delete unread[userId];
            localStorage.setItem('unreadMessages', JSON.stringify(unread));
        }

        try {
            // Use chatAPI to load messages - limit to 10 most recent
            const messages = await loadMessages(userId, 10, 0);
            
            // Simply render the messages - don't set up scrolling here
            renderMessages(messages);
            
            // Let the global scroll handler in messages.js take care of loading more
        } catch (error) {
            console.error('Error loading messages:', error);
            messagesContainer.innerHTML = '<div class="error">Failed to load messages</div>';
        }
    }

    // Update renderMessages to include date headers like WhatsApp
    function renderMessages(messages) {
        if (!messages || messages.length === 0) {
            messagesContainer.innerHTML = '<div class="no-messages">No messages yet. Say hello!</div>';
            return;
        }
        
        const currentUserId = JSON.parse(localStorage.getItem('user')).id;
        messagesContainer.innerHTML = '';
        
        // Sort messages by date (oldest first)
        const sortedMessages = [...messages].sort((a, b) => new Date(a.createdAt) - new Date(b.createdAt));
        
        let currentDate = null;
        
        // Add message elements to DOM
        sortedMessages.forEach(message => {
            const messageDate = new Date(message.createdAt);
            const messageDateStr = messageDate.toDateString();
            
            // If this message is from a different date than the previous one, add a date header
            if (currentDate !== messageDateStr) {
                currentDate = messageDateStr;
                
                // Create a date divider
                const dateHeader = document.createElement('div');
                dateHeader.className = 'date-header';
                
                // Get a user-friendly date string
                const today = new Date().toDateString();
                const yesterday = new Date();
                yesterday.setDate(yesterday.getDate() - 1);
                const yesterdayStr = yesterday.toDateString();
                
                let dateText;
                if (messageDateStr === today) {
                    dateText = 'Today';
                } else if (messageDateStr === yesterdayStr) {
                    dateText = 'Yesterday';
                } else {
                    // For other dates, use a more readable format
                    dateText = messageDate.toLocaleDateString([], {
                        weekday: 'long',
                        month: 'long', 
                        day: 'numeric', 
                        year: messageDate.getFullYear() !== new Date().getFullYear() ? 'numeric' : undefined
                    });
                }
                
                dateHeader.textContent = dateText;
                messagesContainer.appendChild(dateHeader);
            }
            
            // Create the message element
            const isIncoming = message.senderId !== currentUserId;
            const messageElement = document.createElement('div');
            messageElement.className = `message ${isIncoming ? 'received' : 'sent'}`;
            messageElement.dataset.messageId = message.id;
            
            // Show only time for the message, not the full date
            const timeStr = messageDate.toLocaleTimeString([], {hour: '2-digit', minute:'2-digit'});
            
            messageElement.innerHTML = `
                <div class="message-content">${message.content}</div>
                <div class="message-time">${timeStr}</div>
            `;
            messagesContainer.appendChild(messageElement);
        });
        
        // Force scroll to bottom after messages render
        setTimeout(() => {
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        }, 50);
    }

    function sendChatMessage() {
        if (!currentChatUser || !messageInput.value.trim()) return;
        
        // Add a character limit to prevent excessively long messages
        const content = messageInput.value.trim();
        
        // Limit message length
        const MAX_MESSAGE_LENGTH = 250;
        const truncatedContent = content.length > MAX_MESSAGE_LENGTH ? 
            content.substring(0, MAX_MESSAGE_LENGTH) + "..." : 
            content;
        
        if (sendMessage(currentChatUser, truncatedContent)) {
            const messageElement = document.createElement('div');
            messageElement.className = 'message sent';
            const tempId = 'temp-' + Date.now();
            messageElement.dataset.messageId = tempId;
            messageElement.innerHTML = `
                <div class="message-content">${truncatedContent}</div>
                <div class="message-time">${formatMessageDate(new Date())}</div>
            `;
            const noMessages = messagesContainer.querySelector('.no-messages');
            if (noMessages) noMessages.remove();
            messagesContainer.appendChild(messageElement);
            messageInput.value = '';
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
            refreshOnlineUsers();
        }
    }

    // Expose functions to mark messages as unread and to reload users
    window.markMessageUnread = (senderId) => {
        if (currentChatUser === senderId) return false;
        unreadMessages[senderId] = true;
        localStorage.setItem('unreadMessages', JSON.stringify(unreadMessages));
        loadOnlineUsers();
        return true;
    };

    window.loadOnlineUsers = loadOnlineUsers;
};

export const handleNewMessage = (senderId) => {
    if (typeof window.markMessageUnread === 'function') {
        window.markMessageUnread(senderId);
    }
};

export const refreshOnlineUsers = () => {
    if (typeof window.loadOnlineUsers === 'function') {
        window.loadOnlineUsers();
    }
};
