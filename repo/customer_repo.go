package repo

import (
	"context"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/infra"
	infraCache "github.com/iamrz1/ab-auth/infra/cache"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/model"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type CustomerRepo struct {
	DB    infra.DB
	Cache *infraCache.Redis
	Table string
	Log   logger.StructLogger
}

func NewCustomerRepo(db infra.DB, table string, cache *infraCache.Redis, log logger.StructLogger) *CustomerRepo {
	return &CustomerRepo{
		DB:    db,
		Cache: cache,
		Table: table,
		Log:   log,
	}
}

func (pr *CustomerRepo) CreateCustomer(ctx context.Context, doc *model.Customer) error {
	if (*doc) == (model.Customer{}) {
		return rest_error.NewGenericError(http.StatusBadRequest, "Nothing to create")
	}
	err := pr.DB.Insert(ctx, pr.Table, doc)
	if err != nil {
		pr.Log.Errorf("CreateCustomer", "", err.Error())
		return err
	}

	return nil
}

func (pr *CustomerRepo) GetCustomer(ctx context.Context, selector interface{}) (*model.Customer, error) {
	res := model.Customer{}
	err := pr.DB.FindOne(ctx, pr.Table, selector, &res)
	if err != nil {
		pr.Log.Errorf("GetCustomer", "", err.Error())
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
		pr.Log.Errorf("ListCustomers", "", err.Error())
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
		pr.Log.Errorf("UpdateCustomer", "", err.Error())
		return 0, err
	}

	return matched, nil
}

func (pr *CustomerRepo) CountCustomer(ctx context.Context, selector interface{}) (int64, error) {
	n, err := pr.DB.FindAndCount(ctx, pr.Table, selector)
	if err != nil {
		pr.Log.Errorf("CountCustomer", "", err.Error())
		return 0, err
	}

	return n, nil
}

func (pr *CustomerRepo) PurgeOne(ctx context.Context, filter interface{}) (int64, error) {
	purged, err := pr.DB.DeleteOne(ctx, pr.Table, filter)
	if err != nil {
		pr.Log.Errorf("PurgeOne", "", err.Error())
		return 0, err
	}

	return purged, nil
}
