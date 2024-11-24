
`3-Features/frontend/posts/crud.md`:
```markdown
# Post Management

## Operations
```mermaid
graph TD
    A[Create Post] --> B[Add Categories]
    B --> C[Store in DB]
    D[Edit Post] --> E[Update Content]
    F[Delete Post] --> G[Remove Comments]
    H[Read Post] --> I[Load Comments]
    ```

Implementation
```
{
    "PostsComponent": {
        "create": "Create new post",
        "update": "Edit existing post",
        "delete": "Remove post",
        "list": "Show all posts"
    }
}
```
