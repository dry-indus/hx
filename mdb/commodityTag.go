package mdb

import (
	ctx "context"
	"fmt"
	"hx/global"
	"hx/global/context"
	"time"

	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var CommodityTag CommodityTagMod

type CommodityTagMod struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	MerchantId  primitive.ObjectID `bson:"merchantId"`
	CommodityId primitive.ObjectID `bson:"commodityId"`
	TagId       primitive.ObjectID `bson:"tagId"`
	Enable      bool               `bson:"enable"`
	Show        bool               `bson:"show"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

var commodity_tag_collection *qmgo.Collection

func (CommodityTagMod) Collection() *qmgo.Collection {
	if commodity_collection == nil {
		commodity_collection = global.DL_CORE_MDB.Collection("commodityTag")
	}
	return commodity_collection
}

func (this CommodityTagMod) AddMany(c ctx.Context, docs []*CommodityTagMod) (insertedIDs []primitive.ObjectID, err error) {
	r, err := this.Collection().InsertMany(c, docs)

	for _, v := range r.InsertedIDs {
		insertedIDs = append(insertedIDs, v.(primitive.ObjectID))
	}

	return
}

type CommodityTagSlice []*CommodityTagMod

func (c CommodityTagSlice) TagIds() (tagIds []primitive.ObjectID) {
	tagIdm := map[primitive.ObjectID]byte{}

	for _, v := range c {
		if _, ok := tagIdm[v.TagId]; !ok {
			tagIdm[v.TagId] = 1
			tagIds = append(tagIds, v.TagId)
		}
	}

	return
}

func (c CommodityTagSlice) CommodityIds() (commodityIds []primitive.ObjectID) {
	commodityIdm := map[primitive.ObjectID]byte{}

	for _, v := range c {
		if _, ok := commodityIdm[v.CommodityId]; !ok {
			commodityIdm[v.CommodityId] = 1
			commodityIds = append(commodityIds, v.CommodityId)
		}
	}

	return
}

type CommodityTagTerm struct {
	ID          *primitive.ObjectID `bson:"_id,omitempty"`
	MerchantId  *primitive.ObjectID `bson:"merchantId"`
	CommodityId *primitive.ObjectID `bson:"commodityId"`
	TagId       *primitive.ObjectID `bson:"tagId"`
	Enable      *bool               `bson:"enable"`
	Show        *bool               `bson:"show"`
}

func (this CommodityTagTerm) Filter() M {
	filter := M{}

	if this.ID != nil {
		filter["_id"] = this.ID
	}
	if this.MerchantId != nil {
		filter["merchantId"] = this.MerchantId
	}
	if this.CommodityId != nil {
		filter["commodityId"] = this.CommodityId
	}
	if this.TagId != nil {
		filter["tagId"] = this.TagId
	}
	if this.Enable != nil {
		filter["enable"] = this.Enable
	}
	if this.Show != nil {
		filter["show"] = this.Show
	}

	return filter
}

func (this CommodityTagMod) FindByTerm(c context.ContextB, term *CommodityTagTerm) (list CommodityTagSlice, err error) {
	filter := term.Filter()
	if len(filter) == 0 {
		return
	}

	err = this.Collection().Find(c, filter).Sort("-_id").All(&list)
	return
}

func (this CommodityTagMod) DelByTerm(c ctx.Context, term *CommodityTagTerm) error {
	filter := term.Filter()
	if len(filter) == 0 {
		return nil
	}

	return this.Collection().Remove(c, filter)
}

func (this CommodityTagMod) UpsertOne(c ctx.Context, term *CommodityTagTerm, mod *CommodityTagMod) (primitive.ObjectID, error) {

	filter := term.Filter()

	if len(filter) == 0 {
		return primitive.NilObjectID, fmt.Errorf("filter is empty!")
	}

	update := M{
		"$setOnInsert": mod,
		"$set": M{
			"tagId":       mod.TagId,
			"commodityId": mod.CommodityId,
		},
	}

	opt := opts.UpdateOptions{}
	opt.SetUpsert(true)

	r, err := this.Collection().UpdateAll(c, filter, update, opt)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return r.UpsertedID.(primitive.ObjectID), nil
}
