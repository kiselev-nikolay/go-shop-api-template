package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
	"gorm.io/gorm"
)

func New(DB *gorm.DB) *ctrl {
	return &ctrl{DB: DB}
}

type ctrl struct {
	DB *gorm.DB
}

func (c *ctrl) Connect(group *gin.RouterGroup) {
	group.POST("/product", c.Product)
}

type CreatorReq struct {
	ID uint `json:"id"`
}

type CategoryReq struct {
	Name string `json:"name"`
}

type ProductReq struct {
	Code       string        `json:"code"`
	Price      uint          `json:"price"`
	Name       string        `json:"name"`
	Creator    CreatorReq    `json:"creator"`
	Categories []CategoryReq `json:"categories"`
}

type Detail struct {
	Value string `json:"detail"`
}

func (c *ctrl) Product(g *gin.Context) {
	req := &ProductReq{}
	b, _ := g.Request.GetBody()
	x, _ := ioutil.ReadAll(b)
	fmt.Println(string(x))
	bindErr := g.Bind(req)
	if bindErr != nil {
		g.JSON(400, Detail{"wrong body"})
		return
	}
	fmt.Println(req.Creator.ID)
	creatorSearch := c.DB.First(&models.Creator{}, req.Creator.ID)
	if errors.Is(creatorSearch.Error, gorm.ErrRecordNotFound) {
		g.JSON(400, Detail{"creator invalid"})
		return
	}
	cats := make([]models.Category, 0)
	for _, cat := range req.Categories {
		cats = append(cats, models.Category{
			Name: cat.Name,
		})
	}
	p := &models.Product{
		Code:       req.Code,
		Price:      req.Price,
		Name:       req.Name,
		CreatorID:  uint(req.Creator.ID),
		Categories: cats,
	}
	productCreate := c.DB.Create(p)
	if errors.Is(productCreate.Error, gorm.ErrRecordNotFound) {
		g.JSON(400, Detail{"not created"})
		return
	}
	g.JSON(200, Detail{"created"})
}
