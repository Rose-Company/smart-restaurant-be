# Menu Items API - Example Requests

## Overview
Menu Items APIs support full CRUD operations with image upload to DigitalOcean Storage and modifier group associations.

---

## 1. POST /api/admin/menu/items/upload-image - Upload Image to Storage

### Example: Upload single image file
```bash
curl -X POST "http://localhost:8080/api/admin/upload" \
  -F "file=@/path/to/image.jpg"
```

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593200.jpg"
  }
}
```

**Validation Rules:**
- File type: Only `.jpg`, `.jpeg`, `.png` allowed
- Max file size: 10MB
- Returns: Full URL to uploaded image

**Error Response (400 Bad Request - Invalid file type):**
```json
{
  "error": "Only PNG and JPG images are allowed"
}
```

**Error Response (400 Bad Request - File too large):**
```json
{
  "error": "File size must be less than 10MB"
}
```

**Workflow:**
1. Upload image using this endpoint
2. Get back the `url` from response
3. Use this URL in `images` array when creating/updating menu item

---

## 2. GET /api/admin/menu/items - List all menu items

### Example 1: Get all items (default pagination)
```bash
curl -X GET "http://localhost:8080/api/admin/menu/items"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 10,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "name": "Grilled Salmon",
        "category": "Main Course",
        "price": 24.99,
        "status": "Available",
        "last_update": "2025-12-20",
        "chef_recommended": true,
        "image_url": "https://images.unsplash.com/photo-1580476262798-bddd9f4b7369?w=400",
        "description": "Fresh Atlantic salmon grilled to perfection",
        "preparation_time": 15
      }
    ]
  }
}
```

### Example 2: Filter by category
```bash
curl -X GET "http://localhost:8080/api/admin/menu/items?category=Main%20Course"
```

### Example 3: Filter by status
```bash
curl -X GET "http://localhost:8080/api/admin/menu/items?status=available"
```

**Status values:** `available`, `unavailable`, `sold_out`, `all`

### Example 4: Search by name
```bash
curl -X GET "http://localhost:8080/api/admin/menu/items?search=salmon"
```

### Example 5: Sort options
```bash
# Sort by name
curl -X GET "http://localhost:8080/api/admin/menu/items?sort=name"

# Sort by price ascending
curl -X GET "http://localhost:8080/api/admin/menu/items?sort=price_asc"

# Sort by price descending
curl -X GET "http://localhost:8080/api/admin/menu/items?sort=price_desc"

# Sort by last update
curl -X GET "http://localhost:8080/api/admin/menu/items?sort=last_update"
```

### Example 6: Combined filters
```bash
curl -X GET "http://localhost:8080/api/admin/menu/items?page=1&page_size=10&category=appetizer&status=available&search=salad&sort=price_asc"
```

---

## 3. POST /api/admin/menu/items - Create new menu item

### Example 1: Create item with all fields including images and modifiers
```bash
curl -X POST "http://localhost:8080/api/admin/menu/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Grilled Salmon",
    "description": "Fresh Atlantic salmon grilled to perfection with herbs",
    "price": 24.99,
    "preparation_time": 15,
    "category_id": 2,
    "chef_recommended": true,
    "status": "available",
    "images": [
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593200.jpg",
        "is_primary": true
      },
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593201.jpg",
        "is_primary": false
      }
    ],
    "modifiers": [
      {
        "modifier_group_id": "1"
      },
      {
        "modifier_group_id": "2"
      }
    ]
  }'
```

**Response (201 Created):**
```json
{
  "code": 0,
  "message": "Menu item created successfully",
  "data": {
    "id": 11,
    "restaurant_id": 1,
    "category_id": 2,
    "name": "Grilled Salmon",
    "description": "Fresh Atlantic salmon grilled to perfection with herbs",
    "price": 24.99,
    "prep_time_minutes": 15,
    "status": "available",
    "is_chef_recommended": true,
    "is_deleted": false,
    "created_at": "2025-12-26T10:30:00Z",
    "updated_at": "2025-12-26T10:30:00Z"
  }
}
```

### Example 2: Create minimal item (without images and modifiers)
```bash
curl -X POST "http://localhost:8080/api/admin/menu/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Caesar Salad",
    "price": 12.50,
    "category_id": 1,
    "status": "available",
    "preparation_time": 10
  }'
```

**Validation Rules:**
- `name`: Required, max 80 characters
- `category_id`: Required, must exist in menu_categories table
- `price`: Required, must be greater than 0
- `preparation_time`: Optional, 0-240 minutes
- `status`: Required, one of: `available`, `unavailable`, `sold_out`
- `chef_recommended`: Optional boolean, defaults to false
- `description`: Optional text
- `images`: Optional array of image objects
  - `url`: Required if images provided
  - `is_primary`: Optional boolean, defaults to false
- `modifiers`: Optional array of modifier group references
  - `modifier_group_id`: Required if modifiers provided

**Workflow for creating item with images:**
1. Upload each image via `/api/admin/menu/items/upload-image`
2. Collect all returned URLs
3. Create menu item with `images` array containing those URLs
4. Mark one image as `is_primary: true` for the main display image

---

## 4. GET /api/admin/menu/items/:id - Get menu item detail

### Example: Get item by ID
```bash
curl -X GET "http://localhost:8080/api/admin/menu/items/1"
```

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "name": "Grilled Salmon",
    "category": "Main Course",
    "price": 24.99,
    "status": "Available",
    "last_update": "2025-12-20",
    "chef_recommended": true,
    "image_url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593200.jpg",
    "description": "Fresh Atlantic salmon grilled to perfection with herbs and lemon butter",
    "preparation_time": 15,
    "images": [
      {
        "id": "img_001",
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593200.jpg",
        "is_primary": true
      },
      {
        "id": "img_002",
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593201.jpg",
        "is_primary": false
      }
    ],
    "modifiers": [
      {
        "id": "size",
        "modifier_group_id": "1",
        "name": "Size Selection",
        "required": true,
        "selection_type": "Single",
        "options_preview": "Small (+$0), Medium (+$2), Large (+$4)"
      },
      {
        "id": "toppings",
        "modifier_group_id": "2",
        "name": "Extra Toppings",
        "required": false,
        "selection_type": "Multi",
        "options_preview": "Cheese, Bacon, Mushrooms..."
      }
    ]
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "code": 404,
  "message": "Menu item not found",
  "data": null
}
```

---

## 5. PUT /api/admin/menu/items/:id - Update menu item

### Example 1: Update all fields including images
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/items/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Grilled Salmon Supreme",
    "description": "Premium Atlantic salmon with herbs and garlic butter",
    "price": 26.99,
    "preparation_time": 20,
    "category_id": 2,
    "chef_recommended": true,
    "status": "available",
    "images": [
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/new-image.jpg",
        "is_primary": true
      }
    ],
    "modifiers": [
      {
        "modifier_group_id": "1"
      }
    ]
  }'
```

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "restaurant_id": 1,
    "category_id": 2,
    "name": "Grilled Salmon Supreme",
    "description": "Premium Atlantic salmon with herbs and garlic butter",
    "price": 26.99,
    "prep_time_minutes": 20,
    "status": "available",
    "is_chef_recommended": true,
    "is_deleted": false,
    "created_at": "2025-12-26T10:00:00Z",
    "updated_at": "2025-12-26T11:15:00Z"
  }
}
```

### Example 2: Update only price and status
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/items/2" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 13.99,
    "status": "unavailable"
  }'
```

### Example 3: Replace all images
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/items/3" \
  -H "Content-Type: application/json" \
  -d '{
    "images": [
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/updated-1.jpg",
        "is_primary": true
      },
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/updated-2.jpg",
        "is_primary": false
      }
    ]
  }'
```

**Note:** When updating images or modifiers, all existing ones are deleted and replaced with the new data.

**Validation Rules:**
- All fields are optional
- Same validation as create for fields that are provided
- `category_id`: Must exist in menu_categories if provided
- `images`: If provided, replaces ALL existing images
- `modifiers`: If provided, replaces ALL existing modifier associations

---

## 6. DELETE /api/admin/menu/items/:id - Delete menu item (Soft delete)

### Example: Delete menu item
```bash
curl -X DELETE "http://localhost:8080/api/admin/menu/items/1"
```

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "message": "Menu item deleted successfully"
  }
}
```

**Error Response (404 Not Found):**
```json
{
  "code": 404,
  "message": "Menu item not found",
  "data": null
}
```

**Note:** This is a soft delete - sets `is_deleted = true`. The item won't appear in listings but remains in the database.

---

## Complete Workflow Example

### Creating a menu item with images:

```bash
# Step 1: Upload primary image
curl -X POST "http://localhost:8080/api/admin/menu/items/upload-image" \
  -F "file=@salmon-primary.jpg"
# Response: {"data": {"url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593200.jpg"}}

# Step 2: Upload secondary image
curl -X POST "http://localhost:8080/api/admin/menu/items/upload-image" \
  -F "file=@salmon-side.jpg"
# Response: {"data": {"url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593201.jpg"}}

# Step 3: Create menu item with both images
curl -X POST "http://localhost:8080/api/admin/menu/items" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Grilled Salmon",
    "description": "Fresh Atlantic salmon",
    "price": 24.99,
    "preparation_time": 15,
    "category_id": 2,
    "chef_recommended": true,
    "status": "available",
    "images": [
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593200.jpg",
        "is_primary": true
      },
      {
        "url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/1703593201.jpg",
        "is_primary": false
      }
    ],
    "modifiers": [
      {
        "modifier_group_id": "1"
      }
    ]
  }'

# Step 4: Verify creation
curl -X GET "http://localhost:8080/api/admin/menu/items/11"

# Step 5: Update price
curl -X PUT "http://localhost:8080/api/admin/menu/items/11" \
  -H "Content-Type: application/json" \
  -d '{"price": 26.99}'
```

---

## Common Error Responses

### 400 Bad Request - Validation Error
```json
{
  "code": 400,
  "message": "Validation failed: category_id is required",
  "data": null
}
```

### 400 Bad Request - Category Not Found
```json
{
  "code": 400,
  "message": "Category does not exist",
  "data": null
}
```

### 404 Not Found - Item Not Found
```json
{
  "code": 404,
  "message": "Menu item not found",
  "data": null
}
```

### 500 Internal Server Error
```json
{
  "code": 500,
  "message": "Internal server error",
  "data": null
}
```

---

## Notes

- **Image Storage**: Images are uploaded to DigitalOcean Spaces bucket `smart-restaurant` with prefix `menu-items/`
- **Primary Image**: The `image_url` field in list response shows the primary image only
- **Soft Delete**: Deleted items have `is_deleted = true` and are filtered from all listings
- **Status Display**: Database stores lowercase (`available`, `sold_out`) but API returns title case (`Available`, `Sold Out`)
- **Preparation Time**: Stored as `prep_time_minutes` in DB, exposed as `preparation_time` in API
- **Chef Recommended**: JSON field uses `chef_recommended` (snake_case in responses, no underscore in requests)
- **Modifiers**: Associations are recreated on update (all old associations deleted, new ones created)
- **Images**: Complete replacement on update - provide all images you want to keep
