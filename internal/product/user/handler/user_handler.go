package handler

import (
	"server/internal/product/user/dto"
	"server/internal/product/user/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserHandler struct {
	uSvc *service.UserService
}

func NewUserHandler(uSvc *service.UserService) *UserHandler {
	return &UserHandler{uSvc: uSvc}
}

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
