## 1.CORS Setup
Added setupCORS middleware function
Configures headers for cross-origin requests
Handles OPTIONS preflight requests
2. Route Configuration
Defined apiRoutes struct array
Added routes for posts, auth, comments, etc.
Each route specifies:
Path
Handler function
Rate limiting flag
3. Middleware Application
Applied CORS to all routes
Added recovery handler
Optional rate limiting per route
4. Static File Serving
Added SPA fallback route
Configured static file server
Strip prefix handling for static files
The code now properly handles cross-origin requests and applies consistent middleware across all API routes while maintaining proper route organization.