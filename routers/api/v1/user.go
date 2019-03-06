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

type AddUserForm struct {
	Name     string `form:"name" valid:"Required;MaxSize(255)"`
	Email    string `form:"email" valid:"Required;MaxSize(255);MinSize(5)"`
	Password string `form:"password" valid:"Required;MaxSize(255);MinSize(6)"`
}

func AddUser(c *gin.Context) {
	appG := app.Gin{c}
	name := c.Query("name")
	email := c.Query("email")
	password := c.Query("password")
	form := AddUserForm{
		Name:     name,
		Email:    email,
		Password: password,
	}
	httpCode, errCode := app.BindAndValid(c, form)
	if httpCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
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

type EditUserForm struct {
	Name     string `form:"name" valid:"Required;MaxSize(255);MinSize(1)"`
	Email    string `form:"email" valid:"Required;MaxSize(255);MinSize(5)"`
	Password string `form:"password" valid:"Required;MaxSize(255);MinSize(6)"`
	State    int    `form:"state" valid:"Required;Min(0);Max(1)"`
}

func EditUser(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	email := c.Query("email")
	password := c.Query("password")
	state := com.StrTo(c.Query("state")).MustInt()
	form := EditUserForm{
		Name:     name,
		Email:    email,
		Password: password,
		State:    state,
	}
	httpCode, errCode := app.BindAndValid(c, form)
	if httpCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	userService := user_service.User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
		State:    state,
	}
	exist, err := userService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_USER_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_USER_FAIL, nil)
		return
	}
	err = userService.Edit()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_EDIT_USER_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteUser To delete user according `id` parameter.
func DeleteUser(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	validation := validation.Validation{}
	validation.Min(id, 0, "id").Message("ID不得小於1")
	if validation.HasErrors() {
		appG.Response(http.StatusBadRequest, e.ERROR_DELETE_USER_FAIL, nil)
		return
	}
	userService := user_service.User{ID: id}
	exist, err := userService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USER_FAIL, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusOK, e.ERROR_DELETE_USER_FAIL, nil)
		return
	}

	err = userService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_USER_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
