package handler

import (
	"server/internal/product/commodity/dto"
	"server/internal/product/commodity/model"
	"server/internal/product/commodity/service"
	"server/pkg/response"
	"strconv"

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
	var req dto.CreateCommodityRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Error(err.Error())
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	commodity := &model.Commodity{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	err = h.cSvc.CreateCommodity(commodity)
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
	res := make([]dto.CommodityResponse, 0, len(commodities))
	for _, cdt := range commodities {
		res = append(res, dto.CommodityResponse{
			ID:    cdt.ID,
			Name:  cdt.Name,
			Price: cdt.Price,
			Stock: cdt.Stock,
		})
	}
	response.Success(c, res)
	log.Info("user", "list commodity success")
	return
}

func (h *CommodityHandler) UpdateCommodity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	var req dto.UpdateCommodityRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}
	commodity := &model.Commodity{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}
	err = h.cSvc.UpdateCommodity(commodity)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityUpdateFailed, err.Error())
		return
	}
	response.SuccessWithMessage(c, "update success", nil)
	log.Info("update commodity success:", req.Name)
	return
}

func (h *CommodityHandler) DeleteCommodity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}
	err = h.cSvc.RemoveCommodity(id)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityDeleteFailed, err.Error())
		return
	}
	response.SuccessWithMessage(c, "delete success", nil)
	log.Info("delete commodity success:", id)
	return
}

func (h *CommodityHandler) FindCommodityByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		response.BadRequest(c, response.CodeInvalidParams, "name parameter is required")
		return
	}

	commodities, err := h.cSvc.FindCommodityByName(name)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityQueryFailed, err.Error())
		return
	}
	res := make([]dto.CommodityResponse, 0, len(commodities))
	for _, cdt := range commodities {
		res = append(res, dto.CommodityResponse{
			ID:    cdt.ID,
			Name:  cdt.Name,
			Price: cdt.Price,
			Stock: cdt.Stock,
		})
	}
	response.Success(c, res)
	log.Info("user", "find commodity by name success:", name)
	return
}
