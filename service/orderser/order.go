package orderser

import (
	"hx/global"
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
	commodityIDs := []primitive.ObjectID{}
	for _, v := range r.Commoditys {
		commodityIDs = append(commodityIDs, v.ID)
	}

	commodityM, err := mdb.Commodity.FindM(c, &mdb.CommodityTerm{Ids: commodityIDs})
	if err != nil {
		c.Errorf("mdb.Commodity.Find failed! err: %v", err)
		return nil, err
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

		singleChoiceSpM := map[primitive.ObjectID]*usermod.SP{}
		multipleChoiceSpM := map[primitive.ObjectID]*usermod.SP{}
		mustChoiceSpM := map[primitive.ObjectID]*usermod.SP{}
		for _, s := range v.SPs {
			switch s.ChoiceOpt {
			case global.SingleChoice:
				singleChoiceSpM[s.ID] = s
			case global.MultipleChoice:
				multipleChoiceSpM[s.ID] = s
			case global.MustChoice:
				mustChoiceSpM[s.ID] = s
			}
		}

		for _, s := range v.SPs {
			if s.ChoiceOpt != global.MustChoice && !s.Selected {
				continue
			}

			if s.ChoiceOpt == global.MustChoice && !s.Selected {
				v.Invaild = true
				v.InvaildMsg = "必选此项"
				continue
			}

			if s.ChoiceOpt == global.SingleChoice && len(multipleChoiceSpM) != 0 {
				v.Invaild = true
				v.InvaildMsg = "仅支持单选"
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

	var cmdyInvaild bool
	var invaildMsg string
	if !totalPrice.Equal(r.TotalPrice) {
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
		orderPicBuffer = GenOrderPicture(c)
		jumpUrl = GenJumpUrl(c)
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

func GenOrderPicture(c context.UserContext) string {
	return ""
}

func GenJumpUrl(c context.UserContext) string {
	return c.Merchant().Telegram
}
