package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"qibla-backend-chat/usecase/viewmodel"
	"strings"

	"qibla-backend-chat/pkg/jwe"
	"qibla-backend-chat/pkg/jwt"
	"qibla-backend-chat/pkg/str"
	"qibla-backend-chat/usecase"

	"database/sql"
	"github.com/go-playground/universal-translator"
	validator "gopkg.in/go-playground/validator.v9"
)

// Handler ...
type Handler struct {
	ContractUC *usecase.ContractUC
	DB         *sql.DB
	EnvConfig  map[string]string
	Validate   *validator.Validate
	Translator ut.Translator
	Jwt        jwt.Credential
	Jwe        jwe.Credential
}

// Bind bind the API request payload (body) into request struct.
func (h Handler) Bind(r *http.Request, input interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&input)

	return err
}

// emptyJSONArr ...
func emptyJSONArr() []map[string]interface{} {
	return []map[string]interface{}{}
}

// SendSuccess send success into response with 200 http code.
func SendSuccess(w http.ResponseWriter, payload interface{}, meta interface{}) {
	RespondWithJSON(w, 200, 200, "Success", payload, meta)
}

// SendBadRequest send bad request into response with 400 http code.
func SendBadRequest(w http.ResponseWriter, message string) {
	RespondWithJSON(w, 400, 400, message, emptyJSONArr(), emptyJSONArr())
}

// SendRequestValidationError Send validation error response to consumers.
func (h Handler) SendRequestValidationError(w http.ResponseWriter, validationErrors validator.ValidationErrors) {
	errorResponse := map[string][]string{}
	errorTranslation := validationErrors.Translate(h.Translator)
	for _, err := range validationErrors {
		errKey := str.Underscore(err.StructField())
		errorResponse[errKey] = append(
			errorResponse[errKey],
			strings.Replace(errorTranslation[err.Namespace()], err.StructField(), "[]", -1),
		)
	}

	RespondWithJSON(w, 400, 405, "validation error", errorResponse, emptyJSONArr())
}

// RespondWithJSON write json response format
func RespondWithJSON(w http.ResponseWriter, httpCode int, statCode int, message string, payload interface{}, meta interface{}) {
	respPayload := map[string]interface{}{
		"stat_code": statCode,
		"stat_msg":  message,
		"meta":      meta,
		"data":      payload,
	}

	response, _ := json.Marshal(respPayload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(response)
}

// requestIDFromContextInterface ...
func requestIDFromContextInterface(ctx context.Context, key string) map[string]interface{} {
	return ctx.Value(key).(map[string]interface{})
}

// getHeaderReqID ...
func getHeaderReqID(r *http.Request) (res string) {
	if len(r.Header["req_id"]) > 0 {
		res = r.Header["req_id"][0]
	}

	return res
}

// getUserDetail ...
func getUserDetail(r *http.Request) (res viewmodel.UserVM) {
	user := requestIDFromContextInterface(r.Context(), "user")

	res = viewmodel.UserVM{
		ID:                user["id"].(string),
		Username:          user["username"].(string),
		Email:             user["email"].(string),
		Name:              user["name"].(string),
		IsActive:          user["is_active"].(bool),
		ProfilePicture:    user["profile_picture"].(string),
		ProfilePictureURL: user["profile_picture_url"].(string),
		RoleID:            user["role_id"].(string),
		RoleName:          user["role_name"].(string),
		OdooUserID:        user["odoo_user_id"].(int64),
	}

	return res
}
