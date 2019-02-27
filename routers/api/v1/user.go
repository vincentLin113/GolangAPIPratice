package v1

import (
	"net/http"
	"vincent-gin-go/pkg/app"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/service/user_service"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// GetUser 獲取指定ID的User
func GetUser(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	validation := validation.Validation{}
	validation.Min(id, 0, "id").Message("ID不得小於1")

	// 1. 驗證參數
	if validation.HasErrors() {
		// 出現錯誤
		app.MarkErrors(validation.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, validation.Errors)
	}

	// 2. 找對應ID的User
	userService := user_service.User{ID: id}
	exist, err := userService.ExistByID()
	if err != nil {
		// 找不到對應ID的User
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	if !exist {
		// 若不存在而非error, http為正常
		appG.Response(http.StatusOK, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	user, err := userService.Get()
	if user.State == 0 {
		appG.Response(http.StatusOK, e.ERROR_GET_USER_BAN_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, user)
}

func AddUser(c *gin.Context) {
	appG := app.Gin{c}
	name := c.Query("name")
	email := c.Query("email")
	password := c.Query("password")
	validation := validation.Validation{}
	validation.Required(name, "name").Message("Name field is required")
	validation.Email(email, "email").Message("Email格式錯誤")
	validation.MinSize(email, 5, "email").Message("Email字數過短")
	validation.MinSize(password, 6, "password").Message("密碼過短")

	if validation.HasErrors() {
		app.MarkErrors(validation.Errors)
		appG.Response(http.StatusBadRequest, e.ERROR_ADD_USER_FAIL, nil)
		return
	}

	// 驗證無錯, 則新增用戶
	userService := user_service.User{
		Name:     name,
		Email:    email,
		Password: password,
	}
	if userService.ExistByEmail() {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_USER_DUPLICATED_EMAIL_ERROR, nil)
		return
	} else if userService.ExistByName() {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_USER_DUPLICATED_NAME_ERROR, nil)
		return
	}
	err := userService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_USER_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, userService)
}
