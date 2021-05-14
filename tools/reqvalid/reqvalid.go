package reqvalid

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoContentType = errors.New("no content type")
)

func Validate(g *gin.Context) error {
	if g.Request.Method == "POST" {
		if g.ContentType() == "" {
			return ErrNoContentType
		}
	}
	return nil
}
