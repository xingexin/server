package handler

import (
	"server/internal/product/model"
	"server/internal/product/service"
	"server/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CommodityHandler struct {
	cSvc *service.CommodityService
}

func NewCommodityHandler(cSvc *service.CommodityService) *CommodityHandler {
	return &CommodityHandler{cSvc: cSvc}
}

func (h *CommodityHandler) CreateCommodity(c *gin.Context) {
	req := &model.Commodity{}
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Error(err.Error())
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cSvc.CreateCommodity(req)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityCreateFailed, err.Error())
		return
	}
	response.SuccessWithMessage(c, "create success", nil)
	return
}

func (h *CommodityHandler) ListCommodity(c *gin.Context) {
	commodities, err := h.cSvc.ListCommodity()
	if err != nil {
		response.BadRequest(c, response.CodeCommodityQueryFailed, err.Error())
		return
	}
	res := make([]interface{}, 0, len(commodities))
	for _, cdt := range commodities {
		res = append(res, gin.H{
			"id":    cdt.ID,
			"name":  cdt.Name,
			"price": cdt.Price,
			"stock": cdt.Stock,
		})
	}
	response.Success(c, res)
	log.Info("user", "list commodity success")
	return
}

func (h *CommodityHandler) UpdateCommodity(c *gin.Context) {
	req := &model.Commodity{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cSvc.UpdateCommodity(req)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityUpdateFailed, err.Error())
		return
	}
	response.SuccessWithMessage(c, "update success", nil)
	log.Info("update commodity success:", req.Name)
	return
}

func (h *CommodityHandler) DeleteCommodity(c *gin.Context) {
	req := struct {
		ID int `json:"id"`
	}{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	err = h.cSvc.RemoveCommodity(req.ID)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityDeleteFailed, err.Error())
		return
	}
	response.SuccessWithMessage(c, "delete success", nil)
	log.Info("delete commodity success:", req.ID)
	return
}

func (h *CommodityHandler) FindCommodityByName(c *gin.Context) {
	req := struct {
		Name string `json:"name"`
	}{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	commodities, err := h.cSvc.FindCommodityByName(req.Name)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityQueryFailed, err.Error())
		return
	}
	res := make([]interface{}, 0, len(commodities))
	for _, cdt := range commodities {
		res = append(res, gin.H{
			"id":    cdt.ID,
			"name":  cdt.Name,
			"price": cdt.Price,
			"stock": cdt.Stock,
		})
	}
	response.Success(c, res)
	log.Info("user", "find commodity by name success:", req.Name)
	return
}
