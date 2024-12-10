import { PostsList } from './lists.js';
import { PostCreator } from './create.js';
import { CategoriesComponent } from '../categories/index.js';

export class PostsComponent {
    constructor(app) {
        this.app = app;
        this.api = app.api;
        this.container = document.getElementById('posts-container');
        
        this.postsList = new PostsList(this.container, this.api, app);
        this.postCreator = new PostCreator(this.container, this.api);
        this.categories = new CategoriesComponent(app);
        
        // Bind methods
        this.mount = this.mount.bind(this);
        this.createPost = this.createPost.bind(this);
        this.filterByCategory = this.filterByCategory.bind(this);
    }

    async mount() {
        try {
            // Create layout
            this.container.innerHTML = `
                <div class="posts-layout">
                    <div class="sidebar"></div>
                    <div class="main-content"></div>
                </div>
            `;

            const sidebar = this.container.querySelector('.sidebar');
            const mainContent = this.container.querySelector('.main-content');

            // Mount categories
            await this.categories.mount();
            sidebar.appendChild(this.categories.container);

            // Mount posts
            await this.postsList.loadPosts();
            mainContent.appendChild(this.postsList.container);
        } catch (error) {
            this.renderError(error.message);
        }
    }

    async createPost() {
        try {
            this.postCreator.onPostCreated = this.mount;
            await this.postCreator.showCreateForm();
        } catch (error) {
            this.renderError('Failed to load create post form');
        }
    }

    async filterByCategory(categoryId) {
        try {
            await this.postsList.loadPosts(categoryId);
        } catch (error) {
            this.renderError('Failed to filter posts');
        }
    }

    renderError(message) {
        this.container.innerHTML = `
            <div class="error-message">
                ${this.escapeHtml(message)}
            </div>
        `;
    }

    escapeHtml(str) {
        if (!str) return '';
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
}