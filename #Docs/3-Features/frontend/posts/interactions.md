# Post Interactions

## Like System
```mermaid
sequenceDiagram
    User->>Post: Click Like
    Post->>API: POST /api/like/{id}
    API->>Database: Update Count
    Database->>API: New Count
    API->>Post: Update UI
```
## Comment System
```
graph TD
    A[Add Comment] --> B[Store in DB]
    B --> C[Update UI]
    D[Load Comments] --> E[Display Thread]
```

## Implementation
```
{
    "interactions": {
        "like": "Handle post likes",
        "comment": "Manage comments",
        "share": "Share functionality"
    }
}
```