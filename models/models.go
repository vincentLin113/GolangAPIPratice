package models

import (
	"fmt"
	"log"
	"vincent-gin-go/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
}

// 預設取得Auth的帳號密碼組
var auth = &Auth{
	Username: "vincent_test",
	Password: "001122",
}

func init() {
	var (
		err                              error
		dbType, dbName, user, host, port string
	)
	setting.Setup()
	dbType = setting.DatabaseSetting.Type
	dbName = setting.DatabaseSetting.Name
	user = setting.DatabaseSetting.User
	host = setting.DatabaseSetting.Host
	port = setting.DatabaseSetting.Port
	db, err = gorm.Open(dbType, fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", host, user, dbName, port))
	if err != nil {
		log.Fatalf("Open db error: %v", err)
	}
	db.AutoMigrate(&Tag{}, &Article{}, &Auth{})
	fmt.Println("Create `Tag` and `Article` and `Auth`")

	if !checkDefaultAuthInfo() {
		db.Create(&auth)
		fmt.Println("Create default auth information☺️")
	}
}

func checkDefaultAuthInfo() bool {
	var _au Auth
	var maps = make(map[string]interface{})
	maps["username"] = &auth.Username
	maps["password"] = &auth.Password
	db.Model(&Auth{}).Where(maps).Find(&_au)
	return _au.ID > 0
}
