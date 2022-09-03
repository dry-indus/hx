package mdb

import (
	"hx/global"
	"hx/global/context"
	"hx/util"
	"time"

	"github.com/qiniu/qmgo"
)

var Tag TagMod

type TagMod struct {
	ID          string
	Name        string
	Value       string
	CommodityID string // 商品ID
	MerchantId  string
	Attribute   Attribute // public private
	Type        TagType   // server age nationality
	CreatedAt   time.Time
}

func (this TagMod) GenMD5() string {
	return util.MD5O(this)
}

type TagType string

const (
	Server      TagType = "server"
	Age         TagType = "age"
	Nationality TagType = "nationality"
)

type Attribute string

const (
	Public  Attribute = "public"
	Private Attribute = "private"
)

var tag_collection *qmgo.Collection

func (TagMod) Collection() *qmgo.Collection {
	if tag_collection == nil {
		tag_collection = global.DL_CORE_MDB.Collection("tag")
	}
	return tag_collection
}

func (this TagMod) FindByIDs(c context.ContextB, ids []string) (list []*TagMod, err error) {
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

func (this TagMod) FindByIDm(c context.ContextB, ids []string) (map[string]*TagMod, error) {
	list, err := this.FindByIDs(c, ids)
	if err != nil {
		return nil, err
	}

	m := map[string]*TagMod{}
	for _, v := range list {
		m[v.ID] = v
	}

	return m, nil

}
