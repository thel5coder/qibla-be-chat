package viewmodel

// UserVM ....
type UserVM struct {
	ID                string `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Name              string `json:"name"`
	Password          string `json:"password"`
	RoleID            string `json:"role_id"`
	RoleName          string `json:"role_name"`
	OdooUserID        int64  `json:"odoo_user_id"`
	ProfilePicture    string `json:"profile_picture"`
	ProfilePictureURL string `json:"profile_picture_url"`
	IsActive          bool   `json:"is_active"`
	CreatedAt         string `json:"createdAt"`
	UpdatedAt         string `json:"updatedAt"`
	DeletedAt         string `json:"deletedAt"`
}
