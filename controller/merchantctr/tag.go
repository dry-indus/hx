package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/tagser"
)

var Tag TagCtr

type TagCtr struct{}

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
