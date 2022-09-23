package mdb

import (
	ctx "context"
	"hx/global"
	"hx/global/context"
	"hx/model/common"
	"time"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type M bson.M

var Commodity CommodityMod

type CommodityMod struct {
	ID         primitive.ObjectID      `bson:"id"`
	Name       string                  `bson:"name"`
	Category   global.MerchantCategory `bson:"category"` // 品类
	MerchantId primitive.ObjectID      `bson:"merchantId"`
	PicURL     string                  `bson:"picUrl"`
	TagIds     []primitive.ObjectID    `bson:"tagIds"`
	SPIds      []primitive.ObjectID    `bson:"spIds"` // 规格和价格id
	Status     CommodityStatus         `bson:"status"`
	Weight     int                     `bson:"weight"` // 权重
	Count      int                     `bson:"count"`  // 商品数量
	CreatedAt  time.Time               `bson:"createdAt"`
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

func (this CommodityMod) Add(c ctx.Context, mod CommodityMod) (primitive.ObjectID, error) {
	r, err := this.Collection().InsertOne(c, mod)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return r.InsertedID.(primitive.ObjectID), nil
}

func (this CommodityMod) Del(c ctx.Context, id primitive.ObjectID) error {
	return this.Collection().RemoveId(c, id)
}

type CommodityUpdateDoc struct {
	Name     *string
	Category *global.MerchantCategory // 品类
	PicURL   *string
	TagIds   []primitive.ObjectID
	Status   *CommodityStatus
}

func (this CommodityMod) UpdateById(c ctx.Context, id primitive.ObjectID, doc *CommodityUpdateDoc) error {

	m := M{}
	if doc.Name != nil {
		m["name"] = *doc.Name
	}
	if doc.Category != nil {
		m["category"] = *doc.Category
	}
	if doc.PicURL != nil {
		m["picUrl"] = *doc.PicURL
	}
	if doc.TagIds != nil {
		m["tagIds"] = doc.TagIds
	}
	if doc.Status != nil {
		m["status"] = *doc.Status
	}

	if len(m) == 0 {
		return nil
	}

	update := bson.M{"$set": m}
	err := this.Collection().UpdateId(c, id, update)
	return err
}

func (this CommodityMod) Page(c context.ContextB, term *CommodityPageTerm, page *common.Page) (list []*CommodityMod, hasNext bool, err error) {
	return this.page(c, term, page.PageNumber, page.PageSize, "-weight")
}

type CommodityPageTerm struct {
	Ids        []primitive.ObjectID
	MerchantId *primitive.ObjectID
	TagIds     []primitive.ObjectID
	Status     *CommodityStatus
}

func (this CommodityPageTerm) Filter() M {
	filter := M{}

	if len(this.Ids) != 0 {
		filter["id"] = M{
			"$in": this.Ids,
		}
	}

	if this.MerchantId != nil {
		filter["merchantId"] = this.MerchantId
	}

	if len(this.TagIds) != 0 {
		filter["tagIds"] = M{
			"$all": this.TagIds,
		}
	}

	if this.Status != nil {
		filter["status"] = this.Status
	}

	return filter
}

func (this CommodityMod) page(c context.ContextB, term *CommodityPageTerm, pageNumber, pageSize int64, sortBy string) (list []*CommodityMod, hasNext bool, err error) {

	filter := term.Filter()

	if len(filter) == 0 {
		return
	}

	skip := (pageNumber - 1) * pageSize
	limit := pageSize + 1

	err = this.Collection().Find(c, filter).Skip(skip).Limit(limit).Sort(sortBy).All(&list)

	if len(list) > int(pageSize) {
		hasNext = true
		list = list[:pageSize]
	}

	return
}

func (this CommodityMod) Find(c context.ContextB, term *CommodityPageTerm) (list []*CommodityMod, err error) {
	filter := term.Filter()
	err = this.Collection().Find(c, filter).All(&list)
	return
}

func (this CommodityMod) FindM(c context.ContextB, term *CommodityPageTerm) (map[primitive.ObjectID]*CommodityMod, error) {
	list := []*CommodityMod{}
	err := this.Collection().Find(c, term).All(&list)
	if err != nil {
		return nil, err
	}

	m := map[primitive.ObjectID]*CommodityMod{}
	for _, v := range list {
		m[v.ID] = v
	}

	return m, nil
}

func (this CommodityMod) FindOnlineByIDs(c context.ContextB, ids ...primitive.ObjectID) (list []*CommodityMod, err error) {
	if len(ids) == 0 {
		return []*CommodityMod{}, nil
	}

	status := Online
	term := &CommodityPageTerm{Ids: ids, Status: &status}
	err = this.Collection().Find(c, term).All(&list)

	return
}

func (this CommodityMod) FindOnlineByIDm(c context.ContextB, ids ...primitive.ObjectID) (map[primitive.ObjectID]*CommodityMod, error) {
	list, err := this.FindOnlineByIDs(c, ids...)
	if err != nil {
		return nil, err
	}

	m := map[primitive.ObjectID]*CommodityMod{}
	for _, v := range list {
		m[v.ID] = v
	}

	return m, nil
}

func (this CommodityMod) FindByTagIds(c context.ContextB, tagIds []primitive.ObjectID) (list []*CommodityMod, err error) {
	if len(tagIds) == 0 {
		return
	}

	filter := M{
		"tagIds": M{
			"$all": tagIds,
		},
	}

	err = this.Collection().Find(c, filter).All(&list)
	return
}
