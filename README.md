# Real-Time Forum

A single-page application forum with real-time messaging capabilities.

## Architecture
```mermaid
graph TD
    subgraph Frontend
        A[SPA - Single Page App]
        B[WebSocket Client]
        C[REST Client]
    end
    
    subgraph Backend
        D[Go Server]
        E[WebSocket Hub]
        F[SQLite DB]
    end
    
    A --> B & C
    B <-->|Real-time| E
    C <-->|HTTP| D
    D & E --> F

### Features
- User authentication (register/login)
- Posts with categories
- Comments and likes system
- Real-time private messaging
- Online/offline status
- Dynamic content loading

### Setup
1. Clone repository
```
git clone <repository-url>
```

2. Start with Docker
```
docker-compose up --build
```

3. Access application
```
http://localhost:8080
```

### Documentation
Detailed documentation available in #Docs folder:

- Setup & Architecture (1-Setup/)
- Communication Protocols (2-Communication/)
- Core Features (3-Features/)
- Database Design (4-Database/)

### Tech Stack
- Frontend: Vanilla JavaScript, HTML5, CSS3
- Backend: Go 1.16+
- Database: SQLite
- Real-time: WebSocket
- Container: Docker
