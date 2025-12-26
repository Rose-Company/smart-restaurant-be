# Menu Categories API - Example Requests

## 1. GET /api/admin/menu/categories - List all menu categories

### Example 1: Get all categories (default pagination)
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 5,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "name": "Appetizers",
        "description": "Start your meal with our delicious starters",
        "item_count": 8,
        "is_active": true,
        "display_order": 1
      },
      {
        "id": 2,
        "name": "Main Course",
        "description": "Hearty and satisfying main dishes",
        "item_count": 15,
        "is_active": true,
        "display_order": 2
      },
      {
        "id": 3,
        "name": "Desserts",
        "description": "Sweet treats to end your meal",
        "item_count": 6,
        "is_active": true,
        "display_order": 3
      }
    ]
  }
}
```

### Example 2: Get categories with pagination
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?page=1&page_size=10"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 25,
    "page": 1,
    "page_size": 10,
    "items": [...]
  }
}
```

### Example 3: Search categories by name
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?search=salad"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 2,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 4,
        "name": "Salads",
        "description": "Fresh and healthy salad options",
        "item_count": 5,
        "is_active": true,
        "display_order": 4
      },
      {
        "id": 12,
        "name": "Caesar Salad Special",
        "description": "Our signature salad",
        "item_count": 3,
        "is_active": true,
        "display_order": 5
      }
    ]
  }
}
```

### Example 4: Filter by status (active categories only)
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?status=active"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 8,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "name": "Appetizers",
        "description": "Start your meal with our delicious starters",
        "item_count": 8,
        "is_active": true,
        "display_order": 1
      }
    ]
  }
}
```

### Example 5: Filter inactive categories
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?status=inactive"
```

### Example 6: Sort by display order (default)
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?sort=displayOrder"
```

### Example 7: Sort by name
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?sort=name"
```

### Example 8: Sort by item count
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?sort=itemCount"
```

### Example 9: Combined filters
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories?page=1&page_size=10&search=app&status=active&sort=displayOrder"
```

---

## 2. POST /api/admin/menu/categories - Create new category

### Example 1: Create category with all fields
```bash
curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Salads",
    "description": "Fresh and healthy salad options",
    "display_order": 4,
    "status": "active"
  }'
```

**Response (201 Created):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 6,
    "restaurant_id": 1,
    "name": "Salads",
    "description": "Fresh and healthy salad options",
    "display_order": 4,
    "status": "active",
    "created_at": "2025-12-26T10:30:00Z",
    "updated_at": "2025-12-26T10:30:00Z"
  }
}
```

### Example 2: Create category without display_order (auto-assign to last)
```bash
curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Beverages",
    "description": "Drinks and refreshments",
    "status": "active"
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 7,
    "restaurant_id": 1,
    "name": "Beverages",
    "description": "Drinks and refreshments",
    "display_order": 10,
    "status": "active",
    "created_at": "2025-12-26T10:35:00Z",
    "updated_at": "2025-12-26T10:35:00Z"
  }
}
```

### Example 3: Create category with specific restaurant_id
```bash
curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "restaurant_id": 5,
    "name": "Appetizers",
    "description": "Start your meal with our delicious starters",
    "status": "active"
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 8,
    "restaurant_id": 5,
    "name": "Appetizers",
    "description": "Start your meal with our delicious starters",
    "display_order": 1,
    "status": "active",
    "created_at": "2025-12-26T10:40:00Z",
    "updated_at": "2025-12-26T10:40:00Z"
  }
}
```

### Example 4: Create inactive category
```bash
curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Seasonal Specials",
    "description": "Limited time seasonal items",
    "status": "inactive"
  }'
```

### Example 5: Create minimal category (name and status only)
```bash
curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Side Dishes",
    "status": "active"
  }'
```

**Validation Rules:**
- `name`: Required, max 50 characters
- `description`: Optional, text
- `display_order`: Optional (auto-assigned if not provided), minimum 0
- `status`: Required, must be "active" or "inactive"
- `restaurant_id`: Optional (defaults to 1 if not provided)

**Error Response (400 Bad Request):**
```json
{
  "code": 400,
  "message": "Validation failed",
  "data": null
}
```

---

## 3. GET /api/admin/menu/categories/:id - Get category detail

### Example 1: Get category by ID
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories/1"
```

**Response (200 OK):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "name": "Appetizers",
    "description": "Start your meal with our delicious starters",
    "item_count": 8,
    "is_active": true,
    "display_order": 1
  }
}
```

### Example 2: Get category with no items
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories/5"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 5,
    "name": "New Category",
    "description": "Recently added category",
    "item_count": 0,
    "is_active": true,
    "display_order": 6
  }
}
```

**Error Response (404 Not Found):**
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories/999"
```

```json
{
  "code": 404,
  "message": "Category not found",
  "data": null
}
```

**Error Response (400 Bad Request - Invalid ID):**
```bash
curl -X GET "http://localhost:8080/api/admin/menu/categories/abc"
```

```json
{
  "code": 400,
  "message": "Invalid category ID",
  "data": null
}
```

---

## 4. PUT /api/admin/menu/categories/:id - Update category

### Example 1: Update all fields
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/categories/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Appetizers & Starters",
    "description": "Delicious starters to begin your dining experience",
    "display_order": 1,
    "status": "active"
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
    "name": "Appetizers & Starters",
    "description": "Delicious starters to begin your dining experience",
    "display_order": 1,
    "status": "active",
    "created_at": "2025-12-26T10:00:00Z",
    "updated_at": "2025-12-26T11:15:00Z"
  }
}
```

### Example 2: Update only name
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/categories/2" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Main Dishes"
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 2,
    "restaurant_id": 1,
    "name": "Main Dishes",
    "description": "Hearty and satisfying main dishes",
    "display_order": 2,
    "status": "active",
    "created_at": "2025-12-26T10:00:00Z",
    "updated_at": "2025-12-26T11:20:00Z"
  }
}
```

### Example 3: Update description only
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/categories/3" \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Sweet and delightful desserts"
  }'
```

### Example 4: Update display order
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/categories/4" \
  -H "Content-Type: application/json" \
  -d '{
    "display_order": 10
  }'
```

### Example 5: Update status to inactive
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/categories/5" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

**Validation Rules:**
- All fields are optional
- `name`: Max 50 characters if provided
- `display_order`: Minimum 0 if provided
- `status`: Must be "active" or "inactive" if provided

**Error Response (404 Not Found):**
```json
{
  "code": 404,
  "message": "Category not found",
  "data": null
}
```

---

## 5. PATCH /api/admin/menu/categories/:id/status - Toggle category status

### Example 1: Set category to inactive
```bash
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/1/status" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
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
    "name": "Appetizers",
    "description": "Start your meal with our delicious starters",
    "display_order": 1,
    "status": "inactive",
    "created_at": "2025-12-26T10:00:00Z",
    "updated_at": "2025-12-26T11:30:00Z"
  }
}
```

### Example 2: Set category to active
```bash
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/1/status" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": true
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "restaurant_id": 1,
    "name": "Appetizers",
    "description": "Start your meal with our delicious starters",
    "display_order": 1,
    "status": "active",
    "created_at": "2025-12-26T10:00:00Z",
    "updated_at": "2025-12-26T11:35:00Z"
  }
}
```

### Example 3: Deactivate multiple categories
```bash
# Deactivate category 3
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/3/status" \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'

# Deactivate category 5
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/5/status" \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'
```

**Error Response (404 Not Found):**
```bash
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/999/status" \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'
```

```json
{
  "code": 404,
  "message": "Category not found",
  "data": null
}
```

**Error Response (400 Bad Request):**
```bash
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/1/status" \
  -H "Content-Type: application/json" \
  -d '{}'
```

```json
{
  "code": 400,
  "message": "is_active field is required",
  "data": null
}
```

---

## Common Error Responses

### 400 Bad Request - Validation Error
```json
{
  "code": 400,
  "message": "Validation failed: name is required",
  "data": null
}
```

### 404 Not Found - Category Not Found
```json
{
  "code": 404,
  "message": "Category not found",
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

## Testing Workflow

### Complete Test Sequence
```bash
# 1. Create new categories
curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{"name": "Appetizers", "description": "Starters", "status": "active"}'

curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{"name": "Main Course", "description": "Main dishes", "status": "active"}'

curl -X POST "http://localhost:8080/api/admin/menu/categories" \
  -H "Content-Type: application/json" \
  -d '{"name": "Desserts", "description": "Sweets", "status": "active"}'

# 2. List all categories
curl -X GET "http://localhost:8080/api/admin/menu/categories"

# 3. Get category detail
curl -X GET "http://localhost:8080/api/admin/menu/categories/1"

# 4. Update category
curl -X PUT "http://localhost:8080/api/admin/menu/categories/1" \
  -H "Content-Type: application/json" \
  -d '{"name": "Appetizers & Starters", "description": "Updated description"}'

# 5. Toggle status to inactive
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/1/status" \
  -H "Content-Type: application/json" \
  -d '{"is_active": false}'

# 6. Verify status change
curl -X GET "http://localhost:8080/api/admin/menu/categories?status=inactive"

# 7. Toggle back to active
curl -X PATCH "http://localhost:8080/api/admin/menu/categories/1/status" \
  -H "Content-Type: application/json" \
  -d '{"is_active": true}'

# 8. Search categories
curl -X GET "http://localhost:8080/api/admin/menu/categories?search=app"

# 9. Sort by name
curl -X GET "http://localhost:8080/api/admin/menu/categories?sort=name"
```

---

## Notes

- **Restaurant ID**: If `restaurant_id` is not provided in POST request, it defaults to `1`
- **Display Order**: If `display_order` is not provided in POST request, it's automatically assigned to be the last position (max + 1)
- **Status Field**: API uses boolean `is_active` in responses but string `status` ("active"/"inactive") in database and requests
- **Item Count**: The `item_count` field shows the number of menu items in each category (only non-deleted items)
- **Pagination**: Default page size is 20, maximum is 100
- **Search**: Case-insensitive search on category name
- **Sorting Options**: `displayOrder` (default), `name`, `itemCount`
