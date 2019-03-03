package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

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
	User          User   `json: "user"`
	UserID        int    `json: "user_id"; sql:DEFAULT: 0`
}

func ExistArticleByName(name string) bool {
	var article Article
	database.Where("name = ?", name).First(&article)
	return article.ID > 0
}

func ExistArticleById(id int) (bool, error) {
	var article Article
	err := database.Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	return article.ID > 0, err
}

func GetArticleTotalCount(data interface{}) (int, error) {
	var count int
	err := database.Model(Article{}).Where(data).Count(&count).Error
	return count, err
}

// GetAllArticles `Get all article in the database`
func GetAllArticles(pageNum int, pageSize int, maps interface{}) ([]*Article, error) {
	var articles = []*Article{}
	err := database.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&articles).Error
	return articles, err
}

func GetArticle(id int) (*Article, error) {
	// 	能够达到关联，首先是gorm本身做了大量的约定俗成
	//  Article有一个结构体成员是TagID，就是外键。gorm会通过类名+ID的方式去找到这两个类之间的关联关系
	//  Article有一个结构体成员是Tag，就是我们嵌套在Article里的Tag结构体，我们可以通过Related进行关联查询
	var article Article
	tag := Tag{}
	user := User{}
	err := database.Where("id = ? AND deleted_on = ?", id, 0).First(&article).Error
	err = database.Model(&article).Related(&tag).Error
	err = database.Model(&article).Related(&user).Error
	article.Tag = tag
	article.User = user
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &article, nil
}

func AddArticle(data map[string]interface{}) error {
	err := database.Create(&Article{
		TagID:         data["tag_id"].(int),
		Title:         data["title"].(string),
		Desc:          data["desc"].(string),
		Content:       data["content"].(string),
		CreatedBy:     data["created_by"].(string),
		StateCode:     data["state_code"].(int),
		CoverImageUrl: data["cover_image_url"].(string),
		UserID:        data["user_id"].(int),
	}).Error
	return err
}

func EditArticle(id int, data interface{}) error {
	if err := database.Model(&Article{}).Where("id = ?", id).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteArticle(id int) bool {
	database.Where("id = ?", id).Delete(Article{})
	return true
}
