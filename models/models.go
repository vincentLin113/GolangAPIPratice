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

func init() {
	var (
		err                              error
		dbType, dbName, user, host, port string
	)
	setting.Setup()
	dbType = setting.DatabaseSetting.Type
	dbName = setting.DatabaseSetting.Name
	user = setting.DatabaseSetting.User
	// password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	port = setting.DatabaseSetting.Port
	db, err = gorm.Open(dbType, fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", host, user, dbName, port))
	if err != nil {
		log.Fatalf("Open db error: %v", err)
	}
}
