package v1

import (
	"net/http"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/setting"
	"vincent-gin-go/util"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

// Get all tags
func GetTags(c *gin.Context) {
	name := c.Param("name")
	maps := make(map[string]interface{})
	data := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS
	pageNum := util.GetPage(c)
	data["lists"], _ = models.GetTags(pageNum, setting.AppSetting.PageSize, maps)
	data["total"], _ = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": data,
	})
}

// Add tag
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")
	valid := validation.Validation{}
	valid.Required(name, "name").Message("名稱不能為空")
	valid.MaxSize(name, 100, "name").Message("名稱最長為100字符")
	valid.Required(createdBy, "created_by").Message("創建人不能為空")
	valid.MaxSize(createdBy, 100, "created_by").Message("創建人最長為100字元")
	valid.Range(state, 0, 1, "state").Message("狀態只允許0, 1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		exist, _ := models.ExistTagByName(name)
		if !exist {
			// 驗證未錯且不存在此Tag
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			// 驗證未錯但不存在
			code = e.ERROR_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": make(map[string]string),
	})
}

// EditTag To edit tag
func EditTag(c *gin.Context) {
	// Get `id, name, modifiedBy`
	// valid them
	// error
	// c.JSON
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")
	state := com.StrTo(c.DefaultQuery("state", "-1")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不得為空")
	valid.Min(id, 1, "id").Message("ID需大於0")
	valid.Required(name, "name").Message("名字不得為空")
	valid.Required(modifiedBy, "modified_by").Message("修改者不得為空")

	code := e.INVALID_PARAMS

	if !valid.HasErrors() {
		exist, _ := models.ExistTagById(id)
		if !exist {
			// 驗證無措 但 不存在
			code = e.ERROR_EXIST_TAG_FAIL
		} else {
			// 驗證無錯 且 存在
			code = e.SUCCESS
			data := make(map[string]interface{})
			if name != "" {
				data["name"] = name
			}
			if modifiedBy != "" {
				data["modified_by"] = modifiedBy
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
	})

	return
}

// DeleteTag To delete tag
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("Id不得為空")
	valid.Min(id, 1, "id").Message("ID必須大於1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		exist, _ := models.ExistTagById(id)
		if !exist {
			// 驗證無錯 但Tag不存在
			code = e.ERROR_DELETE_TAG_FAIL
		} else {
			code = e.SUCCESS
			_, err := models.DeleteTag(id)
			if err != nil {
				code = e.ERROR_DELETE_TAG_FAIL
			}
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
	})
}
