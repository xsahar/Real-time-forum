# Database Testing Documentation

## Unit Tests

### Database Initialization Test
Location: `tests/unit/database/init_test.go`

#### Test Purpose
Validates database setup and initialization process.

#### What Gets Created
1. **Temporary Test Database**
   - File: `./test_forum.db`
   - Automatically cleaned up after tests

2. **Database Tables**
   - users
   - categories
   - posts
   - comments
   - likes
   - post_categories

3. **Initial Data**
   - Category records seeded
   - Connection verified
   - Schema validated

#### Test Flow
```mermaid
graph TD
    A[Setup Test DB] --> B[Init Database]
    B --> C[Verify Tables]
    C --> D[Check Categories]
    D --> E[Cleanup]