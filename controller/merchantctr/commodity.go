package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/commodityser"
)

var Commodity CommodityCtr

type CommodityCtr struct{}

// @Tags        商户-商品
// @Summary     商品列表
// @Description 商品列表
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.CommodityListRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.CommodityListResponse} "成功"
// @Failure     500   {object} response.HTTPResponse
// @Router      /v1/merchant/commodity/list [post]
func (CommodityCtr) List(c context.MerchantContext) {
	var r merchantmod.CommodityListRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := commodityser.Commodity.List(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品
// @Summary     添加商品
// @Description 添加商品
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.CommodityAddRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.CommodityAddResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                        "失败"
// @Router      /v1/merchant/commodity/add [post]
func (CommodityCtr) Add(c context.MerchantContext) {
	var r merchantmod.CommodityAddRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := commodityser.Commodity.Add(c, c.Merchant().ID, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品
// @Summary     编辑商品
// @Description 编辑商品
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.CommodityModifyRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.CommodityModifyResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                           "失败"
// @Router      /v1/merchant/commodity/modify [post]
func (CommodityCtr) Modify(c context.MerchantContext) {
	var r merchantmod.CommodityModifyRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := commodityser.Commodity.Modify(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品
// @Summary     删除商品
// @Description 删除商品
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.CommodityDelRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.CommodityDelResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                        "失败"                                       "内部服务错误"
// @Router      /v1/merchant/commodity/del [post]
func (CommodityCtr) Del(c context.MerchantContext) {
	var r merchantmod.CommodityDelRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := commodityser.Commodity.Del(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品
// @Summary     发布商品
// @Description 商品发布后，用户可见
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.CommodityPublishRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.CommodityPublishResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                            "失败"
// @Router      /v1/merchant/commodity/publish [post]
func (CommodityCtr) Publish(c context.MerchantContext) {
	var r merchantmod.CommodityPublishRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := commodityser.Commodity.Publish(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品
// @Summary     隐藏商品
// @Description 商品隐藏后，用户不可见
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.CommodityHideRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.CommodityHideResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                         "失败"
// @Router      /v1/merchant/commodity/hide [post]
func (CommodityCtr) Hide(c context.MerchantContext) {
	var r merchantmod.CommodityHideRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := commodityser.Commodity.Hide(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}
