# QR API - Example Requests

---

## 1. Admin Tables QR Code APIs

### 1.1 POST /api/admin/tables/:id/qr/generate - Generate QR code for a table

**Request:**
```bash
curl -X POST "http://164.90.145.135:8080/api/admin/tables/7/qr/generate"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "url": "https://smart-restaurant-fe.vercel.app/menu?table=7&token=38ee46f51a472e3cc06f65fac29c35a5aea48feef547f5663b3bd7405c577f40"
  }
}
```

### 1.2 GET /api/admin/tables/:id/qr/download - Download QR code for a single table
**Request:**
```bash
curl -X GET "http://164.90.145.135:8080/api/admin/tables/3/qr/download?token=bcec21b9b7fc4f2fef389f32cc54df2074d24cc1f2f5209ca080603efa1c04ab" --output table_3_qr.png
```

**Response:**

Binary PNG file download (table_3_qr.png)

Content-Type: image/png

### 1.3 GET /api/admin/tables/qr/download-all - Download all QR codes
**Request:**
```bash
curl -X GET "http://164.90.145.135:8080/api/admin/tables/qr/download-all" --output all_tables_qr.zip
```

**Response:**

Binary ZIP file (all_tables_qr.zip)

Contains PNG files for each table: table_1.png, table_2.png, …

Content-Type: application/zip

### 1.5 GET /api/admin/tables/:id/qr - Get qr code by table id

**Request:**
```bash
curl -X GET "http://164.90.145.135:8080/api/admin/tables/7/qr"
```

**Response:**
```json
{
  "code": 0,
  "message": "",
  "data": {
    "create_at": "2025-12-20T07:58:51.290515Z",
    "expire_at": "2025-12-21T07:58:51.290515Z",
    "url": "https://smart-restaurant-fe.vercel.app/menu?table=7&token=643c199c107249bc3ddc999f19fc0d7d46cb47a856d493177b4c8702d3b54512"
  }
}
```

