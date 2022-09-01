package initilize

import (
	"flag"
	"fmt"
	"hx/global"
	"hx/util"

	"github.com/jinzhu/configor"
	"github.com/shima-park/agollo"
)

func init() {
	initApollo()
	initLog()
	initMongo()
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

	ago, err := agollo.New(
		conf.IP,
		conf.AppID,
		agollo.Cluster(conf.Cluster),
		agollo.PreloadNamespaces(global.Namespaces...),
		agollo.AccessKey(conf.AccessKey),
		agollo.BackupFile(conf.BackupFile),
	)
	if err != nil {
		panic(err)
	}

	for n, p := range global.Namespacem {
		ns := ago.GetNameSpace(n)
		s := Decode(ns, p)
		fmt.Printf("%v using config: %v\n", n, s)
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
				s := Decode(resp.NewValue, ptr)
				fmt.Printf("%v using config: %v\n", resp.Namespace, s)
			}
		}
	}()

}

func Decode(c agollo.Configurations, ptr interface{}) string {
	b, err := util.JSON.Marshal(c)
	if err != nil {
		return fmt.Sprint(err)
	}
	err = util.JSON.Unmarshal(b, &ptr)
	if err != nil {
		return fmt.Sprint(err)
	}

	fmt.Printf("--->Decode: %v, %v", string(b), ptr)
	return string(b)
}
