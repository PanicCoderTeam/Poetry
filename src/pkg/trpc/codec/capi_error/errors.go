package capi_error

const INTERNAL_ERROR_CODE ErrorCode = "InternalError"
const INVAILD_PARAM_CODE ErrorCode = "InvaildParamError"
const RESOURCE_NOT_FOUND_CODE ErrorCode = "ResourceNotFoundError"
const REQUEST_NOT_AUTH_CODE ErrorCode = "RequestNotAuthError"
const RESOURCE_OUT_OF_LIMIT ErrorCode = "ResourceOutOfLimitError"
const USER_ALREADY_EXIST_CODE ErrorCode = "UserAlreadyExistError"
const USER_ALREADY_IN_GAME_CODE ErrorCode = "UserAlreadyInGameError"

var errorInfoMap = map[ErrorCode]*CapiErrorInfo{
	INTERNAL_ERROR_CODE: {
		ErrorCode: INTERNAL_ERROR_CODE,
		ErrorMsg: []string{
			"内部错误",
			"Internal Error",
		},
	},
	INVAILD_PARAM_CODE: {
		ErrorCode: INVAILD_PARAM_CODE,
		ErrorMsg: []string{
			"参数错误",
			"Invaild Param Error",
		},
	},
	RESOURCE_NOT_FOUND_CODE: {
		ErrorCode: RESOURCE_NOT_FOUND_CODE,
		ErrorMsg: []string{
			"资源不存在",
			"Resource Not Found Error",
		},
	},
	REQUEST_NOT_AUTH_CODE: {
		ErrorCode: REQUEST_NOT_AUTH_CODE,
		ErrorMsg: []string{
			"请求未授权",
			"Request Not Auth Error",
		},
	},
	RESOURCE_OUT_OF_LIMIT: {
		ErrorCode: RESOURCE_OUT_OF_LIMIT,
		ErrorMsg: []string{
			"资源超出限制",
			"Resource Out Of Limit Error",
		},
	},
	USER_ALREADY_EXIST_CODE: {
		ErrorCode: USER_ALREADY_EXIST_CODE,
		ErrorMsg: []string{
			"用户已存在",
			"User Already Exist Error",
		},
	},
	USER_ALREADY_IN_GAME_CODE: {
		ErrorCode: USER_ALREADY_IN_GAME_CODE,
		ErrorMsg: []string{
			"用户已在游戏中",
			"User Already In Game Error",
		},
	},
}
