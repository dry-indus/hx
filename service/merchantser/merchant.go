package merchantser

import (
	"hx/global/context"
	"hx/mdb"
)

var Merchant MerchantSer

type MerchantSer struct{}

func (MerchantSer) FindByName(c context.ContextB, name string) (*mdb.MerchantMod, error) {
	merchant, err := mdb.Merchant.FindOneByName(c, name)
	if err != nil {
		c.Errorf("mdb.Merchant.FindOneByName failed! name: %v, err: %v", name, err)
		return nil, err
	}

	return merchant, nil
}
