package orderser

import (
	"hx/global/context"
	"hx/mdb"
	"hx/model/usermod"
	"hx/util"
)

var (
	Order OrderServer
)

type OrderServer struct {
}

func (this OrderServer) Info(c context.ContextB, r usermod.OrderInfoRequest) (*usermod.OrderInfoResponse, error) {
	commoditys, err := mdb.Commodity.FindOnlineByIDs(c, r.CommodityID)
	if err != nil {
		c.Errorf("mdb.Commodity.FindByIDs failed! err: %v", err)
		return nil, err
	}

	details := []*usermod.CommodityDetails{}
	for _, v := range commoditys {
		ts, err := mdb.Tag.FindByIDs(c, v.TagIds)
		if err != nil {
			c.Errorf("mdb.Tag.FindByIDs failed! err: %v", err)
			continue
		}

		tagContents := []string{}
		selects := []*usermod.Select{}

		for _, t := range ts {
			switch t.Type {
			case mdb.Server:
				sel := &usermod.Select{
					CommodityID: v.ID,
					SelectID:    t.ID,
					SelectName:  t.Name,
					SelectPrice: t.Value,
					MD5:         t.GenMD5(),
				}
				selects = append(selects, sel)
			case mdb.Age:
				tagContents = append(tagContents, t.Value)
			case mdb.Nationality:
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
		invaild := v.Status != mdb.Online
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

	onlineCommodityM, err := mdb.Commodity.FindOnlineByIDm(c, commodityIDs)
	if err != nil {
		c.Errorf("mdb.Commodity.FindOnlineByIDm failed! err: %v", err)
		return nil, err
	}

	tagM, err := mdb.Tag.FindByIDm(c, selTagIds)
	if err != nil {
		c.Errorf("mdb.Tag.FindByIDm failed! err: %v", err)
		return nil, err
	}

	var cmdyInvaild bool
	for _, id := range commodityIDs {
		commodity := onlineCommodityM[id]
		if commodity == nil {
			cmdyInvaild = true
			break
		}

		if commodity.Status != mdb.Online {
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
