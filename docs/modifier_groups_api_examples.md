# Modifier Groups API Documentation

## Overview
Modifier Groups API cho phép quản lý các nhóm tùy chỉnh (như Size, Topping, etc.) cho Menu Items.

---

## Curl Example

### Get All Modifier Groups (with pagination)
```bash
curl -X GET "http://localhost:8080/api/admin/menu/modifier-groups?page=1&page_size=10" \
  -H "Content-Type: application/json"
```

### Search Modifier Groups
```bash
curl -X GET "http://localhost:8080/api/admin/menu/modifier-groups?page=1&page_size=10&search=size" \
  -H "Content-Type: application/json"
```

### Filter by Status
```bash
curl -X GET "http://localhost:8080/api/admin/menu/modifier-groups?page=1&page_size=10&status=active" \
  -H "Content-Type: application/json"
```

### Create Modifier Group
```bash
curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Size",
    "selection_type": "single",
    "is_required": true,
    "status": "active"
  }'
```

### Update Modifier Group
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/modifier-groups/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Size Updated",
    "selection_type": "single",
    "is_required": true,
    "status": "active"
  }'
```

### Delete Modifier Group
```bash
curl -X DELETE "http://localhost:8080/api/admin/menu/modifier-groups/1" \
  -H "Content-Type: application/json"
```

---

## Curl Example (Detailed)

### Get All Modifier Groups
```bash
curl -X GET "http://localhost:8080/api/admin/menu/modifier-groups?page=1&page_size=10" \
  -H "Content-Type: application/json"
```

### Get Modifier Group by ID
```bash
curl -X GET "http://localhost:8080/api/admin/menu/modifier-groups/1" \
  -H "Content-Type: application/json"
```

---

## Endpoints

### 1. Get All Modifier Groups
**GET** `/api/admin/menu/modifier-groups`

**Description:** Lấy danh sách tất cả modifier groups (có phân trang)

**Query Parameters:**
- `page` (number, default: 1): Trang hiện tại
- `page_size` (number, default: 10): Số lượng items trên mỗi trang
- `search` (string, optional): Tìm kiếm theo tên
- `status` (string, optional): Lọc theo status (active, inactive)
- `selection_type` (string, optional): Lọc theo selection_type (single, multiple)
- `sort` (string, optional): Sắp xếp (name, display_order)

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Size",
      "description": "Pizza size options",
      "display_order": 1,
      "is_required": true,
      "allow_multiple_selection": false,
      "created_at": "2025-12-19T10:00:00Z",
      "updated_at": "2025-12-19T10:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Toppings",
      "description": "Extra toppings",
      "display_order": 2,
      "is_required": false,
      "allow_multiple_selection": true,
      "created_at": "2025-12-19T10:05:00Z",
      "updated_at": "2025-12-19T10:05:00Z"
    }
  ]
}
```

---

### 2. Create Modifier Group
**POST** `/api/admin/menu/modifier-groups`

**Description:** Tạo modifier group mới

**Request Body:**
```json
{
  "name": "Size",
  "description": "Pizza size options",
  "display_order": 1,
  "is_required": true,
  "allow_multiple_selection": false
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Size",
    "description": "Pizza size options",
    "display_order": 1,
    "is_required": true,
    "allow_multiple_selection": false,
    "created_at": "2025-12-19T10:00:00Z",
    "updated_at": "2025-12-19T10:00:00Z"
  }
}
```

**Error Response:**
```json
{
  "code": 1,
  "message": "error",
  "data": "Modifier group name is required"
}
```

---

### 3. Update Modifier Group
**PUT** `/api/admin/menu/modifier-groups/:id`

**Description:** Cập nhật thông tin modifier group

**URL Parameters:**
- `id` (string, required): ID của modifier group

**Request Body:**
```json
{
  "name": "Size Updated",
  "description": "Pizza size options updated",
  "display_order": 2,
  "is_required": false,
  "allow_multiple_selection": true
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Size Updated",
    "description": "Pizza size options updated",
    "display_order": 2,
    "is_required": false,
    "allow_multiple_selection": true,
    "created_at": "2025-12-19T10:00:00Z",
    "updated_at": "2025-12-19T10:15:00Z"
  }
}
```

---

### 4. Delete Modifier Group
**DELETE** `/api/admin/menu/modifier-groups/:id`

**Description:** Xóa modifier group

**URL Parameters:**
- `id` (string, required): ID của modifier group

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": ""
}
```

**Error Response:**
```json
{
  "code": 1,
  "message": "error",
  "data": "Modifier group not found"
}
```

---

### 5. Create Modifier Options
**POST** `/api/admin/menu/modifier-groups/:id/options`

**Description:** Thêm options vào modifier group (ví dụ: Small, Medium, Large cho Size)

**URL Parameters:**
- `id` (string, required): ID của modifier group

**Request Body:**
```json
{
  "name": "Small",
  "description": "Small size",
  "price_modifier": 0,
  "display_order": 1,
  "is_available": true
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "modifier_group_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Small",
    "description": "Small size",
    "price_modifier": 0,
    "display_order": 1,
    "is_available": true,
    "created_at": "2025-12-19T10:20:00Z",
    "updated_at": "2025-12-19T10:20:00Z"
  }
}
```

---

## Example Usage

### Create Size Modifier Group with Options
```bash
# 1. Create modifier group
curl -X POST http://localhost:8080/api/admin/menu/modifier-groups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Size",
    "description": "Pizza size options",
    "display_order": 1,
    "is_required": true,
    "allow_multiple_selection": false
  }'

# Response: Get modifier_group_id = "550e8400-e29b-41d4-a716-446655440000"

# 2. Add options to the group
curl -X POST http://localhost:8080/api/admin/menu/modifier-groups/550e8400-e29b-41d4-a716-446655440000/options \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small",
    "description": "Small size",
    "price_modifier": 0,
    "display_order": 1,
    "is_available": true
  }'

curl -X POST http://localhost:8080/api/admin/menu/modifier-groups/550e8400-e29b-41d4-a716-446655440000/options \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Medium",
    "description": "Medium size",
    "price_modifier": 2.50,
    "display_order": 2,
    "is_available": true
  }'

curl -X POST http://localhost:8080/api/admin/menu/modifier-groups/550e8400-e29b-41d4-a716-446655440000/options \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Large",
    "description": "Large size",
    "price_modifier": 5.00,
    "display_order": 3,
    "is_available": true
  }'

# 3. Get all modifier groups
curl -X GET http://localhost:8080/api/admin/menu/modifier-groups
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 0 | success | Request successful |
| 1 | error | General error |
| 400 | Bad Request | Invalid request data |
| 404 | Not Found | Modifier group not found |
| 500 | Internal Server Error | Server error |

---

## Notes
- Modifier groups là optional cho menu items
- `is_required`: Nếu true, customer phải chọn option từ group này
- `allow_multiple_selection`: Cho phép chọn nhiều options từ cùng 1 group
- `price_modifier`: Giá thêm cho option này (ví dụ: +2.50 cho Medium size)
