import { APIService } from './services/api.js';
import { WebSocketService } from './services/websocket.js';
import { PostsComponent } from './components/posts/index.js';  // Changed from posts.js
import { ChatComponent } from './components/chat/index.js';    // Changed from chat.js
import { UsersComponent } from './components/users/index.js';

export class App {
    constructor() {
        this.api = new APIService();
        this.components = {};
        this.setupRoutes();
        this.setupAuthListeners();

    }

    cleanup() {
        // Disconnect WebSocket
        if (this.ws) {
            this.ws.disconnect();
            this.ws = null;
        }

        // Clear auth data
        localStorage.removeItem('auth_token');

        // Clear components
        Object.values(this.components).forEach(component => {
            component?.cleanup?.();
        });
        this.components = {};

        // Reset API service
        this.api = null;

        console.log('App cleanup completed');
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
        console.log('Storing auth token:', token ? 'present' : 'missing');
        localStorage.setItem('auth_token', token);
        console.log('Token stored, verifying:', localStorage.getItem('auth_token') ? 'present' : 'missing');
        
        document.querySelector('.main-nav').style.display = 'flex';
        await this.handleRoute();
    }

    setupRoutes() {
        this.routes = {
            '/': {
                component: 'posts',
                requireAuth: true,
                handler: async () => {
                    document.getElementById('content').style.display = 'block';
                    await this.showComponent('posts');
                }
            },
            '/login': {
                component: 'users',
                requireAuth: false,
                handler: async () => {
                    document.getElementById('auth-container').style.display = 'block';
                    await this.showComponent('users', 'login');
                }
            },
            '/register': {
                component: 'users',
                requireAuth: false,
                handler: async () => {
                    await this.showComponent('users', 'register');
                }
            },
            '/chat': {
                component: 'chat',
                requireAuth: true,
                handler: async () => {
                    await this.showComponent('chat');
                }
            }
        };

        window.addEventListener('popstate', () => this.handleRoute());
    }

    async initialize() {
        try {
            const isAuthenticated = await this.checkAuth();
            
            this.components = {
                posts: new PostsComponent(this),
                chat: isAuthenticated ? new ChatComponent(this) : null,
                users: new UsersComponent(this)
            };
    
            if (isAuthenticated) {
                this.ws = new WebSocketService();
                document.querySelector('.main-nav').style.display = 'flex';
                document.getElementById('content').style.display = 'block';
                await this.showComponent('posts'); // Show posts on auth
            }
    
            await this.showComponent('posts');
        } catch (error) {
            console.error('Authentication check failed:', error);
            history.pushState(null, '', '/login');
            await this.showComponent('users', 'login');
        }

        console.log('App initialized successfully');
        return true;
    } catch (error) {
        console.error('Initialization failed:', error);
        return false;
    }


    async handleLogin(event) {
        event.preventDefault();
        
        const form = event.target;
        const formData = new FormData(form);
        const data = {
            identifier: formData.get('identifier'),
            password: formData.get('password')
        };

        try {
            const response = await this.api.post('login', data);
            if (response.status === 'success' && response.data?.token) {
                localStorage.setItem('auth_token', response.data.token);
                window.location.href = '/';
                return;
            }
            throw new Error(response.error || 'Login failed');
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    }

    async handleRoute() {
        const path = window.location.pathname;
        const route = this.routes[path] || this.routes['/'];
    
        if (route.requireAuth && !(await this.checkAuth())) {
            history.pushState(null, '', '/login');
            await this.showComponent('users', 'login');
            return;
        }
    
        // Hide all containers first
        ['content', 'auth-container', 'chat-container'].forEach(id => {
            const container = document.getElementById(id);
            if (container) container.style.display = 'none';
        });
    
        // Show navigation if authenticated
        const nav = document.querySelector('.main-nav');
        if (await this.checkAuth()) {
            nav.style.display = 'flex';
        } else {
            nav.style.display = 'none';
        }
    
        await route.handler();
    }

    async showComponent(name, mode = null) {
        const component = this.components[name];
        if (!component) {
            console.error('Component ${name}not found');
            return;
        }

        ['content', 'auth-container', 'chat-container'].forEach(id => {
            const container = document.getElementById(id);
            if (container) container.style.display = 'none';
        });

        // Show appropriate container based on component
        const containerId = name === 'users' ? 'auth-container' : 'content';
        const container = document.getElementById(containerId);
        if (container) {
            container.style.display = 'block';
        }

        await component.mount?.();
        if (mode && component.setMode) {
            component.setMode(mode);
        }
    }

    async checkAuth() {
        try {
            const token = localStorage.getItem('auth_token');
            console.log('Checking auth token:', token ? 'present' : 'missing');

            if (!token) {
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