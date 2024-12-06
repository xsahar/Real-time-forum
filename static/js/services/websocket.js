// static/js/services/websocket.js
export class WebSocketService {
    constructor() {
        this.baseURL = `ws://${window.location.host}/ws`;
        this.handlers = new Map();
        this.connected = false;
        this.retryCount = 0;
        this.maxRetries = 3;
        this.messageQueue = [];
        
        this.connect();
    }

    connect() {
        try {
            const token = localStorage.getItem('auth_token');
            if (!token) {
                throw new Error('No auth token found');
            }

            this.ws = new WebSocket(`${this.baseURL}?token=${token}`);
            this.setupEventHandlers();
        } catch (error) {
            console.error('WebSocket connection failed:', error);
        }
    }

    setupEventHandlers() {
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.connected = true;
            this.retryCount = 0;
            this.processMessageQueue();
        };

        this.ws.onclose = (event) => {
            this.connected = false;
            if (event.code === 1000) {
                console.log('WebSocket closed normally');
                return;
            }
            this.handleReconnect();
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.connected = false;
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                const handler = this.handlers.get(message.type);
                if (handler) handler(message.data);
            } catch (error) {
                console.error('Failed to handle message:', error);
            }
        };
    }

    handleReconnect() {
        if (this.retryCount >= this.maxRetries) {
            console.error('Max reconnection attempts reached');
            return;
        }

        const delay = Math.min(1000 * Math.pow(2, this.retryCount), 5000);
        this.retryCount++;

        console.log(`Reconnecting in ${delay}ms...`);
        setTimeout(() => this.connect(), delay);
    }

    send(type, data) {
        const message = JSON.stringify({ type, data });
        
        if (!this.connected) {
            this.messageQueue.push(message);
            return;
        }

        this.ws.send(message);
    }

    processMessageQueue() {
        while (this.messageQueue.length > 0 && this.connected) {
            const message = this.messageQueue.shift();
            this.ws.send(message);
        }
    }

    on(type, handler) {
        this.handlers.set(type, handler);
    }

    disconnect() {
        if (this.ws) {
            this.ws.close(1000, 'User logout');
            this.connected = false;
            this.handlers.clear();
            this.messageQueue = [];
        }
    }
}