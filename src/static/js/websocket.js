// src/static/js/websocket.js
export class WebSocketClient {
    constructor() {
        this.socket = null;
        this.messageHandlers = new Set();
        this.statusHandlers = new Set();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
    }

    connect() {
        const token = localStorage.getItem('token');
        if (!token) return;

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const url = `${protocol}//${window.location.host}/api/ws?token=${token}`;

        try {
            this.socket = new WebSocket(url);

            this.socket.onopen = () => {
                console.log('WebSocket connection established');
                this.reconnectAttempts = 0;
            };

            this.socket.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    console.log('WebSocket message received:', message.type);

                    switch (message.type) {
                        case 'message':
                            this.messageHandlers.forEach(handler => handler(message.payload));
                            break;
                        case 'status':
                            this.statusHandlers.forEach(handler => handler(message.payload));
                            break;
                        case 'error':
                            console.error('WebSocket error message:', message.payload);
                            break;
                        default:
                            console.warn('Unknown message type:', message.type);
                            break;
                    }
                } catch (error) {
                    console.error('Error parsing WebSocket message:', error);
                }
            };

            this.socket.onclose = (event) => {
                console.log('WebSocket connection closed', event.code, event.reason);

                // If closed with authentication error, don't reconnect
                if (event.code === 1008 || event.code === 1011) {
                    console.error('WebSocket authentication failed, not reconnecting');
                    return;
                }

                if (this.reconnectAttempts < this.maxReconnectAttempts) {
                    // Attempt to reconnect with exponential backoff
                    const delay = Math.min(1000 * Math.pow(2, this.reconnectAttempts), 30000);
                    console.log(`Attempting to reconnect in ${delay}ms`);
                    setTimeout(() => {
                        this.reconnectAttempts++;
                        this.connect();
                    }, delay);
                }
            };

            this.socket.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
        } catch (error) {
            console.error('Error creating WebSocket connection:', error);
        }
    }

    onMessage(handler) {
        this.messageHandlers.add(handler);
        return () => this.messageHandlers.delete(handler);
    }

    onStatus(handler) {
        this.statusHandlers.add(handler);
        return () => this.statusHandlers.delete(handler);
    }

    sendMessage(receiverId, content) {
        if (this.socket?.readyState === WebSocket.OPEN) {
            console.log(`Sending message to ${receiverId}: ${content.substring(0, 20)}${content.length > 20 ? '...' : ''}`);

            const message = {
                type: 'message',
                payload: {
                    receiverId,
                    content
                }
            };

            try {
                this.socket.send(JSON.stringify(message));
                console.log('Message sent to WebSocket server');
                return true;
            } catch (error) {
                console.error('Error sending WebSocket message:', error);
                return false;
            }
        } else {
            console.warn(`Cannot send message - WebSocket not connected (state: ${this.socket?.readyState})`);
            return false;
        }
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            this.socket = null;
        }
    }
}