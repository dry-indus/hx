package main

import (
	"flag"
	"hx/global"
	_ "hx/initilize"
	"hx/router"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	_commitID = flag.String("commitID", "无", "git commit id")
)

func main() {
	defer global.Close()
	go run()

	pwd, _ := exec.Command("pwd").Output()
	ls, _ := exec.Command("ls").Output()
	global.DL_LOGGER.Infof("server start! GIT: %v, pwd: %s, ls: %s", *_commitID, pwd, ls)
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	global.DL_LOGGER.Infof("close server")
}

func run() {
	router.Run()
}
