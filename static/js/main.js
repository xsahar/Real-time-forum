// entery point of the application  
import { App } from './app.js';
import { handleNavigation } from './utils/navigation.js';
import { addEventListeners } from './utils/events.js';

let app = null;

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', async () => {
    try {
        if (!app) {
            console.log('Initializing app...');
            
            // Create app instance
            app = new App();
            
            // Initialize core functionality
            await app.initialize();
            
            // Setup navigation and event handlers
            handleNavigation(app);
            addEventListeners(app);
            
            // Handle initial route
            await app.handleRoute();
            
        }
    } catch (error) {
        console.error('App initialization failed:', error);
        // Show user-friendly error message
        document.getElementById('app').innerHTML = `
            <div class="error-message">
                Sorry, something went wrong. Please try refreshing the page.
            </div>
        `;
    }
});

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    if (app) {
        // Cleanup WebSocket connections
        app.ws?.disconnect();
        
        // Clear authentication
        localStorage.removeItem('auth_token');
        
        // Reset app instance
        app = null;
    }
});

// Handle errors globally
window.addEventListener('error', (event) => {
    console.error('Global error:', event.error);
    // Prevent default browser error handling
    event.preventDefault();
});

// Handle unhandled promise rejections
window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason);
    // Prevent default browser error handling
    event.preventDefault();
});