package request

// NewParticipantRequest ....
type NewParticipantRequest struct {
	RoomID string `json:"room_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}
