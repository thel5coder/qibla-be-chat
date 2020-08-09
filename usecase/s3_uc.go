package usecase

import (
	"mime/multipart"
	"qibla-backend-chat/pkg/logruslogger"
)

// S3UC ...
type S3UC struct {
	*ContractUC
}

// UploadFile ...
func (uc S3UC) UploadFile(path string, file *multipart.FileHeader) (res string, err error) {
	ctx := "S3UC.UploadFile"

	res, err = uc.S3.Upload(path, file)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "upload", uc.ReqID)
		return res, err
	}

	return res, err
}

// GetURL ...
func (uc S3UC) GetURL(key string) (res string, err error) {
	ctx := "S3UC.GetURL"

	res, err = uc.S3.GetURL(key)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "get_url", uc.ReqID)
		return res, err
	}

	return res, err
}

// GetURLNoErr ...
func (uc S3UC) GetURLNoErr(key string) (res string) {
	res, _ = uc.GetURL(key)

	return res
}

// Delete ...
func (uc S3UC) Delete(key string) (err error) {
	ctx := "S3UC.Delete"

	_, err = uc.S3.Delete(key)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "delete", uc.ReqID)
		return err
	}

	return err
}
