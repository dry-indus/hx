package userctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/usermod"
)

var Land LandCtr

type LandCtr struct{}

func (LandCtr) Page(c context.UserContext) {
	var r usermod.HomeListRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}
	// c.Gin().Get()
}

func (LandCtr) Redirect(c context.UserContext) {
	resp := usermod.LandPageResponse{
		Language:     "ZH",
		RedirectPath: "xxx",
	}
	// c.Session().
	response.Redirect(c.Gin()).Success(resp)
}
