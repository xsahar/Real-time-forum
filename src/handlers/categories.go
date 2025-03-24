package handlers

import (
    "encoding/json"
    "net/http"
)

var Categories = []string{
    "General Discussion",
    "Technology",
    "Programming",
    "Web Development",
    "Business & Finance",
    "Entrepreneurship",
    "Sports",
    "Fitness & Health",
    "Entertainment",
    "Movies & TV",
    "Music",
    "Science",
    "Space & Astronomy",
    "Gaming",
    "Art & Design",
    "Career Development",
    "Fashion & Style",
    "Pets & Animals",
    "Environment",
    "Cars & Automotive",
}

func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Categories)
}