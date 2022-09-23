package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/spser"
)

var SP SPCtr

type SPCtr struct{}

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

func (SPCtr) Modify(c context.MerchantContext) {
	var r merchantmod.SPModifyRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := spser.SP.Modify(c, &r)
	if err != nil {
		if err == spser.ErrSPUsed {
			response.Tip(c.Gin(), err.Error()).Failed(err)
			return
		}
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

func (SPCtr) Del(c context.MerchantContext) {
	var r merchantmod.SPDelRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := spser.SP.Del(c, &r)
	if err != nil {
		if err == spser.ErrSPUsed {
			response.Tip(c.Gin(), err.Error()).Failed(err)
			return
		}
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}
