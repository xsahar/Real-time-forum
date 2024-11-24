export function handleNavigation(app) {
    // Handle clicks on navigation links
    document.addEventListener('click', e => {
        const link = e.target.closest('[data-link]');
        if (link) {
            e.preventDefault();
            navigateTo(link.href, app);
        }
    });

    // Handle browser back/forward
    window.addEventListener('popstate', () => {
        app.handleRoute();
    });
}

export function navigateTo(url, app) {
    history.pushState(null, null, url);
    app.handleRoute();
}