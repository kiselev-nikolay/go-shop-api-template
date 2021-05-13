package catalogue

import "github.com/gin-gonic/gin"

func Views(group *gin.RouterGroup) {
	group.GET("/product", viewProduct)
}

type productRes struct {
	name string
}

func viewProduct(g *gin.Context) {
	g.JSON(200, productRes{
		name: "hey",
	})
}
