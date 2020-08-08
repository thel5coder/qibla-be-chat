package usecase

import (
	"qibla-backend-chat/model"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/usecase/viewmodel"
)

// UserUC ...
type UserUC struct {
	*ContractUC
}

// FindByID ...
func (uc UserUC) FindByID(id int) (res viewmodel.UserVM, err error) {
	ctx := "UserUC.FindByID"

	userModel := model.NewUserModel(uc.DB)
	data, err := userModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.UserVM{
		ID:           data.ID,
		CompanyID:    int(data.CompanyID.Int64),
		RoleID:       int(data.RoleID.Int64),
		Name:         data.Name.String,
		Email:        data.Email.String,
		EmailValidAt: data.EmailValidAt.String,
		Phone:        data.Phone.String,
		PhoneValidAt: data.PhoneValidAt.String,
		Password:     data.Password.String,
		Photo:        data.Photo.String,
		Status:       data.Status.Bool,
		CreatedAt:    data.CreatedAt.String,
		UpdatedAt:    data.UpdatedAt.String,
		DeletedAt:    data.DeletedAt.String,
	}

	return res, err
}
