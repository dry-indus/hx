package orderser

import (
	"hx/global/context"
	"hx/mdb"
	"hx/model/usermod"
	"hx/util"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Order OrderServer
)

type OrderServer struct{}

func (this OrderServer) Submit(c context.UserContext, r usermod.SubmitOrderRequest) (*usermod.SubmitOrderResponse, error) {
	invaild, invaildMsg, err := this.checkCommoditys(c, r)
	if err != nil {
		return nil, err
	}

	var orderPicBuffer string
	var jumpUrl string
	if !invaild {
		orderPicBuffer = this.GenOrderPicture(c)
		jumpUrl = this.GenJumpUrl(c)
	}

	resp := &usermod.SubmitOrderResponse{
		OrderId:        util.UUID().String(),
		Invaild:        invaild,
		InvaildMsg:     invaildMsg,
		OrderPicBuffer: orderPicBuffer,
		JumpUrl:        jumpUrl,
		Commoditys:     r.Commoditys,
		TotalPrice:     r.TotalPrice,
	}

	return resp, nil
}

func (this OrderServer) checkCommoditys(c context.UserContext, r usermod.SubmitOrderRequest) (invaild bool, invaildMsg string, err error) {
	commodityIDs := []primitive.ObjectID{}
	for _, v := range r.Commoditys {
		commodityIDs = append(commodityIDs, v.ID)
	}

	commodityM, err := mdb.Commodity.FindM(c, &mdb.CommodityTerm{Ids: commodityIDs})
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return
	}

	totalPrice := decimal.Zero
	for _, v := range r.Commoditys {
		cmdy := commodityM[v.ID]
		if cmdy == nil {
			v.Invaild = true
			v.InvaildMsg = "商品不存在"
			continue
		}

		if cmdy.Count == 0 {
			v.Invaild = true
			v.InvaildMsg = "商品已售罄"
			continue
		}

		for _, s := range v.SPs {
			if !s.Selected {
				continue
			}

			sp, _ := mdb.SpecificationsPricing.FindById(c, s.ID)
			if sp == nil {
				v.Invaild = true
				v.InvaildMsg = "选项已失效"
				continue
			}

			if s.MD5 != mdb.SPMD5(sp) {
				v.Invaild = true
				v.InvaildMsg = "选项已失效"
				continue
			}

			if !s.TotalPricing.Equal(s.Pricing.Mul(s.BuyCount)) {
				v.Invaild = true
				v.InvaildMsg = "计价错误"
				continue
			}

			totalPrice = totalPrice.Add(s.TotalPricing)
		}
	}

	if !totalPrice.Equal(r.TotalPrice) {
		invaild = true
		invaildMsg = "计价错误"
	}

	for _, v := range r.Commoditys {
		if v.Invaild {
			invaild = true
			invaildMsg = v.InvaildMsg
			break
		}
	}

	return
}

func (this OrderServer) GenOrderPicture(c context.UserContext) string {
	return ""
}

func (this OrderServer) GenJumpUrl(c context.UserContext) string {
	return c.Merchant().TgName
}
