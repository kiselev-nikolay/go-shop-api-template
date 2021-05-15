package reqvalid_test

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kiselev-nikolay/go-shop-api-template/tools/reqvalid"
	"gotest.tools/assert"
)

func TestValidateWare(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ge := gin.New()
	ge.Use(reqvalid.ValidateWare())
	ge.Any("/", func(g *gin.Context) { g.Status(204) })

	t.Run("Just get", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		ge.ServeHTTP(rw, req)
		assert.Equal(t, 204, rw.Result().StatusCode)
	})
	t.Run("Just post", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", nil)
		ge.ServeHTTP(rw, req)
		assert.Equal(t, 400, rw.Result().StatusCode)
		b, err := ioutil.ReadAll(rw.Body)
		assert.NilError(t, err)
		assert.Equal(t, `{"detail":"no content type header"}`, string(b))
	})
	t.Run("Wrong post", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Content-Type", "?")
		ge.ServeHTTP(rw, req)
		assert.Equal(t, 400, rw.Result().StatusCode)
		b, err := ioutil.ReadAll(rw.Body)
		assert.NilError(t, err)
		assert.Equal(t, `{"detail":"invalid content type header"}`, string(b))
	})
	t.Run("Unsupported type post", func(t *testing.T) {
		rw := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Content-Type", "text/php")
		ge.ServeHTTP(rw, req)
		assert.Equal(t, 400, rw.Result().StatusCode)
		b, err := ioutil.ReadAll(rw.Body)
		assert.NilError(t, err)
		assert.Equal(t, `{"detail":"unexpected content type header"}`, string(b))
	})
	t.Run("Nice post", func(t *testing.T) {
		rw := httptest.NewRecorder()
		body := bytes.NewBufferString(`{"msg": "hello, my name is Nikolai from http://nikolai.works"}`)
		req := httptest.NewRequest("POST", "/", body)
		req.Header.Set("Content-Type", "application/json ")
		ge.ServeHTTP(rw, req)
		assert.Equal(t, 204, rw.Result().StatusCode)
	})
}
