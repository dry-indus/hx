package mdb

import (
	"hx/global"
	"hx/global/context"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Merchant MerchantMod

type MerchantMod struct {
	ID        primitive.ObjectID      `bson:"_id,omitempty"`
	Name      string                  `bson:"name"`
	Password  string                  `bson:"password"`
	Prtrait   string                  `bson:"prtrait"` // 头像
	TgID      int64                   `bson:"tgId"`
	TgName    string                  `bson:"tgName"`
	StoreName string                  `bson:"storeName"`
	Category  global.MerchantCategory `bson:"category"` // 品类
	Star      int                     `bson:"star"`
	CreatedAt time.Time               `bson:"createdAt"`
}

var merchant_collection *qmgo.Collection

func (MerchantMod) Collection() *qmgo.Collection {
	if merchant_collection == nil {
		merchant_collection = global.DL_CORE_MDB.Collection("merchant")
	}
	return merchant_collection
}

func (this MerchantMod) Create(c context.ContextB, mod *MerchantMod) (err error) {
	result, err := this.Collection().InsertOne(c, &mod)
	mod.ID = result.InsertedID.(primitive.ObjectID)
	return
}

func (this MerchantMod) FindOneByName(c context.ContextB, name string) (merchant *MerchantMod, err error) {
	return this.FindOne(c, &MerchantTerm{Name: &name})
}

type MerchantTerm struct {
	Id     *primitive.ObjectID
	Name   *string
	TgName *string
}

func (this MerchantTerm) Filter() M {
	filter := M{}

	if this.Id != nil {
		filter["_id"] = this.Id
	}
	if this.Name != nil {
		filter["name"] = this.Name
	}
	if this.TgName != nil {
		filter["tgName"] = this.TgName
	}

	return filter
}

func (this MerchantMod) FindOne(c context.ContextB, term *MerchantTerm) (merchant *MerchantMod, err error) {
	filter := term.Filter()
	if len(filter) == 0 {
		return
	}

	err = this.Collection().Find(c, filter).One(&merchant)
	return
}

func (this MerchantMod) Count(c context.ContextB, term *MerchantTerm) (count int64, err error) {
	filter := term.Filter()
	if len(filter) == 0 {
		return
	}

	return this.Collection().Find(c, filter).Count()
}

func (this MerchantMod) FindByIDs(c context.ContextB, ids []primitive.ObjectID) (list []*MerchantMod, err error) {
	if len(ids) == 0 {
		return []*MerchantMod{}, nil
	}

	filter := M{
		"_id": M{
			"$in": ids,
		},
	}

	err = this.Collection().Find(c, filter).Sort("-_id").All(&list)

	return
}
