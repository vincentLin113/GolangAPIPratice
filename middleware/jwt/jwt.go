package jwt

import (
	"net/http"
	"time"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/util"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	method := func(c *gin.Context) {
		var code int
		var data interface{}
		code = e.SUCCESS
		token := c.GetHeader("token")
		if token == "" {
			// 若無輸入Token
			code = e.ERROR_AUTH_NEED_TOKEN
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMessage(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
	return method
}
