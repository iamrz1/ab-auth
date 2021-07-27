package private

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/iamrz1/ab-auth/api/middleware"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/service"
	"github.com/iamrz1/ab-auth/utils"
	rLog "github.com/iamrz1/rest-log"
	"net/http"
)

func newAddressRouter(svc *service.Config, rLogger rLog.Logger) *addressRouter {
	return &addressRouter{
		Services: svc,
		Log:      rLogger,
	}
}

type addressRouter struct {
	Services *service.Config
	Log      rLog.Logger
}

func (pr *privateRouter) addressRouter() *chi.Mux {
	r := chi.NewRouter()

	ar := newAddressRouter(pr.Services, pr.Log)
	// address apis, private
	r.With(middleware.AuthenticatedOnly).Post("/", ar.addNewAddressHandler)
	r.With(middleware.AuthenticatedOnly).Patch("/{id}", ar.updateAddressHandler)
	r.With(middleware.AuthenticatedOnly).Delete("/{id}", ar.removeAddressHandler) //empty body
	r.With(middleware.AuthenticatedOnly).Get("/all", ar.getAddressesHandler)
	r.With(middleware.AuthenticatedOnly).Get("/primary", ar.getPrimaryAddressHandler)
	r.With(middleware.AuthenticatedOnly).Post("/primary/{id}", ar.setPrimaryAddressHandler) //empty body

	//r.With(middleware.SecretOnly).Get("/secret/get-primary-address/{username}", getPrimaryAddressWithSecretHandler)

	return r
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
func (pr *addressRouter) addNewAddressHandler(w http.ResponseWriter, r *http.Request) {
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
func (pr *addressRouter) updateAddressHandler(w http.ResponseWriter, r *http.Request) {
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
func (pr *addressRouter) removeAddressHandler(w http.ResponseWriter, r *http.Request) {
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
func (pr *addressRouter) getAddressesHandler(w http.ResponseWriter, r *http.Request) {
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
func (pr *addressRouter) getPrimaryAddressHandler(w http.ResponseWriter, r *http.Request) {
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
func (pr *addressRouter) setPrimaryAddressHandler(w http.ResponseWriter, r *http.Request) {
	var req = &model.AddressUpdateReq{}
	id := chi.URLParam(r, "id")
	username := r.Header.Get(utils.UsernameKey)
	if username == "" {
		utils.HandleObjectError(w, rest_error.NewGenericError(http.StatusUnauthorized, "Missing username"))
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
