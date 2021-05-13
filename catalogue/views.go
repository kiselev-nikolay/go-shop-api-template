package catalogue

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewViews(DB *gorm.DB) *Views {
	return &Views{DB: DB}
}

type Views struct {
	DB *gorm.DB
}

func (v *Views) Connect(group *gin.RouterGroup) {
	group.GET("/product", v.Product)
}

type ProductRes struct {
	Name string
}

func (v *Views) Product(g *gin.Context) {
	id, err := strconv.Atoi(g.Query("id"))
	if err != nil {
		g.JSON(400, ProductRes{})
		return
	}
	p := &Product{}
	result := v.DB.First(p, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		g.JSON(404, ProductRes{})
		return
	}
	g.JSON(200, ProductRes{
		Name: p.Name,
	})
}
