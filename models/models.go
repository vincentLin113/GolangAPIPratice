package models

import (
	"fmt"
	"log"
	"os"
	"time"
	"vincent-gin-go/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

// 預設取得Auth的帳號密碼組
var auth = Auth{
	Username: "vincent_test",
	Password: "001122",
}

func Setup() {
	var err error
	if setting.IsLocalTest() {
		dbType := setting.DatabaseSetting.Type
		dbName := setting.DatabaseSetting.Name
		user := setting.DatabaseSetting.User
		// host := setting.DatabaseSetting.Host
		port := setting.DatabaseSetting.Port
		db, err = gorm.Open(dbType, fmt.Sprintf("host=127.0.0.1 user=%s dbname=%s port=%s sslmode=disable", user, dbName, port))
	} else {
		db, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}

	if err != nil {
		log.Fatalf("Open db error: %v", err)
	}
	db.AutoMigrate(&Tag{}, &Article{}, &Auth{})
	fmt.Println("Create `Tag` and `Article` and `Auth`")
	// 此為單數Model的開關
	// db.SingularTable(true)
	db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	if !checkDefaultAuthInfo() {
		db.Create(&auth)
		fmt.Println("Create default auth information☺️")
	}
}

func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			// 製造時間 若為空白
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			// 修改時間若為空白
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

func checkDefaultAuthInfo() bool {
	var _au Auth
	db.Table("auths").Where("password = ?", &auth.Password).Find(&_au)
	return _au.ID > 0
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		// 檢查是否手動指定了 delete_option
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}
		// 獲取我們約定的刪除字段, 若存在則UPDATE軟刪除, 不存在則DELETE硬刪除
		deletedOnField, hasDeletedField := scope.FieldByName("DeletedOn")
		if !scope.Search.Unscoped && hasDeletedField {
			// example: "UPDATE \"articles\" SET \"deleted_on\"=$1  WHERE (id = $2)"
			sqlText := fmt.Sprintf("UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedOnField.DBName),
				scope.AddToVars(time.Now().Unix()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)
			scope.Raw(sqlText).Exec()
		} else {
			sqlText := fmt.Sprintf("DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)
			scope.Raw(sqlText).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
