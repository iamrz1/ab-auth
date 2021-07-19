package middleware

import (
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
	"strings"
)

func JWTTokenOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn == "" {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)
		r.Header.Set(utils.AuthorizationKey, jwtTkn)

		next.ServeHTTP(w, r)
	})
}

func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn == "" {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

		jwtTkn = stripBearerFromToken(jwtTkn)

		claims, err := utils.VerifyToken(jwtTkn, false)
		if err != nil {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, err.Error()))
			return
		}

		if !isTokenFresh(claims.Username, claims.IssuedAt) {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Session expired"))
			return
		}

		r.Header.Set(utils.UsernameKey, claims.Username)

		next.ServeHTTP(w, r)
	})
}

func isTokenFresh(username string, issueTime int64) bool {
	lastResetAt, err := utils.GetLastResetAt(username)
	if err != nil {
		return false
	}

	if lastResetAt > issueTime {
		return false
	}

	return true

}

func stripBearerFromToken(token string) string {
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
		//log.Println("token contains Bearer")
	}

	if strings.HasPrefix(token, "bearer ") {
		token = strings.TrimPrefix(token, "bearer ")
		//log.Println("token contains bearer")
	}

	return token
}
