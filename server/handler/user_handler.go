package handler

import (
	"net/http"

	"qibla-backend-chat/usecase"
)

// UserHandler ...
type UserHandler struct {
	Handler
}

// GetByTokenHandler ...
func (h *UserHandler) GetByTokenHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	userID := user["id"].(int)

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.FindByID(userID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
