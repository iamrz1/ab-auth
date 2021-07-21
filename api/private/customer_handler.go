package private

import (
	"encoding/json"
	"github.com/go-chi/chi"
	_ "github.com/iamrz1/ab-auth/docs"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/utils"
	"net/http"
)

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
func (pr *privateRouter) verifyAccessToken(w http.ResponseWriter, r *http.Request) {
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
func (pr *privateRouter) getCustomerProfile(w http.ResponseWriter, r *http.Request) {
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
func (pr *privateRouter) updateCustomerProfile(w http.ResponseWriter, r *http.Request) {
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
func (pr *privateRouter) refreshToken(w http.ResponseWriter, r *http.Request) {
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
func (pr *privateRouter) updatePassword(w http.ResponseWriter, r *http.Request) {
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

// addNewAddressHandler godoc
// @Summary Add a customer address
// @Description Add a customer address as long as the total address count for the customer is not greater than 5
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.AddressCreateReq true "Some fields are mandatory"
// @Success 201 {object} response.AddressListSuccessRes
// @Failure 400 {object} response.EmptyListErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyListErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyListErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/address [post]
func (pr *privateRouter) addNewAddressHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	var req = &model.AddressCreateReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}

	req.Username = username

	res, err := pr.Services.CustomerService.AddAddress(r.Context(), req.ToAddress())
	if err != nil {
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONList(w, http.StatusCreated, "Address created", res, nil, true)
}

// updateAddressHandler godoc
// @Summary Update address by id
// @Description Update an address for customer using address id
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Param  Body body model.AddressUpdateReq true "Some fields are mandatory"
// @Success 201 {object} response.AddressListSuccessRes
// @Failure 400 {object} response.EmptyListErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyListErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyListErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/address/{id} [patch]
func (pr *privateRouter) updateAddressHandler(w http.ResponseWriter, r *http.Request) {
	var req = &model.AddressUpdateReq{}
	id := chi.URLParam(r, "id")
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}

	req.ID = id
	req.Username = username

	res, err := pr.Services.CustomerService.UpdateAddress(r.Context(), req.ToAddress())
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONList(w, http.StatusOK, "Address updated", res, nil, true)
}

// removeAddressHandler godoc
// @Summary Remove a customer address
// @Description Remove an address for customer using address id
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 201 {object} response.AddressListSuccessRes
// @Failure 400 {object} response.EmptyListErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyListErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyListErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/address/{id} [delete]
func (pr *privateRouter) removeAddressHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	var req = &model.Address{ID: id, Username: username}

	res, err := pr.Services.CustomerService.RemoveAddress(r.Context(), req)
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONList(w, http.StatusOK, "Address removed", res, nil, true)
}

// getAddressesHandler godoc
// @Summary Get customer's addresses
// @Description Get all the addresses of the requesting customer
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 201 {object} response.AddressListSuccessRes
// @Failure 400 {object} response.EmptyListErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyListErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyListErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/address/all [get]
func (pr *privateRouter) getAddressesHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	res, err := pr.Services.CustomerService.GetAddresses(r.Context(), username)
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONList(w, http.StatusOK, "Address updated", res, nil, true)
}

// getPrimaryAddressHandler godoc
// @Summary Get primary address
// @Description Get the primary address of the requesting customer
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 201 {object} response.AddressSuccessRes
// @Failure 400 {object} response.EmptyErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyErrorRes "Unauthorized access attempt."
// @Failure 417 {object} response.EmptyErrorRes "User is yet to set a primary address"
// @Failure 500 {object} response.EmptyErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/address/primary [get]
func (pr *privateRouter) getPrimaryAddressHandler(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	res, err := pr.Services.CustomerService.GetPrimaryAddress(r.Context(), username)
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONList(w, http.StatusOK, "Address removed", res, nil, true)
}

// setPrimaryAddressHandler godoc
// @Summary Set a primary address
// @Description Set an address using address id as the primary address for customer, remove the previous primary address if needed
// @Tags Customers
// @Produce  json
// @Param authorization header string true "Set access token here"
// @Success 201 {object} response.AddressListSuccessRes
// @Failure 400 {object} response.EmptyListErrorRes "Invalid request body, or missing required fields."
// @Failure 401 {object} response.EmptyListErrorRes "Unauthorized access attempt."
// @Failure 500 {object} response.EmptyListErrorRes "API sever or db unreachable."
// @Router /api/v1/private/customers/address/primary/{id} [post]
func (pr *privateRouter) setPrimaryAddressHandler(w http.ResponseWriter, r *http.Request) {
	var req = &model.AddressUpdateReq{}
	id := chi.URLParam(r, "id")
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleObjectError(w, rest_error.NewValidationError("Invalid JSON", err))
		return
	}

	req.ID = id
	req.Username = username

	res, err := pr.Services.CustomerService.SetPrimaryAddress(r.Context(), req.ToAddress())
	if err != nil {
		pr.Log.Error("updateAddressHandler", "", err.Error())
		utils.HandleListError(w, err)
		return
	}

	utils.ServeJSONList(w, http.StatusOK, "Address updated", res, nil, true)
}

//
//// getPrimaryAddressWithSecretHandler gets the primary address for customer
//func getPrimaryAddressWithSecretHandler(w http.ResponseWriter, r *http.Request) {
//	log.Println(">>>Calling RPC from getPrimaryAddressWithSecretHandler")
//	username := chi.URLParam(r, "username")
//
//	if len(username) < 1 {
//		err := utils.ServeJSON(w, "E_INTERNAL", http.StatusInternalServerError, "", nil, false)
//
//		if err != nil {
//			log.Println(err)
//			return
//		}
//	}
//	res, err := grpc_client.GetAuthClients().AddressClient.GetPrimaryAddress(r.Context(), &pb.Address{Username: username})
//	if err != nil {
//		log.Println("Got error from Calling RPC: GetPrimaryAddress", err)
//		e, ok := status.FromError(err)
//		if !ok {
//			utils.ServeJSON(w, "E_INTERNAL", http.StatusInternalServerError, "", nil, false)
//			return
//		}
//
//		if e.Code() == codes.Internal {
//			utils.ServeJSON(w, utils.Code(e.Code().String()), utils.HTTPStatusFromCode(e.Code()), "Something went wrong", nil, false)
//			return
//		}
//
//		utils.ServeJSON(w, utils.Code(e.Code().String()), utils.HTTPStatusFromCode(e.Code()), e.Message(), nil, false)
//		return
//	}
//
//	if res == nil {
//		utils.ServeJSON(w, "", http.StatusOK, "Fetched address successfully", model.Address{}, true)
//		return
//	}
//
//	log.Println("<< GetPrimaryAddress successful")
//	utils.ServeJSON(w, "", http.StatusOK, "Fetched address successfully", NewAddressPbToResponse(res), true)
//}
//

func (pr *privateRouter) purgeCustomer(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	req := model.CustomerDeleteReq{Username: username}

	data, err := pr.Services.CustomerService.PurgeCustomer(r.Context(), &req)
	if err != nil {
		utils.HandleObjectError(w, err)
		return
	}

	utils.ServeJSONObject(w, http.StatusOK, "Purged customer successfully", &data, nil, true)
}
