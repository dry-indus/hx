package response

import (
	"github.com/gin-gonic/gin"
)

var (
	OK   = 200
	Fail = 500

	InvalidParam        = func(c *gin.Context, msg ...string) Action { return Action{c, 1000, defaultStr(msg, "Invalid Param")} }
	InternalServerError = func(c *gin.Context, msg ...string) Action {
		return Action{c, 2000, defaultStr(msg, "Internal Server Error")}
	}
	Tip      = func(c *gin.Context, msg ...string) Action { return Action{c, 3000, defaultStr(msg, "Tip")} }
	Reload   = func(c *gin.Context, msg ...string) Action { return Action{c, 4000, defaultStr(msg, "Reload")} }
	Relogin  = func(c *gin.Context, msg ...string) Action { return Action{c, 5000, defaultStr(msg, "Relogin")} }
	Redirect = func(c *gin.Context, msg ...string) Action { return Action{c, 6000, defaultStr(msg, "Redirect")} }
)

func defaultStr(v []string, def string) string {
	if len(v) != 0 {
		return v[0]
	}
	return def
}

type Action struct {
	c    *gin.Context
	Code int
	Msg  string
}

func Success(c *gin.Context, data ...interface{}) {
	response(c, OK, OK, "success", data...)
}

func (this Action) Success(data ...interface{}) {
	response(this.c, OK, this.Code, this.Msg, data...)
}

func (this Action) Failed(data ...interface{}) {
	response(this.c, Fail, this.Code, this.Msg, data...)
}

type HTTPResponse struct {
	// | 业务响应码 | 响应信息 | 描述 |
	// | ---------- | -------- | ---- |
	// | 1000           | Invalid Param         | 无效参数 |
	// | 2000           | Internal Server Error         | 服务器内部错误 |
	// | 3000           | Tip         | 弹出信息  [查看示例](https://nutui.jd.com/#/zh-cn/component/notify)   |
	// | 4000           | Reload         | 重新加载页面 |
	// | 5000           | Relogin         | 重新登陆 |
	// | 6000           | Redirect         | 重定向  |
	Status int `json:"status" enums:"1000,2000,3000,4000,5000,6000"`
	// 信息
	Message string `json:"message"`
	// 数据
	Data interface{} `json:"data"`
}

func response(c *gin.Context, status, action int, msg string, data ...interface{}) {
	var dat interface{}
	dat = struct{}{}
	for _, v := range data {
		dat = v
		break
	}

	c.JSON(status, HTTPResponse{
		action,
		msg,
		dat,
	})
}
