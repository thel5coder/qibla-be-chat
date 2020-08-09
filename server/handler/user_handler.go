package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"

	"qibla-backend-chat/usecase"
)

// UserHandler ...
type UserHandler struct {
	Handler
}

// GetByTokenHandler ...
func (h *UserHandler) GetByTokenHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	userID := user["id"].(string)

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.FindByID(userID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// GetTravelPackageHandler ...
func (h *UserHandler) GetTravelPackageHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	userID := user["id"].(string)

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.FindTravelPackage(userID)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// GetJamaahHandler ...
func (h *UserHandler) GetJamaahHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	userID := user["id"].(string)

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if id == 0 {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.FindJamaah(userID, int64(id))
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// LoginHandler ...
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	userUc := usecase.UserUC{ContractUC: h.ContractUC}
	res, err := userUc.Login(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
