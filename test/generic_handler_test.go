package test

import (
	_ "github.com/iamrz1/ab-auth/docs"
	infraMongo "github.com/iamrz1/ab-auth/infra/mongo"
)

const (
	name = "Product Y"
	slug = "product-y"
)

var db *infraMongo.Mongo

//func getServer() (*chi.Mux, error) {
//	err := config.LoadConfig()
//	if err != nil {
//		log.Println("could not load one or more config")
//		return nil, err
//	}
//	//rLogger
//	rLogger := logger.DefaultOutStructLogger
//	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
//	defer cancel()
//
//	cfg := config.GetConfig()
//
//	gracefulTimeout := time.Second * time.Duration(cfg.GracefulTimeout)
//	db, err = infraMongo.New(ctx, cfg.DSN, cfg.Database, gracefulTimeout)
//	if err != nil {
//		return nil, err
//	}
//
//	log.Println("db initialized")
//	//cfg := config.getConfig()
//
//	//addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
//	svc := service.SetupServiceConfig(cfg, db, rLogger)
//	handler, err := api.SetupRouter(cfg, svc, rLogger)
//	if err != nil {
//		log.Println("cant setup router:", err)
//		return nil, err
//	}
//
//	return handler, nil
//}
//
//func closeResource() {
//	db.Close(context.Background())
//}
//
//func TestCreateObject(t *testing.T) {
//	product := model.Generic{
//		Name:         name,
//		StringField:  "dfdg",
//		DecimalField: 0.877,
//		IntegerField: 23,
//	}
//
//	b, err := json.Marshal(&product)
//	assert.NoError(t, err, "marshal failed")
//
//	req := httptest.NewRequest(http.MethodPost, "/api/v1/public/generics", bytes.NewReader(b))
//	rec := httptest.NewRecorder()
//
//	mux, err := getServer()
//	assert.NoError(t, err, "failed to get server")
//
//	defer closeResource()
//
//	mux.ServeHTTP(rec, req)
//
//	resp := rec.Result()
//	body, err := ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	d := map[string]interface{}{}
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.EqualValuesf(t, resp.StatusCode, http.StatusCreated, "failed with response body: %v", d)
//}
//
//func TestGetObjects(t *testing.T) {
//	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/generics", nil)
//	rec := httptest.NewRecorder()
//
//	mux, err := getServer()
//	assert.NoError(t, err, "failed to get server")
//	defer closeResource()
//
//	mux.ServeHTTP(rec, req)
//
//	resp := rec.Result()
//	body, err := ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	d := map[string]interface{}{}
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
//
//	res := response.GenericListSuccessRes{}
//	merr := json.Unmarshal(body, &res)
//	if merr != nil {
//		t.Fatal(merr)
//	}
//
//	//assert.EqualValues(t, 1, len(res.Data.Generics))
//}
//
//func TestGetSingleObject(t *testing.T) {
//	req := httptest.NewRequest(http.MethodGet, "/api/v1/public/generics/"+slug, nil)
//	rec := httptest.NewRecorder()
//
//	mux, err := getServer()
//	assert.NoError(t, err, "failed to get server")
//	defer closeResource()
//
//	mux.ServeHTTP(rec, req)
//
//	resp := rec.Result()
//	body, err := ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	d := map[string]interface{}{}
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
//}
//
//func TestUpdateObject(t *testing.T) {
//	product := model.Generic{
//		IntegerField: 100,
//	}
//
//	b, err := json.Marshal(&product)
//	assert.NoError(t, err, "marshal failed")
//
//	req := httptest.NewRequest(http.MethodPatch, "/api/v1/public/generics/"+slug, bytes.NewReader(b))
//	rec := httptest.NewRecorder()
//
//	mux, err := getServer()
//	assert.NoError(t, err, "failed to get server")
//	defer closeResource()
//
//	mux.ServeHTTP(rec, req)
//
//	resp := rec.Result()
//	body, err := ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	d := map[string]interface{}{}
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
//}
//
//func TestPurgeSingleObject(t *testing.T) {
//	req := httptest.NewRequest(http.MethodDelete, "/api/v1/public/generics/purge/"+slug, nil)
//	rec := httptest.NewRecorder()
//
//	mux, err := getServer()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer closeResource()
//
//	mux.ServeHTTP(rec, req)
//
//	resp := rec.Result()
//	body, err := ioutil.ReadAll(resp.Body)
//	assert.NoError(t, err, "failed to read response body")
//
//	d := map[string]interface{}{}
//	err = json.Unmarshal(body, &d)
//	if err != nil {
//		t.Fatal(err)
//	}
//	assert.EqualValuesf(t, resp.StatusCode, http.StatusOK, "failed with response body: %v", d)
//}
