package usecase

import (
	"errors"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/model"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/usecase/viewmodel"
)

// UserUC ...
type UserUC struct {
	*ContractUC
}

// FindByID ...
func (uc UserUC) FindByID(id string) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.FindByID"

	userModel := model.NewUserModel(uc.DB)
	data, err := userModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.UserVM{
		ID:         data.ID,
		Username:   data.Username.String,
		Email:      data.Email.String,
		Name:       data.Name.String,
		Password:   data.Password.String,
		RoleID:     data.RoleID.String,
		RoleName:   data.RoleName.String,
		OdooUserID: data.OdoUserID.Int64,
		IsActive:   data.IsActive.Bool,
		CreatedAt:  data.CreatedAt.String,
		UpdatedAt:  data.UpdatedAt.String,
		DeletedAt:  data.DeletedAt.String,
	}

	return res, err
}

// FindByOdooUserID ...
func (uc UserUC) FindByOdooUserID(odooUserID int64) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.FindByOdooUserID"

	userModel := model.NewUserModel(uc.DB)
	data, err := userModel.FindByOdooUserID(odooUserID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.UserVM{
		ID:         data.ID,
		Username:   data.Username.String,
		Email:      data.Email.String,
		Name:       data.Name.String,
		Password:   data.Password.String,
		RoleID:     data.RoleID.String,
		RoleName:   data.RoleName.String,
		OdooUserID: data.OdoUserID.Int64,
		IsActive:   data.IsActive.Bool,
		CreatedAt:  data.CreatedAt.String,
		UpdatedAt:  data.UpdatedAt.String,
		DeletedAt:  data.DeletedAt.String,
	}

	return res, err
}

// FindTravelPackage ...
func (uc UserUC) FindTravelPackage(id string) (res []viewmodel.TravelPackageVM, err error) {
	ctx := "UserUC.FindTravelPackage"

	user, err := uc.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	// Get detail user to get travel package id
	odooUc := OdooUC{ContractUC: uc.ContractUC}
	partner, err := odooUc.FindByIDPartner(user.OdooUserID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_partner", uc.ReqID)
		return res, err
	}

	// Getall  travel package by user
	res, err = odooUc.FindAllPackage(partner.PackageListID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_package", uc.ReqID)
		return res, err
	}

	return res, err
}

// FindJamaah ...
func (uc UserUC) FindJamaah(id string, travelPackageID int64) (res []viewmodel.UserVM, err error) {
	ctx := "UserUC.FindJamaah"

	_, err = uc.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	// Get detail user to get travel package id
	odooUc := OdooUC{ContractUC: uc.ContractUC}
	travelPackage, err := odooUc.FindByIDTravelPackage(travelPackageID)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_travel_package", uc.ReqID)
		return res, err
	}

	for _, t := range travelPackage.UserList {
		userTemp, _ := uc.FindByOdooUserID(t)
		if userTemp.ID != "" && userTemp.ID != id {
			res = append(res, userTemp)
		}
	}

	return res, err
}

// Login ...
func (uc UserUC) Login(id string) (res viewmodel.LoginVM, err error) {
	ctx := "AdminUC.Login"

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"id": id,
	}
	jwePayload, err := uc.ContractUC.Jwe.Generate(payload)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwe", uc.ReqID)
		return res, errors.New(helper.JWT)
	}
	res.Token, res.ExpiredDate, err = uc.ContractUC.Jwt.GetToken(jwePayload)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "jwt", uc.ReqID)
		return res, errors.New(helper.JWT)
	}

	return res, err
}
