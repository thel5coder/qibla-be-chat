package usecase

import (
	"qibla-backend-chat/model"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/usecase/viewmodel"
)

// RoleUC ...
type RoleUC struct {
	*ContractUC
}

// FindByID ...
func (uc RoleUC) FindByID(id string) (res viewmodel.RoleVM, err error) {
	ctx := "RoleUC.FindByID"

	roleModel := model.NewRoleModel(uc.DB)
	data, err := roleModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.RoleVM{
		ID:        data.ID,
		Code:      data.Code.String,
		Name:      data.Name.String,
		Status:    data.Status.Bool,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// FindByCode ...
func (uc RoleUC) FindByCode(code string) (res viewmodel.RoleVM, err error) {
	ctx := "RoleUC.FindByCode"

	roleModel := model.NewRoleModel(uc.DB)
	data, err := roleModel.FindByCode(code)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.RoleVM{
		ID:        data.ID,
		Code:      data.Code.String,
		Name:      data.Name.String,
		Status:    data.Status.Bool,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}
