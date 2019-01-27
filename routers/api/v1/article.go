package v1

import (
	"net/http"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/app"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/service/article_service"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必須大於0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	article, err := articleService.Get()
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必須大於0")
	var data = make(map[string]interface{})

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		// 若無錯誤
		success, _ := models.ExistArticleById(id)
		if success {
			article, _ := models.GetArticle(id)
			data["article"] = article
			code = e.SUCCESS
		} else {
			code = e.ERROR_GET_ARTICLE_FAIL
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": data,
	})
}

func AddArticle(c *gin.Context) {
	tag_id := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.Query("state")).MustInt()

	valid := validation.Validation{}
	valid.Min(tag_id, 1, "tag_id").Message("標籤ID需大於0")
	valid.Required(title, "title").Message("Title is required field")
	valid.Required(desc, "desc").Message("Description is required field")
	valid.Required(content, "content").Message("Content is required field")
	valid.Required(createdBy, "created_by").Message("Created_by is required field")
	valid.Range(state, 0, 1, "state").Message("State is only accepted 0 or 1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagById(tag_id) {
			// No error be found and tag is exist
			var data = make(map[string]interface{})
			data["tag_id"] = tag_id
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state_code"] = state
			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_ADD_ARTICLE_FAIL
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

// EditArticle Required field `tag_id`, `title`, `desc`, `content`, `created_by`, `state`
func EditArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	tag_id := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("文章ID不得小於1")
	valid.Min(tag_id, 1, "tag_id").Message("標籤ID需大於0")
	valid.Required(title, "title").Message("Title is required field")
	valid.Required(desc, "desc").Message("Desc is required field")
	valid.Required(content, "content").Message("Content is required field")
	valid.Required(createdBy, "created_by").Message("created_by is required field")
	valid.Range(state, 0, 1, "state").Message("State value only accepted 0 or 1")

	var data = make(map[string]interface{})
	var maps = make(map[string]interface{})
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		// 驗證無錯
		success, _ := models.ExistArticleById(id)
		if success {
			data["tag_id"] = tag_id
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state_code"] = state
			err := models.EditArticle(id, data)
			if err != nil {
				code = e.ERROR_EDIT_ARTICLE_FAIL
			} else {
				code = e.SUCCESS
				maps["data"] = data
			}
		} else {
			// 找不到對應ID的文章
			code = e.ERROR_EDIT_ARTICLE_FAIL
		}
	} else {
		for _, err := range valid.Errors {
			logging.Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": maps,
	})
}

func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID不得小於1")
	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		// 驗證無錯
		success, _ := models.ExistArticleById(id)
		if success {
			// 找到對應文章
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_DELETE_ARTICLE_FAIL
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
