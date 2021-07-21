package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/infra"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	"github.com/iamrz1/ab-auth/model"
	rLog "github.com/iamrz1/rest-log"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"time"
)

type CustomerRepo struct {
	DB    infra.DB
	Cache *infraCache.Redis
	Table string
	Log   rLog.Logger
}

func NewCustomerRepo(db infra.DB, table string, cache *infraCache.Redis, log rLog.Logger) *CustomerRepo {
	return &CustomerRepo{
		DB:    db,
		Cache: cache,
		Table: table,
		Log:   log,
	}
}

func (pr *CustomerRepo) HoldCustomerRegistrationInCache(otp string, doc *model.CustomerSignupReq) error {
	if (*doc) == (model.CustomerSignupReq{}) {
		return rest_error.NewGenericError(http.StatusBadRequest, "Nothing to create")
	}

	data, err := json.Marshal(doc)
	if err != nil {
		pr.Log.Error("HoldCustomerRegistrationInCache", "", err.Error())
		return err
	}

	scmd := pr.Cache.Client.Set(fmt.Sprintf("%s_%s", doc.Username, otp), data, time.Minute*6)
	err = scmd.Err()
	if err != nil {
		pr.Log.Error("HoldCustomerRegistrationInCache", "", err.Error())
		return err
	}

	return nil
}

func (pr *CustomerRepo) GetCustomerRegistrationFromCache(username, otp string) (*model.CustomerSignupReq, error) {
	res := model.CustomerSignupReq{}
	scmd := pr.Cache.Client.Get(fmt.Sprintf("%s_%s", username, otp))
	err := scmd.Err()
	if err != nil {
		pr.Log.Error("GetCustomerRegistrationFromCache", "", err.Error())
		if err == redis.Nil {
			return nil, fmt.Errorf("%s", "OTP expired")
		}
		return nil, err
	}

	b, err := scmd.Bytes()

	err = json.Unmarshal(b, &res)
	if err != nil {
		pr.Log.Error("GetCustomerRegistrationFromCache", "", err.Error())
		return nil, err
	}

	return &res, nil
}

func (pr *CustomerRepo) CreateCustomer(ctx context.Context, doc *model.Customer) error {
	if (*doc) == (model.Customer{}) {
		return rest_error.NewGenericError(http.StatusBadRequest, "Nothing to create")
	}
	err := pr.DB.Insert(ctx, pr.Table, doc)
	if err != nil {
		pr.Log.Error("CreateCustomer", "", err.Error())
		return err
	}

	return nil
}

func (pr *CustomerRepo) GetCustomer(ctx context.Context, selector interface{}) (*model.Customer, error) {
	res := model.Customer{}
	err := pr.DB.FindOne(ctx, pr.Table, selector, &res)
	if err != nil {
		pr.Log.Error("GetCustomer", "", err.Error())
		return nil, err
	}

	return &res, nil
}

func (pr *CustomerRepo) ListCustomers(ctx context.Context, selector interface{}, listOptions *model.ListOptions) ([]*model.Customer, error) {
	res := make([]*model.Customer, 0)
	if listOptions == nil {
		listOptions = &model.ListOptions{}
	}
	if listOptions.Sort == nil {
		listOptions.Sort = bson.M{"_id": -1}
	}
	err := pr.DB.List(ctx, pr.Table, selector, listOptions.Page, listOptions.Limit, &res, listOptions.Sort)
	if err != nil {
		pr.Log.Error("ListCustomers", "", err.Error())
		return nil, err
	}

	return res, nil
}

func (pr *CustomerRepo) UpdateCustomer(ctx context.Context, filter, doc *model.Customer) (int64, error) {
	if (*doc) == (model.Customer{}) {
		log.Println("nothing to update")
		return 0, rest_error.NewGenericError(http.StatusBadRequest, "Nothing to update")
	}
	matched, err := pr.DB.Update(ctx, pr.Table, filter, doc)
	if err != nil {
		pr.Log.Error("updateCustomerProfile", "", err.Error())
		return 0, err
	}

	return matched, nil
}

func (pr *CustomerRepo) CountCustomer(ctx context.Context, selector interface{}) (int64, error) {
	n, err := pr.DB.FindAndCount(ctx, pr.Table, selector)
	if err != nil {
		pr.Log.Error("CountCustomer", "", err.Error())
		return 0, err
	}

	return n, nil
}

func (pr *CustomerRepo) PurgeOne(ctx context.Context, filter interface{}) (int64, error) {
	purged, err := pr.DB.DeleteOne(ctx, pr.Table, filter)
	if err != nil {
		pr.Log.Error("PurgeOne", "", err.Error())
		return 0, err
	}

	return purged, nil
}
