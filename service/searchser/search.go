package searchser

import (
	"hx/global"
	"hx/global/context"
	"hx/mdb"
	"hx/model/common"
	"hx/model/searchmod"
	"net/url"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	Search SearchSer
)

type SearchSer struct{}

func (this SearchSer) SearchStore(c context.ContextB, keywords string, page common.Page) (suggest []string, result []*searchmod.Store) {
	collection := "store"
	bucket := "default"

	suggestResultCh := make(chan *global.SonicSearcResult)
	suggestEvent := &global.SonicSuggestEvent{
		Collection: collection,
		Bucket:     bucket,
		Word:       keywords,
		Limit:      5,
		Result:     suggestResultCh,
		Trace:      c.Trace(),
	}

	searchResultCh := make(chan *global.SonicSearcResult)
	searchEvent := &global.SonicSearchEvent{
		Collection: collection,
		Bucket:     bucket,
		Terms:      keywords,
		Limit:      int(page.Limit()),
		Offset:     int(page.Skip()),
		Lang:       c.Lang(),
		Result:     searchResultCh,
		Trace:      c.Trace(),
	}

	var suggestResult *global.SonicSearcResult
	var searchResult *global.SonicSearcResult

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		global.SONIC_SUGGEST_CH <- suggestEvent
		suggestResult = <-suggestResultCh
		c.Debugf("suggest event finish")
	}()

	go func() {
		defer wg.Done()
		global.SONIC_SEARCH_CH <- searchEvent
		searchResult = <-searchResultCh
		c.Debugf("search event finish")
	}()

	suggest = suggestResult.Results

	if searchResult != nil {
		mids := []primitive.ObjectID{}
		for _, v := range searchResult.Results {
			if id, err := primitive.ObjectIDFromHex(v); err == nil {
				mids = append(mids, id)
			}
		}

		list, _ := mdb.Merchant.FindByIDs(c, mids)
		for _, v := range list {
			store := &searchmod.Store{
				StoreName: v.StoreName,
				Prtrait:   v.Prtrait,
				Category:  v.Category,
				Star:      v.Star,
				URL:       getStoreURL(v.Name),
			}
			result = append(result, store)
		}
	}

	return
}

func getStoreURL(merchantName string) string {
	u := url.URL{Scheme: "https", Host: global.Application.Domian, Path: merchantName}
	return u.String()
}
