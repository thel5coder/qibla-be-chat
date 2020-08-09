package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"qibla-backend-chat/server/request"
	"qibla-backend-chat/usecase"

	validator "gopkg.in/go-playground/validator.v9"
)

// ParticipantHandler ...
type ParticipantHandler struct {
	Handler
}

// AddParticipantHandler ...
func (h *ParticipantHandler) AddParticipantHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	req := request.NewParticipantRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	participantUc := usecase.ParticipantUC{ContractUC: h.ContractUC}
	res, err := participantUc.NewParticipant(&user, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// RemoveParticipantHandler ...
func (h *ParticipantHandler) RemoveParticipantHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	req := request.NewParticipantRequest{}
	if err := h.Handler.Bind(r, &req); err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if err := h.Handler.Validate.Struct(req); err != nil {
		h.SendRequestValidationError(w, err.(validator.ValidationErrors))
		return
	}

	participantUc := usecase.ParticipantUC{ContractUC: h.ContractUC}
	res, err := participantUc.DeleteParticipant(&user, &req)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// LeaveParticipantHandler ...
func (h *ParticipantHandler) LeaveParticipantHandler(w http.ResponseWriter, r *http.Request) {
	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Get detail user
	user := getUserDetail(r)

	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	participantUc := usecase.ParticipantUC{ContractUC: h.ContractUC}
	res, err := participantUc.DeleteParticipantLeave(id, &user)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
