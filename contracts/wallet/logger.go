package wallet

import (
	"github.com/op/go-logging"
	"os"
)

var wtLogger *logging.Logger = logging.MustGetLogger("wallet")

func init() {
	bk := logging.NewLogBackend(os.Stdout, "", 0)
	var fmt logging.Formatter = logging.MustStringFormatter(
		`%{color} %{time:2006-01-02T15:04:05} [%{module}] %{shortfunc} > %{level:.4s} %{id:03x} %{color:reset}: %{message} `,
	)
	var bkFormatter logging.Backend = logging.NewBackendFormatter(bk, fmt)
	var bkLeveled logging.LeveledBackend = logging.AddModuleLevel(bkFormatter)
	bkLeveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(bkLeveled)
}
