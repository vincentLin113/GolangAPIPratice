package models

import (
	"fmt"
	"log"
	"os"
	"time"
	"vincent-gin-go/pkg/setting"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	gormigrate "gopkg.in/gormigrate.v1"
)

var database *gorm.DB

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
		database, err = gorm.Open(dbType, fmt.Sprintf("host=127.0.0.1 user=%s dbname=%s port=%s sslmode=disable", user, dbName, port))
	} else {
		database, err = gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	}

	if err != nil {
		log.Fatalf("Open db error: %v", err)
	}
	database.AutoMigrate(&Tag{}, &Article{}, &Auth{})
	fmt.Println("Create `Tag` and `Article` and `Auth`")
	// 此為單數Model的開關
	// db.SingularTable(true)
	database.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	database.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	database.Callback().Delete().Replace("gorm:delete", deleteCallback)
	database.DB().SetMaxIdleConns(10)
	database.DB().SetMaxOpenConns(100)
	if !checkDefaultAuthInfo() {
		database.Create(&auth)
		fmt.Println("Create default auth information☺️")
	}
	databaseMigrate()
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
	database.Table("auths").Where("password = ?", &auth.Password).Find(&_au)
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

func databaseMigrate() {
	m := gormigrate.New(database, gormigrate.DefaultOptions, []*gormigrate.Migration{
		// create persons table
		{
			ID: "201902261439",
			Migrate: func(tx *gorm.DB) error {
				return tx.AutoMigrate(&User{}).Error
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.DropTable("Users").Error
			},
		},
		{
			ID: "201902261638",
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(User{}).Error
				return err
			},
		},
		{
			ID: "201902261707",
			Migrate: func(tx *gorm.DB) error {
				err := tx.Table("users").Where("state is NULL").UpdateColumn("state", "0").Error
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				fmt.Println("\n#### UPDATE USER.STATE column fail ###")
				return nil
			},
		},
		{
			ID: "2019002270944",
			Migrate: func(tx *gorm.DB) error {
				err := tx.Table("users").Where("password = ?", "").UpdateColumn("password", "123456abcd").Error
				if err == nil {
					fmt.Println("\n### UPDATE USER.PASSWORD successfully###")
				}
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				fmt.Println("\n ### UPDATE USER.PASSWORD column fail ###")
				return nil
			},
		},
		{
			ID: "201902271027",
			Migrate: func(tx *gorm.DB) error {
				type Article struct {
					Model
					Tag           Tag    `json: "tag"`
					TagID         int    `json: "tag_id"`
					Title         string `json: "title"`
					Desc          string `json: "desc"`
					Content       string `json: "content"`
					CreatedBy     string `json: "created_by"`
					ModifiedBy    string `json: "modified_by"`
					StateCode     int    `json: "stateCode"`
					CoverImageUrl string `json:"cover_image_url"`
					User_ID       int    `json: "user_id"; sql:DEFAULT: 0`
				}
				err := tx.AutoMigrate(&Article{}).Error
				if err == nil {
					fmt.Println("\n### ADD NEW COLUMN `User_ID` to `User`###")
					err = tx.Table("articles").Where("user_id is NULL").UpdateColumn("user_id", "0").Error
					if err == nil {
						fmt.Println("\n####UPDATE ARTICLE.USER_ID column successfully###")
					}
				}
				tx.Table("users").DropColumn("user_id")
				return err
			},
			Rollback: func(tx *gorm.DB) error {
				fmt.Println("\n###ADD COLUMN `User_ID` is fail###")
				return nil
			},
		},
	})

	if err := m.Migrate(); err != nil {
		log.Fatalln("Could not migrate: ", err)
	}
	log.Println("### Migration did run successfully ###")
}
