package e

var MsgFlags = map[int]string{
	SUCCESS:                               "ok",
	ERROR:                                 "fail",
	INVALID_PARAMS:                        "请求参数错误",
	ERROR_EXIST_TAG:                       "已存在该标签名称",
	ERROR_EXIST_TAG_FAIL:                  "获取已存在标签失败",
	ERROR_NOT_EXIST_TAG:                   "该标签不存在",
	ERROR_GET_TAGS_FAIL:                   "获取所有标签失败",
	ERROR_COUNT_TAG_FAIL:                  "统计标签失败",
	ERROR_ADD_TAG_FAIL:                    "新增标签失败",
	ERROR_EDIT_TAG_FAIL:                   "修改标签失败",
	ERROR_DELETE_TAG_FAIL:                 "删除标签失败",
	ERROR_EXPORT_TAG_FAIL:                 "导出标签失败",
	ERROR_IMPORT_TAG_FAIL:                 "导入标签失败",
	ERROR_NOT_EXIST_ARTICLE:               "该文章不存在",
	ERROR_ADD_ARTICLE_FAIL:                "新增文章失败",
	ERROR_DELETE_ARTICLE_FAIL:             "删除文章失败",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:        "检查文章是否存在失败",
	ERROR_EDIT_ARTICLE_FAIL:               "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:              "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:               "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:                "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:         "生成文章海报失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:           "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:        "Token已超时",
	ERROR_AUTH_TOKEN:                      "Token生成失败",
	ERROR_AUTH:                            "Token错误",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:          "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:         "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT:       "校验图片错误，图片格式或大小有问题",
	ERROR_AUTH_NEED_TOKEN:                 "參數需加入Token",
	ERROR_GET_USER_FAIL:                   "GET_USER_FAIL",
	ERROR_ADD_USER_FAIL:                   "ADD_USER_FAIL",
	ERROR_EDIT_USER_FAIL:                  "EDIT_USER_FAIL",
	ERROR_DELETE_USER_FAIL:                "DELETE_USER_FAIL",
	ERROR_ADD_USER_DUPLICATED_EMAIL_ERROR: "USER_EMAIL_DUPLICATED",
	ERROR_ADD_USER_DUPLICATED_NAME_ERROR:  "USER_NAME_DUPLICATED",
	ERROR_GET_USER_BAN_FAIL:               "USER_BAN_FAIL",
	ERROR_GET_USER_DELETED_FAIL:           "USER_BE_DELETED",
}

func GetMessage(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
