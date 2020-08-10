package usecase

import (
	"errors"
	"mime/multipart"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/model"
	"qibla-backend-chat/mongomodel"
	"qibla-backend-chat/pkg/interfacepkg"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/pkg/pusher"
	"qibla-backend-chat/usecase/viewmodel"
	"time"
)

// ChatUC ...
type ChatUC struct {
	*ContractUC
}

// FindAllByRoom ...
func (uc ChatUC) FindAllByRoom(roomID, message, lastID string, limit int) (res []viewmodel.ChatVM, err error) {
	ctx := "ChatUC.FindAllByRoom"

	limit = uc.LimitMax(limit)

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	data, err := chatModel.FindAllByRoom(roomID, message, lastID, limit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all", uc.ReqID)
		return res, err
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	for _, r := range data {
		res = append(res, viewmodel.ChatVM{
			ID:         r.ID,
			RoomID:     r.Type,
			Message:    r.Message,
			Payload:    r.Payload,
			PayloadURL: s3Uc.GetURLNoErr(r.Payload),
			Type:       r.Type,
			UserID:     r.UserID,
			CreatedAt:  r.CreatedAt,
			UpdatedAt:  r.UpdatedAt,
			DeletedAt:  r.DeletedAt,
		})
	}

	return res, err
}

// FindAllByRoomUser ...
func (uc ChatUC) FindAllByRoomUser(userID, roomID, message, lastID string, limit int) (res []viewmodel.ChatVM, err error) {
	ctx := "ChatUC.FindAllByRoomUser"

	roomUc := RoomUC{ContractUC: uc.ContractUC}
	_, err = roomUc.FindByIDParticipant(roomID, userID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "invalid_room", uc.ReqID)
		return res, errors.New(helper.InvalidRoom)
	}

	res, err = uc.FindAllByRoom(roomID, message, lastID, limit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all", uc.ReqID)
		return res, err
	}

	return res, err
}

// FindByID ...
func (uc ChatUC) FindByID(id string) (res viewmodel.ChatVM, err error) {
	ctx := "ChatUC.FindByID"

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	data, err := chatModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	res = viewmodel.ChatVM{
		ID:         data.ID,
		RoomID:     data.Type,
		Message:    data.Message,
		Payload:    data.Payload,
		PayloadURL: s3Uc.GetURLNoErr(data.Payload),
		Type:       data.Type,
		UserID:     data.UserID,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
		DeletedAt:  data.DeletedAt,
	}

	return res, err
}

// FindLast ...
func (uc ChatUC) FindLast(roomID string) (res viewmodel.ChatVM, err error) {
	ctx := "ChatUC.FindLast"

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	data, err := chatModel.FindLast(roomID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find", uc.ReqID)
		return res, err
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	res = viewmodel.ChatVM{
		ID:         data.ID,
		RoomID:     data.Type,
		Message:    data.Message,
		Payload:    data.Payload,
		PayloadURL: s3Uc.GetURLNoErr(data.Payload),
		Type:       data.Type,
		UserID:     data.UserID,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
		DeletedAt:  data.DeletedAt,
	}

	return res, err
}

// Create ...
func (uc ChatUC) Create(res *mongomodel.ChatEntity) (err error) {
	ctx := "ChatUC.Create"

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	res.ID, err = chatModel.Store(res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "store", uc.ReqID)
		return err
	}

	return err
}

// Delete ...
func (uc ChatUC) Delete(id string) (err error) {
	ctx := "ChatUC.Delete"

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	_, err = chatModel.Delete(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "destroy", uc.ReqID)
		return err
	}

	return err
}

// NewChat ...
func (uc ChatUC) NewChat(userData *viewmodel.UserVM, roomID, message string, file *multipart.FileHeader) (res viewmodel.ChatVM, err error) {
	ctx := "ChatUC.NewChat"

	// Check user in room
	roomUc := RoomUC{ContractUC: uc.ContractUC}
	_, err = roomUc.FindByIDParticipant(roomID, userData.ID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_room", uc.ReqID)
		return res, err
	}

	// Upload file if any
	payload := ""
	payloadURL := ""
	fileType := mongomodel.ChatTypeText
	if file != nil {
		s3Uc := S3UC{ContractUC: uc.ContractUC}
		fileKey, err := s3Uc.UploadFile(model.FileChat+"/"+userData.ID, file)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "upload_file", uc.ReqID)
			return res, err
		}

		payload = fileKey
		payloadURL = s3Uc.GetURLNoErr(fileKey)
		fileType = helper.GetChatFileType(payload)
	}

	// Record into chat table
	now := time.Now().UTC()
	chatBody := mongomodel.ChatEntity{
		RoomID:    roomID,
		Message:   message,
		Payload:   payload,
		Type:      fileType,
		UserID:    userData.ID,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	err = uc.Create(&chatBody)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "add_chat", uc.ReqID)
		return res, err
	}

	res = viewmodel.ChatVM{
		ID:         chatBody.ID,
		RoomID:     chatBody.RoomID,
		Message:    chatBody.Message,
		Payload:    chatBody.Payload,
		PayloadURL: payloadURL,
		Type:       chatBody.Type,
		UserID:     chatBody.UserID,
		CreatedAt:  chatBody.CreatedAt,
		UpdatedAt:  chatBody.UpdatedAt,
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventNewChat, roomID, userData.ID, res)

	return res, err
}

// DeleteChat ...
func (uc ChatUC) DeleteChat(id string, userData *viewmodel.UserVM) (res viewmodel.ChatVM, err error) {
	ctx := "ChatUC.DeleteChat"

	res, err = uc.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_chat", uc.ReqID)
		return res, err
	}
	if res.UserID != userData.ID {
		logruslogger.Log(logruslogger.WarnLevel, interfacepkg.Marshall(userData), ctx, "invalid_user", uc.ReqID)
		return res, errors.New(helper.InvalidUser)
	}

	// Delete chat from database
	err = uc.Delete(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete_chat", uc.ReqID)
		return res, err
	}

	// Delete payload
	if res.Payload != "" {
		s3Uc := S3UC{ContractUC: uc.ContractUC}
		s3Uc.Delete(res.Payload)
	}

	// Trigger pusher
	pusherUc := PusherUC{ContractUC: uc.ContractUC}
	pusherUc.SendAllParticipant(pusher.EventDeleteChat, res.RoomID, userData.ID, res)

	return res, err
}
