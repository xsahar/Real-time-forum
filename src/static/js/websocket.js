import { refreshOnlineUsers } from './chat.js';

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
                
                // Refresh users list when connection is established to see current online users
                refreshOnlineUsers();
            };

            this.socket.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    console.log('WebSocket message received:', message.type);

                    switch (message.type) {
                        case 'message':
                            // Check if this is a message we received (not sent by us)
                            const currentUser = JSON.parse(localStorage.getItem('user'));
                            if (!currentUser) break;
                            
                            const isSender = message.payload.senderId === currentUser.id;
                            
                            // Process the message with all handlers
                            this.messageHandlers.forEach(handler => handler(message.payload));
                            
                            // Immediately handle notification for received messages
                            if (!isSender) {
                                // This is an INCOMING message to us
                                console.log("Received message from:", message.payload.senderId);
                                
                                // Immediately try to handle notification
                                this.refreshWithRetries(message.payload);
                            }
                            break;
                        
                        case 'status':
                            this.statusHandlers.forEach(handler => handler(message.payload));
                            // Refresh users list when status changes
                            refreshOnlineUsers();
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
    
       refreshWithRetries(messagePayload) {
        setTimeout(() => {
            refreshOnlineUsers();
        }, 300);
    }
    
    addNotificationDot(senderId) {
        try {
            // Get all matching user items (there should be just one)
            const userItems = document.querySelectorAll(`.user-item[data-user-id="${senderId}"]`);
            
            if (userItems.length > 0) {
                userItems.forEach(userItem => {
                    // Check if notification dot already exists
                    if (!userItem.querySelector('.notification-dot')) {
                        const notificationDot = document.createElement('span');
                        notificationDot.className = 'notification-dot';
                        userItem.appendChild(notificationDot);
                        console.log(`ADDED notification dot for user ${senderId}`);
                    }
                });
                return true; // Successfully added notification
            } else {
                // If the user item doesn't exist in DOM yet, return false
                console.log(`User item for ${senderId} not found in DOM`);
                return false;
            }
        } catch (error) {
            console.error('Error adding notification dot:', error);
            return false;
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
                
                // Refresh users list after sending a message
                setTimeout(refreshOnlineUsers, 300);
                
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