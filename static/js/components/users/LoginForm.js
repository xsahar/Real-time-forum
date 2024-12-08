// static/js/components/users/LoginForm.js
import { BaseForm } from './BaseForm.js';

export class LoginForm extends BaseForm {
    constructor(container, onSubmit, onSwitchForm) {
        super(container, onSubmit, onSwitchForm);
        this.isLoading = false;
    }

    render() {
        return `
            <div class="auth-form">
                <h2>Login</h2>
                <div id="form-errors" class="error-message" style="display: none;"></div>
                <form id="login-form" novalidate>
                    <div class="form-group">
                        <input type="text" 
                               name="identifier" 
                               placeholder="Email or Nickname" 
                               required
                               minlength="3"
                               maxlength="50"
                               autocomplete="username"
                               class="form-input">
                    </div>
                    <div class="form-group">
                        <input type="password" 
                               name="password" 
                               placeholder="Password" 
                               required
                               minlength="6"
                               maxlength="100"
                               autocomplete="current-password"
                               class="form-input">
                    </div>
                    <button type="submit" class="submit-btn" data-testid="login-submit">Login</button>
                    <p class="auth-switch">
                        No account? <a href="#" id="show-register">Register</a>
                    </p>
                </form>
            </div>
        `;
    }

    setupListeners() {
        const form = this.container.querySelector('#login-form');
        const registerLink = this.container.querySelector('#show-register');

        form.addEventListener('submit', this.handleSubmit.bind(this));
        registerLink?.addEventListener('click', this.handleRegisterClick.bind(this));
        
        // Setup input validation
        form.querySelectorAll('input').forEach(input => {
            input.addEventListener('input', () => this.validateOnInput(input));
            input.addEventListener('blur', () => this.validateOnBlur(input));
        });
    }


    sanitizeFormData(data) {
        return {
            identifier: data.identifier?.trim() || '',
            password: data.password || ''
        };
    }

    async handleSubmit(e) {
        e.preventDefault();
        e.stopPropagation();
    
        if (this.isLoading) return;
    
        const form = e.target;
        this.setLoading(true);
        this.clearError();
    
        try {
            // Debug form fields directly
            const identifierField = form.querySelector('input[name="identifier"]');
            const passwordField = form.querySelector('input[name="password"]');
    
            console.log('Form fields:', {
                identifierValue: identifierField?.value,
                identifierExists: !!identifierField,
                passwordExists: !!passwordField,
                passwordLength: passwordField?.value?.length
            });
    
            // Get form data
            const data = {
                identifier: identifierField?.value?.trim(),
                password: passwordField?.value
            };
    
            // Debug data before validation
            console.log('Form data before validation:', {
                hasIdentifier: !!data.identifier,
                identifierLength: data.identifier?.length,
                hasPassword: !!data.password,
                passwordLength: data.password?.length
            });
    
            // Validate data
            if (!this.validateFormData(data)) {
                return;
            }
    
            const response = await this.onSubmit(data);
            
            if (response?.status === 'success' && response.data?.token) {
                await this.handleSuccessfulLogin(response.data);
                return;
            }
    
            throw new Error(response?.error || 'Login failed');
        } catch (error) {
            console.error('Login error:', error);
            this.showError(error.message);
        } finally {
            this.setLoading(false);
        }
    }
    
    handleRegisterClick(e) {
        e.preventDefault();
        if (this.isLoading) return;
        this.onSwitchForm();
    }

    validateOnInput(input) {
        input.classList.remove('invalid');
        this.clearError();
    }

    validateOnBlur(input) {
        const value = input.value.trim();
        const isValid = value.length >= input.minLength && value.length <= input.maxLength;
        
        if (!isValid) {
            input.classList.add('invalid');
            this.showError(this.getValidationMessage(input.name));
        }
    }

    getValidationMessage(fieldName) {
        const messages = {
            identifier: 'Please enter a valid email or nickname (3-50 characters)',
            password: 'Password must be between 6-100 characters'
        };
        return messages[fieldName] || 'Invalid input';
    }

    validateFormData(data) {
        if (!data.identifier || data.identifier.length < 3 || data.identifier.length > 50) {
            this.showError(this.getValidationMessage('identifier'));
            return false;
        }
        if (!data.password || data.password.length < 6 || data.password.length > 100) {
            this.showError(this.getValidationMessage('password'));
            return false;
        }
        return true;
    }

    async handleSuccessfulLogin(data) {
        try {
            if (!data?.token) {
                throw new Error('No authentication token received');
            }

            localStorage.setItem('auth_token', data.token);
            
            window.dispatchEvent(new CustomEvent('login-success', {
                detail: { token: data.token, user: data.user }
            }));

            window.location.replace('/');
        } catch (error) {
            console.error('Login completion error:', error);
            throw new Error(`Failed to complete login: ${error.message}`);
        }
    }

    setLoading(isLoading) {
        this.isLoading = isLoading;
        const submitBtn = this.container.querySelector('.submit-btn');
        const inputs = this.container.querySelectorAll('input');

        if (submitBtn) {
            submitBtn.disabled = isLoading;
            submitBtn.textContent = isLoading ? 'Logging in...' : 'Login';
        }

        inputs.forEach(input => input.disabled = isLoading);
    }
}