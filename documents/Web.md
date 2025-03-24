The WebSocket connection showing "101 Type: plain" in the network inspector is not an issue - it's actually the expected behavior. This is a good sign that your WebSocket connection is working correctly.

Why WebSocket shows "101" and "plain" (Not an Issue)
The "101" status code represents a successful protocol switch from HTTP to WebSocket. Here's what's happening:

Initial HTTP Request: When establishing a WebSocket connection, your browser first makes a standard HTTP request to the server.

Upgrade Protocol: In this request, it includes special headers asking to "upgrade" the connection from HTTP to WebSocket protocol.

101 Status Code: If the server agrees to the upgrade, it responds with HTTP status code 101 ("Switching Protocols").

Connection Established: After this handshake, the connection is no longer using HTTP but has switched to the WebSocket protocol.

The "plain" content type is also normal for WebSocket connections since they don't use standard HTTP content types like "application/json" after the connection is established.

Comparison with Your Other API Calls
Your other API endpoints (like /api/users/online) show JSON content type because:

They are standard REST API calls using HTTP
They return data formatted as JSON
They complete after a single request/response cycle
In contrast, WebSockets:

Start as HTTP but upgrade to WebSocket protocol (hence the 101)
Maintain a persistent connection
Don't use HTTP content types for the actual messages
Allow bi-directional communication
What WebSocket Messages Look Like
After the 101 upgrade, the WebSocket connection stays open. When you send or receive messages over this connection, they won't appear as separate network requests in the inspector.

The actual WebSocket messages (where your JSON data is being sent) happen over this persistent connection. To see these messages:

Look in the "Messages" or "Frames" tab in Chrome DevTools Network panel when inspecting the WebSocket connection
Look at the browser console logs where you're printing received messages
Conclusion
The 101 status with plain content type for your WebSocket connection is normal and indicates the connection was successfully established. Your API is working correctly, with REST endpoints returning JSON and WebSockets handling real-time communication.