package viewmodel

// FileVM ....
type FileVM struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Path      string `json:"path"`
	TempPath  string `json:"temp_path"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
