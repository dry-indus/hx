package merchantctr

import (
	"hx/global/context"
	"hx/global/response"
	"hx/model/merchantmod"
	"hx/service/fileser"
	"hx/util"
)

var File FileCtr

type FileCtr struct{}

// @Tags        文件-上传
// @Summary     上传文件
// @Description 异步文件上传，该请求会立即返回当前上传状态。
// @Accept      multipart/form-data
// @Produce     json
// @param       file     formData file   true  "文件"
// @param       hoken    header   string false "hoken"
// @param       language header   string false "语言" default(zh-CN)
// @Security    Auth
// @Success     200 {object} response.HTTPResponse{data=merchantmod.FileUploadResponse} "成功"
// @Failure     500 {object} response.HTTPResponse
// @Router      /file/upload [post]
func (FileCtr) Upload(c context.MerchantContext) {
	form, _ := c.Gin().MultipartForm()
	files := form.File["upload[]"]
	if len(files) == 0 {
		if file, _ := c.Gin().FormFile("file"); file != nil {
			files = append(files, file)
		}
	}

	taskId := util.UUID().String()
	for i, file := range files {
		c.Infof("Upload %d/%d fileName: %s, size: %v", i+1, len(files), file.Filename, file.Size)
		f, err := file.Open()
		if err != nil {
			c.Warningf("Upload %d/%d Open file failed! fileName: %s", i+1, len(files), file.Filename)
			continue
		}
		fileser.File.MerchantUpload(c, taskId, file.Filename, file.Size, f)
	}

	resp := &merchantmod.FileUploadResponse{
		TaskId: taskId,
	}

	response.Success(c.Gin(), resp)

}

// @Tags        文件-上传
// @Summary     获取指定任务id的文件上传状态
// @Description 状态包括当前文件的总大小和已经上传的大小，以及上传的错误信息
// @Accept      json
// @Produce     json
// @Param       param    body   merchantmod.FileStatusRequest true  "参数"
// @param       hoken    header string                        false "hoken"
// @param       language header string                        false "语言" default(zh-CN)
// @Security    Auth
// @Success     200 {object} response.HTTPResponse{data=merchantmod.FileStatusResponse} "成功"
// @Failure     500 {object} response.HTTPResponse
// @Router      /file/status [post]
func (FileCtr) Status(c context.MerchantContext) {
	var r merchantmod.FileStatusRequest
	err := c.Gin().ShouldBindJSON(&r)
	if err != nil {
		response.InvalidParam(c.Gin()).Failed(err)
		return
	}

	status := fileser.File.UploadStatus(c, r.TaskId)

	resp := &merchantmod.FileStatusResponse{
		TaskId: r.TaskId,
		Status: status,
	}

	response.Success(c.Gin(), resp)
}
