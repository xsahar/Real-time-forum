// static/js/components/users/index.js
import { LoginForm } from './LoginForm.js';
import { RegisterForm } from './RegisterForm.js';

export class UsersComponent {
    constructor(app) {
        this.app = app;
        this.api = app.api;
        this.container = document.getElementById('auth-container');
        this.navLinks = document.querySelector('.main-nav');
        this.setupEventListeners();
    }

    setupEventListeners() {
        // Global event delegation for form switching
        this.container.addEventListener('click', (e) => {
            if (e.target.id === 'show-register') {
                e.preventDefault();
                console.log('Switching to register form');
                this.showRegisterForm();
            } else if (e.target.id === 'show-login') {
                e.preventDefault();
                console.log('Switching to login form');
                this.showLoginForm();
            }
        });
    }

    async mount() {
        try {
            this.navLinks.style.display = 'none';
            const isAuth = await this.checkAuth();

            if (isAuth) {
                this.navLinks.style.display = 'flex';
                window.location.href = '/';
                return;
            }
            this.showLoginForm();
        } catch (error) {
            console.error('Mount error:', error);
        }
    }

    showLoginForm() {
        const loginForm = new LoginForm(
            this.container, 
            this.handleLogin.bind(this),
            this.showRegisterForm.bind(this)
        );
        loginForm.mount();
    }

    showRegisterForm() {
        console.log('Showing register form');
        const registerForm = new RegisterForm(
            this.container,
            this.handleRegister.bind(this),
            this.showLoginForm.bind(this)
        );
        registerForm.mount();
    }

    async handleLogin(formData) {
        try {
            const response = await this.api.post('login', formData);
            if (response.status === 'success' && response.token) {
                localStorage.setItem('auth_token', response.token);
                this.navLinks.style.display = 'flex';
                window.location.href = '/';
                return true;
            }
            throw new Error(response.error || 'Login failed');
        } catch (error) {
            throw new Error(error.message || 'Login failed');
        }
    }

    async handleRegister(formData) {
        try {
            const response = await this.api.post('register', formData);
            if (response.status === 'success') {
                this.showLoginForm();
                return true;
            }
            throw new Error(response.error || 'Registration failed');
        } catch (error) {
            throw new Error(error.message || 'Registration failed');
        }
    }

    async checkAuth() {
        try {
            const token = localStorage.getItem('auth_token');
            if (!token) return false;
            
            const response = await this.api.get('auth/status');
            return response.status === 'success';
        } catch (error) {
            localStorage.removeItem('auth_token');
            return false;
        }
    }
}