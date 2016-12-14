package database

import (
	"github.com/op/go-logging"
	"database/sql"
	"baas/app-wallet/consolesrvc/common"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var dbLogger *logging.Logger = common.NewLogger("database")//logging.MustGetLogger("database")
var db *sql.DB
const(
	DATETIME_FORMAT = "2006-01-02 15:04:05"
)

var (
	ERROR_DB_NOT_CONNECTED = "failed connecting to database"
	ERROR_DB_PREPARED = "failed preparing sql statement"
	ERROR_DB_EXECUTE = "failed executing sql statement"
	ERROR_DB_QUERY = "failed quering sql statement"
)

func InitDB(dsn string) *sql.DB{
	db = new(sql.DB)
	var err error
	if db, err = sql.Open("mysql", dsn); err != nil {
		dbLogger.Fatal(err)
	}
	return db
}

func GetDB()*sql.DB {
	if db == nil { // todo: could be removed out; Just for test
		viper.SetConfigName("wallet")
		viper.AddConfigPath("$GOPATH/src/baas/app-wallet/consolesrvc")
		viper.AddConfigPath(".")
		viper.ReadInConfig()
		return InitDB(viper.GetString("database.mysql.dsn"))
	}
	return db
}
