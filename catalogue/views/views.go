package views

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
	"gorm.io/gorm"
)

func New(DB *gorm.DB) *views {
	return &views{DB: DB}
}

type views struct {
	DB *gorm.DB
}

func (v *views) Connect(group *gin.RouterGroup) {
	group.GET("/product", v.Product)
}

type CreatorRes struct {
	Name string `json:"name"`
}

type CategoryRes struct {
	Name string `json:"name"`
}

type ProductRes struct {
	Code       string        `json:"code"`
	Price      string        `json:"price"`
	Name       string        `json:"name"`
	Creator    CreatorRes    `json:"creator"`
	Categories []CategoryRes `json:"categories"`
}

func (v *views) Product(g *gin.Context) {
	id, err := strconv.Atoi(g.Query("id"))
	if err != nil {
		g.JSON(400, ProductRes{})
		return
	}
	p := &models.Product{}
	result := v.DB.First(p, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		g.JSON(404, ProductRes{})
		return
	}
	g.JSON(200, ProductRes{
		Name: p.Name,
	})
}
