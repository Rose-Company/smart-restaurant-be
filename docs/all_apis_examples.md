# 📚 Complete API Examples - TASK-001 to TASK-030

Complete documentation of all Smart Restaurant APIs with request/response examples.

---

## 📋 Table of Contents

1. **TASK-001 to TASK-009**: Order Management APIs
2. **TASK-010 to TASK-012**: Bill Management APIs
3. **TASK-013 to TASK-014**: Payment Management APIs
4. **TASK-015**: Discount Management APIs
5. **TASK-016 to TASK-020**: Customer Profile APIs
6. **TASK-021 to TASK-027**: Staff Management APIs
7. **TASK-028 to TASK-030**: Staff Profile APIs

---

# ORDER MANAGEMENT APIs (TASK-001 to TASK-009)

## TASK-001: Create Order
**POST** `/api/orders`  
**Authentication:** Optional (Guest or Authenticated)

### Request
```bash
curl -X POST "http://localhost:8080/api/orders" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "table_id": 5,
    "items": [
      {
        "menu_item_id": 12,
        "quantity": 2,
        "modifiers": [
          {
            "modifier_group_id": 3,
            "modifier_option_id": 8,
            "additional_price": 0
          }
        ],
        "notes": "No onions"
      },
      {
        "menu_item_id": 15,
        "quantity": 1,
        "modifiers": [],
        "notes": ""
      }
    ],
    "customer_notes": "Extra sauce on the side",
    "dining_mode": "in-restaurant"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Order created successfully",
  "data": {
    "id": 45,
    "table_id": 5,
    "order_number": "ORD-2025-00045",
    "status": "pending",
    "total_amount": 350000,
    "items_count": 2,
    "items": [
      {
        "id": 89,
        "menu_item_id": 12,
        "quantity": 2,
        "unit_price": 150000,
        "modifiers": [
          {
            "id": 201,
            "modifier_group_id": 3,
            "modifier_option_id": 8,
            "name": "Extra cheese",
            "additional_price": 0
          }
        ],
        "notes": "No onions",
        "status": "pending"
      }
    ],
    "customer_notes": "Extra sauce on the side",
    "created_at": "2025-01-19T10:30:45Z",
    "estimated_ready_time": "2025-01-19T10:45:00Z"
  }
}
```

---

## TASK-002: Get Orders List
**GET** `/api/orders`  
**Authentication:** Optional

### Request
```bash
curl -X GET "http://localhost:8080/api/orders?page=1&page_size=10&status=pending&table_id=5&category=Main%20Courses&sort=created_at.desc" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Query Parameters
- `page` (int): Page number (default: 1)
- `page_size` (int): Items per page (default: 10)
- `status` (string): Filter by status (pending, confirmed, preparing, ready, completed, cancelled)
- `table_id` (int): Filter by table
- `category` (string): Filter by menu category (e.g., "Main Courses", "Appetizers", "Desserts")
- `sort` (string): Sort order (created_at_asc, created_at.desc, total_amount_asc, total_amount_desc)

### Response (Success)
```json
{
  "code": 200,
  "message": "Orders retrieved successfully",
  "data": {
    "total": 25,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 45,
        "table_id": 5,
        "order_number": "ORD-2025-00045",
        "table_name": "Table 5",
        "customer_id": "user-123",
        "customer_name": "John Doe",
        "status": "preparing",
        "total_amount": 350000,
        "items_count": 2,
        "items": [
          {
            "id": 89,
            "menu_item_id": 12,
            "menu_item_name": "Grilled Salmon",
            "menu_item_image": "https://cdn.example.com/menu/salmon.jpg",
            "quantity": 2,
            "unit_price": 150000,
            "subtotal": 300000,
            "special_instructions": "No onions",
            "status": "preparing",
            "modifiers": [
              {
                "modifier_group_id": 3,
                "modifier_group_name": "Sauce",
                "modifier_option_id": 8,
                "modifier_option_name": "Garlic butter",
                "price": 0
              }
            ]
          },
          {
            "id": 90,
            "menu_item_id": 15,
            "menu_item_name": "Caesar Salad",
            "menu_item_image": "https://cdn.example.com/menu/caesar-salad.jpg",
            "quantity": 1,
            "unit_price": 50000,
            "subtotal": 50000,
            "special_instructions": null,
            "status": "ready",
            "modifiers": []
          }
        ],
        "created_at": "2025-01-19T10:30:45Z",
        "updated_at": "2025-01-19T10:40:00Z",
        "estimated_ready_time": "2025-01-19T10:45:00Z",
        "waiter_id": "staff-001",
        "waiter_name": "John Waiter"
      }
    ]
  }
}
```

---

## TASK-003: Get Order Details
**GET** `/api/orders/:id`  
**Authentication:** Optional

### Request
```bash
curl -X GET "http://localhost:8080/api/orders/45" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Order details retrieved successfully",
  "data": {
    "id": 45,
    "table_id": 5,
    "order_number": "ORD-2025-00045",
    "status": "preparing",
    "total_amount": 350000,
    "items": [
      {
        "id": 89,
        "menu_item_id": 12,
        "name": "Grilled Salmon",
        "quantity": 2,
        "unit_price": 150000,
        "subtotal": 300000,
        "modifiers": [
          {
            "id": 201,
            "name": "Extra cheese",
            "additional_price": 0
          }
        ],
        "notes": "No onions",
        "status": "preparing",
        "prepared_at": null
      }
    ],
    "timeline": [
      {
        "id": 1,
        "status": "pending",
        "changed_at": "2025-01-19T10:30:45Z",
        "changed_by": "Customer"
      },
      {
        "id": 2,
        "status": "confirmed",
        "changed_at": "2025-01-19T10:31:00Z",
        "changed_by": "Waiter"
      }
    ],
    "waiter_id": "user-123",
    "customer_notes": "Extra sauce on the side",
    "created_at": "2025-01-19T10:30:45Z",
    "estimated_ready_time": "2025-01-19T10:45:00Z"
  }
}
```

---

## TASK-004: Update Order Status
**PATCH** `/api/orders/:id/status`  
**Authentication:** Required (Admin/Kitchen Staff/Waiter)

### Request
```bash
curl --location --request PATCH 'http://localhost:8080/api/orders/1/status' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer YOUR_TOKEN' \
  --data '{
    "status": "confirmed",
    "updated_by": "waiter",
    "reason": "Order confirmed and sent to kitchen"
  }'
```

### Allowed Status Transitions
- pending → confirmed
- confirmed → preparing
- preparing → ready
- ready → completed
- any → cancelled

### Response (Success)
```json
{
  "code": 200,
  "message": "Order status updated successfully",
  "data": {
    "id": 45,
    "order_number": "ORD-2025-00045",
    "status": "confirmed",
    "previous_status": "pending",
    "updated_at": "2025-01-19T10:31:00Z",
    "updated_by": "waiter",
    "updated_by_name": "John Waiter",
    "estimated_ready_time": "2025-01-19T10:45:00Z"
  }
}
```

---

## TASK-005: Update Order Item Status (Multiple Orders)
**PATCH** `/api/orders/items/:itemId/status`  
**Authentication:** Required (Kitchen Staff)

### Request
```bash
curl --location --request PATCH 'http://localhost:8080/api/orders/items/89/status' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer YOUR_TOKEN' \
  --data '{
    "order_ids": [45, 48, 52],
    "status": "ready"
  }'
```

### URL Parameters
- `:itemId` - The menu item ID to mark as ready (required in URL path)

### Request Body
- `order_ids` - Array of order IDs where this item needs to be marked as ready
- `status` - New status (ready, completed, etc.)

### Note
Since the same menu item can appear in multiple orders (different tables), this endpoint allows you to mark a specific item as ready across multiple orders at once. You only need the item ID in the URL and provide the list of order IDs in the request body.

### Response (Success)
```json
{
  "code": 200,
  "message": "Order item status updated successfully",
  "data": {
    "total_updated": 3,
    "updated_items": [
      {
        "item_id": 89,
        "order_id": 45,
        "menu_item_name": "Caesar Salad",
        "status": "ready",
        "previous_status": "preparing",
        "updated_at": "2025-01-19T10:43:00Z",
        "updated_by": "kitchen",
        "updated_by_name": "Chef John"
      },
      {
        "item_id": 89,
        "order_id": 48,
        "menu_item_name": "Caesar Salad",
        "status": "ready",
        "previous_status": "preparing",
        "updated_at": "2025-01-19T10:43:00Z",
        "updated_by": "kitchen",
        "updated_by_name": "Chef John"
      },
      {
        "item_id": 89,
        "order_id": 52,
        "menu_item_name": "Caesar Salad",
        "status": "ready",
        "previous_status": "preparing",
        "updated_at": "2025-01-19T10:43:00Z",
        "updated_by": "kitchen",
        "updated_by_name": "Chef John"
      }
    ]
  }
}
```

---

## TASK-005b: Update Multiple Items Status (Single Order)
**PATCH** `/api/orders/:id/items/status`  
**Authentication:** Required (Kitchen Staff)

### Request
```bash
curl --location --request PATCH 'http://localhost:8080/api/orders/45/items/status' \
  --header 'Content-Type: application/json' \
  --header 'Authorization: Bearer YOUR_TOKEN' \
  --data '{
    "items": [
      {
        "menu_item_id": 89,
        "status": "ready"
      },
      {
        "menu_item_id": 92,
        "status": "ready"
      },
      {
        "menu_item_id": 105,
        "status": "completed"
      }
    ]
  }'
```

### URL Parameters
- `:id` - The order ID (required in URL path)

### Request Body
- `items` - Array of items with:
  - `menu_item_id` - The menu item ID
  - `status` - New status for this item

### Note
This endpoint updates multiple different items in a single order. Use this when you want to mark several items as ready/completed in one specific order, unlike TASK-005 which updates the same item across multiple orders.

### Response (Success)
```json
{
  "code": 200,
  "message": "Order items status updated successfully",
  "data": {
    "total_updated": 3,
    "updated_items": [
      {
        "item_id": 156,
        "order_id": 45,
        "menu_item_name": "Caesar Salad",
        "status": "ready",
        "previous_status": "preparing",
        "updated_at": "2025-01-19T10:43:00Z",
        "updated_by": "kitchen",
        "updated_by_name": "Chef John"
      },
      {
        "item_id": 157,
        "order_id": 45,
        "menu_item_name": "French Onion Soup",
        "status": "ready",
        "previous_status": "preparing",
        "updated_at": "2025-01-19T10:43:00Z",
        "updated_by": "kitchen",
        "updated_by_name": "Chef John"
      },
      {
        "item_id": 158,
        "order_id": 45,
        "menu_item_name": "Lobster Bisque",
        "status": "completed",
        "previous_status": "ready",
        "updated_at": "2025-01-19T10:43:00Z",
        "updated_by": "kitchen",
        "updated_by_name": "Chef John"
      }
    ]
  }
}
```

---

## TASK-006: Update Order (Add Notes/Metadata) ==> ĐANG LỖI
**PATCH** `/api/orders/:id`   
**Authentication:** Optional

### Request
```bash
curl -X PATCH "http://localhost:8080/api/orders/45" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "customer_notes": "Please make it quick, we are in a hurry",
    "special_requests": "No spicy, extra lemon",
    "dining_mode": "in-restaurant"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Order updated successfully",
  "data": {
    "id": 45,
    "customer_notes": "Please make it quick, we are in a hurry",
    "special_requests": "No spicy, extra lemon",
    "updated_at": "2025-01-19T10:35:00Z"
  }
}
```

---

## TASK-007: Cancel Order
**POST** `/api/orders/:id/cancel`  
**Authentication:** Optional

### Request
```bash
curl -X POST "http://localhost:8080/api/orders/45/cancel" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "reason": "Changed our mind",
    "refund_amount": 350000
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Order cancelled successfully",
  "data": {
    "id": 45,
    "status": "cancelled",
    "reason": "Changed our mind",
    "refund_amount": 350000,
    "cancelled_at": "2025-01-19T10:35:00Z"
  }
}
```

---

## TASK-008: Send Kitchen Alert ==> chưa có màn FE
**POST** `/api/orders/:id/alert`  
**Authentication:** Required (Waiter)

### Request
```bash
curl -X POST "http://localhost:8080/api/orders/45/alert" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "alert_type": "urgent",
    "message": "Customer at table 5 is waiting for their order",
    "priority": "high"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Kitchen alert sent successfully",
  "data": {
    "alert_id": 156,
    "order_id": 45,
    "alert_type": "urgent",
    "message": "Customer at table 5 is waiting for their order",
    "sent_at": "2025-01-19T10:35:00Z",
    "acknowledged_at": null
  }
}
```

---

## TASK-009: Create Order Review ==> chưa ráp được nha
**POST** `/api/orders/:id/review`   
**Authentication:** Required (Customer)

### Request
```bash
curl -X POST "http://localhost:8080/api/orders/45/review" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "rating": 5,
    "comment": "Excellent food and service!",
    "food_quality": 5,
    "service_quality": 4,
    "value_for_money": 5
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Review submitted successfully",
  "data": {
    "review_id": 89,
    "order_id": 45,
    "rating": 5,
    "comment": "Excellent food and service!",
    "food_quality": 5,
    "service_quality": 4,
    "value_for_money": 5,
    "created_at": "2025-01-19T10:45:00Z"
  }
}
```

---

# BILL MANAGEMENT APIs (TASK-010 to TASK-012)

## TASK-010: Create Bill from Order
**POST** `/api/bills`  
**Authentication:** Optional

### Request
```bash
curl -X POST "http://localhost:8080/api/bills" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "order_id": 45,
    "discount_code": "SAVE10",
    "tax_percentage": 10
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Bill created successfully",
  "data": {
    "id": 67,
    "bill_number": "BILL-2025-00067",
    "order_id": 45,
    "subtotal": 350000,
    "tax_amount": 35000,
    "discount_amount": 35000,
    "total_amount": 350000,
    "items": [
      {
        "item_id": 89,
        "name": "Grilled Salmon",
        "quantity": 2,
        "unit_price": 150000,
        "subtotal": 300000
      }
    ],
    "discount_code": "SAVE10",
    "discount_type": "percentage",
    "discount_value": 10,
    "status": "unpaid",
    "created_at": "2025-01-19T10:45:00Z"
  }
}
```

---

## TASK-011: Get Bill Details
**GET** `/api/bills/:id`  
**Authentication:** Optional

### Request
```bash
curl -X GET "http://localhost:8080/api/bills/67" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Bill details retrieved successfully",
  "data": {
    "id": 67,
    "bill_number": "BILL-2025-00067",
    "order_id": 45,
    "table_id": 5,
    "subtotal": 350000,
    "tax_amount": 35000,
    "tax_percentage": 10,
    "discount_amount": 35000,
    "discount_code": "SAVE10",
    "service_charge": 0,
    "total_amount": 350000,
    "items": [
      {
        "item_id": 89,
        "name": "Grilled Salmon",
        "quantity": 2,
        "unit_price": 150000,
        "subtotal": 300000,
        "tax": 30000
      }
    ],
    "status": "unpaid",
    "payment_status": "pending",
    "created_at": "2025-01-19T10:45:00Z",
    "due_date": "2025-01-19T11:00:00Z"
  }
}
```

---

## TASK-012: Update Bill (Add Discount, Mark Paid)
**PATCH** `/api/bills/:id`  
**Authentication:** Optional

### Request - Add Discount
```bash
curl -X PATCH "http://localhost:8080/api/bills/67" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "discount_code": "WELCOME20",
    "discount_type": "percentage",
    "discount_value": 20,
    "add_service_charge": true,
    "service_charge_percentage": 5
  }'
```

### Request - Mark Paid
```bash
curl -X PATCH "http://localhost:8080/api/bills/67" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "status": "paid",
    "payment_method": "cash",
    "paid_amount": 350000
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Bill updated successfully",
  "data": {
    "id": 67,
    "subtotal": 350000,
    "discount_amount": 70000,
    "tax_amount": 28000,
    "service_charge": 17500,
    "total_amount": 325500,
    "status": "paid",
    "payment_method": "cash",
    "paid_at": "2025-01-19T10:50:00Z"
  }
}
```

---

# PAYMENT MANAGEMENT APIs (TASK-013 to TASK-014)

## TASK-013: Process Payment
**POST** `/api/payments`  
**Authentication:** Optional

### Request
```bash
curl -X POST "http://localhost:8080/api/payments" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "bill_id": 67,
    "payment_method": "credit_card",
    "amount": 325500,
    "card_number": "4111111111111111",
    "card_holder": "John Doe",
    "card_exp_month": 12,
    "card_exp_year": 2026,
    "card_cvv": "123"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Payment processed successfully",
  "data": {
    "payment_id": 234,
    "bill_id": 67,
    "transaction_id": "TXN-2025-001234",
    "payment_method": "credit_card",
    "amount": 325500,
    "status": "success",
    "processed_at": "2025-01-19T10:51:00Z",
    "reference_code": "REF-2025-001234"
  }
}
```

---

## TASK-014: Check Payment Status
**GET** `/api/payments/:id/status`  
**Authentication:** Optional

### Request
```bash
curl -X GET "http://localhost:8080/api/payments/234/status" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Payment status retrieved successfully",
  "data": {
    "payment_id": 234,
    "bill_id": 67,
    "transaction_id": "TXN-2025-001234",
    "amount": 325500,
    "status": "success",
    "payment_method": "credit_card",
    "processed_at": "2025-01-19T10:51:00Z",
    "receipt_url": "https://receipt.example.com/REF-2025-001234"
  }
}
```

---

# DISCOUNT MANAGEMENT APIs (TASK-015)

## TASK-015: Validate Discount Code
**POST** `/api/discounts/validate`  
**Authentication:** Optional

### Request
```bash
curl -X POST "http://localhost:8080/api/discounts/validate" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "code": "SAVE10",
    "order_total": 350000,
    "table_id": 5
  }'
```

### Response (Success - Valid Code)
```json
{
  "code": 200,
  "message": "Discount code is valid",
  "data": {
    "discount_code": "SAVE10",
    "discount_type": "percentage",
    "discount_value": 10,
    "discount_amount": 35000,
    "min_order_amount": 100000,
    "max_discount_amount": 100000,
    "valid": true,
    "expires_at": "2025-12-31T23:59:59Z",
    "usage_count": 15,
    "max_usage": 100
  }
}
```

### Response (Error - Invalid Code)
```json
{
  "code": 400,
  "message": "Discount code is invalid or expired",
  "data": {
    "discount_code": "INVALID123",
    "valid": false,
    "reason": "Code not found"
  }
}
```

---

# CUSTOMER PROFILE APIs (TASK-016 to TASK-020)

## TASK-016: Get Customer Profile
**GET** `/api/customer/profile`  
**Authentication:** Required (Logged-in Customer)

### Request
```bash
curl -X GET "http://localhost:8080/api/customer/profile" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "user-123",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+84912345678",
    "avatar_url": "https://cdn.example.com/avatars/user-123.jpg",
    "street_address": "123 Main St",
    "city": "Ho Chi Minh",
    "state": "HCMC",
    "postal_code": "70000",
    "country": "Vietnam",
    "total_orders": 25,
    "loyalty_points": 2500,
    "member_since": "2024-01-15T00:00:00Z",
    "preferred_language": "vi",
    "status": "active"
  }
}
```

---

## TASK-017: Update Customer Profile
**PUT** `/api/customer/profile`  
**Authentication:** Required

### Request
```bash
curl -X PUT "http://localhost:8080/api/customer/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+84912345678",
    "street_address": "456 Oak Ave",
    "city": "Da Nang",
    "state": "DN",
    "postal_code": "50000",
    "country": "Vietnam"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": "user-123",
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "+84912345678",
    "street_address": "456 Oak Ave",
    "city": "Da Nang",
    "state": "DN",
    "postal_code": "50000",
    "country": "Vietnam",
    "updated_at": "2025-01-19T11:00:00Z"
  }
}
```

---

## TASK-018: Upload Avatar
**POST** `/api/customer/avatar`  
**Authentication:** Required

### Request (Form Data)
```bash
curl -X POST "http://localhost:8080/api/customer/avatar" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "avatar=@/path/to/avatar.jpg"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Avatar uploaded successfully",
  "data": {
    "id": "user-123",
    "avatar_url": "https://cdn.example.com/avatars/user-123-new.jpg",
    "avatar_updated_at": "2025-01-19T11:05:00Z"
  }
}
```

---

## TASK-019: Change Password
**PATCH** `/api/customer/password`  
**Authentication:** Required

### Request
```bash
curl -X PATCH "http://localhost:8080/api/customer/password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "old_password": "OldPassword123!",
    "new_password": "NewPassword456!",
    "confirm_password": "NewPassword456!"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Password changed successfully",
  "data": {
    "id": "user-123",
    "password_changed_at": "2025-01-19T11:10:00Z"
  }
}
```

---

## TASK-020: Get Customer Reviews
**GET** `/api/customer/reviews`  
**Authentication:** Required

### Request
```bash
curl -X GET "http://localhost:8080/api/customer/reviews?page=1&page_size=10" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Reviews retrieved successfully",
  "data": {
    "total": 5,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": 89,
        "order_id": 45,
        "rating": 5,
        "comment": "Excellent food and service!",
        "food_quality": 5,
        "service_quality": 4,
        "value_for_money": 5,
        "created_at": "2025-01-19T10:45:00Z"
      }
    ]
  }
}
```

---

# STAFF MANAGEMENT APIs (TASK-021 to TASK-027)

## TASK-021: List Staff
**GET** `/api/admin/staff`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X GET "http://localhost:8080/api/admin/staff?page=1&page_size=10&role=waiter&status=active&search=John" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Query Parameters
- `page` (int): Page number
- `page_size` (int): Items per page
- `role` (string): Filter by role (waiter, kitchen_staff, manager, cashier)
- `status` (string): Filter by status (active, inactive)
- `search` (string): Search by name or email

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff list retrieved successfully",
  "data": {
    "total": 15,
    "page": 1,
    "page_size": 10,
    "items": [
      {
        "id": "staff-001",
        "name": "John Waiter",
        "email": "john@restaurant.com",
        "phone": "+84912345678",
        "role": "waiter",
        "status": "active",
        "assigned_tables": [1, 2, 3, 4, 5],
        "avatar_url": "https://cdn.example.com/avatars/staff-001.jpg",
        "created_at": "2024-06-15T00:00:00Z",
        "last_login": "2025-01-19T10:00:00Z"
      }
    ]
  }
}
```

---

## TASK-022: Create Staff Account
**POST** `/api/admin/staff`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X POST "http://localhost:8080/api/admin/staff" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Jane Kitchen",
    "email": "jane@restaurant.com",
    "phone": "+84987654321",
    "password": "SecurePassword123!",
    "role": "kitchen_staff",
    "status": "active",
    "department": "Kitchen",
    "position": "Sous Chef",
    "hire_date": "2025-01-15",
    "employment_status": "full_time",
    "shift_type": "morning",
    "weekly_hours": 40,
    "assigned_tables": [1, 2]
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff account created successfully",
  "data": {
    "id": "staff-002",
    "name": "Jane Kitchen",
    "email": "jane@restaurant.com",
    "phone": "+84987654321",
    "role": "kitchen_staff",
    "status": "active",
    "created_at": "2025-01-19T11:15:00Z"
  }
}
```

---

## TASK-023: Get Staff Details
**GET** `/api/admin/staff/:id`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X GET "http://localhost:8080/api/admin/staff/staff-001" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff details retrieved successfully",
  "data": {
    "id": "staff-001",
    "name": "John Waiter",
    "email": "john@restaurant.com",
    "phone": "+84912345678",
    "role": "waiter",
    "status": "active",
    "avatar_url": "https://cdn.example.com/avatars/staff-001.jpg",
    "assigned_tables": [1, 2, 3, 4, 5],
    "assigned_tables_details": [
      {
        "id": 1,
        "name": "Table 1",
        "capacity": 4,
        "status": "occupied"
      }
    ],
    "statistics": {
      "total_orders_served": 245,
      "average_rating": 4.8,
      "orders_today": 12,
      "active_orders": 3
    },
    "created_at": "2024-06-15T00:00:00Z",
    "last_login": "2025-01-19T10:00:00Z"
  }
}
```

---

## TASK-024: Update Staff Account
**PUT** `/api/admin/staff/:id`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X PUT "http://localhost:8080/api/admin/staff/staff-001" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "John Senior Waiter",
    "email": "john.senior@restaurant.com",
    "phone": "+84912345679",
    "role": "waiter",
    "status": "active",
    "department": "Service",
    "position": "Senior Waiter",
    "shift_type": "afternoon",
    "weekly_hours": 45,
    "assigned_tables": [1, 2, 3, 4, 5, 6]
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff account updated successfully",
  "data": {
    "id": "staff-001",
    "name": "John Senior Waiter",
    "email": "john.senior@restaurant.com",
    "updated_at": "2025-01-19T11:20:00Z"
  }
}
```

---

## TASK-025: Delete Staff Account
**DELETE** `/api/admin/staff/:id`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X DELETE "http://localhost:8080/api/admin/staff/staff-002" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff account deleted successfully",
  "data": {
    "id": "staff-002",
    "status": "inactive",
    "deleted_at": "2025-01-19T11:25:00Z"
  }
}
```

---

## TASK-026: Send Staff Invitation
**POST** `/api/admin/staff/:id/send-invite`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X POST "http://localhost:8080/api/admin/staff/staff-003/send-invite" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "email": "newstaff@restaurant.com",
    "role": "kitchen_staff",
    "message": "Welcome to our restaurant team!"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff invitation sent successfully",
  "data": {
    "invitation_id": "inv-001",
    "email": "newstaff@restaurant.com",
    "role": "kitchen_staff",
    "token": "invite_xxxxx_1234567890",
    "expires_at": "2025-01-26T11:30:00Z",
    "sent_at": "2025-01-19T11:30:00Z"
  }
}
```

---

## TASK-027: Assign Tables to Waiter
**PATCH** `/api/admin/staff/:id/assign-tables`  
**Authentication:** Required (Admin)

### Request
```bash
curl -X PATCH "http://localhost:8080/api/admin/staff/staff-001/assign-tables" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "table_ids": [1, 2, 3, 7, 8]
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Tables assigned to waiter successfully",
  "data": {
    "waiter_id": "staff-001",
    "assigned_tables": [1, 2, 3, 7, 8],
    "previous_tables": [1, 2, 3, 4, 5],
    "assigned_at": "2025-01-19T11:35:00Z"
  }
}
```

---

# STAFF PROFILE APIs (TASK-028 to TASK-030)

## TASK-028: Get Staff Profile
**GET** `/api/staff/profile`  
**Authentication:** Required (Logged-in Staff)

### Request
```bash
curl -X GET "http://localhost:8080/api/staff/profile" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Staff profile retrieved successfully",
  "data": {
    "id": "staff-001",
    "name": "John Waiter",
    "email": "john@restaurant.com",
    "phone": "+84912345678",
    "role": "waiter",
    "status": "active",
    "avatar_url": "https://cdn.example.com/avatars/staff-001.jpg",
    "assigned_tables": [1, 2, 3, 4, 5],
    "assigned_tables_names": ["Table 1", "Table 2", "Table 3", "Table 4", "Table 5"],
    "statistics": {
      "total_orders_served": 245,
      "average_rating": 4.8,
      "orders_today": 12,
      "active_orders": 3
    },
    "department": "Service",
    "position": "Waiter",
    "hire_date": "2024-06-15",
    "employment_status": "full_time",
    "shift_type": "morning",
    "weekly_hours": 40,
    "created_at": "2024-06-15T00:00:00Z",
    "last_login": "2025-01-19T10:00:00Z"
  }
}
```

---

## TASK-029: Update Own Profile
**PUT** `/api/staff/profile`  
**Authentication:** Required

### Request
```bash
curl -X PUT "http://localhost:8080/api/staff/profile" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "John Waiter Updated",
    "email": "john.updated@restaurant.com",
    "phone": "+84912345679",
    "avatar": "https://cdn.example.com/avatars/staff-001-new.jpg"
  }'
```

### Response (Success)
```json
{
  "code": 200,
  "message": "Profile updated successfully",
  "data": {
    "id": "staff-001",
    "name": "John Waiter Updated",
    "email": "john.updated@restaurant.com",
    "phone": "+84912345679",
    "avatar_url": "https://cdn.example.com/avatars/staff-001-new.jpg",
    "updated_at": "2025-01-19T11:40:00Z"
  }
}
```

---

## TASK-030: Change Password
**PATCH** `/api/staff/password`  
**Authentication:** Required

### Request
```bash
curl -X PATCH "http://localhost:8080/api/staff/password" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "old_password": "OldPassword123!",
    "new_password": "NewPassword456!",
    "confirm_password": "NewPassword456!"
  }'
```

### Password Requirements
- Minimum 8 characters
- Must contain uppercase letters
- Must contain lowercase letters
- Must contain numbers
- Must contain special characters (!@#$%^&*()_+-=[]{}|;:,.<>?)

### Response (Success)
```json
{
  "code": 200,
  "message": "Password changed successfully",
  "data": {
    "id": "staff-001",
    "password_changed_at": "2025-01-19T11:45:00Z"
  }
}
```

---

# Common Error Responses

## 400 Bad Request
```json
{
  "code": 400,
  "message": "Invalid input",
  "data": {
    "field": "email",
    "error": "Email format is invalid"
  }
}
```

## 401 Unauthorized
```json
{
  "code": 401,
  "message": "Not authorized",
  "data": {
    "reason": "Invalid or expired token"
  }
}
```

## 403 Forbidden
```json
{
  "code": 403,
  "message": "Access denied",
  "data": {
    "reason": "Insufficient permissions"
  }
}
```

## 404 Not Found
```json
{
  "code": 404,
  "message": "Resource not found",
  "data": {
    "resource": "Order",
    "id": 999
  }
}
```

## 500 Internal Server Error
```json
{
  "code": 500,
  "message": "Internal server error",
  "data": {
    "error_id": "ERR-2025-001234"
  }
}
```

---

# Authentication

All APIs use JWT Bearer Token authentication (except optional endpoints).

## Format
```
Authorization: Bearer <your_jwt_token>
```

## Example
```bash
curl -X GET "http://localhost:8080/api/customer/profile" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

# Rate Limiting

- **Rate Limit**: 100 requests per minute
- **Headers**: 
  - `X-RateLimit-Limit`: 100
  - `X-RateLimit-Remaining`: 99
  - `X-RateLimit-Reset`: Unix timestamp

---

# Pagination

All list endpoints support pagination with these parameters:
- `page` (default: 1)
- `page_size` (default: 10, max: 100)

Response format:
```json
{
  "total": 50,
  "page": 1,
  "page_size": 10,
  "items": []
}
```

---

# Contact

For API support, contact: api-support@restaurant.example.com
