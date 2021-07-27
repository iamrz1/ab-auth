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
	"log"
	"net/http"
	"strings"
	"time"
)

type merchantService struct {
	CommonRepo   *repo.CommonRepo
	MerchantRepo *repo.MerchantRepo
	AddressRepo  *repo.AddressRepo
	Log          rLog.Logger
	Config       *config.AppConfig
}

func NewMerchantService(cfg *config.AppConfig, cm *repo.CommonRepo, cs *repo.MerchantRepo, logger rLog.Logger) *merchantService {
	return &merchantService{
		CommonRepo:   cm,
		MerchantRepo: cs,
		Log:          logger,
		Config:       cfg,
	}
}

func (gs *merchantService) CreateMerchant(ctx context.Context, req *model.MerchantSignupReq) (string, error) {
	err := utils.ValidatePassword(req.Password)
	if err != nil {
		return "", rest_error.NewValidationError("", err)
	}

	if !utils.IsValidPhoneNumber(req.Username) {
		return "", rest_error.NewValidationError("Phone number is not valid", nil)
	}

	_, err = gs.GetMerchant(ctx, &model.Merchant{Username: req.Username})
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

	err = gs.MerchantRepo.HoldMerchantRegistrationInCache(otp, req)
	if err != nil {
		gs.Log.Error("CreateMerchant", "", err.Error())
		return "", err
	}

	return otp, nil
}

func (gs *merchantService) VerifyMerchantSignUp(ctx context.Context, req *model.MerchantSignupVerificationReq) error {
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

	merchantData, err := gs.MerchantRepo.GetMerchantRegistrationFromCache(req.Username, req.OTP)
	if err != nil {
		gs.Log.Error("VerifyMerchantSignUp", "", err.Error())
		return rest_error.NewValidationError("", err)
	}

	c := &model.Merchant{
		Username:    merchantData.Username,
		FullName:    merchantData.FullName,
		Password:    utils.GetEncodedPassword(merchantData.Password),
		Status:      utils.StatusPending,
		IsVerified:  utils.BoolP(true),
		IsDeleted:   utils.BoolP(false),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
		LastResetAt: time.Now().UTC(),
	}

	err = gs.MerchantRepo.CreateMerchant(ctx, c)
	if err != nil {
		gs.Log.Error("VerifyMerchantSignUp", "", err.Error())
		return err
	}

	return nil
}

func (gs *merchantService) Login(ctx context.Context, req *model.LoginReq) (*model.Token, error) {
	incorrectMsg := "Incorrect username or password"

	ok, err := gs.CommonRepo.LockKey(fmt.Sprintf("%s_%s_password_match", req.Username, "login"), 5)
	if err != nil || !ok {
		return nil, fmt.Errorf("%s", "Please try again in a few seconds")
	}

	if !gs.CommonRepo.EnsureUsageLimit(fmt.Sprintf("%s_%s_password_match_limit", req.Username, "login"), 5, 5*60) {
		// max 5 try in 5 minutes
		return nil, rest_error.NewGenericError(http.StatusTooManyRequests, "Please try again later")
	}

	g, err := gs.MerchantRepo.GetMerchant(ctx, model.Merchant{Username: req.Username})
	if err != nil {
		gs.Log.Error("login", "", err.Error())
		return nil, rest_error.NewValidationError(incorrectMsg, nil)
	}

	if !utils.VerifyPassword(req.Password, g.Password) {
		gs.Log.Error("login", "", "password mismatch")
		return nil, rest_error.NewValidationError(incorrectMsg, nil)
	}

	utils.SetLastResetAt(g.Username, g.LastResetAt.Unix())

	access, refresh := utils.GenerateTokens(g.Username, "", "merchant")

	return &model.Token{AccessToken: access, RefreshToken: refresh}, nil
}

func (gs *merchantService) GetShortProfile(ctx context.Context, req *model.Token) (*model.MerchantShort, error) {
	claims, err := utils.VerifyToken(req.AccessToken, false)
	if err != nil {
		gs.Log.Error("GetShortProfile", "", err.Error())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, err.Error())
	}

	g, err := gs.MerchantRepo.GetMerchant(ctx, model.Merchant{Username: claims.Username})
	if err != nil {
		gs.Log.Error("GetShortProfile", "", err.Error())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, "Invalid user")
	}

	if claims.IssuedAt < g.LastResetAt.Unix() {
		log.Println("claims.IssuedAt:", claims.IssuedAt, "g.BirthDate.Unix():", g.LastResetAt.Unix())
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, "Inert token")
	}

	if claims.UserType != "merchant" {
		return nil, rest_error.NewGenericError(http.StatusUnauthorized, "Restricted to merchants")
	}

	return g.ToShortResponse(), nil
}

func (gs *merchantService) GetMerchant(ctx context.Context, req *model.Merchant) (*model.Merchant, error) {
	g, err := gs.MerchantRepo.GetMerchant(ctx, req)
	if err != nil {
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *merchantService) ListMerchants(ctx context.Context, req *model.MerchantListReq) ([]*model.Merchant, int64, error) {
	selector := &bson.D{}

	if req.Search != "" {
		selector = utils.AppendSearchPattern(selector, "string_field", req.Search, true)
	}

	opts := &model.ListOptions{
		Page:  req.Page,
		Limit: req.Limit,
		Sort:  nil,
	}

	Merchants, err := gs.MerchantRepo.ListMerchants(ctx, selector, opts)
	if err != nil {
		gs.Log.Error("ListMerchants", "", err.Error())
		return nil, 0, err
	}

	for _, g := range Merchants {
		g = g.ToResponse()
	}

	count, err := gs.MerchantRepo.CountMerchant(ctx, selector)
	if err != nil {
		gs.Log.Error("CountMerchant", "", err.Error())
		return nil, 0, err
	}

	return Merchants, count, nil
}

func (gs *merchantService) UpdateMerchant(ctx context.Context, req *model.MerchantProfileUpdateReq) (*model.Merchant, error) {
	filter := &model.Merchant{Username: req.Username}
	_, err := gs.MerchantRepo.GetMerchant(ctx, filter)
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

	updateDoc := &model.Merchant{
		FullName:      strings.TrimSpace(req.FullName),
		Gender:        gender,
		Email:         strings.ToLower(strings.TrimSpace(req.Email)),
		Occupation:    strings.TrimSpace(req.Occupation),
		Organization:  strings.TrimSpace(req.Organization),
		BirthDate:     utils.GetTimeFromISOString(req.BirthDate),
		ProfilePicURL: strings.TrimSpace(req.ProfilePicURL),
		UpdatedAt:     time.Now().UTC(),
	}

	_, err = gs.MerchantRepo.UpdateMerchant(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateMerchantProfile", "", err.Error())
		return nil, err
	}

	g, err := gs.MerchantRepo.GetMerchant(ctx, filter)
	if err != nil {
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *merchantService) UpdatePassword(ctx context.Context, req *model.UpdatePasswordReq) (*model.Merchant, error) {
	filter := &model.Merchant{Username: req.Username}
	c, err := gs.MerchantRepo.GetMerchant(ctx, filter)
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

	updateDoc := &model.Merchant{
		Password:    utils.GetEncodedPassword(req.NewPassword),
		LastResetAt: time.Now().UTC(),
	}

	_, err = gs.MerchantRepo.UpdateMerchant(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateMerchantProfile", "", err.Error())
		return nil, err
	}

	utils.SetLastResetAt(req.Username, updateDoc.LastResetAt.Unix())

	return c.ToResponse(), nil
}

func (gs *merchantService) DeleteMerchant(ctx context.Context, delete *model.MerchantDeleteReq) (*model.Merchant, error) {
	filter := &model.Merchant{Username: delete.Username}
	_, err := gs.MerchantRepo.GetMerchant(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	updateDoc := &model.Merchant{
		IsDeleted: utils.BoolP(true),
	}

	_, err = gs.MerchantRepo.UpdateMerchant(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateMerchantProfile", "", err.Error())
		return nil, err
	}

	g, err := gs.MerchantRepo.GetMerchant(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *merchantService) PurgeMerchant(ctx context.Context, delete *model.MerchantDeleteReq) (*model.Merchant, error) {
	filter := model.Merchant{Username: delete.Username}
	g, err := gs.MerchantRepo.GetMerchant(ctx, filter)
	if err != nil {
		if err == infra.ErrNotFound {
			return nil, rest_error.NewValidationError("", infra.ErrNotFound)
		}
		return nil, err
	}

	_, err = gs.MerchantRepo.PurgeOne(ctx, filter)
	if err != nil {
		gs.Log.Error("PurgeOne", "", err.Error())
		return nil, err
	}

	return g.ToResponse(), nil
}

func (gs *merchantService) ForgotPassword(ctx context.Context, req *model.ForgotPasswordReq) (string, error) {
	if !utils.IsValidPhoneNumber(req.Username) {
		return "", rest_error.NewValidationError("Phone number is not valid", nil)
	}

	_, err := gs.GetMerchant(ctx, &model.Merchant{Username: req.Username})
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

func (gs *merchantService) SetPassword(ctx context.Context, req *model.SetPasswordReq) error {
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

	filter := &model.Merchant{Username: req.Username}

	updateDoc := &model.Merchant{
		Password:    utils.GetEncodedPassword(req.Password),
		LastResetAt: time.Now().UTC(),
	}

	_, err = gs.MerchantRepo.UpdateMerchant(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateMerchantProfile", "", err.Error())
		return err
	}

	utils.SetLastResetAt(req.Username, updateDoc.LastResetAt.Unix())

	return nil
}

func (gs *merchantService) ChangePassword(ctx context.Context, req *model.SetPasswordReq) error {
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

	filter := &model.Merchant{Username: req.Username}

	updateDoc := &model.Merchant{
		Password:    utils.GetEncodedPassword(req.Password),
		LastResetAt: time.Now().UTC(),
	}

	_, err = gs.MerchantRepo.UpdateMerchant(ctx, filter, updateDoc)
	if err != nil {
		gs.Log.Error("updateMerchantProfile", "", err.Error())
		return err
	}

	utils.SetLastResetAt(req.Username, updateDoc.LastResetAt.Unix())

	return nil
}
