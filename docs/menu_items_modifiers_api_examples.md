# Menu Items Modifier Groups API Documentation

## Overview
API này cho phép gán Modifier Groups vào Menu Items và quản lý modifier relationships.

---

## Curl Example

### Assign Modifier Group to Menu Item
```bash
curl -X POST "http://localhost:8080/api/admin/menu/items/1/modifier-groups" \
  -H "Content-Type: application/json" \
  -d '{
    "modifier_group_id": "1"
  }'
```

### Delete Modifier Group from Menu Item
```bash
curl -X DELETE "http://localhost:8080/api/admin/menu/items/1/modifier-groups/1" \
  -H "Content-Type: application/json"
```

---

## Endpoints

### 1. Assign Modifier Group to Menu Item
**POST** `/api/menu/items/:id/modifier-groups`

**Description:** Gán một modifier group vào menu item. Khi assign thành công, customer sẽ có thể chọn options từ group này khi order.

**URL Parameters:**
- `id` (string, required): ID của menu item

**Request Body:**
```json
{
  "modifier_group_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response:**
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440100",
    "menu_item_id": "550e8400-e29b-41d4-a716-446655440001",
    "modifier_group_id": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2025-12-19T10:30:00Z"
  }
}
```

**Error Response:**
```json
{
  "code": 1,
  "message": "error",
  "data": "Menu item not found"
}
```

---

### 2. Delete Modifier Group from Menu Item
**DELETE** `/api/menu/items/:id/modifier-groups/:groupId`

**Description:** Xóa gán modifier group từ menu item

**URL Parameters:**
- `id` (string, required): ID của menu item
- `groupId` (string, required): ID của modifier group

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
  "data": "Modifier group assignment not found"
}
```

---



When customer orders pizza, request body would look like:

```json
{
  "items": [
    {
      "menu_item_id": "550e8400-e29b-41d4-a716-446655440001",
      "quantity": 2,
      "modifiers": [
        {
          "modifier_group_id": "550e8400-e29b-41d4-a716-446655440000",
          "selected_option_id": "550e8400-e29b-41d4-a716-446655440002",
          "selected_option_name": "Medium (12\")",
          "price_modifier": 2.50
        },
        {
          "modifier_group_id": "550e8400-e29b-41d4-a716-446655440003",
          "selected_option_ids": [
            "550e8400-e29b-41d4-a716-446655440010",
            "550e8400-e29b-41d4-a716-446655440011"
          ],
          "selected_option_names": [
            "Pepperoni",
            "Mushroom"
          ],
          "price_modifier": 2.50
        }
      ]
    }
  ]
}
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 0 | success | Request successful |
| 1 | error | General error |
| 400 | Bad Request | Invalid request data |
| 404 | Not Found | Menu item or modifier group not found |
| 409 | Conflict | Modifier group already assigned |
| 500 | Internal Server Error | Server error |

---

## Business Rules

- Một menu item có thể có nhiều modifier groups
- Một modifier group có thể được gán vào nhiều menu items
- Khi delete modifier group assignment, không ảnh hưởng đến modifier group (có thể được dùng cho items khác)
- Khi customer order, phải select options cho tất cả required modifier groups
- Có thể select multiple options từ một group nếu `allow_multiple_selection: true`
- `price_modifier` cộng vào giá menu item

---

## Notes
- POST endpoint nằm trong `/api/menu/items` (public)
- DELETE endpoint nằm trong `/api/menu/items` (public)
- Quản lý modifier groups nằm trong `/api/admin/menu/modifier-groups` (admin only)
