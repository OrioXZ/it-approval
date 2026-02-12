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

type Status struct {
	Code    string  `json:"code" gorm:"primaryKey"`
	Label   string  `json:"label"`
	Seq     int     `json:"seq"`
	Color   *string `json:"color"`
	IsFinal string  `json:"is_final"`
}

func (Request) TableName() string { return "requests" }
func (Status) TableName() string  { return "master_status" }
