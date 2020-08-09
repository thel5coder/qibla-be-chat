package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"qibla-backend-chat/server/request"
	"qibla-backend-chat/usecase"
	"strconv"

	validator "gopkg.in/go-playground/validator.v9"
)

// RoomHandler ...
type RoomHandler struct {
	Handler
}

// FindAllByParticipantHandler ...
func (h *RoomHandler) FindAllByParticipantHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	message := r.URL.Query().Get("message")
	lastID := r.URL.Query().Get("last_id")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	roomUc := usecase.RoomUC{ContractUC: h.ContractUC}
	res, err := roomUc.FindAllByParticipant(user.ID, message, lastID, limit)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// FindByIDHandler ...
func (h *RoomHandler) FindByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	roomUc := usecase.RoomUC{ContractUC: h.ContractUC}
	res, err := roomUc.FindByIDParticipant(id, user.ID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// CreateHandler ...
func (h *RoomHandler) CreateHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	req := request.NewRoomRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	roomUc := usecase.RoomUC{ContractUC: h.ContractUC}
	res, err := roomUc.NewRoom(&user, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// UpdateHandler ...
func (h *RoomHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	req := request.UpdateRoomRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	roomUc := usecase.RoomUC{ContractUC: h.ContractUC}
	res, err := roomUc.UpdateRoom(id, user.ID, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// DeleteHandler ...
func (h *RoomHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	roomUc := usecase.RoomUC{ContractUC: h.ContractUC}
	res, err := roomUc.DeleteRoom(id, user.ID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
