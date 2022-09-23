/**
 * @Author: lzw5399
 * @Date: 2020/8/2 14:55
 * @Desc: 初始化数据库并产生数据库全局变量, 这里默认是postgres
 */
package initilize

import (
	"context"
	"hx/global"

	"github.com/qiniu/qmgo"
)

func initMongo() {
	ctx := context.Background()
	coreMongoConf := global.CoreMongo

	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: coreMongoConf.Uri})
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
