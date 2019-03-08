package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"vincent-gin-go/pkg/app"
	"vincent-gin-go/pkg/e"
	"vincent-gin-go/pkg/logging"
	"vincent-gin-go/pkg/upload"

	"github.com/astaxie/beego/validation"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMessage(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()
		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMessage(code),
		"data": data,
	})
}

func UploadImageByURL(c *gin.Context) {
	appG := app.Gin{c}
	imageURLText := c.Query("imageUrl")
	valid := validation.Validation{}
	valid.Required(imageURLText, "imageURLText").Message("ImageURLText is required field")
	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		return
	}
	response, err := http.Get(imageURLText)
	if err != nil {
		fmt.Printf("GET IMAGE ERROR: %v", err)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}
	image := response.Body
	imageName := upload.GetImageName(imageURLText)
	imageFullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := imageFullPath + imageName
	file, err := os.Create(src)
	err := io.Copy(file, response.Body)
	if err != nil {
		panic(err)
	}
	// if !upload.CheckImageExt(imageName) || upload.CheckImageSize(image) {

	// }
}
