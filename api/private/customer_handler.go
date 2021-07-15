package private

import (
	"encoding/json"
	_ "github.com/iamrz1/ab-auth/docs"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

// VerifyAccessToken godoc
// @Summary Verify customer's access token
// @Description VerifyAccessToken lets apps to verify that a provided token is in-fact valid
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Value of access token"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 401 {object} response.CustomerErrorRes
// @Router /api/v1/private/customers/verify-token [get]
func (pr *privateRouter) VerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	utils.ServeJSONObject(w, "", http.StatusOK, "Token verified", nil, nil, true)
}

// GetCustomerProfile godoc
// @Summary Get basic profile
// @Description Returns user's profile using access token
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 200 {object} response.CustomerSuccessRes
// @Failure 400 {object} response.CustomerErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.CustomerErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.CustomerErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/profile [get]
func (pr *privateRouter) GetCustomerProfile(w http.ResponseWriter, r *http.Request) {
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

	utils.ServeJSONObject(w, "", http.StatusOK, "Successful", &data, nil, true)
}

// UpdateCustomerProfile godoc
// @Summary Update basic profile
// @Description Update user's basic profile info
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.CustomerProfileUpdateReq true "Some fields are mandatory"
// @Success 200 {object} response.CustomerSuccessRes
// @Failure 400 {object} response.CustomerErrorRes
// @Failure 404 {object} response.CustomerErrorRes
// @Failure 500 {object} response.CustomerErrorRes
// @Router /api/v1/private/customers/profile [patch]
func (pr *privateRouter) UpdateCustomerProfile(w http.ResponseWriter, r *http.Request) {
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

	utils.ServeJSONObject(w, "", http.StatusOK, "Profile updated", &data, nil, true)
}

// RefreshToken godoc
// @Summary Refresh customer's access token
// @Description Generate new access and refresh tokens using current refresh token
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param authorization header string true "Value of refresh token"
// @Success 200 {object} response.TokenSuccessRes
// @Failure 401 {object} response.CustomerErrorRes
// @Router /api/v1/private/customers/refresh-token [get]
func (pr *privateRouter) RefreshToken(w http.ResponseWriter, r *http.Request) {
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

	utils.ServeJSONObject(w, "", http.StatusOK, "Token refreshed", &token, nil, true)
}

//func (pr *privateRouter) PurgeCustomer(w http.ResponseWriter, r *http.Request) {
//	slug := chi.URLParam(r, "slug")
//
//	req := model.CustomerDeleteReq{Username: slug}
//
//	data, err := pr.Services.CustomerService.PurgeCustomer(r.Context(), &req)
//	if err != nil {
//		utils.HandleObjectError(w, err)
//		return
//	}
//
//	utils.ServeJSONObject(w, "", http.StatusOK, "Purged generic object successfully", &data, nil, true)
//}
