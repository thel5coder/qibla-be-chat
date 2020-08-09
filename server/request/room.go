package request

// NewRoomRequest ....
type NewRoomRequest struct {
	Type              string                      `json:"type" validate:"required,oneof=private group"`
	Name              string                      `json:"name"`
	ProfilePicture    string                      `json:"profile_picture"`
	Description       string                      `json:"description"`
	UserParticipantID string                      `json:"user_participant_id"`
	Participants      []NewRoomParticipantRequest `json:"participants"`
}

// NewRoomParticipantRequest ....
type NewRoomParticipantRequest struct {
	UserID     string `json:"user_id" validate:"required"`
	OdooUserID int64  `json:"odoo_user_id"`
}

// UpdateRoomRequest ....
type UpdateRoomRequest struct {
	Name           string `json:"name"`
	ProfilePicture string `json:"profile_picture"`
	Description    string `json:"description"`
}
