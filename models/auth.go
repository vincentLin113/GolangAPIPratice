package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	db.Table("auths").First(&auth)
	var auths []Auth
	db.Find(&auths)
	return auth.ID != 0
}
