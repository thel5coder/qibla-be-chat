package middleware

import (
	"context"
	"fmt"
	"strings"

	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	apiHandler "qibla-backend-chat/server/handler"
	"qibla-backend-chat/usecase"
)

type jwtClaims struct {
	jwt.StandardClaims
}

// VerifyMiddlewareInit ...
type VerifyMiddlewareInit struct {
	*usecase.ContractUC
}

// VerifyPermissionInit ...
type VerifyPermissionInit struct {
	*usecase.ContractUC
	Menu string
}

func userContextInterface(ctx context.Context, req *http.Request, subject string, body map[string]interface{}) context.Context {
	return context.WithValue(ctx, subject, body)
}

// VerifyUserTokenCredential ...
func (m VerifyMiddlewareInit) VerifyUserTokenCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &jwtClaims{}

		tokenAuthHeader := r.Header.Get("Authorization")
		if !strings.Contains(tokenAuthHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
		tokenAuth := strings.Replace(tokenAuthHeader, "Bearer ", "", -1)

		_, err := jwt.ParseWithClaims(tokenAuth, claims, func(token *jwt.Token) (interface{}, error) {
			if jwt.SigningMethodHS256 != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			secret := m.ContractUC.EnvConfig["TOKEN_SECRET"]
			return []byte(secret), nil
		})
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, "Invalid Token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		if claims.ExpiresAt < time.Now().Unix() {
			apiHandler.RespondWithJSON(w, 401, 401, "Expired Token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Decrypt payload
		jweRes, err := m.ContractUC.Jwe.Rollback(claims.Id)
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, "Error when load the payload!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		userID := jweRes["id"].(string)

		// Check id in table
		userUc := usecase.UserUC{ContractUC: m.ContractUC}
		user, err := userUc.FindByID(userID)
		if user.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}
		if !user.IsActive {
			apiHandler.RespondWithJSON(w, 401, 401, "Inactive User!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["username"] = user.Username
		jweRes["email"] = user.Email
		jweRes["name"] = user.Name
		jweRes["is_active"] = user.IsActive
		jweRes["role_id"] = user.RoleID
		jweRes["role_name"] = user.RoleName
		jweRes["odoo_user_id"] = user.OdooUserID

		ctx := userContextInterface(r.Context(), r, "user", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
