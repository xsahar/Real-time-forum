## API Endpoints Structure

1. Base URL Definition
```
constructor() {
    this.baseURL = '/api';  // All endpoints start with this
}
```

2. Endpoint USage

- Login: /api/login
- Register: /api/register
- Auth Status: /api/auth/status
- Posts: /api/posts
etc.

3. Template Literal Syntax
```
`${this.baseURL}/${endpoint}`
// Examples:
// endpoint = 'login' → '/api/login'
// endpoint = 'register' → '/api/register'
```

4. commmon Endpoints
```
const API_ENDPOINTS = {
    LOGIN: 'login',
    REGISTER: 'register',
    AUTH_STATUS: 'auth/status',
    POSTS: 'posts',
    USERS: 'users'
}
```
