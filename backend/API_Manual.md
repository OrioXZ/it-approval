# 📌 IT Approval Backend – API Manual (Gin + Gorm)

เอกสารนี้อธิบายการใช้งาน REST API ที่ให้บริการโดยระบบ **IT Approval Backend**

---

## 🧭 Base Info

- **Base URL**: `http://localhost:8080`  
- **Content-Type (สำหรับ body)**: `application/json`
- **Time format**: `RFC3339` (ตัวอย่าง: `2026-02-13T01:23:45Z`)  
- **Pagination (สำหรับ list)**: `limit` / `offset`

---

## ✅ Health Check

### `GET /health`

ใช้ตรวจสอบว่า service ทำงานอยู่หรือไม่

**Example**
```http
GET /health
```

**200 OK**
```json
{
  "status": "ok"
}
```

---

## 🏷️ Statuses

### `GET /statuses`

ดึงรายการสถานะทั้งหมด (เรียงตาม `seq ASC`)

**Example**
```http
GET /statuses
```

**200 OK (array)**
```json
[
  {
    "id": 1,
    "code": "PENDING",
    "name": "Pending",
    "seq": 1,
    "isFinal": "N",
    "status": "Y"
  }
]
```

**500 Internal Server Error**
```json
{ "error": "<error message>" }
```


---

## 🧾 Requests

### 1) 🔍 Get Requests (List)

#### `GET /requests?limit=50&offset=0`

ดึงรายการ request แบบมี pagination และคืนค่า `total` มาด้วย

**Query Params**
- `limit` (optional): จำนวนรายการต่อหน้า  
  - default = `50`
  - min = `1`
  - max = `200`
- `offset` (optional): ข้ามรายการกี่ตัว
  - default = `0`
  - min = `0`

**Example**
```http
GET /requests?limit=50&offset=0
```

**200 OK**
```json
{
  "items": [
    {
      "id": 1,
      "title": "Request VPN access",
      "statusCode": "PENDING",
      "createdAt": "2026-02-13T01:23:45Z",
      "updatedAt": "2026-02-13T01:23:45Z",
      "decidedAt": null,
      "decidedBy": null,
      "decidedReason": null
    }
  ],
  "limit": 50,
  "offset": 0,
  "total": 1
}
```

**500 Internal Server Error**
```json
{ "error": "<error message>" }
```

---

### 2) ➕ Create Request

#### `POST /requests`

สร้าง request ใหม่

**Request Body**
- `title` *(required)*: ชื่อเรื่อง
- `status_code` *(optional)*: ถ้าไม่ส่งมา จะ default เป็น `"PENDING"`
  - มีการ validate ว่า `status_code` ต้องมีอยู่จริงในตาราง status

**Example**
```http
POST /requests
Content-Type: application/json

{
  "title": "Request VPN access",
  "status_code": "PENDING"
}
```

**200 OK**
```json
{
  "id": 1,
  "title": "Request VPN access",
  "statusCode": "PENDING",
  "createdAt": "2026-02-13T01:23:45Z",
  "updatedAt": "2026-02-13T01:23:45Z"
}
```

**400 Bad Request**
- JSON ไม่ถูกต้อง
```json
{ "error": "invalid json" }
```
- status_code ไม่ถูกต้อง / ไม่พบในระบบ
```json
{ "error": "invalid status_code" }
```

**500 Internal Server Error**
```json
{ "error": "<error message>" }
```

---

### 3) ✏️ Patch Request Status

#### `PATCH /requests/:id/status`

อัปเดตสถานะของ request (ใช้ transaction)

**Path Params**
- `id` *(required)*: เลข id ของ request

**Request Body (ทุก field เป็น optional)**
- `status_code`: สถานะใหม่
- `decided_reason`: เหตุผล (ใช้ตอนปิดงาน/สถานะ final)
- `decided_by`: คนตัดสินใจ (ใช้ตอนปิดงาน/สถานะ final)

**Business Rules (สำคัญ)**
1) ถ้า request เดิมเป็น **final status** แล้ว → **ห้ามแก้** (ตอบ `409 Conflict`)  
2) ถ้าจะอัปเดตเป็นสถานะที่เป็น **final** → ต้องส่ง `decided_reason` และ `decided_by` มาด้วย (ไม่งั้นตอบ `400`)  
3) ถ้าไม่ได้ส่ง field ไหนมาเลย (body ว่าง หรือทุก field เป็น null) → จะคืนข้อมูลเดิม `200 OK` โดยไม่แก้ไข

> การเป็น “final status” ตรวจจากตาราง status โดยดู `isFinal == "Y"`

**Example (update status only)**
```http
PATCH /requests/1/status
Content-Type: application/json

{
  "status_code": "IN_REVIEW"
}
```

**Example (set final status)**
```http
PATCH /requests/1/status
Content-Type: application/json

{
  "status_code": "APPROVED",
  "decided_reason": "All requirements met",
  "decided_by": "atipong"
}
```

**200 OK**
```json
{
  "id": 1,
  "title": "Request VPN access",
  "statusCode": "APPROVED",
  "decidedAt": "2026-02-13T01:25:00Z",
  "decidedBy": "atipong",
  "decidedReason": "All requirements met",
  "createdAt": "2026-02-13T01:23:45Z",
  "updatedAt": "2026-02-13T01:25:00Z"
}
```

**400 Bad Request**
- id ไม่ถูกต้อง
```json
{ "error": "invalid id" }
```
- JSON ไม่ถูกต้อง
```json
{ "error": "invalid json" }
```
- status_code ไม่รู้จัก (หาในตาราง status ไม่เจอ)
```json
{ "error": "unknown status_code" }
```
- จะ set เป็น final แต่ส่งข้อมูลไม่ครบ
```json
{ "error": "final status requires decided_reason and decided_by" }
```

**404 Not Found**
```json
{ "error": "request not found" }
```

**409 Conflict**
```json
{ "error": "cannot update request with final status" }
```

**500 Internal Server Error**
- DB error ทั่วไป
```json
{ "error": "database error" }
```
- commit ไม่สำเร็จ
```json
{ "error": "commit failed" }
```
- หรือข้อความ error จาก db ในบางกรณี
```json
{ "error": "<error message>" }
```

---

