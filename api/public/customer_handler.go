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
// @Description Signup a new new customer for a valid non-existing phone number
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.CustomerSignupReq true "All fields are mandatory"
// @Success 201 {object} response.EmptySuccessRes
// @Failure 400 {object} response.CustomerErrorRes
// @Failure 404 {object} response.CustomerErrorRes
// @Failure 500 {object} response.CustomerErrorRes
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

	_, err = utils.VerifyCaptcha(req.CaptchaID, req.CaptchaValue)
	if err != nil {
		utils.HandleObjectError(w, rest_error.NewValidationError("", err))
		return
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

	//otp, requestID, err := utils.GenerateTimedRandomDigits(req.Username,5,pr.Services.CustomerService.Config.OtpTtlMinutes)

	utils.ServeJSONObject(w, "", http.StatusCreated, "OTP sent", nil, meta, true)
}

// VerifySignUp godoc
// @Summary Verify a new customer using otp
// @Description VerifySignUp uses user defined otp and matches it with existing reference in cache to verify a signup
// @Tags Customers
// @Accept  json
// @Produce  json
// @Param  Body body model.CustomerSignupVerificationReq true "All fields are mandatory"
// @Success 200 {object} response.EmptySuccessRes
// @Failure 400 {object} response.CustomerErrorRes
// @Failure 404 {object} response.CustomerErrorRes
// @Failure 500 {object} response.CustomerErrorRes
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

	var meta map[string]string

	//otp, requestID, err := utils.GenerateTimedRandomDigits(req.Username,5,pr.Services.CustomerService.Config.OtpTtlMinutes)

	utils.ServeJSONObject(w, "", http.StatusOK, "Verified", nil, meta, true)
}

//// GetCustomer godoc
//// @Summary Get a single object
//// @Description Returns a single object matching the given id
//// @Tags Customers
//// @Produce  json
//// @Param  slug path string true "slug of the target object"
//// @Success 200 {object} response.CustomerSuccessRes
//// @Failure 400 {object} response.CustomerErrorRes "Invalid request body, or missing required fields."
//// @Failure 401 {object} response.CustomerErrorRes "Unauthorized access attempt."
//// @Failure 500 {object} response.CustomerErrorRes "API sever or db unreachable."
//// @Router /api/v1/public/generics/{slug} [get]
//func (pr *publicRouter) GetCustomer(w http.ResponseWriter, r *http.Request) {
//	slug := chi.URLParam(r, "slug")
//
//	req := &model.Customer{Username: slug}
//
//	data, err := pr.Services.CustomerService.GetCustomer(r.Context(), req)
//	if err != nil {
//		utils.HandleObjectError(w, err)
//		return
//	}
//
//	utils.ServeJSONObject(w, "", http.StatusOK, "Fetched generic object successfully", &data, nil, true)
//}
//
//// UpdateCustomer godoc
//// @Summary Update generic object
//// @Description Update an existing generic object
//// @Tags Customers
//// @Accept  json
//// @Produce  json
//// @Param  Body body model.CustomerUpdateReq true "Some fields are mandatory"
//// @Success 200 {object} response.CustomerSuccessRes
//// @Failure 400 {object} response.CustomerErrorRes
//// @Failure 404 {object} response.CustomerErrorRes
//// @Failure 500 {object} response.CustomerErrorRes
//// @Router /api/v1/public/generics/{slug} [patch]
//func (pr *publicRouter) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
//	slug := chi.URLParam(r, "slug")
//
//	req := model.CustomerUpdateReq{}
//	err := json.NewDecoder(r.Body).Decode(&req)
//	if err != nil {
//		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
//		return
//	}
//	req.Username = slug
//
//	err = model.Validate(req)
//	if err != nil {
//		utils.HandleObjectError(w, rest_error.NewValidationError("missing required field(s)", err))
//		return
//	}
//
//	data, err := pr.Services.CustomerService.UpdateCustomer(r.Context(), &req)
//	if err != nil {
//		utils.HandleObjectError(w, err)
//		return
//	}
//
//	utils.ServeJSONObject(w, "", http.StatusOK, "Updated generic object successfully", &data, nil, true)
//}
//
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
//	utils.ServeJSONObject(w, "", http.StatusOK, "Purged generic object successfully", &data, nil, true)
//}
