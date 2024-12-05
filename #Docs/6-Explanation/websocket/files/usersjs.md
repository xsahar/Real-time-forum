### File Structure & Purpose
1. BaseForm.js
Base class for form handling
Provides common functionality: mounting, error display, clearing errors
Used by both LoginForm and RegisterForm
2. LoginForm.js
Handles login form rendering and interactions
Contains form HTML template with email/nickname and password fields
Has "Register" link that triggers form switch
3. RegisterForm.js
Handles registration form rendering and validation
Contains form HTML template with all user fields
Validates form data before submission
Has "Login" link to switch back
4. index.js (UsersComponent)
Main controller for user authentication
Manages switching between login and register forms
Handles API calls for login/register
Manages auth state and navigation
### Registration Link Flow
1. When "Register" is clicked in LoginForm:
```
Login Form Click -> onSwitchForm() -> UsersComponent.showRegisterForm() -> Mount Register Form
```
2. The flow works because:
LoginForm has click listener on register link
Calls this.onSwitchForm() (passed from UsersComponent)
UsersComponent's showRegisterForm() creates new RegisterForm
RegisterForm mounts itself in container

Key Interactions
1. Form Switching
```
Login <-> Register switching handled by UsersComponent
Both forms get onSubmit and onSwitchForm callbacks
BaseForm provides common mounting logic
```
2. Form Validation
```
RegisterForm validates all fields
LoginForm checks required fields
Both use BaseForm's error display
```

3. API Integration
```
Forms collect data
UsersComponent makes API calls
Responses trigger navigation/error display
```
