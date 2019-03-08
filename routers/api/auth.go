package api

import (
	"fmt"
	"net/http"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/app"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/util"
	v1 "vincent-gin-go/routers/api/v1"
	"vincent-gin-go/service/user_service"

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
	appG := app.Gin{c}
	email := c.Query("email")
	token := c.Query("token")
	validation := validation.Validation{}
	validation.Email(email, "email").Message("EMAIL_FORMAT_ERROR")
	validation.Required(token, "token").Message("REQUIRED_FIELD_ERROR")
	if validation.HasErrors() {
		app.MarkErrors(validation.Errors)
		appG.Response(http.StatusBadRequest, e.ERROR_ACTIVATE_USER_FEILD, nil)
		return
	}
	userService := user_service.User{Email: email}
	exist := userService.ExistByEmail()
	if !exist {
		appG.Response(http.StatusInternalServerError, e.ERROR_ACTIVATE_USER_FEILD, nil)
		return
	}
	user, err := userService.GetByEmail()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ACTIVATE_USER_FEILD, nil)
		return
	}
	fmt.Println("\nUSERTOKEN: ", user.Token)
	fmt.Println("\nTOKEN: ", token)
	if user.Token != token {
		appG.Response(http.StatusInternalServerError, e.ERROR_ACTIVATE_USER_FEILD, nil)
		return
	}
	userService.ID = user.ID
	err = userService.Activate()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ACTIVATE_USER_FEILD, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
