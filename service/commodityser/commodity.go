package commodityser

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
	Commodity Commodityser
)

type Commodityser struct{}

func (this Commodityser) List(c context.UserContext, r *merchantmod.CommodityListRequest) (*merchantmod.CommodityListResponse, error) {
	allTags, err := mdb.Tag.FindByMerchantId(c, c.Merchant().ID)
	if err != nil {
		c.Errorf("mdb.Tag.FindByMerchantId failed! err: %v", err)
		return nil, err
	}

	tags := []*merchantmod.Tag{}
	tagM := map[primitive.ObjectID]*merchantmod.Tag{}
	for _, v := range allTags {
		tag := &merchantmod.Tag{
			ID:   v.ID,
			Name: v.Name,
		}
		tagM[v.ID] = tag
		tags = append(tags, tag)
	}

	term := &mdb.CommodityTerm{MerchantId: &c.Merchant().ID}
	commoditys, hasNext, err := mdb.Commodity.Page(c, term, &r.Page)
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	cs := []*merchantmod.Commodity{}
	for _, v := range commoditys {
		truely := true
		cts, _ := mdb.CommodityTag.FindByTerm(c, &mdb.CommodityTagTerm{CommodityId: &v.ID, Show: &truely, Enable: &truely})
		tags, _ := mdb.Tag.FindByIDs(c, cts.TagIds())
		ts := []*merchantmod.Tag{}
		for _, tag := range tags {
			t := &merchantmod.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
			ts = append(ts, t)
		}

		spmods, err := mdb.SpecificationsPricing.FindByCommodityId(c, v.ID)
		if err != nil {
			c.Errorf("mdb.SpecificationsPricing.FindByCommodityId failed! err: %v", err)
		}
		sps := []*merchantmod.SP{}
		for _, v := range spmods {
			sp := &merchantmod.SP{
				ID:             v.ID,
				Specifications: v.Specifications,
				Pricing:        v.Pricing,
				PicURL:         v.PicURL,
			}
			sps = append(sps, sp)
		}

		commodity := &merchantmod.Commodity{
			ID:        v.ID,
			Name:      v.Name,
			PicURL:    v.PicURL,
			Tags:      ts,
			SPs:       sps,
			Category:  v.Category,
			Show:      v.Show,
			Online:    v.Online,
			Weight:    v.Weight,
			Count:     v.Count,
			CreatedAt: v.CreatedAt,
		}
		cs = append(cs, commodity)
	}

	resp := &merchantmod.CommodityListResponse{
		List:    cs,
		AllTags: tags,
		HasNext: hasNext,
	}

	return resp, nil
}

func (this Commodityser) Add(c context.MerchantContext, merchantId primitive.ObjectID, r *merchantmod.CommodityAddRequest) (*merchantmod.CommodityAddResponse, error) {
	var count int
	commodityIds := []primitive.ObjectID{}
	for _, v := range r.Commoditys {
		commodityId, err := this.AddOne(c, v)
		c.Debugf("AddOne finish! commodityId: %s, err: %v", commodityId, err)
		if !commodityId.IsZero() {
			count++
			commodityIds = append(commodityIds, commodityId)
		}
	}

	resp := &merchantmod.CommodityAddResponse{
		Count: count,
		Ids:   commodityIds,
	}

	return resp, nil
}

func (Commodityser) AddOne(c context.MerchantContext, add *merchantmod.CommodityAdd) (primitive.ObjectID, error) {
	s, err := global.DL_CORE_CLI.Session()
	if err != nil {
		c.Errorf("AddOne create session* failed! err: %v", err)
		return primitive.NilObjectID, err
	}
	defer s.EndSession(c)

	var pushEvent *global.SonicIngestEvent

	callback := func(sessCtx ctx.Context) (interface{}, error) {
		now := time.Now()

		commodityId, err := mdb.Commodity.Add(sessCtx, mdb.CommodityMod{
			Name:       add.Name,
			Category:   c.Merchant().Category,
			MerchantId: c.Merchant().ID,
			PicURL:     add.PicURL,
			Show:       global.Show,
			Online:     global.Online,
			CreatedAt:  now,
		})
		if err != nil {
			c.Errorf("Commodity.Add failed! err: %v", err)
			return nil, err
		}

		spMods := []*mdb.SpecificationsPricingMod{}
		for _, sp := range add.SPs {
			spMod := mdb.SpecificationsPricingMod{
				CommodityId:    commodityId,
				Specifications: sp.Specifications,
				Pricing:        sp.Pricing,
			}
			spMods = append(spMods, &spMod)
		}

		_, err = mdb.SpecificationsPricing.AddMany(sessCtx, spMods)
		if err != nil {
			c.Errorf("SpecificationsPricing.AddMany failed! err: %v", err)
			return nil, err
		}

		addNameTag := &mdb.TagMod{
			MerchantId: c.Merchant().ID,
			Name:       add.Name,
			CreatedAt:  now,
		}

		// 商品名存为不显示的标签
		upsertTagId, err := mdb.Tag.UpsertOne(sessCtx, &mdb.TagTerm{MerchantId: &c.Merchant().ID, Name: &add.Name}, addNameTag)
		if err != nil {
			c.Errorf("Tag.CountByTerm failed! err: %v", err)
			return nil, err
		}

		// 关联商品和标签
		ctms := []*mdb.CommodityTagMod{{MerchantId: c.Merchant().ID, CommodityId: commodityId, TagId: upsertTagId, Enable: true, Show: false, CreatedAt: now}}
		for _, v := range add.Tags {
			ctm := &mdb.CommodityTagMod{
				MerchantId:  c.Merchant().ID,
				CommodityId: commodityId,
				TagId:       v.ID,
				Enable:      v.Selected,
				Show:        true,
				CreatedAt:   now,
			}
			ctms = append(ctms, ctm)
		}

		_, err = mdb.CommodityTag.AddMany(sessCtx, ctms)
		if err != nil {
			c.Errorf("CommodityTag.AddMany failed! err: %v", err)
			return nil, err
		}

		pushEvent = &global.SonicIngestEvent{
			Collection: "tag",
			Bucket:     c.Merchant().ID.String(),
			Records:    []gosonic.IngestBulkRecord{{Text: add.Name, Object: upsertTagId.String()}},
			Lang:       c.Lang(),
			Trace:      c.Trace(),
		}

		return commodityId, nil
	}

	commodityId, err := s.StartTransaction(c, callback)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if pushEvent != nil {
		go func() { global.SONIC_INGESTER_CH <- pushEvent }()
	}

	return commodityId.(primitive.ObjectID), nil
}

func (Commodityser) Modify(c context.MerchantContext, r *merchantmod.CommodityModifyRequest) (*merchantmod.CommodityModifyResponse, error) {
	oldCommodity, err := mdb.Commodity.FindOneById(c, r.Id)
	if err != nil {
		c.Errorf("Commodity.FindOneById failed! id: %s, err: %v", r.Id, err)
		return nil, err
	}

	s, err := global.DL_CORE_CLI.Session()
	if err != nil {
		c.Errorf("Modify create session failed! err: %v", err)
		return nil, err
	}
	defer s.EndSession(c)

	var popEvent *global.SonicIngestEvent
	var pushEvent *global.SonicIngestEvent

	callback := func(sessCtx ctx.Context) (interface{}, error) {
		now := time.Now()
		err := mdb.Commodity.UpdateById(c, r.Id, &mdb.CommodityUpdateDoc{
			Name:     r.Name,
			Category: &c.Merchant().Category,
			PicURL:   r.PicURL,
		})
		if err != nil {
			c.Errorf("Commodity.UpdateById failed! err: %v", err)
			return nil, err
		}

		if r.Name != nil && oldCommodity.Name != *r.Name {

			modifyNameTag := &mdb.TagMod{
				MerchantId: c.Merchant().ID,
				Name:       *r.Name,
				CreatedAt:  now,
			}

			// 商品名存为不显示的标签
			upsertTagId, err := mdb.Tag.UpsertOne(sessCtx, &mdb.TagTerm{MerchantId: &c.Merchant().ID, Name: &oldCommodity.Name}, modifyNameTag)
			if err != nil {
				c.Errorf("Tag.UpsertOne failed! err: %v", err)
				return nil, err
			}

			ctm := &mdb.CommodityTagMod{
				MerchantId:  c.Merchant().ID,
				CommodityId: r.Id,
				TagId:       upsertTagId,
				Enable:      true,
				Show:        false,
				CreatedAt:   now,
			}

			_, err = mdb.CommodityTag.UpsertOne(sessCtx, &mdb.CommodityTagTerm{MerchantId: &c.Merchant().ID, CommodityId: &r.Id, TagId: &upsertTagId}, ctm)
			if err != nil {
				c.Errorf("CommodityTag.UpsertOne failed! err: %v", err)
				return nil, err
			}

			popEvent = &global.SonicIngestEvent{
				Method:     global.Pop,
				Collection: "tag",
				Bucket:     c.Merchant().ID.String(),
				Records:    []gosonic.IngestBulkRecord{{Text: oldCommodity.Name, Object: upsertTagId.String()}},
				Lang:       c.Lang(),
				Trace:      c.Trace(),
			}

			pushEvent = &global.SonicIngestEvent{
				Method:     global.Push,
				Collection: "tag",
				Bucket:     c.Merchant().ID.String(),
				Records:    []gosonic.IngestBulkRecord{{Text: *r.Name, Object: upsertTagId.String()}},
				Lang:       c.Lang(),
				Trace:      c.Trace(),
			}
		}

		return nil, nil
	}

	_, err = s.StartTransaction(c, callback)
	if err != nil {
		return nil, err
	}

	if pushEvent != nil {
		go func() { global.SONIC_INGESTER_CH <- pushEvent }()
	}

	if popEvent != nil {
		go func() { global.SONIC_INGESTER_CH <- popEvent }()
	}

	resp := &merchantmod.CommodityModifyResponse{r.Id}

	return resp, nil
}

func (Commodityser) Del(c context.MerchantContext, r *merchantmod.CommodityDelRequest) (*merchantmod.CommodityDelResponse, error) {
	s, err := global.DL_CORE_CLI.Session()
	if err != nil {
		c.Errorf("Del create session failed! err: %v", err)
		return nil, err
	}
	defer s.EndSession(c)

	callback := func(sessCtx ctx.Context) (interface{}, error) {
		err := mdb.Commodity.Del(sessCtx, r.Id)
		if err != nil {
			c.Errorf("Commodity.Del failed! err: %v", err)
			return nil, err
		}

		err = mdb.CommodityTag.DelByTerm(sessCtx, &mdb.CommodityTagTerm{MerchantId: &c.Merchant().ID, CommodityId: &r.Id})
		if err != nil {
			c.Errorf("CommodityTag.DelByTerm failed! err: %v", err)
			return nil, err
		}

		return nil, nil
	}

	_, err = s.StartTransaction(c, callback)
	if err != nil {
		return nil, err
	}

	resp := &merchantmod.CommodityDelResponse{Id: r.Id}

	return resp, nil
}

func (this Commodityser) Publish(c context.ContextB, r *merchantmod.CommodityPublishRequest) (*merchantmod.CommodityPublishResponse, error) {
	err := this.UpdateShow(c, r.Id, global.Show)
	if err != nil {
		return nil, err
	}

	resp := &merchantmod.CommodityPublishResponse{Id: r.Id}
	return resp, nil
}

func (this Commodityser) Hide(c context.ContextB, r *merchantmod.CommodityHideRequest) (*merchantmod.CommodityHideResponse, error) {
	err := this.UpdateShow(c, r.Id, global.Hide)
	if err != nil {
		return nil, err
	}

	resp := &merchantmod.CommodityHideResponse{Id: r.Id}
	return resp, nil
}

func (Commodityser) UpdateShow(c context.ContextB, id primitive.ObjectID, status global.CommodityStatus) error {
	err := mdb.Commodity.UpdateById(c, id, &mdb.CommodityUpdateDoc{
		Show: &status,
	})
	if err != nil {
		c.Errorf("Commodity UpdateById failed! err: %v", err)
		return err
	}

	return nil
}
