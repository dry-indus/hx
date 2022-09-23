package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/commodityser"
)

var Commodity CommodityCtr

type CommodityCtr struct{}

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
