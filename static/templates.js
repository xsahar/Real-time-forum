// static/templates.js
const templates = {
    home: `
        <section id="filter">
            <h2>Filter Posts</h2>
            <form id="filter-form" action="/" method="GET">
                <div class="filter-group">
                    <label for="category-select">Category:</label>
                    <select name="category" id="category-select">
                        <option value="">All Categories</option>
                        <!-- Categories will be dynamically loaded -->
                    </select>
                </div>

                <div id="post-type-filter" class="filter-group" style="display: none;">
                    <label>Post Type:</label>
                    <div class="radio-group">
                        <input type="radio" id="allPosts" name="filter" value="all" checked>
                        <label for="allPosts">All Posts</label>
                        <input type="radio" id="myPosts" name="filter" value="my-posts">
                        <label for="myPosts">My Posts</label>
                        <input type="radio" id="likedPosts" name="filter" value="liked-posts">
                        <label for="likedPosts">Liked Posts</label>
                    </div>
                </div>

                <button type="submit">Apply Filter</button>
            </form>
        </section>

        <section id="posts">
            <h2>Recent Posts</h2>
            <div id="posts-container" class="grid">
                <!-- Posts will be dynamically loaded -->
            </div>
        </section>
    `,
    login: `
        <div class="auth-form-container">
            <h2>Login</h2>
            <form id="login-form" class="auth-form" action="/login" method="POST">
                <input type="text" name="identifier" placeholder="Nickname or Email" required>
                <input type="password" name="password" placeholder="Password" required>
                <button type="submit">Login</button>
            </form>
            <p class="auth-link">Don't have an account? <a href="#" id="show-register">Register</a></p>
        </div>
    `,
    register: `
        <div class="auth-form-container">
            <h2>Register</h2>
            <form id="registration-form" class="auth-form" action="/register" method="POST">
            <input type="text" name="email" placeholder="Email" required>
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <input type="text" name="nickname" placeholder="Nickname" required>
            <input type="text" name="age" placeholder="Age" required>
            <input type="text" name="gender" placeholder="Gender" required>
            <input type="text" name="first_name" placeholder="First Name" required>
            <input type="text" name="last_name" placeholder="Last Name" required>
            <button type="submit">Register</button>
            </form>
            <p class="auth-link">Already have an account? <a href="#" id="show-login">Login</a></p>
        </div>
    `,
    newPost: `
        <div class="post-form-container">
            <h2>New Post</h2>
            <form id="new-post-form" class="post-form" action="/new-post" method="POST">
                <input type="text" name="title" placeholder="Title" required>
                <textarea name="content" placeholder="Content" required></textarea>
                <button type="submit">Post</button>
            </form>
        </div>
    `,
    post: `
        <header>
            <nav class="container">
                <div class="logo">Forum</div>
                <ul id="nav-links">
                    <!-- Navigation links will be dynamically inserted here -->
                </ul>
            </nav>
        </header>

        <main class="container">
            <div id="error-messages"></div>
            <article class="post card">
                <h2 id="post-title"></h2>
                <p id="post-content"></p>
                <div class="post-meta">
                    <span class="category" id="post-categories">
                        <i class="fas fa-tag"></i>
                        <!-- Categories will be dynamically loaded -->
                    </span>            
                    <span class="author" id="post-author"></span>
                </div>
                <div class="post-actions" id="post-actions">
                    <!-- Post actions will be dynamically loaded -->
                </div>
            </article>

            <section id="comments" class="card">
                <h3>Comments</h3>
                <div id="comments-container">
                    <!-- Comments will be dynamically loaded -->
                </div>
            </section>
       
            <section id="add-comment" class="card" style="display: none;">
                <h3>Add a Comment</h3>
                <form id="comment-form">
                    <textarea name="content" placeholder="Write your comment here..." required minlength="10" maxlength="500"></textarea>
                    <button type="submit">Submit Comment</button>
                </form>
            </section>

            <section id="login-to-comment" class="card" style="display: none;">
                <p>Please <a href="#" id="login-link">log in</a> to post comments.</p>
            </section>
        </main>

        <footer>
            <p>&copy; 2024 Forum. All rights reserved.</p>
        </footer>
    `,
    categories: `
        <header>
            <nav class="container">
                <div class="logo">Forum</div>
                <ul id="nav-links">
                    <!-- Navigation links will be dynamically inserted here -->
                </ul>
            </nav>
        </header>

        <main class="container">
            <h1>Categories</h1>
            <ul class="categories-list" id="categories-list">
                <!-- Categories will be dynamically loaded -->
            </ul>
        </main>

        <section id="add-category" style="display: none;">
            <h3>Add a Custom Category</h3>
            <form id="add-category-form">
                <textarea name="category" required></textarea>
                <button type="submit">Submit Category</button>
            </form>
        </section>

        <footer>
            <p>&copy; 2024 Forum. All rights reserved.</p>
        </footer>
    `
};