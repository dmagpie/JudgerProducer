package router

import (
	"JudgerProducer/config"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var secret string

func init() {
	s, err := config.GetConfig("SECRET")
	if err != nil {
		panic(err)
	}
	secret = s
	if secret == "" {
		panic(errors.New("empty secret is set"))
	}
}

// BasicAuth 基本鉴权
func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusNonAuthoritativeInfo, gin.H{
				"errmsg": "need-authorization",
			})
			return
		}

		if token != secret {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errmsg": "authorization-failed",
			})
			return
		}

		c.Next()
	}
}
