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
