package middleware

import (
	"qibla-backend-chat/model"
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

// VerifyJwtTokenCredential ...
func (m VerifyMiddlewareInit) VerifyJwtTokenCredential(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims := &jwtClaims{}
		claims := jwt.MapClaims{}

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

		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			apiHandler.RespondWithJSON(w, 401, 401, "Expired Token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Check id in user table
		userUc := usecase.UserUC{ContractUC: m.ContractUC}
		user, err := userUc.FindByID(int(claims["id"].(float64)))
		if err != nil {
			apiHandler.RespondWithJSON(w, 401, 401, "Invalid user token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		claimRes := map[string]interface{}{
			"id":         user.ID,
			"company_id": user.CompanyID,
		}

		ctx := userContextInterface(r.Context(), r, "user", claimRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifySuperadminTokenCredential ...
func (m VerifyMiddlewareInit) VerifySuperadminTokenCredential(next http.Handler) http.Handler {
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

		// Check if the token provided is a valid admin token
		if jweRes["role"] == nil {
			apiHandler.RespondWithJSON(w, 401, 401, "Invalid admin token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}
		if jweRes["role"].(string) != "admin" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not an admin token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		adminID := jweRes["id"].(string)

		// Check id in table
		adminUc := usecase.AdminUC{ContractUC: m.ContractUC}
		admin, err := adminUc.FindByID(adminID)
		if admin.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		// Check if role is superadmin
		if admin.RoleCode != model.RoleCodeSuperadmin {
			apiHandler.RespondWithJSON(w, 401, 401, "Invalid Role!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["code"] = admin.Code
		jweRes["name"] = admin.Name
		jweRes["roleCode"] = admin.RoleCode
		jweRes["roleName"] = admin.RoleName

		ctx := userContextInterface(r.Context(), r, "admin", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// VerifyAdminTokenCredential ...
func (m VerifyMiddlewareInit) VerifyAdminTokenCredential(next http.Handler) http.Handler {
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

		// Check if the token provided is a valid admin token
		if jweRes["role"] == nil {
			apiHandler.RespondWithJSON(w, 401, 401, "Invalid admin token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}
		if jweRes["role"].(string) != "admin" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not an admin token!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		adminID := jweRes["id"].(string)

		// Check id in table
		adminUc := usecase.AdminUC{ContractUC: m.ContractUC}
		admin, err := adminUc.FindByID(adminID)
		if admin.ID == "" {
			apiHandler.RespondWithJSON(w, 401, 401, "Not found!", []map[string]interface{}{}, []map[string]interface{}{})
			return
		}

		jweRes["code"] = admin.Code
		jweRes["name"] = admin.Name
		jweRes["roleCode"] = admin.RoleCode
		jweRes["roleName"] = admin.RoleName

		ctx := userContextInterface(r.Context(), r, "admin", jweRes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
