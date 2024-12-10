export class PostsList {
    constructor(container, api, app) {
        this.container = container;
        this.api = api;
        this.app = app;
    }

    async loadPosts(categoryId = null) {
        try {
            // Show loading state
            this.container.innerHTML = '<div class="loading">Loading posts...</div>';
    
            // Build endpoint with optional category filter
            const endpoint = categoryId ? `posts?category=${categoryId}` : 'posts';
            console.log('Loading posts from:', endpoint);
    
            const response = await this.api.get(endpoint);
            console.log('Posts response:', response);
    
            if (response.status === 'success') {
                console.log(`Loaded ${response.data?.length || 0} posts`);
                this.renderPosts(response.data);
            } else {
                throw new Error(response.error || 'Failed to load posts');
            }
        } catch (error) {
            console.error('Failed to load posts:', error);
            this.container.innerHTML = `
                <div class="error-message">
                    Failed to load posts: ${error.message}
                </div>
            `;
            throw error;
        }
    }

    renderPosts(posts) {
        if (!posts?.length) {
            this.container.innerHTML = '<p>No posts available</p>';
            return;
        }
    
        this.container.innerHTML = posts
            .map(post => this.renderPost(post))
            .join('');
        
        this.setupEventListeners();
    }

    renderPost(post) {
        return `
            <article class="post" data-id="${post.id}">
                <h2>${this.escapeHtml(post.title)}</h2>
                <div class="post-categories">
                    ${post.categories?.map(cat => 
                        `<span class="category">${this.escapeHtml(cat)}</span>`
                    ).join('') || ''}
                </div>
                <p>${this.escapeHtml(post.content)}</p>
                <div class="post-meta">
                    <span>By ${this.escapeHtml(post.author)}</span>
                    <div class="post-actions">
                        <button class="like-btn" data-post-id="${post.id}">👍 ${post.likes || 0}</button>
                        <button class="dislike-btn" data-post-id="${post.id}">👎 ${post.dislikes || 0}</button>
                        <button class="comment-btn" data-post-id="${post.id}">💬 Comment</button>
                    </div>
                </div>
                <div class="comments-section" id="comments-${post.id}" style="display: none;"></div>
            </article>
        `;
    }

    setupEventListeners() {
        this.container.addEventListener('click', async (e) => {
            if (!this.app.checkAuth()) {
                alert('Please login to interact with posts');
                return;
            }

            const postId = e.target.closest('.post')?.dataset.id;
            if (!postId) return;

            if (e.target.classList.contains('like-btn')) {
                await this.handleLike(postId, 1);
            } else if (e.target.classList.contains('dislike-btn')) {
                await this.handleLike(postId, -1);
            } else if (e.target.classList.contains('comment-btn')) {
                await this.showComments(postId);
            }
        });

        this.setupCommentListeners();
    }

    setupCommentListeners() {
        this.container.addEventListener('submit', async (e) => {
            if (e.target.classList.contains('comment-form')) {
                e.preventDefault();
                const postId = e.target.dataset.postId;
                const content = e.target.content.value;
                await this.submitComment(postId, content);
            }
        });
    }

    // ... (comment methods and utilities)
}