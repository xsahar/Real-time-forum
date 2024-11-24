package database

import (
    "Real-time-forum/src/models"
    // "database/sql"
)

func GetAllPosts() ([]models.Post, error) {
    query := `
        SELECT 
            p.id, p.title, p.content, p.user_id,
            u.nickname as author,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = 1) as likes,
            (SELECT COUNT(*) FROM likes WHERE post_id = p.id AND value = -1) as dislikes,
            p.created_at
        FROM posts p
        JOIN users u ON p.user_id = u.id 
        ORDER BY p.created_at DESC
    `

    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        err := rows.Scan(
            &post.ID, &post.Title, &post.Content, &post.UserID,
            &post.Author, &post.Likes, &post.Dislikes,
            &post.CreatedAt,
        )
        if err != nil {
            return nil, err
        }

        // Get categories for post
        post.Categories, _ = getPostCategories(post.ID)
        posts = append(posts, post)
    }

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