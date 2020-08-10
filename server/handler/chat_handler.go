package handler

import (
	"net/http"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/pkg/str"
	"qibla-backend-chat/usecase"

	"github.com/go-chi/chi"
)

// ChatHandler ...
type ChatHandler struct {
	Handler
}

// GetAllByRoomHandler ...
func (h *ChatHandler) GetAllByRoomHandler(w http.ResponseWriter, r *http.Request) {
	h.ContractUC.ReqID = getHeaderReqID(r)
	user := getUserDetail(r)

	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		SendBadRequest(w, helper.InvalidRoom)
		return
	}
	message := r.URL.Query().Get("message")
	lastID := r.URL.Query().Get("last_id")
	limit := str.StringToInt(r.URL.Query().Get("limit"))

	chatUC := usecase.ChatUC{ContractUC: h.ContractUC}
	res, err := chatUC.FindAllByRoomUser(user.ID, roomID, message, lastID, limit)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// CreateHandler ...
func (h *ChatHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	h.ContractUC.ReqID = getHeaderReqID(r)
	user := getUserDetail(r)

	roomID := r.PostFormValue("room_id")
	message := r.PostFormValue("message")
	if roomID == "" {
		SendBadRequest(w, helper.InvalidRoom)
		return
	}

	// Read file
	file, header, err := r.FormFile("file")
	if err == nil {
		defer file.Close()
	}

	if message == "" && file == nil {
		SendBadRequest(w, helper.InvalidBody)
		return
	}

	chatUc := usecase.ChatUC{ContractUC: h.ContractUC}
	res, err := chatUc.NewChat(&user, roomID, message, header)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// DeleteHandler ..
func (h *ChatHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	h.ContractUC.ReqID = getHeaderReqID(r)
	user := getUserDetail(r)

	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	chatUc := usecase.ChatUC{ContractUC: h.ContractUC}
	res, err := chatUc.DeleteChat(id, &user)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
