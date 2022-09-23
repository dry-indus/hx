package main

import (
	"fmt"
	"hx/global"
	"os"
	"os/signal"
	"syscall"

	_ "hx/docs"
	_ "hx/initilize"
	"hx/router"
)

var (
	CommitID string
)

// @title          HaiXian API
// @version        1.0
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url  http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7777
func main() {
	go run()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	global.DL_LOGGER.Infof("close server")
}

func run() {
	r := router.InitRouter()
	defer global.Close()

	port := global.Application.Port
	global.DL_LOGGER.Infof("server listening at port: %v, GIT: %v", port, CommitID)
	fmt.Println()
	_ = r.Run(":" + port)
}
