package usecase

import (
	"errors"
	"fmt"
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

// FindByIDArray ...
func (uc OdooUC) FindByIDArray(id []int64, model string, res interface{}) (err error) {
	ctx := "OdooUC.FindByIDArray"

	var options *odoo.Options
	err = uc.Odoo.Read(model, id, options, res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "search", uc.ReqID)
		return err
	}

	return err
}

// FindByIDTravelPackage ...
func (uc OdooUC) FindByIDTravelPackage(id int64) (res viewmodel.TravelPackageVM, err error) {
	ctx := "OdooUC.FindByIDTravelPackage"

	var data []viewmodel.TravelPackageVM
	err = uc.FindByID(int64(id), odoopkg.TravelPackage, &data)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_id", uc.ReqID)
		return res, err
	}
	if len(data) < 1 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "empty", uc.ReqID)
		return res, errors.New(helper.RecordNotFound)
	}

	// Get first data in array
	res = data[0]

	// Append jamaah id into user list
	for _, rj := range res.JamaahList {
		if fmt.Sprintf("%T", rj) == "int64" {
			res.UserList = append(res.UserList, rj.(int64))
		}
	}

	// Append tour guide into list of id
	for _, rg := range res.GuideList {
		if fmt.Sprintf("%T", rg) == "int64" {
			res.GuideListID = append(res.GuideListID, rg.(int64))
		}
	}

	// Get guide by guide list id
	var guide []viewmodel.GuideVM
	uc.FindByIDArray(res.GuideListID, odoopkg.Guide, &guide)
	for _, g := range guide {
		for _, gp := range g.PartnerID {
			if fmt.Sprintf("%T", gp) == "int64" {
				res.UserList = append(res.UserList, gp.(int64))
			}
		}
	}

	return res, err
}

// FindByIDPartner ...
func (uc OdooUC) FindByIDPartner(id int64) (res viewmodel.PartnerVM, err error) {
	ctx := "OdooUC.FindByIDPartner"

	var data []viewmodel.PartnerVM
	err = uc.FindByID(int64(id), odoopkg.Partner, &data)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_id", uc.ReqID)
		return res, err
	}
	if len(data) < 1 {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "empty", uc.ReqID)
		return res, errors.New(helper.RecordNotFound)
	}

	// Get first data in array
	res = data[0]

	// Convert package list into integer
	for _, d := range res.PackageList {
		if fmt.Sprintf("%T", d) == "int64" {
			res.PackageListID = append(res.PackageListID, d.(int64))
		}
	}

	return res, err
}

// FindAllPackage ...
func (uc OdooUC) FindAllPackage(id []int64) (res []viewmodel.TravelPackageVM, err error) {
	ctx := "OdooUC.FindAllPackage"

	err = uc.FindByIDArray(id, odoopkg.TravelPackage, &res)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all", uc.ReqID)
		return res, err
	}

	for i, r := range res {
		// Append jamaah id into user list
		for _, rj := range r.JamaahList {
			if fmt.Sprintf("%T", rj) == "int64" {
				res[i].UserList = append(res[i].UserList, rj.(int64))
			}
		}

		// Append tour guide into list of id
		for _, rg := range r.GuideList {
			if fmt.Sprintf("%T", rg) == "int64" {
				res[i].GuideListID = append(res[i].GuideListID, rg.(int64))
			}
		}

		// Get guide by guide list id
		var guide []viewmodel.GuideVM
		uc.FindByIDArray(res[i].GuideListID, odoopkg.Guide, &guide)
		for _, g := range guide {
			for _, gp := range g.PartnerID {
				if fmt.Sprintf("%T", gp) == "int64" {
					res[i].UserList = append(res[i].UserList, gp.(int64))
				}
			}
		}
	}

	return res, err
}
