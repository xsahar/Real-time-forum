// main.js
import { App } from './app.js';
import { handleNavigation } from './utils/navigation.js';
import { addEventListeners } from './utils/events.js';

let app = null;

document.addEventListener('DOMContentLoaded', async () => {
    try {
        if (!app) {
            console.log('Initializing app...');
            app = new App();
            await app.initialize();
            
            handleNavigation(app);
            addEventListeners(app);
            
            await app.handleRoute();
            console.log('App initialized successfully');
        }
    } catch (error) {
        console.error('App initialization failed:', error);
    }
});

// Cleanup on page unload
window.addEventListener('beforeunload', () => {
    app = null;
});