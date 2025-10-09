package product

import (
	"net/http"
	"server/internal/product/model"
	"server/internal/product/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	uSvc *service.UserService
	cSvc *service.CommodityService
}

func NewHandler(uSvc *service.UserService, cSvc *service.CommodityService) *Handler {
	return &Handler{uSvc: uSvc, cSvc: cSvc}
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
	err = h.uSvc.Register(req.Account, req.Password, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account or account already exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"create": "create success"})
	log.Info("user register success:", req.Account)
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

func (h *Handler) CreateCommodity(c *gin.Context) {
	req := &model.Commodity{}
	err := c.ShouldBindJSON(&req)

	if err != nil {
		println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	err = h.cSvc.CreateCommodity(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"create": "create success"})
	return
}

func (h *Handler) ListCommodity(c *gin.Context) {
	commodities, err := h.cSvc.ListCommodity()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, commodities)
	return
}

func (h *Handler) UpdateCommodity(c *gin.Context) {
	req := &model.Commodity{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	err = h.cSvc.UpdateCommodity(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"update": "update success"})
	log.Info("update commodity success:", req.Name)
	return
}
