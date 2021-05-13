package catalogue_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue"
	"github.com/kiselev-nikolay/go-test-docker-dependencies/testdep"
	"gotest.tools/assert"
)

func MustCreateDB() (*gorm.DB, func()) {
	dbPort, err := testdep.FindFreePort()
	if err != nil {
		log.Fatalln(err)
	}
	dockerPg := testdep.Postgres{
		Port:     dbPort,
		User:     "test",
		Password: "test",
		Database: "test",
	}
	stop, err := dockerPg.Run(10)
	if err != nil {
		log.Fatalln(err)
	}
	if err != nil {
		log.Fatalln(err)
	}
	db, err := gorm.Open(postgres.Open(dockerPg.ConnString()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return db, func() {
		err := stop()
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func TestViews(t *testing.T) {
	db, stop := MustCreateDB()
	defer stop()

	err := db.AutoMigrate(&catalogue.Creator{})
	assert.NilError(t, err)
	err = db.AutoMigrate(&catalogue.Category{})
	assert.NilError(t, err)
	err = db.AutoMigrate(&catalogue.Product{})
	assert.NilError(t, err)

	creator := &catalogue.Creator{
		Name: "test creator",
	}
	db.Create(creator)

	product := catalogue.Product{
		Code:      "abc",
		Price:     100,
		Name:      "test",
		CreatorID: creator.ID,
		Categories: []catalogue.Category{
			{Name: "test cat 1"},
			{Name: "test cat 2"},
			{Name: "test cat 3"},
		},
	}
	db.Create(&product)
	db.Save(&product)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	v := catalogue.NewViews(db)
	router.GET("/product", v.Product)

	t.Run("Wrong request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/product", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &catalogue.ProductRes{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 400, w.Result().StatusCode)
		assert.DeepEqual(t, res, &catalogue.ProductRes{})
	})

	t.Run("Simple get", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/product?id=%d", product.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &catalogue.ProductRes{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 200, w.Result().StatusCode)
		assert.DeepEqual(t, res, &catalogue.ProductRes{Name: "test"})
	})

	t.Run("Not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/product?id=9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &catalogue.ProductRes{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 404, w.Result().StatusCode)
		assert.DeepEqual(t, res, &catalogue.ProductRes{})
	})
}
