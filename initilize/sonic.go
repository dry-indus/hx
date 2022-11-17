package initilize

import (
	"hx/global"
	"hx/util"
	"strings"

	iso639_3 "github.com/barbashov/iso639-3"
	gosonic "github.com/expectedsh/go-sonic/sonic"
)

func initSonic() {
	cfg := global.Sonic
	global.SONIC_INGESTER_CH = make(chan *global.SonicIngestEvent)

	go func() {
		for e := range global.SONIC_INGESTER_CH {
			parallelRoutines := len(e.Records)
			if parallelRoutines > cfg.ParallelRoutines {
				parallelRoutines = cfg.ParallelRoutines
			}
			lang := getLang(e.Lang)

			ingester, err := gosonic.NewIngester(cfg.Host, cfg.Port, cfg.Password)
			if err != nil {
				global.DL_LOGGER.Errorf("sonic new ingester failed! trace: %s, errs: %s", e.Trace, util.MustMarshalToString(err))
				continue
			}

			switch e.Method {
			case global.BulkPush:
				errs := ingester.BulkPush(e.Collection, e.Bucket, parallelRoutines, e.Records, lang)
				global.DL_LOGGER.Debugf("sonic bulk push finish! trace: %s, records:%s, lang: %s, errs: %v", e.Trace, util.MustMarshalToString(e.Records), lang, errs)
				if len(errs) != 0 {
					global.DL_LOGGER.Errorf("sonic bulk push failed! trace: %s, records:%s, lang: %s, errs: %v", e.Trace, util.MustMarshalToString(e.Records), lang, errs)
				}
			case global.BulkPop:
				errs := ingester.BulkPop(e.Collection, e.Bucket, parallelRoutines, e.Records)
				global.DL_LOGGER.Debugf("sonic bulk pop finish! trace: %s, records:%s, lang: %s, errs: %v", e.Trace, util.MustMarshalToString(e.Records), lang, errs)
				if len(errs) != 0 {
					global.DL_LOGGER.Errorf("sonic bulk pop failed! trace: %s, records:%s, lang: %s, errs: %v", e.Trace, util.MustMarshalToString(e.Records), lang, errs)
				}
			case global.Push:
				err := ingester.Push(e.Collection, e.Bucket, e.Records[0].Object, e.Records[0].Text, lang)
				global.DL_LOGGER.Debugf("sonic push finish! trace: %s, records:%s, lang: %s, err: %v", e.Trace, util.MustMarshalToString(e.Records), lang, err)
				if err != nil {
					global.DL_LOGGER.Errorf("sonic push failed! trace: %s, records:%s, lang: %s, err: %s", e.Trace, util.MustMarshalToString(e.Records), lang, err)
				}
			case global.Pop:
				err := ingester.Pop(e.Collection, e.Bucket, e.Records[0].Object, e.Records[0].Text)
				global.DL_LOGGER.Debugf("sonic pop finish! trace: %s, records:%s, lang: %s, err: %v", e.Trace, util.MustMarshalToString(e.Records), lang, err)
				if err != nil {
					global.DL_LOGGER.Errorf("sonic pop failed! trace: %s, records:%s, lang: %s, err: %s", e.Trace, util.MustMarshalToString(e.Records), lang, err)
				}
			default:
				global.DL_LOGGER.Errorf("not fond sonic method: %v! trace: %s", e.Method, e.Trace)
			}

			ingester.Quit()
		}
	}()

	global.SONIC_SEARCH_CH = make(chan *global.SonicSearchEvent)

	go func() {
		for e := range global.SONIC_SEARCH_CH {
			search, err := gosonic.NewSearch(cfg.Host, cfg.Port, cfg.Password)
			if err != nil {
				global.DL_LOGGER.Errorf("sonic new search failed! trace: %s, errs: %s", e.Trace, util.MustMarshalToString(err))
				continue
			}
			result := &global.SonicSearcResult{}
			result.Results, result.Err = search.Query(e.Collection, e.Bucket, e.Terms, e.Limit, e.Offset, getLang(e.Lang))
			global.DL_LOGGER.Debugf("sonic search query finish! trace: %s, results: %s, errs: %v ", e.Trace, result.Results, result.Err)
			e.Result <- result

			search.Quit()
		}
	}()

	global.SONIC_SUGGEST_CH = make(chan *global.SonicSuggestEvent)

	go func() {
		for e := range global.SONIC_SUGGEST_CH {
			suggest, err := gosonic.NewSearch(cfg.Host, cfg.Port, cfg.Password)
			if err != nil {
				global.DL_LOGGER.Errorf("sonic new suggest failed! trace: %s, errs: %s", e.Trace, util.MustMarshalToString(err))
				continue
			}

			result := &global.SonicSearcResult{}
			result.Results, result.Err = suggest.Suggest(e.Collection, e.Bucket, e.Word, e.Limit)
			global.DL_LOGGER.Debugf("sonic search suggest finish! trace: %s, results: %s, errs: %v ", e.Trace, result.Results, result.Err)
			e.Result <- result

			suggest.Quit()
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
