export const initAuth = async () => {
    const loginForm = document.getElementById('login-form');
    const registerForm = document.getElementById('register-form');
    const showRegisterLink = document.getElementById('show-register');
    const showLoginLink = document.getElementById('show-login');
    const authSection = document.getElementById('auth-section');
    const mainSection = document.getElementById('main-section');
    const logoutBtn = document.getElementById('logout-btn');

    // Validate required elements
    if (!authSection || !mainSection) {
        console.error('Required auth sections not found');
        return { isAuthenticated: false };
    }

    // Form toggle handlers
    if (showRegisterLink) {
        showRegisterLink.addEventListener('click', (e) => {
            e.preventDefault();
            loginForm?.classList.add('hidden');
            registerForm?.classList.remove('hidden');
        });
    }

    if (showLoginLink) {
        showLoginLink.addEventListener('click', (e) => {
            e.preventDefault();
            registerForm?.classList.add('hidden');
            loginForm?.classList.remove('hidden');
        });
    }

    // Register form handler
    if (registerForm) {
        registerForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const formData = {
                nickname: document.getElementById('reg-nickname')?.value.trim(),
                age: parseInt(document.getElementById('reg-age')?.value, 10),
                gender: document.getElementById('reg-gender')?.value,
                firstName: document.getElementById('reg-firstname')?.value.trim(),
                lastName: document.getElementById('reg-lastname')?.value.trim(),
                email: document.getElementById('reg-email')?.value.trim(),
                password: document.getElementById('reg-password')?.value
            };

            // Validation
            if (!formData.nickname || !formData.email || !formData.password) {
                alert('Please fill in all required fields');
                return;
            }

            if (formData.password.length < 6) {
                alert('Password must be at least 6 characters long');
                return;
            }

            if (isNaN(formData.age) || formData.age < 13) {
                alert('Please enter a valid age (13 or older)');
                return;
            }

            try {
                const response = await fetch('/api/auth/register', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(formData)
                });

                const data = await response.json();

                if (!response.ok) {
                    throw new Error(data.message || 'Registration failed');
                }

                console.log('Registration successful');
                registerForm.reset();
                registerForm.classList.add('hidden');
                loginForm?.classList.remove('hidden');
                alert('Registration successful! Please login.');
            } catch (error) {
                console.error('Registration error:', error);
                alert(error.message || 'Registration failed. Please try again.');
            }
        });
    }

    // Login form handler
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();

            const formData = {
                identifier: document.getElementById('login-identifier')?.value.trim(),
                password: document.getElementById('login-password')?.value
            };

            if (!formData.identifier || !formData.password) {
                alert('Please fill in all fields');
                return;
            }

            try {
                const response = await fetch('/api/auth/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(formData)
                });

                const data = await response.json();

                if (!response.ok || !data.token || !data.user) {
                    throw new Error(data.message || 'Login failed');
                }

                // Store auth data
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));

                // Let app.js handle UI updates by returning auth state
                location.reload(); // Force page reload to reinitialize
                return {
                    isAuthenticated: true,
                    user: data.user,
                    token: data.token
                };

            } catch (error) {
                console.error('Login error:', error);
                alert(error.message || 'Login failed. Please try again.');
                return { isAuthenticated: false };
            }
        });
    }

    // Logout handler
    if (logoutBtn) {
        logoutBtn.addEventListener('click', () => {
            try {
                // Clear auth data
                localStorage.removeItem('token');
                localStorage.removeItem('user');

                // Reset forms
                loginForm?.reset();
                registerForm?.reset();

                // Show login form
                loginForm?.classList.remove('hidden');
                registerForm?.classList.add('hidden');

                location.reload(); // Force page reload
                console.log('Logout successful');
                return { isAuthenticated: false };
            } catch (error) {
                console.error('Logout error:', error);
                alert('Error during logout. Please refresh the page.');
                return { isAuthenticated: false };
            }
        });
    }

    // Check initial auth state
    try {
        const token = localStorage.getItem('token');
        const user = JSON.parse(localStorage.getItem('user') || 'null');

        if (token && user?.id) {
            return {
                isAuthenticated: true,
                user,
                token
            };
        }
        return { isAuthenticated: false };
    } catch (error) {
        console.error('Auth initialization error:', error);
        return { isAuthenticated: false };
    }
};