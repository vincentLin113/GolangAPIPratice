package v1

import (
	"net/http"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/app"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/setting"
	"vincent-gin-go/service/article_service"
	"vincent-gin-go/service/tag_service"
	"vincent-gin-go/util"

	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func GetArticles(c *gin.Context) {
	appG := app.Gin{c}
	validation := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		validation.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	tagId := -1
	if arg := c.PostForm("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		validation.Min(tagId, 1, "tag_id").Message("Tag_ID不得小於1")
	}

	// 檢查驗證
	if validation.HasErrors() {
		app.MarkErrors(validation.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total
	appG.Response(http.StatusOK, e.SUCCESS, data)
}

func GetArticle(c *gin.Context) {
	appG := app.Gin{c}
	id := com.StrTo(c.Param("id")).MustInt()
	validation := validation.Validation{}
	validation.Min(id, 1, "id").Message("ID不能小於1")

	if validation.HasErrors() {
		// 若有錯誤, 使用共有方法處理
		app.MarkErrors(validation.Errors)
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exist, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	if !exist {
		// 若不存在, 則http為正常
		appG.Response(http.StatusOK, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}
	article, err := articleService.Get()
	appG.Response(http.StatusOK, e.SUCCESS, article)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state_code" valid:"Range(0, 1)"`
}

func AddArticle(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddArticleForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	_, err := tagService.ExistById()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}
	if err := articleService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state_code" valid:"Range(0, 1)"`
}

// EditArticle Required field `tag_id`, `title`, `desc`, `content`, `created_by`, `state`
func EditArticle(c *gin.Context) {
	appG := app.Gin{c}
	form := EditArticleForm{
		ID: com.StrTo(c.Param("id")).MustInt(),
	}

	httpCode, errCode := app.BindAndValid(c, form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		ID:            form.ID,
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		ModifiedBy:    form.ModifiedBy,
	}
	exists, err := articleService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{
		ID: form.TagID,
	}
	tagExist, err := tagService.ExistById()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}
	if !tagExist {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	err = articleService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
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
