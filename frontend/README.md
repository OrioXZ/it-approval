# IT Approval Frontend

Frontend application for IT Approval System.

This application allows users to view requests by status and perform approval or rejection actions.

---

## Tech Stack

- Angular
- TypeScript
- RxJS
- REST API integration

---

## Features

- View requests by status (Pending / Approved / Rejected)
- Approve request with reason
- Reject request with reason
- Status-based filtering
- Backend integration via REST API

---

## Prerequisites

- Node.js (>= 18 recommended)
- npm
- Angular CLI

Install Angular CLI if not installed:

```bash
npm install -g @angular/cli
```

---

## Setup & Run

Install dependencies:

```bash
npm install
```

Run development server:

```bash
ng serve
```

Open browser:

```
http://localhost:4200
```

---

## Backend Requirement

Backend service must be running at:

```
http://localhost:8080
```

Make sure:

- Backend server is running
- CORS is enabled
- Database is initialized

---

## API Integration

### Get Requests

```
GET /requests?status=PENDING|APPROVED|REJECTED
```

Returns request list by status.

---

### Update Request Status

```
PATCH /requests/{id}/status
```

Payload:

```json
{
  "status_code": "APPROVED | REJECTED",
  "decided_reason": "string",
  "decided_by": "string"
}
```

---

## Design Notes

- Service layer handles API communication
- Backend uses snake_case fields, mapped in frontend service
- Status workflow: `PENDING → APPROVED / REJECTED`
- Angular Signals used for state management

---

