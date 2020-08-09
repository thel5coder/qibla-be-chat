package usecase

import (
	"qibla-backend-chat/mongomodel"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/usecase/viewmodel"
)

// ChatUC ...
type ChatUC struct {
	*ContractUC
}

// FindAllByRoom ...
func (uc ChatUC) FindAllByRoom(roomID, message, lastID string, limit int) (res []viewmodel.ChatVM, err error) {
	ctx := "ChatUC.FindAllByRoom"

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	data, err := chatModel.FindAllByRoom(roomID, message, lastID, limit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all", uc.ReqID)
		return res, err
	}

	for _, r := range data {
		res = append(res, viewmodel.ChatVM{
			ID:        r.ID,
			RoomID:    r.Type,
			Message:   r.Message,
			Payload:   r.UserID,
			Type:      r.Type,
			UserID:    r.UserID,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			DeletedAt: r.DeletedAt,
		})
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

	res = viewmodel.ChatVM{
		ID:        data.ID,
		RoomID:    data.Type,
		Message:   data.Message,
		Payload:   data.UserID,
		Type:      data.Type,
		UserID:    data.UserID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt,
	}

	return res, err
}

// Create ...
func (uc ChatUC) Create(body *mongomodel.ChatEntity) (res string, err error) {
	ctx := "ChatUC.Create"

	chatModel := mongomodel.NewChatModel(uc.MongoDB, uc.MongoDBName)
	res, err = chatModel.Store(body)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "store", uc.ReqID)
		return res, err
	}

	return res, err
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
