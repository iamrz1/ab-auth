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

type MerchantRepo struct {
	DB    infra.DB
	Cache *infraCache.Redis
	Table string
	Log   rLog.Logger
}

func NewMerchantRepo(db infra.DB, table string, cache *infraCache.Redis, log rLog.Logger) *MerchantRepo {
	return &MerchantRepo{
		DB:    db,
		Cache: cache,
		Table: table,
		Log:   log,
	}
}

func (pr *MerchantRepo) HoldMerchantRegistrationInCache(otp string, doc *model.MerchantSignupReq) error {
	if (*doc) == (model.MerchantSignupReq{}) {
		return rest_error.NewGenericError(http.StatusBadRequest, "Nothing to create")
	}

	data, err := json.Marshal(doc)
	if err != nil {
		pr.Log.Error("HoldMerchantRegistrationInCache", "", err.Error())
		return err
	}

	scmd := pr.Cache.Client.Set(fmt.Sprintf("%s_%s", doc.Username, otp), data, time.Minute*6)
	err = scmd.Err()
	if err != nil {
		pr.Log.Error("HoldMerchantRegistrationInCache", "", err.Error())
		return err
	}

	return nil
}

func (pr *MerchantRepo) GetMerchantRegistrationFromCache(username, otp string) (*model.MerchantSignupReq, error) {
	res := model.MerchantSignupReq{}
	scmd := pr.Cache.Client.Get(fmt.Sprintf("%s_%s", username, otp))
	err := scmd.Err()
	if err != nil {
		pr.Log.Error("GetMerchantRegistrationFromCache", "", err.Error())
		if err == redis.Nil {
			return nil, fmt.Errorf("%s", "OTP expired")
		}
		return nil, err
	}

	b, err := scmd.Bytes()

	err = json.Unmarshal(b, &res)
	if err != nil {
		pr.Log.Error("GetMerchantRegistrationFromCache", "", err.Error())
		return nil, err
	}

	return &res, nil
}

func (pr *MerchantRepo) CreateMerchant(ctx context.Context, doc *model.Merchant) error {
	if (*doc) == (model.Merchant{}) {
		return rest_error.NewGenericError(http.StatusBadRequest, "Nothing to create")
	}
	err := pr.DB.Insert(ctx, pr.Table, doc)
	if err != nil {
		pr.Log.Error("CreateMerchant", "", err.Error())
		return err
	}

	return nil
}

func (pr *MerchantRepo) GetMerchant(ctx context.Context, selector interface{}) (*model.Merchant, error) {
	res := model.Merchant{}
	err := pr.DB.FindOne(ctx, pr.Table, selector, &res)
	if err != nil {
		pr.Log.Error("GetMerchant", "", err.Error())
		return nil, err
	}

	return &res, nil
}

func (pr *MerchantRepo) ListMerchants(ctx context.Context, selector interface{}, listOptions *model.ListOptions) ([]*model.Merchant, error) {
	res := make([]*model.Merchant, 0)
	if listOptions == nil {
		listOptions = &model.ListOptions{}
	}
	if listOptions.Sort == nil {
		listOptions.Sort = bson.M{"_id": -1}
	}
	err := pr.DB.List(ctx, pr.Table, selector, listOptions.Page, listOptions.Limit, &res, listOptions.Sort)
	if err != nil {
		pr.Log.Error("ListMerchants", "", err.Error())
		return nil, err
	}

	return res, nil
}

func (pr *MerchantRepo) UpdateMerchant(ctx context.Context, filter, doc *model.Merchant) (int64, error) {
	if (*doc) == (model.Merchant{}) {
		log.Println("nothing to update")
		return 0, rest_error.NewGenericError(http.StatusBadRequest, "Nothing to update")
	}
	matched, err := pr.DB.Update(ctx, pr.Table, filter, doc)
	if err != nil {
		pr.Log.Error("updateMerchantProfile", "", err.Error())
		return 0, err
	}

	return matched, nil
}

func (pr *MerchantRepo) CountMerchant(ctx context.Context, selector interface{}) (int64, error) {
	n, err := pr.DB.FindAndCount(ctx, pr.Table, selector)
	if err != nil {
		pr.Log.Error("CountMerchant", "", err.Error())
		return 0, err
	}

	return n, nil
}

func (pr *MerchantRepo) PurgeOne(ctx context.Context, filter interface{}) (int64, error) {
	purged, err := pr.DB.DeleteOne(ctx, pr.Table, filter)
	if err != nil {
		pr.Log.Error("PurgeOne", "", err.Error())
		return 0, err
	}

	return purged, nil
}
