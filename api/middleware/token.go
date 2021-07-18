package middleware

import (
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

func JWTTokenOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get(utils.AuthorizationKey)
		if jwtTkn == "" {
			utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing access token"))
			return
		}

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
