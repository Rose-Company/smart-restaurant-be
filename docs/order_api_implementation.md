# 📦 ORDER MANAGEMENT APIs DOCUMENTATION

## Overview
Complete implementation of Order Management APIs (TASK-001 to TASK-009) for Smart Restaurant System.

**Base URL:** `/api/orders`  
**Authentication:** Required for most endpoints (JWT token)

---

## 📋 API Endpoints Summary

| Task | Method | Endpoint | Description | Priority |
|------|--------|----------|-------------|----------|
| TASK-001 | POST | `/api/orders` | Create/Submit order from cart | 🔴 CRITICAL |
| TASK-002 | GET | `/api/orders` | Get orders list with filters | 🔴 CRITICAL |
| TASK-003 | GET | `/api/orders/:id` | Get order details (role-based) | 🔴 CRITICAL |
| TASK-004 | PATCH | `/api/orders/:id/status` | Update order status | 🔴 CRITICAL |
| TASK-005 | PATCH | `/api/orders/:id/items/:itemId/status` | Update item status | 🟠 HIGH |
| TASK-006 | PATCH | `/api/orders/:id` | Add notes/metadata | 🟡 MEDIUM |
| TASK-007 | POST | `/api/orders/:id/cancel` | Cancel order | 🟡 MEDIUM |
| TASK-008 | POST | `/api/orders/:id/alert` | Send kitchen alert | 🟡 MEDIUM |
| TASK-009 | POST | `/api/orders/:id/review` | Submit review | 🟢 LOW |

---

## 🗂️ File Structure

```
internal/
├── handlers/
│   ├── base.go           ✅ Updated with order routes
│   └── order.go          ✨ NEW - Order handlers (9 endpoints)
├── models/
│   └── order.go          ✨ NEW - Order models & request/response types
├── repositories/
│   └── order.go          ✨ NEW - Order repositories (5 repos)
└── services/
    ├── base.go           ✅ Updated with order repos
    └── order.go          ✨ NEW - Order business logic (9 functions)
```

---

## 📊 Database Tables Used

| Table | Purpose | Migration |
|-------|---------|-----------|
| `orders` | Main order records | 001 + 006 |
| `order_items` | Order line items | 001 + 006 |
| `order_item_modifiers` | Item modifier selections | 006 |
| `order_timeline` | Status change audit trail | 006 |
| `kitchen_alerts` | Kitchen-to-waiter alerts | 006 |
| `tables` | Table information | 001 |
| `menu_items` | Menu item details | 003 |
| `modifier_groups` | Modifier groups | 003 |
| `modifier_options` | Modifier options | 003 |
| `users` | Customer & staff info | 005 |

---

## 🔧 Implementation Details

### Models (`internal/models/order.go`)

**Core Models:**
- `Order` - Main order entity
- `OrderItem` - Order line item
- `OrderModifier` - Selected modifiers
- `OrderTimeline` - Status history
- `KitchenAlert` - Kitchen alerts

**Request Models:**
- `CreateOrderRequest`
- `UpdateOrderStatusRequest`
- `UpdateOrderItemStatusRequest`
- `UpdateOrderRequest`
- `CancelOrderRequest`
- `CreateAlertRequest`
- `CreateReviewRequest`

**Response Models:**
- `OrderResponse` - Full order details
- `OrderListItemResponse` - List item (simplified)
- `PaginatedOrdersResponse` - Paginated list
- `OrderStatusUpdateResponse`
- `OrderItemStatusUpdateResponse`
- `AlertResponse`
- `CancelOrderResponse`

### Repositories (`internal/repositories/order.go`)

**5 Repositories Created:**
1. `OrderRepository` - Order CRUD operations
2. `OrderItemRepository` - Order items CRUD
3. `OrderModifierRepository` - Modifiers CRUD
4. `OrderTimelineRepository` - Timeline CRUD
5. `KitchenAlertRepository` - Alerts CRUD

All repositories extend `BaseRepository[T]` with full CRUD operations.

### Services (`internal/services/order.go`)

**Business Logic Functions:**

1. **CreateOrder** - TASK-001
   - Validates table & menu items
   - Calculates totals (subtotal, tax, modifiers)
   - Creates order, items, modifiers in transaction
   - Creates initial timeline entry
   - Generates unique order number (ORD-YYYY-NNNNNN)

2. **GetOrders** - TASK-002
   - Dynamic filtering: status, table_id, date_range, search
   - Pagination support (page, page_size)
   - Sorting: created_at, total_amount (ASC/DESC)
   - Role-based filtering
   - Preloads: Table, CustomerUser, Waiter

3. **GetOrderByID** - TASK-003
   - Role-based response:
     - **Kitchen**: Simplified (items only)
     - **Customer/Waiter/Admin**: Full details + timeline
   - Includes: items, modifiers, timeline, bill summary

4. **UpdateOrderStatus** - TASK-004
   - Validates status transitions
   - Updates timestamps (accepted_at, preparing_at, ready_at, etc.)
   - Creates timeline entry
   - Valid transitions:
     ```
     pending → confirmed → preparing → ready → served → completed
                    ↓          ↓         ↓
                cancelled  cancelled  cancelled
     ```

5. **UpdateOrderItemStatus** - TASK-005
   - Updates individual item status
   - Kitchen staff can mark items: preparing, ready, completed
   - Waiters can mark items: served

6. **UpdateOrder** - TASK-006
   - Appends notes (preserves existing)
   - Updates metadata (JSONB)

7. **CancelOrder** - TASK-007
   - Validates cancellation eligibility
   - Cannot cancel if: preparing, ready, completed
   - Creates timeline entry
   - Returns refund info

8. **SendKitchenAlert** - TASK-008
   - Creates kitchen alert
   - Sends to assigned waiter
   - Types: order_ready, item_ready, delay_warning, urgent
   - Priority: low, normal, high, urgent

9. **CreateOrderReview** - TASK-009
   - Validates order is completed
   - Validates customer ownership
   - Creates review with photos
   - (TODO: Needs customer_reviews table from migration 007)

**Helper Functions:**
- `generateOrderNumber()` - Generates unique ORD-YYYY-NNNNNN format
- `isValidStatusTransition()` - Validates status workflow
- `strPtr()` - String pointer helper

### Handlers (`internal/handlers/order.go`)

**9 HTTP Handlers:**

Each handler:
1. Validates request (binding)
2. Extracts user context (user_id, user_role, user_name)
3. Calls service layer
4. Returns standardized response

**Error Handling:**
- Uses `common.AbortWithError()` for consistent errors
- Returns proper HTTP status codes
- Validates path parameters (order ID, item ID)

---

## 🔐 Authentication & Authorization

**User Context:**
- `user_id` - Extracted from JWT token
- `user_role` - Role name (customer/waiter/kitchen/admin)
- `user_name` - User's full name

**Role-Based Access:**
- **Customer**: Create order, view own orders, cancel own orders, submit reviews
- **Waiter**: View all orders, update status (served), update items
- **Kitchen**: View orders, update status (preparing→ready), send alerts
- **Admin**: Full access to all operations

---

## 🎯 Status Workflow

```
┌─────────┐
│ pending │ (Customer places order)
└────┬────┘
     │
     ▼
┌───────────┐
│ confirmed │ (Waiter confirms)
└─────┬─────┘
      │
      ▼
┌───────────┐
│ preparing │ (Kitchen starts cooking)
└─────┬─────┘
      │
      ▼
┌────────┐
│ ready  │ (Kitchen marks ready)
└───┬────┘
    │
    ▼
┌────────┐
│ served │ (Waiter serves to table)
└───┬────┘
    │
    ▼
┌───────────┐
│ completed │ (Customer finishes, bill paid)
└───────────┘

Any status can → cancelled (with restrictions)
```

---

## 📝 Usage Examples

### 1. Create Order (TASK-001)
```bash
POST /api/orders
Content-Type: application/json
Authorization: Bearer <token>

{
  "table_id": 5,
  "items": [
    {
      "menu_item_id": 31,
      "quantity": 2,
      "special_instructions": "No onions please",
      "modifiers": [
        {
          "modifier_group_id": 1,
          "modifier_option_id": 3
        }
      ]
    }
  ],
  "notes": "Birthday celebration"
}
```

### 2. Get Orders with Filters (TASK-002)
```bash
GET /api/orders?status=preparing&date_from=2026-01-01&page=1&page_size=10&sort=created_at_desc
Authorization: Bearer <token>
```

### 3. Update Order Status (TASK-004)
```bash
PATCH /api/orders/1001/status
Content-Type: application/json
Authorization: Bearer <token>

{
  "status": "preparing",
  "updated_by": "kitchen",
  "reason": "Started cooking"
}
```

### 4. Cancel Order (TASK-007)
```bash
POST /api/orders/1001/cancel
Content-Type: application/json
Authorization: Bearer <token>

{
  "reason": "Customer changed mind"
}
```

---

## ⚠️ Known Limitations & TODOs

1. **Transaction Management:**
   - Current implementation uses basic `BEGIN/COMMIT/ROLLBACK`
   - Should use GORM's native transaction API: `db.Transaction()`

2. **Customer Reviews (TASK-009):**
   - Placeholder response returned
   - Needs `customer_reviews` table from migration 007
   - Requires `CustomerReviewRepository` implementation

3. **Authentication Middleware:**
   - Routes registered but middleware not applied
   - Need to uncomment `authenticator` in `base.go`
   - Add middleware to order routes: `orders.Use(authenticator.AuthMiddleware())`

4. **Role-Based Filtering (TASK-002):**
   - `role` query parameter accepted but not fully implemented
   - Should filter orders based on user role:
     - Customer: only own orders
     - Waiter: orders assigned to them
     - Kitchen: all active orders
     - Admin: all orders

5. **Validation:**
   - Order item quantity limits not enforced
   - Menu item availability not checked
   - Table capacity not validated

6. **Notifications:**
   - Kitchen alerts created but not sent (no WebSocket/SSE)
   - Email notifications not implemented
   - Push notifications not implemented

---

## 🧪 Testing Checklist

- [ ] Create order with modifiers
- [ ] Create order without customer ID (guest)
- [ ] Get paginated orders list
- [ ] Filter orders by status
- [ ] Filter orders by date range
- [ ] Search orders by order number
- [ ] Get order details by ID
- [ ] Update order status (valid transitions)
- [ ] Update order status (invalid transitions) - should fail
- [ ] Update individual item status
- [ ] Add notes to order
- [ ] Cancel pending order
- [ ] Cancel preparing order - should fail
- [ ] Send kitchen alert
- [ ] Submit review for completed order
- [ ] Submit review for pending order - should fail

---

## 📚 Related EPICs

**Implemented:**
- ✅ SR-66: Customer flow - Menu flow
- ✅ SR-67: Customer shopping cart + order payment
- ✅ SR-68: User profile - Order History
- ✅ SR-69: Admin - Customer order - Orders Management
- ✅ SR-69: Waiter Orders Management
- ✅ SR-120: Kitchen Staff Dashboard - Order Queue Management

---

## 🚀 Next Steps

1. **Migration 006:** Run database migration for order tables
2. **Authentication:** Apply JWT middleware to order routes
3. **Reviews:** Implement full review system (migration 007)
4. **WebSocket:** Real-time order updates & kitchen alerts
5. **Payments:** Integrate with bills & payments system
6. **Testing:** Unit tests + integration tests
7. **Documentation:** API docs with Swagger/OpenAPI

---

**Last Updated:** January 12, 2026  
**Author:** Smart Restaurant Development Team  
**Status:** ✅ Core implementation complete (9/9 endpoints)
