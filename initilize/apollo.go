package initilize

import (
	"flag"
	"fmt"
	"hx/global"
	"hx/service/tgser"
	"hx/util"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/shima-park/agollo"
)

func init() {
	initTime()
	initApollo()
	initLog()
	initRedis()
	initMongo()
	initSession()
	initTgBot(tgser.Tg)
	initOSS()
	initSonic()
}

var _config = flag.String("config", "./config/dev_settings.yaml", "start-up config file")

func initApollo() {
	flag.Parse()

	var conf struct {
		AppID      string
		Cluster    string
		IP         string
		AccessKey  string
		BackupFile string
	}

	configFile := *_config
	fmt.Println("use config file: ", configFile)

	err := configor.Load(&conf, configFile)
	if err != nil {
		panic("resolve settings failed...")
	}

	global.AppName = conf.AppID
	global.ENV = strings.ToUpper(conf.Cluster)

	ago, err := agollo.New(
		conf.IP,
		conf.AppID,
		agollo.Cluster(conf.Cluster),
		agollo.PreloadNamespaces(global.Namespaces...),
		agollo.AccessKey(conf.AccessKey),
		// agollo.BackupFile(conf.BackupFile),
	)
	if err != nil {
		panic(err)
	}

	for n, p := range global.Namespacem {
		var err error
		ns := ago.GetNameSpace(n)
		if strings.Contains(n, ".json") {
			err = Decode(ns["content"], p)
		} else {
			err = Decode(ns, p)
		}

		fmt.Printf("%v using config: %v, err: %v\n", n, ns, err)
	}

	go func() {
		//启动apollo长轮询
		errorCh := ago.Start()
		defer ago.Stop()

		watchCh := ago.Watch()

		for {
			select {
			case err := <-errorCh:
				fmt.Println("Error:", err)
			case resp := <-watchCh:
				ptr := global.Namespacem[resp.Namespace]
				if ptr == nil {
					fmt.Printf("don't use namespace: %s...\n", resp.Namespace)
					continue
				}

				var err error
				nv := resp.NewValue
				if strings.Contains(resp.Namespace, ".json") {
					err = Decode(nv["content"], ptr)
				} else {
					err = Decode(nv, ptr)
				}
				fmt.Printf("%v using config: %v, err: %v\n", resp.Namespace, nv, err)
			}
		}
	}()

}

func Decode(o interface{}, ptr interface{}) error {
	b, err := util.JSON.Marshal(o)
	if err != nil {
		return err
	}

	err = util.JSON.Unmarshal(b, &ptr)
	if err != nil {
		return err
	}

	return nil
}
