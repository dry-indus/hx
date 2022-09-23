package merchantmod

import (
	"hx/mdb"
	"hx/model/common"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SelectedID struct {
	Id       primitive.ObjectID
	Selected bool
}

type Commodity struct {
	ID     primitive.ObjectID
	Name   string
	PicURL string
	Tags   []*Tag
	SPs    []*SP
}

type Tag struct {
	ID   primitive.ObjectID
	Name string
}

type SP struct {
	ID             primitive.ObjectID
	Specifications string
	Pricing        decimal.Decimal
	PicURL         string
}

type CommodityListRequest struct {
	*common.Page
}

type CommodityListResponse struct {
	List    []*Commodity
	AllTags []*Tag
	HasNext bool
}

type CommodityAddRequest struct {
	Commoditys []*CommodityAdd
}

type CommodityAdd struct {
	Name     string
	Category int `binding:"required"` // 品类
	PicURL   string
	Tags     []SelectedID
	SPs      []struct {
		Specifications string
		Pricing        decimal.Decimal
	} `binding:"required"` // 规格和价格
	Attribute mdb.Attribute `binding:"required"` // public: 公共商品 private: 私有商品
}

type CommodityAddResponse struct {
	Count int
}

type CommodityModifyRequest struct {
	Id        primitive.ObjectID `binding:"required"`
	Name      *string
	Category  *int // 品类
	PicURL    *string
	Tags      []SelectedID
	Attribute *mdb.Attribute // public: 公共商品 private: 私有商品
}

type CommodityModifyResponse struct {
}

type CommodityDelRequest struct {
	Id primitive.ObjectID `binding:"required"`
}

type CommodityDelResponse struct {
}

type CommodityPublishRequest struct {
	Id primitive.ObjectID `binding:"required"`
}

type CommodityPublishResponse struct {
}

type CommodityHideRequest struct {
	Id primitive.ObjectID `binding:"required"`
}

type CommodityHideResponse struct {
}
