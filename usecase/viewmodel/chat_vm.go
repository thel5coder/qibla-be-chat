package viewmodel

// ChatVM ....
type ChatVM struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	Message   string `json:"message"`
	Payload   string `json:"payload"`
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
}
