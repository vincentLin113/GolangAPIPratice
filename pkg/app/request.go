package app

import (
	"vincent-gin-go/pkg/logging"

	"github.com/astaxie/beego/validation"
)

// MarkErrors: 處理Validation的錯誤
func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
	return
}
