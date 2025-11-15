package handler

import (
	"server/internal/product/user/dto"
	"server/internal/product/user/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// UserHandler 处理用户相关的HTTP请求
type UserHandler struct {
	uSvc *service.UserService
}

// NewUserHandler 创建一个新的用户处理器实例
func NewUserHandler(uSvc *service.UserService) *UserHandler {
	return &UserHandler{uSvc: uSvc}
}

// Register 处理用户注册请求
// 业务流程：
// 1. 解析请求体中的注册信息（账号、密码、姓名）
// 2. 调用Service层进行注册（包含：账号重复检查、密码加密、创建用户记录）
// 3. 返回注册结果
//
// 注意：
// - 账号必须唯一，重复账号会返回错误
// - 密码在Service层会进行bcrypt加密后存储
// - 建议前端对密码强度进行校验
func (h *UserHandler) Register(c *gin.Context) {
	// 解析请求体，绑定到RegisterRequest结构体
	var req dto.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	// 调用Service层进行用户注册
	// 内部流程：检查账号是否存在 -> 加密密码 -> 创建用户记录
	err = h.uSvc.Register(req.Account, req.Password, req.Name)
	if err != nil {
		response.BadRequest(c, response.CodeUserAlreadyExists, "invalid account or account already exist")
		return
	}

	response.SuccessWithMessage(c, "create success", nil)
	log.Info("user register success:", req.Account)
	return
}

// Login 处理用户登录请求，返回JWT令牌
// 业务流程：
// 1. 解析请求体中的登录信息（账号、密码）
// 2. 调用Service层进行身份验证（密码校验）
// 3. 验证成功后生成JWT令牌并返回
//
// JWT令牌包含的Claims：
// - userID: 用户ID
// - account: 用户账号
// - exp: 过期时间（默认24小时）
//
// 注意：
// - 客户端需将令牌存储并在后续请求的Authorization头中携带
// - 格式：Authorization: Bearer <token>
func (h *UserHandler) Login(c *gin.Context) {
	// 解析请求体，绑定到LoginRequest结构体
	var req dto.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	log.Info("user", req.Account, " try to log in")

	// 调用Service层进行登录验证
	// 内部流程：查询用户 -> bcrypt密码验证 -> 生成JWT令牌
	token, err := h.uSvc.Login(req.Account, req.Password)
	if err != nil {
		response.Unauthorized(c, response.CodeInvalidPassword, "invalid account or password")
		return
	}

	// 返回JWT令牌
	res := dto.LoginResponse{Token: token}
	response.Success(c, res)
	log.Info("user login success:", req.Account)
	return
}
