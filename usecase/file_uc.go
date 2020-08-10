package usecase

import (
	"errors"
	"mime/multipart"
	"qibla-backend-chat/model"
	"qibla-backend-chat/pkg/logruslogger"
	"qibla-backend-chat/usecase/viewmodel"
	"time"
)

// FileUC ...
type FileUC struct {
	*ContractUC
}

// FindOne ...
func (uc FileUC) FindOne(id string) (res viewmodel.FileVM, err error) {
	ctx := "FileUC.FindOne"

	fileModel := model.NewFileModel(uc.DB)
	data, err := fileModel.FindByID(id)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_id", uc.ReqID)
		return res, err
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	tempPath, err := s3Uc.GetURL(data.Path.String)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_url", uc.ReqID)
		return res, err
	}

	res = viewmodel.FileVM{
		ID:        data.ID,
		Type:      data.Type.String,
		Path:      data.Path.String,
		TempPath:  tempPath,
		UserID:    data.UserID.String,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// FindOneUnassigned check if image is unassigned
func (uc FileUC) FindOneUnassigned(id, types, userUpload string) (res viewmodel.FileVM, err error) {
	ctx := "FileUC.FindOneUnassigned"

	fileModel := model.NewFileModel(uc.DB)
	data, err := fileModel.FindUnassignedByID(id, types, userUpload)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_id", uc.ReqID)
		return res, err
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	tempPath, err := s3Uc.GetURL(data.Path.String)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_url", uc.ReqID)
		return res, err
	}

	res = viewmodel.FileVM{
		ID:        data.ID,
		Type:      data.Type.String,
		Path:      data.Path.String,
		TempPath:  tempPath,
		UserID:    data.UserID.String,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// FindOneAssigned check if image is assigned
func (uc FileUC) FindOneAssigned(id, types string) (res viewmodel.FileVM, err error) {
	ctx := "FileUC.FindOneAssigned"

	fileModel := model.NewFileModel(uc.DB)
	data, err := fileModel.FindAssignedByID(id, types)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_id", uc.ReqID)
		return res, errors.New("Invalid " + types + " file")
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	tempPath, err := s3Uc.GetURL(data.Path.String)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_url", uc.ReqID)
		return res, err
	}

	res = viewmodel.FileVM{
		ID:        data.ID,
		Type:      data.Type.String,
		Path:      data.Path.String,
		TempPath:  tempPath,
		UserID:    data.UserID.String,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		DeletedAt: data.DeletedAt.String,
	}

	return res, err
}

// Create ...
func (uc FileUC) Create(types, url, userUpload string, deleteUnusedFile bool) (res viewmodel.FileVM, err error) {
	ctx := "FileUC.Create"

	fileModel := model.NewFileModel(uc.DB)

	// Delete all unused files first
	if deleteUnusedFile {
		err = uc.DeleteAllUnused(userUpload, types)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete_unused", uc.ReqID)
			return res, err
		}
	}

	now := time.Now().UTC()
	res = viewmodel.FileVM{
		Type:      types,
		Path:      url,
		UserID:    userUpload,
		CreatedAt: now.Format(time.RFC3339),
		UpdatedAt: now.Format(time.RFC3339),
	}
	res.ID, err = fileModel.Store(res, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "create", uc.ReqID)
		return res, err
	}

	// Get temp url
	s3Uc := S3UC{ContractUC: uc.ContractUC}
	res.TempPath, err = s3Uc.GetURL(res.Path)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_url", uc.ReqID)
		return res, err
	}

	return res, err
}

// Delete ...
func (uc FileUC) Delete(id string) (res viewmodel.FileVM, err error) {
	ctx := "FileUC.Delete"

	now := time.Now().UTC()
	fileModel := model.NewFileModel(uc.DB)
	res.ID, err = fileModel.Destroy(id, now)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete", uc.ReqID)
		return res, err
	}

	return res, err
}

// DeleteAllUnused ...
func (uc FileUC) DeleteAllUnused(userID, types string) (err error) {
	ctx := "FileUC.DeleteAllUnused"
	fileModel := model.NewFileModel(uc.DB)
	unusedFile, err := fileModel.FindAllByUserID(userID, types)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "find_all_unused", uc.ReqID)
		return err
	}

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	roomUc := RoomUC{ContractUC: uc.ContractUC}
	now := time.Now().UTC()
	for _, r := range unusedFile {
		// Find file used
		user, _ := roomUc.FindByProfilePicture(userID, r.Path.String)
		if user.ID != "" {
			continue
		}

		err = s3Uc.Delete(r.Path.String)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "s3", uc.ReqID)
		}

		_, err = fileModel.Destroy(r.ID, now)
		if err != nil {
			logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "user_file", uc.ReqID)
		}
	}
	err = nil

	return err
}

// Upload ...
func (uc FileUC) Upload(types, userID string, file *multipart.FileHeader, deleteUnusedFile bool) (res viewmodel.FileVM, err error) {
	ctx := "FileUC.Upload"

	s3Uc := S3UC{ContractUC: uc.ContractUC}
	fileKey, err := s3Uc.UploadFile(types+"/"+userID, file)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "upload_file", uc.ReqID)
		return res, err
	}

	res, err = uc.Create(types, fileKey, userID, deleteUnusedFile)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "create", uc.ReqID)
		return res, err
	}

	return res, err
}
