package models

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) bool {
	var auth Auth
	database.Table("users").Where("email = ? AND password = ?", username, password).Find(&auth)
	return auth.ID != 0
}
