package service

import (
	"context"
	"github.com/iamrz1/ab-auth/config"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/infra"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/model"
	"github.com/iamrz1/ab-auth/repo"
	"github.com/iamrz1/ab-auth/utils"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strings"
	"time"
)

type customerService struct {
	CommonRepo   *repo.CommonRepo
	CustomerRepo *repo.CustomerRepo
	Log          logger.StructLogger
	Config       *config.AppConfig
}

func NewCustomerService(cfg *config.AppConfig, cm *repo.CommonRepo, cs *repo.CustomerRepo, logger logger.StructLogger) *customerService {
	return &customerService{
		CommonRepo:   cm,
		CustomerRepo: cs,
		Log:          logger,
		Config:       cfg,
	}
}

func (gs *customerService) CreateCustomer(ctx context.Context, req *model.CustomerSignupReq) (string, error) {
	err := validatePassword(req.Password)
	if err != nil {
		return "", rest_error.NewValidationError("", err)
	}

	if !utils.IsValidPhoneNumber(req.Username) {
		return "", rest_error.NewValidationError("Phone number is not valid", nil)
	}

	_, err = gs.GetCustomer(ctx, &model.Customer{Username: req.Username})
	if err != nil {
		if err != infra.ErrNotFound {
			return "", rest_error.NewValidationError("", infra.ErrNotFound)
		}
	} else {
		return "", rest_error.NewValidationError("User already exists", err)
	}

	otp := utils.GetRandomDigits(5)
	// todo: send otp
	//gs.CommonRepo.SaveOTP(req.Username, otp)

	err = gs.CustomerRepo.HoldCustomerRegistrationInCache(otp, req)
	if err != nil {
		gs.Log.Errorf("CreateCustomer", "", err.Error())
		return "", err
	}

	return otp, nil
}

func (gs *customerService) VerifyCustomerSignUp(ctx context.Context, req *model.CustomerSignupVerificationReq) error {
	if !utils.IsValidPhoneNumber(req.Username) {
		return rest_error.NewValidationError("Phone number is not valid", nil)
	}

	customerData, err := gs.CustomerRepo.GetCustomerRegistrationFromCache(req.Username, req.OTP)
	if err != nil {
		gs.Log.Errorf("VerifyCustomerSignUp", "", err.Error())
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
		gs.Log.Errorf("VerifyCustomerSignUp", "", err.Error())
		return err
	}

	return nil
}

func (gs *customerService) Login(ctx context.Context, req *model.LoginReq) (*model.Token, error) {
	incorrectMsg := "Incorrect username or password"

	g, err := gs.CustomerRepo.GetCustomer(ctx, model.Customer{Username: req.Username})
	if err != nil {
		gs.Log.Errorf("Login", "", err.Error())
		return nil, rest_error.NewValidationError(incorrectMsg, nil)
	}

	if !utils.VerifyPassword(req.Password, g.Password) {
		gs.Log.Errorf("Login", "", "password mismatch")
		return nil, rest_error.NewValidationError(incorrectMsg, nil)
	}

	access, refresh := utils.GenerateTokens(g.Username, "", "customer")

	return &model.Token{AccessToken: access, RefreshToken: refresh}, nil
}

//func (gs *customerService) Refresh(ctx context.Context, req *model.Token) (*model.Token, error) {
//	//incorrectMsg := "Incorrect username or password"
//
//	//g, err := gs.CustomerRepo.GetCustomer(ctx, model.Customer{Username: req.Username})
//	//if err != nil {
//	//	gs.Log.Errorf("Login", "", err.Error())
//	//	return nil, rest_error.NewValidationError(incorrectMsg, nil)
//	//}
//	//
//	//if !utils.VerifyPassword(req.Password, g.Password) {
//	//	gs.Log.Errorf("Login", "", "password mismatch")
//	//	return nil, rest_error.NewValidationError(incorrectMsg, nil)
//	//}
//
//	access, refresh := utils.GenerateTokens(g.Username, "", "customer")
//
//	return &model.Token{AccessToken: access, RefreshToken: refresh}, nil
//}

func (gs *customerService) GetShortProfile(ctx context.Context, req *model.Token) (*model.CustomerShort, error) {
	claims, err := utils.VerifyToken(req.AccessToken, false)
	if err != nil {
		gs.Log.Errorf("VerifyAccessToken", "", err.Error())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, err.Error())
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, model.Customer{Username: claims.Username})
	if err != nil {
		gs.Log.Errorf("VerifyAccessToken", "", err.Error())
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
		gs.Log.Errorf("ListCustomers", "", err.Error())
		return nil, 0, err
	}

	for _, g := range Customers {
		g = g.ToResponse()
	}

	count, err := gs.CustomerRepo.CountCustomer(ctx, selector)
	if err != nil {
		gs.Log.Errorf("CountCustomer", "", err.Error())
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

	updateDoc := &model.Customer{
		FullName:     strings.TrimSpace(req.FullName),
		Gender:       strings.ToLower(strings.TrimSpace(req.Gender)),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		Occupation:   strings.TrimSpace(req.Occupation),
		Organization: strings.TrimSpace(req.Organization),
		BirthDate:    utils.GetTimeFromISOString(req.BirthDate),
		UpdatedAt:    time.Now().UTC(),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Errorf("UpdateCustomerProfile", "", err.Error())
		return nil, err
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		return nil, err
	}

	return g.ToResponse(), nil
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
		gs.Log.Errorf("UpdateCustomerProfile", "", err.Error())
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
		gs.Log.Errorf("PurgeOne", "", err.Error())
		return nil, err
	}

	return g.ToResponse(), nil
}
