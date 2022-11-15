package initilize

import (
	"hx/global"
	"hx/util"
	"strings"

	iso639_3 "github.com/barbashov/iso639-3"
	gosonic "github.com/expectedsh/go-sonic/sonic"
)

func initSonic() {
	var err error
	cfg := global.Sonic
	ingester, err := gosonic.NewIngester(cfg.Host, cfg.Port, cfg.Password)
	if err != nil {
		panic(err)
	}

	global.SONIC_INGESTER_CH = make(chan *global.SonicBulkPushEvent)

	go func() {

		for e := range global.SONIC_INGESTER_CH {
			parallelRoutines := len(e.Records)
			if parallelRoutines > cfg.ParallelRoutines {
				parallelRoutines = cfg.ParallelRoutines
			}
			errs := ingester.BulkPush(e.Collection, e.Bucket, parallelRoutines, e.Records, getLang(e.Lang))
			global.DL_LOGGER.Debugf("sonic bulk push finish! trace: %s, errs: %s", e.Trace, util.MustMarshalToString(errs))
			if len(errs) != 0 {
				global.DL_LOGGER.Errorf("sonic bulk push failed! trace: %s, errs: %s", e.Trace, util.MustMarshalToString(errs))
			}
		}
	}()

	search, err := gosonic.NewSearch(cfg.Host, cfg.Port, cfg.Password)
	if err != nil {
		panic(err)
	}

	global.SONIC_SEARCH_CH = make(chan *global.SonicSearchEvent)

	go func() {
		for e := range global.SONIC_SEARCH_CH {
			result := &global.SonicSearcResult{}
			result.Results, result.Err = search.Query(e.Collection, e.Bucket, e.Terms, e.Limit, e.Offset, getLang(e.Lang))
			global.DL_LOGGER.Debugf("sonic search query finish! trace: %s, results: %s, errs: %v ", e.Trace, result.Results, result.Err)
			e.Result <- result
		}
	}()

	suggest, err := gosonic.NewSearch(cfg.Host, cfg.Port, cfg.Password)
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range global.SONIC_SUGGEST_CH {
			result := &global.SonicSearcResult{}
			result.Results, result.Err = suggest.Suggest(e.Collection, e.Bucket, e.Word, e.Limit)
			global.DL_LOGGER.Debugf("sonic search suggest finish! trace: %s, results: %s, errs: %v ", e.Trace, result.Results, result.Err)
			e.Result <- result
		}
	}()
}

func getLang(lang string) gosonic.Lang {
	lan := iso639_3.FromPart1Code(lang)
	if lan == nil {
		return gosonic.LangCmn
	}

	if strings.Contains(lan.Name, "Chinese") {
		return gosonic.LangCmn
	}

	return gosonic.Lang(lan.Part3)
}
