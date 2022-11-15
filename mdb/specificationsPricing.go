package mdb

import (
	ctx "context"
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/util"
	"time"

	"github.com/qiniu/qmgo"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SpecificationsPricing SpecificationsPricingMod

type SpecificationsPricingMod struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CommodityId    primitive.ObjectID `bson:"commodityId"`
	Specifications string             `bson:"specifications"`
	Count          int                `bson:"count"`
	Pricing        decimal.Decimal    `bson:"pricing"`
	PicURL         string             `bson:"picUrl"`
	CreatedAt      time.Time          `bson:"createdAt"`
}

var specifications_pricing_collection *qmgo.Collection

func SPMD5(sp *SpecificationsPricingMod) string {
	s := fmt.Sprintf("%s_%s", sp.Specifications, sp.Pricing)
	return util.MD5V([]byte(s))
}

func (SpecificationsPricingMod) Collection() *qmgo.Collection {
	if specifications_pricing_collection == nil {
		specifications_pricing_collection = global.DL_CORE_MDB.Collection("specifications_pricing")
	}
	return specifications_pricing_collection
}

func (this SpecificationsPricingMod) AddOne(c ctx.Context, mod *SpecificationsPricingMod) (primitive.ObjectID, error) {
	r, err := this.Collection().InsertOne(c, mod)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}

func (this SpecificationsPricingMod) AddMany(c ctx.Context, mods []*SpecificationsPricingMod) ([]primitive.ObjectID, error) {
	if len(mods) == 0 {
		return nil, nil
	}
	r, err := this.Collection().InsertMany(c, mods)
	if err != nil {
		return nil, err
	}

	ids := []primitive.ObjectID{}
	for _, v := range r.InsertedIDs {
		ids = append(ids, v.(primitive.ObjectID))
	}

	return ids, err
}

func (this SpecificationsPricingMod) DelByID(c ctx.Context, id primitive.ObjectID) error {
	return this.Collection().RemoveId(c, id)
}

type SpecificationsPricingUpdateDoc struct {
	Specifications *string
	Pricing        *decimal.Decimal
	PicURL         *string
}

func (this SpecificationsPricingMod) UpdateById(c ctx.Context, id primitive.ObjectID, doc *SpecificationsPricingUpdateDoc) error {
	m := M{}
	if doc.Specifications != nil {
		m["specifications"] = *doc.Specifications
	}
	if doc.Pricing != nil {
		m["pricing"] = *doc.Pricing
	}
	if doc.PicURL != nil {
		m["picUrl"] = *doc.PicURL
	}

	if len(m) == 0 {
		return nil
	}

	update := bson.M{"$set": m}
	err := this.Collection().UpdateId(c, id, update)
	return err
}

func (this SpecificationsPricingMod) FindByCommodityId(c context.ContextB, commodityId primitive.ObjectID) (list []*SpecificationsPricingMod, err error) {
	filter := M{
		"commodityId": commodityId,
	}

	err = this.Collection().Find(c, filter).All(&list)

	return
}

func (this SpecificationsPricingMod) FindById(c context.ContextB, id primitive.ObjectID) (mod *SpecificationsPricingMod, err error) {
	filter := M{
		"_id": id,
	}

	err = this.Collection().Find(c, filter).One(&mod)

	return
}

func (this SpecificationsPricingMod) CountByCommodityId(c context.ContextB, commodityId primitive.ObjectID) (count int64, err error) {
	filter := M{
		"commodityId": commodityId,
	}

	count, err = this.Collection().Find(c, filter).Count()

	return
}
