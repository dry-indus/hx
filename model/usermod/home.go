package usermod

import (
	"hx/global"
	"hx/model/common"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeListRequest struct {
	common.Page `json:",inline" binding:"required" validate:"required"`
}

type HomeListResponse struct {
	// List 商品列表
	List []*Commodity
	// HasNext true: 有下一页，否则没有下一页
	HasNext bool
}

type HomeSearchRequest struct {
	common.Page `json:",inline" binding:"required" validate:"required"`
	// TagIDs 标签id
	TagIDs []primitive.ObjectID `json:"tagIDs"`
	// CommodityIDs 商品id
	CommodityIDs []primitive.ObjectID `json:"commodityIDs"`
}

type HomeSearchResponse struct {
	// List 商品列表
	List []*Commodity
	// HasNext true: 有下一页，否则没有下一页
	HasNext bool
}

type Commodity struct {
	// ID 商品id
	ID primitive.ObjectID `json:"id"`
	// Name 商品名称
	Name string `json:"name"`
	// Name 商品缩略图
	PicURL string `json:"picUrl"`
	// Tags 商品标签
	Tags []*Tag `json:"tags"`
	// SPs 商品规格与定价,至少有一项
	SPs []*SP `json:"sps"`
	// Category 1:餐饮,2:服饰
	Category global.MerchantCategory `json:"category" enums:"1,2"`
	// Online 1:Online,2:Offline
	Online global.CommodityStatus `json:"online" enums:"0,1"`
	// Show 3:Show,4:Hide
	Show global.CommodityStatus `json:"show" enums:"2,3"`
	// Weight 权重，控制显示顺序 desc
	Weight int `json:"weight"`
	// Count 商品数量
	Count int `json:"count"`
	// CreatedAt 商品创建时间
	CreatedAt time.Time `json:"createdAt"`
	// Invaild true: 无效,不可选。否则可选
	Invaild bool
	// InvaildMsg 失效信息
	InvaildMsg string
}

type Tag struct {
	// ID 标签id
	ID primitive.ObjectID `json:"id"`
	// Name 标签名
	Name string `json:"name"`
}

type SP struct {
	// Id 规格和定价id
	ID primitive.ObjectID `json:"id"`
	// Specifications 商品规格
	// example: 一个，一份，一碗，一件
	Specifications string `json:"specifications" binding:"required" validate:"required" `
	// Pricing 商品定价
	// example: 10，10.5
	Pricing decimal.Decimal `json:"pricing" binding:"required" validate:"required"`
	// BuyCount 购买数量
	BuyCount decimal.Decimal `json:"buyCount" binding:"required" validate:"required"`
	// TotalPricing 商品总价
	// example: TotalPricing = Pricing * BuyCount
	TotalPricing decimal.Decimal `json:"totalPricing" binding:"required" validate:"required"`
	// PicURL 规格和价格缩略图
	PicURL string `json:"picUrl"`
	// Selected true: 已选，否则：未选
	Selected bool `json:"selected"`
	// MD5 规格和定价的指纹
	MD5 string `json:"md5"`
}

type SubmitOrderRequest struct {
	// Commoditys 商品
	Commoditys []*Commodity `json:"commoditys"`
	// TotalPrice 总价
	TotalPrice decimal.Decimal `json:"totalPrice"`
}

type SubmitOrderResponse struct {
	// OrderId 订单id
	OrderId string `json:"orderId"`
	// Invaild true: 无效,不显示订单图。否则显示
	Invaild bool `json:"invaild"`
	// InvaildMsg 失效信息
	InvaildMsg string `json:"invaildMsg"`
	// OrderPicBuffer 订单缩略图
	OrderPicBuffer string `json:"orderPicBuffer"`
	// JumpUrl 跳转连接
	JumpUrl string `json:"jumpUrl"`
	// Commoditys 商品
	Commoditys []*Commodity `json:"commoditys"`
	// TotalPrice 总价
	TotalPrice decimal.Decimal `json:"totalPrice"`
}
