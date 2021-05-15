package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
	"github.com/kiselev-nikolay/go-shop-api-template/common/detailres"
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
	ID uint `json:"id" binding:"required"`
}

type CategoryReq struct {
	Name string `json:"name" binding:"required"`
}

type ProductReq struct {
	Code       string        `json:"code" binding:"required"`
	Price      uint          `json:"price" binding:"required"`
	Name       string        `json:"name" binding:"required"`
	Creator    CreatorReq    `json:"creator" binding:"required"`
	Categories []CategoryReq `json:"categories"`
}

type ProductRes struct {
	ID uint `json:"id"`
}

func (c *ctrl) Product(g *gin.Context) {
	req := &ProductReq{}
	bindErr := g.Bind(req)
	if bindErr != nil {
		g.JSON(400, detailres.New("wrong body"+bindErr.Error()))
		return
	}
	creatorSearch := c.DB.First(&models.Creator{}, req.Creator.ID)
	if errors.Is(creatorSearch.Error, gorm.ErrRecordNotFound) {
		g.JSON(400, detailres.New("creator invalid"))
		return
	}
	cats := make([]models.Category, 0)
	for _, cat := range req.Categories {
		dbCat := &models.Category{
			Name: cat.Name,
		}
		c.DB.FirstOrCreate(dbCat, "name = ?", cat.Name)
		cats = append(cats, *dbCat)
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
		g.JSON(400, detailres.New("not created"))
		return
	}
	g.JSON(200, &ProductRes{ID: p.ID})
}
