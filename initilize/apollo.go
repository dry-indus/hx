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
		ns := ago.GetNameSpace(n)
		err := Decode(n, ns, p)
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

				err := Decode(resp.Namespace, resp.NewValue, ptr)
				fmt.Printf("%v using config: %v, err: %v\n", resp.Namespace, resp.NewValue, err)
			}
		}
	}()

}

func Decode(ns string, cfg agollo.Configurations, ptr interface{}) error {
	if strings.Contains(ns, ".json") {
		return util.JSON.UnmarshalFromString(cfg["content"].(string), &ptr)
	}

	b, err := util.JSON.Marshal(cfg)
	if err != nil {
		return err
	}

	err = util.JSON.Unmarshal(b, &ptr)
	if err != nil {
		return err
	}

	return nil
}
