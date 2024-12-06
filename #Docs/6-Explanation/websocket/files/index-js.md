### Index.js

1. Optional Chaining (?.)
```
if (response?.status === 'success')
```

- ?. is a safe way to access nested properties
- Prevents "Cannot read property 'status' of undefined" errors
- Returns undefined if response is null/undefined
- Similar to writing: if (response && response.status === 'success')

2. History API Usage
```
history.pushState(null, '', '/login')
```
- Updates URL without page reload
- Manages browser history stack
- Parameters:
state object (null)
unused title ('')
new URL path ('/login')

3. Null Coalescing (||)
```
throw new Error(response?.error || 'Login failed')
```

- Provides fallback error message
- If response?.error is undefined/null, uses 'Login failed'
- Ensures user always sees meaningful error message

4. Token Management
```
localStorage.setItem('auth_token', response.data.token)
localStorage.removeItem('auth_token')
```
- Consistent token storage/removal
- Used for authentication state
- Cleared on errors/logout

5. Debug Logging
```
console.log('Processing login:', { 
    identifier: formData.identifier,
    hasPassword: !!formData.password 
})
```
- Enhanced debugging
- Sanitized password logging
- Shows important state changes

6. Response Validation
```
if (response?.status === 'success' && response.data?.token)
```
- Multiple conditions checked
- Ensures both status and token exist
- Prevents undefined access errors