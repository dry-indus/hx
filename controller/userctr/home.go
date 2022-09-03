package userctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/usermod"
	"hx/service/orderser.go"
	"hx/service/userser"
)

var Home HomeCtr

type HomeCtr struct{}

// @Tags        用户-首页
// @Summary     商品列表
// @Description 首页核心接口，展示商品列表
// @Accept      json
// @Produce     json
// @Param       param       body     usermod.HomeListRequest                              true "参数"
// @Param       xx_merchant path     string                                               true "Merchant Name"
// @Success     200         {object} response.HTTPResponse{Data=usermod.HomeListResponse} "成功"
// @Failure     500         {object} response.HTTPResponse                                "请求失败"
// @Router      /v1/user/home/list/{xx_merchant} [post]
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

func (HomeCtr) OrderInfo(c context.UserContext) {
	var r usermod.OrderInfoRequest
	if err := c.Gin().ShouldBindJSON(&r); err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := orderser.Order.Info(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

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
