package mdb

import (
	ctx "context"
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/util"
	"time"

	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Tag TagMod

type TagMod struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	MerchantId primitive.ObjectID `bson:"merchantId"`
	Name       string             `bson:"name"`
	CreatedAt  time.Time          `bson:"createdAt"`
}

func (this TagMod) GenMD5() string {
	return util.MD5O(this)
}

var tag_collection *qmgo.Collection

func (TagMod) Collection() *qmgo.Collection {
	if tag_collection == nil {
		tag_collection = global.DL_CORE_MDB.Collection("tag")
	}

	return tag_collection
}

func (this TagMod) AddOne(c ctx.Context, mod *TagMod) (primitive.ObjectID, error) {
	r, err := this.Collection().InsertOne(c, mod)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}

func (this TagMod) UpsertOne(c ctx.Context, term *TagTerm, mod *TagMod) (primitive.ObjectID, error) {

	filter := term.Filter()

	if len(filter) == 0 {
		return primitive.NilObjectID, fmt.Errorf("filter is empty!")
	}

	update := M{
		"$setOnInsert": mod,
		"$set":         M{"name": mod.Name},
	}

	opt := opts.UpdateOptions{}
	opt.SetUpsert(true)

	r, err := this.Collection().UpdateAll(c, filter, update, opt)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return r.UpsertedID.(primitive.ObjectID), nil
}

func (this TagMod) AddMany(c ctx.Context, mods []*TagMod) ([]primitive.ObjectID, error) {
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

type TagTerm struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	MerchantId *primitive.ObjectID `bson:"merchantId"`
	Name       *string             `bson:"name"`
}

func (this TagTerm) Filter() M {
	filter := M{}

	if this.ID != nil {
		filter["_id"] = this.ID
	}
	if this.MerchantId != nil {
		filter["merchantId"] = this.MerchantId
	}
	if this.Name != nil {
		filter["name"] = this.Name
	}

	return filter
}

func (this TagMod) FindByID(c ctx.Context, id primitive.ObjectID) (tag *TagMod, err error) {
	return this.FindOneByTerm(c, &TagTerm{ID: &id})
}

func (this TagMod) FindOneByTerm(c ctx.Context, term *TagTerm) (tag *TagMod, err error) {
	filter := term.Filter()
	err = this.Collection().Find(c, filter).One(&tag)

	return
}

func (this TagMod) CountByTerm(c context.ContextB, term *TagTerm) (count int64, err error) {
	filter := term.Filter()
	count, err = this.Collection().Find(c, filter).Count()

	return
}

func (this TagMod) DelByID(c ctx.Context, id primitive.ObjectID) error {
	return this.Collection().RemoveId(c, id)
}

func (this TagMod) FindByIDs(c context.ContextB, ids []primitive.ObjectID) (list []*TagMod, err error) {
	if len(ids) == 0 {
		return []*TagMod{}, nil
	}

	filter := M{
		"_id": M{
			"$in": ids,
		},
	}

	err = this.Collection().Find(c, filter).All(&list)

	return
}

func (this TagMod) FindByIDm(c context.ContextB, ids []primitive.ObjectID) (map[primitive.ObjectID]*TagMod, error) {
	list, err := this.FindByIDs(c, ids)
	if err != nil {
		return nil, err
	}

	m := map[primitive.ObjectID]*TagMod{}
	for _, v := range list {
		m[v.ID] = v
	}

	return m, nil
}

func (this TagMod) FindByMerchantId(c context.ContextB, merchantId primitive.ObjectID) (list []*TagMod, err error) {
	filter := M{
		"merchantId": merchantId,
	}

	err = this.Collection().Find(c, filter).All(&list)

	return
}
