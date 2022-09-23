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
	ID             primitive.ObjectID
	CommodityId    primitive.ObjectID
	Specifications string
	Pricing        decimal.Decimal
	PicURL         string
	ChoiceOpt      ChoiceOpt
	CreatedAt      time.Time
}

type ChoiceOpt int

const (
	SingleChoice   ChoiceOpt = 0
	MultipleChoice ChoiceOpt = 1
	MustChoice     ChoiceOpt = 2
)

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

func (this SpecificationsPricingMod) AddOne(c ctx.Context, mod *SpecificationsPricingMod) error {
	_, err := this.Collection().InsertOne(c, mod)
	return err
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
	ChoiceOpt      *ChoiceOpt
}

func (this SpecificationsPricingMod) UpdateById(c ctx.Context, id primitive.ObjectID, doc *SpecificationsPricingUpdateDoc) error {
	m := M{}
	if doc.Specifications != nil {
		m["Specifications"] = *doc.Specifications
	}
	if doc.Pricing != nil {
		m["Pricing"] = *doc.Pricing
	}
	if doc.PicURL != nil {
		m["PicURL"] = *doc.PicURL
	}
	if doc.ChoiceOpt != nil {
		m["ChoiceOpt"] = *doc.ChoiceOpt
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
		"CommodityId": commodityId,
	}

	err = this.Collection().Find(c, filter).All(&list)

	return
}

func (this SpecificationsPricingMod) FindById(c context.ContextB, id primitive.ObjectID) (mod *SpecificationsPricingMod, err error) {
	filter := M{
		"ID": id,
	}

	err = this.Collection().Find(c, filter).One(&mod)

	return
}
