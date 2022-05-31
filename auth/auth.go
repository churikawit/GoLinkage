
package auth

import (
	"fmt"
	"time"
	"net/http"
	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
)

func HandleAccessToken(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
	})
	ss, err := token.SignedString([]byte("mysignature"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": ss,
	})
}

func Protect(token string) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _,ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("mysignature"), nil
	})

	return err
}
