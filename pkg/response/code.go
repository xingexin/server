package response

// 错误码定义
// 格式：模块(2位) + 类型(2位) + 序号(2位)
const (
	// 通用错误码 (10xxxx)
	CodeSuccess       = 0      // 成功
	CodeInternalError = 100000 // 服务器内部错误
	CodeInvalidJSON   = 100001 // JSON 格式错误
	CodeInvalidParams = 100002 // 参数错误
	CodeUnauthorized  = 100003 // 未授权

	// 用户模块错误码 (20xxxx)
	CodeUserNotFound      = 201001 // 用户不存在
	CodeUserAlreadyExists = 201002 // 用户已存在
	CodeInvalidPassword   = 201003 // 账号或密码错误
	CodeTokenInvalid      = 201004 // Token 无效

	// 商品模块错误码 (30xxxx)
	CodeCommodityNotFound     = 301001 // 商品不存在
	CodeCommodityCreateFailed = 301002 // 商品创建失败
	CodeCommodityUpdateFailed = 301003 // 商品更新失败
	CodeCommodityDeleteFailed = 301004 // 商品删除失败
	CodeCommodityQueryFailed  = 301005 // 商品查询失败
)

// 错误消息映射表
var msgMap = map[int]string{
	CodeSuccess:       "操作成功",
	CodeInternalError: "服务器内部错误",
	CodeInvalidJSON:   "JSON 格式错误",
	CodeInvalidParams: "参数错误",
	CodeUnauthorized:  "未授权",

	CodeUserNotFound:      "用户不存在",
	CodeUserAlreadyExists: "用户已存在",
	CodeInvalidPassword:   "账号或密码错误",
	CodeTokenInvalid:      "Token 无效",

	CodeCommodityNotFound:     "商品不存在",
	CodeCommodityCreateFailed: "商品创建失败",
	CodeCommodityUpdateFailed: "商品更新失败",
	CodeCommodityDeleteFailed: "商品删除失败",
	CodeCommodityQueryFailed:  "商品查询失败",
}

// GetMsg 根据错误码获取错误消息
func GetMsg(code int) string {
	msg, ok := msgMap[code]
	if !ok {
		return "未知错误"
	}
	return msg
}
