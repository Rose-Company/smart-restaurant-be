# Tables API - Example Requests

## 1. GET /api/admin/tables - List all tables

### Example 1: Get all tables (default pagination)
```bash
curl -X GET "http://localhost:8080/api/admin/tables"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 12,
    "page": 1,
    "page_size": 20,
    "items": [
      {
        "id": 1,
        "table_number": "T-01",
        "capacity": 4,
        "location": "Main Hall",
        "status": "active"
      },
      {
        "id": 2,
        "table_number": "T-02",
        "capacity": 2,
        "location": "Main Hall",
        "status": "occupied",
        "order_data": {
          "active_orders": 3,
          "total_bill": 156.0
        }
      }
    ],
    "extra": null
  }
}
```

### Example 2: Get tables with pagination
```bash
curl -X GET "http://localhost:8080/api/admin/tables?page=1&page_size=10"
```

### Example 3: Search tables by table number
```bash
curl -X GET "http://localhost:8080/api/admin/tables?search=T-0"
```

### Example 4: Filter by status (occupied tables only)
```bash
curl -X GET "http://localhost:8080/api/admin/tables?status=occupied"
```

### Example 5: Filter by zone/location
```bash
curl -X GET "http://localhost:8080/api/admin/tables?zone=VIP"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "total": 2,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 15,
        "table_number": "T-05",
        "capacity": 8,
        "location": "VIP",
        "status": "active"
      },
      {
        "id": 19,
        "table_number": "VIP-02",
        "capacity": 6,
        "location": "VIP",
        "status": "occupied"
      }
    ],
    "extra": null
  }
}
```

### Example 6: Sort by table number
```bash
curl -X GET "http://localhost:8080/api/admin/tables?sort=table_number"
```

### Example 7: Sort by capacity (descending)
```bash
curl -X GET "http://localhost:8080/api/admin/tables?sort=capacity"
```

### Example 8: Sort by recently created
```bash
curl -X GET "http://localhost:8080/api/admin/tables?sort=recently_created"
```

### Example 9: Combined filters
```bash
curl -X GET "http://localhost:8080/api/admin/tables?zone=Main%20Hall&status=active&sort=tableNumber&page=1&page_size=5"
```

---

## 2. GET /api/admin/tables/:id - Get single table details

### Example 1: Get table by ID
```bash
curl -X GET "http://localhost:8080/api/admin/tables/1"
```

**Response (Active table):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "table_number": "T-01",
    "capacity": 4,
    "location": "Main Hall",
    "status": "active"
  }
}
```

### Example 2: Get occupied table (with order data)
```bash
curl -X GET "http://localhost:8080/api/admin/tables/2"
```

**Response (Occupied table):**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 2,
    "table_number": "T-02",
    "capacity": 2,
    "location": "Main Hall",
    "status": "occupied",
    "order_data": {
      "active_orders": 3,
      "total_bill": 156.0
    }
  }
}
```

### Example 3: Table not found
```bash
curl -X GET "http://localhost:8080/api/admin/tables/999"
```

**Response:**
```json
{
  "code": 404,
  "message": "record not found",
  "data": null
}
```

---

## 3. POST /api/admin/tables - Create new table

### Example 1: Create a new table in Main Hall
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "T-20",
    "capacity": 4,
    "location": "Main Hall",
    "status": "active"
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 20,
    "table_number": "T-20",
    "capacity": 4,
    "location": "Main Hall",
    "status": "active",
    "created_at": "2025-12-19T10:30:00Z",
    "updated_at": "2025-12-19T10:30:00Z"
  }
}
```

### Example 2: Create a VIP table
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "VIP-10",
    "capacity": 8,
    "location": "VIP",
    "status": "active"
  }'
```

### Example 3: Create a Patio table
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "P-15",
    "capacity": 6,
    "location": "Patio",
    "status": "active"
  }'
```

### Example 4: Error - Duplicate table number
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "T-01",
    "capacity": 4,
    "location": "Main Hall",
    "status": "active"
  }'
```

**Response:**
```json
{
  "code": 400,
  "message": "table number already exists",
  "data": null
}
```

### Example 5: Error - Invalid status
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "T-25",
    "capacity": 4,
    "location": "Main Hall",
    "status": "invalid_status"
  }'
```

**Response:**
```json
{
  "code": 400,
  "message": "Key: 'CreateTableRequest.Status' Error:Field validation for 'Status' failed on the 'oneof' tag",
  "data": null
}
```

### Example 6: Error - Missing required fields
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "capacity": 4,
    "location": "Main Hall"
  }'
```

**Response:**
```json
{
  "code": 400,
  "message": "Key: 'CreateTableRequest.TableNumber' Error:Field validation for 'TableNumber' failed on the 'required' tag",
  "data": null
}
```

---

## 4. PUT /api/admin/tables/:id - Update table

### Example 1: Update all fields
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/1" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "T-01-NEW",
    "capacity": 6,
    "location": "VIP",
    "status": "active"
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "table_number": "T-01-NEW",
    "capacity": 6,
    "location": "VIP",
    "status": "active",
    "created_at": "2025-12-19T08:00:00Z",
    "updated_at": "2025-12-19T10:45:00Z"
  }
}
```

### Example 2: Update only table number
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/2" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "T-02-UPDATED"
  }'
```

### Example 3: Update only capacity
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/3" \
  -H "Content-Type: application/json" \
  -d '{
    "capacity": 8
  }'
```

### Example 4: Update only location
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/4" \
  -H "Content-Type: application/json" \
  -d '{
    "location": "Patio"
  }'
```

### Example 5: Update only status
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/5" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

### Example 6: Update multiple fields
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/6" \
  -H "Content-Type: application/json" \
  -d '{
    "capacity": 10,
    "location": "VIP",
    "status": "active"
  }'
```

### Example 7: Error - Duplicate table number
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/7" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "T-01"
  }'
```

**Response:**
```json
{
  "code": 400,
  "message": "table number already exists",
  "data": null
}
```

### Example 8: Error - Table not found
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/999" \
  -H "Content-Type: application/json" \
  -d '{
    "capacity": 6
  }'
```

**Response:**
```json
{
  "code": 404,
  "message": "record not found",
  "data": null
}
```

---

## 5. PATCH /api/admin/tables/:id/status - Update table status

### Example 1: Mark table as occupied
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/1/status" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "occupied"
  }'
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "id": 1,
    "table_number": "T-01",
    "capacity": 4,
    "location": "Main Hall",
    "status": "occupied",
    "created_at": "2025-12-19T08:00:00Z",
    "updated_at": "2025-12-19T11:00:00Z"
  }
}
```

### Example 2: Mark table as active (available)
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/2/status" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "active"
  }'
```

### Example 3: Mark table as inactive (maintenance/closed)
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/3/status" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

### Example 4: Error - Invalid status value
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/4/status" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "invalid"
  }'
```

**Response:**
```json
{
  "code": 400,
  "message": "Key: 'UpdateTableStatusRequest.Status' Error:Field validation for 'Status' failed on the 'oneof' tag",
  "data": null
}
```

### Example 5: Error - Missing status field
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/5/status" \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Response:**
```json
{
  "code": 400,
  "message": "Key: 'UpdateTableStatusRequest.Status' Error:Field validation for 'Status' failed on the 'required' tag",
  "data": null
}
```

---

## Status Values Reference

- `active` - Table is available and ready for use
- `occupied` - Table is currently being used by customers
- `inactive` - Table is not available (maintenance, closed, etc.)

## Common Error Responses

### 400 Bad Request
```json
{
  "code": 400,
  "message": "validation error message",
  "data": null
}
```

### 404 Not Found
```json
{
  "code": 404,
  "message": "record not found",
  "data": null
}
```

### 500 Internal Server Error
```json
{
  "code": 500,
  "message": "internal server error",
  "data": null
}
```

---

## Testing with Postman

You can import these examples into Postman:

1. Create a new collection named "Tables API"
2. Add environment variable: `baseUrl = http://localhost:8080`
3. Use `{{baseUrl}}/api/admin/tables` in your requests

## Testing Workflow

### 1. Initial Setup
```bash
# Get all tables to see current state
curl -X GET "http://localhost:8080/api/admin/tables"
```

### 2. Create a new table
```bash
curl -X POST "http://localhost:8080/api/admin/tables" \
  -H "Content-Type: application/json" \
  -d '{
    "table_number": "TEST-01",
    "capacity": 4,
    "location": "Test Area",
    "status": "active"
  }'
```

### 3. Get the created table details
```bash
curl -X GET "http://localhost:8080/api/admin/tables/[ID_FROM_STEP_2]"
```

### 4. Update the table
```bash
curl -X PUT "http://localhost:8080/api/admin/tables/[ID_FROM_STEP_2]" \
  -H "Content-Type: application/json" \
  -d '{
    "capacity": 6,
    "location": "VIP"
  }'
```

### 5. Change table status to occupied
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/[ID_FROM_STEP_2]/status" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "occupied"
  }'
```

### 6. Filter occupied tables
```bash
curl -X GET "http://localhost:8080/api/admin/tables?status=occupied"
```

### 7. Change back to active
```bash
curl -X PATCH "http://localhost:8080/api/admin/tables/[ID_FROM_STEP_2]/status" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "active"
  }'
```
