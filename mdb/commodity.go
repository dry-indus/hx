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
	ID         primitive.ObjectID      `bson:"_id,omitempty"`
	Name       string                  `bson:"name"`
	Category   global.MerchantCategory `bson:"category"` // 品类
	MerchantId primitive.ObjectID      `bson:"merchantId"`
	PicURL     string                  `bson:"picUrl"`
	TagIds     []primitive.ObjectID    `bson:"tagIds"`
	SPIds      []primitive.ObjectID    `bson:"spIds"` // 规格和价格id
	Online     global.CommodityStatus  `bson:"online"`
	Show       global.CommodityStatus  `bson:"show"`
	Weight     int                     `bson:"weight"` // 权重
	Count      int                     `bson:"count"`  // 商品数量
	CreatedAt  time.Time               `bson:"createdAt"`
}

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
	Online   *global.CommodityStatus
	Show     *global.CommodityStatus
}

func (this CommodityUpdateDoc) Update() M {
	m := M{}
	if this.Name != nil {
		m["name"] = *this.Name
	}
	if this.Category != nil {
		m["category"] = *this.Category
	}
	if this.PicURL != nil {
		m["picUrl"] = *this.PicURL
	}
	if this.TagIds != nil {
		m["tagIds"] = this.TagIds
	}
	if this.Online != nil {
		m["online"] = *this.Online
	}
	if this.Online != nil {
		m["show"] = *this.Show
	}

	if len(m) == 0 {
		return nil
	}

	return M{"$set": m}
}

func (this CommodityMod) UpdateById(c ctx.Context, id primitive.ObjectID, doc *CommodityUpdateDoc) error {
	update := doc.Update()
	if len(update) == 0 {
		return nil
	}

	err := this.Collection().UpdateId(c, id, update)
	return err
}

func (this CommodityMod) Page(c context.ContextB, term *CommodityTerm, page *common.Page) (list []*CommodityMod, hasNext bool, err error) {
	return this.page(c, term, page.PageNumber, page.PageSize, "-weight")
}

type CommodityTerm struct {
	Ids        []primitive.ObjectID
	MerchantId *primitive.ObjectID
	TagIds     []primitive.ObjectID
	Online     *global.CommodityStatus
	Show       *global.CommodityStatus
}

func (this CommodityTerm) Filter() M {
	filter := M{}

	if len(this.Ids) != 0 {
		filter["_id"] = M{
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
	if this.Online != nil {
		filter["online"] = *this.Online
	}
	if this.Online != nil {
		filter["show"] = *this.Show
	}

	return filter
}

func (this CommodityMod) page(c context.ContextB, term *CommodityTerm, pageNumber, pageSize int64, sortBy string) (list []*CommodityMod, hasNext bool, err error) {

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

func (this CommodityMod) Find(c context.ContextB, term *CommodityTerm) (list []*CommodityMod, err error) {
	filter := term.Filter()
	err = this.Collection().Find(c, filter).All(&list)
	return
}

func (this CommodityMod) FindM(c context.ContextB, term *CommodityTerm) (map[primitive.ObjectID]*CommodityMod, error) {
	list, err := this.Find(c, term)
	if err != nil {
		return nil, err
	}

	m := map[primitive.ObjectID]*CommodityMod{}
	for _, v := range list {
		m[v.ID] = v
	}

	return m, nil
}

func (this CommodityMod) FindShowByIDs(c context.ContextB, ids ...primitive.ObjectID) (list []*CommodityMod, err error) {
	if len(ids) == 0 {
		return []*CommodityMod{}, nil
	}

	status := global.Show
	term := &CommodityTerm{Ids: ids, Show: &status}
	err = this.Collection().Find(c, term).All(&list)

	return
}

func (this CommodityMod) FindShowyIDm(c context.ContextB, ids ...primitive.ObjectID) (map[primitive.ObjectID]*CommodityMod, error) {
	list, err := this.FindShowByIDs(c, ids...)
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
