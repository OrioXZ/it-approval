# IT Approval System (Test No.3) - Backend

Backend API for a simple approval workflow:
- List requests
- Update request status (approve/reject)
- Maintain master status metadata (label/seq/color/is_final)

Tech Stack:
- Go
- Gin
- GORM
- SQLite

## Project Structure (backend)
- `cmd/api` - API entrypoint
- `internal/db` - DB connection, models, migrations runner
- `internal/handlers` (or `handlers`) - HTTP handlers / routes
- `migrations` - SQL migration files
- `db` - runtime sqlite database file (ignored by git)

> Note: SQLite database file is ignored via `.gitignore` (e.g. `db/*.db`).

---

## How to Run

### Prerequisites
- Go 1.20+ (or your installed version)

### Run
```bash
cd backend
go mod tidy
go run ./cmd/api
```
