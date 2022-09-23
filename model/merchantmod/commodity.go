package merchantmod

import (
	"hx/model/common"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Commodity struct {
	ID     primitive.ObjectID `json:"id"`
	Name   string             `json:"name"`
	PicURL string             `json:"picUrl"`
	Tags   []*Tag             `json:"tags"`
	SPs    []*SP              `json:"sps"`
}

type Tag struct {
	// ID 标签id
	ID primitive.ObjectID `json:"id"`
	// Name 标签名
	Name string `json:"name"`
	// Selected true:已选择，否则未选择
	Selected bool `json:"selected"`
}

type SP struct {
	ID             primitive.ObjectID `json:"id"`
	Specifications string             `json:"specifications"`
	Pricing        decimal.Decimal    `json:"pricing"`
	PicURL         string             `json:"picUrl"`
}

type CommodityListRequest struct {
	// Page 页面
	*common.Page `json:"page"`
}

type CommodityListResponse struct {
	// List 商品列表
	List []*Commodity `json:"list"`
	// AllTags 商户设置的所有标签
	AllTags []*Tag `json:"allTags"`
	// HasNext true: 有下一页
	HasNext bool `json:"hasNext"`
}

type CommodityAddRequest struct {
	// Commoditys 需要添加的商品列表
	Commoditys []*CommodityAdd `json:"commoditys" binding:"required" validate:"required"`
}

type CommodityAdd struct {
	// Name 商品名称
	Name string `json:"name"`
	// PicURL 商品缩略图
	PicURL string `json:"picURL" binding:"required" validate:"required"`
	// Tags 商品标签
	Tags []*Tag `json:"tags"`
	// SPs 商品规格和定价,至少有一个
	SPs []struct {
		// Specifications 商品规格
		// example: 一个，一份，一碗，一件
		Specifications string `json:"specifications" binding:"required" validate:"required" `
		// Pricing 商品定价
		// example: 10，10.5
		Pricing decimal.Decimal `json:"pricing" binding:"required" validate:"required"`
	} `json:"sps" binding:"required" validate:"required"` // 规格和价格
}

type CommodityAddResponse struct {
	Count int `json:"count"`
}

type CommodityModifyRequest struct {
	Id     primitive.ObjectID `json:"id" binding:"required" validate:"required"`
	Name   *string            `json:"name"`
	PicURL *string            `json:"picURL"`
	// Tags 重设的标签列表, 仅设置selected:true 的标签
	// example: [{"id":"id","selected":true}]
	Tags []*Tag `json:"tags"`
}

type CommodityModifyResponse struct {
}

type CommodityDelRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type CommodityDelResponse struct {
}

type CommodityPublishRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type CommodityPublishResponse struct {
}

type CommodityHideRequest struct {
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type CommodityHideResponse struct {
}
