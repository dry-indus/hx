/**
 * @Author: lzw5399
 * @Date: 2020/8/2 14:55
 * @Desc: 初始化数据库并产生数据库全局变量, 这里默认是postgres
 */
package initilize

import (
	"context"
	"hx/global"
	"hx/util"
	"reflect"

	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

func initMongo() {
	ctx := context.Background()
	coreMongoConf := global.CoreMongo

	opt := opts.Client().
		SetAppName(global.AppName).
		SetRegistry(bson.NewRegistryBuilder().
			RegisterDecoder(reflect.TypeOf(decimal.Decimal{}), util.Decimal{}).
			RegisterEncoder(reflect.TypeOf(decimal.Decimal{}), util.Decimal{}).
			Build())

	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: coreMongoConf.Uri}, options.ClientOptions{ClientOptions: opt})
	if err != nil {
		panic(err)
	}

	global.DL_CLOSE = append(global.DL_CLOSE, func() error {
		err := client.Close(ctx)
		global.DL_LOGGER.Infof("mongo is close! err: %v", err)
		return err
	})

	global.DL_CORE_CLI = client
	global.DL_CORE_MDB = client.Database("m_4000")
}
