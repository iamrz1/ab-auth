package public

import (
	"encoding/json"
	_ "github.com/iamrz1/ab-auth/docs"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

// Signup godoc
// @Summary Signup a new customer
// @Description Signup a new customer for a valid non-existing phone number
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.CustomerSignupReq true "All fields are mandatory"
// @Success 201 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/signup [post]
func (pr *publicRouter) Signup(w http.ResponseWriter, r *http.Request) {
	req := model.CustomerSignupReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}
	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Missing required field(s)", err))
		return
	}

	if pr.Services.CustomerService.Config.Environment == utils.EnvProduction || req.CaptchaValue != utils.DefaultCaptchaValue {
		_, err = utils.VerifyCaptcha(req.CaptchaID, req.CaptchaValue)
		if err != nil {
			utils.HandleObjectError(w, rest_error.NewValidationError("", err))
			return
		}
	}

	otp, err := pr.Services.CustomerService.CreateCustomer(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	var meta map[string]string

	if pr.Services.CustomerService.Config.Environment != utils.EnvProduction {
		meta = map[string]string{"otp": otp}
	}

	utils.ServeJSONObject(w, http.StatusCreated, "OTP sent", nil, meta, true)
}

// VerifySignUp godoc
// @Summary Verify a new customer using otp
// @Description Use user defined otp to match it with existing reference in cache to verify a signup
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.CustomerSignupVerificationReq true "All fields are mandatory"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/verify-signup [post]
func (pr *publicRouter) VerifySignUp(w http.ResponseWriter, r *http.Request) {
	req := model.CustomerSignupVerificationReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}
	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Missing required field(s)", err))
		return
	}

	err = pr.Services.CustomerService.VerifyCustomerSignUp(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Verified", nil, nil, true)
}

// Login godoc
// @Summary Login as a customer
// @Description Login uses user defined username and password to authenticate a user.
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.LoginReq true "All fields are mandatory"
// @Success 200 {object} response.TokenSuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/login [post]
func (pr *publicRouter) Login(w http.ResponseWriter, r *http.Request) {
	req := model.LoginReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}
	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Missing required field(s)", err))
		return
	}

	res, err := pr.Services.CustomerService.Login(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Logged in", res, nil, true)
}

// ForgotPassword godoc
// @Summary Request OTP to reset password
// @Description ForgotPassword uses username and captcha to send otp to user's registered number
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.ForgotPasswordReq true "All fields are mandatory"
// @Success 201 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/forgot-password [post]
func (pr *publicRouter) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	req := model.ForgotPasswordReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}

	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Missing required field(s)", err))
		return
	}

	if pr.Services.CustomerService.Config.Environment == utils.EnvProduction || req.CaptchaValue != utils.DefaultCaptchaValue {
		_, err = utils.VerifyCaptcha(req.CaptchaID, req.CaptchaValue)
		if err != nil {
			utils.HandleObjectError(w, rest_error.NewValidationError("", err))
			return
		}
	}

	otp, err := pr.Services.CustomerService.ForgotPassword(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	var meta map[string]string

	if pr.Services.CustomerService.Config.Environment != utils.EnvProduction {
		meta = map[string]string{"otp": otp}
	}

	utils.ServeJSONObject(w, http.StatusCreated, "OTP sent", nil, meta, true)
}

// SetPassword godoc
// @Summary Set user's password with OTP
// @Description Set new password using OTP received during forgot-password
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.SetPasswordReq true "All fields are mandatory"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/set-password [post]
func (pr *publicRouter) SetPassword(w http.ResponseWriter, r *http.Request) {
	req := model.SetPasswordReq{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}

	err = model.Validate(req)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("Missing required field(s)", err))
		return
	}

	err = pr.Services.CustomerService.SetPassword(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Password set", nil, nil, true)
}

//func (pr *publicRouter) PurgeCustomer(w http.ResponseWriter, r *http.Request) {
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
//	utils.ServeJSONObject(w, http.StatusOK, "Purged generic object successfully", &data, nil, true)
//}
