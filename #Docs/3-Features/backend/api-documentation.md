# API Documentation

## Authentication

### Register User
```http
POST /api/register
Content-Type: application/json

{
    "email": "string",
    "nickname": "string",
    "password": "string",
    "age": "number",
    "gender": "string",
    "firstName": "string",
    "lastName": "string"
}

Response 200:
{
    "status": "success",
    "message": "Registration successful"
}

```
### Login
POST /api/login
Content-Type: application/json

{
    "identifier": "string", // email or nickname
    "password": "string"
}

Response 200:
{
    "status": "success",
    "token": "string"
}

## Posts
### Create Post
POST /api/posts
Authorization: Bearer <token>
Content-Type: application/json

{
    "title": "string",
    "content": "string",
    "categories": ["number"]
}

Response 200:
{
    "status": "success",
    "data": {
        "id": "number",
        "title": "string",
        "content": "string"
    }
}

### List Posts
GET /api/posts
Authorization: Bearer <token>

Response 200:
{
    "status": "success",
    "data": [
        {
            "id": "number",
            "title": "string",
            "content": "string",
            "author": "string",
            "likes": "number",
            "dislikes": "number",
            "categories": ["string"]
        }
    ]
}

### Post Interactions

POST /api/post/{id}/comment
Authorization: Bearer <token>
Content-Type: application/json

{
    "content": "string"
}

POST /api/likes
Authorization: Bearer <token>
Content-Type: application/json

{
    "post_id": "number",
    "value": "number" // 1 for like, -1 for dislike
}