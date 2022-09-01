package global

import (
	"fmt"

	"github.com/qiniu/qmgo"
	"github.com/sirupsen/logrus"
)

var (
	DL_CLOSE    []func() error
	DL_CORE_MDB *qmgo.Database
	DL_LOGGER   *logrus.Logger
)

func Close() {
	for _, f := range DL_CLOSE {
		err := f()
		if err != nil {
			if DL_LOGGER != nil {
				DL_LOGGER.Error(err)
			} else {
				fmt.Println(err)
			}
		}
	}
}
