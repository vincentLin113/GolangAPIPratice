package article_service

import (
	"encoding/json"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/gredis"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/service/cache_service"
)

type Article struct {
	ID            int
	TagID         int
	Title         string
	Desc          string
	Content       string
	CoverImageUrl string
	State         int
	CreatedBy     string
	ModifiedBy    string

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

func (a *Article) Add() error {
	return models.AddArticle(a.getMaps())
}

// getMaps 從Model本身產出Dictionary(map[string]interface{})
func (a *Article) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state_code"] = a.State
	}
	if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}
	return maps
}
