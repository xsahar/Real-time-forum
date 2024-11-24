# Like/Dislike System

## Overview
Real-time like/dislike functionality for posts with instant count updates.

## Features
- Like/Dislike posts
- Toggle between like/dislike
- Real-time counter updates
- User-specific tracking



### URL Endpoints
```http
POST /api/like/{postID}    # Add/toggle like
POST /api/dislike/{postID} # Add/toggle dislike

Response Formats
Success (200 OK)
```
{
    "status": "success",
    "message": "Like added successfully",
    "data": {
        "likes": 1,
        "dislikes": 0
    }
}
```

Error (4xx/5xx)

```
{
    "status": "error",
    "message": "Error description"
}
```

