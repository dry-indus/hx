package merchantmod

type FileUploadResponse struct {
	// TaskId 上传任务ID
	TaskId string `json:"taskId"`
	// 文件和上传状态的映射
	Status map[string]*UploadStatus `json:"status"`
}

type UploadStatus struct {
	// TaskId 上传任务ID
	TaskId string `json:"taskId"`
	// 文件名
	FileName string `json:"fileName"`
	// 上传后获取的文件URL
	URL string `json:"url"`
	// 已经上传的Bytes
	ConsumedBytes int64 `json:"consumedBytes"`
	// 需要上传的总Bytes
	TotalBytes int64 `json:"totalBytes"`
	// 已经读写的Bytes
	RwBytes int64 `json:"rwBytes"`
	// 上传的错误信息
	Err string `json:"err"`
	// 状态更新的时间，UnixNano时间戳
	At int64 `json:"at"`
	// true: 上传完成
	IsCompleted bool `json:"isCompleted"`
}

type FileStatusRequest struct {
	// TaskId 上传任务ID
	TaskId string `json:"taskId" binding:"required" validate:"required"`
}

type FileStatusResponse struct {
	// TaskId 上传任务ID
	TaskId string `json:"taskId"`
	// 文件和上传状态的映射
	Status map[string]*UploadStatus `json:"status"`
}
