package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/merchantser"
)

var Auth AuthCtr

type AuthCtr struct{}

func (AuthCtr) Login(c context.MerchantContext) {
	var r merchantmod.LoginRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}
	resp, err := merchantser.Auth.Login(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

func (AuthCtr) Logout(c context.MerchantContext) {
	var r merchantmod.LogoutRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := merchantser.Auth.Logout(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}

func (AuthCtr) Register(c context.MerchantContext) {
	var r merchantmod.RegisterRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	resp, err := merchantser.Auth.Register(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	response.Success(c.Gin(), resp)
}
