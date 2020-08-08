package request

// AdminRequest ....
type AdminRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
	RoleID   string `json:"role_id" validate:"required"`
	Status   bool   `json:"status"`
}

// AdminLoginRequest ....
type AdminLoginRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required"`
}
