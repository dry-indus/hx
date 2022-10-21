package initilize

import (
	"hx/global"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func initOSS() {
	c := global.Oss
	client, err := oss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret,
		oss.Timeout(c.ConnectTimeoutSec, c.ReadWriteTimeout),
		oss.EnableMD5(true),
		oss.EnableCRC(true),
	)
	if err != nil {
		panic(err)
	}

	// Get bucket
	bucket, err := client.Bucket(c.BucketName)
	if err != nil {
		panic(err)
	}
	global.DL_OSS_BUCKET = bucket
}
