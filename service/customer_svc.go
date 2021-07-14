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

	c := &model.Customer{
		Username:   req.Username,
		FullName:   req.FullName,
		Password:   utils.GetPasswordHash(req.Password),
		Status:     utils.StatusActive,
		IsVerified: utils.BoolP(false),
		IsDeleted:  utils.BoolP(false),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	_, err = gs.GetCustomer(ctx, &model.Customer{Username: c.Username})
	if err != nil {
		if err != infra.ErrNotFound {
			return "", rest_error.NewValidationError("", infra.ErrNotFound)
		}
	} else {
		return "", rest_error.NewValidationError("User already exists", err)
	}

	err = gs.CustomerRepo.CreateCustomer(ctx, c)
	if err != nil {
		gs.Log.Errorf("CreateCustomer", "", err.Error())
		return "", err
	}

	otp := utils.GetRandomDigits(5)
	// todo: send otp
	gs.CommonRepo.SaveOTP(req.Username, otp)

	return otp, nil
}

func (gs *customerService) VerifyCustomerSignUp(ctx context.Context, req *model.CustomerSignupVerificationReq) error {
	if !utils.IsValidPhoneNumber(req.Username) {
		return rest_error.NewValidationError("Phone number is not valid", nil)
	}

	err := gs.CommonRepo.VerifyOTP(req.Username, "signup", req.OTP)
	if err != nil {
		return rest_error.NewValidationError("", err)
	}

	filter := &model.Customer{Username: req.Username}
	update := &model.Customer{
		IsVerified: utils.BoolP(true),
		UpdatedAt:  time.Now().UTC(),
	}

	_, err = gs.CustomerRepo.UpdateCustomer(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (gs *customerService) GetCustomer(ctx context.Context, req *model.Customer) (*model.Customer, error) {
	g, err := gs.CustomerRepo.GetCustomer(ctx, req)
	if err != nil {
		return nil, err
	}

	return g.ToCustomerResponse(), nil
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
		g = g.ToCustomerResponse()
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
		gs.Log.Errorf("UpdateCustomer", "", err.Error())
		return nil, err
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		return nil, err
	}

	return g.ToCustomerResponse(), nil
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
		gs.Log.Errorf("UpdateCustomer", "", err.Error())
		return nil, err
	}

	g, err := gs.CustomerRepo.GetCustomer(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	return g.ToCustomerResponse(), nil
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

	return g.ToCustomerResponse(), nil
}
