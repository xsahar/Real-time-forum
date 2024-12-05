export function addEventListeners(app) {
    // Global event delegation
    document.addEventListener('click', e => {
        handleButtonClicks(e, app);
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
        app.components.users?.logout();
    }

    // Navigation links
    if (e.target.matches('[data-link]')) {
        e.preventDefault();
        const href = e.target.getAttribute('href');
        history.pushState(null, '', href);
        app.handleRoute();
    }
}