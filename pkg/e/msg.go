package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "請求參數錯誤",
	ERROR_EXIST_TAG:                 "已存在該標籤名稱",
	ERROR_NOT_EXIST_TAG:             "該標籤不存在",
	ERROR_NOT_EXIST_ARTICLE:         "該文章不存在",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token失敗",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超時",
	ERROR_AUTH_TOKEN:                "Token生成失敗",
	ERROR_AUTH:                      "Token錯誤",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存圖片失敗",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "檢查圖片失敗",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校驗圖片錯誤，圖片格式或大小有問題",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "检查文章是否存在失败",
	ERROR_EDIT_ARTICLE_FAIL:         "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:        "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:         "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:          "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海报失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
