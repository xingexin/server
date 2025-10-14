package handler

import (
	"net/http"
	"server/internal/product/service"

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
	req := struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	err = h.uSvc.Register(req.Account, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account or account already exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"create": "create success"})
	log.Info("user register success:", req.Account)
	return
}

func (h *UserHandler) Login(c *gin.Context) {
	req := struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	log.Info("user", req.Account, " try to log in")
	token, err := h.uSvc.Login(req.Account, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid account or password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
	log.Info("user login success:", req.Account)
	return
}
