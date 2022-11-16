package userctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/usermod"
	"hx/service/orderser"
	"hx/service/userser"
)

var Home HomeCtr

type HomeCtr struct{}

// @Tags        用户-首页
// @Summary     商品列表
// @Description 首页核心接口，展示商品列表
// @Accept      json
// @Produce     json
// @Param       param    body     usermod.HomeListRequest                              true  "参数"
// @Param       merchant header   string                                               false "Merchant Name" default(default)
// @param       language header   string                                               false "语言"            default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=usermod.HomeListResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                "请求失败"
// @Router      /home/list [post]
func (HomeCtr) List(c context.UserContext) {
	var r usermod.HomeListRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := userser.Home.List(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        用户-搜索商品
// @Summary     搜索商品
// @Description 搜索商品
// @Accept      json
// @Produce     json
// @Param       param    body     usermod.HomeSearchRequest                              true  "参数"
// @Param       merchant header   string                                                 false "Merchant Name" default(default)
// @param       language header   string                                                 false "语言"            default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=usermod.HomeSearchResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                  "请求失败"
// @Router      /home/search [post]
func (HomeCtr) Search(c context.UserContext) {
	var r usermod.HomeSearchRequest
	if err := c.Gin().ShouldBindJSON(&r); err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := userser.Home.Search(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        用户-提交订单
// @Summary     提交订单
// @Description 提交并审核订单，为有效订单提供缩略图，为无效订单提供失效信息
// @Accept      json
// @Produce     json
// @Param       param    body     usermod.SubmitOrderRequest                              true  "参数"
// @Param       merchant header   string                                                  false "Merchant Name" default(default)
// @param       language header   string                                                  false "语言"            default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=usermod.SubmitOrderResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                   "请求失败"
// @Router      /home/order/submit [post]
func (HomeCtr) SubmitOrder(c context.UserContext) {
	var r usermod.SubmitOrderRequest
	if err := c.Gin().ShouldBindJSON(&r); err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := orderser.Order.Submit(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}
