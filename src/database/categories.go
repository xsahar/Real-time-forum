// src/database/categories.go
package database

import (
	// "database/sql"
	"Real-time-forum/src/models"
	"errors"
	"fmt"
	"log"
	"strings"
)

var (
    ErrCategoryNotFound = errors.New("category not found")
    ErrCategoryExists   = errors.New("category already exists")
    ErrInvalidCategory = errors.New("invalid category data")
)

// CreateCategory inserts a new category into the database
func CreateCategory(category models.Category) (*models.Category, error) {
    if err := validateCategory(category); err != nil {
        return nil, err
    }

    // Check for existing category with same name
    exists, err := categoryExists(category.Name)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, ErrCategoryExists
    }

    query := `INSERT INTO categories (name, description) VALUES (?, ?)`
    
    result, err := DB.Exec(query, category.Name, category.Description)
    if err != nil {
        return nil, fmt.Errorf("failed to create category: %v", err)
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("failed to get last insert id: %v", err)
    }

    category.ID = int(id)
    return &category, nil
}

// GetCategories retrieves all categories, optionally filtered by search term
// Add debug logging
func GetCategories(search string) ([]models.Category, error) {
    log.Printf("Getting categories with search: %v", search)

    var query string
    var args []interface{}

    if search != "" {
        query = `SELECT id, name, description, created_at FROM categories 
                 WHERE name LIKE ? OR description LIKE ?
                 ORDER BY name ASC`
        args = []interface{}{"%"+search+"%", "%"+search+"%"}
    } else {
        query = `SELECT id, name, description, created_at FROM categories ORDER BY name ASC`
    }

    log.Printf("Executing query: %s", query)
    rows, err := DB.Query(query, args...)
    if err != nil {
        log.Printf("Database error: %v", err)
        return nil, fmt.Errorf("failed to query categories: %v", err)
    }
    defer rows.Close()

    var categories []models.Category
    for rows.Next() {
        var cat models.Category
        if err := rows.Scan(
            &cat.ID,
            &cat.Name,
            &cat.Description, // Will now handle NULL properly
            &cat.CreatedAt,
        ); err != nil {
            log.Printf("Scan error: %v", err)
            return nil, fmt.Errorf("failed to scan category: %v", err)
        }
        categories = append(categories, cat)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating categories: %v", err)
    }

    log.Printf("Found %d categories", len(categories))
    return categories, nil
}

// UpdateCategory updates an existing category
func UpdateCategory(category models.Category) (*models.Category, error) {
    if err := validateCategory(category); err != nil {
        return nil, err
    }

    // Check if category exists
    exists, err := categoryExistsById(category.ID)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, ErrCategoryNotFound
    }

    // Start transaction
    tx, err := DB.Begin()
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %v", err)
    }

    query := `UPDATE categories SET name = ?, description = ? WHERE id = ?`
    
    result, err := tx.Exec(query, category.Name, category.Description, category.ID)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to update category: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        tx.Rollback()
        return nil, ErrCategoryNotFound
    }

    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %v", err)
    }

    return &category, nil
}

// DeleteCategory removes a category by ID
func DeleteCategory(id int) error {
    // Check if category exists
    exists, err := categoryExistsById(id)
    if err != nil {
        return err
    }
    if !exists {
        return ErrCategoryNotFound
    }

    // Start transaction
    tx, err := DB.Begin()
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %v", err)
    }

    query := `DELETE FROM categories WHERE id = ?`
    
    result, err := tx.Exec(query, id)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to delete category: %v", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to get rows affected: %v", err)
    }

    if rowsAffected == 0 {
        tx.Rollback()
        return ErrCategoryNotFound
    }

    if err = tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %v", err)
    }

    return nil
}

// Helper functions
func validateCategory(category models.Category) error {
    if strings.TrimSpace(category.Name) == "" {
        return ErrInvalidCategory
    }
    if len(category.Name) > 50 {
        return fmt.Errorf("category name too long (max 50 characters)")
    }
    return nil
}

func categoryExists(name string) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)`
    err := DB.QueryRow(query, name).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check category existence: %v", err)
    }
    return exists, nil
}

func categoryExistsById(id int) (bool, error) {
    var exists bool
    query := `SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)`
    err := DB.QueryRow(query, id).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check category existence: %v", err)
    }
    return exists, nil
}