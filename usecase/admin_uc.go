package usecase

import (
	"qibla-backend-chat/helper"
	"qibla-backend-chat/model"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/pkg/str"
	"qibla-backend-chat/server/request"
	"qibla-backend-chat/usecase/viewmodel"
	"errors"
	"time"
)

// AdminUC ...
type AdminUC struct {
	*ContractUC
}

// GenerateCode randomize code & check uniqueness from DB
func (uc AdminUC) GenerateCode() (res string, err error) {
	adminModel := model.NewAdminModel(uc.DB)
	res = str.RandAlphanumericString(8)
	for {
		data, _ := adminModel.FindByCode(res)
		if data.ID == "" {
			break
		}
		res = str.RandAlphanumericString(8)
	}

	return res, err
}

// Login ...
func (uc AdminUC) Login(data request.AdminLoginRequest) (res viewmodel.AdminLoginVM, err error) {
	ctx := "AdminUC.Login"

	// Decrypt password input
	data.Password, err = uc.AesFront.Decrypt(data.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "decrypt", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	admin, err := uc.FindByEmail(data.Email)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_by_email", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}

	if admin.Password != data.Password {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "invalid_password", uc.ReqID)
		return res, errors.New(helper.InvalidCredentials)
	}
	if !admin.Status {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "inactive_admin", uc.ReqID)
		return res, errors.New(helper.InactiveAdmin)
	}

	// Jwe the payload & Generate jwt token
	payload := map[string]interface{}{
		"id":   admin.ID,
		"role": "admin",
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

// FindAll ...
func (uc AdminUC) FindAll(page, limit int) (res []viewmodel.AdminVM, pagination viewmodel.PaginationVM, err error) {
	ctx := "AdminUC.FindAll"

	limit = uc.LimitMax(limit)
	limit, offset := uc.PaginationPageOffset(page, limit)

	adminModel := model.NewAdminModel(uc.DB)
	data, count, err := adminModel.FindAll(offset, limit)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, pagination, err
	}
	pagination = PaginationRes(page, count, limit)

	for _, r := range data {
		res = append(res, viewmodel.AdminVM{
			ID:        r.ID,
			Code:      r.Code.String,
			Name:      r.Name.String,
			Email:     r.Email.String,
			RoleID:    r.RoleID.String,
			RoleCode:  r.Role.Code.String,
			RoleName:  r.Role.Name.String,
			Status:    r.Status.Bool,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
			DeletedAt: r.DeletedAt.String,
		})
	}

	return res, pagination, err
}

// FindByID ...
func (uc AdminUC) FindByID(id string) (res viewmodel.AdminVM, err error) {
	ctx := "AdminUC.FindByID"

	adminModel := model.NewAdminModel(uc.DB)
	data, err := adminModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.AdminVM{
		ID:        data.ID,
		Code:      data.Code.String,
		Name:      data.Name.String,
		Email:     data.Email.String,
		Password:  uc.Aes.DecryptNoErr(data.Password.String),
		RoleID:    data.RoleID.String,
		RoleCode:  data.Role.Code.String,
		RoleName:  data.Role.Name.String,
		Status:    data.Status.Bool,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// FindByCode ...
func (uc AdminUC) FindByCode(code string) (res viewmodel.AdminVM, err error) {
	ctx := "AdminUC.FindByCode"

	adminModel := model.NewAdminModel(uc.DB)
	data, err := adminModel.FindByCode(code)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.AdminVM{
		ID:        data.ID,
		Code:      data.Code.String,
		Name:      data.Name.String,
		Email:     data.Email.String,
		Password:  uc.Aes.DecryptNoErr(data.Password.String),
		RoleID:    data.RoleID.String,
		RoleCode:  data.Role.Code.String,
		RoleName:  data.Role.Name.String,
		Status:    data.Status.Bool,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// FindByEmail ...
func (uc AdminUC) FindByEmail(email string) (res viewmodel.AdminVM, err error) {
	ctx := "AdminUC.FindByEmail"

	adminModel := model.NewAdminModel(uc.DB)
	data, err := adminModel.FindByEmail(email)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	res = viewmodel.AdminVM{
		ID:        data.ID,
		Code:      data.Code.String,
		Name:      data.Name.String,
		Email:     data.Email.String,
		Password:  uc.Aes.DecryptNoErr(data.Password.String),
		RoleID:    data.RoleID.String,
		RoleCode:  data.Role.Code.String,
		RoleName:  data.Role.Name.String,
		Status:    data.Status.Bool,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// Create ...
func (uc AdminUC) Create(data request.AdminRequest) (res viewmodel.AdminVM, err error) {
	ctx := "AdminUC.Create"

	// Check duplicate email
	adminModel := model.NewAdminModel(uc.DB)
	admin, _ := adminModel.FindByEmail(data.Email)
	if admin.ID != "" {
		logruslogger.Log(logruslogger.WarnLevel, data.Email, ctx, "duplicate_email", uc.ReqID)
		return res, errors.New(helper.DuplicateEmail)
	}

	// Decrypt password input
	data.Password, err = uc.AesFront.Decrypt(data.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "decrypt", uc.ReqID)
		return res, err
	}

	// Generate code
	code, err := uc.GenerateCode()
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "generate_code", uc.ReqID)
		return res, err
	}

	// Encrypt password
	password, err := uc.Aes.Encrypt(data.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "encrypt_password", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.AdminVM{
		Code:      code,
		Name:      data.Name,
		Email:     data.Email,
		Password:  password,
		RoleID:    data.RoleID,
		Status:    data.Status,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	res.ID, err = adminModel.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Update ...
func (uc AdminUC) Update(id string, data request.AdminRequest) (res viewmodel.AdminVM, err error) {
	ctx := "AdminUC.Update"

	// Check duplicate email
	adminModel := model.NewAdminModel(uc.DB)
	admin, _ := adminModel.FindByEmail(data.Email)
	if admin.ID != "" && admin.ID != id {
		logruslogger.Log(logruslogger.WarnLevel, data.Email, ctx, "duplicate_email", uc.ReqID)
		return res, errors.New(helper.DuplicateEmail)
	}

	// Decrypt password input
	data.Password, err = uc.AesFront.Decrypt(data.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "decrypt", uc.ReqID)
		return res, err
	}

	// Encrypt password
	password, err := uc.Aes.Encrypt(data.Password)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "encrypt_password", uc.ReqID)
		return res, err
	}

	now := time.Now().UTC()
	res = viewmodel.AdminVM{
		Name:      data.Name,
		Email:     data.Email,
		Password:  password,
		RoleID:    data.RoleID,
		Status:    data.Status,
		UpdatedAt: now.Format(time.RFC3339),
	}
	res.ID, err = adminModel.Update(id, res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}

// Delete ...
func (uc AdminUC) Delete(id string) (res viewmodel.AdminVM, err error) {
	ctx := "AdminUC.Delete"

	now := time.Now().UTC()
	adminModel := model.NewAdminModel(uc.DB)
	res.ID, err = adminModel.Destroy(id, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "query", uc.ReqID)
		return res, err
	}

	return res, err
}
