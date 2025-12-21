# Menu API - Example Requests

---

## 1. GET /api/menu - Load menu by table token
**Request:**
```bash
curl -X POST "http://164.90.145.135:8080/api/menu?table=3&token=bcec21b9b7fc4f2fef389f32cc54df2074d24cc1f2f5209ca080603efa1c04ab"
```

**Response (valid token, not expired):**
```json
{
  "menu": true
}
```

**Response (invalid or expired token):**
```json
{
  "menu": false
}
```
