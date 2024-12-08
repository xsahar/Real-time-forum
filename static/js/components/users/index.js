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
                history.pushState(null, '', '/register');
                this.showRegisterForm();
            } else if (e.target.id === 'show-login') {
                e.preventDefault();
                history.pushState(null, '', '/login');
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
                history.pushState(null, '', '/');
                await this.app.handleRoute();
                return;
            }
            this.showLoginForm();
        } catch (error) {
            console.error('Mount error:', error);
            this.showLoginForm();
        }
    }

    showLoginForm() {
        console.log('Showing login form');
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
            console.log('Processing login:',{
                identifier: formData.identifier,
                hasPassword: !!formData.password
            });

            const response = await this.api.post('login', formData);
            if (response.status === 'success' && response.data?.token) {
                localStorage.setItem('auth_token', response.data.token);
                this.navLinks.style.display = 'flex';
                history.pushState(null, '', '/');
                await this.app.handleRoute();
                return true;
            }
            throw new Error(response?.error || 'Login failed');
        } catch (error) {
            console.error('Login error:', error);
            throw error;
        }
    }

    async handleRegister(formData) {
        try {
            console.log('Processing registration');
            const response = await this.api.post('register', formData);
            if (response?.status === 'success') {
                history.pushState(null, '', '/login');
                this.showLoginForm();
                return true;
            }
            throw new Error(response.error || 'Registration failed');
        } catch (error) {
            console.error('Registration error:', error);
            throw new Error(error.message || 'Registration failed');
        }
    }

    async checkAuth() {
        try {
            const token = localStorage.getItem('auth_token');
            if (!token){
                 return false;
            }
            
            const response = await this.api.get('auth/status');
        return response?.status === 'success';
    } catch (error) {
        console.error('Auth check error:', error);
        localStorage.removeItem('auth_token');
        return false;
    }
    }
}