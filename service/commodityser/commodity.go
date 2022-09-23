package commodityser

import (
	ctx "context"
	"hx/global"
	"hx/global/context"
	"hx/mdb"
	"hx/model/merchantmod"
	"time"

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

	term := &mdb.CommodityPageTerm{MerchantId: &c.Merchant().ID}
	cs, hasNext, err := mdb.Commodity.Page(c, term, r.Page)
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	commoditys := []*merchantmod.Commodity{}

	for _, v := range cs {
		tags := []*merchantmod.Tag{}
		for _, tagId := range v.TagIds {
			tagMod := tagM[tagId]
			if tagMod == nil {
				continue
			}
			tag := &merchantmod.Tag{
				ID:   tagMod.ID,
				Name: tagMod.Name,
			}
			tags = append(tags, tag)
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
			ID:     v.ID,
			Name:   v.Name,
			PicURL: v.PicURL,
			Tags:   tags,
			SPs:    sps,
		}
		commoditys = append(commoditys, commodity)
	}

	resp := &merchantmod.CommodityListResponse{
		List:    commoditys,
		AllTags: tags,
		HasNext: hasNext,
	}

	return resp, nil
}

func (this Commodityser) Add(c context.MerchantContext, merchantId primitive.ObjectID, r *merchantmod.CommodityAddRequest) (*merchantmod.CommodityAddResponse, error) {
	var count int
	for _, v := range r.Commoditys {
		commodityId, err := this.AddOne(c, v)
		c.Infof("AddOne finish! commodityId: %s, err: %v", commodityId, err)
		if !commodityId.IsZero() {
			count++
		}
	}

	resp := &merchantmod.CommodityAddResponse{
		Count: count,
	}

	return resp, nil
}

func (Commodityser) AddOne(c context.MerchantContext, add *merchantmod.CommodityAdd) (primitive.ObjectID, error) {
	s, err := global.DL_CORE_CLI.Session()
	if err != nil {
		c.Errorf("AddOne create sessio* failed! err: %v", err)
		return primitive.ObjectID{}, err
	}
	defer s.EndSession(c)

	callback := func(sessCtx ctx.Context) (interface{}, error) {
		now := time.Now()

		tagIds := []primitive.ObjectID{}
		for _, v := range add.Tags {
			if v.Selected {
				tagIds = append(tagIds, v.ID)
			}
		}

		commodityId, err := mdb.Commodity.Add(sessCtx, mdb.CommodityMod{
			Name:       add.Name,
			Category:   c.Merchant().Category,
			MerchantId: c.Merchant().ID,
			PicURL:     add.PicURL,
			TagIds:     tagIds,
			Status:     mdb.Hide,
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

		return commodityId, nil
	}

	commodityId, err := s.StartTransaction(c, callback)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return commodityId.(primitive.ObjectID), nil
}

func (Commodityser) Modify(c context.MerchantContext, r *merchantmod.CommodityModifyRequest) (*merchantmod.CommodityModifyResponse, error) {
	tagIds := []primitive.ObjectID{}
	for _, v := range r.Tags {
		if v.Selected {
			tagIds = append(tagIds, v.ID)
		}
	}

	err := mdb.Commodity.UpdateById(c, r.Id, &mdb.CommodityUpdateDoc{
		Name:     r.Name,
		Category: &c.Merchant().Category,
		PicURL:   r.PicURL,
		TagIds:   tagIds,
	})
	if err != nil {
		c.Errorf("Commodity UpdateById failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.CommodityModifyResponse{}

	return resp, nil
}

func (Commodityser) Del(c context.ContextB, r *merchantmod.CommodityDelRequest) (*merchantmod.CommodityDelResponse, error) {
	err := mdb.Commodity.Del(c, r.Id)
	if err != nil {
		c.Errorf("Commodity Del failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.CommodityDelResponse{}

	return resp, nil
}

func (this Commodityser) Publish(c context.ContextB, r *merchantmod.CommodityPublishRequest) (*merchantmod.CommodityPublishResponse, error) {
	err := this.UpdateStatus(c, r.Id, mdb.Show)
	if err != nil {
		return nil, err
	}

	resp := &merchantmod.CommodityPublishResponse{}
	return resp, nil
}

func (this Commodityser) Hide(c context.ContextB, r *merchantmod.CommodityHideRequest) (*merchantmod.CommodityHideResponse, error) {
	err := this.UpdateStatus(c, r.Id, mdb.Hide)
	if err != nil {
		return nil, err
	}

	resp := &merchantmod.CommodityHideResponse{}
	return resp, nil
}

func (Commodityser) UpdateStatus(c context.ContextB, id primitive.ObjectID, status mdb.CommodityStatus) error {
	err := mdb.Commodity.UpdateById(c, id, &mdb.CommodityUpdateDoc{
		Status: &status,
	})
	if err != nil {
		c.Errorf("Commodity UpdateById failed! err: %v", err)
		return err
	}

	return nil
}
