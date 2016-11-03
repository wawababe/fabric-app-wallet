package common

import (
	"os"
	"github.com/op/go-logging"
)

func NewLogger(module string)(*logging.Logger){
	bk := logging.NewLogBackend(os.Stdout, "", 0)
	var format = logging.MustStringFormatter(
		`%{color} %{time:2006-01-02 15:04:05} [%{module}] %{shortfunc} > %{level:.4s} %{id:03x} %{color:reset}: %{message}`,
	)
	bkFormatter := logging.NewBackendFormatter(bk, format)
	bkLeveled := logging.AddModuleLevel(bkFormatter)
	bkLeveled.SetLevel(logging.DEBUG, module)
	logging.SetBackend(bkLeveled)
	return logging.MustGetLogger(module)
}

