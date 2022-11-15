package landingctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/landingmod"
	"hx/service/searchser"
)

var Store StoreCtr

type StoreCtr struct{}

// @Tags        落地页-商铺搜索
// @Summary     根据关键字搜索店铺，返回搜索建议，和搜索结果
// @Description 搜索建议：用于补全关键字；搜索结果：店铺列表
// @Accept      json
// @Produce     json
// @Param       param    body     landingmod.StoreSearchRequest                              true  "参数"
// @param       language header   string                                                     false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse{Data=landingmod.StoreSearchResponse} "成功"
// @Failure     500      {object} response.HTTPResponse                                      "失败"
// @Router      /store/search/ [post]
func (this StoreCtr) Search(c context.ContextB) {
	var r landingmod.StoreSearchRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	suggest, result := searchser.Search.SearchStore(c, r.Keywords, r.Page)

	resp := landingmod.StoreSearchResponse{
		Keywords: r.Keywords,
		Suggest:  suggest,
		Result:   result,
	}

	response.Success(c.Gin(), resp)
}

// @Tags        测试-插入搜索数据
// @Summary     插入搜索数据,
// @Description 主要用来验证搜索接口的关键字建议
// @Accept      json
// @Produce     json
// @Param       param    body     landingmod.SearchPushRequest                              true  "参数"
// @param       language header   string                                                     false "语言" default(zh-CN)
// @Success     200      {object} response.HTTPResponse "成功"
// @Failure     500      {object} response.HTTPResponse                                      "失败"
// @Router      /test/push/ [post]
func (this StoreCtr) SearchPush(c context.ContextB) {
	var r landingmod.SearchPushRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	searchser.Search.Push(c, r.Key, r.Val)

	response.Success(c.Gin())
}
