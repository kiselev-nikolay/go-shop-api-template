package user

import (
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	key = "gsat:user:authored"
)

var (
	ErrNoContentType = errors.New("no content type header")
)

type User struct {
	Name         string
	PasswordHash string
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUser(g *gin.Context) (user *User, ok bool) {
	value, ok := g.Get(key)
	if !ok {
		return nil, false
	}
	user, castOk := value.(*User)
	if !castOk {
		return nil, false
	}
	return user, true
}

func UserWare() gin.HandlerFunc {
	return func(g *gin.Context) {
		// Todo ...
		g.Next()
	}
}

func MockUserWare() gin.HandlerFunc {
	return func(g *gin.Context) {
		hash, _ := hashPassword("123456")
		g.Set(key, &User{
			Name:         "Nikolai Kiselev",
			PasswordHash: hash,
		})
		g.Next()
	}
}
