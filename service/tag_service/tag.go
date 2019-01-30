package tag_service

import (
	"encoding/json"
	"strconv"
	"strings"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/gredis"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/service/cache_service"
)

type Tag struct {
	ID         int
	Name       string
	CreatedBy  string
	ModifiedBy string
	State      int
	PageNum    int
	PageSize   int
}

func (t *Tag) GetTagsKey() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.Name != "" {
		keys = append(keys, t.Name)
	}
	if t.State >= 0 {
		keys = append(keys, strconv.Itoa(t.State))
	}
	if t.PageNum > 0 {
		keys = append(keys, strconv.Itoa(t.PageNum))
	}
	if t.PageSize > 0 {
		keys = append(keys, strconv.Itoa(t.PageSize))
	}

	return strings.Join(keys, "_")
}

// 根據Tag的ID檢查是否存在
func (t *Tag) ExistById() (bool, error) {
	return models.ExistTagById(t.ID)
}

func (t *Tag) ExistByName() (bool, error) {
	return models.ExistTagByName(t.Name)
}

func (t *Tag) Add() error {
	return models.AddTag(t.Name, t.State, t.CreatedBy)
}

func (t *Tag) Edit() error {
	data := make(map[string]interface{})
	data["modified_by"] = t.ModifiedBy
	data["name"] = t.Name
	if t.State >= 0 {
		data["state"] = t.State
	}
	return models.EditTag(t.ID, data)
}

func (t *Tag) Delete() (bool, error) {
	return models.DeleteTag(t.ID)
}

func (t *Tag) Count() (int, error) {
	return models.GetTagTotal(t.getMaps())
}

func (t *Tag) GetAll() ([]models.Tag, error) {
	var (
		tags, cacheTags []models.Tag
	)

	cache := cache_service.Tag{
		State:    t.State,
		Name:     t.Name,
		PageNum:  t.PageNum,
		PageSize: t.PageSize,
	}
	key := cache.GetTagKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheTags)
			return cacheTags, nil
		}
	}

	tags, err := models.GetTags(t.PageNum, t.PageSize, t.getMaps())
	if err != nil {
		return nil, err
	}
	gredis.Set(key, tags, 3600)
	return tags, nil
}

// getMaps 獲得由自身的值所產生的Dictionary
func (t *Tag) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if t.Name != "" {
		maps["name"] = t.Name
	}
	if t.State >= 0 {
		maps["state"] = t.State
	}
	return maps
}
