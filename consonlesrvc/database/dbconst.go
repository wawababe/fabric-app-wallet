package database

import (
	"os"
	"baas/app-wallet/chaincode/github.com/op/go-logging"
	"database/sql"
)

var dbLogger *logging.Logger = logging.MustGetLogger("database")
var db *sql.DB

const (
	DSN string = "root:101812@/app_wallet"
	DATETIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	ERROR_DB_NOT_CONNECTED = "failed connecting to database"
	ERROR_DB_PREPARED = "failed preparing sql statement"
	ERROR_DB_EXECUTE = "failed executing sql statement"
	ERROR_DB_QUERY = "failed quering sql statement"
)

func init(){
	bk := logging.NewLogBackend(os.Stdout, "", 0)
	var format = logging.MustStringFormatter(
		`%{color} %{time:2006-01-02T15:04:05} [%{module}] %{shortfunc} > %{level:.4s} %{id: 03x} %{color: reset}: %{message}`,
	)
	bkFormatter := logging.NewBackendFormatter(bk, format)
	bkLeveled := logging.AddModuleLevel(bkFormatter)
	bkLeveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(bkLeveled)
	db = new(sql.DB)

	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(err)
	}


}

func GetDB()*sql.DB {
	return db
}