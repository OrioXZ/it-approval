# API Manual - IT Approval System (Test No.3)

## Base URL
- Local: `http://localhost:8080`

## Content-Type
- `application/json`

## Common Response (Error)
| HTTP | Meaning |
|------|---------|
| 400  | Invalid request (bad id, bad json, invalid status_code, missing fields) |
| 404  | Not found |
| 409  | Conflict (request already in final status) |
| 500  | Server/DB error |

Example:
```json
{ "error": "invalid json" }


1) Health
GET /health

Check server status.

Response 200
{ "status": "ok" }

2) Master Status
GET /statuses

Get all master statuses (used by FE for label/order/color mapping).

Response 200
[
  {
    "code": "PENDING",
    "label": "Pending",
    "seq": 1,
    "color": "#AAAAAA",
    "is_final": "N"
  },
  {
    "code": "APPROVED",
    "label": "Approved",
    "seq": 2,
    "color": "#00AA00",
    "is_final": "Y"
  },
  {
    "code": "REJECTED",
    "label": "Rejected",
    "seq": 3,
    "color": "#AA0000",
    "is_final": "Y"
  }
]

3) Requests
GET /requests

List requests.

Query Params (Optional)

limit (default 50, max 200)

offset (default 0)

Example

GET /requests?limit=50&offset=0

Response 200
{
  "data": [
    {
      "id": 1,
      "title": "Request A",
      "status_code": "PENDING",
      "decided_reason": null,
      "decided_at": null,
      "decided_by": null,
      "created_at": "2026-02-13T01:00:00Z",
      "updated_at": "2026-02-13T01:00:00Z"
    }
  ],
  "limit": 50,
  "offset": 0,
  "total": 1
}

POST /requests

Create a new request.

Request Body
{
  "title": "Request A"
}

Behavior

Default status_code = PENDING

Response 201
{
  "id": 1,
  "title": "Request A",
  "status_code": "PENDING",
  "decided_reason": null,
  "decided_at": null,
  "decided_by": null,
  "created_at": "2026-02-13T01:00:00Z",
  "updated_at": "2026-02-13T01:00:00Z"
}

Error 400 (missing title / invalid json)
{ "error": "invalid json" }

PATCH /requests/:id/status

Update request status (approve/reject) and decision fields.

Request Body (partial update supported)
{
  "status_code": "APPROVED",
  "decided_reason": "Looks good",
  "decided_by": "atipong"
}

Rules

status_code must exist in master_status

If target status is final (is_final = 'Y'), then:

decided_reason is required

decided_by is required

decided_at will be set automatically

If current request status is already final, further updates are rejected (409)

Response 200
{
  "id": 1,
  "title": "Request A",
  "status_code": "APPROVED",
  "decided_reason": "Looks good",
  "decided_at": "2026-02-13T01:10:00Z",
  "decided_by": "atipong",
  "created_at": "2026-02-13T01:00:00Z",
  "updated_at": "2026-02-13T01:10:00Z"
}

Error 404 (not found)
{ "error": "not found" }

Error 409 (already final)
{ "error": "cannot update request with final status" }

4) Quick Test (Postman / curl)
Approve

PATCH /requests/1/status

{
  "status_code": "APPROVED",
  "decided_reason": "ok",
  "decided_by": "atipong"
}

Reject

PATCH /requests/1/status

{
  "status_code": "REJECTED",
  "decided_reason": "not ok",
  "decided_by": "atipong"
}

Try patch again after final (should be 409)

PATCH /requests/1/status

{ "status_code": "APPROVED" }