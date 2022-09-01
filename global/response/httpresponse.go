/**
 * @Author: lzw5399
 * @Date: 2020/7/31 23:21
 * @Desc: 格式化错误响应
 */
package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	OK      Action = 1000
	Tip     Action = 2000
	Reload  Action = 3000
	Relogin Action = 4000
)

type Action int

func Response(c *gin.Context, code Action, msg string, data ...interface{}) {
	var dat interface{}
	dat = struct{}{}
	for _, v := range data {
		dat = v
		break
	}
	result(c, int(code), dat, msg)
}

func Success(c *gin.Context, data ...interface{}) {
	Response(c, OK, "Success", data...)
}

func Failed(c *gin.Context, data ...interface{}) {
	Response(c, OK, "Failed", data...)
}

type HTTPResponse struct {
	Status  int
	Message string
	Data    interface{}
}

func result(c *gin.Context, status int, data interface{}, msg string) {
	c.JSON(http.StatusOK, HTTPResponse{
		status,
		msg,
		data,
	})
}
