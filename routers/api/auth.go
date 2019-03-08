package api

import (
	"net/http"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/util"
	v1 "vincent-gin-go/routers/api/v1"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// 建struct
// 驗證username, password

type auth struct {
	Username string `valid: Required;MaxSize(50)`
	Password string `valid: Required;MaxSize(50)`
}

func Login(c *gin.Context) {
	// Get email, password
	email := c.Query("email")
	password := c.Query("password")
	valid := validation.Validation{}

	valid.Required(email, "email").Message("Email is required field")
	valid.Required(password, "password").Message("Password is required field")
	code := e.INVALID_PARAMS
	var data = make(map[string]interface{})
	if !valid.HasErrors() {
		// 驗證無錯
		if models.CheckAuth(email, password) {
			// User is exist
			// Generate token
			token, err := util.GenerateToken(email, password)
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
		for err := range valid.Errors {
			logging.Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": data,
	})
}

func SignUp(c *gin.Context) {
	v1.AddUser(c)
}

func ActivateUser(c *gin.Context) {
	
}