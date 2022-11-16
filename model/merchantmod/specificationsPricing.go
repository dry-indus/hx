package merchantmod

import (
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SPAddRequest struct {
	// CommodityId 商品id
	CommodityId primitive.ObjectID `json:"commodityId" binding:"required" validate:"required"`
	// Specifications 商品规格
	Specifications string `json:"specifications" binding:"required" validate:"required"`
	// Pricing 商品定价
	Pricing decimal.Decimal `json:"pricing" binding:"required" validate:"required"`
	// PicURL 规格与定价缩略图
	PicURL string `json:"picUrl"`
}

type SPAddResponse struct {
	//Id 规格与定价id
	Id primitive.ObjectID `json:"id"`
}

type SPModifyRequest struct {
	//Id 规格与定价id
	Id primitive.ObjectID `json:"id" binding:"required" validate:"required"`
	// Specifications 商品规格
	Specifications *string `json:"specifications" binding:"required" validate:"required"`
	// Pricing 商品定价
	Pricing *decimal.Decimal `json:"pricing" binding:"required" validate:"required"`
	// PicURL 规格与定价缩略图
	PicURL *string `json:"picUrl"`
}

type SPModifyResponse struct {
	//Id 规格与定价id
	Id primitive.ObjectID `json:"id"`
}

type SPDelRequest struct {
	//Id 规格与定价id
	Id primitive.ObjectID `binding:"required" validate:"required"`
}

type SPDelResponse struct {
}
