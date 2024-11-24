[`#Docs/2-Communication/rest-api.md`](#Docs/2-Communication/rest-api.md ):
```markdown
# REST API Documentation

## Authentication Endpoints
- POST `/api/register`
- POST `/api/login`
- POST `/api/logout`

## Content Endpoints

graph LR
    A[Posts] --> B[Create]
    A --> C[Read]
    A --> D[Update]
    A --> E[Delete]
```

## API Structure

Posts Management

GET /api/posts
POST /api/create-post
GET /api/post/:id
POST /api/like
Comments

POST /api/comment
GET /api/comments/:postId
Categories

GET /api/categories
POST /api/add-category

### Response Format
```
{
    "status": "success|error",
    "data": {},
    "message": "Response message"
}
```
# Posts API Documentation

## Endpoints

### Create Post
```http
POST /api/create-post
Content-Type: application/json

{
    "title": "string",
    "content": "string"
}

Response (201 Created):
{
    "status": "success",
    "message": "Post created successfully",
    "data": post_id
}

GET /api/post/{id}

Response (200 OK):
{
    "status": "success",
    "data": {
        "id": number,
        "title": "string",
        "content": "string",
        "author": "string",
        "likes": number,
        "dislikes": number,
        "categories": array,
        "created_at": "datetime"
    }
}

PUT /api/post/{id}
Content-Type: application/json

{
    "title": "string",
    "content": "string"
}

Response (200 OK):
{
    "status": "success",
    "message": "Post updated successfully"
}

DELETE /api/post/{id}

Response (200 OK):
{
    "status": "success",
    "message": "Post deleted successfully"
}

Error Responses
{
    "status": "error",
    "message": "Error description"
}

Test Coverage
Create post validation
Get post with details
Update post content
Delete post cascade
Error handling
Authentication checks