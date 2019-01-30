package models

import (
	"log"
)

type Tag struct {
	Model
	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	err := db.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	return tags, err
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := db.Model(&Tag{}).Where(maps).Count(&count).Error
	return count, err
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := db.Select("id").Where("name = ?", name).First(&tag).Error
	if tag.ID > 0 {
		return true, err
	}
	return false, err
}

func ExistTagById(id int) (bool, error) {
	var tag Tag
	err := db.Where("id = ?", id).First(&tag).Error
	if tag.ID > 0 {
		return true, err
	}
	return false, err
}

func AddTag(name string, state int, createdBy string) error {
	err := db.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error
	return err
}

func EditTag(id int, data interface{}) error {
	var tag Tag
	err := db.Where("id = ?", id).First(&tag).Updates(data).Error
	if err != nil {
		log.Fatalf("EditTag error: %v", err)
	}
	return err
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
// func (tag *Tag) BeforeCreate(scope *gorm.Scope) error {
// 	scope.SetColumn("CreatedOn", time.Now().Unix())
// 	return nil
// }

// // Gorm回調方法
// func (tag *Tag) BeforeUpdate(scope *gorm.Scope) error {
// 	scope.SetColumn("ModifiedOn", time.Now().Unix())
// 	return nil
// }
