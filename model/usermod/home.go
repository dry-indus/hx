package usermod

import (
	"hx/model/common"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeListRequest struct {
	*common.Page
}

type HomeListResponse struct {
	List    []*Commodity
	HasNext bool
}

type HomeSearchRequest struct {
	TagIDs       []primitive.ObjectID
	CommodityIDs []primitive.ObjectID
	*common.Page
}

type HomeSearchResponse struct {
	List    []*Commodity
	HasNext bool
}

type Commodity struct {
	ID         primitive.ObjectID
	Name       string
	PicURL     string
	Tags       []*Tag
	SPs        []*SP
	Invaild    bool   // true: 无效,不可选。否则可选
	InvaildMsg string // 失效信息
}

type Tag struct {
	ID   primitive.ObjectID
	Name string
}

type SP struct {
	ID             primitive.ObjectID
	Specifications string
	Pricing        decimal.Decimal
	Selected       bool
	MD5            string
}

type SubmitOrderRequest struct {
	Commoditys []*Commodity
	TotalPrice decimal.Decimal
}

type SubmitOrderResponse struct {
	OrderId        string
	Invaild        bool   // true: 无效,不显示订单图。否则显示
	InvaildMsg     string // 失效信息
	OrderPicBuffer string
	JumpUrl        string
	Commoditys     []*Commodity
	TotalPrice     decimal.Decimal
}
