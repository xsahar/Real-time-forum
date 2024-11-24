# Posts and Comments System

## Post Management
```mermaid
graph TD
    A[Create Post] --> B[Select Category]
    B --> C[Submit Post]
    C --> D[Display in Feed]
    D --> E[Enable Comments]


### Data Flow

sequenceDiagram
    User->>Frontend: Create Post/Comment
    Frontend->>Backend: Send Data
    Backend->>Database: Store
    Backend->>Frontend: Confirm
    Frontend->>UI: Update Display
```  

### Features:
1. Posts
- Category selection
- Title and content
- Author information
- Timestamp
- Like/Dislike system

2. Comments
- Nested under posts
- Author attribution
- Timestamp
- Show/Hide functionality