package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
)

var Land LandCtr

type LandCtr struct{}

func (LandCtr) Redirect(c context.MerchantContext) {
	resp := merchantmod.LandPageResponse{
		Language:     "ZH",
		RedirectPath: "xxx",
	}

	response.Redirect(c.Gin()).Success(resp)
}
