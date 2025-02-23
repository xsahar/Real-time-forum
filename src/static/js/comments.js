export const initComments = () => {
    document.addEventListener('click', async (e) => {
        if (e.target.classList.contains('view-comments-btn')) {
            const postId = e.target.dataset.postId;
            await loadComments(postId);
        }
    });

    const form = document.querySelector('.comment-form');
    if (form) {
        form.addEventListener('submit', async (e) => {
            e.preventDefault();
            const textarea = form.querySelector('textarea');
            const submitBtn = form.querySelector('button');

            if (!textarea.value.trim()) {
                alert('Please write a comment');
                return;
            }

            submitBtn.disabled = true;
            submitBtn.textContent = 'Posting...';

            try {
                const comment = await createComment(postId, textarea.value.trim());
                const commentsList = document.querySelector('.comments-list');
                const noComments = commentsList.querySelector('.no-comments');

                if (noComments) {
                    noComments.remove();
                }

                const commentElement = document.createElement('div');
                commentElement.className = 'comment current-user-comment';
                commentElement.innerHTML = `
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
                `;

                commentsList.insertBefore(commentElement, commentsList.firstChild);
                textarea.value = '';
            } catch (error) {
                alert(error.message);
            } finally {
                submitBtn.disabled = false;
                submitBtn.textContent = 'Post Comment';
            }
        });
    }
};

export async function loadComments(postId) {
    try {
        const token = localStorage.getItem('token');
        const user = JSON.parse(localStorage.getItem('user'));

        if (!token || !user) {
            throw new Error('Authentication required');
        }

        const response = await fetch(`/api/posts/${postId}/comments`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });

        if (!response.ok) {
            throw new Error('Failed to load comments');
        }

        const comments = await response.json();

        // Sort comments - current user's comments first, then by date
        comments.sort((a, b) => {
            if (a.userId === user.id && b.userId !== user.id) return -1;
            if (a.userId !== user.id && b.userId === user.id) return 1;
            return new Date(b.createdAt) - new Date(a.createdAt);
        });

        return comments.map(comment => `
            <div class="comment ${comment.userId === user.id ? 'current-user-comment' : ''}">
                <div class="comment-meta">
                    <span class="author">
                        ${comment.author ? `${comment.author.nickname} (${comment.author.firstName} ${comment.author.lastName})` : 'Unknown'}
                    </span>
                    <span class="date">${new Date(comment.createdAt).toLocaleString()}</span>
                </div>
                <div class="comment-content">${comment.content}</div>
            </div>
        `).join('');

    } catch (error) {
        console.error('Error loading comments:', error);
        throw error;
    }
}

export async function createComment(postId, content) {
    try {
        const user = JSON.parse(localStorage.getItem('user'));
        const token = localStorage.getItem('token');

        if (!user || !token) {
            throw new Error('Authentication required');
        }

        const response = await fetch(`/api/posts/${postId}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
                content,
                userId: user.id,
                postId
            })
        });

        const data = await response.json();

        if (!response.ok || !data.id) {
            throw new Error(data.message || 'Failed to create comment');
        }

        return data; // Will include author info
    } catch (error) {
        console.error('Error creating comment:', error);
        throw error;
    }
}