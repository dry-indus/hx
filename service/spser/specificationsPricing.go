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
		CreatedAt:      time.Now(),
	}
	id, err := mdb.SpecificationsPricing.AddOne(c, mod)
	if err != nil {
		c.Errorf("SP AddOne failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.SPAddResponse{Id: id}

	return resp, nil
}

func (this SPSer) Modify(c context.ContextB, r *merchantmod.SPModifyRequest) (*merchantmod.SPModifyResponse, error) {
	update := &mdb.SpecificationsPricingUpdateDoc{
		Specifications: r.Specifications,
		Pricing:        r.Pricing,
		PicURL:         r.PicURL,
	}
	err := mdb.SpecificationsPricing.UpdateById(c, r.Id, update)
	if err != nil {
		c.Errorf("SP UpdateById failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.SPModifyResponse{Id: r.Id}

	return resp, nil
}

var ErrOneMustBeRetained = fmt.Errorf("one must be retained!")

func (SPSer) Del(c context.ContextB, r *merchantmod.SPDelRequest) (*merchantmod.SPDelResponse, error) {
	sp, err := mdb.SpecificationsPricing.FindById(c, r.Id)
	if err != nil {
		c.Errorf("SP FindById failed! err: %v", err)
		return nil, err
	}

	count, err := mdb.SpecificationsPricing.CountByCommodityId(c, sp.CommodityId)
	if err != nil {
		c.Errorf("SP DelByID failed! err: %v", err)
		return nil, err
	}

	if count <= 1 {
		return nil, ErrOneMustBeRetained
	}

	err = mdb.SpecificationsPricing.DelByID(c, r.Id)
	if err != nil {
		c.Errorf("SP DelByID failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.SPDelResponse{}

	return resp, nil
}
