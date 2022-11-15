package landingctr

import (
	"hx/global"
	"hx/global/context"
	"hx/global/response"
	"hx/model/landingmod"
)

var Pre PreCtr

type PreCtr struct{}

// @Tags        落地页-配置
// @Summary     配置
// @Description 配置
// @Accept      json
// @Produce     json
// @param       language header   string                                                  false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{data=landingmod.SetttingResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                   "失败"
// @Router      /pre/setting [post]
func (this PreCtr) Settting(c context.ContextB) {

	entrys := []landingmod.Entry{}

	for _, v := range global.Landing.Entrys {
		e := landingmod.Entry{
			Name:           v.Name,
			URL:            v.URL,
			BackgroundRPGA: v.BackgroundRPGA,
			Show:           v.Show,
		}
		entrys = append(entrys, e)
	}

	resp := landingmod.SetttingResponse{
		Entrys: entrys,
	}

	response.Success(c.Gin(), resp)
}
