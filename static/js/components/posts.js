
export class PostsComponent {
    constructor(app) {
        this.app = app;
        this.api = app.api;
        this.container = document.getElementById('posts-container');
    }

    async mount() {
        try {
            await this.loadPosts();
        } catch (error) {
            this.renderError(error.message);
        }
    }

    async loadPosts() {
        try {
            const response = await this.api.get('posts');
            if (response.status === 'success') {
                this.renderPosts(response.data);
            } else {
                throw new Error(response.error || 'Failed to load posts');
            }
        } catch (error) {
            console.error('Failed to load posts:', error);
            throw error;
        }
    }

    renderError(message) {
        this.container.innerHTML = `
            <div class="error-message">
                ${this.escapeHtml(message)}
            </div>
        `;
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
        }
    
        async showComments(postId) {
            const commentsSection = document.getElementById(`comments-${postId}`);
            const isHidden = commentsSection.style.display === 'none';
            
            if (isHidden) {
                try {
                    const response = await this.api.get(`post/${postId}/comments`);
                    if (response.status === 'success') {
                        commentsSection.innerHTML = this.renderComments(response.data);
                        commentsSection.style.display = 'block';
                    }
                } catch (error) {
                    console.error('Failed to load comments:', error);
                }
            } else {
                commentsSection.style.display = 'none';
            }
        }
    
        renderComments(comments) {
            if (!comments?.length) {
                return `
                    <div class="comments-container">
                        <p>No comments yet</p>
                        <form class="comment-form" data-post-id="${postId}">
                            <textarea name="content" placeholder="Add a comment..." required></textarea>
                            <button type="submit">Post Comment</button>
                        </form>
                    </div>
                `;
            }
    
            return `
                <div class="comments-container">
                    ${comments.map(comment => `
                        <div class="comment">
                            <p>${this.escapeHtml(comment.content)}</p>
                            <small>By ${this.escapeHtml(comment.author)} on ${new Date(comment.created_at).toLocaleString()}</small>
                        </div>
                    `).join('')}
                    <form class="comment-form" data-post-id="${postId}">
                        <textarea name="content" placeholder="Add a comment..." required></textarea>
                        <button type="submit">Post Comment</button>
                    </form>
                </div>
            `;
        }
    
    escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }


    async handleLike(postId, value) {
        try {
            await this.api.post(`likes`, { post_id: postId, value });
            await this.loadPosts();
        } catch (error) {
            console.error('Failed to update like:', error);
        }
    }

    async creatPost(post) {
        try{
            const categories = await this.api.get('categories');
            const form = `
            <form id = "post-form" class="post-form">
            <input type="text" name="title" placeholder="Title" required>
            <textarea name="content" placeholder="Post Content" required></textarea>
            <div class="categories-selection">
            <h3>Select Categories</h3>
            ${categories.data.map(cat => `
            <label>
            <input type="checkbox" name="categories" value="$
            {cat.id}">
            ${this.escapeHtml(cat.name)}
            </label>
            `).join('')}
            </div>
            <button type="submit">Create Post</button>
            </form>
            `;
            this.container.innerHTML = form;
        } catch (error) {
            console.error('Failed to load categories:', error);
        }
}
}
