package mdb

import (
	ctx "context"
	"hx/global"
	"hx/global/context"
	"hx/util"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Tag TagMod

type TagMod struct {
	ID         primitive.ObjectID
	Name       string
	MerchantId primitive.ObjectID
	CreatedAt  time.Time
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
		return primitive.ObjectID{}, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
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

func (this TagMod) DelByID(c ctx.Context, id primitive.ObjectID) error {
	return this.Collection().RemoveId(c, id)
}

func (this TagMod) FindByIDs(c context.ContextB, ids []primitive.ObjectID) (list []*TagMod, err error) {
	if len(ids) == 0 {
		return []*TagMod{}, nil
	}

	filter := M{
		"id": M{
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
		"MerchantId": merchantId,
	}

	err = this.Collection().Find(c, filter).All(&list)

	return
}
