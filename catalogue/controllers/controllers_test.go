package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/controllers"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
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
	c := controllers.New(db)
	router.POST("/product", c.Product)

	t.Run("Create with invalid creator", func(t *testing.T) {
		data, _ := json.Marshal(controllers.ProductReq{
			Code:  "scp096",
			Price: 200,
			Name:  "Strange Human Size Box",
			Creator: controllers.CreatorReq{
				ID: 29,
			},
			Categories: []controllers.CategoryReq{},
		})
		req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(data))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &controllers.Detail{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 400, w.Result().StatusCode)
		assert.DeepEqual(t, res, &controllers.Detail{"creator invalid"})
	})

	t.Run("Create", func(t *testing.T) {
		creator := &models.Creator{
			Name: "unknown",
		}
		db.Create(creator)

		data, _ := json.Marshal(controllers.ProductReq{
			Code:  "scp096",
			Price: 200,
			Name:  "Strange Human Size Box",
			Creator: controllers.CreatorReq{
				ID: creator.ID,
			},
			Categories: []controllers.CategoryReq{},
		})
		req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(data))
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &controllers.Detail{}
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		fmt.Println(string(bodyData))
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 200, w.Result().StatusCode)
		assert.DeepEqual(t, res, &controllers.Detail{"created"})
	})
}
