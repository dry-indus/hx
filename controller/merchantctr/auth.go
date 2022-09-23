package merchantctr

import (
	"hx/global"
	"hx/global/context"
	"hx/global/response"
	"hx/mdb"
	"hx/model/merchantmod"
	"hx/service/merchantser"
)

var Auth AuthCtr

type AuthCtr struct{}

// @Tags        商户-鉴权
// @Summary     登陆
// @Description 商户登陆
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.LoginRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.LoginResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                 "请求失败"
// @Failure     1000  {object} response.HTTPResponse                                 "参数错误"
// @Failure     2000  {object} response.HTTPResponse                                 "内部服务错误"
// @Router      /v1/merchant/auth/login [post]
func (this AuthCtr) Login(c context.MerchantContext) {
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

	this.flushSession(c, merchant)

	resp := &merchantmod.RegisterResponse{
		Name:     merchant.Name,
		Telegram: merchant.Telegram,
		Category: merchant.Category,
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-鉴权
// @Summary     注销
// @Description 商户注销
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.LogoutRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.LogoutResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                  "请求失败"
// @Failure     1000  {object} response.HTTPResponse                                  "参数错误"
// @Failure     2000  {object} response.HTTPResponse                                  "内部服务错误"
// @Router      /v1/merchant/auth/logout [post]
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

	merchantser.Auth.RemoveToken(c, c.Merchant().Name)

	s := c.Session()
	s.Options.MaxAge = -1
	s.Save(c.Gin().Request, c.Gin().Writer)

	response.Success(c.Gin(), resp)
}

// @Tags        商户-鉴权
// @Summary     注册
// @Description 商户注册
// @Accept      json
// @Produce     json
// @Param       param body     merchantmod.RegisterRequest                              true "参数"
// @Success     200   {object} response.HTTPResponse{Data=merchantmod.RegisterResponse} "成功"
// @Failure     500   {object} response.HTTPResponse                                    "请求失败"
// @Failure     1000  {object} response.HTTPResponse                                    "参数错误"
// @Failure     2000  {object} response.HTTPResponse                                    "内部服务错误"
// @Router      /v1/merchant/auth/register [post]
func (this AuthCtr) Register(c context.MerchantContext) {
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

	this.flushSession(c, merchant)

	resp := &merchantmod.RegisterResponse{
		Name:     merchant.Name,
		Telegram: merchant.Telegram,
		Category: merchant.Category,
	}

	response.Success(c.Gin(), resp)
}

func (AuthCtr) flushSession(c context.MerchantContext, merchant *mdb.MerchantMod) {
	s := c.Session()

	{
		token := merchantser.Auth.FlushToken(c, merchant.Name)
		s.Values[global.MERCHANT_TOKEN] = token
	}
	{
		s.Values[global.ACCOUNT] = merchant.Name
	}
	{
		lang := c.Gin().Query(global.LANGUAGE)
		if len(lang) == 0 {
			lang = global.Application.DefaultLanguage
		}
		s.Values[global.LANGUAGE] = lang
	}

	s.Save(c.Gin().Request, c.Gin().Writer)
}
