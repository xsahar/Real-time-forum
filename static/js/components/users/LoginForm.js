// static/js/components/users/LoginForm.js
import { BaseForm } from './BaseForm.js';

export class LoginForm extends BaseForm {
    render() {
        return `
            <div class="auth-form">
                <h2>Login</h2>
                <div id="form-errors" class="error-message" style="display: none;"></div>
                <form id="login-form">
                    <div class="form-group">
                        <input type="text" 
                               name="identifier" 
                               placeholder="Email or Nickname" 
                               required
                               minlength="3"
                               autocomplete="username">
                    </div>
                    <div class="form-group">
                        <input type="password" 
                               name="password" 
                               placeholder="Password" 
                               required
                               minlength="6"
                               autocomplete="current-password">
                    </div>
                    <button type="submit" class="submit-btn">Login</button>
                    <p class="auth-switch">
                        No account? <a href="#" id="show-register">Register</a>
                    </p>
                </form>
            </div>
        `;
    }

    setupListeners() {
        const form = this.container.querySelector('#login-form');
        const submitBtn = form.querySelector('button[type="submit"]');
        const registerLink = this.container.querySelector('#show-register');

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            e.stopPropagation();

            this.clearError();
            this.setLoading(true);

            try {
                const formData = new FormData(form);
                const data = {
                    identifier: formData.get('identifier')?.trim(),
                    password: formData.get('password')
                };

                if (!this.validateFormData(data)) {
                    return;
                }

                const response = await this.onSubmit(data);
                console.log('Login response:', {
                    status: response?.status,
                    hasData: !!response?.data,
                    hasToken: !!response?.data?.token
                });
                
                if (response?.status === 'success' && response.data?.token) {
                    await this.handleSuccessfulLogin(response.data);
                    return;
                }

                throw new Error(response?.error || 'Login failed');
            } catch (error) {
                console.error('Login error:', error);
                this.showError(error.message);
                form.querySelector('input[name="password"]').value = '';
            } finally {
                this.setLoading(false);
            }
        });

        registerLink.addEventListener('click', (e) => {
            e.preventDefault();
            e.stopPropagation();
            this.onSwitchForm();
        });

        this.setupInputValidation(form);
    }

    setupInputValidation(form) {
        form.querySelectorAll('input').forEach(input => {
            input.addEventListener('input', () => {
                input.setCustomValidity('');
                if (input.checkValidity()) {
                    this.clearError();
                }
            });

            input.addEventListener('invalid', (e) => {
                e.preventDefault();
                const messages = {
                    identifier: 'Please enter a valid email or nickname (min 3 characters)',
                    password: 'Password must be at least 6 characters'
                };
                this.showError(messages[input.name] || 'Invalid input');
            });
        });
    }

    validateFormData(data) {
        if (!data.identifier || data.identifier.length < 3) {
            this.showError('Please enter a valid email or nickname (min 3 characters)');
            return false;
        }
        if (!data.password || data.password.length < 6) {
            this.showError('Password must be at least 6 characters');
            return false;
        }
        return true;
    }

    async handleSuccessfulLogin(data) {
        try {
            // Store token
            localStorage.setItem('auth_token', data.token);
            
            // Verify token storage
            const storedToken = localStorage.getItem('auth_token');
            if (!storedToken) {
                throw new Error('Failed to store authentication token');
            }

            // Notify app
            window.dispatchEvent(new CustomEvent('login-success', {
                detail: { 
                    token: storedToken,
                    user: data.user
                }
            }));

            // Navigate
            history.pushState(null, '', '/');
            window.dispatchEvent(new PopStateEvent('popstate'));
        } catch (error) {
            throw new Error('Failed to complete login: ' + error.message);
        }
    }

    setLoading(isLoading) {
        const submitBtn = this.container.querySelector('button[type="submit"]');
        if (submitBtn) {
            submitBtn.disabled = isLoading;
            submitBtn.textContent = isLoading ? 'Logging in...' : 'Login';
        }
    }

    showError(message) {
        super.showError(message);
        this.setLoading(false);
    }

    clearError() {
        super.clearError();
        this.setLoading(false);
    }
}