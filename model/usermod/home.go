package usermod

import "hx/model/common"

type HomeListRequest struct {
	*common.Page
}

type HomeListResponse struct {
	List    []*Commodity
	HasNext bool
}

type Commodity struct {
	CommodityID string
	PicURL      string
	Tags        []*Tags
}

type HomeSearchRequest struct {
	TagIDs []string
	*common.Page
}

type Tags struct {
	TagID   string
	TagName string
}

type HomeSearchResponse struct {
	List    []*Commodity
	HasNext bool
}

type OrderInfoRequest struct {
	CommodityID []string
}

type OrderInfoResponse struct {
	Details []*CommodityDetails
}

type CommodityDetails struct {
	CommodityID string
	PicURL      string
	TagNames    []string
	Selects     []*Select
	Invaild     bool   // true: 无效,不可选。否则可选
	InvaildMsg  string // 失效信息
}

type Select struct {
	CommodityID string
	SelectID    string
	SelectName  string
	SelectPrice string
	MD5         string
}

type SubmitOrderRequest struct {
	Selects    []*Select
	TotalPrice string
}

type SubmitOrderResponse struct {
	OrderId        string
	Invaild        bool // true: 无效,不显示订单图。否则显示
	OrderPicBuffer string
	JumpUrl        string
}
