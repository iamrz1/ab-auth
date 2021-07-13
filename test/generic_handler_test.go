package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/iamrz1/auth/api"
	"github.com/iamrz1/auth/config"
	_ "github.com/iamrz1/auth/docs"
	infraMongo "github.com/iamrz1/auth/infra/mongo"
	"github.com/iamrz1/auth/logger"
	"github.com/iamrz1/auth/model"
	"github.com/iamrz1/auth/model/response"
	"github.com/iamrz1/auth/service"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	name = "Product Y"
	slug = "product-y"
)

var db *infraMongo.Mongo

func getServer() (*chi.Mux, error) {
	err := config.LoadConfig()
	if err != nil {
		log.Println("could not load one or more config")
		return nil, err
	}
	//logStruct
	logStruct := logger.DefaultOutStructLogger
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	cfg := config.GetConfig()

	gracefulTimeout := time.Second * time.Duration(cfg.GracefulTimeout)
	db, err = infraMongo.New(ctx, cfg.DSN, cfg.Database, gracefulTimeout)
	if err != nil {
		return nil, err
	}

	log.Println("db initialized")
	//cfg := config.getConfig()

	//addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	svc := service.SetupServiceConfig(cfg, db, logStruct)
	handler, err := api.SetupRouter(cfg, svc, logStruct)
	if err != nil {
		log.Println("cant setup router:", err)
		return nil, err
	}

	return handler, nil
}

func closeResource() {
	db.Close(context.Background())
}

func TestCreateObject(t *testing.T) {
	product := model.Generic{
		Name:         name,
		StringField:  "dfdg",
		DecimalField: 0.877,
		IntegerField: 23,
	}

	b, err := json.Marshal(&product)
	assert.NoError(t, err, "marshal failed")

	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/generics", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	mux, err := getServer()
	assert.NoError(t, err, "failed to get server")

	defer closeResource()

	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	d := map[string]interface{}{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValuesf(t, resp.StatusCode, http.StatusCreated, "failed with response body: %v", d)
}

func TestGetObjects(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/generics", nil)
	rec := httptest.NewRecorder()

	mux, err := getServer()
	assert.NoError(t, err, "failed to get server")
	defer closeResource()

	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	d := map[string]interface{}{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)

	res := response.GenericListSuccessRes{}
	merr := json.Unmarshal(body, &res)
	if merr != nil {
		t.Fatal(merr)
	}

	//assert.EqualValues(t, 1, len(res.Data.Generics))
}

func TestGetSingleObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/generics/"+slug, nil)
	rec := httptest.NewRecorder()

	mux, err := getServer()
	assert.NoError(t, err, "failed to get server")
	defer closeResource()

	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	d := map[string]interface{}{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
}

func TestUpdateObject(t *testing.T) {
	product := model.Generic{
		IntegerField: 100,
	}

	b, err := json.Marshal(&product)
	assert.NoError(t, err, "marshal failed")

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/public/generics/"+slug, bytes.NewReader(b))
	rec := httptest.NewRecorder()

	mux, err := getServer()
	assert.NoError(t, err, "failed to get server")
	defer closeResource()

	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	d := map[string]interface{}{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
}

func TestPurgeSingleObject(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/public/generics/purge/"+slug, nil)
	rec := httptest.NewRecorder()

	mux, err := getServer()
	if err != nil {
		t.Fatal(err)
	}
	defer closeResource()

	mux.ServeHTTP(rec, req)

	resp := rec.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err, "failed to read response body")

	d := map[string]interface{}{}
	err = json.Unmarshal(body, &d)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
}
