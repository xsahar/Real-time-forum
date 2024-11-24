# WebSocket Documentation

## 1. Overview: HTTP vs WebSocket
```
graph TD
    subgraph HTTP
        A[Request] -->|1. Ask| B[Server]
        B -->|2. Reply| A
        A -->|3. Connection Closed| C((End))
    end
    
    subgraph WebSocket
        D[Client] <-->|Persistent Connection| E[Server]
        E -->|Real-time Updates| D
    end

```
## 2. Connection Flow
```
stateDiagram-v2
    [*] --> Connect: HTTP Upgrade
    Connect --> Authenticate: Validate Session
    Authenticate --> Active: Success
    Active --> MessageFlow: Start Communication
    MessageFlow --> Active: Continue
    Active --> Disconnected: Connection Lost
    Disconnected --> Connect: Retry
```

### Implementation Order:
```
graph TD
    A[Hub] --> B[Client Management]
    B --> C[Message Handling]
    C --> D[Real-time Features]
    D --> E[Integration]
```

### 3. Message Types & Flow
```
sequenceDiagram
    participant U1 as User 1
    participant S as Server
    participant U2 as User 2
    
    U1->>S: Connect
    U2->>S: Connect
    Note over S: Broadcast Online Status
    
    U1->>S: Send Message
    S->>U2: Deliver Message
    
    U1->>S: Start Typing
    S->>U2: Typing Indicator
    
    U1->>S: Disconnect
    Note over S: Broadcast Offline Status

 ```
### 4. Data Structure
```
{
    "type": "message|status|typing",
    "data": {
        "content": "message text",
        "timestamp": "2024-03-20T10:30:00Z"
    },
    "from": "user_id",
    "to": "user_id"
}
```
### 5. Key Components
```
graph TD
    A[Hub] -->|Manages| B[Clients]
    A -->|Handles| C[Messages]
    A -->|Tracks| D[Status]
    
    B -->|Send| E[Real-time Updates]
    B -->|Receive| F[User Actions]
    
    C -->|Broadcast| G[All Users]
    C -->|Direct| H[Specific User]
```

## 6. Features

1. Real-time Communication
- Private messaging
- Status updates
- Typing indicators

2. State Management
- Online/Offline tracking
- Message history
- Connection recovery

3. Security
- Session validation
- Input sanitization
- Rate limiting

### progress:
1. Authentication working ✅
2. CRUD operations for posts ✅
3. Documentation structure ✅