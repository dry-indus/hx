package userctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/usermod"
	"hx/service/userser"
)

var Home HomeCtr

type HomeCtr struct{}

func (HomeCtr) List(c context.UserContext) {
	var r usermod.HomeListRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	resp, err := userser.Home.List(c, r)
	if err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	response.Success(c.Gin(), resp)
}

func (HomeCtr) Search(c context.UserContext) {
	var r usermod.HomeSearchRequest
	if err := c.Gin().ShouldBindJSON(&r); err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	resp, err := userser.Home.Search(c, r)
	if err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	response.Success(c.Gin(), resp)
}

func (HomeCtr) OrderInfo(c context.UserContext) {
	var r usermod.OrderInfoRequest
	if err := c.Gin().ShouldBindJSON(&r); err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	resp, err := userser.Order.Info(c, r)
	if err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	response.Success(c.Gin(), resp)
}

func (HomeCtr) SubmitOrder(c context.UserContext) {
	var r usermod.SubmitOrderRequest
	if err := c.Gin().ShouldBindJSON(&r); err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	resp, err := userser.Order.Submit(c, r)
	if err != nil {
		response.Failed(c.Gin(), err)
		return
	}

	response.Success(c.Gin(), resp)
}
