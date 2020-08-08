package handler

import (
	"github.com/go-chi/chi"
	"net/http"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/pkg/odoo"
	"qibla-backend-chat/usecase/viewmodel"
	"strconv"

	"qibla-backend-chat/usecase"
)

// OdooHandler ...
type OdooHandler struct {
	Handler
}

// GetFieldHandler ...
func (h *OdooHandler) GetFieldHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	odooUc := usecase.OdooUC{ContractUC: h.ContractUC}
	res, err := odooUc.GetField(id)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// FindAllTravelPackageHandler ...
func (h *OdooHandler) FindAllTravelPackageHandler(w http.ResponseWriter, r *http.Request) {
	odooUc := usecase.OdooUC{ContractUC: h.ContractUC}
	var res []viewmodel.TravelPackageVM
	err := odooUc.FindAll(odoo.TravelPackage, &res)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}

// FindByIDTravelPackageHandler ...
func (h *OdooHandler) FindByIDTravelPackageHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if id == 0 {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	odooUc := usecase.OdooUC{ContractUC: h.ContractUC}
	var res []viewmodel.TravelPackageVM
	err := odooUc.FindByID(int64(id), odoo.TravelPackage, &res)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if len(res) < 1 {
		SendBadRequest(w, helper.RecordNotFound)
		return
	}

	SendSuccess(w, res[0], nil)
	return
}

// FindByIDPartnerHandler ...
func (h *OdooHandler) FindByIDPartnerHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if id == 0 {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	odooUc := usecase.OdooUC{ContractUC: h.ContractUC}
	var res []viewmodel.PartnerVM
	err := odooUc.FindByID(int64(id), odoo.Partner, &res)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if len(res) < 1 {
		SendBadRequest(w, helper.RecordNotFound)
		return
	}

	SendSuccess(w, res[0], nil)
	return
}

// FindByIDGuideHandler ...
func (h *OdooHandler) FindByIDGuideHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if id == 0 {
		SendBadRequest(w, "Parameter must be filled")
		return
	}

	odooUc := usecase.OdooUC{ContractUC: h.ContractUC}
	var res []viewmodel.GuideVM
	err := odooUc.FindByID(int64(id), odoo.Guide, &res)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}
	if len(res) < 1 {
		SendBadRequest(w, helper.RecordNotFound)
		return
	}

	SendSuccess(w, res[0], nil)
	return
}
