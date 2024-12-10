export class PostCreator {
    constructor(container, api) {
        this.container = container;
        this.api = api;
    }

    async showCreateForm() {
        try {
            const categories = await this.api.get('categories');
            const form = this.renderCreateForm(categories.data);
            this.container.innerHTML = form;
            this.setupFormListener();
        } catch (error) {
            console.error('Failed to load categories:', error);
            throw error;
        }
    }

    renderCreateForm(categories) {
        return `
            <form id="post-form" class="post-form">
                <input type="text" name="title" placeholder="Title" required>
                <textarea name="content" placeholder="Post Content" required></textarea>
                <div class="categories-selection">
                    <h3>Select Categories</h3>
                    ${categories.map(cat => `
                        <label>
                            <input type="checkbox" name="categories" value="${cat.id}">
                            ${this.escapeHtml(cat.name)}
                        </label>
                    `).join('')}
                </div>
                <button type="submit">Create Post</button>
            </form>
        `;
    }

    setupFormListener() {
        const form = document.getElementById('post-form');
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const formData = new FormData(e.target);
            const post = {
                title: formData.get('title'),
                content: formData.get('content'),
                categories: formData.getAll('categories')
            };
            await this.submitPost(post);
        });
    }

    async submitPost(post) {
        try {
            const response = await this.api.post('posts', post);
            if (response.status === 'success') {
                await this.onPostCreated?.();
            }
        } catch (error) {
            console.error('Failed to create post:', error);
            throw error;
        }
    }

    escapeHtml(str) {
        if (!str) return '';
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
}