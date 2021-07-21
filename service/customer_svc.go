package service

import (
	"context"
	"fmt"
	"github.com/iamrz1/ab-auth/config"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/infra"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/repo"
	"github.com/iamrz1/ab-auth/utils"
	rLog "github.com/iamrz1/rest-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
	"time"
)

type customerService struct {
	CommonRepo   *repo.CommonRepo
	CustomerRepo *repo.CustomerRepo
	AddressRepo  *repo.AddressRepo
	Log          rLog.Logger
	Config       *config.AppConfig
}

func NewCustomerService(cfg *config.AppConfig, cm *repo.CommonRepo, cs *repo.CustomerRepo, ar *repo.AddressRepo, logger rLog.Logger) *customerService {
	return &customerService{
		CommonRepo:   cm,
		CustomerRepo: cs,
		AddressRepo:  ar,
		Log:          logger,
		Config:       cfg,
	}
}

func (gs *customerService) CreateCustomer(ctx context.Context, req *model.CustomerSignupReq) (string, error) {
	err := utils.ValidatePassword(req.Password)
	if err != nil {
		return "", rest_error.NewValidationError("", err)
	}

	if !utils.IsValidPhoneNumber(req.Username) {
		return "", rest_error.NewValidationError("Phone number is not valid", nil)
	}

	_, err = gs.GetCustomer(ctx, &model.Customer{Username: req.Username})
	if err != nil {
		if err != infra.ErrNotFound {
			return "", err
		}
	} else {
		return "", rest_error.NewValidationError("User already exists", err)
	}

	otp, err := gs.CommonRepo.GetOTP(req.Username, "signup", 5, 24*60*60, 10)
	if err != nil {
		return "", rest_error.NewGenericError(http.StatusTooManyRequests, err.Error())
	}
	// todo: send otp

	err = gs.CustomerRepo.HoldCustomerRegistrationInCache(otp, req)
	if err != nil {
		gs.Log.Error("CreateCustomer", "", err.Error())
		return "", err
	}

	return otp, nil
}

func (gs *customerService) VerifyCustomerSignUp(ctx context.Context, req *model.CustomerSignupVerificationReq) error {
	if !utils.IsValidPhoneNumber(req.Username) {
		return rest_error.NewValidationError("Phone number is not valid", nil)
	}

	ok, err := gs.CommonRepo.LockKey(fmt.Sprintf("%s_%s_otp_match", req.Username, "signup"), 5)
	if err != nil || !ok {
		return fmt.Errorf("%s", "Please try again in a few seconds")
	}

	if !gs.CommonRepo.EnsureUsageLimit(fmt.Sprintf("%s_%s_otp_gen_limit", req.Username, "signup"), 5, 5*60) {
		return rest_error.NewGenericError(http.StatusTooManyRequests, "OTP verification failed")
	}

	customerData, err := gs.CustomerRepo.GetCustomerRegistrationFromCache(req.Username, req.OTP)
	if err != nil {
		gs.Log.Error("VerifyCustomerSignUp", "", err.Error())
		return rest_error.NewValidationError("", err)
	}

	c := &model.Customer{
		Username:    customerData.Username,
		FullName:    customerData.FullName,
		Password:    utils.GetEncodedPassword(customerData.Password),
		Status:      utils.StatusActive,
		IsVerified:  utils.BoolP(true),
		IsDeleted:   utils.BoolP(false),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		LastResetAt: time.Now().UTC(),
	}

	err = gs.CustomerRepo.CreateCustomer(ctx, c)
	if err != nil {
		gs.Log.Error("VerifyCustomerSignUp", "", err.Error())
		return err
	}

	return nil
}

func (gs *customerService) Login(ctx context.Context, req *model.LoginReq) (*model.Token, error) {
	incorrectMsg := "Incorrect username or password"

	ok, err := gs.CommonRepo.LockKey(fmt.Sprintf("%s_%s_password_match", req.Username, "login"), 5)
	if err != nil || !ok {
		return nil, fmt.Errorf("%s", "Please try again in a few seconds")
	}

	if !gs.CommonRepo.EnsureUsageLimit(fmt.Sprintf("%s_%s_password_match_limit", req.Username, "login"), 5, 5*60) {
		// max 5 try in 5 minutes
		return nil, rest_error.NewGenericError(http.StatusTooManyRequests, "Please try again later")
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, model.Customer{Username: req.Username})
	if err != nil {
		gs.Log.Error("login", "", err.Error())
		return nil, rest_error.NewValidationError(incorrectMsg, nil)
	}

	if !utils.VerifyPassword(req.Password, g.Password) {
		gs.Log.Error("login", "", "password mismatch")
		return nil, rest_error.NewValidationError(incorrectMsg, nil)
	}

	utils.SetLastResetAt(g.Username, g.LastResetAt.Unix())

	access, refresh := utils.GenerateTokens(g.Username, "", "customer")

	return &model.Token{AccessToken: access, RefreshToken: refresh}, nil
}

func (gs *customerService) GetShortProfile(ctx context.Context, req *model.Token) (*model.CustomerShort, error) {
	claims, err := utils.VerifyToken(req.AccessToken, false)
	if err != nil {
		gs.Log.Error("GetShortProfile", "", err.Error())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, err.Error())
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, model.Customer{Username: claims.Username})
	if err != nil {
		gs.Log.Error("GetShortProfile", "", err.Error())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, "Invalid user")
	}

	if claims.IssuedAt < g.LastResetAt.Unix() {
		log.Println("claims.IssuedAt:", claims.IssuedAt, "g.BirthDate.Unix():", g.LastResetAt.Unix())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, "Inert token")
	}

	if claims.UserType != "customer" {
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, "Restricted to customers")
	}

	return g.ToShortResponse(), nil
}

func (gs *customerService) GetCustomer(ctx context.Context, req *model.Customer) (*model.Customer, error) {
	g, err := gs.CustomerRepo.GetCustomer(ctx, req)
	if err != nil {
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *customerService) ListCustomers(ctx context.Context, req *model.CustomerListReq) ([]*model.Customer, int64, error) {
	selector := &bson.D{}

	if req.Search != "" {
		selector = utils.AppendSearchPattern(selector, "string_field", req.Search, true)
	}

	opts := &model.ListOptions{
		Page:  req.Page,
		Limit: req.Limit,
		Sort:  nil,
	}

	Customers, err := gs.CustomerRepo.ListCustomers(ctx, selector, opts)
	if err != nil {
		gs.Log.Error("ListCustomers", "", err.Error())
		return nil, 0, err
	}

	for _, g := range Customers {
		g = g.ToResponse()
	}

	count, err := gs.CustomerRepo.CountCustomer(ctx, selector)
	if err != nil {
		gs.Log.Error("CountCustomer", "", err.Error())
		return nil, 0, err
	}

	return Customers, count, nil
}

func (gs *customerService) UpdateCustomer(ctx context.Context, req *model.CustomerProfileUpdateReq) (*model.Customer, error) {
	filter := &model.Customer{Username: req.Username}
	_, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}
	gender := strings.ToLower(strings.TrimSpace(req.Gender))
	if gender != "" && !utils.IsGenderValid(gender) {
		return nil, rest_error.NewValidationError("Invalid gender", nil)
	}

	updateDoc := &model.Customer{
		FullName:      strings.TrimSpace(req.FullName),
		Gender:        gender,
		Email:         strings.ToLower(strings.TrimSpace(req.Email)),
		Occupation:    strings.TrimSpace(req.Occupation),
		Organization:  strings.TrimSpace(req.Organization),
		BirthDate:     utils.GetTimeFromISOString(req.BirthDate),
		ProfilePicURL: strings.TrimSpace(req.ProfilePicURL),
		UpdatedAt:     time.Now().UTC(),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateCustomerProfile", "", err.Error())
		return nil, err
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *customerService) UpdatePassword(ctx context.Context, req *model.UpdatePasswordReq) (*model.Customer, error) {
	filter := &model.Customer{Username: req.Username}
	c, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	ok, err := gs.CommonRepo.LockKey(fmt.Sprintf("%s_%s_password_match", req.Username, "update"), 5)
	if err != nil || !ok {
		return nil, fmt.Errorf("%s", "Please try again in a few seconds")
	}

	if !gs.CommonRepo.EnsureUsageLimit(fmt.Sprintf("%s_%s_password_match_limit", req.Username, "update"), 5, 5*60) {
		// max 5 try in 5 minutes
		return nil, rest_error.NewGenericError(http.StatusTooManyRequests, "Please try again later")
	}

	if !utils.VerifyPassword(req.CurrentPassword, c.Password) {
		gs.Log.Error("updatePassword", "", "password mismatch")
		return nil, rest_error.NewValidationError("Incorrect password", nil)
	}

	updateDoc := &model.Customer{
		Password:    utils.GetEncodedPassword(req.NewPassword),
		LastResetAt: time.Now().UTC(),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateCustomerProfile", "", err.Error())
		return nil, err
	}

	utils.SetLastResetAt(req.Username, updateDoc.LastResetAt.Unix())

	return c.ToResponse(), nil
}

func (gs *customerService) DeleteCustomer(ctx context.Context, delete *model.CustomerDeleteReq) (*model.Customer, error) {
	filter := &model.Customer{Username: delete.Username}
	_, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	updateDoc := &model.Customer{
		IsDeleted: utils.BoolP(true),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateCustomerProfile", "", err.Error())
		return nil, err
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *customerService) PurgeCustomer(ctx context.Context, delete *model.CustomerDeleteReq) (*model.Customer, error) {
	filter := model.Customer{Username: delete.Username}
	g, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	_, err = gs.CustomerRepo.PurgeOne(ctx, filter)
	if err != nil {
		gs.Log.Error("PurgeOne", "", err.Error())
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *customerService) ForgotPassword(ctx context.Context, req *model.ForgotPasswordReq) (string, error) {
	if !utils.IsValidPhoneNumber(req.Username) {
		return "", rest_error.NewValidationError("Phone number is not valid", nil)
	}

	_, err := gs.GetCustomer(ctx, &model.Customer{Username: req.Username})
	if err != nil {
		if err != infra.ErrNotFound {
			return "", err
		}
		return "", nil // lets just pretend that the user exists and throw off random api calls
	}

	otp, err := gs.CommonRepo.GetOTP(req.Username, "forgot", 2, 12*60*60, 10)
	if err != nil {
		return "", rest_error.NewGenericError(http.StatusTooManyRequests, err.Error())
	}

	err = gs.CommonRepo.SetOTP(req.Username, "forgot", otp, 5*60)
	if err != nil {
		return "", rest_error.NewValidationError("", err)
	}
	// todo: send otp

	return otp, nil
}

func (gs *customerService) SetPassword(ctx context.Context, req *model.SetPasswordReq) error {
	if !utils.IsValidPhoneNumber(req.Username) {
		return rest_error.NewValidationError("Phone number is not valid", nil)
	}

	ok, err := gs.CommonRepo.LockKey(fmt.Sprintf("%s_%s_otp_match", req.Username, "forgot"), 5)
	if err != nil || !ok {
		return fmt.Errorf("%s", "Please try again in a few seconds")
	}

	if !gs.CommonRepo.EnsureUsageLimit(fmt.Sprintf("%s_%s_otp_match_limit", req.Username, "forgot"), 5, 5*60) {
		// max 5 try in 5 minutes
		return rest_error.NewGenericError(http.StatusTooManyRequests, "OTP verification failed")
	}

	err = gs.CommonRepo.MatchOTP(req.Username, "forgot", req.OTP)
	if err != nil {
		return rest_error.NewValidationError("", err)
	}

	filter := &model.Customer{Username: req.Username}

	updateDoc := &model.Customer{
		Password:    utils.GetEncodedPassword(req.Password),
		LastResetAt: time.Now().UTC(),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateCustomerProfile", "", err.Error())
		return err
	}

	utils.SetLastResetAt(req.Username, updateDoc.LastResetAt.Unix())

	return nil
}

func (gs *customerService) ChangePassword(ctx context.Context, req *model.SetPasswordReq) error {
	if !utils.IsValidPhoneNumber(req.Username) {
		return rest_error.NewValidationError("Phone number is not valid", nil)
	}

	ok, err := gs.CommonRepo.LockKey(fmt.Sprintf("%s_%s_otp_match", req.Username, "forgot"), 5)
	if err != nil || !ok {
		return fmt.Errorf("%s", "Please try again in a few seconds")
	}

	if !gs.CommonRepo.EnsureUsageLimit(fmt.Sprintf("%s_%s_otp_match_limit", req.Username, "forgot"), 5, 5*60) {
		// max 5 try in 5 minutes
		return rest_error.NewGenericError(http.StatusTooManyRequests, "Please try again later")
	}

	err = gs.CommonRepo.MatchOTP(req.Username, "forgot", req.OTP)
	if err != nil {
		return rest_error.NewValidationError("", err)
	}

	filter := &model.Customer{Username: req.Username}

	updateDoc := &model.Customer{
		Password:    utils.GetEncodedPassword(req.Password),
		LastResetAt: time.Now().UTC(),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateCustomerProfile", "", err.Error())
		return err
	}

	utils.SetLastResetAt(req.Username, updateDoc.LastResetAt.Unix())

	return nil
}

//address
func (gs *customerService) AddAddress(ctx context.Context, req *model.Address) ([]*model.Address, error) {
	filter := model.Address{Username: req.Username, IsDeleted: utils.BoolP(false)}
	n, err := gs.AddressRepo.GetAddressCount(ctx, filter)
	if err != nil {
		return nil, err
	}
	if n >= utils.MaxAddressAllowed {
		return nil, rest_error.NewValidationError(fmt.Sprintf("Maximum %d addresses are allowed", utils.MaxAddressAllowed), nil)
	}

	if n == 0 {
		req.IsPrimary = utils.BoolP(true)
	} else {
		req.IsPrimary = utils.BoolP(false)
	}

	err = gs.AddressRepo.AddAddress(ctx, req)
	if err != nil {
		return nil, err
	}

	list, err := gs.AddressRepo.GetAddresses(ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (gs *customerService) UpdateAddress(ctx context.Context, req *model.Address) ([]*model.Address, error) {
	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, rest_error.NewValidationError("Invalid address ID", nil)
	}
	req.ID = ""
	filter := bson.E{Key: "_id", Value: objID}
	n, err := gs.AddressRepo.UpdateAddress(ctx, filter, req)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, rest_error.NewValidationError("Nothing to update", nil)
	}

	getFilter := model.Address{Username: req.Username, IsDeleted: utils.BoolP(false)}
	list, err := gs.AddressRepo.GetAddresses(ctx, getFilter, nil)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (gs *customerService) RemoveAddress(ctx context.Context, req *model.Address) ([]*model.Address, error) {
	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, rest_error.NewValidationError("Invalid address ID", nil)
	}
	filter := bson.E{Key: "_id", Value: objID}

	doc := &model.Address{IsDeleted: utils.BoolP(true)}
	n, err := gs.AddressRepo.UpdateAddress(ctx, filter, doc)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, rest_error.NewValidationError("Address not found", nil)
	}

	getFilter := model.Address{Username: req.Username, IsDeleted: utils.BoolP(false)}
	list, err := gs.AddressRepo.GetAddresses(ctx, getFilter, nil)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (gs *customerService) GetAddresses(ctx context.Context, username string) ([]*model.Address, error) {
	getFilter := model.Address{Username: username, IsDeleted: utils.BoolP(false)}
	list, err := gs.AddressRepo.GetAddresses(ctx, getFilter, nil)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (gs *customerService) GetPrimaryAddress(ctx context.Context, username string) (*model.Address, error) {
	getFilter := model.Address{Username: username, IsDeleted: utils.BoolP(false), IsPrimary: utils.BoolP(true)}
	address, err := gs.AddressRepo.GetAddress(ctx, getFilter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewGenericError(http.StatusExpectationFailed, "No default address")
		}
		return nil, err
	}

	return address, nil
}

func (gs *customerService) SetPrimaryAddress(ctx context.Context, req *model.Address) ([]*model.Address, error) {
	getFilter := model.Address{Username: req.Username, IsDeleted: utils.BoolP(false)}
	list, err := gs.AddressRepo.GetAddresses(ctx, getFilter, nil)
	if err != nil {
		return nil, err
	}

	found := false
	isAlreadyPrimary := false
	oldPrimaryID := ""
	for _, address := range list {
		if address.ID == req.ID {
			found = true
			if address.IsPrimary != nil && (*(address).IsPrimary) {
				isAlreadyPrimary = true
			}
			address.IsPrimary = utils.BoolP(true)
		} else {
			if address.IsPrimary != nil && (*(address).IsPrimary) {
				oldPrimaryID = address.ID
				address.IsPrimary = utils.BoolP(false)
			}
		}
	}

	if !found {
		return nil, rest_error.NewValidationError("Unknown address ID", nil)
	}

	if isAlreadyPrimary {
		return list, nil
	}

	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return nil, rest_error.NewValidationError("Invalid address ID", nil)
	}
	filter := bson.E{Key: "_id", Value: objID}
	_, err = gs.AddressRepo.UpdateAddress(ctx, filter, &model.Address{IsPrimary: utils.BoolP(true)})
	if err != nil {
		return nil, err
	}

	if oldPrimaryID != "" {
		objID, err = primitive.ObjectIDFromHex(oldPrimaryID)
		filter = bson.E{Key: "_id", Value: objID}
		_, err = gs.AddressRepo.UpdateAddress(ctx, filter, &model.Address{IsPrimary: utils.BoolP(false)})
		if err != nil {
			return nil, err
		}
	}

	return list, nil
}
