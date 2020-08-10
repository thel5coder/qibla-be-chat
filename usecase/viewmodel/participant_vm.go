package viewmodel

// ParticipantVM ....
type ParticipantVM struct {
	ID                string `json:"id"`
	RoomID            string `json:"room_id"`
	UserID            string `json:"user_id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	RoleID            string `json:"role_id"`
	RoleName          string `json:"role_name"`
	OdooUserID        int64  `json:"odoo_user_id"`
	ProfilePicture    string `json:"profile_picture"`
	ProfilePictureURL string `json:"profile_picture_url"`
	Type              string `json:"type"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
	DeletedAt         string `json:"deleted_at"`
}
