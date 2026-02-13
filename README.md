# IT Approval System

## Overview

This project is a simple IT request approval system that allows users
to:

-   View IT requests
-   Approve requests
-   Reject requests with a reason
-   Filter requests by status (Pending / Approved / Rejected)

The system consists of:

-   **Frontend** --- Angular application
-   **Backend** --- Golang (Gin) REST API
-   **Database** --- SQLite

This project is organized as a **monorepo** containing both frontend and
backend for easier setup and review.

------------------------------------------------------------------------

## Repository Structure

    it-approval/
    ├── backend/    # Golang Gin API + SQLite database
    ├── frontend/   # Angular application
    └── README.md

------------------------------------------------------------------------

## Tech Stack

### Frontend

-   Angular
-   TypeScript
-   RxJS

### Backend

-   Golang
-   Gin Web Framework
-   GORM
-   SQLite

------------------------------------------------------------------------

## Features

-   View IT requests
-   Approve request
-   Reject request with reason
-   Status-based filtering
-   Backend integration via REST API
-   Transaction handling for status updates
-   Validation for final status updates

------------------------------------------------------------------------

## Prerequisites

-   Node.js (\>= 18 recommended)
-   npm
-   Angular CLI
-   Go (\>= 1.21 recommended)

Install Angular CLI if not installed:

``` bash
npm install -g @angular/cli
```

------------------------------------------------------------------------

## How to Run

### Start Backend

``` bash
cd backend
go run ./cmd/api

```

Backend runs at:

http://localhost:8080

### Start Frontend

``` bash
cd frontend
npm install
ng serve
```

Frontend runs at:

http://localhost:4200

------------------------------------------------------------------------

## API Endpoints

### Get Requests

GET /requests\
GET /requests?status=PENDING\
GET /requests?status=APPROVED\
GET /requests?status=REJECTED

### Update Request Status

PATCH /requests/{id}/status

#### Approve Example

``` json
{
  "status_code": "APPROVED",
  "decided_reason": "admin",
  "decided_by": "admin"
}
```

#### Reject Example

``` json
{
  "status_code": "REJECTED",
  "decided_reason": "not enough budget",
  "decided_by": "admin"
}
```

------------------------------------------------------------------------

## Business Rules

-   Final status requires:
    -   decidedReason
    -   decidedBy
-   Status updates use database transactions.
-   Requests can be filtered by status.

------------------------------------------------------------------------

## Future Improvements

-   Authentication
-   Unit tests
-   Pagination
-   Docker support
