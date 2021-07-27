package public

import (
	"encoding/json"
	"github.com/go-chi/chi"
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

func (pr *publicRouter) merchantRouter() *chi.Mux {
	r := chi.NewRouter()
	cr := newMerchantRouter(pr.Services, pr.Log)

	r.Post("/signup", cr.signup)
	r.Post("/verify-signup", cr.verifySignUp)
	r.Post("/login", cr.login)
	r.Post("/forgot-password", cr.forgotPassword)
	r.Post("/set-password", cr.setPassword)
	return r
}

// signup godoc
// @Summary Signup a new customer
// @Description Signup a new customer for a valid non-existing phone number
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param  Body body model.MerchantSignupReq true "All fields are mandatory"
// @Success 201 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/signup [post]
func (pr *merchantRouter) signup(w http.ResponseWriter, r *http.Request) {
	req := model.MerchantSignupReq{}

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

	if pr.Services.MerchantService.Config.Environment == utils.EnvProduction || req.CaptchaValue != utils.DefaultCaptchaValue {
		_, err = utils.VerifyCaptcha(req.CaptchaID, req.CaptchaValue)
		if err != nil {
			utils.HandleObjectError(w, rest_error.NewValidationError("", err))
			return
		}
	}

	otp, err := pr.Services.MerchantService.CreateMerchant(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	var meta map[string]string

	if pr.Services.MerchantService.Config.Environment != utils.EnvProduction {
		meta = map[string]string{"otp": otp}
	}

	utils.ServeJSONObject(w, http.StatusCreated, "OTP sent", nil, meta, true)
}

// verifySignUp godoc
// @Summary Verify a new customer using otp
// @Description Use customer defined otp to match it with existing reference in cache to verify a signup
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param  Body body model.MerchantSignupVerificationReq true "All fields are mandatory"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/verify-signup [post]
func (pr *merchantRouter) verifySignUp(w http.ResponseWriter, r *http.Request) {
	req := model.MerchantSignupVerificationReq{}

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

	err = pr.Services.MerchantService.VerifyMerchantSignUp(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Verified", nil, nil, true)
}

// login godoc
// @Summary Login as a customer
// @Description Login uses customer defined username and password to authenticate a customer.
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param  Body body model.LoginReq true "All fields are mandatory"
// @Success 200 {object} response.TokenSuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/login [post]
func (pr *merchantRouter) login(w http.ResponseWriter, r *http.Request) {
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

	res, err := pr.Services.MerchantService.Login(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Logged in", res, nil, true)
}

// forgotPassword godoc
// @Summary Request OTP to reset password
// @Description Use username and captcha to send otp to customer's registered number
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param  Body body model.ForgotPasswordReq true "All fields are mandatory"
// @Success 201 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/forgot-password [post]
func (pr *merchantRouter) forgotPassword(w http.ResponseWriter, r *http.Request) {
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

	if pr.Services.MerchantService.Config.Environment == utils.EnvProduction || req.CaptchaValue != utils.DefaultCaptchaValue {
		_, err = utils.VerifyCaptcha(req.CaptchaID, req.CaptchaValue)
		if err != nil {
			utils.HandleObjectError(w, rest_error.NewValidationError("", err))
			return
		}
	}

	otp, err := pr.Services.MerchantService.ForgotPassword(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	var meta map[string]string

	if pr.Services.MerchantService.Config.Environment != utils.EnvProduction {
		meta = map[string]string{"otp": otp}
	}

	utils.ServeJSONObject(w, http.StatusCreated, "OTP sent", nil, meta, true)
}

// setPassword godoc
// @Summary Set customer's password with OTP
// @Description Set new password using OTP received during forgot-password
// @Tags Merchants
// @Accept  json
// @Produce  json
// @Param  Body body model.SetPasswordReq true "All fields are mandatory"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 400 {object} response.EmptyErrorRes
// @Failure 404 {object} response.EmptyErrorRes
// @Failure 500 {object} response.EmptyErrorRes
// @Router /api/v1/public/customers/set-password [post]
func (pr *merchantRouter) setPassword(w http.ResponseWriter, r *http.Request) {
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

	err = pr.Services.MerchantService.SetPassword(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Password set", nil, nil, true)
}
