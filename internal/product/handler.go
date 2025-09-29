package product

import (
	"net/http"
	"server/internal/product/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *service.UserService
}

func NewHandler(s *service.UserService) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Register(c *gin.Context) {
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
	err = h.svc.Register(req.Account, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account or account already exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"create": "create success"})
	return
}

func (h *Handler) Login(c *gin.Context) {
	req := struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	token, err := h.svc.Login(req.Account, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid account or password"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
	return
}
