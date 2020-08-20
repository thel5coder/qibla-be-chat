package usecase

import (
	"errors"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/mongomodel"
	"qibla-backend-chat/pkg/interfacepkg"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/pkg/number"
	"qibla-backend-chat/pkg/pusher"
	"qibla-backend-chat/server/request"
	"qibla-backend-chat/usecase/viewmodel"
	"strings"
	"time"
)

// RoomUC ...
type RoomUC struct {
	*ContractUC
}

// FindAllByParticipant ...
func (uc RoomUC) FindAllByParticipant(participantID, message, lastID string, limit int) (res []viewmodel.RoomVM, err error) {
	ctx := "RoomUC.FindAllByParticipant"

	limit = uc.LimitMax(limit)

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	data, err := roomModel.FindAllByParticipant(participantID, message, lastID, limit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all", uc.ReqID)
		return res, err
	}

	chatUc := ChatUC{ContractUC: uc.ContractUC}
	participantUc := ParticipantUC{ContractUC: uc.ContractUC}
	s3Uc := S3UC{ContractUC: uc.ContractUC}
	for _, r := range data {
		// Set name for private chat
		if r.Type == mongomodel.RoomTypePrivate {
			nameArr := strings.Split(r.Name, ";")
			if len(nameArr) != 2 {
				r.Name = "Chat-" + r.Name
			} else if participantID == r.UserID {
				r.Name = nameArr[1]
			} else {
				r.Name = nameArr[0]
			}
		}

		// Find Last Chat
		lastChat, _ := chatUc.FindLast(r.ID)

		// Find participant
		status := true
		if r.Type == mongomodel.RoomTypePrivate {
			participant, _ := participantUc.SelectAllByRoom(r.ID)
			if len(participant) < 2 {
				status = false
			}
		}

		res = append(res, viewmodel.RoomVM{
			ID:                r.ID,
			Type:              r.Type,
			Name:              r.Name,
			ProfilePicture:    r.ProfilePicture,
			ProfilePictureURL: s3Uc.GetURLNoErr(r.ProfilePicture),
			Description:       r.Description,
			UserID:            r.UserID,
			UserParticipantID: r.UserParticipantID,
			LastChat:          lastChat,
			Status:            status,
			CreatedAt:         r.CreatedAt,
			UpdatedAt:         r.UpdatedAt,
			DeletedAt:         r.DeletedAt,
		})
	}

	return res, err
}

// BuildBody ...
func (uc RoomUC) BuildBody(userID string, data *mongomodel.RoomEntity, res *viewmodel.RoomVM) {
	// Set name for private chat
	if data.Type == mongomodel.RoomTypePrivate {
		nameArr := strings.Split(data.Name, ";")
		if len(nameArr) != 2 {
			data.Name = "Chat-" + data.Name
		} else if userID == data.UserID {
			data.Name = nameArr[1]
		} else {
			data.Name = nameArr[0]
		}
	}

	// Find Last Chat
	chatUc := ChatUC{ContractUC: uc.ContractUC}
	lastChat, _ := chatUc.FindLast(data.ID)

	// Find participant
	status := true
	participantUc := ParticipantUC{ContractUC: uc.ContractUC}
	participant, _ := participantUc.SelectAllByRoom(data.ID)
	if data.Type == mongomodel.RoomTypePrivate && len(participant) < 2 {
		status = false
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	res.ID = data.ID
	res.Type = data.Type
	res.Name = data.Name
	res.ProfilePicture = data.ProfilePicture
	res.ProfilePictureURL = s3Uc.GetURLNoErr(data.ProfilePicture)
	res.Description = data.Description
	res.UserID = data.UserID
	res.UserParticipantID = data.UserParticipantID
	res.LastChat = lastChat
	res.Participants = participant
	res.Status = status
	res.CreatedAt = data.CreatedAt
	res.UpdatedAt = data.UpdatedAt
	res.DeletedAt = data.DeletedAt
}

// FindByID ...
func (uc RoomUC) FindByID(id, userID string) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.FindByID"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	data, err := roomModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}

	uc.BuildBody(userID, &data, &res)

	return res, err
}

// FindByIDParticipant ...
func (uc RoomUC) FindByIDParticipant(id, userID string) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.FindByIDParticipant"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	data, err := roomModel.FindByIDParticipant(id, userID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}
	if data.ID == "" {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "empty", uc.ReqID)
		return res, errors.New(helper.RecordNotExist)
	}

	uc.BuildBody(userID, &data, &res)

	return res, err
}

// FindByProfilePicture ...
func (uc RoomUC) FindByProfilePicture(userID, profilePicture string) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.FindByProfilePicture"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	data, err := roomModel.FindByProfilePicture(userID, profilePicture)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}

	uc.BuildBody(userID, &data, &res)

	return res, err
}

// FindPrivateByUser ...
func (uc RoomUC) FindPrivateByUser(userID, userParticipantID string) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.FindPrivateByUser"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	data, err := roomModel.FindPrivateByUser(userID, userParticipantID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}

	uc.BuildBody(userID, &data, &res)

	return res, err
}

// Create ...
func (uc RoomUC) Create(res *mongomodel.RoomEntity) (err error) {
	ctx := "RoomUC.Create"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	res.ID, err = roomModel.Store(res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "store", uc.ReqID)
		return err
	}

	return err
}

// Update ...
func (uc RoomUC) Update(res *mongomodel.RoomEntity) (err error) {
	ctx := "RoomUC.Update"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	res.ID, err = roomModel.Update(res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "update", uc.ReqID)
		return err
	}

	return err
}

// Delete ...
func (uc RoomUC) Delete(id string) (err error) {
	ctx := "RoomUC.Delete"

	roomModel := mongomodel.NewRoomModel(uc.MongoDB, uc.MongoDBName)
	_, err = roomModel.Delete(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "destroy", uc.ReqID)
		return err
	}

	return err
}

// CheckValidParticipant ...
func (uc RoomUC) CheckValidParticipant(userData *viewmodel.UserVM, data *request.NewRoomRequest) (err error) {
	ctx := "RoomUC.CheckValidParticipant"

	// Check  valid participant and add self id if not exist in participant
	addSelf := true
	privateCreatorName := "-"
	privateParticipantName := "-"
	userUc := UserUC{ContractUC: uc.ContractUC}
	for i, p := range data.Participants {
		user, err := userUc.FindByID(p.UserID)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_participant", uc.ReqID)
			return err
		}
		if data.Type == mongomodel.RoomTypePrivate && p.UserID != userData.ID {
			privateParticipantName = user.Name
			data.UserParticipantID = p.UserID
		}

		if p.UserID == userData.ID {
			addSelf = false
			privateCreatorName = user.Name
		}

		// Fill odoo userid
		data.Participants[i].OdooUserID = user.OdooUserID
	}
	if addSelf {
		data.Participants = append(data.Participants, request.NewRoomParticipantRequest{
			UserID:     userData.ID,
			OdooUserID: userData.OdooUserID,
		})
		privateCreatorName = userData.Name
	}

	// Add name if type is private
	if data.Type == mongomodel.RoomTypePrivate {
		data.Name = privateCreatorName + ";" + privateParticipantName
	}

	// Remove duplicate user id
	data.Participants = helper.UniqueRoomParticipantData(&data.Participants)

	return err
}

// CheckOddoParticipant ...
func (uc RoomUC) CheckOddoParticipant(userData *viewmodel.UserVM, participants *[]request.NewRoomParticipantRequest) (err error) {
	ctx := "RoomUC.CheckOddoParticipant"

	// Get detail user to get travel package id
	odooUc := OdooUC{ContractUC: uc.ContractUC}
	partner, err := odooUc.FindByIDPartner(userData.OdooUserID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_partner", uc.ReqID)
		return err
	}

	// Getall  travel package by user
	odooPackage, err := odooUc.FindAllPackage(partner.PackageListID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_package", uc.ReqID)
		return err
	}

	// Append odoo user id into array of int64
	var participantArr []int64
	for _, p := range *participants {
		participantArr = append(participantArr, p.OdooUserID)
	}

	isValid := false
	for _, o := range odooPackage {
		if number.IntArrInArr(participantArr, o.UserList) {
			isValid = true
		}
	}

	if !isValid {
		logruslogger.Log(logruslogger.WarnLevel, interfacepkg.Marshall(participantArr), ctx, "find_package", uc.ReqID)
		return errors.New(helper.InvalidParticipant)
	}

	return err
}

// NewRoom ...
func (uc RoomUC) NewRoom(userData *viewmodel.UserVM, data *request.NewRoomRequest) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.NewRoom"

	// Check  valid participant and add self id if not exist in participant
	err = uc.CheckValidParticipant(userData, data)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_valid_participant", uc.ReqID)
		return res, err
	}

	// Verify body
	if data.Type == mongomodel.RoomTypeGroup && data.Name == "" {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "data_name", uc.ReqID)
		return res, errors.New(helper.InvalidGroupName)
	}
	if data.Type == mongomodel.RoomTypePrivate && len(data.Participants) != 2 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "data_participant", uc.ReqID)
		return res, errors.New(helper.InvalidParticipant)
	}

	// If type is private check room exist or not
	if data.Type == mongomodel.RoomTypePrivate {
		existRoom, _ := uc.FindPrivateByUser(userData.ID, data.UserParticipantID)
		if existRoom.ID != "" && existRoom.Status {
			logruslogger.Log(logruslogger.WarnLevel, interfacepkg.Marshall(existRoom), ctx, "room_exist", uc.ReqID)
			return res, errors.New(helper.RoomExist)
		}
	}

	// Check validity participant from odoo
	err = uc.CheckOddoParticipant(userData, &data.Participants)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "check_odoo_participant", uc.ReqID)
		return res, err
	}

	// Add room
	now := time.Now().UTC()
	roomBody := mongomodel.RoomEntity{
		Type:              data.Type,
		Name:              data.Name,
		ProfilePicture:    data.ProfilePicture,
		Description:       data.Description,
		UserParticipantID: data.UserParticipantID,
		UserID:            userData.ID,
		CreatedAt:         now.Format(time.RFC3339),
		UpdatedAt:         now.Format(time.RFC3339),
	}
	err = uc.Create(&roomBody)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "add_room", uc.ReqID)
		return res, err
	}

	// Add participant
	participantUc := ParticipantUC{ContractUC: uc.ContractUC}
	for _, p := range data.Participants {
		types := mongomodel.ParticipantTypeUser
		if p.UserID == userData.ID {
			types = mongomodel.ParticipantTypeAdmin
		}

		participantBody := mongomodel.ParticipantEntity{
			RoomID:    roomBody.ID,
			UserID:    p.UserID,
			Type:      types,
			CreatedAt: now.Format(time.RFC3339),
			UpdatedAt: now.Format(time.RFC3339),
		}
		participantUc.Create(&participantBody)
	}

	// Build response struct
	s3Uc := S3UC{ContractUC: uc.ContractUC}
	res = viewmodel.RoomVM{
		ID:                roomBody.ID,
		Type:              roomBody.Type,
		Name:              roomBody.Name,
		ProfilePicture:    roomBody.ProfilePicture,
		ProfilePictureURL: s3Uc.GetURLNoErr(roomBody.ProfilePicture),
		Description:       roomBody.Description,
		UserParticipantID: roomBody.UserParticipantID,
		UserID:            userData.ID,
		CreatedAt:         roomBody.CreatedAt,
		UpdatedAt:         roomBody.UpdatedAt,
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventNewRoomUser, res.ID, userData.ID, res)

	return res, err
}

// UpdateRoom ...
func (uc RoomUC) UpdateRoom(id, userID string, body *request.UpdateRoomRequest) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.UpdateRoom"

	now := time.Now().UTC()
	roomBody := mongomodel.RoomEntity{
		ID:             id,
		Name:           body.Name,
		ProfilePicture: body.ProfilePicture,
		Description:    body.Description,
		UserID:         userID,
		CreatedAt:      now.Format(time.RFC3339),
		UpdatedAt:      now.Format(time.RFC3339),
	}
	err = uc.Update(&roomBody)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "update", uc.ReqID)
		return res, err
	}

	// Build response struct
	s3Uc := S3UC{ContractUC: uc.ContractUC}
	res = viewmodel.RoomVM{
		ID:                roomBody.ID,
		Name:              roomBody.Name,
		ProfilePicture:    roomBody.ProfilePicture,
		ProfilePictureURL: s3Uc.GetURLNoErr(roomBody.ProfilePicture),
		Description:       roomBody.Description,
		UserID:            userID,
		CreatedAt:         roomBody.CreatedAt,
		UpdatedAt:         roomBody.UpdatedAt,
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventUpdateRoomUser, id, userID, res)

	return res, err
}

// DeleteRoom ...
func (uc RoomUC) DeleteRoom(id, userID string) (res viewmodel.RoomVM, err error) {
	ctx := "RoomUC.DeleteRoom"

	res, err = uc.FindByID(id, userID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_id", uc.ReqID)
		return res, err
	}
	if res.Type == mongomodel.RoomTypeGroup && res.UserID != userID {
		logruslogger.Log(logruslogger.WarnLevel, userID, ctx, "invalid_user", uc.ReqID)
		return res, errors.New(helper.InvalidRoomType)
	}

	if res.Type == mongomodel.RoomTypeGroup {
		err = uc.Delete(id)
	} else {
		participantUc := ParticipantUC{ContractUC: uc.ContractUC}
		err = participantUc.DeleteByRoomParticipant(id, userID)
	}
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete", uc.ReqID)
		return res, err
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventDeleteRoomUser, id, userID, res)

	return res, err
}
