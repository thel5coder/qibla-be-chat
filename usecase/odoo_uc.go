package usecase

import (
	"errors"
	odoo "github.com/skilld-labs/go-odoo"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/pkg/logruslogger"
	odoopkg "qibla-backend-chat/pkg/odoo"
	"qibla-backend-chat/usecase/viewmodel"
)

// OdooUC ...
type OdooUC struct {
	*ContractUC
}

// GetField ...
func (uc OdooUC) GetField(name string) (res map[string]interface{}, err error) {
	ctx := "OdooUC.GetField"

	var options *odoo.Options
	res, err = uc.Odoo.FieldsGet(name, options)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_field", uc.ReqID)
		return res, err
	}

	return res, err
}

// FindAll ...
func (uc OdooUC) FindAll(model string, res interface{}) (err error) {
	ctx := "OdooUC.FindAll"

	var criteria *odoo.Criteria
	var options *odoo.Options
	err = uc.Odoo.SearchRead(model, criteria, options, res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "search", uc.ReqID)
		return err
	}

	return err
}

// FindByIDPartner ...
func (uc OdooUC) FindByIDPartner(id int64) (res viewmodel.PartnerVM, err error) {
	ctx := "OdooUC.FindByIDPartner"

	var options *odoo.Options
	var data []viewmodel.PartnerVM
	err = uc.Odoo.Read(odoopkg.Partner, []int64{id}, options, &data)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "search", uc.ReqID)
		return res, err
	}
	if len(data) < 1 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "empty", uc.ReqID)
		return res, errors.New(helper.UserNotFound)
	}

	return data[0], err
}

// FindByID ...
func (uc OdooUC) FindByID(id int64, model string, res interface{}) (err error) {
	ctx := "OdooUC.FindByID"

	var options *odoo.Options
	err = uc.Odoo.Read(model, []int64{id}, options, res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "search", uc.ReqID)
		return err
	}

	return err
}
