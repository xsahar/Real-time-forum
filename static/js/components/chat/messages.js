export class MessageManager {
    constructor(container, ws) {
        this.container = container;
        this.ws = ws;
        this.messageList = container.querySelector('.message-list');
        this.messageForm = container.querySelector('.message-form');
        this.messages = new Map();
        this.setupListeners();
    }

    setupListeners() {
        this.messageForm.addEventListener('submit', (e) => {
            e.preventDefault();
            const input = e.target.querySelector('input');
            if (input.value.trim()) {
                this.sendMessage(input.value);
                input.value = '';
            }
        });
    }

    sendMessage(content) {
        if (!this.activeContact) return;
        
        const message = {
            content,
            to: this.activeContact,
            timestamp: new Date()
        };

        this.ws.send('chat', message);
        this.appendMessage(message, true);
    }

    receiveMessage(message) {
        this.appendMessage(message, false);
    }

    appendMessage(message, isOwn) {
        const messageElement = document.createElement('div');
        messageElement.className = `message ${isOwn ? 'own' : ''}`;
        messageElement.innerHTML = `
            <div class="message-content">${message.content}</div>
            <div class="message-time">${new Date(message.timestamp).toLocaleTimeString()}</div>
        `;
        this.messageList.appendChild(messageElement);
        this.messageList.scrollTop = this.messageList.scrollHeight;
    }
}