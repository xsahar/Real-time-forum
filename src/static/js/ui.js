export const initUI = () => {
    const themeToggle = document.getElementById('theme-toggle');
    if (themeToggle) {
        const currentTheme = localStorage.getItem('theme') || 'light';
        document.body.classList.toggle('dark-theme', currentTheme === 'dark');
        
        themeToggle.addEventListener('click', () => {
            document.body.classList.toggle('dark-theme');
            localStorage.setItem('theme', 
                document.body.classList.contains('dark-theme') ? 'dark' : 'light'
            );
        });
    }

    // Navigation handling
    const navHome = document.getElementById('nav-home');
    const postsFeed = document.getElementById('posts-feed');
    const createPostSection = document.getElementById('create-post-section');

    if (navHome && postsFeed && createPostSection) {
        navHome.addEventListener('click', (e) => {
            e.preventDefault();
            postsFeed.classList.remove('hidden');
            createPostSection.classList.add('hidden');
        });
    }
};