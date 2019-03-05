package article_service

import (
	"encoding/json"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/gredis"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/service/cache_service"
	"vincent-gin-go/service/user_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	ModifiedBy    string
	UserID        int

	PageNum  int
	PageSize int
}

func (a *Article) ExistByID() (bool, error) {
	return models.ExistArticleById(a.ID)
}

// Get 獲取Article(From Redis or Database)
func (a *Article) Get() (*models.Article, error) {
	var cacheArticle *models.Article
	cache := cache_service.Article{ID: a.ID}
	// 獲取獨特辨識碼
	key := cache.GetArticleKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheArticle)
			return cacheArticle, nil
		}
	}
	// 走到這邊 代表Redis無緩存
	article, err := models.GetArticle(a.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, article, 3600)
	return article, nil
}

// Count 調用Model中的方法
func (a *Article) Count() (int, error) {
	return models.GetArticleTotalCount(a.getMaps())
}

// GetAll 獲取全部Article
func (a *Article) GetAll() ([]*models.Article, error) {
	var (
		articles, cacheArticles []*models.Article
	)
	cache := cache_service.Article{
		TagID:    a.TagID,
		State:    a.State,
		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetArticlesKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			// 取得Redis所緩存的Articles
			json.Unmarshal(data, &cacheArticles)
			return cacheArticles, nil
		}
	}
	articles, err := models.GetAllArticles(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, articles, 3600)
	return articles, nil
}

func (a *Article) Add() (err error, errorCode int, errorMsg string) {
	userService := user_service.User{ID: a.UserID}
	_, err = userService.ExistByID()
	if err != nil {
		errorCode = e.ERROR_GET_USER_FAIL
		errorMsg = e.GetMessage(e.ERROR_GET_USER_FAIL)
		return
	}
	_, err, code, msg := userService.Get()
	if err != nil {
		errorCode = e.ERROR_GET_USER_FAIL
		errorMsg = e.GetMessage(e.ERROR_GET_USER_FAIL)
		return
	}
	if code != 0 && msg != "" {
		errorCode = code
		errorMsg = msg
		return
	}
	// if user.DeletedOn > 0 {
	// 	err = e.UserBeDeleted()
	// 	errorCode = e.ERROR_GET_USER_DELETED_FAIL
	// 	errorMsg = e.GetMessage(e.ERROR_GET_USER_DELETED_FAIL)
	// 	return
	// }
	// if user.State == 0 {
	// 	err = e.UserBeStoped()
	// 	errorCode = e.ERROR_GET_USER_BAN_FAIL
	// 	errorMsg = e.GetMessage(e.ERROR_GET_USER_BAN_FAIL)
	// 	return
	// }
	err = models.AddArticle(a.getMaps())
	errorCode = 0
	errorMsg = ""
	return
}

func (a *Article) Edit() error {
	var err error
	err = models.EditArticle(a.ID, map[string]interface{}{
		"tag_id":          a.TagID,
		"title":           a.Title,
		"desc":            a.Desc,
		"content":         a.Content,
		"cover_image_url": a.CoverImageUrl,
		"state":           a.State,
		"modified_by":     a.ModifiedBy,
		"user_id":         a.UserID,
	})
	return err
}

// getMaps 從Model本身產出Dictionary(map[string]interface{})
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	if a.State != -1 {
		maps["state_code"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	if a.UserID != -1 {
		maps["user_id"] = a.UserID
		// userService := user_service.User{ID: a.UserID}
		// user, err, errCode, errMsg := userService.Get()
		// if err == nil {
		// 	maps["created_by"] = user.Name
		// }
	}
	maps["content"] = a.Content
	maps["desc"] = a.Desc
	maps["title"] = a.Title
	maps["cover_image_url"] = a.CoverImageUrl
	return maps
}
