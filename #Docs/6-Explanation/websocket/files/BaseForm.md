### Class Structure
1. Constructor
```
constructor(container, onSubmit, onSwitchForm)
```
- container: DOM element to render form in
- onSubmit: Callback for form submission
- onSwitchForm: Callback to switch between forms
2. Abstract Methods (Must be implemented by child classes)
- render(): Returns form HTML
- setupListeners(): Sets up event handlers
- validateFormData(): Validates form input

3. Utility Methods
- mount(): Renders form and sets up listeners
- showError(): Displays error messages
- clearError(): Clears error messages
- resetForm(): Resets form state

### Usage Flow
1. Create form instance (Login/Register)
2. Mount form to container
3. Handle user interactions
4. Validate input
5. Submit or show errors
6. Switch between forms

### Benefits
- Code reuse
- Consistent error handling
- Standardized form lifecycle
- Enforced implementation
- Centralized DOM manipulation

### Error Handling Pattern
```
showError(message) {
    const errorDiv = this.container.querySelector('#form-errors');
    if (errorDiv) {
        errorDiv.textContent = message;
        errorDiv.style.display = 'block';
    }
}
```
This base class ensures both login and registration forms follow the same patterns and behaviors.