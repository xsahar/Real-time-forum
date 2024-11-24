export class UsersComponent {
    constructor(app) {
        this.app = app;
        this.api = app.api;
        this.container = document.getElementById('auth-container');
        this.navLinks = document.querySelector('.main-nav');
    }

    async mount() {
        // Hide navigation initially
        this.navLinks.style.display = 'none';
        const isAuth = await this.checkAuth();

        if (isAuth) {
            this.navLinks.style.display = 'flex';
            window.location.href = '/';
            return;
        }
        this.renderLoginForm();
        this.setupEventListeners();
    }

    async checkAuth() {
        try {
            const response = await this.api.get('auth/status');
            return response.status === 'success';
        } catch {
            localStorage.removeItem('auth_token');
            return false;
        }
    }

    renderLoginForm() {
        this.container.innerHTML = `
            <div class="auth-form">
                <h2>Login</h2>
                <form id="login-form">
                    <div class="form-group">
                        <input type="text" name="identifier" placeholder="Email or Nickname" required>
                    </div>
                    <div class="form-group">
                        <input type="password" name="password" placeholder="Password" required>
                    </div>
                    <button type="submit" class="submit-btn">Login</button>
                    <p class="auth-switch">
                        No account? <a href="#" id="show-register">Register</a>
                    </p>
                </form>
            </div>
        `;
    }

    renderRegisterForm() {
        this.container.innerHTML = `
            <div class="auth-form">
                <h2>Create Account</h2>
                <form id="register-form">
                    <div class="form-group">
                        <input type="text" name="nickname" placeholder="Nickname" required 
                               pattern="[A-Za-z0-9_]{3,}" title="3+ characters, letters, numbers, underscore">
                    </div>
                    <div class="form-group">
                        <input type="number" name="age" placeholder="Age" required min="13" max="120">
                    </div>
                    <div class="form-group">
                        <select name="gender" required>
                            <option value="">Select Gender</option>
                            <option value="male">Male</option>
                            <option value="female">Female</option>
                            <option value="other">Other</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <input type="text" name="firstName" placeholder="First Name" required
                               pattern="[A-Za-z]{2,}" title="2+ letters">
                    </div>
                    <div class="form-group">
                        <input type="text" name="lastName" placeholder="Last Name" required
                               pattern="[A-Za-z]{2,}" title="2+ letters">
                    </div>
                    <div class="form-group">
                        <input type="email" name="email" placeholder="Email" required>
                    </div>
                    <div class="form-group">
                        <input type="password" name="password" placeholder="Password" required
                               minlength="6" title="6+ characters">
                    </div>
                    <button type="submit" class="submit-btn">Create Account</button>
                    <p class="auth-switch">
                        Have an account? <a href="#" id="show-login">Login</a>
                    </p>
                </form>
            </div>
        `;
    }

    setupEventListeners() {
        // Form toggle handlers
        this.container.addEventListener('click', (e) => {
            e.preventDefault(); // Prevent default for all clicks in container
            const target = e.target;

            if (target.id === 'show-register') {
                this.renderRegisterForm();
            } else if (target.id === 'show-login') {
                this.renderLoginForm();
            }
        });

        // Form submission handlers
        this.container.addEventListener('submit', async (e) => {
            e.preventDefault();
            const form = e.target;
            const submitBtn = form.querySelector('.submit-btn');
            
            // Disable button and show loading state
            submitBtn.disabled = true;
            submitBtn.textContent = 'Processing...';
            
            try {
                const formData = new FormData(form);
                const data = Object.fromEntries(formData);

                if (form.id === 'login-form') {
                    console.log('Attempting login with:', data);
                    await this.login(data);
                    this.showSuccess('Login successful!');
                    // Add small delay before redirect
                    setTimeout(() => {
                        this.navLinks.style.display = 'flex';
                        window.location.href = '/';
                    }, 500);
                } else if (form.id === 'register-form') {
                    console.log('Attempting registration with:', data);
                    await this.register(data);
                    this.showSuccess('Registration successful! Please login.');
                    setTimeout(() => this.renderLoginForm(), 1500);
                }
            } catch (error) {
                console.error('Form submission error:', error);
                this.showError(form, error.message);
                // Reset button text on error
                submitBtn.textContent = form.id === 'login-form' ? 'Login' : 'Create Account';
            } finally {
                submitBtn.disabled = false;
            }
        });
    }

    showError(form, message) {
        this.removeMessages();
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error-message';
        errorDiv.textContent = message;
        form.insertBefore(errorDiv, form.firstChild);
    }

    showSuccess(message) {
        this.removeMessages();
        const successDiv = document.createElement('div');
        successDiv.className = 'success-message';
        successDiv.textContent = message;
        this.container.insertBefore(successDiv, this.container.firstChild);
    }

    removeMessages() {
        this.container.querySelectorAll('.error-message, .success-message')
            .forEach(el => el.remove());
    }

    async login(credentials) {
        try {
            const response = await this.api.post('login', credentials);
            console.log('Login response:', response);
            
            if (response.status === 'success' && response.token) {
                localStorage.setItem('auth_token', response.token);
                return true;
            } else {
                throw new Error(response.error || 'Login failed. Please check your credentials.');
            }
        } catch (error) {
            console.error('Login error:', error);
            throw new Error('Login failed. Please try again.');
        }
    }

    async register(userData) {
        try {
            // Validate age
            if (parseInt(userData.age) < 13) {
                throw new Error('Must be at least 13 years old');
            }
            // Validate password
            if (userData.password.length < 6) {
                throw new Error('Password must be at least 6 characters');
            }
            // Additional validations
            if (!userData.nickname || userData.nickname.length < 3) {
                throw new Error('Nickname must be at least 3 characters');
            }
            if (!userData.email || !userData.email.includes('@')) {
                throw new Error('Please enter a valid email address');
            }

            console.log('Sending registration request:', userData);
            const response = await this.api.post('register', userData);
            
            if (response.status !== 'success') {
                throw new Error(response.error || 'Registration failed');
            }
            
            return true;
        } catch (error) {
            console.error('Registration error:', error);
            throw new Error(error.message || 'Registration failed. Please try again.');
        }
    }

    logout() {
        localStorage.removeItem('auth_token');
        window.location.href = '/login';
    }
}