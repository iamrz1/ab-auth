package service

import (
	"github.com/iamrz1/ab-auth/utils"
	"testing"
)

const (
	name = "Service Model X"
	slug = "service-model-x"
)

//var gs *GenericService
//var ctx context.Context
//var doc = model.Generic{
//	Name:         name,
//	Slug:         slug,
//	StringField:  "svc123",
//	DecimalField: 0.01,
//	IntegerField: 10,
//}
//
//func init() {
//	ctx = context.Background()
//	godotenv.Load("../../.env")
//	err := config.LoadConfig()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	cfg := config.GetConfig()
//
//	gracefulTimeout := time.Second * time.Duration(cfg.GracefulTimeout)
//	db, err := infraMongo.New(ctx, cfg.DSN, cfg.Database, gracefulTimeout)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	gs = SetupServiceConfig(cfg, db, logger.DefaultOutStructLogger).GenericService
//}
//
//func TestGenericService_CreateGeneric(t *testing.T) {
//	req := model.GenericCreateReq{
//		Name:         doc.Name,
//		StringField:  doc.StringField,
//		DecimalField: doc.DecimalField,
//		IntegerField: doc.IntegerField,
//	}
//	t.Run("valid data", func(t *testing.T) {
//
//		res, err := gs.CreateGeneric(ctx, &req)
//		assert.NoError(t, err, "failed to create object")
//		assert.EqualValues(t, doc, *res)
//	})
//
//	t.Run("invalid data", func(t *testing.T) {
//		req.Name = ""
//		_, err := gs.CreateGeneric(ctx, &req)
//		assert.Error(t, err, "expected validation error")
//	})
//}
//
//func TestGenericService_GetGeneric(t *testing.T) {
//	t.Run("valid data", func(t *testing.T) {
//		filter := model.Generic{
//			Slug: slug,
//		}
//		res, err := gs.GetGeneric(ctx, &filter)
//		assert.NoError(t, err, "failed to get object")
//
//		assert.EqualValues(t, doc, *res)
//	})
//
//	t.Run("invalid data", func(t *testing.T) {
//		filter := model.Generic{
//			Slug: slug + "i",
//		}
//		_, err := gs.GetGeneric(ctx, &filter)
//		assert.Error(t, err, "expected not found err")
//	})
//}
//
//func TestGenericService_ListGenerics(t *testing.T) {
//	t.Run("valid data", func(t *testing.T) {
//		filter := model.GenericListReq{Search: "svc"}
//		res, count, err := gs.ListGenerics(ctx, &filter)
//		assert.NoError(t, err, "failed to get objects")
//		assert.EqualValues(t, 1, count)
//		assert.EqualValues(t, doc, *res[0])
//
//	})
//
//	t.Run("invalid data", func(t *testing.T) {
//		filter := model.GenericListReq{Search: "bcd"}
//		_, count, err := gs.ListGenerics(ctx, &filter)
//		assert.NoError(t, err, "failed to get objects")
//		assert.EqualValues(t, count, 0)
//	})
//}
//
//func TestGenericService_UpdateGeneric(t *testing.T) {
//	t.Run("valid data", func(t *testing.T) {
//		update := model.GenericUpdateReq{
//			Slug:         slug,
//			DecimalField: 1.01,
//		}
//
//		res, err := gs.UpdateGeneric(ctx, &update)
//		assert.NoError(t, err, "failed to update object")
//		assert.EqualValues(t, update.DecimalField, res.DecimalField)
//	})
//
//	t.Run("revert to old data", func(t *testing.T) {
//		update := model.GenericUpdateReq{
//			Slug:         slug,
//			DecimalField: 0.01,
//		}
//
//		res, err := gs.UpdateGeneric(ctx, &update)
//		assert.NoError(t, err, "failed to update object")
//		assert.EqualValues(t, update.DecimalField, res.DecimalField)
//	})
//
//	t.Run("invalid data", func(t *testing.T) {
//		update := model.GenericUpdateReq{
//			Slug: slug,
//		}
//
//		_, err := gs.UpdateGeneric(ctx, &update)
//		assert.Error(t, err, "shoul dnot update with empty object")
//	})
//}
//
//func TestGenericService_PurgeOne(t *testing.T) {
//	t.Run("valid data", func(t *testing.T) {
//		filter := model.GenericDeleteReq{
//			Slug: slug,
//		}
//		res, err := gs.PurgeGeneric(ctx, &filter)
//		assert.NoError(t, err, "failed to create object")
//		assert.EqualValues(t, doc, *res)
//	})
//
//	t.Run("invalid data", func(t *testing.T) {
//		filter := model.GenericDeleteReq{
//			Slug: slug + "i",
//		}
//		_, err := gs.PurgeGeneric(ctx, &filter)
//		assert.Error(t, err, "expected not found err")
//	})
//}

func TestValidatePassword(t *testing.T) {
	err := utils.ValidatePassword("Evaly2020!")
	if err != nil {
		t.Fatal(err)
	}

	err = utils.ValidatePassword("evaly2020!")
	if err != nil {
		t.Log(err)
	}

	err = utils.ValidatePassword("Evaly2020!বাংলা")
	if err != nil {
		t.Log(err)
	}

}
