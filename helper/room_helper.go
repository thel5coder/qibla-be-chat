package helper

import (
	"qibla-backend-chat/server/request"
)

// UniqueRoomParticipantData ...
func UniqueRoomParticipantData(slice *[]request.NewRoomParticipantRequest) []request.NewRoomParticipantRequest {
	keys := make(map[string]bool)
	list := []request.NewRoomParticipantRequest{}
	for _, entry := range *slice {
		if _, value := keys[entry.UserID]; !value {
			keys[entry.UserID] = true
			list = append(list, entry)
		}
	}

	return list
}
