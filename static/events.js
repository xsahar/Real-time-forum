document.addEventListener('DOMContentLoaded', () => {
    document.getElementById('main-content').innerHTML = templates.home;

    document.getElementById('home-link').addEventListener('click', async (e) => {
        e.preventDefault();
        document.getElementById('main-content').innerHTML = templates.home;
        addFilterFormListener();
    });

    document.getElementById('register-link').addEventListener('click', async (e) => {
        e.preventDefault();
        document.getElementById('main-content').innerHTML = templates.register;
    });

    document.getElementById('login-link').addEventListener('click', async (e) => {
        e.preventDefault();
        document.getElementById('main-content').innerHTML = templates.login;
    });

    document.getElementById('new-post-link').addEventListener('click', async (e) => {
        e.preventDefault();
        document.getElementById('main-content').innerHTML = templates.newPost;
    });

    addFilterFormListener();
});

function addFilterFormListener() {
    const filterForm = document.getElementById('filter-form');
    if (filterForm) {
        filterForm.addEventListener('submit', function(event) {
            event.preventDefault(); // Prevent form submission and page reload

            // Get selected filter criteria from form inputs
            const category = document.getElementById('category-select').value;
            const postType = document.querySelector('input[name="filter"]:checked').value;

            // Perform filtering logic based on category and post type
            // Update post display dynamically without reloading the page
            filterPosts(category, postType);
        });
    }
}

function filterPosts(category, postType) {
    // Perform filtering logic here based on category and post type
    // Update the display of posts in the posts-container section
    
    // Example: Filter posts based on selected criteria
    const allPosts = document.querySelectorAll('.post');
    allPosts.forEach(post => {
        if ((category === '' || post.dataset.category === category) && (postType === 'all' || post.dataset.type === postType)) {
            post.style.display = 'block'; // Show post
        } else {
            post.style.display = 'none'; // Hide post
        }
    });
}