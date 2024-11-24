import { App } from './app.js';
import { handleNavigation } from './utils/navigation.js';
import { addEventListeners } from './utils/events.js';

// Initialize app when DOM is ready
document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Create app instance
        const app = new App();
        
        // Initialize app
        await app.initialize();
        
        // Setup navigation handling
        handleNavigation(app);
        
        // Setup global event listeners
        addEventListeners(app);
        
    } catch (error) {
        console.error('Failed to initialize app:', error);
    }
});