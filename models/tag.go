package models

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) (tag []Tag) {
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tag).Error
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	return
}

func GetTagTotal(maps interface{}) (count int) {
	db.Model(&Tag{}).Where(maps).Count(&count)
	return
}

func ExistTagByName(name string) bool {
	var tag Tag
	db.Select("id").Where("name = ?", name).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func ExistTagById(id int) bool {
	var tag Tag
	db.Where("id = ?", id).First(&tag)
	if tag.ID > 0 {
		return true
	}
	return false
}

func AddTag(name string, state int, createdBy string) bool {
	db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	})
	return true
}

func EditTag(id int, data interface{}) bool {
	var tag Tag
	err := db.Where("id = ?", id).First(&tag).Updates(data).Error
	if err != nil {
		log.Fatalf("EditTag error: %v", err)
		return false
	}
	return true
}

func DeleteTag(id int) (bool, error) {
	err := db.Where("id = ?", id).Delete(&Tag{}).Error
	return true, err
}

// gorm所支持的回调方法：

// 创建：BeforeSave、BeforeCreate、AfterCreate、AfterSave
// 更新：BeforeSave、BeforeUpdate、AfterUpdate、AfterSave
// 删除：BeforeDelete、AfterDelete
// 查询：AfterFind
// Gorm回調方法
func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

// Gorm回調方法
func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}
