package db

import (
	"hx/global"
	"hx/global/context"
	"hx/model/common"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type M bson.M

var Commodity CommodityMod

type CommodityMod struct {
	ID         string
	MerchantId string
	PicURL     string
	TagIds     []string
	Attribute  Attribute // public private
	Status     CommodityStatus
	CreatedAt  time.Time
}

type CommodityStatus int

const (
	Online  CommodityStatus = 0
	Offline CommodityStatus = 1
	Show    CommodityStatus = 2
	Hide    CommodityStatus = 3
)

var commodity_collection *qmgo.Collection

func (CommodityMod) Collection() *qmgo.Collection {
	if commodity_collection == nil {
		commodity_collection = global.DL_CORE_MDB.Collection("commodity")
	}
	return commodity_collection
}

func (this CommodityMod) FindOnline(c context.ContextB, merchantId string, tagIds []string, page *common.Page) (list []*CommodityMod, hasNext bool, err error) {
	filter := M{
		"MerchantId": merchantId,
		"Status":     Online,
	}

	if len(tagIds) != 0 {
		filter["TagIds"] = M{
			"$all": tagIds,
		}
	}

	skip := (page.PageNumber - 1) * page.PageSize
	limit := page.PageSize + 1

	err = this.Collection().Find(c, filter).Skip(skip).Limit(limit).Sort("-CreatedAt").All(&list)

	if len(list) > int(page.PageSize) {
		hasNext = true
		list = list[:page.PageSize]
	}

	return
}

func (this CommodityMod) FindOnlineByIDs(c context.ContextB, ids []string) (list []*CommodityMod, err error) {
	if len(ids) == 0 {
		return []*CommodityMod{}, nil
	}

	filter := M{
		"id": M{
			"$in": ids,
		},
		"Status": Online,
	}

	err = this.Collection().Find(c, filter).All(&list)

	return
}

func (this CommodityMod) FindOnlineByIDm(c context.ContextB, ids []string) (map[string]*CommodityMod, error) {
	list, err := this.FindOnlineByIDs(c, ids)
	if err != nil {
		return nil, err
	}

	m := map[string]*CommodityMod{}
	for _, v := range list {
		m[v.ID] = v
	}

	return m, nil
}
