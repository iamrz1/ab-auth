package private

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/api/middleware"
	_ "github.com/iamrz1/ab-auth/docs"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/service"
	"github.com/iamrz1/ab-auth/utils"
	rLog "github.com/iamrz1/rest-log"
	"net/http"
)

func newMerchantRouter(svc *service.Config, rLogger rLog.Logger) *merchantRouter {
	return &merchantRouter{
		Services: svc,
		Log:      rLogger,
	}
}

type merchantRouter struct {
	Services *service.Config
	Log      rLog.Logger
}

func (pr *privateRouter) merchantRouter() *chi.Mux {
	r := chi.NewRouter()

	cr := newMerchantRouter(pr.Services, pr.Log)

	r.With(middleware.AuthenticatedMerchantOnly).Get("/profile", cr.getMerchantProfile)
	r.With(middleware.AuthenticatedMerchantOnly).Patch("/profile", cr.updateMerchantProfile)
	r.With(middleware.AuthenticatedMerchantOnly).Get("/verify-token", cr.verifyAccessToken)
	r.With(middleware.JWTTokenOnly).Get("/refresh-token", cr.refreshToken)
	r.With(middleware.AuthenticatedMerchantOnly).Put("/password", cr.updatePassword)

	r.Mount("/address", pr.addressRouter())

	return r
}

// verifyAccessToken godoc
// @Summary Verify merchant's access token
// @Description verifyAccessToken lets apps to verify that a provided token is in-fact valid
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param authorization header string true "Value of access token"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 401 {object} response.EmptyErrorRes
// @Router /api/v1/private/merchants/verify-token [get]
func (pr *merchantRouter) verifyAccessToken(w http.ResponseWriter, r *http.Request) {
	utils.ServeJSONObject(w, http.StatusOK, "Token verified", nil, nil, true)
}

// getMerchantProfile godoc
// @Summary Get basic profile
// @Description Returns merchant's profile using access token
// @Tags Merchants
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 200 {object} response.MerchantSuccessRes
// @Failure 400 {object} response.EmptyErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyErrorRes "API sever or db unreachable."
// @Router /api/v1/private/merchants/profile [get]
func (pr *merchantRouter) getMerchantProfile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	req := &model.Merchant{Username: username}

	data, err := pr.Services.MerchantService.GetMerchant(r.Context(), req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Successful", &data, nil, true)
}

// updateMerchantProfile godoc
// @Summary Update basic profile
// @Description Update merchant's basic profile info
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.MerchantProfileUpdateReq true "Some fields are mandatory"
// @Success 200 {object} response.MerchantSuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/private/merchants/profile [patch]
func (pr *merchantRouter) updateMerchantProfile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	req := model.MerchantProfileUpdateReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}
	req.Username = username

	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("missing required field(s)", err))
		return
	}

	data, err := pr.Services.MerchantService.UpdateMerchant(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Profile updated", &data, nil, true)
}

// refreshToken godoc
// @Summary Refresh merchant's access token
// @Description Generate new access and refresh tokens using current refresh token
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param authorization header string true "Value of refresh token"
// @Success 200 {object} response.TokenSuccessRes
// @Failure 401 {object} response.EmptyErrorRes
// @Router /api/v1/private/merchants/refresh-token [get]
func (pr *merchantRouter) refreshToken(w http.ResponseWriter, r *http.Request) {
	jwtTkn := r.Header.Get(utils.AuthorizationKey)
	if jwtTkn == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing refresh token"))
		return
	}

	claims, err := utils.VerifyToken(jwtTkn, true)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, err.Error()))
		return
	}

	cus, err := pr.Services.MerchantService.GetMerchant(r.Context(), &model.Merchant{Username: claims.Username})
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, err.Error()))
		return
	}

	if claims.IssuedAt < cus.LastResetAt.Unix() {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Session expired"))
		return
	}

	access, refresh := utils.GenerateTokens(cus.Username, "", "merchant")
	token := model.Token{AccessToken: access, RefreshToken: refresh}

	utils.ServeJSONObject(w, http.StatusOK, "Token refreshed", &token, nil, true)
}

// updatePassword godoc
// @Summary Update existing password
// @Description Update to a new password using merchant's existing password
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.UpdatePasswordReq true "Some fields are mandatory"
// @Success 200 {object} response.MerchantSuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/private/merchants/password [put]
func (pr *merchantRouter) updatePassword(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	req := model.UpdatePasswordReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}
	req.Username = username

	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("missing required field(s)", err))
		return
	}

	data, err := pr.Services.MerchantService.UpdatePassword(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Password updated", &data, nil, true)
}

func (pr *merchantRouter) purgeMerchant(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	req := model.MerchantDeleteReq{Username: username}

	data, err := pr.Services.MerchantService.PurgeMerchant(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Purged merchant successfully", &data, nil, true)
}
