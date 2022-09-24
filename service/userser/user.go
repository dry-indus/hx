package userser

import (
	"hx/global"
	"hx/global/context"
	"hx/mdb"
	"hx/model/common"
	"hx/model/usermod"
)

var (
	Home HomeServer
)

type HomeServer struct {
}

func (this HomeServer) List(c context.UserContext, r usermod.HomeListRequest) (*usermod.HomeListResponse, error) {
	status := global.Show
	term := &mdb.CommodityPageTerm{MerchantId: &c.Merchant().ID, Show: &status}
	commoditys, hasNext, err := this.search(c, term, &r.Page)
	if err != nil {
		c.Errorf("search failed! err: %v", err)
		return nil, err
	}

	resp := &usermod.HomeListResponse{
		List:    commoditys,
		HasNext: hasNext,
	}

	return resp, nil
}

func (this HomeServer) Search(c context.UserContext, r usermod.HomeSearchRequest) (*usermod.HomeSearchResponse, error) {
	status := global.Show
	term := &mdb.CommodityPageTerm{MerchantId: &c.Merchant().ID, Ids: r.CommodityIDs, TagIds: r.TagIDs, Show: &status}
	commoditys, hasNext, err := this.search(c, term, r.Page)
	if err != nil {
		c.Errorf("search failed! err: %v", err)
		return nil, err
	}

	resp := &usermod.HomeSearchResponse{
		List:    commoditys,
		HasNext: hasNext,
	}

	return resp, nil
}

func (this HomeServer) search(c context.UserContext, term *mdb.CommodityPageTerm, page *common.Page) (commoditys []*usermod.Commodity, hasNext bool, err error) {
	list, hasNext, err := mdb.Commodity.Page(c, term, page)
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return
	}

	for _, v := range list {
		ts, err := mdb.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			c.Errorf("mdb.Tag.FindByIDs failed! err: %v", err)
			continue
		}

		tags := []*usermod.Tag{}
		for _, t := range ts {
			tag := &usermod.Tag{
				ID:   t.ID,
				Name: t.Name,
			}
			tags = append(tags, tag)
		}

		spmods, err := mdb.SpecificationsPricing.FindByCommodityId(c, v.ID)
		if err != nil {
			c.Errorf("mdb.SpecificationsPricing.FindByCommodityId failed! err: %v", err)
			continue
		}

		sps := []*usermod.SP{}
		for _, s := range spmods {
			sp := &usermod.SP{
				ID:             s.ID,
				Specifications: s.Specifications,
				Pricing:        s.Pricing,
				MD5:            mdb.SPMD5(s),
			}
			sps = append(sps, sp)
		}

		commodity := &usermod.Commodity{
			ID:     v.ID,
			Name:   v.Name,
			PicURL: v.PicURL,
			Tags:   tags,
			SPs:    sps,
		}

		commoditys = append(commoditys, commodity)
	}

	return commoditys, hasNext, nil
}
