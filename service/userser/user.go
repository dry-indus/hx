package userser

import (
	"hx/global/context"
	"hx/model/common"
	"hx/model/db"
	"hx/model/usermod"
	"hx/util"
)

var (
	Home  HomeServer
	Order OrderServer
)

type HomeServer struct {
	common.Logger
}

func (this HomeServer) List(c context.UserContext, r usermod.HomeListRequest) (*usermod.HomeListResponse, error) {
	list, hasNext, err := db.Commodity.FindOnline(c, c.Merchant().ID, nil, r.Page)
	if err != nil {
		this.Errorf("db.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	commoditys := []*usermod.Commodity{}
	for _, v := range list {

		ts, err := db.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			this.Errorf("db.Tag.FindByIDs failed! err: %v", err)
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
	list, hasNext, err := db.Commodity.FindOnline(c, c.Merchant().ID, r.TagIDs, r.Page)
	if err != nil {
		this.Errorf("db.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	commoditys := []*usermod.Commodity{}
	for _, v := range list {
		ts, err := db.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			this.Errorf("db.Tag.FindByIDs failed! err: %v", err)
			continue
		}

		tags := []*usermod.Tags{}
		for _, t := range ts {
			if t.Type == db.Server {
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

type OrderServer struct {
	common.Logger
}

func (this OrderServer) Info(c context.ContextB, r usermod.OrderInfoRequest) (*usermod.OrderInfoResponse, error) {
	commoditys, err := db.Commodity.FindOnlineByIDs(c, r.CommodityID)
	if err != nil {
		this.Errorf("db.Commodity.FindByIDs failed! err: %v", err)
		return nil, err
	}

	details := []*usermod.CommodityDetails{}
	for _, v := range commoditys {
		ts, err := db.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			this.Errorf("db.Tag.FindByIDs failed! err: %v", err)
			continue
		}

		tagContents := []string{}
		selects := []*usermod.Select{}

		for _, t := range ts {
			switch t.Type {
			case db.Server:
				sel := &usermod.Select{
					CommodityID: v.ID,
					SelectID:    t.ID,
					SelectName:  t.Name,
					SelectPrice: t.Value,
					MD5:         t.GenMD5(),
				}
				selects = append(selects, sel)
			case db.Age:
				tagContents = append(tagContents, t.Value)
			case db.Nationality:
				tagContents = append(tagContents, t.Value)
			default:
				if len(t.Value) != 0 {
					tagContents = append(tagContents, t.Value)
				} else {
					tagContents = append(tagContents, t.Name)
				}
			}
		}

		if len(tagContents) > 6 {
			tagContents = tagContents[:6]
		}

		var invaildMsg string
		invaild := v.Status != db.Online
		if invaild {
			invaildMsg = "商品走丢了！"
		}

		d := &usermod.CommodityDetails{
			CommodityID: v.ID,
			PicURL:      v.PicURL,
			TagNames:    tagContents,
			Selects:     selects,
			Invaild:     invaild,
			InvaildMsg:  invaildMsg,
		}
		details = append(details, d)
	}

	resp := &usermod.OrderInfoResponse{
		Details: details,
	}

	return resp, nil
}

func (this OrderServer) Submit(c context.ContextB, r usermod.SubmitOrderRequest) (*usermod.SubmitOrderResponse, error) {
	selTagIds := []string{}
	commodityIDs := []string{}
	selectM := map[string]*usermod.Select{}

	for _, v := range r.Selects {
		commodityIDs = append(commodityIDs, v.CommodityID)
		selTagIds = append(selTagIds, v.SelectID)
		selectM[v.SelectID] = v
	}

	onlineCommodityM, err := db.Commodity.FindOnlineByIDm(c, commodityIDs)
	if err != nil {
		this.Errorf("db.Commodity.FindOnlineByIDm failed! err: %v", err)
		return nil, err
	}

	tagM, err := db.Tag.FindByIDm(c, selTagIds)
	if err != nil {
		this.Errorf("db.Tag.FindByIDm failed! err: %v", err)
		return nil, err
	}

	var cmdyInvaild bool
	for _, id := range commodityIDs {
		commodity := onlineCommodityM[id]
		if commodity == nil {
			cmdyInvaild = true
			break
		}

		if commodity.Status != db.Online {
			cmdyInvaild = true
			break
		}
	}

	var tagInvaild bool
	for _, id := range selTagIds {
		tag := tagM[id]
		sel := selectM[id]
		if tag == nil || sel == nil {
			tagInvaild = true
			break
		}
		if sel.MD5 != tag.GenMD5() {
			tagInvaild = true
			break
		}
	}

	resp := &usermod.SubmitOrderResponse{
		OrderId:        util.UUID().String(),
		Invaild:        cmdyInvaild || tagInvaild,
		OrderPicBuffer: GenOrderPicture(),
		JumpUrl:        GenJumpUrl(),
	}

	return resp, nil
}

func GenOrderPicture() string {
	return ""
}

func GenJumpUrl() string {
	return ""
}
