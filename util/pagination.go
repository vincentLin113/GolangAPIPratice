package util

import (
	"vincent-gin-go/pkg/setting"

	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

// GetPage: Query`page`
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * setting.AppSetting.PageSize
	}
	return result
}
