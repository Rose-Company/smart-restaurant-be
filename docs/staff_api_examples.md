# Staff API - Example Requests

## GET /api/staff/tables - Get Tables for Staff View

### Overview
Get tables with active orders for waiter/kitchen staff view. Includes detailed order and item information with filtering by `is_ready_to_bill` and `is_help_needed` flags.

**Endpoint:** `/api/staff/tables`  
**Method:** GET  
**Authentication:** Required (Bearer Token)  
**Base URL:** `http://localhost:8080/api/staff/tables`

---

## Query Parameters

| Parameter | Type | Optional | Description | Example |
|-----------|------|----------|-------------|---------|
| `page` | int | ✅ | Page number (default: 1) | 1 |
| `page_size` | int | ✅ | Items per page (default: 10) | 5 |
| `is_ready_to_bill` | bool | ✅ | Filter by billing status | true |
| `is_help_needed` | bool | ✅ | Filter by help status | false |

---

## Examples

### 1️⃣ Get all occupied tables with active orders
```bash
curl -X GET "http://localhost:8080/api/staff/tables" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Tables retrieved successfully",
  "data": {
    "total": 2,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 5,
        "table_number": "Table 05",
        "capacity": 4,
        "location": "Main Hall",
        "status": "occupied",
        "orders": [
          {
            "id": 3,
            "order_number": "ORD-2025-00003",
            "status": "preparing",
            "total_amount": 250000,
            "is_ready_to_bill": false,
            "is_help_needed": false,
            "items_count": 2,
            "items": [
              {
                "id": 7,
                "item_name": "Burger",
                "quantity": 2,
                "unit_price": 100000,
                "status": "completed"
              },
              {
                "id": 8,
                "item_name": "Fries",
                "quantity": 1,
                "unit_price": 50000,
                "status": "preparing"
              }
            ],
            "created_at": "2026-01-19T14:30:00Z",
            "customer_name": "John Doe"
          }
        ],
        "active_orders_count": 1,
        "total_bill": 250000,
        "created_at": "2026-01-19T14:25:00Z",
        "updated_at": "2026-01-19T14:35:00Z"
      },
      {
        "id": 8,
        "table_number": "Table 08",
        "capacity": 2,
        "location": "Window Seat",
        "status": "occupied",
        "orders": [
          {
            "id": 4,
            "order_number": "ORD-2025-00004",
            "status": "ready",
            "total_amount": 180000,
            "is_ready_to_bill": true,
            "is_help_needed": false,
            "items_count": 1,
            "items": [
              {
                "id": 9,
                "item_name": "Salmon",
                "quantity": 1,
                "unit_price": 150000,
                "status": "completed"
              }
            ],
            "created_at": "2026-01-19T14:20:00Z",
            "customer_name": "Jane Smith"
          }
        ],
        "active_orders_count": 1,
        "total_bill": 180000,
        "created_at": "2026-01-19T14:15:00Z",
        "updated_at": "2026-01-19T14:40:00Z"
      }
    ]
  }
}
```

---

### 2️⃣ Filter: Tables ready for billing
```bash
curl -X GET "http://localhost:8080/api/staff/tables?is_ready_to_bill=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Purpose:** Shows tables with at least one order where `is_ready_to_bill=true`  
**Use Case:** Waiter checks which tables need to be billed

---

### 3️⃣ Filter: Tables that need help
```bash
curl -X GET "http://localhost:8080/api/staff/tables?is_help_needed=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Purpose:** Shows tables with at least one order where `is_help_needed=true`  
**Use Case:** Kitchen/Waiter checks which tables need assistance

---

### 4️⃣ Filter: Kitchen-ready tables (both flags false)
```bash
curl -X GET "http://localhost:8080/api/staff/tables?is_ready_to_bill=false&is_help_needed=false" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Purpose:** Shows tables where ALL orders have BOTH flags as false  
**Use Case:** Kitchen checks which tables are ready to serve

---

### 5️⃣ Get with pagination
```bash
curl -X GET "http://localhost:8080/api/staff/tables?page=1&page_size=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:** First 5 tables

---

### 6️⃣ Combine multiple filters
```bash
curl -X GET "http://localhost:8080/api/staff/tables?page=2&page_size=10&is_ready_to_bill=true" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Filter Logic

| Filter Params | Result | Use Case |
|---------------|--------|----------|
| None | All occupied tables with active orders | General view |
| `is_ready_to_bill=true` | Tables with ≥1 order ready for billing | Billing screen |
| `is_help_needed=true` | Tables with ≥1 order needing help | Help requests |
| `is_ready_to_bill=false&is_help_needed=false` | Tables where all orders are kitchen-ready | Kitchen view |

---

## Response Fields Explanation

### Table Object
- `id` - Table ID
- `table_number` - Display number (e.g., "Table 05")
- `capacity` - Number of seats
- `location` - Zone/area
- `status` - Table status (occupied)
- `orders` - Array of active orders on this table
- `active_orders_count` - Count of non-completed orders
- `total_bill` - Total amount for all active orders
- `created_at` - Table creation time
- `updated_at` - Last update time

### Order Object (nested in table)
- `id` - Order ID
- `order_number` - Unique order number (ORD-2025-XXXXX)
- `status` - Order status (pending, preparing, ready, completed, etc.)
- `total_amount` - Total order amount (VND)
- `is_ready_to_bill` - Flag indicating order is ready to be billed
- `is_help_needed` - Flag indicating customer needs help
- `items_count` - Number of items in order
- `items` - Array of order items
- `created_at` - Order creation time
- `customer_name` - Customer name for this order

### Order Item Object (nested in order)
- `id` - Order item ID
- `item_name` - Menu item name
- `quantity` - Quantity ordered
- `unit_price` - Price per unit
- `status` - Item status (pending, preparing, completed, cancelled)

---

## Error Cases

### 401 Unauthorized
```json
{
  "code": 401,
  "error_code": "unauthorized",
  "message": "Unauthorized access",
  "data": null
}
```

### 403 Forbidden
```json
{
  "code": 403,
  "error_code": "forbidden",
  "message": "You don't have permission to access this resource",
  "data": null
}
```

### 400 Bad Request
```json
{
  "code": 400,
  "error_code": "bad_request",
  "message": "Invalid query parameters",
  "data": null
}
```

---

## Use Case Examples

### 👨‍💼 Waiter workflow:
1. Get all tables: `GET /api/staff/tables`
2. Filter billing-ready: `GET /api/staff/tables?is_ready_to_bill=true`
3. Mark table as served, then process payment

### 👨‍🍳 Kitchen workflow:
1. Get all tables: `GET /api/staff/tables`
2. Filter kitchen-ready: `GET /api/staff/tables?is_ready_to_bill=false&is_help_needed=false`
3. Serve items to tables

### 🚨 Help requests:
1. Get help-needed tables: `GET /api/staff/tables?is_help_needed=true`
2. Address customer issues
3. Clear help flag via order update API

---

## Rate Limiting & Performance

- ✅ Paginated responses (default 10 per page)
- ✅ Optimized queries to avoid N+1 problems
- ✅ Real-time data (no caching)
- ⏱️ Typical response time: < 500ms

---

## Related APIs

- **Update Order:** `PATCH /api/orders/:id` - Update order notes/metadata
- **Update Order Status:** `PATCH /api/orders/:id/status` - Change order status
- **Get Order Details:** `GET /api/orders/:id` - Get full order information
- **Update Item Status:** `PATCH /api/orders/:id/items/status` - Mark items as done
- **Get Table Detail:** `GET /api/staff/tables/:id` - Get table detail with all order items

---

# GET /api/staff/tables/:id - Get Table Detail for Staff

### Overview
Get detailed view of a specific table with all order items (completed, preparing, cancelled) grouped by order ID. Shows complete order item information including which order each item belongs to.

**Endpoint:** `/api/staff/tables/:id`  
**Method:** GET  
**Authentication:** Required (Bearer Token)  
**Base URL:** `http://localhost:8080/api/staff/tables/:id`

---

## URI Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `id` | int | Table ID | 5 |

---

## Examples

### 1️⃣ Get table detail with all items
```bash
curl -X GET "http://localhost:8080/api/staff/tables/5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response (200 OK):**
```json
{
  "code": 200,
  "message": "Table detail retrieved successfully",
  "data": {
    "id": 5,
    "table_number": "05",
    "capacity": 4,
    "location": "Main Hall",
    "status": "occupied",
    "guest_count": 4,
    "total_bill": 450000,
    "all_orders_count": 2,
    "order_items": [
      {
        "id": 7,
        "order_id": 10,
        "item_name": "Grilled Salmon",
        "quantity": 1,
        "unit_price": 150000,
        "status": "completed"
      },
      {
        "id": 8,
        "order_id": 10,
        "item_name": "Caesar Salad",
        "quantity": 2,
        "unit_price": 50000,
        "status": "completed"
      },
      {
        "id": 9,
        "order_id": 11,
        "item_name": "Burger",
        "quantity": 2,
        "unit_price": 75000,
        "status": "serving"
      },
      {
        "id": 10,
        "order_id": 11,
        "item_name": "Fries",
        "quantity": 1,
        "unit_price": 50000,
        "status": "serving"
      }
    ],
    "created_at": "2026-01-19T14:25:00Z"
  }
}
```

---

### 2️⃣ Get table detail - Kitchen view
```bash
curl -X GET "http://localhost:8080/api/staff/tables/8" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Purpose:** Shows all items on a table for kitchen staff to verify what's being prepared  
**Use Case:** Kitchen checks what items are needed for table 8

---

## Response Fields Explanation

### Table Detail Object
- `id` - Table ID
- `table_number` - Display number (e.g., "05", "08")
- `capacity` - Number of seats
- `location` - Zone/area (Main Hall, Window Seat, etc.)
- `status` - Table status (occupied, active, inactive)
- `guest_count` - Table capacity (number of guests can be seated)
- `total_bill` - Sum of all orders' totals (all orders, not just active)
- `all_orders_count` - Total number of orders placed on this table (includes completed/cancelled)
- `order_items` - Array of ALL order items (all statuses)
- `created_at` - Table creation time

### Order Item Detail Object
- `id` - Order item unique ID
- `order_id` - Which order this item belongs to (KEY FIELD for grouping)
- `item_name` - Menu item name
- `quantity` - Quantity ordered
- `unit_price` - Price per unit (VND)
- `status` - Item status (pending, preparing, completed, serving, cancelled)

---

## Item Status Reference

| Status | Meaning | Color |
|--------|---------|-------|
| `pending` | Waiting to be prepared | ⚪ Gray |
| `preparing` | Currently being prepared | 🟠 Orange |
| `completed` | Ready to be served | 🟢 Green |
| `serving` | Being served to customer | 🔵 Blue |
| `cancelled` | Item was cancelled | ❌ Red |

---

## Data Retrieval Strategy

**This endpoint optimizes queries:**
- ✅ Gets table data once
- ✅ Gets all orders for table in ONE query
- ✅ Gets all order items in ONE query
- ✅ Groups items by order_id in memory
- ⏱️ Typical response: < 200ms

**No N+1 queries** - All data loaded efficiently regardless of order/item count

---

## Use Case Examples

### 🍳 Kitchen displays order items for table:
```bash
# Kitchen staff clicks on table 5
GET /api/staff/tables/5

# Returns all items (by order_id):
# Order 10: Salmon (1) + Salad (2)
# Order 11: Burger (2) + Fries (1)
```

### 👨‍💼 Waiter verifies items before serving:
```bash
# Waiter wants to verify items on table before billing
GET /api/staff/tables/5

# Reviews all items with order_id
# Confirms which items were ordered together
```

### 📋 Manager checks table history:
```bash
# Manager reviews what was ordered on table 5
GET /api/staff/tables/5

# Sees all items including cancelled ones
# Calculates total spent
```

---

## Error Cases

### 401 Unauthorized
```json
{
  "code": 401,
  "error_code": "unauthorized",
  "message": "Unauthorized access",
  "data": null
}
```

### 404 Not Found
```json
{
  "code": 404,
  "error_code": "table_not_found",
  "message": "Table not found",
  "data": null
}
```

### 400 Bad Request
```json
{
  "code": 400,
  "error_code": "invalid_table_id",
  "message": "Invalid table ID format",
  "data": null
}
```

---

## Comparison: `/api/staff/tables` vs `/api/staff/tables/:id`

| Feature | `/api/staff/tables` | `/api/staff/tables/:id` |
|---------|-------------------|------------------------|
| Purpose | List occupied tables | Show one table detail |
| Data | Orders grouped by table | All items with order_id |
| Filtering | By status flags | No filtering |
| Use Case | Browse all tables | View specific table |
| Response | Multiple tables, paginated | Single table, all items |
| Items Structure | Nested in orders | Flat list with order_id |

---

## Related APIs

- **List Tables:** `GET /api/staff/tables` - List all occupied tables with active orders
- **Update Order:** `PATCH /api/orders/:id` - Update order notes/metadata
- **Update Item Status:** `PATCH /api/orders/:id/items/status` - Mark items as done
- **Get Order Details:** `GET /api/orders/:id` - Get full order information

