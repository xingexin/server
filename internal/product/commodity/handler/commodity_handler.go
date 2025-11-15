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

// CommodityHandler 处理商品相关的HTTP请求
type CommodityHandler struct {
	cSvc *service.CommodityService
}

// NewCommodityHandler 创建一个新的商品处理器实例
func NewCommodityHandler(cSvc *service.CommodityService) *CommodityHandler {
	return &CommodityHandler{cSvc: cSvc}
}

// CreateCommodity 处理创建商品请求
// 业务流程：
// 1. 解析请求体中的商品信息（名称、价格、库存）
// 2. 构建商品模型对象
// 3. 调用Service层创建商品（设置创建时间、更新时间）
// 4. 返回创建结果
//
// 注意：
// - 创建时间和更新时间由Service层自动设置
// - 库存初始值需要手动指定，后续可通过Redis缓存管理
func (h *CommodityHandler) CreateCommodity(c *gin.Context) {
	// 解析请求体，绑定到CreateCommodityRequest结构体
	var req dto.CreateCommodityRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		log.Error(err.Error())
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	// 构建商品模型对象
	commodity := &model.Commodity{
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	// 调用Service层创建商品
	err = h.cSvc.CreateCommodity(commodity)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityCreateFailed, err.Error())
		return
	}

	response.SuccessWithMessage(c, "create success", nil)
	return
}

// ListCommodity 处理获取商品列表请求
// 业务流程：
// 1. 调用Service层查询所有商品
// 2. 将商品模型列表转换为响应DTO列表
// 3. 返回商品列表
//
// 返回内容包括：
// - 商品ID、名称、价格、库存
//
// 注意：
// - 此接口返回所有商品，未实现分页功能
// - 库存数据来自MySQL，非实时库存（实时库存在Redis中）
func (h *CommodityHandler) ListCommodity(c *gin.Context) {
	// 调用Service层查询所有商品
	commodities, err := h.cSvc.ListCommodity()
	if err != nil {
		response.BadRequest(c, response.CodeCommodityQueryFailed, err.Error())
		return
	}

	// 将商品模型列表转换为响应DTO列表
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

// UpdateCommodity 处理更新商品请求
// 业务流程：
// 1. 从URL路径中提取商品ID
// 2. 解析请求体中的更新信息
// 3. 调用Service层更新商品（保留创建时间，更新修改时间）
// 4. 返回更新结果
//
// 注意：
// - 更新时间由Service层自动设置为当前时间
// - 创建时间会被保留，不会被覆盖
// - 更新库存不会自动同步到Redis，需手动刷新缓存
func (h *CommodityHandler) UpdateCommodity(c *gin.Context) {
	// 从URL路径参数中获取商品ID（如：/commodity/123）
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidParams, "invalid id parameter")
		return
	}

	// 解析请求体，获取要更新的字段
	var req dto.UpdateCommodityRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, response.CodeInvalidJSON, "invalid JSON")
		return
	}

	// 构建商品模型对象
	commodity := &model.Commodity{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	// 调用Service层更新商品
	err = h.cSvc.UpdateCommodity(commodity)
	if err != nil {
		response.BadRequest(c, response.CodeCommodityUpdateFailed, err.Error())
		return
	}

	response.SuccessWithMessage(c, "update success", nil)
	log.Info("update commodity success:", req.Name)
	return
}

// DeleteCommodity 处理删除商品请求
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

// FindCommodityByName 处理根据名称查找商品请求
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
