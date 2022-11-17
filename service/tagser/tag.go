package tagser

import (
	ctx "context"
	"hx/global"
	"hx/global/context"
	"hx/mdb"
	"hx/model/merchantmod"
	"time"

	gosonic "github.com/expectedsh/go-sonic/sonic"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Tag TagSer
)

type TagSer struct {
}

func (this TagSer) Add(c context.MerchantContext, merchantId primitive.ObjectID, r *merchantmod.TagAddRequest) (*merchantmod.TagAddResponse, error) {
	mod := &mdb.TagMod{
		Name:       r.Name,
		MerchantId: merchantId,
		CreatedAt:  time.Now(),
	}

	id, err := mdb.Tag.AddOne(c, mod)
	if err != nil {
		c.Errorf("Tag AddOne failed! err: %v", err)
		return nil, err
	}

	event := &global.SonicIngestEvent{
		Method:     global.Push,
		Collection: "tag",
		Bucket:     c.Merchant().ID.String(),
		Records:    []gosonic.IngestBulkRecord{{Text: r.Name, Object: id.String()}},
		Lang:       c.Lang(),
		Trace:      c.Trace(),
	}

	go func() { global.SONIC_INGESTER_CH <- event }()

	resp := &merchantmod.TagAddResponse{Id: id}

	return resp, nil
}

func (TagSer) Del(c context.MerchantContext, r *merchantmod.TagDelRequest) (*merchantmod.TagDelResponse, error) {
	tag, _ := mdb.Tag.FindByID(c, r.Id)
	if tag == nil {
		return &merchantmod.TagDelResponse{Id: r.Id}, nil
	}

	s, err := global.DL_CORE_CLI.Session()
	if err != nil {
		c.Errorf("Del create session* failed! err: %v", err)
		return nil, err
	}
	defer s.EndSession(c)

	callback := func(sessCtx ctx.Context) (interface{}, error) {
		err = mdb.Tag.DelByID(sessCtx, r.Id)
		if err != nil {
			c.Errorf("Tag.Del failed! err: %v", err)
			return nil, err
		}

		err = mdb.CommodityTag.DelByTerm(sessCtx, &mdb.CommodityTagTerm{TagId: &r.Id})
		if err != nil {
			c.Errorf("CommodityTag.DelByTerm failed! err: %v", err)
			return nil, err
		}
		return r.Id, nil
	}

	tagId, err := s.StartTransaction(c, callback)
	if err != nil {
		return nil, err
	}

	event := &global.SonicIngestEvent{
		Method:     global.Pop,
		Collection: "tag",
		Bucket:     c.Merchant().ID.String(),
		Records:    []gosonic.IngestBulkRecord{{Text: tag.Name, Object: tag.ID.String()}},
		Lang:       c.Lang(),
		Trace:      c.Trace(),
	}

	go func() { global.SONIC_INGESTER_CH <- event }()

	resp := &merchantmod.TagDelResponse{Id: tagId.(primitive.ObjectID)}

	return resp, nil
}

func (TagSer) Stat(c context.ContextB, r *merchantmod.TagStatRequest) (*merchantmod.TagStatResponse, error) {
	tag, _ := mdb.Tag.FindByID(c, r.Id)
	truely := true
	cts, _ := mdb.CommodityTag.FindByTerm(c, &mdb.CommodityTagTerm{TagId: &r.Id, Show: &truely, Enable: &truely})
	commoditys, _ := mdb.Commodity.Find(c, &mdb.CommodityTerm{Ids: cts.CommodityIds()})

	cs := []*merchantmod.Commodity{}
	for _, v := range commoditys {
		c := &merchantmod.Commodity{
			ID:     v.ID,
			Name:   v.Name,
			PicURL: v.PicURL,
		}
		cs = append(cs, c)
	}

	resp := &merchantmod.TagStatResponse{
		Tag: &merchantmod.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		},
		Commoditys: cs,
	}

	return resp, nil
}
