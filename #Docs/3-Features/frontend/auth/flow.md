# Authentication Flow

## Login Process
```mermaid
sequenceDiagram
    User->>Frontend: Enter Credentials
    Frontend->>API: POST /api/login
    API->>Database: Validate
    Database->>API: User Data
    API->>Frontend: JWT Token
    Frontend->>User: Redirect Home
```

## Registration Process
```
sequenceDiagram
    User->>Frontend: Fill Registration Form
    Frontend->>API: POST /api/register
    API->>Database: Create User
    Database->>API: Success/Error
    API->>Frontend: Response
    Frontend->>User: Show Result

```

```
sequenceDiagram
    participant U as User
    participant AC as Auth Component
    participant S as Server
    
    U->>AC: Enter Credentials
    AC->>S: POST /api/login
    S->>AC: Return Token
    AC->>U: Show Success
    Note over AC: Store Token
    AC->>U: Redirect Home
    ```