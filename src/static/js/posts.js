import { loadCategories, loadPosts, handlePostSubmission } from './postsAPI.js';
import { bindViewComments } from './postsDetail.js';

let postsFeed = null;

export const initPosts = async () => {
    const createPostForm = document.getElementById('create-post-form');
    const createPostBtn = document.getElementById('create-post-btn');
    const cancelPostBtn = document.getElementById('cancel-post');
    const createPostSection = document.getElementById('create-post-section');
    postsFeed = document.getElementById('posts-feed');

    if (!createPostForm || !createPostSection || !postsFeed) {
        console.error('Missing required elements');
        return;
    }

    // Load categories and posts
    await loadCategories();
    await loadPosts();
    
    // Explicitly rebind view comments after posts are loaded
    bindViewComments();

    // Bind UI events for post creation
    createPostBtn.addEventListener('click', (e) => {
        e.preventDefault();
        createPostSection.classList.remove('hidden');
        postsFeed.classList.add('hidden');
    });
    
    cancelPostBtn.addEventListener('click', () => {
        createPostSection.classList.add('hidden');
        postsFeed.classList.remove('hidden');
        createPostForm.reset();
    });

    createPostForm.addEventListener('submit', async (e) => {
        await handlePostSubmission(e);
        // Explicitly rebind view comments after post creation
        bindViewComments();
    });
    
    // Set up event delegation for all current and future post interactions
    document.addEventListener('click', (e) => {
        // If the click was on a view comments button
        if (e.target.classList.contains('view-comments-btn')) {
            const postId = e.target.dataset.postId;
            if (postId && typeof window.showPostDetail === 'function') {
                window.showPostDetail(postId);
            }
        }
    });
};

// Make the function globally available
window.refreshPostsWithCommentBinding = async () => {
    await loadPosts();
    bindViewComments();
};
