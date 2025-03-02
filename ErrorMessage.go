package common

var CodeMapMessage = map[int]string{
	Success: "操作成功",

	//公共错误
	CommonSignError:             "签名错误",
	CommonParamError:            "请求参数不合法",
	CommonPermissionDenied:      "账号没有访问接口的权限",
	CommonTokenError:            "token错误",
	CommonPictureError:          "缺少上传图片",
	CommonUploadError:           "上传失败",
	CommonPushError:             "推送错误",
	CommonSystemError:           "系统错误",
	CommonTokenTimeOut:          "登录超时,请重新登录",
	CommonCaptchaCreateError:    "验证码生成失败",
	CommonTokenCreateError:      "token生成失败",
	CommonNotCarriedToken:       "未携带token",
	CommonStatusParamError:      "状态参数错误",
	CommonTimeDateError:         "时间日期格式错误",
	CommonUploadTooBig:          "图片超出大小限制",
	CommonAnalysisPictureFail:   "解析图片失败",
	CommonImageWrongStyle:       "无效的图片格式",
	CommonSortParamMust:         "排序参数为必填",
	CommonTwoFactorNoPassed:     "二次验证未检验",
	CommonRebuildTimeError:      "无法重建70天之前的报表",
	CommonRebuildTodayError:     "无法重建当天的报表",
	CommonInitNotCompletedError: "系统资源初始化尚未完成",
	CommonDataNotExist:          "数据不存在",
	CommonDbError:               "系统错误",
	CommonDbInsertError:         "数据保存失败",
	CommonDbUpdateError:         "数据更新失败",
}
