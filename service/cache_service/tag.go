package cache_service

import (
	"strconv"
	"strings"
	"vincent-gin-go/pkg/e"
)

type Tag struct {
	ID       int
	State    int
	Name     string
	PageNum  int
	PageSize int
}

// GetTagKey 獲取Cache of tag的獨特碼
func (t *Tag) GetTagKey() string {
	return e.CACHE_TAG + "_" + strconv.Itoa(t.ID)
}

func (t *Tag) GetTagsTag() string {
	keys := []string{
		e.CACHE_TAG,
		"LIST",
	}

	if t.ID > 0 {
		keys = append(keys, strconv.Itoa(t.ID))
	}
	if t.State > 0 {
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
