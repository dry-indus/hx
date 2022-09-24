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
	// Status 1000: Invalid Param,
	// Status 2000: Internal Server Error
	// Status 3000: Tip
	// Status 4000: Reload
	// Status 5000: Relogin
	// Status 6000: Redirect
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
