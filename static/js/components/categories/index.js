// static/js/components/categories/index.js
export class CategoriesComponent {
    constructor(app) {
        this.app = app;
        this.api = app.api;
        this.container = document.createElement('div');
        this.container.className = 'categories-section';
    }

    async mount() {
        try {
            const response = await this.api.get('categories');
            console.log('Categories response:', response)
            if (response.status === 'success') {
                this.render(response.data);
                this.setupListeners();
            }
        } catch (error) {
            console.error('Failed to load categories:', error);
        }
    }

    render(categories) {
        console.log('Rendering Categories:', categories)
        if (!categories?.length) {
            this.container.innerHTML = '<p>No categories available</p>';
            return;
        }
        this.container.innerHTML = `
            <h3>Categories</h3>
            <div class="categories-filter">
                <select id="category-filter">
                    <option value="">All Categories</option>
                    ${categories.map(cat => 
                        `<option value="${cat.id}">${this.escapeHtml(cat.name)}</option>`
                    ).join('')}
                </select>
            </div>
        `;
    }

    setupListeners() {
        const select = this.container.querySelector('#category-filter');
        select?.addEventListener('change', (e) => {
            const categoryId = e.target.value;
            this.app.components.posts?.filterByCategory(categoryId);
        });
    }

    escapeHtml(str) {
        const div = document.createElement('div');
        div.textContent = str;
        return div.innerHTML;
    }
}