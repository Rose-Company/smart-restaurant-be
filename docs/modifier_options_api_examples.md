# Modifier Options API Documentation

## Overview
Modifier Options API cho phép quản lý các tùy chọn cụ thể trong mỗi Modifier Group (ví dụ: Small, Medium, Large cho Size).

---

## Quick Reference - Copy & Paste

### Create Modifier Option (Add to Group)
```bash
curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small (10\")",
    "price_adjustment": 0,
    "status": "active"
  }'
```

### Create Multiple Options
```bash
curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Medium (12\")",
    "price_adjustment": 2.50,
    "status": "active"
  }'

curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Large (14\")",
    "price_adjustment": 5.00,
    "status": "active"
  }'
```

### Update Modifier Option
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/modifier-options/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small (10\") Updated",
    "price_adjustment": 0.50,
    "status": "active"
  }'
```

### Delete Modifier Option
```bash
curl -X DELETE "http://localhost:8080/api/admin/menu/modifier-options/1" \
  -H "Content-Type: application/json"
```

---

## Curl Example (Detailed)

### Create Modifier Option (Add to Group)
```bash
curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small (10\")",
    "price_adjustment": 0,
    "status": "active"
  }'
```

### Update Modifier Option
```bash
curl -X PUT "http://localhost:8080/api/admin/menu/modifier-options/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small (10\")",
    "price_adjustment": 0,
    "status": "active"
  }'
```

### Delete Modifier Option
```bash
curl -X DELETE "http://localhost:8080/api/admin/menu/modifier-options/1" \
  -H "Content-Type: application/json"
```

---

## Endpoints

### 1. Create Modifier Option
**POST** `/api/admin/menu/modifier-groups/:id/options`

**Description:** Thêm option mới vào modifier group

**URL Parameters:**
- `id` (number, required): ID của modifier group

**Request Body:**
```json
{
  "name": "Small (10\")",
  "price_adjustment": 0,
  "status": "active"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "group_id": 1,
    "name": "Small (10\")",
    "price_adjustment": 0,
    "status": "active"
  }
}
```

---

### 2. Update Modifier Option
**PUT** `/api/admin/menu/modifier-options/:id`

**Description:** Cập nhật thông tin modifier option

**URL Parameters:**
- `id` (number, required): ID của modifier option

**Request Body:**
```json
{
  "name": "Small (10\") Updated",
  "price_adjustment": 0.50,
  "status": "active"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "group_id": 1,
    "name": "Small (10\") Updated",
    "price_adjustment": 0.50,
    "status": "active"
  }
}
```

---

### 3. Delete Modifier Option
**DELETE** `/api/admin/menu/modifier-options/:id`

**Description:** Xóa modifier option

**URL Parameters:**
- `id` (number, required): ID của modifier option

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": "Modifier Option deleted successfully"
}
```

**Error Response:**
```json
{
  "code": 1,
  "message": "error",
  "data": "Modifier option not found"
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
    "selection_type": "single",
    "is_required": true,
    "status": "active"
  }'

# Response will include: "id": 1 (for modifier group)

# 2. Add options to the group
curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small (10\")",
    "price_adjustment": 0,
    "status": "active"
  }'

curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Medium (12\")",
    "price_adjustment": 2.50,
    "status": "active"
  }'

curl -X POST "http://localhost:8080/api/admin/menu/modifier-groups/1/options" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Large (14\")",
    "price_adjustment": 5.00,
    "status": "active"
  }'

# 3. Update an option
curl -X PUT "http://localhost:8080/api/admin/menu/modifier-options/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Small (10\") - Limited",
    "price_adjustment": -0.50,
    "status": "active"
  }'

# 4. Delete an option
curl -X DELETE "http://localhost:8080/api/admin/menu/modifier-options/3"
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 0 | success | Request successful |
| 1 | error | General error |
| 400 | Bad Request | Invalid request data |
| 404 | Not Found | Modifier option not found |
| 500 | Internal Server Error | Server error |

---

## Notes
- Options phải được tạo thông qua endpoint POST `/api/admin/menu/modifier-groups/:id/options`
- `price_adjustment` có thể âm (discount) hoặc dương (upcharge)
- Status: `active` hoặc `inactive`
- Khi xóa option, kiểm tra xem có order nào đang sử dụng option này không
- Tất cả modifier option được sắp xếp theo thứ tự thêm vào (FIFO)
