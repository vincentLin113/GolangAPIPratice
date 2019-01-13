package util

import (
	"time"
	"vincent-gin-go/pkg/setting"

	jwt "github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json: "username"`
	Password string `json: "password"`
	jwt.StandardClaims
}

// GenerateToken 產生憑證
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(setting.AppSetting.JwtExpireHour) * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "vincent-gin-blog",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// ParseToken 解析憑證
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// func (m MapClaims) Valid() 验证基于时间的声明exp, iat, nbf，注意如果没有任何声明在令牌中，仍然会被认为是有效的。并且对于时区偏差没有计算方法
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
