package controllers_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/controllers"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
	"github.com/kiselev-nikolay/go-shop-api-template/common/detailres"
	"github.com/kiselev-nikolay/go-shop-api-template/tools/reqvalid"
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

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(reqvalid.ValidateWare())
	c := controllers.New(db)
	router.POST("/product", c.Product)

	getProductsCount := func() uint {
		var x uint
		db.Raw("SELECT count(*) FROM products").Scan(&x)
		return x
	}

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
		req.Header.Add("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := detailres.New("")
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 400, w.Result().StatusCode)
		assert.DeepEqual(t, detailres.New("creator invalid"), res)
	})

	t.Run("Create", func(t *testing.T) {
		wasCount := getProductsCount()
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
		res := detailres.New("")
		bodyData, err := ioutil.ReadAll(w.Body)
		assert.NilError(t, err)
		err = json.Unmarshal(bodyData, res)
		assert.NilError(t, err)
		assert.Equal(t, 200, w.Result().StatusCode)
		assert.DeepEqual(t, res, detailres.New("created"))
		assert.Equal(t, 200, w.Result().StatusCode)
		nowCount := getProductsCount()
		assert.Equal(t, uint(1), nowCount-wasCount)
	})

	t.Run("Product categories create once", func(t *testing.T) {
		wasCount := getProductsCount()
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
			Categories: []controllers.CategoryReq{
				{Name: "strange"},
				{Name: "vintage"},
				{Name: "for kids"},
				{Name: "under bad monsters"},
			},
		})

		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(data))
			req.Header.Add("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		nowCount := getProductsCount()
		assert.Equal(t, uint(5), nowCount-wasCount)

		var categoriesCount uint
		db.Raw("SELECT count(*) FROM categories").Scan(&categoriesCount)
		assert.Equal(t, uint(4), categoriesCount)
	})
}
