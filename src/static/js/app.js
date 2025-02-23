import { initAuth } from './auth.js';
import { initPosts } from './posts.js';
import { initUI } from './ui.js';
import { initChat } from './chat.js';

// State management
let isInitialized = false;

// Helper functions
const resetAuthState = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
};

const updateUIState = (authSection, mainSection, isAuthenticated) => {
    if (isAuthenticated) {
        authSection.classList.add('hidden');
        mainSection.classList.remove('hidden');
    } else {
        mainSection.classList.add('hidden');
        authSection.classList.remove('hidden');
    }
};

document.addEventListener('DOMContentLoaded', async () => {
    try {
        // Get and validate required elements
        const authSection = document.getElementById('auth-section');
        const mainSection = document.getElementById('main-section');
        const userInfo = document.querySelector('.user-info');

        if (!authSection || !mainSection) {
            throw new Error('Required sections not found');
        }

        // Initialize auth first and wait for result
        const authResult = await initAuth();
        console.log('Auth initialization complete');

        if (authResult?.isAuthenticated) {
            console.log('Valid session found:', { userId: authResult.user.id });

            // Single initialization block
            try {
                initUI();
                await initPosts();

                // Update UI state
                updateUIState(authSection, mainSection, true);
                if (userInfo) {
                    userInfo.textContent = `Welcome, ${authResult.user.nickname}!`;
                }

                // Initialize chat for authenticated users
                const token = localStorage.getItem('token');
                if (token) {
                    initChat();
                }

                isInitialized = true;
                console.log('Application initialization successful');
            } catch (error) {
                console.error('Feature initialization error:', error);
                throw new Error(`Failed to initialize application: ${error.message}`);
            }
        } else {
            console.log('No valid session found, showing auth form');
            updateUIState(authSection, mainSection, false);
            isInitialized = true;
        }

    } catch (error) {
        console.error('Application initialization error:', error);
        resetAuthState();

        const authSection = document.getElementById('auth-section');
        const mainSection = document.getElementById('main-section');

        if (authSection && mainSection) {
            updateUIState(authSection, mainSection, false);
        }

        alert(`Failed to initialize application: ${error.message}\nPlease refresh the page.`);
    }

    // Add error boundary for uncaught errors
    window.addEventListener('error', (event) => {
        console.error('Global error:', event.error);
        if (isInitialized) {
            alert('An unexpected error occurred. Please refresh the page.');
        }
    });
});