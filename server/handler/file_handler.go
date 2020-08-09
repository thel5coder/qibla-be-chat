package handler

import (
	"net/http"
	"qibla-backend-chat/helper"
	"qibla-backend-chat/model"
	"qibla-backend-chat/pkg/str"
	"qibla-backend-chat/usecase"
)

// FileHandler ...
type FileHandler struct {
	Handler
}

// UploadHandler ...
func (h *FileHandler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	user := requestIDFromContextInterface(r.Context(), "user")
	userID := user["id"].(string)

	// Get logrus request ID
	h.ContractUC.ReqID = getHeaderReqID(r)

	// Check file size
	maxUploadSize := str.StringToInt(h.Handler.EnvConfig["FILE_MAX_UPLOAD_SIZE"])
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxUploadSize))
	err := r.ParseMultipartForm(int64(maxUploadSize))
	if err != nil {
		SendBadRequest(w, helper.FileTooBig)
		return
	}

	// Read file type
	fileType := r.PostFormValue("type")
	if fileType != model.FileGroupRoomProfilePicture {
		SendBadRequest(w, helper.InvalidFileType)
		return
	}

	// Read file
	file, header, err := r.FormFile("file")
	if err != nil {
		SendBadRequest(w, helper.FileError)
		return
	}
	defer file.Close()

	// Upload file to local temporary
	fileUc := usecase.FileUC{ContractUC: h.ContractUC}
	res, err := fileUc.Upload(fileType, userID, header)
	if err != nil {
		SendBadRequest(w, err.Error())
		return
	}

	SendSuccess(w, res, nil)
	return
}
