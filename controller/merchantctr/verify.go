package merchantctr

import (
	"hx/global"
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/verifyser"
	"strings"
)

var Verify VerifyCtr

type VerifyCtr struct{}

// @Tags        商户-验证
// @Summary     发送验证码
// @Description 发送验证码
// @Accept      json
// @Produce     json
// @Param       param    body     merchantmod.SendCodeRequest                              true  "参数"
// @param       sence    path     string                                                   true  "验证场景" default(register)
// @param       language header   string                                                   false "语言"   default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=merchantmod.SendCodeResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                    "失败"
// @Router      /verify/{sence}/code/send [post]
func (VerifyCtr) SendCode(c context.MerchantContext) {
	var r merchantmod.SendCodeRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	sence := global.GetSence(c.Gin().Param("sence"))
	if len(sence) == 0 {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	code, err := verifyser.TgVerify.SendCode(c, sence, r.Name, r.TgId, 4)
	if err != nil {
		response.Tip(c.Gin(), "发送失败").Failed(err)
		return
	}

	resp := merchantmod.SendCodeResponse{}
	if !strings.EqualFold(global.ENV, "PRO") {
		resp.VerifyCode = code
	}

	response.Tip(c.Gin(), "发送成功").Success(resp)
}
