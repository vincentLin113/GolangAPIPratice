package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	database.Where(&Auth{Username: username, Password: password}).Find(&auth)
	return auth.ID != 0
}
