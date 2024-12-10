package database

import (
    "Real-time-forum/src/models"
    "log"
    // "database/sql"
)

func GetAllPosts() ([]models.Post, error) {
    log.Printf("GetAllPosts: Starting query")
    
    query := `
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.nickname as author,
            COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = 1), 0) as likes,
            COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = -1), 0) as dislikes
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        ORDER BY p.created_at DESC
    `

    rows, err := DB.Query(query)
    if err != nil {
        log.Printf("GetAllPosts: Query error: %v", err)
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID,
            &post.UserID,
            &post.Title,
            &post.Content,
            &post.CreatedAt,
            &post.Author,
            &post.Likes,
            &post.Dislikes,
        )
        if err != nil {
            log.Printf("GetAllPosts: Row scan error: %v", err)
            continue
        }

        categories, err := getPostCategories(post.ID)
        if err == nil {
            post.Categories = categories
        }

        posts = append(posts, post)
    }

    log.Printf("GetAllPosts: Found %d posts", len(posts))
    return posts, nil
}

func getPostCategories(postID int) ([]string, error) {
    query := `
        SELECT c.name 
        FROM categories c
        JOIN post_categories pc ON c.id = pc.category_id
        WHERE pc.post_id = ?
    `
    
    rows, err := DB.Query(query, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var categories []string
    for rows.Next() {
        var category string
        if err := rows.Scan(&category); err != nil {
            return nil, err
        }
        categories = append(categories, category)
    }

    return categories, nil
}

func GetPostWithDetails(postID int) (*models.Post, error) {
    log.Printf("GetPostWithDetails: Fetching post %d", postID)
    
    var post models.Post
    query := `
        SELECT 
            p.id, p.user_id, p.title, p.content, p.created_at,
            u.nickname as author,
            COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = 1), 0) as likes,
            COALESCE((SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = -1), 0) as dislikes
        FROM posts p
        LEFT JOIN users u ON p.user_id = u.id
        WHERE p.id = ?`
        
    err := DB.QueryRow(query, postID).Scan(
        &post.ID, 
        &post.UserID, 
        &post.Title, 
        &post.Content, 
        &post.CreatedAt,
        &post.Author, 
        &post.Likes, 
        &post.Dislikes,
    )
    if err != nil {
        log.Printf("GetPostWithDetails: Query error: %v", err)
        return nil, err
    }

    categories, err := getPostCategories(post.ID)
    if err != nil {
        log.Printf("GetPostWithDetails: Failed to get categories: %v", err)
    } else {
        post.Categories = categories
    }

    return &post, nil
}

func AddComment(postID, userID int, content string) (int64, error) {
    result, err := DB.Exec(`
        INSERT INTO comments (post_id, user_id, content, created_at)
        VALUES (?, ?, ?, CURRENT_TIMESTAMP)`,
        postID, userID, content)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}

func GetComments(postID int) ([]models.Comment, error) {
    rows, err := DB.Query(`
        SELECT c.id, c.content, u.nickname as author, c.created_at
        FROM comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?
        ORDER BY c.created_at DESC`, postID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var comments []models.Comment
    for rows.Next() {
        var comment models.Comment
        if err := rows.Scan(&comment.ID, &comment.Content, &comment.Author, &comment.CreatedAt); err != nil {
            continue
        }
        comments = append(comments, comment)
    }
    return comments, nil
}