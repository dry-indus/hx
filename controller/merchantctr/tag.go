package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/tagser"
)

var Tag TagCtr

type TagCtr struct{}

// @Tags        商户-标签
// @Summary     添加标签
// @Description 添加标签
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.TagAddRequest                              true  "参数"
// @param       hoken    header   string                                                 false "hoken"
// @param       language header   string                                                 false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{Data=merchantmod.TagAddResponse} "成功"
// @Security    Auth
// @Failure     500 {object} response.HTTPResponse "失败"
// @Router      /commodity/tag/add [post]
func (TagCtr) Add(c context.MerchantContext) {
	var r merchantmod.TagAddRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := tagser.Tag.Add(c, c.Merchant().ID, &r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-标签
// @Summary     删除标签
// @Description 删除标签
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.TagDelRequest                              true  "参数"
// @param       hoken    header   string                                                 false "hoken"
// @param       language header   string                                                 false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{Data=merchantmod.TagDelResponse} "成功"
// @Security    Auth
// @Failure     500 {object} response.HTTPResponse "失败"
// @Router      /commodity/tag/del [post]
func (TagCtr) Del(c context.MerchantContext) {
	var r merchantmod.TagDelRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := tagser.Tag.Del(c, &r)
	if err != nil {
		if err == tagser.ErrTagUsed {
			response.Tip(c.Gin(), err.Error()).Failed(err)
			return
		}
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}
