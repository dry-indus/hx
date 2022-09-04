package merchantctr

import (
	"hx/global"
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
	merchant, err := merchantser.Auth.Login(c, r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	s := c.Session()
	{
		token := merchantser.Auth.FlushToken(c, merchant)
		s.Values[global.MERCHANT_TOKEN] = token
	}
	{
		lang := c.Gin().Query("language")
		if len(lang) == 0 {
			lang = global.Application.DefaultLanguage
		}
		s.Values[global.LANGUAGE_KEY] = lang
	}
	s.Save(c.Gin().Request, c.Gin().Writer)

	resp := &merchantmod.RegisterResponse{
		Name:     merchant.Name,
		Telegram: merchant.Telegram,
		Class:    merchant.Class,
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

	err = merchantser.Auth.Register(c, r)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	merchant, err := merchantser.Merchant.FindByName(c, r.Name)
	if err != nil {
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	token := merchantser.Auth.FlushToken(c, merchant)
	s := c.Session()
	s.Values[global.MERCHANT_TOKEN] = token
	s.Save(c.Gin().Request, c.Gin().Writer)

	resp := &merchantmod.RegisterResponse{
		Name:     merchant.Name,
		Telegram: merchant.Telegram,
		Class:    merchant.Class,
	}

	response.Success(c.Gin(), resp)
}
