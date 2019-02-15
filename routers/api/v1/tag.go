package v1

import (
	"net/http"
	"vincent-gin-go/pkg/app"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/setting"
	"vincent-gin-go/service/tag_service"
	"vincent-gin-go/util"

	"github.com/astaxie/beego/validation"

	"github.com/Unknwon/com"
	"github.com/gin-gonic/gin"
)

// GetTags Get all tags
// 1. 取得參數, 2. 組建tag_service.Tag, 3. 利用組建的service_tag來獲取Database中的資料, 4. JSON Response
func GetTags(c *gin.Context) {
	appG := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}
	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}
	tags, err := tagService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(255)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0, 1)"`
}

// @Summary 新增文章標籤
// @Produce json
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param created_by query int false "CreatedBy"
// @Success 200 {string} json "{"code": 200, "data": {}, "msg": "ok"}"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	// 1. 取得參數
	// 2. 產生AddTagForm的Model
	// 3. 格式驗證
	// 4. 產生tag_service.Tag的Model
	// 5. 使用tag_service來與Database交互
	appG := app.Gin{c}
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	form := AddTagForm{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}
	httpCode, errorCode := app.BindAndValid(c, &form)
	if errorCode != e.SUCCESS {
		appG.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      name,
		CreatedBy: createdBy,
		State:     state,
	}

	err := tagService.Add()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0, 1)"`
}

// EditTag To edit tag
// @Summary 編輯文章標籤
// @Produce json
// @Param name query string true "Name"
// @Param modifiedBy query string true "Modified_By"
// @Param id query int true "ID"
// @Param state query int false "State"
// Success 200 {string} json "{"code": 200, "msg": "ok", "data": nil}"
func EditTag(c *gin.Context) {
	// Get `id, name, modifiedBy`
	// valid them
	// 驗證無錯後, 建立tag_service.Tag
	// 驗證此ID是否存在
	// Handle error
	// JSON Response
	appG := app.Gin{c}
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")
	id := com.StrTo(c.Param("id")).MustInt()
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()
	form := EditTagForm{
		ID:         id,
		Name:       name,
		ModifiedBy: modifiedBy,
	}
	if state > -1 {
		form.State = state
	}
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       name,
		ModifiedBy: modifiedBy,
		State:      state,
	}
	exist, err := tagService.ExistById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG, nil)
		return
	}
	if !exist {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	err = tagService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteTag To delete tag
// @Summary 刪除標籤
// @Produce json
// @Param id query int true "ID"
// Success 200 {string}
func DeleteTag(c *gin.Context) {
	// 取得ID
	// 驗證ID
	// Handle validation errors
	// 組建tag_service.Tag
	// 檢查是否存在
	// Handle errors
	// 嘗試刪除
	// Handle error
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID不得低於1")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}
	tagService := tag_service.Tag{
		ID: id,
	}
	exists, err := tagService.ExistById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	success, err := tagService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	if !success {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
