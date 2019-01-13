package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

type Article struct {
	Model
	// Id      int    `json: "id" gorm: "index"`
	Tag        Tag    `json: "tag"`
	TagID      int    `json: "tag_id"`
	Title      string `json: "title"`
	Desc       string `json: "desc"`
	Content    string `json: "content"`
	CreatedBy  string `json: "created_by"`
	ModifiedBy string `json: "modified_by"`
	DeletedOn  int    `json: "deleted_on"`
	StateCode  int    `json: "stateCode"`
}

func (article *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedOn", time.Now().Unix())
	return nil
}

func (article *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("ModifiedOn", time.Now().Unix())
	return nil
}

func ExistArticleByName(name string) bool {
	var article Article
	db.Where("name = ?", name).First(&article)
	return article.ID > 0
}

func ExistArticleById(id int) bool {
	var article Article
	db.Where("id = ?", id).First(&article)
	return article.ID > 0
}

func GetArticleTotalCount(data interface{}) (count int) {
	db.Model(Article{}).Where(data).Count(&count)
	return
}

// GetAllArticles `Get all article in the database`
func GetAllArticles(pageNum int, pageSize int, maps interface{}) (articles []Article) {
	db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles)
	return
}

func GetArticle(id int) (article Article) {
	// 	能够达到关联，首先是gorm本身做了大量的约定俗成
	//  Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
	//  Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询
	db.Where("id = ?", id).First(&article)
	db.Model(&article).Related(&article.Tag)
	return
}

func AddArticle(data map[string]interface{}) bool {
	db.Create(&Article{
		TagID:     data["tag_id"].(int),
		Title:     data["title"].(string),
		Desc:      data["desc"].(string),
		Content:   data["content"].(string),
		CreatedBy: data["created_by"].(string),
		StateCode: data["state_code"].(int),
	})
	return true
}

func EditArticle(id int, data interface{}) (bool, Article) {
	db.Model(&Article{}).Where("id = ?", id).Updates(data)
	article := GetArticle(id)
	return true, article
}

func DeleteArticle(id int) bool {
	db.Where("id = ?", id).Delete(Article{})
	return true
}
