package viewmodel

// RoomVM ....
type RoomVM struct {
	ID                string          `json:"id"`
	Type              string          `json:"type"`
	Name              string          `json:"name"`
	ProfilePicture    string          `json:"profile_picture"`
	Description       string          `json:"description"`
	UserID            string          `json:"user_id"`
	UserParticipantID string          `json:"user_participant_id"`
	LastChat          ChatVM          `json:"last_chat"`
	Participants      []ParticipantVM `json:"participants"`
	Status            bool            `json:"status"`
	CreatedAt         string          `json:"created_at"`
	UpdatedAt         string          `json:"updated_at"`
	DeletedAt         string          `json:"deleted_at"`
}
