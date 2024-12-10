import { ContactsManager } from './contacts.js';
import { MessageManager } from './messages.js';

export class ChatComponent {
    constructor(app) {
        this.app = app;
        this.ws = app.ws;
        this.container = document.getElementById('chat-container');
        this.setupUI();
        
        // Initialize managers
        this.contacts = new ContactsManager(this.container, app.api);
        this.messages = new MessageManager(this.container, this.ws);
    }

    async mount() {
        this.setupUI();
        this.setupWebSocket();
        await this.contacts.loadContacts();
    }

    setupUI() {
        this.container.innerHTML = `
            <div class="chat-wrapper">
                <div class="contacts-list"></div>
                <div class="chat-messages">
                    <div class="message-list"></div>
                    <form class="message-form">
                        <input type="text" placeholder="Type a message..." required>
                        <button type="submit">Send</button>
                    </form>
                </div>
            </div>
        `;
    }

    setupWebSocket() {
        this.ws.on('chat', (message) => {
            this.messages.receiveMessage(message);
        });

        this.ws.on('status', (data) => {
            this.contacts.updateUserStatus(data);
        });
    }
}