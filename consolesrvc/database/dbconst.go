package database

import (
	"github.com/op/go-logging"
	"database/sql"
	"baas/app-wallet/consonlesrvc/common"
	_ "github.com/go-sql-driver/mysql"
)

var dbLogger *logging.Logger = common.NewLogger("database")//logging.MustGetLogger("database")
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
	db = new(sql.DB)

	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(err)
	}


}

func GetDB()*sql.DB {
	return db
}