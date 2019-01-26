package app

import (
	"vincent-gin-go/pkg/e"

	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"code": httpCode,
		"msg":  e.GetMessage(errCode),
		"data": data,
	})

	return
}
