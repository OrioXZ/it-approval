PRAGMA foreign_keys = ON;

INSERT OR IGNORE INTO master_status (code, label, seq, color, is_final) VALUES
('PENDING',  'รออนุมัติ', 1, 'gray',  'N'),
('APPROVED', 'อนุมัติ',   2, 'green', 'Y'),
('REJECTED', 'ไม่อนุมัติ', 3, 'red',  'Y');

INSERT INTO requests (title, status_code) VALUES
('IT-REQ-0001 ขอสิทธิ์เข้าระบบ A', 'PENDING'),
('IT-REQ-0002 ขอเพิ่มพื้นที่ Storage', 'PENDING'),
('IT-REQ-0003 ขอ VPN Access', 'PENDING'),
('IT-REQ-0004 ขอสร้าง Email Alias', 'PENDING'),
('IT-REQ-0005 ขอเปิดใช้งาน Software License', 'PENDING');

INSERT INTO requests (title, status_code, decided_reason, decided_at, decided_by) VALUES
('IT-REQ-0090 ขอสิทธิ์เข้าระบบ Finance', 'APPROVED', 'ตรวจสอบแล้วตรงตาม policy', datetime('now','-2 day'), 'demo_user'),
('IT-REQ-0091 ขอเครื่องใหม่', 'REJECTED', 'งบประมาณไม่เพียงพอ', datetime('now','-1 day'), 'demo_user');
