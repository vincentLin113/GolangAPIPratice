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

// isValid Check this tag is valid or not.
func (t *Tag) isValid() bool {
	return t.State == 1
}

func GetTags(pageNum int, pageSize int, maps interface{}) ([]Tag, error) {
	var tags []Tag
	err := database.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	return tags, err
}

func GetTagTotal(maps interface{}) (int, error) {
	var count int
	err := database.Model(&Tag{}).Where(maps).Count(&count).Error
	return count, err
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := database.Select("id").Where("name = ?", name).First(&tag).Error
	if tag.ID > 0 {
		return true, err
	}
	return false, err
}

func ExistTagById(id int) (bool, error) {
	var tag Tag
	err := database.Where("id = ?", id).First(&tag).Error
	if tag.ID > 0 {
		return true, err
	}
	return false, err
}

func AddTag(name string, state int, createdBy string) error {
	err := database.Create(&Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}).Error
	return err
}

func EditTag(id int, data interface{}) error {
	var tag Tag
	err := database.Where("id = ?", id).First(&tag).Updates(data).Error
	if err != nil {
		log.Fatalf("EditTag error: %v", err)
	}
	return err
}

func DeleteTag(id int) (bool, error) {
	err := database.Where("id = ?", id).Delete(&Tag{}).Error
	return true, err
}
