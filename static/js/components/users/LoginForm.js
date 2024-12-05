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

    setupListeners() {
        const form = this.container.querySelector('#login-form');
        const submitBtn = form.querySelector('button[type="submit"]');
        const registerLink = this.container.querySelector('#show-register');

        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            this.clearError();
            submitBtn.disabled = true;

            try {
                const formData = new FormData(form);
                const data = {
                    identifier: formData.get('identifier')?.trim(),
                    password: formData.get('password')
                };

                if (!data.identifier || !data.password) {
                    throw new Error('All fields are required');
                }

                await this.onSubmit(data);
            } catch (error) {
                this.showError(error.message);
            } finally {
                submitBtn.disabled = false;
            }
        });

        registerLink.addEventListener('click', (e) => {
            e.preventDefault();
            this.onSwitchForm();
        });
    }
}