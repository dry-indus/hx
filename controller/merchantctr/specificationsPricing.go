package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/spser"
)

var SP SPCtr

type SPCtr struct{}

// @Tags        商户-商品规格与定价
// @Summary     添加商品规格与定价
// @Description 添加商品规格与定价
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.SPAddRequest                              true  "参数"
// @param       hoken    header   string                                                false "hoken"
// @param       language header   string                                                false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=merchantmod.SPAddResponse} "成功"
// @Security    Auth
// @Failure     500 {object} response.HTTPResponse "失败"
// @Router      /commodity/sp/add [post]
func (SPCtr) Add(c context.MerchantContext) {
	var r merchantmod.SPAddRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := spser.SP.Add(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品规格与定价
// @Summary     编辑商品规格与定价
// @Description 编辑商品规格与定价
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.SPModifyRequest                              true  "参数"
// @param       hoken    header   string                                                   false "hoken"
// @param       language header   string                                                   false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=merchantmod.SPModifyResponse} "成功"
// @Security    Auth
// @Failure     500 {object} response.HTTPResponse "失败"
// @Router      /commodity/sp/modify [post]
func (SPCtr) Modify(c context.MerchantContext) {
	var r merchantmod.SPModifyRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := spser.SP.Modify(c, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-商品规格与定价
// @Summary     删除商品规格与定价
// @Description 删除商品规格与定价，每个商品至少保留一项
// @Accept      json
// @Produce     json
// @Param       param body merchantmod.SPDelRequest true "参数"
// @Security    Auth
// @Success     200 {object} response.HTTPResponse{data=merchantmod.SPDelResponse} "成功"
// @Failure     500 {object} response.HTTPResponse                                 "失败"
// @Router      /commodity/sp/del [post]
func (SPCtr) Del(c context.MerchantContext) {
	var r merchantmod.SPDelRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := spser.SP.Del(c, &r)
	if err != nil {
		if err == spser.ErrOneMustBeRetained {
			response.Tip(c.Gin(), err.Error()).Failed(err)
			return
		}
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}
