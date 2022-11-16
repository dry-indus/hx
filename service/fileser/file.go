package fileser

import (
	ctx "context"
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/model/merchantmod"
	"hx/util"
	"io"
	"net/url"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var (
	File = FileServer{}
)

type FileServer struct {
}

func (f *FileServer) MerchantUpload(c context.MerchantContext, taskId, fileName string, length int64, reader io.Reader) {
	f.upload(c, taskId, global.MERCHANT, c.Merchant().Name, fileName, length, reader)
}

func (f *FileServer) upload(c context.ContextB, taskId, role, userName, fileName string, length int64, reader io.Reader) {
	objectKey := GenObjectKey(role, userName, fileName)

	go func() {
		err := global.DL_OSS_BUCKET.PutObject(objectKey, reader,
			oss.Routines(10),
			oss.EnableMd5(),
			oss.ContentLength(length),
			oss.Progress(&OssProgressListener{
				TaskId:     taskId,
				FileName:   fileName,
				ObjectKey:  objectKey,
				TotalBytes: length,
				C:          c,
			}))

		if err != nil {
			c.Errorf("PutObject failed! objectKey: %s, err: %s", objectKey, err)
			return
		}
	}()

	return
}

func (FileServer) UploadStatus(c context.ContextB, taskId string) (status map[string]*merchantmod.UploadStatus) {
	hKey := fmt.Sprintf(global.OSS_PROGRESS_HASH_FMT, taskId)
	val := global.DL_CORE_REDIS.HGetAll(c, hKey).Val()
	if len(val) == 0 {
		return
	}

	status = map[string]*merchantmod.UploadStatus{}
	for k, v := range val {
		s := merchantmod.UploadStatus{}
		util.JSON.UnmarshalFromString(v, &s)
		status[k] = &s
	}

	return
}

func GenObjectKey(role, userName, fileName string) string {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	yymmdd := time.Now().In(cstZone).Format("2006-01-02")
	return fmt.Sprintf("%s/%s/%s/%s", yymmdd, role, userName, fileName)
}

// 定义进度条监听器
type OssProgressListener struct {
	TaskId     string
	FileName   string
	ObjectKey  string
	TotalBytes int64
	C          context.ContextB
}

func (o *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		setUploadStatus(o.C, &merchantmod.UploadStatus{
			TaskId:     o.TaskId,
			FileName:   o.FileName,
			TotalBytes: o.TotalBytes,
		})
	case oss.TransferDataEvent:
		status := getUploadStatus(o.C, o.TaskId, o.FileName)
		status.ConsumedBytes = event.ConsumedBytes
		status.RwBytes = event.RwBytes
		setUploadStatus(o.C, status)
	case oss.TransferCompletedEvent:
		status := getUploadStatus(o.C, o.TaskId, o.FileName)
		uri := url.URL{Scheme: global.Oss.UrlScheme, Host: global.Oss.UrlHost, Path: o.ObjectKey}
		status.URL = uri.String()
		status.IsCompleted = true
		setUploadStatus(o.C, status)
	case oss.TransferFailedEvent:
		status := getUploadStatus(o.C, o.TaskId, o.FileName)
		status.Err = "upload failed!"
		setUploadStatus(o.C, status)
	}

	o.C.Debugf("Progress Changed! event: %s", util.MustMarshalToString(event))
}

func getUploadStatus(c context.ContextB, taskId, fileName string) *merchantmod.UploadStatus {
	hKey := fmt.Sprintf(global.OSS_PROGRESS_HASH_FMT, taskId)
	val := global.DL_CORE_REDIS.HGet(ctx.TODO(), hKey, fileName).Val()
	status := &merchantmod.UploadStatus{}
	util.JSON.UnmarshalFromString(val, &status)
	return status
}

func setUploadStatus(c context.ContextB, status *merchantmod.UploadStatus) {
	status.At = time.Now().UnixNano()
	hKey := fmt.Sprintf(global.OSS_PROGRESS_HASH_FMT, status.TaskId)
	val := util.MustMarshalToString(status)
	ttl := 10 * time.Minute
	todo := ctx.TODO()
	global.DL_CORE_REDIS.HSet(todo, hKey, status.FileName, val)
	global.DL_CORE_REDIS.Expire(todo, hKey, ttl)
}
