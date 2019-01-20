package api

import (
	"fmt"
	"net/http"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/util"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// 建struct
// 驗證username, password

type auth struct {
	Username string `valid: Required;MaxSize(50)`
	Password string `valid: Required;MaxSize(50)`
}

func GetAuth(c *gin.Context) {
	// Get username, password
	username := c.Query("username")
	password := c.Query("password")
	valid := validation.Validation{}

	valid.Required(username, "username").Message("User is required field")
	valid.Required(password, "password").Message("Password is required field")
	code := e.INVALID_PARAMS
	var data = make(map[string]interface{})
	if !valid.HasErrors() {
		// 驗證無錯
		if models.CheckAuth(username, password) {
			// User is exist
			// Generate token
			token, err := util.GenerateToken(username, password)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
			} else {
				data["token"] = token
				code = e.SUCCESS
			}
		} else {
			code = e.ERROR_AUTH
		}
	} else {
		for _, err := range valid.Errors {
			fmt.Printf("CheckAuth validation error: ", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": data,
	})
}
