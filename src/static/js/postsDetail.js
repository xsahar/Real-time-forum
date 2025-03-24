import { createComment, loadComments } from './comments.js';
import { loadPosts } from './postsAPI.js';  // Import loadPosts from postsAPI.js instead

export const bindViewComments = () => {
    document.querySelectorAll('.view-comments-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const postId = btn.dataset.postId;
            showPostDetail(postId);
        });
    });
};

export async function showPostDetail(postId) {
    const postsFeed = document.getElementById('posts-feed');
    const token = localStorage.getItem('token');
    if (!postsFeed || !token) return;

    postsFeed.style.opacity = '0';
    postsFeed.innerHTML = '<div class="loading">Loading post...</div>';
    setTimeout(async () => {
        try {
            const response = await fetch(`/api/posts/${postId}`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });
            if (!response.ok) throw new Error('Failed to load post');
            const post = await response.json();
            const postDate = post.createdAt ? new Date(post.createdAt).toLocaleString() : 'Unknown date';
            const authorName = post.author ? `${post.author.nickname} (${post.author.firstName} ${post.author.lastName})` : 'Unknown';
            postsFeed.innerHTML = `
                <article class="post-detail">
                    <button class="back-btn">← Back to Posts</button>
                    <div class="post-header">
                        <h2>${post.title || 'Untitled Post'}</h2>
                        <span>${post.category || 'Uncategorized'}</span>
                    </div>
                    <div class="post-meta">
                        Posted by ${authorName} on ${postDate}
                    </div>
                    <div class="post-content">${post.content || 'No content available'}</div>
                    <div class="comments-section">
                        <h3>Comments</h3>
                        <div class="comments-list"><div class="loading">Loading comments...</div></div>
                        <form class="comment-form" data-post-id="${postId}">
                            <textarea id="comment-content-${postId}" required placeholder="Write a comment..."></textarea>
                            <button type="submit">Post Comment</button>
                        </form>
                    </div>
                </article>
            `;
            
            // Change this to use loadPosts instead of reloading the page
            document.querySelector('.back-btn').addEventListener('click', () => {
                postsFeed.style.opacity = '0';
                postsFeed.innerHTML = '<div class="loading">Loading posts...</div>';
                
                // Add a slight delay for visual transition
                setTimeout(async () => {
                    try {
                        // Load posts from postsAPI.js
                        await loadPosts();
                        
                        // Rebind view comments buttons
                        bindViewComments();
                    } catch (error) {
                        console.error('Error returning to posts:', error);
                        postsFeed.innerHTML = `<div class="error">Error: ${error.message}</div>`;
                    }
                    postsFeed.style.opacity = '1';
                }, 200);
            });
            
            // Set the global currentPostId
            window.currentPostId = postId;
            
            // Load comments and bind form
            await loadComments(postId);
            bindCommentForm(postId);
        } catch (error) {
            console.error(error);
            postsFeed.innerHTML = `<div class="error">Error: ${error.message}<button class="back-btn">← Back to Posts</button></div>`;
            
            // Handle error state back button too
            document.querySelector('.back-btn').addEventListener('click', () => {
                postsFeed.style.opacity = '0';
                postsFeed.innerHTML = '<div class="loading">Loading posts...</div>';
                setTimeout(async () => {
                    await loadPosts();
                    bindViewComments();
                    postsFeed.style.opacity = '1';
                }, 200);
            });
        }
        postsFeed.style.opacity = '1';
    }, 200);
}

function bindCommentForm(postId) {
    const form = document.querySelector('.comment-form');
    const token = localStorage.getItem('token');
    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        const textarea = document.getElementById(`comment-content-${postId}`);
        const content = textarea.value.trim();
        const submitBtn = form.querySelector('button[type="submit"]');

        if (!content) return alert('Write a comment');

        submitBtn.disabled = true;
        submitBtn.textContent = 'Posting...';

        try {
            await createComment(postId, content);
            form.reset();
            
            // Reload comments after posting
            await loadComments(postId);
            
        } catch (error) {
            console.error(error);
            alert('Error posting comment');
        } finally {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Post Comment';
        }
    });
}
