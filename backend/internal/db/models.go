package db

type Request struct {
	ID            int     `json:"id" gorm:"primaryKey"`
	Title         string  `json:"title"`
	StatusCode    string  `json:"status_code"`
	DecidedReason *string `json:"decided_reason"`
	DecidedAt     *string `json:"decided_at"`
	DecidedBy     *string `json:"decided_by"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}
