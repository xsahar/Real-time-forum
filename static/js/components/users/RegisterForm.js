// static/js/components/users/RegisterForm.js
import { BaseForm } from './BaseForm.js';

export class RegisterForm extends BaseForm {
    render() {
        return `
            <div class="auth-form">
                <h2>Create Account</h2>
                <div id="form-errors" class="error-message" style="display: none;"></div>
                <form id="register-form">
                    <div class="form-group">
                        <input type="text" name="nickname" placeholder="Nickname" 
                               required pattern="[A-Za-z0-9_]{3,}"
                               title="3+ characters, letters, numbers, underscore">
                    </div>
                    <div class="form-group">
                        <input type="number" name="age" placeholder="Age" 
                               required min="13" max="120">
                    </div>
                    <div class="form-group">
                        <select name="gender" required>
                            <option value="">Select Gender</option>
                            <option value="male">Male</option>
                            <option value="female">Female</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <input type="text" name="firstName" placeholder="First Name" 
                               required pattern="[A-Za-z]{2,}">
                    </div>
                    <div class="form-group">
                        <input type="text" name="lastName" placeholder="Last Name" 
                               required pattern="[A-Za-z]{2,}">
                    </div>
                    <div class="form-group">
                        <input type="email" name="email" placeholder="Email" required>
                    </div>
                    <div class="form-group">
                        <input type="password" name="password" placeholder="Password" 
                               required minlength="6">
                    </div>
                    <button type="submit" class="submit-btn">Create Account</button>
                    <p class="auth-switch">
                        Have an account? <a href="#" id="show-login">Login</a>
                    </p>
                </form>
            </div>
        `;
    }

        setupListeners() {
            const form = this.container.querySelector('#register-form');
            const submitBtn = form.querySelector('button[type="submit"]');
    
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                this.clearError();
                submitBtn.disabled = true;
    
                try {
                    const formData = new FormData(form);
                    // Log raw form data
                    console.log('Raw form data:', Object.fromEntries(formData));
    
                    const data = {
                        nickname: formData.get('nickname')?.trim(),
                        email: formData.get('email')?.trim(),
                        password: formData.get('password'),
                        firstName: formData.get('firstName')?.trim(),
                        lastName: formData.get('lastName')?.trim(),
                        gender: formData.get('gender'),
                        age: parseInt(formData.get('age'))
                    };
    
                    // Validate required fields
                    const requiredFields = ['nickname', 'email', 'password', 'firstName', 'lastName', 'gender', 'age'];
                    const missingFields = requiredFields.filter(field => !data[field]);
    
                    if (missingFields.length > 0) {
                        throw new Error(`Missing fields: ${missingFields.join(', ')}`);
                    }
    
                    // Format validation
                    if (data.password.length < 6) throw new Error('Password must be at least 6 characters');
                    if (!data.email.includes('@')) throw new Error('Invalid email format');
                    if (data.nickname.length < 3) throw new Error('Nickname must be at least 3 characters');
                    if (data.age < 13 || data.age > 120) throw new Error('Age must be between 13 and 120');
    
                    // Log formatted data
                    console.log('Formatted data:', {...data, password: '[REDACTED]'});
                    
                    await this.onSubmit(data);
                } catch (error) {
                    console.error('Registration error:', error);
                    this.showError(error.message);
                } finally {
                    submitBtn.disabled = false;
                }
            });
    
            // Add form validation on input
            form.querySelectorAll('input, select').forEach(input => {
                input.addEventListener('input', () => {
                    input.setCustomValidity('');
                    if (input.checkValidity()) {
                        this.clearError();
                    }
                });
            });
        }
    
    validateFormData(data) {
        const validations = [
            { test: () => data.nickname?.length >= 3, message: 'Nickname must be at least 3 characters' },
            { test: () => data.password?.length >= 6, message: 'Password must be at least 6 characters' },
            { test: () => data.email?.includes('@'), message: 'Invalid email format' },
            { test: () => data.age >= 13 && data.age <= 120, message: 'Age must be between 13 and 120' },
            { test: () => data.firstName?.trim().length >= 2, message: 'First name is required' },
            { test: () => data.lastName?.trim().length >= 2, message: 'Last name is required' },
            { test: () => data.gender, message: 'Gender is required' }
        ];

        for (const validation of validations) {
            if (!validation.test()) {
                this.showError(validation.message);
                return false;
            }
        }
        return true;
    }
}