package usecase

import (
	"errors"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/mongomodel"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/pkg/pusher"
	"qibla-backend-chat/server/request"
	"qibla-backend-chat/usecase/viewmodel"
	"time"
)

// ParticipantUC ...
type ParticipantUC struct {
	*ContractUC
}

// SelectAllByRoom ...
func (uc ParticipantUC) SelectAllByRoom(roomID string) (res []viewmodel.ParticipantVM, err error) {
	ctx := "ParticipantUC.SelectAllByRoom"

	participantModel := mongomodel.NewParticipantModel(uc.MongoDB, uc.MongoDBName)
	data, err := participantModel.SelectAllByRoom(roomID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "select_all", uc.ReqID)
		return res, err
	}

	userUc := UserUC{ContractUC: uc.ContractUC}
	for _, r := range data {
		// Get user detail
		user, _ := userUc.FindByID(r.UserID)

		res = append(res, viewmodel.ParticipantVM{
			ID:                r.ID,
			RoomID:            r.Type,
			UserID:            r.UserID,
			Username:          user.Username,
			Email:             user.Email,
			Name:              user.Name,
			RoleID:            user.RoleID,
			RoleName:          user.RoleName,
			OdooUserID:        user.OdooUserID,
			ProfilePicture:    user.ProfilePicture,
			ProfilePictureURL: user.ProfilePictureURL,
			Type:              r.Type,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			DeletedAt:         r.DeletedAt,
		})
	}

	return res, err
}

// FindAllByRoom ...
func (uc ParticipantUC) FindAllByRoom(roomID, lastID string, limit int) (res []viewmodel.ParticipantVM, err error) {
	ctx := "ParticipantUC.FindAllByRoom"

	participantModel := mongomodel.NewParticipantModel(uc.MongoDB, uc.MongoDBName)
	data, err := participantModel.FindAllByRoom(roomID, lastID, limit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all", uc.ReqID)
		return res, err
	}

	userUc := UserUC{ContractUC: uc.ContractUC}
	for _, r := range data {
		// Get user detail
		user, _ := userUc.FindByID(r.UserID)

		res = append(res, viewmodel.ParticipantVM{
			ID:                r.ID,
			RoomID:            r.Type,
			UserID:            r.UserID,
			Username:          user.Username,
			Email:             user.Email,
			Name:              user.Name,
			RoleID:            user.RoleID,
			RoleName:          user.RoleName,
			OdooUserID:        user.OdooUserID,
			ProfilePicture:    user.ProfilePicture,
			ProfilePictureURL: user.ProfilePictureURL,
			Type:              r.Type,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			DeletedAt:         r.DeletedAt,
		})
	}

	return res, err
}

// FindByRoomParticipant ...
func (uc ParticipantUC) FindByRoomParticipant(roomID, userID string) (res viewmodel.ParticipantVM, err error) {
	ctx := "ParticipantUC.FindByRoomParticipant"

	participantModel := mongomodel.NewParticipantModel(uc.MongoDB, uc.MongoDBName)
	data, err := participantModel.FindByRoomParticipant(roomID, userID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}

	// Get user detail
	userUc := UserUC{ContractUC: uc.ContractUC}
	user, _ := userUc.FindByID(data.UserID)

	res = viewmodel.ParticipantVM{
		ID:                data.ID,
		RoomID:            data.Type,
		UserID:            data.UserID,
		Username:          user.Username,
		Email:             user.Email,
		Name:              user.Name,
		RoleID:            user.RoleID,
		RoleName:          user.RoleName,
		OdooUserID:        user.OdooUserID,
		ProfilePicture:    user.ProfilePicture,
		ProfilePictureURL: user.ProfilePictureURL,
		Type:              data.Type,
		CreatedAt:         data.CreatedAt,
		UpdatedAt:         data.UpdatedAt,
		DeletedAt:         data.DeletedAt,
	}

	return res, err
}

// Create ...
func (uc ParticipantUC) Create(body *mongomodel.ParticipantEntity) (err error) {
	ctx := "ParticipantUC.Create"

	participantModel := mongomodel.NewParticipantModel(uc.MongoDB, uc.MongoDBName)
	body.ID, err = participantModel.Store(body)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "store", uc.ReqID)
		return err
	}

	return err
}

// Delete ...
func (uc ParticipantUC) Delete(id string) (err error) {
	ctx := "ParticipantUC.Delete"

	participantModel := mongomodel.NewParticipantModel(uc.MongoDB, uc.MongoDBName)
	_, err = participantModel.Delete(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "destroy", uc.ReqID)
		return err
	}

	return err
}

// DeleteByRoomParticipant ...
func (uc ParticipantUC) DeleteByRoomParticipant(roomID, userID string) (err error) {
	ctx := "ParticipantUC.DeleteByRoomParticipant"

	participantModel := mongomodel.NewParticipantModel(uc.MongoDB, uc.MongoDBName)
	err = participantModel.DeleteByRoomParticipant(roomID, userID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "destroy", uc.ReqID)
		return err
	}

	return err
}

// NewParticipant ...
func (uc ParticipantUC) NewParticipant(userData *viewmodel.UserVM, participant *request.NewParticipantRequest) (res viewmodel.ParticipantVM, err error) {
	ctx := "ParticipantUC.NewParticipant"

	roomUc := RoomUC{ContractUC: uc.ContractUC}
	room, err := roomUc.FindByID(participant.RoomID, userData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_room", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}
	if room.Type != mongomodel.RoomTypeGroup {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "room_type", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}
	if room.UserID != userData.ID {
		logruslogger.Log(logruslogger.WarnLevel, userData.ID, ctx, "not_creator", uc.ReqID)
		return res, errors.New(helper.InvalidUser)
	}

	// Check if user already in room
	for _, p := range room.Participants {
		if p.UserID == participant.UserID {
			logruslogger.Log(logruslogger.WarnLevel, "", ctx, "already_in_room", uc.ReqID)
			return res, errors.New(helper.AlreadyInRoom)
		}
	}

	// Find participant user detail
	userUc := UserUC{ContractUC: uc.ContractUC}
	user, err := userUc.FindByID(participant.UserID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_participant", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}

	// Validate user in odoo
	newParticipants := []request.NewRoomParticipantRequest{
		{
			UserID:     user.ID,
			OdooUserID: user.OdooUserID,
		},
		{
			UserID:     userData.ID,
			OdooUserID: userData.OdooUserID,
		},
	}
	err = roomUc.CheckOddoParticipant(&user, &newParticipants)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_odoo_participant", uc.ReqID)
		return res, errors.New(helper.InvalidUser)
	}

	// Add record
	now := time.Now().UTC()
	participantBody := mongomodel.ParticipantEntity{
		RoomID:    participant.RoomID,
		UserID:    participant.UserID,
		Type:      mongomodel.ParticipantTypeUser,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	err = uc.Create(&participantBody)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "create", uc.ReqID)
		return res, err
	}

	res = viewmodel.ParticipantVM{
		ID:                participantBody.ID,
		RoomID:            participantBody.RoomID,
		UserID:            participantBody.UserID,
		Username:          user.Username,
		Email:             user.Email,
		Name:              user.Name,
		RoleID:            user.RoleID,
		RoleName:          user.RoleName,
		OdooUserID:        user.OdooUserID,
		ProfilePicture:    user.ProfilePicture,
		ProfilePictureURL: user.ProfilePictureURL,
		Type:              participantBody.Type,
		CreatedAt:         now.Format(time.RFC3339),
		UpdatedAt:         now.Format(time.RFC3339),
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventNewParticipant, res.ID, userData.ID, res)

	return res, err
}

// DeleteParticipant ...
func (uc ParticipantUC) DeleteParticipant(userData *viewmodel.UserVM, participant *request.NewParticipantRequest) (res viewmodel.ParticipantVM, err error) {
	ctx := "ParticipantUC.DeleteParticipant"

	roomUc := RoomUC{ContractUC: uc.ContractUC}
	room, err := roomUc.FindByID(participant.RoomID, userData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_room", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}
	if room.Type != mongomodel.RoomTypeGroup {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "room_type", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}
	if room.UserID != userData.ID {
		logruslogger.Log(logruslogger.WarnLevel, userData.ID, ctx, "not_creator", uc.ReqID)
		return res, errors.New(helper.InvalidUser)
	}
	if room.UserID == participant.UserID {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "is_admin_user", uc.ReqID)
		return res, errors.New(helper.InvalidUser)
	}

	res, err = uc.FindByRoomParticipant(participant.RoomID, participant.UserID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_participant", uc.ReqID)
		return res, err
	}

	err = uc.DeleteByRoomParticipant(participant.RoomID, participant.UserID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete_room", uc.ReqID)
		return res, err
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventRemoveParticipant, res.ID, userData.ID, res)

	return res, err
}

// DeleteParticipantLeave ...
func (uc ParticipantUC) DeleteParticipantLeave(roomID string, userData *viewmodel.UserVM) (res viewmodel.ParticipantVM, err error) {
	ctx := "ParticipantUC.DeleteParticipantLeave"

	roomUc := RoomUC{ContractUC: uc.ContractUC}
	room, err := roomUc.FindByID(roomID, userData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_room", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}
	if room.Type != mongomodel.RoomTypeGroup {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "room_type", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}
	if room.UserID == userData.ID {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "is_admin_user", uc.ReqID)
		return res, errors.New(helper.InvalidUser)
	}

	res, err = uc.FindByRoomParticipant(roomID, userData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_participant", uc.ReqID)
		return res, err
	}

	err = uc.DeleteByRoomParticipant(roomID, userData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete_room", uc.ReqID)
		return res, err
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventLeaveParticipant, res.ID, userData.ID, res)

	return res, err
}
