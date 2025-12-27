# Load Menu API Documentation

## Overview
Load Menu API cho phép customer lấy menu thông qua QR code (table + token). Endpoint này hỗ trợ pagination, search, và filter theo category.

---

## Endpoint

### Load Menu by Table QR Code
**GET** `/api/menu`

**Description:** Lấy menu cho một table cụ thể. Yêu cầu `table` ID và `token` từ QR code.

**Query Parameters:**
- `table` (required) - Table ID
- `token` (required) - QR token (phải hợp lệ và chưa hết hạn)
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20)
- `search` (optional) - Search by item name
- `sort` (optional) - Sort field: `name`, `price_asc`, `price_desc`, `last_update` (default: id)
- `category` (optional) - Filter by category name

**Response (200 OK):**
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
        "last_update": "2025-12-25",
        "chef_recommended": true,
        "image_url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/uuid.jpg",
        "description": "Fresh Atlantic salmon grilled to perfection with herbs",
        "preparation_time": 15
      },
      {
        "id": 2,
        "name": "Caesar Salad",
        "category": "Appetizers",
        "price": 12.99,
        "status": "Available",
        "last_update": "2025-12-25",
        "chef_recommended": false,
        "image_url": "https://smart-restaurant.sfo3.digitaloceanspaces.com/menu-items/uuid.jpg",
        "description": "Fresh romaine with parmesan and croutons",
        "preparation_time": 10
      }
    ]
  }
}
```

**Error Response (400 Bad Request - Missing Parameters):**
```json
{
  "menu": false,
  "error": "table and token required"
}
```

**Error Response (400 Bad Request - Invalid Table ID):**
```json
{
  "menu": false,
  "error": "invalid table id"
}
```

**Error Response (403 Forbidden - Invalid or Expired Token):**
```json
{
  "menu": false
}
```

---

## Example Usage

### Basic Request - Get Menu for Table
```bash
# Get menu for table 5 with valid token
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz" \
  -H "Content-Type: application/json"

# Pretty print
curl -s -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz" | jq '.'
```

### With Pagination
```bash
# Get page 2 with 10 items per page
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&page=2&page_size=10"

# Get first 5 items
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&page_size=5"
```

### With Search
```bash
# Search for items containing "salmon"
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&search=salmon"

# Search with pagination
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&search=salad&page=1&page_size=20"
```

### With Category Filter
```bash
# Get only items from "Main Course" category
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&category=Main%20Course"

# Get appetizers
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&category=Appetizers"

# Category filter with search
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&category=Main%20Course&search=salmon"
```

### With Sorting
```bash
# Sort by name ascending
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&sort=name"

# Sort by price ascending
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&sort=price_asc"

# Sort by price descending
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&sort=price_desc"

# Sort by last update
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&sort=last_update"
```

### Complex Query - Combined Filters
```bash
# Get page 1 of "Main Course" items, search "salmon", sorted by price
curl -X GET "http://localhost:8080/api/menu?table=5&token=abc123xyz&category=Main%20Course&search=salmon&sort=price_asc&page=1&page_size=10"
```

---


```jsx
import React, { useState, useEffect } from 'react';

export function MenuView({ tableId, qrToken }) {
  const [menu, setMenu] = useState(null);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState('');
  const [category, setCategory] = useState('');
  const [page, setPage] = useState(1);

  const fetchMenu = async (searchTerm = '', categoryFilter = '', pageNum = 1) => {
    const params = new URLSearchParams({
      table: tableId,
      token: qrToken,
      search: searchTerm,
      category: categoryFilter,
      page: pageNum,
      page_size: 20,
      sort: 'name'
    });

    try {
      const response = await fetch(`http://localhost:8080/api/menu?${params}`);
      const data = await response.json();
      
      if (data.code === 0) {
        setMenu(data.data);
      } else {
        console.error('Failed to load menu:', data.message);
      }
    } catch (error) {
      console.error('Error fetching menu:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchMenu(search, category, page);
  }, [search, category, page, tableId, qrToken]);

  if (loading) return <div>Loading menu...</div>;
  if (!menu) return <div>Failed to load menu</div>;

  return (
    <div className="menu">
      <div className="filters">
        <input
          type="text"
          placeholder="Search menu items..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
        />
        <select 
          value={category} 
          onChange={(e) => {
            setCategory(e.target.value);
            setPage(1);
          }}
        >
          <option value="">All Categories</option>
          <option value="Appetizers">Appetizers</option>
          <option value="Main Course">Main Course</option>
          <option value="Beverages">Beverages</option>
        </select>
      </div>

      <div className="items">
        {menu.items.map(item => (
          <div key={item.id} className="item">
            <img src={item.image_url} alt={item.name} />
            <h3>{item.name}</h3>
            <p className="category">{item.category}</p>
            <p>{item.description}</p>
            <div className="meta">
              <span className="price">${item.price.toFixed(2)}</span>
              <span className="status">{item.status}</span>
              {item.chef_recommended && <span className="badge">Chef Recommended</span>}
            </div>
          </div>
        ))}
      </div>

      <div className="pagination">
        <button onClick={() => setPage(Math.max(1, page - 1))} disabled={page === 1}>
          Previous
        </button>
        <span>Page {menu.page} of {Math.ceil(menu.total / menu.page_size)}</span>
        <button 
          onClick={() => setPage(page + 1)} 
          disabled={page >= Math.ceil(menu.total / menu.page_size)}
        >
          Next
        </button>
      </div>
    </div>
  );
}
```

---

## Error Codes

| Code | Message | Description |
|------|---------|-------------|
| 0 | Success | Request successful |
| 400 | Bad Request | Missing or invalid parameters |
| 403 | Forbidden | Invalid or expired token |
| 500 | Internal Server Error | Server error |

---

## Data Structure

```
Response
├── code: number (0 = success)
├── message: string
└── data: BaseListResponse
    ├── total: number (total items available)
    ├── page: number (current page)
    ├── page_size: number (items per page)
    └── items: MenuItemResponse[]
        ├── id: number
        ├── name: string
        ├── category: string
        ├── price: number
        ├── status: string (Available, Unavailable, Sold Out)
        ├── last_update: string (YYYY-MM-DD format)
        ├── chef_recommended: boolean
        ├── image_url: string
        ├── description: string
        └── preparation_time: number (minutes)
```

---

## Query Parameter Combinations

| Use Case | Query String |
|----------|-------------|
| Get first 20 items | `?page=1&page_size=20` |
| Search by name | `?search=salmon` |
| Filter by category | `?category=Main%20Course` |
| Sort by price | `?sort=price_asc` |
| All together | `?search=salmon&category=Main%20Course&sort=price_asc&page=1` |

---

## Notes
- `table` và `token` là bắt buộc và phải hợp lệ
- Token phải chưa hết hạn (kiểm tra `qr_token_expires_at`)
- Chỉ trả về items từ restaurant của table đó
- Status hiển thị dạng Title Case: `Available`, `Unavailable`, `Sold Out`
- Search không phân biệt hoa thường
- Category filter cũng không phân biệt hoa thường
- Default sort là theo ID ascending nếu không chỉ định
- Page size max recommended là 50 để tránh quá tải
