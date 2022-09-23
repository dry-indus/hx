package merchantmod

import (
	"hx/mdb"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SPAddRequest struct {
	CommodityId    primitive.ObjectID `binding:"required"`
	Specifications string             `binding:"required"`
	Pricing        decimal.Decimal    `binding:"required"`
	PicURL         string
	ChoiceOpt      mdb.ChoiceOpt
}

type SPAddResponse struct {
}

type SPModifyRequest struct {
	Id             primitive.ObjectID `binding:"required"`
	Specifications *string            `binding:"required"`
	Pricing        *decimal.Decimal   `binding:"required"`
	PicURL         *string
	ChoiceOpt      *mdb.ChoiceOpt
}

type SPModifyResponse struct {
}

type SPDelRequest struct {
	Id primitive.ObjectID `binding:"required"`
}

type SPDelResponse struct {
}
