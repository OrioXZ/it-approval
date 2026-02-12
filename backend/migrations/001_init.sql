PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS master_status (
  code       TEXT PRIMARY KEY,
  label      TEXT NOT NULL,
  seq        INTEGER NOT NULL,
  color      TEXT NULL,
  is_final   TEXT NOT NULL DEFAULT 'N' CHECK (is_final IN ('Y','N'))
);

CREATE TABLE IF NOT EXISTS requests (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  title          TEXT NOT NULL,
  status_code    TEXT NOT NULL,
  decided_reason TEXT NULL,
  decided_at     TEXT NULL,
  decided_by     TEXT NULL,
  created_at     TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at     TEXT NOT NULL DEFAULT (datetime('now')),
  FOREIGN KEY (status_code) REFERENCES master_status(code)
);

CREATE INDEX IF NOT EXISTS idx_requests_status_code ON requests(status_code);
CREATE INDEX IF NOT EXISTS idx_requests_created_at ON requests(created_at);

CREATE TRIGGER IF NOT EXISTS trg_requests_updated_at
AFTER UPDATE ON requests
FOR EACH ROW
BEGIN
  UPDATE requests
  SET updated_at = datetime('now')
  WHERE id = NEW.id;
END;
