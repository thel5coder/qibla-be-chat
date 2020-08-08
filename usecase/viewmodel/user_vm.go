package viewmodel

// UserVM ....
type UserVM struct {
	ID           int    `json:"id"`
	CompanyID    int    `json:"companyId"`
	RoleID       int    `json:"roleId"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	EmailValidAt string `json:"emailValidAt"`
	Phone        string `json:"phone"`
	PhoneValidAt string `json:"phoneValidAt"`
	Password     string `json:"password"`
	Photo        string `json:"photo"`
	Status       bool   `json:"status"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
	DeletedAt    string `json:"deletedAt"`
}
