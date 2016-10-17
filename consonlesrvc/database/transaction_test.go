package database

import (
	"testing"
	"database/sql"
	"github.com/hyperledger/fabric/core/util"
)

func TestAddTransaction(t *testing.T) {
	var db *sql.DB
	var err error
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal("Failed opening database")
	}

	var testNum = 2

	var trans = make([]Transaction, testNum)
	for i, _ := range trans {
		trans[i].TXUUID = util.GenerateUUID()
		trans[i].PayerUUID = util.GenerateUUID()
		trans[i].PayeeUUID = util.GenerateUUID()
		trans[i].Amount = 30
		dbLogger.Debugf("transaction: %#v", trans[i])
	}
	trans[0].PayerUUID = "5cdb617c-2712-480a-a02b-facd8c86e579"
	trans[1].TXUUID = trans[0].TXUUID

	var tests = []struct{
		newline bool
		sep string
		arg *Transaction
		want int64
	}{
		{false, " ", &trans[0], 1},
		{true, " ", &trans[1], 0},
	}

	for i, testitem := range tests {
		rowsAff, _ := AddTransaction(db, testitem.arg)
		if rowsAff != testitem.want {
			t.Errorf("Test #%d: Add transaction %#v, affected rows = %d, but want %d", i, testitem.arg, rowsAff, testitem.want)
		}
	}
}

func TestGetTransactionByTransUUID(t *testing.T) {
	var db *sql.DB
	var err error
	var us *Transaction
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
	}

	var txuuid string
	//useruuid = "5cdb617c-2712-480a-a02b-facd8c86e579"
	txuuid = "7382e24f-bd1c-4e41-ab86-e819e823b75b"
	us, err = GetTransaction(db, txuuid)
	if us == nil || err != nil{
		t.Errorf("Failed retrieving transaction")
	}
	dbLogger.Debugf("Get transaction: %#v", *us)
}

func TestGetTransactionsByPayeruuid(t *testing.T) {
	var db *sql.DB
	var err error
	var txs []*Transaction
	if db, err = sql.Open("mysql", DSN); err != nil {
		dbLogger.Fatal(ERROR_DB_NOT_CONNECTED)
	}

	var useruuid string = "5cdb617c-2712-480a-a02b-facd8c86e579"
	txs, err = GetTransactionsByPayeruuid(db, useruuid)
	if err != nil {
		t.Errorf("Failed retrieving user accounts by useruuid %s: %v", useruuid, err)
	}
	for i, txitem := range txs {
		dbLogger.Debugf("Accounts #%d: %v", i, *txitem)
	}
}
