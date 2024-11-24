# Database Schema Documentation

## Entity Relationships
```mermaid
erDiagram
    USERS ||--o{ POSTS : creates
    USERS ||--o{ COMMENTS : writes
    POSTS ||--o{ COMMENTS : has
    CATEGORIES ||--o{ POSTS : contains

Table Overview

1. Users
Purpose: Store user information and authentication data
Key Fields:
Primary Key: id (auto-increment)
Unique: email, nickname
Authentication: password (hashed)
Profile: age, gender, names

2. Posts
Purpose: Store user-created content
Key Fields:
Primary Key: id
Foreign Keys: user_id, category_id
Content: title, content
Metadata: created_at

3. Comments
Purpose: Enable post discussions
Key Fields:
Links: post_id, user_id
Content: text
Tracking: created_at

4. Likes
Purpose: Track user interactions
Key Fields:
Constraints: One like per user per post
Value: -1 (dislike) or 1 (like)

5. Categories
Purpose: Organize posts by topics
Key Fields:
Unique: name
Optional: description

Data Integrity

1.Constraints
Foreign key relationships
Unique field validation
Required field checks

2.Cascade Rules
Post deletion: Removes associated comments/likes
User deletion: Removes all user content

# Like System Schema

## Likes Table
```sql
CREATE TABLE likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER,
    user_id INTEGER,
    value INTEGER CHECK(value IN (-1,1)),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
)

Key Constraints
One reaction per user per post
Value must be -1 (dislike) or 1 (like)
Cascading deletes with posts/users

