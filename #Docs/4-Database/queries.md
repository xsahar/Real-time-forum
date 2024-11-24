# Database Schema

## Entity Relationship

erDiagram
    USERS ||--o{ POSTS : creates
    USERS ||--o{ COMMENTS : writes
    POSTS ||--o{ COMMENTS : has
    CATEGORIES ||--o{ POSTS : contains

    USERS {
        int id
        string email
        string nickname
        string password
    }

    POSTS {
        int id
        int user_id
        string content
        timestamp created_at
    }

Key Relations
User Management
Content Storage
Interaction Tracking