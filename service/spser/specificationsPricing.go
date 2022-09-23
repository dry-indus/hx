package spser

import (
	"fmt"
	"hx/global/context"
	"hx/mdb"
	"hx/model/merchantmod"
	"time"
)

var (
	SP SPSer
)

type SPSer struct{}

func (this SPSer) Add(c context.ContextB, r *merchantmod.SPAddRequest) (*merchantmod.SPAddResponse, error) {
	mod := &mdb.SpecificationsPricingMod{
		CommodityId:    r.CommodityId,
		Specifications: r.Specifications,
		Pricing:        r.Pricing,
		PicURL:         r.PicURL,
		ChoiceOpt:      r.ChoiceOpt,
		CreatedAt:      time.Now(),
	}
	err := mdb.SpecificationsPricing.AddOne(c, mod)
	if err != nil {
		c.Errorf("SP AddOne failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.SPAddResponse{}

	return resp, nil
}

func (this SPSer) Modify(c context.ContextB, r *merchantmod.SPModifyRequest) (*merchantmod.SPModifyResponse, error) {
	update := &mdb.SpecificationsPricingUpdateDoc{
		Specifications: r.Specifications,
		Pricing:        r.Pricing,
		PicURL:         r.PicURL,
		ChoiceOpt:      r.ChoiceOpt,
	}
	err := mdb.SpecificationsPricing.UpdateById(c, r.Id, update)
	if err != nil {
		c.Errorf("SP UpdateById failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.SPModifyResponse{}

	return resp, nil
}

var ErrSPUsed = fmt.Errorf("SpecificationsPricing is used")

func (SPSer) Del(c context.ContextB, r *merchantmod.SPDelRequest) (*merchantmod.SPDelResponse, error) {
	err := mdb.SpecificationsPricing.DelByID(c, r.Id)
	if err != nil {
		c.Errorf("SP DelByID failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.SPDelResponse{}

	return resp, nil
}
