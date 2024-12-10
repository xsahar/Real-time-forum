export class ContactsManager {
    constructor(container, api) {
        this.container = container;
        this.api = api;
        this.contactsList = container.querySelector('.contacts-list');
    }

    async loadContacts() {
        try {
            const response = await this.api.get('users/online');
            this.renderContacts(response.data);
        } catch (error) {
            console.error('Failed to load contacts:', error);
        }
    }

    renderContacts(contacts) {
        this.contactsList.innerHTML = contacts.map(contact => `
            <div class="contact ${contact.online ? 'online' : ''}" data-id="${contact.id}">
                <span>${contact.nickname}</span>
                <span class="status-dot"></span>
            </div>
        `).join('');
    }

    updateUserStatus(data) {
        const contact = this.contactsList.querySelector(`[data-id="${data.userId}"]`);
        if (contact) {
            contact.classList.toggle('online', data.status === 'online');
        }
    }
}