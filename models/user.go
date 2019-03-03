package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	State    int    `json:state;sql:"DEFAULT:1"`
}

// ExistUserByName 檢查此名稱是否有人
func ExistUserByName(name string) bool {
	var user User
	database.Where("name = ?", name).First(&user)
	return user.ID > 0
}

// ExistUserByEmail 檢查此信箱是否有人
func ExistUserByEmail(email string) bool {
	var user User
	database.Where("email = ?", email).First(&user)
	return user.ID > 0
}

// ExistUserByID 檢查此ID是否有人
func ExistUserByID(id int) (bool, error) {
	var user User
	err := database.Where("id = ?", id).First(&user).Error
	return user.ID > 0, err
}

// AddUser 新增用戶
func AddUser(data map[string]interface{}) error {
	err := database.Create(&User{
		Name:     data["name"].(string),
		Email:    data["email"].(string),
		Password: data["password"].(string),
		State:    data["state"].(int),
	}).Error
	return err
}

// GetUser 獲取用戶
func GetUser(id int) (*User, error) {
	var user User
	err := database.Where("id = ?", id).First(&user).Error
	// GORM中找不到也算一種錯誤
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func GetAllUser(maps interface{}, pageSize int, pageNum int) ([]*User, error) {
	users := []*User{}
	err := database.Preload("Article").Where(maps).Offset(pageNum).Limit(pageSize).Find(&users).Error
	return users, err
}

// EditUser 編輯用戶
func EditUser(id int, data map[string]interface{}) error {
	if err := database.Model(&User{}).Where("id = ?", id).Update(data).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

// DeleteUser 刪除用戶
func DeleteUser(id int) error {
	err := database.Where("id = ?", id).Delete(User{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return nil
}

// GetUserCount 獲取用戶數量
func GetUserCount(data interface{}) (int, error) {
	var count int
	err := database.Model(User{}).Where(data).Count(&count).Error
	return count, err
}
