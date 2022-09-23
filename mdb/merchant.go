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
	ID        primitive.ObjectID      `bson:"id"`
	Name      string                  `bson:"name"`
	Password  string                  `bson:"password"`
	Telegram  string                  `bson:"telegram"`
	Category  global.MerchantCategory `bson:"category"` // 品类
	CreatedAt time.Time               `bson:"createdAt"`
}

var merchant_collection *qmgo.Collection

func (MerchantMod) Collection() *qmgo.Collection {
	if merchant_collection == nil {
		merchant_collection = global.DL_CORE_MDB.Collection("merchant")
	}
	return merchant_collection
}

func (this MerchantMod) Create(c context.ContextB, mod MerchantMod) (id primitive.ObjectID, err error) {
	result, err := this.Collection().InsertOne(c, mod)
	id = result.InsertedID.(primitive.ObjectID)

	return
}

func (this MerchantMod) FindOneByName(c context.ContextB, name string) (merchant *MerchantMod, err error) {
	filter := M{
		"name": name,
	}

	err = this.Collection().Find(c, filter).One(&merchant)
	return
}
