// export class ChatComponent {
//     constructor(app) {
//         this.app = app;
//         this.ws = app.ws;
//         this.container = document.getElementById('chat-container');
//         this.messages = new Map(); // Store messages by user
//     }

//     async mount() {
//         this.setupUI();
//         this.setupWebSocket();
//         await this.loadContacts();
//     }

//     setupUI() {
//         this.container.innerHTML = `
//             <div class="chat-wrapper">
//                 <div class="contacts-list"></div>
//                 <div class="chat-messages">
//                     <div class="message-list"></div>
//                     <form class="message-form">
//                         <input type="text" placeholder="Type a message..." required>
//                         <button type="submit">Send</button>
//                     </form>
//                 </div>
//             </div>
//         `;

//         this.contactsList = this.container.querySelector('.contacts-list');
//         this.messageList = this.container.querySelector('.message-list');
//         this.messageForm = this.container.querySelector('.message-form');

//         this.messageForm.addEventListener('submit', (e) => {
//             e.preventDefault();
//             const input = e.target.querySelector('input');
//             if (input.value.trim()) {
//                 this.sendMessage(input.value);
//                 input.value = '';
//             }
//         });
//     }

//     setupWebSocket() {
//         this.ws.on('chat', (message) => {
//             this.receiveMessage(message);
//         });

//         this.ws.on('status', (data) => {
//             this.updateUserStatus(data);
//         });
//     }

//     async loadContacts() {
//         try {
//             const response = await this.app.api.get('users/online');
//             this.renderContacts(response.data);
//         } catch (error) {
//             console.error('Failed to load contacts:', error);
//         }
//     }

//     renderContacts(contacts) {
//         this.contactsList.innerHTML = contacts.map(contact => `
//             <div class="contact ${contact.online ? 'online' : ''}" data-id="${contact.id}">
//                 <span>${contact.nickname}</span>
//                 <span class="status-dot"></span>
//             </div>
//         `).join('');
//     }

//     sendMessage(content) {
//         if (!this.activeContact) return;
        
//         const message = {
//             content,
//             to: this.activeContact,
//             timestamp: new Date()
//         };

//         this.ws.send('chat', message);
//         this.appendMessage(message, true);
//     }

//     receiveMessage(message) {
//         this.appendMessage(message, false);
//     }

//     appendMessage(message, isOwn) {
//         const messageElement = document.createElement('div');
//         messageElement.className = `message ${isOwn ? 'own' : ''}`;
//         messageElement.innerHTML = `
//             <div class="message-content">${message.content}</div>
//             <div class="message-time">${new Date(message.timestamp).toLocaleTimeString()}</div>
//         `;
//         this.messageList.appendChild(messageElement);
//         this.messageList.scrollTop = this.messageList.scrollHeight;
//     }

//     updateUserStatus(data) {
//         const contact = this.contactsList.querySelector(`[data-id="${data.userId}"]`);
//         if (contact) {
//             contact.classList.toggle('online', data.status === 'online');
//         }
//     }
// }