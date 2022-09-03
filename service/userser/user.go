package userser

import (
	"hx/global/context"
	"hx/mdb"
	"hx/model/usermod"
)

var (
	Home HomeServer
)

type HomeServer struct {
}

func (this HomeServer) List(c context.UserContext, r usermod.HomeListRequest) (*usermod.HomeListResponse, error) {
	list, hasNext, err := mdb.Commodity.FindOnline(c, c.Merchant().ID, nil, r.Page)
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	commoditys := []*usermod.Commodity{}
	for _, v := range list {

		ts, err := mdb.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			c.Errorf("mdb.Tag.FindByIDs failed! err: %v", err)
			continue
		}

		tags := []*usermod.Tags{}
		for _, t := range ts {
			tag := &usermod.Tags{
				TagID:   t.ID,
				TagName: t.Name,
			}
			tags = append(tags, tag)
		}

		if len(tags) > 10 {
			tags = tags[:10]
		}

		commodity := &usermod.Commodity{
			CommodityID: v.ID,
			PicURL:      v.PicURL,
			Tags:        tags,
		}

		commoditys = append(commoditys, commodity)
	}

	resp := &usermod.HomeListResponse{
		List:    commoditys,
		HasNext: hasNext,
	}

	return resp, nil
}

func (this HomeServer) Search(c context.UserContext, r usermod.HomeSearchRequest) (*usermod.HomeSearchResponse, error) {
	list, hasNext, err := mdb.Commodity.FindOnline(c, c.Merchant().ID, r.TagIDs, r.Page)
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	commoditys := []*usermod.Commodity{}
	for _, v := range list {
		ts, err := mdb.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			c.Errorf("mdb.Tag.FindByIDs failed! err: %v", err)
			continue
		}

		tags := []*usermod.Tags{}
		for _, t := range ts {
			if t.Type == mdb.Server {
				continue
			}

			tag := &usermod.Tags{
				TagID:   t.ID,
				TagName: t.Name,
			}
			tags = append(tags, tag)
		}

		commodity := &usermod.Commodity{
			CommodityID: v.ID,
			PicURL:      v.PicURL,
			Tags:        tags,
		}

		commoditys = append(commoditys, commodity)
	}

	resp := &usermod.HomeSearchResponse{
		List:    commoditys,
		HasNext: hasNext,
	}

	return resp, nil
}
