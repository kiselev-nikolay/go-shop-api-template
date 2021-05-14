package views_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/views"
	"github.com/kiselev-nikolay/go-shop-api-template/tools/testtools"
	"gotest.tools/assert"
)

func TestViews(t *testing.T) {
	db, stop := testtools.MustCreateDB()
	defer stop()

	err := db.AutoMigrate(&models.Creator{})
	assert.NilError(t, err)
	err = db.AutoMigrate(&models.Category{})
	assert.NilError(t, err)
	err = db.AutoMigrate(&models.Product{})
	assert.NilError(t, err)

	creator := &models.Creator{
		Name: "test creator",
	}
	db.Create(creator)

	product := models.Product{
		Code:      "abc",
		Price:     100,
		Name:      "test",
		CreatorID: creator.ID,
		Categories: []models.Category{
			{Name: "test cat 1"},
			{Name: "test cat 2"},
			{Name: "test cat 3"},
		},
	}
	db.Create(&product)
	db.Save(&product)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	v := views.New(db)
	router.GET("/product", v.Product)

	t.Run("Wrong request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/product", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &views.ProductRes{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 400, w.Result().StatusCode)
		assert.DeepEqual(t, res, &views.ProductRes{})
	})

	t.Run("Simple get", func(t *testing.T) {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/product?id=%d", product.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &views.ProductRes{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 200, w.Result().StatusCode)
		assert.DeepEqual(t, res, &views.ProductRes{Name: "test"})
	})

	t.Run("Not found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/product?id=9999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &views.ProductRes{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 404, w.Result().StatusCode)
		assert.DeepEqual(t, res, &views.ProductRes{})
	})
}
