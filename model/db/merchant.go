package db

import (
	"hx/global"
	"hx/global/context"
	"time"

	"github.com/qiniu/qmgo"
)

var Merchant MerchantMod

type MerchantMod struct {
	ID        string
	Name      string
	Password  string
	Telegram  string
	CreatedAt time.Time
}

var merchant_collection *qmgo.Collection

func (MerchantMod) Collection() *qmgo.Collection {
	if merchant_collection == nil {
		merchant_collection = global.DL_CORE_MDB.Collection("merchant")
	}
	return merchant_collection
}

func (this MerchantMod) FindOneByName(c context.ContextB, name string) (merchant *MerchantMod, err error) {
	filter := M{
		"Name": name,
	}

	err = this.Collection().Find(c, filter).One(&merchant)
	return
}
