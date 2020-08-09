package usecase

import (
	"qibla-backend-chat/pkg/logruslogger"
)

// PusherUC ...
type PusherUC struct {
	*ContractUC
}

// SendAllParticipant ...
func (uc PusherUC) SendAllParticipant(event, roomID, userID string, body interface{}) (err error) {
	ctx := "PusherUC.SendAllParticipant"

	// Get all participant
	participantUc := ParticipantUC{ContractUC: uc.ContractUC}
	participant, err := participantUc.SelectAllByRoom(roomID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_participant", uc.ReqID)
		return err
	}

	// Send notification to all user except the user hwo make the action
	for _, p := range participant {
		if p.UserID != userID {
			err = uc.Pusher.Send(event+p.UserID, body)
			if err != nil {
				logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "pusher", uc.ReqID)
			}
		}
	}

	return err
}
