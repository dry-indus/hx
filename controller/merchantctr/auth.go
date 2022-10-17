package merchantctr

import (
	"errors"
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/merchantser"
	"hx/service/verifyser"
)

var Auth AuthCtr

type AuthCtr struct{}

// @Tags        商户-鉴权
// @Summary     登陆
// @Description 商户登陆
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.LoginRequest                              true  "参数"
// @param       language header   string                                                false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{Data=merchantmod.LoginResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                 "失败"
// @Router      /auth/login [post]
func (this AuthCtr) Login(c context.MerchantContext) {
	var r merchantmod.LoginRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}
	merchant, err := merchantser.Auth.Login(c, r)
	if err != nil {
		response.Tip(c.Gin(), err.Error()).Failed(err)
		return
	}

	resp := &merchantmod.RegisterResponse{
		Name:     merchant.Name,
		Prtrait:  merchant.Prtrait,
		Category: merchant.Category,
	}

	response.Success(c.Gin(), resp)
}

// @Tags        商户-鉴权
// @Summary     注销
// @Description 商户注销
// @Accept      json
// @Produce     json
// @Param       param    body   merchantmod.LogoutRequest true  "参数"
// @param       language header string                    false "语言" default(zh-CN)
// @Security    Auth
// @Success     200 {object} response.HTTPResponse{Data=merchantmod.LogoutResponse} "成功"
// @Failure     500 {object} response.HTTPResponse                                  "失败"
// @Router      /auth/logout [post]
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

// @Tags        商户-鉴权
// @Summary     注册
// @Description 商户注册
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.RegisterRequest                              true  "参数"
// @param       language header   string                                                   false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{Data=merchantmod.RegisterResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                    "失败"
// @Router      /auth/register [post]
func (this AuthCtr) Register(c context.MerchantContext) {
	var r merchantmod.RegisterRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	merchant, err := merchantser.Auth.Register(c, r)
	if err != nil {
		if errors.Is(err, merchantser.ErrPwdNotMatch) ||
			errors.Is(err, merchantser.ErrAccountExists) ||
			errors.Is(err, merchantser.ErrTgExists) ||
			errors.Is(err, merchantser.ErrPwdNotMatch) ||
			errors.Is(err, verifyser.ErrCodeNotMatch) ||
			errors.Is(err, verifyser.ErrTgNameNotMatch) {
			response.InvalidParam(c.Gin(), err.Error()).Failed(nil)
			return
		}
		if errors.Is(err, verifyser.ErrTgId) {
			response.InvalidParam(c.Gin(), "邀请码无效").Failed(nil)
			return
		}
		response.InternalServerError(c.Gin()).Failed(err)
		return
	}

	resp := &merchantmod.RegisterResponse{
		Name:     merchant.Name,
		Prtrait:  merchant.Prtrait,
		Category: merchant.Category,
	}

	response.Success(c.Gin(), resp)
}
