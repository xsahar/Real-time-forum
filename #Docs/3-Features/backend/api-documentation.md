
3. [`#Docs/api-documentation.md`](#Docs/api-documentation.md ):
```markdown
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
```

### Login
```
POST /api/login
Content-Type: application/json

{
  "email": "string",
  "password": "string"
}
```
### Posts, Create Post
```
POST /api/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "string",
  "content": "string",
  "category": "number"
}
```
## Error Response Format
```json
{
    "status": "error",
    "message": "Error description"
}


GitHub Copilot
Let's document the error handling implementation and testing in the appropriate files:

api-documentation.md:
Standard Error Types
Server Error (500)

Internal server errors
Panic recovery
Unexpected errors
Not Found Error (404)

Resource not found
Invalid routes
Headers
Content-Type: application/json
Status Code: Appropriate HTTP status