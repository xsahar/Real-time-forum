// Main application class that handles routing, component management, and authentication
import { APIService } from '/static/js/services/api.js';
import { WebSocketService } from './services/websocket.js';
import { ChatComponent } from './components/chat.js';
import { PostsComponent } from './components/posts.js';
import { UsersComponent } from './components/users/index.js';

export class App {
    constructor() {
        // Initialize API service for backend communication
        this.api = new APIService();
        // Store component instances
        this.components = {};
        this.setupRoutes();
        this.setupAuthListeners();
    }

    setupAuthListeners() {
        window.addEventListener('login-success', (event) => {
            const { token, user } = event.detail;
            if (token) {
                console.log('Auth token received');
                this.handleAuthentication(token, user);
            }
        });
    }

    async handleAuthentication(token, user) {
        // Add debug log
        console.log('Storing auth token:', token ? 'present' : 'missing');
        localStorage.setItem('auth_token', token);
        console.log('Token stored, verifying:', localStorage.getItem('auth_token') ? 'present' : 'missing');
        
        document.querySelector('.main-nav').style.display = 'flex';
        await this.handleRoute();
    }

    setupRoutes() {
        // Define application routes and their handlers
        this.routes = {
            '/': {
                component: 'posts',
                requireAuth: true,
                // Async: Waits for posts component to load and mount
                handler: async () => {
                    await this.showComponent('posts');
                    document.getElementById('content').style.display = 'block';
                }
            },
            '/login': {
                component: 'users',
                requireAuth: false,
                // Async: Handles login form display and mounting
                handler: async () => {
                    await this.showComponent('users', 'login');
                }
            },
            '/register': {
                component: 'users',
                requireAuth: false,
                // Async: Handles registration form display and mounting
                handler: async () => {
                    await this.showComponent('users', 'register');
                }
            },
            '/chat': {
                component: 'chat',
                requireAuth: true,
                // Async: Initializes chat component and WebSocket
                handler: async () => {
                    await this.showComponent('chat');
                }
            }
        };

        // Listen for browser navigation events
        window.addEventListener('popstate', () => this.handleRoute());
    }

    // Async: Initializes application state and components
    async initialize() {
        try {
            // Check if user is authenticated
            const isAuthenticated = await this.checkAuth();
            
            // Initialize components based on auth state
            this.components = {
                posts: new PostsComponent(this),
                chat: isAuthenticated ? new ChatComponent(this) : null,
                users: new UsersComponent(this)
            };

            // Setup WebSocket for real-time communication if authenticated
            if (isAuthenticated) {
                this.ws = new WebSocketService();
                document.querySelector('.main-nav').style.display = 'flex';
            }

            // Handle initial route
            await this.handleRoute();
            console.log('App initialized successfully');
        } catch (error) {
            console.error('Initialization failed:', error);
        }
    }

    // Async: Handles route changes and authentication checks
    async handleRoute() {
        const path = window.location.pathname;
        const route = this.routes[path] || this.routes['/'];

        // Redirect to login if auth required but not authenticated
        if (route.requireAuth && !(await this.checkAuth())) {
            history.pushState(null, '', '/login');
            await this.handleRoute();
            return;
        }

        await route.handler();
    }

    // Async: Manages component visibility and mounting
    async showComponent(name, mode = null) {
        const component = this.components[name];
        if (!component) return;

        // Hide all other components
        Object.keys(this.components).forEach(key => {
            const container = document.getElementById(`${key}-container`);
            if (container) {
                container.style.display = 'none';
            }
        });

        // Show and initialize requested component
        const container = document.getElementById(`${name}-container`);
        if (container) {
            container.style.display = 'block';
            if (mode) {
                component.setMode?.(mode);
            }
            await component.mount?.();
        }
    }

    // Async: Verifies user authentication status
    async checkAuth() {
        try {
            const token = localStorage.getItem('auth_token');
            console.log('Checking auth token:', token ? 'present' : 'missing');

            if (!token){
             return false;   
            } 

            console.log('Sending auth status request');
            const response = await this.api.get('auth/status');
            console.log('Auth status response:', response);
            
            return response.status === 'success';
        } catch (error) {
            console.error('Auth check error:', error);
            localStorage.removeItem('auth_token');
            return false;
        }
}
}