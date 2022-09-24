package tagser

import (
	"fmt"
	"hx/global/context"
	"hx/mdb"
	"hx/model/merchantmod"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Tag TagSer
)

type TagSer struct {
}

func (this TagSer) Add(c context.ContextB, merchantId primitive.ObjectID, r *merchantmod.TagAddRequest) (*merchantmod.TagAddResponse, error) {
	mod := &mdb.TagMod{
		Name:       r.Name,
		MerchantId: merchantId,
		CreatedAt:  time.Now(),
	}

	id, err := mdb.Tag.AddOne(c, mod)
	if err != nil {
		c.Errorf("Tag AddOne failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.TagAddResponse{Id: id}

	return resp, nil
}

var ErrTagUsed = fmt.Errorf("tag is used")

func (TagSer) Del(c context.ContextB, r *merchantmod.TagDelRequest) (*merchantmod.TagDelResponse, error) {
	if list, _ := mdb.Commodity.FindByTagIds(c, []primitive.ObjectID{r.Id}); len(list) > 0 {
		return nil, ErrTagUsed
	}

	err := mdb.Tag.DelByID(c, r.Id)
	if err != nil {
		c.Errorf("Tag Del failed! err: %v", err)
		return nil, err
	}

	resp := &merchantmod.TagDelResponse{}

	return resp, nil
}
