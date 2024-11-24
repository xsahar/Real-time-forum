export function addEventListeners(app) {
    // Global event delegation
    document.addEventListener('click', e => {
        handleButtonClicks(e, app);
        handleFormSubmissions(e, app);
    });

    // Add navigation handling
    window.addEventListener('popstate', () => {
        app.handleRoute();
    });
}

function handleButtonClicks(e, app) {
    // Like button
    if (e.target.matches('.like-btn')) {
        const postId = e.target.closest('[data-post-id]').dataset.postId;
        app.components.posts?.handleLike(postId);
    }
    
    // Logout button
    if (e.target.matches('#logout-btn')) {
        e.preventDefault();
        app.handleLogout();
    }

    // Navigation links
    if (e.target.matches('[data-link]')) {
        e.preventDefault();
        const href = e.target.getAttribute('href');
        history.pushState(null, '', href);
        app.handleRoute();
    }
}

function handleFormSubmissions(e, app) {
    if (e.target.matches('form')) {
        e.preventDefault();
        const formData = new FormData(e.target);
        const formId = e.target.id;
        
        switch(formId) {
            case 'login-form':
                app.components.users?.handleLogin(Object.fromEntries(formData));
                break;
            case 'register-form':
                app.components.users?.handleRegister(Object.fromEntries(formData));
                break;
            case 'post-form':
                app.components.posts?.handleNewPost(Object.fromEntries(formData));
                break;
            case 'chat-form':
                app.components.chat?.handleNewMessage(Object.fromEntries(formData));
                break;
        }
    }
}