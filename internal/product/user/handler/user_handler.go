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
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
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
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	log.Info("user", req.Account, " try to log in")
	token, err := h.uSvc.Login(req.Account, req.Password)
	if err != nil {
		response.Unauthorized(c, response.CodeInvalidPassword, "invalid account or password")
		return
	}
	res := dto.LoginResponse{Token: token}
	response.Success(c, res)
	log.Info("user login success:", req.Account)
	return
}
