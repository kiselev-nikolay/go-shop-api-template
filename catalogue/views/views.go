package views

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/catalogue/models"
	"github.com/kiselev-nikolay/go-shop-api-template/common/detailres"
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
	Price      uint          `json:"price"`
	Name       string        `json:"name"`
	Creator    CreatorRes    `json:"creator"`
	Categories []CategoryRes `json:"categories"`
}

func (v *views) Product(g *gin.Context) {
	id, err := strconv.Atoi(g.Query("id"))
	if err != nil {
		g.JSON(400, detailres.New("id param is missing or invalid"))
		return
	}
	p := &models.Product{}
	q := v.DB.First(p, id)
	if errors.Is(q.Error, gorm.ErrRecordNotFound) {
		g.JSON(404, detailres.New("product not found"))
		return
	}
	creator := &models.Creator{}
	creatorQ := v.DB.First(creator, p.CreatorID)
	if errors.Is(creatorQ.Error, gorm.ErrRecordNotFound) {
		g.JSON(404, detailres.New("creator not found"))
		return
	}
	findCategoriesErr := v.DB.Model(&p).Association("Categories").Find(&p.Categories)
	if findCategoriesErr != nil {
		g.JSON(404, detailres.New("categories not found"))
		return
	}
	r := ProductRes{
		Code:  p.Code,
		Price: p.Price,
		Name:  p.Name,
		Creator: CreatorRes{
			Name: creator.Name,
		},
		Categories: make([]CategoryRes, 0, len(p.Categories)),
	}
	for _, cat := range p.Categories {
		r.Categories = append(r.Categories, CategoryRes{
			Name: cat.Name,
		})
	}
	g.JSON(200, r)
}
