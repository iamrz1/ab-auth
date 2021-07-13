package middleware

import (
	"github.com/iamrz1/ab-auth/utils"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
)

const (
	tokenKey            string = "token"
	authTypeKey         string = "authType"
	tokenOnlyAuthType   string = "tokenOnly"
	PlatformEvaly              = "Evaly"
	ModlueAuthorization        = "Authorization"
	ActionRead                 = "Read"
	ActionWrite                = "Write"
	RequestIDHeader            = "X-Request-Id"
)

func StdLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.String())
		log.Println()
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func NoToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := metadata.AppendToOutgoingContext(r.Context())
		requestId := r.Header.Get(RequestIDHeader)
		if requestId == "" {
			requestId = uuid.NewV4().String()
		}
		ctx = metadata.AppendToOutgoingContext(ctx, RequestIDHeader, requestId)

		jwtTkn := r.Header.Get("Authorization")
		if jwtTkn != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, "Authorization", jwtTkn)
		}

		secretKey := r.Header.Get(utils.KeyForSecretKey)
		if secretKey != "" {
			ctx = metadata.AppendToOutgoingContext(ctx, utils.KeyForSecretKey, secretKey)
		}

		ip := r.Header.Get(utils.RealUserIpKey)
		if ip == "" {
			ip = "0.0.0.0"
		}

		ctx = metadata.AppendToOutgoingContext(ctx, utils.RealUserIpKey, ip)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TokenOnly(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get("Authorization")
		if jwtTkn == "" {
			resp := utils.Response{
				Message: "Missing authorization token",
				Status:  http.StatusUnauthorized,
			}
			resp.ServeJSON(w)
			return
		}

		ctx := metadata.AppendToOutgoingContext(r.Context(), "Authorization", jwtTkn)
		requestId := r.Header.Get(RequestIDHeader)
		if requestId == "" {
			requestId = uuid.NewV4().String()
		}
		ctx = metadata.AppendToOutgoingContext(ctx, RequestIDHeader, requestId)

		ip := r.Header.Get(utils.RealUserIpKey)
		if ip == "" {
			ip = "0.0.0.0"
		}

		ctx = metadata.AppendToOutgoingContext(ctx, utils.RealUserIpKey, ip)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SecretOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretKey := r.Header.Get(utils.KeyForSecretKey)
		if secretKey == "" {
			resp := utils.Response{
				Message: "Missing authorization secret",
				Status:  http.StatusUnauthorized,
			}
			log.Println("Missing authorization secret")
			resp.ServeJSON(w)
			return
		}

		requestId := r.Header.Get(RequestIDHeader)
		if requestId == "" {
			requestId = uuid.NewV4().String()
			r.Header.Set(RequestIDHeader, requestId)
		}

		ip := r.Header.Get(utils.RealUserIpKey)
		if ip == "" {
			ip = "0.0.0.0"
			r.Header.Set(utils.RealUserIpKey, ip)
		}

		next.ServeHTTP(w, r)
	})
}

func AuthenticatedOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtTkn := r.Header.Get("Authorization")
		if jwtTkn == "" {
			resp := utils.Response{
				Message: "Missing authorization token",
				Status:  http.StatusUnauthorized,
			}
			resp.ServeJSON(w)
			return
		}

		//TODO: implement when auth is ready

		requestId := r.Header.Get(RequestIDHeader)
		if requestId == "" {
			requestId = uuid.NewV4().String()
		}

		ip := r.Header.Get(utils.RealUserIpKey)
		if ip == "" {
			ip = "0.0.0.0"
		}

		next.ServeHTTP(w, r)
	})
}
