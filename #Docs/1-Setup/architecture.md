# System Architecture

## Overview
```mermaid
graph TD
    subgraph Frontend
        UI[Single Page App]
        WC[WebSocket Client]
        RC[REST Client]
    end

    subgraph Backend
        WS[WebSocket Server]
        API[REST API]
        DB[(SQLite)]
    end

    UI --> WC & RC
    WC <--> WS
    RC <--> API
    WS & API --> DB