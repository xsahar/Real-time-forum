export class WebSocketService {
    constructor() {
        this.handlers = new Map();
        this.retryCount = 0;
        this.maxRetries = 5;
        this.connect();
    }

    connect() {
        try {
            const token = localStorage.getItem('auth_token');
            if (!token) {
                console.warn('No auth token found for WebSocket connection');
                return;
            }

            const wsUrl = `ws://${window.location.host}/ws?token=${token}`;
            this.ws = new WebSocket(wsUrl);
            this.setupEventHandlers();
        } catch (error) {
            console.error('WebSocket connection failed:', error);
            this.handleConnectionError();
        }
    }

    setupEventHandlers() {
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.retryCount = 0;
        };

        this.ws.onclose = (event) => {
            if (event.code === 1000) {
                console.log('WebSocket closed normally');
                return;
            }
            this.handleConnectionError();
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                const handler = this.handlers.get(message.type);
                if (handler) handler(message.data);
            } catch (error) {
                console.error('Failed to handle WebSocket message:', error);
            }
        };
    }

    handleConnectionError() {
        if (this.retryCount >= this.maxRetries) {
            console.error('Max reconnection attempts reached');
            return;
        }

        const delay = Math.min(1000 * Math.pow(2, this.retryCount), 30000);
        console.log(`Attempting to reconnect in ${delay/1000} seconds...`);
        
        setTimeout(() => {
            this.retryCount++;
            this.connect();
        }, delay);
    }

    on(type, handler) {
        this.handlers.set(type, handler);
    }

    send(type, data) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify({ type, data }));
        }
    }
}