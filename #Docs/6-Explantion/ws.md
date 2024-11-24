# WebSocket Server Implementation Guide
## Architecture Overview

The system consists of two main components:

1. Hub - Central message coordinator
2. Client - Individual connection handler

### Hub Component (hub.go)
Think of the Hub as a central control room.

### Primary Functions
- Track connected clients (like a guest list)
- Manage message broadcasting (like a PA system)
- Handle connections/disconnections
- Coordinate communication channels

### Key Components

1. Client Management
- Maintains client map
- Uses mutex for thread safety
- Handles registration/unregistration

2. Message System
Supports multiple message types:
- Chat
- Status
- Typing
 Handles targeted and broadcast messages
Manages status updates

3. Communication Channels
- broadcast: Message distribution
- register: New connections
- unregister: Disconnections

## Client Component (client.go)

### Overview

Represents a single WebSocket connection handler.

### Key Features

1. Identity Management
- Unique UserID per client
- Hub connection reference

2. Communication Handling
- Dedicated message channel
- ReadPump for incoming messages
- WritePump for outgoing messages

### Timing Constants

1. writeWait (10s)
- Maximum message write time
- Prevents server hanging
- Handles slow clients

2. pongWait (60s)
- Maximum ping-pong interval
- Connection health check
- Dead connection detection

3. pingPeriod (54s)
- Ping message frequency
- Set to 90% of pongWait
- Maintains connection health

## Operation Flow

### Connection Lifecycle

1. New Connection
- Client creation
- Hub registration
- Status broadcast

2. Active Operation
- Message routing through Hub
- Client message distribution
- Connection health monitoring

3. Disconnection
- Client unregistration
- Resource cleanup
- Status update broadcast

### Real-world Analogies

- `writeWait:` Waiter awaiting customer response
- `pongWait:` Phone line connection check
- `pingPeriod:` Regular "are you there?" check

### Critical Timing Benefits
- Stable connection maintenance
- Quick disconnection detection
- Resource leak prevention
- Reliable real-time communication