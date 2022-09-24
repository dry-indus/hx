package merchantmod

import (
	"hx/global"
	"hx/model/common"
	"time"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
}

type Tag struct {
	// ID 标签id
	ID primitive.ObjectID `json:"id"`
	// Name 标签名
	Name string `json:"name"`
	// Selected true:已选择，否则未选择
	Selected bool `json:"selected"`
}

//SP 规格和定价
type SP struct {
	// Id 规格和定价id
	ID primitive.ObjectID `json:"id"`
	// Specifications 商品规格
	// example: 一个，一份，一碗，一件
	Specifications string `json:"specifications" binding:"required" validate:"required" `
	// Pricing 商品定价
	// example: 10，10.5
	Pricing decimal.Decimal `json:"pricing" binding:"required" validate:"required"`
	// PicURL 规格和价格缩略图
	PicURL string `json:"picUrl"`
	// ChoiceOpt 选择设置 0:单选；1:多选；2:必选
	ChoiceOpt global.ChoiceOpt `json:"choiceOpt" enums:"0,1,2" `
}

type CommodityListRequest struct {
	common.Page `json:",inline" binding:"required" validate:"required"`
}

type CommodityListResponse struct {
	// List 商品列表
	List []*Commodity `json:"list"`
	// AllTags 商户设置的所有标签
	AllTags []*Tag `json:"allTags"`
	// HasNext true: 有下一页，否则没有下一页
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
	// SPs 商品规格和定价,至少有一项
	SPs []*SP `json:"sps" binding:"required" validate:"required"` // 规格和价格
}

type CommodityAddResponse struct {
	// Ids 新增商品id
	Ids []primitive.ObjectID `json:"ids"`
	// Count 新增商品数量
	Count int `json:"count"`
}

type CommodityModifyRequest struct {
	// Id 商品id
	Id primitive.ObjectID `json:"id" binding:"required" validate:"required"`
	// PicURL 商品名称
	Name *string `json:"name"`
	// PicURL 商品缩略图
	PicURL *string `json:"picURL"`
	// Tags 重设的标签列表, 仅设置selected:true 的标签
	// example: [{"id":"id","selected":true}]
	Tags []*Tag `json:"tags"`
}

type CommodityModifyResponse struct {
	// Id 商品id
	Id primitive.ObjectID `json:"id"`
}

type CommodityDelRequest struct {
	// Id 商品id
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type CommodityDelResponse struct {
}

type CommodityPublishRequest struct {
	// Id 商品id
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type CommodityPublishResponse struct {
}

type CommodityHideRequest struct {
	// Id 商品id
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type CommodityHideResponse struct {
}
