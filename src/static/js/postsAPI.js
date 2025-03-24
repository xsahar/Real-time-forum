import { debounce } from './uiUtils.js';

export const loadCategories = async () => {
    const categorySelect = document.getElementById('create-post-category');
    if (!categorySelect) {
        console.error('Category select not found');
        return;
    }
    if (categorySelect.options.length > 1) return;

    try {
        const response = await fetch('/api/categories');
        if (!response.ok) throw new Error('Failed to load categories');
        const categories = await response.json();
        categorySelect.innerHTML = '<option value="">Select Category</option>';
        categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category;
            option.textContent = category;
            categorySelect.appendChild(option);
        });
    } catch (error) {
        console.error(error);
        categorySelect.innerHTML = '<option>Error loading categories</option>';
    }
};

export const loadPosts = async (category = '') => {
    const postsFeed = document.getElementById('posts-feed');
    if (!postsFeed) return;

    postsFeed.innerHTML = '<div class="loading">Loading posts...</div>';
    try {
        const token = localStorage.getItem('token');
        if (!token) throw new Error('No token');
        const response = await fetch(`/api/posts${category ? `?category=${category}` : ''}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) throw new Error('Failed to load posts');
        const posts = await response.json();
        postsFeed.innerHTML = posts.length ? posts.map(post => `
            <article class="post-card" data-post-id="${post.id}">
                <div class="post-header">
                    <h3>${post.title}</h3>
                    <span>${post.category}</span>
                </div>
                <div class="post-meta">
                    Posted by ${post.author?.nickname || 'Unknown'} on ${new Date(post.createdAt).toLocaleDateString()}
                </div>
                <div class="post-actions">
                    <button class="view-comments-btn" data-post-id="${post.id}">View Comments</button>
                </div>
            </article>
        `).join('') : '<p>No posts yet</p>';
    } catch (error) {
        console.error(error);
        postsFeed.innerHTML = '<div class="error">No posts yet</div>';
    }
};

export async function handlePostSubmission(e) {
    e.preventDefault();
    const submitBtn = e.target.querySelector('button[type="submit"]');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Creating...';

    try {
        const user = JSON.parse(localStorage.getItem('user'));
        const token = localStorage.getItem('token');
        if (!user || !token) throw new Error('Authentication required');

        const formData = {
            title: document.getElementById('create-post-title').value.trim(),
            content: document.getElementById('create-post-content').value.trim(),
            category: document.getElementById('create-post-category').value,
            userId: user.id
        };
        if (!formData.title || !formData.content || !formData.category) {
            throw new Error('Fill in all fields');
        }

        const timeoutPromise = new Promise((_, reject) =>
            setTimeout(() => reject(new Error('Timeout')), 5000)
        );
        const fetchPromise = fetch('/api/posts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(formData)
        });

        const response = await Promise.race([fetchPromise, timeoutPromise]);
        if (!response.ok) throw new Error('Failed to create post');
        await response.json();

        e.target.reset();
        document.getElementById('create-post-section').classList.add('hidden');
        document.getElementById('posts-feed').classList.remove('hidden');
        await debounce(loadPosts, 500)();

    } catch (error) {
        console.error(error);
        alert(error.message);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Create Post';
    }
}
