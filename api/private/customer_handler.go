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

func newCustomerRouter(svc *service.Config, rLogger rLog.Logger) *customerRouter {
	return &customerRouter{
		Services: svc,
		Log:      rLogger,
	}
}

type customerRouter struct {
	Services *service.Config
	Log      rLog.Logger
}

func (pr *privateRouter) customerRouter() *chi.Mux {
	r := chi.NewRouter()

	cr := newCustomerRouter(pr.Services, pr.Log)

	r.With(middleware.AuthenticatedCustomerOnly).Get("/profile", cr.getCustomerProfile)
	r.With(middleware.AuthenticatedCustomerOnly).Patch("/profile", cr.updateCustomerProfile)
	r.With(middleware.AuthenticatedCustomerOnly).Get("/verify-token", cr.verifyAccessToken)
	r.With(middleware.JWTTokenOnly).Get("/refresh-token", cr.refreshToken)
	r.With(middleware.AuthenticatedCustomerOnly).Put("/password", cr.updatePassword)

	r.Mount("/address", pr.addressRouter())

	return r
}

// verifyAccessToken godoc
// @Summary Verify customer's access token
// @Description verifyAccessToken lets apps to verify that a provided token is in-fact valid
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Value of access token"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 401 {object} response.EmptyErrorRes
// @Router /api/v1/private/customers/verify-token [get]
func (pr *customerRouter) verifyAccessToken(w http.ResponseWriter, r *http.Request) {
	utils.ServeJSONObject(w, http.StatusOK, "Token verified", nil, nil, true)
}

// getCustomerProfile godoc
// @Summary Get basic profile
// @Description Returns customer's profile using access token
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 200 {object} response.CustomerSuccessRes
// @Failure 400 {object} response.EmptyErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/profile [get]
func (pr *customerRouter) getCustomerProfile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	req := &model.Customer{Username: username}

	data, err := pr.Services.CustomerService.GetCustomer(r.Context(), req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Successful", &data, nil, true)
}

// updateCustomerProfile godoc
// @Summary Update basic profile
// @Description Update customer's basic profile info
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.CustomerProfileUpdateReq true "Some fields are mandatory"
// @Success 200 {object} response.CustomerSuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/private/customers/profile [patch]
func (pr *customerRouter) updateCustomerProfile(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	req := model.CustomerProfileUpdateReq{}
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

	data, err := pr.Services.CustomerService.UpdateCustomer(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Profile updated", &data, nil, true)
}

// refreshToken godoc
// @Summary Refresh customer's access token
// @Description Generate new access and refresh tokens using current refresh token
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Value of refresh token"
// @Success 200 {object} response.TokenSuccessRes
// @Failure 401 {object} response.EmptyErrorRes
// @Router /api/v1/private/customers/refresh-token [get]
func (pr *customerRouter) refreshToken(w http.ResponseWriter, r *http.Request) {
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

	cus, err := pr.Services.CustomerService.GetCustomer(r.Context(), &model.Customer{Username: claims.Username})
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, err.Error()))
		return
	}

	if claims.IssuedAt < cus.LastResetAt.Unix() {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Session expired"))
		return
	}

	access, refresh := utils.GenerateTokens(cus.Username, "", "customer")
	token := model.Token{AccessToken: access, RefreshToken: refresh}

	utils.ServeJSONObject(w, http.StatusOK, "Token refreshed", &token, nil, true)
}

// updatePassword godoc
// @Summary Update existing password
// @Description Update to a new password using customer's existing password
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.UpdatePasswordReq true "Some fields are mandatory"
// @Success 200 {object} response.CustomerSuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/private/customers/password [put]
func (pr *customerRouter) updatePassword(w http.ResponseWriter, r *http.Request) {
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

	data, err := pr.Services.CustomerService.UpdatePassword(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Password updated", &data, nil, true)
}

func (pr *customerRouter) purgeCustomer(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	req := model.CustomerDeleteReq{Username: username}

	data, err := pr.Services.CustomerService.PurgeCustomer(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Purged customer successfully", &data, nil, true)
}
