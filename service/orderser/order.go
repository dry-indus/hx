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

func (this OrderServer) Submit(c context.ContextB, r usermod.SubmitOrderRequest) (*usermod.SubmitOrderResponse, error) {
	commodityIDs := []primitive.ObjectID{}
	for _, v := range r.Commoditys {
		commodityIDs = append(commodityIDs, v.ID)
	}

	commodityM, err := mdb.Commodity.FindM(c, &mdb.CommodityPageTerm{Ids: commodityIDs})
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return nil, err
	}

	totalPrice := decimal.Zero
	for _, v := range r.Commoditys {
		cmdy := commodityM[v.ID]
		if cmdy == nil {
			v.Invaild = false
			v.InvaildMsg = "商品不存在"
		}

		if cmdy.Count == 0 {
			v.Invaild = false
			v.InvaildMsg = "商品已售罄"
		}

		for _, s := range v.SPs {
			if !s.Selected {
				continue
			}

			sp, _ := mdb.SpecificationsPricing.FindById(c, s.ID)
			if sp == nil {
				v.Invaild = false
				v.InvaildMsg = "商品已失效"
			}

			if s.MD5 != mdb.SPMD5(sp) {
				v.Invaild = false
				v.InvaildMsg = "商品已失效"
			}

			totalPrice = totalPrice.Add(s.Pricing)
		}
	}

	var cmdyInvaild bool
	var invaildMsg string
	if totalPrice != r.TotalPrice {
		cmdyInvaild = true
		invaildMsg = "计价错误"
	}

	for _, v := range r.Commoditys {
		if v.Invaild {
			cmdyInvaild = true
			invaildMsg = v.InvaildMsg
			break
		}
	}

	var orderPicBuffer string
	var jumpUrl string
	if !cmdyInvaild {
		orderPicBuffer = GenOrderPicture()
		jumpUrl = GenJumpUrl()
	}

	resp := &usermod.SubmitOrderResponse{
		OrderId:        util.UUID().String(),
		Invaild:        cmdyInvaild,
		InvaildMsg:     invaildMsg,
		OrderPicBuffer: orderPicBuffer,
		JumpUrl:        jumpUrl,
		Commoditys:     r.Commoditys,
		TotalPrice:     totalPrice,
	}

	return resp, nil
}

func GenOrderPicture() string {
	return ""
}

func GenJumpUrl() string {
	return ""
}
