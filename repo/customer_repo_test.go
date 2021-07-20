package repo

import (
	"context"
	"github.com/iamrz1/ab-auth/config"
	infraMongo "github.com/iamrz1/ab-auth/infra/mongo"
	"github.com/iamrz1/ab-auth/logger"
	"github.com/iamrz1/ab-auth/model"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

const (
	name = "User X"
	slug = "repo-model-x"
)

var cr CustomerRepo
var ctx context.Context
var doc = model.Customer{
	Username:            "01746410745",
	FullName:            name,
	Password:            "1111111",
	RecoveryPhoneNumber: "01746410745",
	Gender:              "male",
	BirthDateString:     "",
	IsDeleted:           nil,
}

func init() {
	ctx = context.Background()
	godotenv.Load("../../.env")
	err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	cfg := config.GetConfig()

	gracefulTimeout := time.Second * time.Duration(cfg.GracefulTimeout)
	db, err := infraMongo.New(ctx, cfg.DSN, cfg.Database, gracefulTimeout)
	if err != nil {
		log.Println(err)
		return
	}

	cr = CustomerRepo{
		DB:    db,
		Table: cfg.CustomerTable,
		Log:   logger.GetDefaultStructLogger(true),
	}
}

func TestCustomerRepo_CreateCustomer(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		err := cr.CreateCustomer(ctx, &doc)
		assert.NoError(t, err, "failed to create object")
	})

	t.Run("invalid data", func(t *testing.T) {
		err := cr.CreateCustomer(ctx, &model.Customer{})
		assert.Error(t, err, "expected empty object error")
	})
}

func TestCustomerRepo_GetCustomer(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug,
		}
		res, err := cr.GetCustomer(ctx, &filter)
		assert.NoError(t, err, "failed to get object")

		assert.EqualValues(t, doc, *res)
	})

	t.Run("invalid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug + "i",
		}
		_, err := cr.GetCustomer(ctx, &filter)
		assert.Error(t, err, "expected not found err")
	})
}

func TestCustomerRepo_ListCustomers(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		filter := model.Customer{Username: slug}
		res, err := cr.ListCustomers(ctx, &filter, nil)
		assert.NoError(t, err, "failed to get objects")
		for _, r := range res {
			log.Println(*r)
		}
		log.Println(res)
		assert.EqualValues(t, 1, len(res))
		assert.EqualValues(t, doc, *res[0])

	})

	t.Run("invalid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug + "i",
		}
		_, err := cr.GetCustomer(ctx, &filter)
		assert.Error(t, err, "expected not found err")
	})
}

func TestCustomerRepo_UpdateCustomer(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug,
		}

		update := model.Customer{RecoveryPhoneNumber: "101"}
		matched, err := cr.UpdateCustomer(ctx, &filter, &update)
		assert.NoError(t, err, "failed to update object")
		assert.EqualValues(t, 1, matched)

		res, err := cr.GetCustomer(ctx, &filter)
		assert.NoError(t, err, "failed to get object")
		assert.EqualValues(t, update.RecoveryPhoneNumber, res.RecoveryPhoneNumber)
	})

	t.Run("invalid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug + "i",
		}

		update := model.Customer{RecoveryPhoneNumber: "201"}
		matched, err := cr.UpdateCustomer(ctx, &filter, &update)
		assert.NoError(t, err, "failed to update object")
		assert.EqualValues(t, 0, matched)
	})
}

func TestCustomerRepo_CountCustomer(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		filter := model.Customer{Username: slug}
		count, err := cr.CountCustomer(ctx, &filter)
		assert.NoError(t, err, "failed to get object")
		assert.EqualValues(t, 1, count)
	})

	t.Run("invalid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug + "i",
		}
		count, err := cr.CountCustomer(ctx, &filter)
		assert.NoError(t, err, "failed to get object")
		assert.EqualValues(t, 0, count)
	})
}

func TestCustomerRepo_PurgeOne(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug,
		}
		n, err := cr.PurgeOne(ctx, &filter)
		assert.NoError(t, err, "failed to purge object")
		assert.EqualValues(t, 1, n)
		assert.NotEqualValues(t, 0, n)
		time.Sleep(time.Second)
	})

	t.Run("invalid data", func(t *testing.T) {
		filter := model.Customer{
			Username: slug + "i",
		}
		n, err := cr.PurgeOne(ctx, &filter)
		assert.NoError(t, err, "expected not found err")
		assert.EqualValues(t, 0, n)
		assert.NotEqualValues(t, 1, n)
	})
}
