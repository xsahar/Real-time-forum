import { createComment } from './comments.js';

// Module-level variables
let categoriesLoaded = false;
let postsFeed = null; // Add this line

export const initPosts = async () => {
    const createPostForm = document.getElementById('create-post-form');
    const createPostBtn = document.getElementById('create-post-btn');
    const cancelPostBtn = document.getElementById('cancel-post');
    const createPostSection = document.getElementById('create-post-section');
    postsFeed = document.getElementById('posts-feed'); // Update this line

    // Validate required elements
    if (!createPostForm || !createPostSection || !postsFeed) {
        console.error('Required elements not found:', {
            form: !!createPostForm,
            section: !!createPostSection,
            feed: !!postsFeed
        });
        return;
    }

    // Load categories only once at module level
    if (!categoriesLoaded) {
        await loadCategories();
        categoriesLoaded = true;
    }

    // Load posts
    await loadPosts();

    // Create post button handler
    if (createPostBtn) {
        createPostBtn.addEventListener('click', (e) => {
            e.preventDefault();
            console.log('Opening create post form');
            createPostSection.classList.remove('hidden');
            postsFeed.classList.add('hidden');
        });
    }

    // Cancel button handler
    if (cancelPostBtn) {
        cancelPostBtn.addEventListener('click', () => {
            console.log('Canceling post creation');
            createPostSection.classList.add('hidden');
            postsFeed.classList.remove('hidden');
            createPostForm.reset();
        });
    }

    // Form submission handler - Update this part
    createPostForm.addEventListener('submit', handlePostSubmission);
};

// New separate function for post submission
async function handlePostSubmission(e) {
    e.preventDefault();
    console.log('Starting post creation...');

    const submitBtn = e.target.querySelector('button[type="submit"]');
    submitBtn.disabled = true;
    submitBtn.textContent = 'Creating...';

    try {
        const user = JSON.parse(localStorage.getItem('user'));
        const token = localStorage.getItem('token');

        if (!user || !token) {
            throw new Error('Authentication required');
        }

        const formData = {
            title: document.getElementById('create-post-title').value.trim(),
            content: document.getElementById('create-post-content').value.trim(),
            category: document.getElementById('create-post-category').value,
            userId: user.id
        };

        // Validate form data
        if (!formData.title || !formData.content || !formData.category) {
            throw new Error('Please fill in all fields');
        }

        console.log('Sending form data:', formData);

        // Use Promise.race to prevent duplicate submissions
        const timeoutPromise = new Promise((_, reject) =>
            setTimeout(() => reject(new Error('Request timeout')), 5000)
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

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Failed to create post');
        }

        const data = await response.json();
        if (!data.id || !data.title) {
            console.error('Invalid response data:', data);
            throw new Error('Invalid server response');
        }

        // Success handling
        console.log('Post created successfully:', data);
        e.target.reset();
        document.getElementById('create-post-section').classList.add('hidden');
        document.getElementById('posts-feed').classList.remove('hidden');

        // Debounce loadPosts call
        await debounce(loadPosts, 500)();

    } catch (error) {
        console.error('Post creation error:', error);
        alert(error.message);
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = 'Create Post';
    }
}

// Add debounce utility
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        return new Promise((resolve) => {
            const later = () => {
                clearTimeout(timeout);
                resolve(func(...args));
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        });
    };
}

// Update loadPosts to use debouncing
const debouncedLoadPosts = debounce(loadPosts, 500);

export const loadCategories = async () => {
    const categorySelect = document.getElementById('create-post-category');
    if (!categorySelect) {
        console.error('Category select element not found');
        return;
    }

    try {
        // Check if categories are already loaded
        if (categorySelect.options.length > 1) {
            console.log('Categories already loaded');
            return;
        }

        console.log('Loading categories...');
        const response = await fetch('/api/categories');

        if (!response.ok) {
            throw new Error('Failed to load categories');
        }

        const categories = await response.json();
        console.log('Categories loaded:', categories);

        if (!Array.isArray(categories)) {
            throw new Error('Invalid categories data');
        }

        // Clear existing options except the first one
        while (categorySelect.options.length > 1) {
            categorySelect.remove(1);
        }

        // Add new options
        categories.forEach(category => {
            const option = document.createElement('option');
            option.value = category;
            option.textContent = category;
            categorySelect.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading categories:', error);
        categorySelect.innerHTML = '<option value="">Error loading categories</option>';
    }
};


async function loadPosts(category = '') {
    if (!postsFeed) {
        postsFeed = document.getElementById('posts-feed');
    }
    if (!postsFeed) {
        console.error('Posts feed element not found');
        return;
    }

    postsFeed.innerHTML = '<div class="loading">Loading posts...</div>';

    try {
        const token = localStorage.getItem('token');
        if (!token) {
            throw new Error('No authentication token found');
        }

        const response = await fetch(`/api/posts${category ? `?category=${category}` : ''}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Failed to load posts');
        }

        const posts = await response.json();
        console.log('Posts loaded:', posts.length);

        postsFeed.innerHTML = posts.length ? posts.map(post => `
            <article class="post-card" data-post-id="${post.id}">
                <div class="post-header">
                    <h3 class="post-title">${post.title}</h3>
                    <span class="post-category">${post.category}</span>
                </div>
                <div class="post-meta">
                    Posted by ${post.author?.nickname || 'Unknown'} 
                    on ${new Date(post.createdAt).toLocaleDateString()}
                </div>
                <div class="post-actions">
                    <button class="view-comments-btn" data-post-id="${post.id}">
                        View Comments
                    </button>
                </div>
            </article>
        `).join('') : '<p class="no-posts">No posts yet</p>';

        // Update click handlers
        document.querySelectorAll('.view-comments-btn').forEach(btn => {
            btn.addEventListener('click', () => {
                const postId = btn.dataset.postId;
                showPostDetail(postId);
            });
        });
    } catch (error) {
        console.error('Error loading posts:', error);
        postsFeed.innerHTML = `
            <div class="error-message">
                <p>Failed to load posts: ${error.message}</p>
                <button onclick="loadPosts()">Try Again</button>
            </div>
        `;
    }
}

function smoothTransition(element, content, callback) {
    element.style.opacity = '0';
    setTimeout(() => {
        element.innerHTML = content;
        element.style.opacity = '1';
        if (callback) callback();
    }, 200);
}

async function showPostDetail(postId) {
    if (!postsFeed) {
        postsFeed = document.getElementById('posts-feed');
    }
    try {
        const user = JSON.parse(localStorage.getItem('user'));
        const token = localStorage.getItem('token');

        if (!token || !user) {
            throw new Error('Authentication required');
        }

        postsFeed.style.transition = 'opacity 0.2s ease-in-out';
        smoothTransition(postsFeed, '<div class="loading">Loading post...</div>');

        // Fetch post details
        const postResponse = await fetch(`/api/posts/${postId}`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!postResponse.ok) {
            throw new Error('Failed to load post');
        }

        const post = await postResponse.json();
        const postDate = post.createdAt ? new Date(post.createdAt).toLocaleString() : 'Unknown date';
        const authorName = post.author ?
            `${post.author.nickname} (${post.author.firstName} ${post.author.lastName})` :
            'Unknown';

        // Create post detail view
        smoothTransition(postsFeed, `
            <article class="post-detail">
                <button class="back-btn">← Back to Posts</button>
                
                <div class="post-header">
                    <h2 class="post-title">${post.title || 'Untitled Post'}</h2>
                    <span class="post-category">${post.category || 'Uncategorized'}</span>
                </div>
                <div class="post-meta">
                    Posted by ${authorName} on ${postDate}
                </div>
                <div class="post-content">${post.content || 'No content available'}</div>

                <div class="comments-section">
                    <h3>Comments</h3>
                    <div class="comments-list">
                        <div class="loading">Loading comments...</div>
                    </div>
                    <form class="comment-form" data-post-id="${postId}">
                        <textarea 
                            id="comment-content-${postId}"
                            name="comment-content"
                            required 
                            placeholder="Write a comment..."
                        ></textarea>
                        <button type="submit">Post Comment</button>
                    </form>
                </div>
            </article>
        `, async () => {
            const backBtn = document.querySelector('.back-btn');
            if (backBtn) {
                backBtn.addEventListener('click', () => loadPosts());
            }

            // Load comments immediately
            try {
                const commentsResponse = await fetch(`/api/posts/${postId}/comments`, {
                    headers: { 'Authorization': `Bearer ${token}` }
                });

                if (!commentsResponse.ok) {
                    throw new Error('Failed to load comments');
                }

                const comments = await commentsResponse.json();
                const commentsList = document.querySelector('.comments-list');

                if (commentsList) {
                    commentsList.innerHTML = Array.isArray(comments) && comments.length ?
                        comments.map(comment => `
                            <div class="comment ${comment.userId === user.id ? 'current-user-comment' : ''}">
                                <div class="comment-meta">
                                    <span class="author">
                                        ${comment.author ?
                                `${comment.author.nickname} (${comment.author.firstName} ${comment.author.lastName})` :
                                'Unknown'}
                                    </span>
                                    <span class="date">
                                        ${new Date(comment.createdAt).toLocaleString()}
                                    </span>
                                </div>
                                <div class="comment-content">${comment.content}</div>
                            </div>
                        `).join('') :
                        '<p class="no-comments">No comments yet</p>';
                }
            } catch (error) {
                console.error('Error loading comments:', error);
                document.querySelector('.comments-list').innerHTML =
                    '<div class="error-message">Failed to load comments</div>';
            }

            // Setup comment form
            const form = document.querySelector('.comment-form');
            if (form) {
                form.addEventListener('submit', async (e) => {
                    e.preventDefault();
                    const textarea = document.getElementById(`comment-content-${postId}`);
                    const content = textarea?.value.trim();
                    const submitBtn = form.querySelector('button[type="submit"]');

                    if (!content) {
                        alert('Please write a comment');
                        return;
                    }

                    submitBtn.disabled = true;
                    submitBtn.textContent = 'Posting...';

                    try {
                        await createComment(postId, content);
                        form.reset();
                        // Refresh comments directly instead of using viewCommentsBtn
                        const commentsResponse = await fetch(`/api/posts/${postId}/comments`, {
                            headers: { 'Authorization': `Bearer ${token}` }
                        });

                        if (!commentsResponse.ok) {
                            throw new Error('Failed to load comments');
                        }

                        const comments = await commentsResponse.json();
                        const commentsList = document.querySelector('.comments-list');

                        if (commentsList) {
                            commentsList.innerHTML = Array.isArray(comments) && comments.length ?
                                comments.map(comment => `
                                    <div class="comment ${comment.userId === user.id ? 'current-user-comment' : ''}">
                                        <div class="comment-meta">
                                            <span class="author">
                                                ${comment.author ?
                                        `${comment.author.nickname} (${comment.author.firstName} ${comment.author.lastName})` :
                                        'Unknown'}
                                            </span>
                                            <span class="date">
                                                ${new Date(comment.createdAt).toLocaleString()}
                                            </span>
                                        </div>
                                        <div class="comment-content">${comment.content}</div>
                                    </div>
                                `).join('') :
                                '<p class="no-comments">No comments yet</p>';
                        }
                    } catch (error) {
                        console.error('Error posting comment:', error);
                        alert(`Failed to post comment: ${error.message}`);
                    } finally {
                        submitBtn.disabled = false;
                        submitBtn.textContent = 'Post Comment';
                    }
                });
            }
        });

    } catch (error) {
        console.error('Post detail error:', error);
        smoothTransition(postsFeed, `
            <div class="error-message">
                <p>Error: ${error.message}</p>
                <button class="back-btn">← Back to Posts</button>
            </div>
        `);
    }
}

export { debouncedLoadPosts as loadPosts, showPostDetail };