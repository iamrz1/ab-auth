package repo

import (
	"context"
	"fmt"
	rest_error "github.com/iamrz1/ab-auth/error"
	"github.com/iamrz1/ab-auth/infra"
	"github.com/iamrz1/ab-auth/model"
	rLog "github.com/iamrz1/rest-log"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
)

type AddressRepo struct {
	DB           infra.DB
	AddressTable string
	BDGeoTable   string
	Log          rLog.Logger
}

func NewAddressRepo(db infra.DB, addressTable, bdLocationTable string, log rLog.Logger) *AddressRepo {
	return &AddressRepo{
		DB:           db,
		Log:          log,
		AddressTable: addressTable,
		BDGeoTable:   bdLocationTable,
	}
}

func (ar *AddressRepo) AddAddress(ctx context.Context, address *model.Address) error {
	err := ar.DB.Insert(ctx, ar.AddressTable, address)
	if err != nil {
		ar.Log.Error("GetAddresses", "", fmt.Sprintf("insert err: %s", err.Error()))
		return err
	}

	return nil
}

func (ar *AddressRepo) GetAddress(ctx context.Context, filter interface{}) (*model.Address, error) {
	res := &model.Address{}
	err := ar.DB.FindOne(ctx, ar.AddressTable, filter, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ar *AddressRepo) GetAddressCount(ctx context.Context, filter interface{}) (int64, error) {
	n, err := ar.DB.FindAndCount(ctx, ar.AddressTable, filter)
	if err != nil {
		ar.Log.Error("CountCustomer", "", err.Error())
		return 0, err
	}

	return n, nil
}

func (ar *AddressRepo) GetAddresses(ctx context.Context, filter interface{}, listOptions *model.ListOptions) ([]*model.Address, error) {
	res := make([]*model.Address, 0)
	if listOptions == nil {
		listOptions = &model.ListOptions{}
	}
	if listOptions.Sort == nil {
		listOptions.Sort = bson.M{"_id": -1}
	}
	err := ar.DB.List(ctx, ar.AddressTable, filter, listOptions.Page, listOptions.Limit, &res, listOptions.Sort)
	if err != nil {
		ar.Log.Error("GetAddresses", "", err.Error())
		return nil, err
	}

	return res, nil
}

func (ar *AddressRepo) UpdateAddress(ctx context.Context, filter interface{}, doc *model.Address) (int64, error) {
	if (*doc) == (model.Address{}) {
		log.Println("nothing to update")
		return 0, rest_error.NewGenericError(http.StatusBadRequest, "Nothing to update")
	}
	matched, err := ar.DB.Update(ctx, ar.AddressTable, filter, doc)
	if err != nil {
		ar.Log.Error("updateCustomerProfile", "", err.Error())
		return 0, err
	}

	return matched, nil
}

func (ar *AddressRepo) PurgeAddress(ctx context.Context, filter interface{}) (int64, error) {
	purged, err := ar.DB.DeleteOne(ctx, ar.AddressTable, filter)
	if err != nil {
		ar.Log.Error("PurgeOne", "", err.Error())
		return 0, err
	}

	return purged, nil
}

func (ar *AddressRepo) GetBdLocations(ctx context.Context, filter interface{}, listOptions *model.ListOptions) ([]*model.BDLocation, error) {
	res := make([]*model.BDLocation, 0)
	if listOptions == nil {
		listOptions = &model.ListOptions{}
	}
	if listOptions.Sort == nil {
		listOptions.Sort = bson.M{"_id": -1}
	}
	err := ar.DB.List(ctx, ar.BDGeoTable, filter, listOptions.Page, listOptions.Limit, &res, listOptions.Sort)
	if err != nil {
		ar.Log.Error("GetAddresses", "", err.Error())
		return nil, err
	}

	return res, nil
}
