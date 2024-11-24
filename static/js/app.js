import { APIService } from './services/api.js';
import { WebSocketService } from './services/websocket.js';
import { ChatComponent } from './components/chat.js';
import { PostsComponent } from './components/posts.js';
import { UsersComponent } from './components/users.js';

export class App {
    constructor() {
        this.api = new APIService();
        this.components = {};
        this.setupRoutes();
    }

    setupRoutes() {
        // Define routes
        this.routes = {
            '/': () => this.showComponent('posts'),
            '/chat': () => this.showComponent('chat'),
            '/login': () => this.showComponent('auth', 'login'),
            '/register': () => this.showComponent('auth', 'register'),
            '/profile': () => this.showComponent('users')
        };

        // Handle navigation
        window.addEventListener('popstate', () => this.handleRoute());
    }

    async initialize() {
        try {
            // Check auth status
            const isAuthenticated = await this.checkAuth();
            
            // Initialize components
            this.components = {
                posts: new PostsComponent(this),
                chat: isAuthenticated ? new ChatComponent(this) : null,
                users: new UsersComponent(this)
            };

            // Initialize WebSocket if authenticated
            if (isAuthenticated) {
                this.ws = new WebSocketService();
            }

            // Handle initial route
            this.handleRoute();

        } catch (error) {
            console.error('Initialization failed:', error);
        }
    }

    async handleRoute() {
        const path = window.location.pathname;
        const route = this.routes[path] || this.routes['/'];
        
        // Hide all containers
        Object.keys(this.components).forEach(name => {
            const container = document.getElementById(`${name}-container`);
            if (container) {
                container.style.display = 'none';
            }
        });

        // Execute route handler
        await route();
    }

    async showComponent(name, mode = null) {
        const component = this.components[name];
        if (!component) return;

        const container = document.getElementById(`${name}-container`);
        if (container) {
            container.style.display = 'block';
            if (mode) {
                component.setMode?.(mode);
            }
            await component.mount?.();
        }
    }

    async checkAuth() {
        const token = localStorage.getItem('auth_token');
        return !!token;
    }
}