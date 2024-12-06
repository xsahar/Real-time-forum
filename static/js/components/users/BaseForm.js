// static/js/components/users/BaseForm.js
export class BaseForm {
    constructor(container, onSubmit, onSwitchForm) {
        this.container = container;
        this.onSubmit = onSubmit;
        this.onSwitchForm = onSwitchForm;
    }

    mount() {
        this.container.innerHTML = this.render();
        this.setupListeners();
    }

    render() {
        throw new Error('Render method must be implemented by child class');
    }

    setupListeners() {
        throw new Error('SetupListeners method must be implemented by child class');
    }

    showError(message) {
        const errorDiv = this.container.querySelector('#form-errors');
        if (errorDiv) {
            errorDiv.textContent = message;
            errorDiv.style.display = 'block';
        }
    }

    clearError() {
        const errorDiv = this.container.querySelector('#form-errors');
        if (errorDiv) {
            errorDiv.textContent = '';
            errorDiv.style.display = 'none';
        }
    }

    validateFormData(data) {
        throw new Error('ValidateFormData method must be implemented by child class');
    }

    resetForm() {
        const form = this.container.querySelector('form');
        if (form) {
            form.reset();
            this.clearError();
        }
    }
}