Unit Tests (/tests/unit/...)

Tests individual components in isolation
Mocks dependencies
Faster execution
Tests single functions/methods
Example: Testing Hub methods independently

Integration Tests (/tests/integration/...)
Tests multiple components working together
Uses real dependencies
Slower execution
Tests complete workflows
Example: Testing WebSocket client-server communication

go test command:
```
go test -v ./tests/integration/database -cover
```
Components:

1. go test - Base command to run Go tests
2. -v - Verbose flag, shows:
Individual test cases
Running time
Pass/fail status
Test output
3. database - Path to test files
Runs all *_test.go files in this directory
Only integration tests for database package
4. -cover - Enables code coverage reporting
Shows percentage of code executed by tests
Helps identify untested code paths

Example output:
```
=== RUN   TestDatabaseIntegration
=== RUN   TestDatabaseIntegration/Create_Tables
=== RUN   TestDatabaseIntegration/Initialize_Categories
--- PASS: TestDatabaseIntegration (0.15s)
    --- PASS: TestDatabaseIntegration/Create_Tables (0.05s)
    --- PASS: TestDatabaseIntegration/Initialize_Categories (0.10s)
coverage: 75.6% of statements
```
Alternative flags:

-coverprofile=coverage.out - Saves coverage data to file
-short - Runs short tests only
-timeout=30s - Sets test timeout
-run=TestName - Runs specific test