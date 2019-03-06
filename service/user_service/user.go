package user_service

import (
	"encoding/json"
	"vincent-gin-go/models"
	"vincent-gin-go/pkg/gredis"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/service/cache_service"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
	State    int

	PageNum  int
	PageSize int
}

func (u *User) ExistByID() (bool, error) {
	return models.ExistUserByID(u.ID)
}

func (u *User) ExistByEmail() bool {
	return models.ExistUserByEmail(u.Email)
}

func (u *User) ExistByName() bool {
	return models.ExistUserByName(u.Name)
}

func (u *User) Get() (*models.User, error) {
	var user *models.User
	cache := cache_service.User{ID: u.ID}
	key := cache.GetUserKey()
	if gredis.Exists(key) {
		// 若存在
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &user)
			return user, nil
		}
	}
	// 若不存在
	existUser, err := models.GetUser(u.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, existUser, 3600)
	return existUser, err
}

// Count 獲取User數量
func (u *User) Count() (int, error) {
	return models.GetUserCount(u.getMap())
}

// GetAllUser 取得多個用戶
func (u *User) GetAllUser() ([]*models.User, error) {
	var (
		users, cacheUsers []*models.User
	)
	cache := cache_service.User{
		ID:       u.ID,
		Name:     u.Name,
		PageNum:  u.PageNum,
		PageSize: u.PageSize,
	}
	key := cache.GetUsersKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheUsers)
			return cacheUsers, nil
		}
	}
	// 緩存不存在
	users, err := models.GetAllUser(u.getMap(), u.PageSize, u.PageNum)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, users, 3600)
	return users, nil
}

// Edit 編輯用戶
func (u *User) Edit() error {
	err := models.EditUser(u.ID, map[string]interface{}{
		"name":     u.Name,
		"email":    u.Email,
		"password": u.Password,
		"state":    u.State,
	})
	return err
}

// Add 新增User
func (u *User) Add() error {
	err := models.AddUser(u.getMap())
	return err
}

// Delete To delete this `user`
func (u *User) Delete() error {
	err := models.DeleteUser(u.ID)
	return err
}

func (u *User) getMap() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["name"] = u.Name
	maps["email"] = u.Email
	maps["password"] = u.Password
	return maps
}
