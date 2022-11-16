package landingmod

import (
	"hx/model/common"
	"hx/model/searchmod"
)

type StoreSearchRequest struct {
	// Keywords 搜索关键字
	Keywords    string `json:"keywords" binding:"required" validate:"required"`
	common.Page `json:",inline" binding:"required" validate:"required"`
}

type StoreSearchResponse struct {
	// Keywords 搜索关键字
	Keywords string `json:"keywords"`
	// 搜索关键字建议列表
	Suggest []string `json:"suggest"`
	// 搜索结果
	Result []*searchmod.Store `json:"result"`
}

type SearchPushRequest struct {
	Key string `json:"key" binding:"required" validate:"required"`
	Val string `json:"val" binding:"required" validate:"required"`
}
