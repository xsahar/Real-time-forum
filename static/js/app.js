import { APIService } from '/static/js/services/api.js'; 
import { WebSocketService } from './services/websocket.js';
import { ChatComponent } from './components/chat.js';
import { PostsComponent } from './components/posts.js';
import { UsersComponent } from './components/users/index.js';

export class App {
    constructor() {
        this.api = new APIService();
        this.components = {};
        this.setupRoutes();
    }

    setupRoutes() {
        // Define routes with auth requirements and components
        this.routes = {
            '/': {
                component: 'posts',
                requireAuth: true,
                handler: async () => {
                    await this.showComponent('posts');
                    document.getElementById('content').style.display = 'block';
                }
            },
            '/chat': {
                component: 'chat',
                requireAuth: true,
                handler: async () => {
                    await this.showComponent('chat');
                    document.getElementById('chat-container').style.display = 'block';
                }
            },
            '/new-post': {
                component: 'posts',
                requireAuth: true,
                handler: async () => {
                    await this.showComponent('posts', 'create');
                }
            },
            '/login': {
                component: 'users',
                requireAuth: false,
                handler: async () => {
                    await this.showComponent('users', 'login');
                }
            },
            '/register': {
                component: 'users',
                requireAuth: false,
                handler: async () => {
                    await this.showComponent('users', 'register');
                }
            }
        };

        // Handle navigation
        window.addEventListener('popstate', () => this.handleRoute());
    }

    async initialize() {
        try {
            const isAuthenticated = await this.checkAuth();
            
            // Initialize components based on auth status
            this.components = {
                posts: new PostsComponent(this),
                chat: isAuthenticated ? new ChatComponent(this) : null,
                users: new UsersComponent(this)
            };

            // Setup WebSocket for authenticated users
            if (isAuthenticated) {
                this.ws = new WebSocketService();
                document.querySelector('.main-nav').style.display = 'flex';
            }

            await this.handleRoute();
        } catch (error) {
            console.error('Initialization failed:', error);
        }
    }

    async handleRoute() {
        const path = window.location.pathname;
        const route = this.routes[path] || this.routes['/'];

        // Check auth requirements
        if (route.requireAuth && !(await this.checkAuth())) {
            window.location.href = '/login';
            return;
        }

        // Hide all containers
        Object.keys(this.components).forEach(name => {
            const container = document.getElementById(`${name}-container`);
            if (container) {
                container.style.display = 'none';
            }
        });

        // Execute route handler
        await route.handler();
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
        try {
            const token = localStorage.getItem('auth_token');
            if (!token) return false;
            
            const response = await this.api.get('auth/status');
            return response.status === 'success';
        } catch (error) {
            console.error('Auth check error:', error);
            localStorage.removeItem('auth_token');
            return false;
        }
    }
}