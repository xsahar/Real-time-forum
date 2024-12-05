
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
}