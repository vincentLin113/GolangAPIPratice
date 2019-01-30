package app

import (
	"net/http"
	"vincent-gin-go/pkg/e"

	"github.com/astaxie/beego/validation"

	"github.com/gin-gonic/gin"
)

// BindAndValid 返回http code + errorCode
// form: 輸入含有驗證的struct
func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}
	if !check {
		MarkErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
