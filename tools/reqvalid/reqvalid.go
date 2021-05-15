package reqvalid

import (
	"errors"
	"mime"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kiselev-nikolay/go-shop-api-template/common/detailres"
)

var (
	ErrNoContentType         = errors.New("no content type header")
	ErrInvalidContentType    = errors.New("invalid content type header")
	ErrUnexpectedContentType = errors.New("unexpected content type header")
)

func Validate(g *gin.Context) error {
	if g.Request.Method == "POST" {
		if g.ContentType() == "" {
			return ErrNoContentType
		}
		mediaType, _, err := mime.ParseMediaType(g.ContentType())
		if err != nil {
			return ErrInvalidContentType
		}
		switch mediaType {
		case binding.MIMEJSON,
			binding.MIMEYAML,
			binding.MIMEPROTOBUF,
			binding.MIMEXML, binding.MIMEXML2,
			binding.MIMEMSGPACK, binding.MIMEMSGPACK2,
			binding.MIMEMultipartPOSTForm, binding.MIMEPOSTForm:
			return nil
		default:
			return ErrUnexpectedContentType
		}
	}
	return nil
}

func ValidateWare() gin.HandlerFunc {
	return func(g *gin.Context) {
		err := Validate(g)
		if err != nil {
			g.JSON(400, detailres.New(err.Error()))
			return
		}
		g.Next()
	}
}
