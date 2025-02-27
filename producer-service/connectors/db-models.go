package connectors

// User represents a record in the users table.
type User struct {
	ID           int      `json:"id"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	EmailAddress string   `json:"email_address"`
	CreatedAt    string   `json:"created_at"`
	DeletedAt    string   `json:"deleted_at"`
	MergedAt     string   `json:"merged_at"`
	ParentUserID *float64 `json:"parent_user_id"` // Optional field (nullable)
}
