# Project Structure Documentation

## Core Components

### 1. Backend Organization
```plaintext
src/
├── config/          # Application configuration
├── database/        # Database operations
│   ├── db_init.go   # Database initialization
│   └── setup.go     # Connection management
├── handlers/        # Request handlers
│   ├── auth/        # Authentication
│   ├── posts/       # Post management
│   └── ws/          # WebSocket handling
├── middleware/      # HTTP middleware
├── models/          # Data structures
│   └── structs.go   # Core models
└── utils/          # Shared utilities

Testing Structure
tests/
├── unit/
│   └── database/    # Database tests
└── integration/     # Integration tests

✅ Database Setup

Schema definition
Table creation
Category initialization
Unit tests
🔄 Data Models

Defining structures
Adding validations
JSON handling
⏳ Next Steps

Implement handlers
Add middleware
Setup WebSocket