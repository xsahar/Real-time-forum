// export class WebSocketClient {
//     constructor() {
//         this.socket = null;
//         this.messageHandlers = new Set();
//         this.statusHandlers = new Set();
//     }

//     connect() {
//         const token = localStorage.getItem('token');
//         if (!token) return;

//         const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
//         this.socket = new WebSocket(`${protocol}//${window.location.host}/api/ws`);

//         this.socket.onmessage = (event) => {
//             const message = JSON.parse(event.data);
//             switch (message.type) {
//                 case 'message':
//                     this.messageHandlers.forEach(handler => handler(message.payload));
//                     break;
//                 case 'status':
//                     this.statusHandlers.forEach(handler => handler(message.payload));
//                     break;
//             }
//         };

//         this.socket.onclose = () => {
//             // Attempt to reconnect after 5 seconds
//             setTimeout(() => this.connect(), 5000);
//         };
//     }

//     onMessage(handler) {
//         this.messageHandlers.add(handler);
//     }

//     onStatus(handler) {
//         this.statusHandlers.add(handler);
//     }

//     sendMessage(receiverId, content) {
//         if (this.socket?.readyState === WebSocket.OPEN) {
//             this.socket.send(JSON.stringify({
//                 type: 'message',
//                 payload: {
//                     receiverId,
//                     content
//                 }
//             }));
//         }
//     }
// }