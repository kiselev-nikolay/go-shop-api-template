package catalogue_test

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
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/views"
	"github.com/kiselev-nikolay/go-shop-api-template/tools/reqvalid"
	"github.com/kiselev-nikolay/go-shop-api-template/tools/testtools"
	"gotest.tools/assert"
)

func TestDomain(t *testing.T) {
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

	v := views.New(db)
	v.Connect(router.Group("views"))

	c := controllers.New(db)
	c.Connect(router.Group("ctrls"))

	getProductsCount := func() uint {
		var x uint
		db.Raw("SELECT count(*) FROM products").Scan(&x)
		return x
	}

	assert.Equal(t, uint(0), getProductsCount(), "seems like your database mock isnt clean")

	t.Run("Product", func(t *testing.T) {
		creator := &models.Creator{
			Name: "Game Stop",
		}
		db.Create(creator)

		postProductReq := &controllers.ProductReq{
			Code:  "AoE2-gs",
			Price: 5,
			Name:  "Age of Empires II",
			Creator: controllers.CreatorReq{
				ID: creator.ID,
			},
			Categories: []controllers.CategoryReq{},
		}
		data, _ := json.Marshal(postProductReq)
		postReq, _ := http.NewRequest("POST", "/ctrls/product", bytes.NewBuffer(data))
		postReq.Header.Add("Content-Type", "application/json")
		postRes := httptest.NewRecorder()
		router.ServeHTTP(postRes, postReq)
		assert.Equal(t, 200, postRes.Result().StatusCode)
		assert.Equal(t, uint(1), getProductsCount())
		postBodyData, err := ioutil.ReadAll(postRes.Body)
		assert.NilError(t, err)
		postResBody := &controllers.ProductRes{}
		err = json.Unmarshal(postBodyData, postResBody)
		assert.NilError(t, err)

		getURL := fmt.Sprintf("/views/product?id=%d", postResBody.ID)
		getReq, _ := http.NewRequest("GET", getURL, nil)
		getRes := httptest.NewRecorder()
		router.ServeHTTP(getRes, getReq)
		assert.Equal(t, 200, getRes.Result().StatusCode)
		getBodyData, err := ioutil.ReadAll(getRes.Body)
		fmt.Println(string(getBodyData))
		assert.NilError(t, err)
		getProductRes := &views.ProductRes{}
		err = json.Unmarshal(getBodyData, getProductRes)
		assert.NilError(t, err)
		assert.DeepEqual(t, getProductRes, &views.ProductRes{
			Code:  postProductReq.Code,
			Price: postProductReq.Price,
			Name:  postProductReq.Name,
			Creator: views.CreatorRes{
				Name: creator.Name,
			},
			Categories: []views.CategoryRes{},
		})
	})
}
