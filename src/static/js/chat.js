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
            showChatMessages(userId);
        }
    });

    // Handle message sending
    sendButton.addEventListener('click', () => sendMessage());
    messageInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            sendMessage();
        }
    });

    // Add back button handler
    backButton.addEventListener('click', () => {
        chatMessages.style.display = 'none';
        usersList.style.display = 'block';
        backButton.style.display = 'none';
        currentChatUser = null;
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
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) throw new Error('Failed to load users');

            const users = await response.json();
            renderUsers(users);
        } catch (error) {
            console.error('Error loading users:', error);
            usersList.innerHTML = '<div class="error">Failed to load users</div>';
        }
    }

    function renderUsers(users) {
        usersList.innerHTML = users.map(user => `
            <div class="user-item" data-user-id="${user.id}">
                <span class="user-status ${user.online ? 'online' : 'offline'}"></span>
                <span class="user-name">${user.nickname}</span>
            </div>
        `).join('');
    }

    function showChatMessages(userId) {
        usersList.style.display = 'none';
        chatMessages.style.display = 'flex';
        backButton.style.display = 'block';
        loadChatHistory(userId);
    }

    async function loadChatHistory(userId) {
        try {
            const response = await fetch(`/api/messages/${userId}`, {
                headers: {
                    'Authorization': `Bearer ${token}`
                }
            });

            if (!response.ok) throw new Error('Failed to load messages');

            const messages = await response.json();
            renderMessages(messages);
        } catch (error) {
            console.error('Error loading messages:', error);
            messagesContainer.innerHTML = '<div class="error">Failed to load messages</div>';
        }
    }

    function renderMessages(messages) {
        messagesContainer.innerHTML = messages.map(message => `
            <div class="message ${message.senderId === currentChatUser ? 'received' : 'sent'}">
                <div class="message-content">${message.content}</div>
                <div class="message-time">${new Date(message.createdAt).toLocaleTimeString()}</div>
            </div>
        `).join('');
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    async function sendMessage() {
        if (!currentChatUser || !messageInput.value.trim()) return;

        try {
            const response = await fetch('/api/messages', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({
                    receiverId: currentChatUser,
                    content: messageInput.value.trim()
                })
            });

            if (!response.ok) throw new Error('Failed to send message');

            messageInput.value = '';
            loadChatHistory(currentChatUser);
        } catch (error) {
            console.error('Error sending message:', error);
        }
    }
};