package util

import (
	"fmt"
	"vincent-gin-go/pkg/setting"

	"gopkg.in/gomail.v2"
)

func SendActivationEmail(email string, token string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "keepexcelsior@gmail")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "ACTIVATE YOUR ACCOUNT!")
	address := "http://localhost:8001/activate" + "?email=" + email + "&token=" + token
	m.SetBody("text/html", "<html><head><title>Activate your account\n</title></head><body><table width=\"100%\" cellspacing=\"0\" cellpadding=\"0\"><tr><td><table cellspacing=\"0\" cellpadding=\"0\"><tr><td style=\"border-radius: 2px;\" bgcolor=\"#ED2939\"><a href="+address+" target=\"_blank\" style=\"padding: 8px 12px; border: 1px solid #ED2939;border-radius: 2px;font-family: Helvetica, Arial, sans-serif;font-size: 14px; color: #ffffff;text-decoration: none;font-weight:bold;display: inline-block;\">Click</a></td></tr></table></td></tr></table></body></html>")
	fmt.Printf("##SEND MAIL: %v", email)
	d := gomail.NewPlainDialer("smtp.gmail.com", 587, setting.EmailSetting.Email, setting.EmailSetting.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
